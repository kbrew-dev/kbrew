package raw

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/tabwriter"

	osappsv1 "github.com/openshift/api/apps/v1"
	osversioned "github.com/openshift/client-go/apps/clientset/versioned"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	// Load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/engine"
	"github.com/kbrew-dev/kbrew/pkg/kube"
	"github.com/kbrew-dev/kbrew/pkg/log"
	"github.com/kbrew-dev/kbrew/pkg/yaml"
)

type method string

const (
	install   method = "apply"
	uninstall method = "delete"
	// upgrade   method = "apply" // unused

	evalExpression = `select(.kind  == "%s" and .metadata.name == "%s").%s |= %v`
)

var yamlDelimiter = regexp.MustCompile(`(?m)^---$`)

// App represents K8s app defined with plain YAML manifests
type App struct {
	app      config.App
	log      *log.Logger
	kubeCli  kubernetes.Interface
	osAppCli osversioned.Interface
}

// New returns new instance of raw App
func New(c config.App, log *log.Logger) (*App, error) {
	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to load Kubernetes config")
	}

	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create Kubernetes client")
	}
	osCli, err := osversioned.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create OpenShift client")
	}

	rApp := &App{
		app:      c,
		log:      log,
		kubeCli:  cli,
		osAppCli: osCli,
	}
	return rApp, nil
}

func (r *App) resolveArgs() error {
	//TODO: user global singleton kubeconfig in all modules
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		return errors.Wrapf(err, "Failed to load Kubernetes config")
	}

	e := engine.NewEngine(config)

	// TODO(@sahil.lakhwani): Parse only templated arguments
	if len(r.app.Args) != 0 {
		for arg, value := range r.app.Args {
			v, err := e.Render(fmt.Sprintf("%v", value))
			if err != nil {
				return err
			}
			r.app.Args[arg] = v
		}
	}
	return nil
}

// Install installs the app specified by name, version and namespace.
func (r *App) Install(ctx context.Context, name, namespace, version string, options map[string]string) error {
	manifest, err := getManifest(r.app.Repository.URL)
	if err != nil {
		return err
	}

	if err := r.resolveArgs(); err != nil {
		return err
	}

	patchedManifest, err := patchManifest(manifest, r.app.Args)
	if err != nil {
		return err
	}

	if err := kube.CreateNamespace(ctx, r.kubeCli, namespace); err != nil && !k8sErrors.IsAlreadyExists(err) {
		return err
	}

	// TODO(@prasad): Use go sdks
	out, err := kubectlCommand(ctx, install, name, namespace, patchedManifest)
	if err != nil {
		r.log.Debug(out)
		return err
	}
	r.log.Debug(out)
	r.log.Debugf("Waiting for components to be ready for %s\n", name)
	return r.waitForReady(ctx, namespace)
}

// Uninstall uninstalls the app specified by name and namespace.
func (r *App) Uninstall(ctx context.Context, name, namespace string) error {
	// TODO(@prasad): Use go sdks
	out, err := kubectlCommand(ctx, uninstall, name, namespace, r.app.Repository.URL)
	r.log.Debug(out)
	return err
}

// Search searches the app specified by name.
func (r *App) Search(ctx context.Context, name string) (string, error) {
	return printList(r.app), nil
}

func kubectlCommand(ctx context.Context, m method, name, namespace, manifest string) (string, error) {
	var c *exec.Cmd
	switch m {
	case install:
		c = exec.CommandContext(ctx, "kubectl", string(m), "-f", "-")
		// Pass the manifest on STDIN
		c.Stdin = strings.NewReader(manifest)
	default:
		c = exec.CommandContext(ctx, "kubectl", string(m), "-f", manifest)
	}

	if namespace != "" {
		c.Args = append(c.Args, "--namespace", namespace)
	}
	c.Stderr = os.Stderr
	output, err := c.Output()
	return string(output), err
}

// Workloads returns K8s workload object reference list for the raw app
func (r *App) Workloads(ctx context.Context, namespace string) ([]corev1.ObjectReference, error) {
	resp, err := http.Get(r.app.Repository.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read resource manifest from URL")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read resource manifest from URL")
	}
	return ParseManifestYAML(string(data), namespace)
}

func (r *App) waitForReady(ctx context.Context, namespace string) error {
	workloads, err := r.Workloads(ctx, namespace)
	if err != nil {
		return err
	}
	for _, wRef := range workloads {
		switch wRef.Kind {
		case "Pod":
			if err := kube.WaitForPodReady(ctx, r.kubeCli, wRef.Namespace, wRef.Name); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Pod not in ready state. Namespace: %s, Name: %s", wRef.Namespace, wRef.Name))
			}

		case "Deployment":
			if err := kube.WaitForDeploymentReady(ctx, r.kubeCli, wRef.Namespace, wRef.Name); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Deployment not in ready state. Namespace: %s, Name: %s", wRef.Namespace, wRef.Name))
			}

		case "StatefulSet":
			if err := kube.WaitForStatefulSetReady(ctx, r.kubeCli, wRef.Namespace, wRef.Name); err != nil {
				return errors.Wrap(err, fmt.Sprintf("StatefulSet not in ready state. Namespace: %s, Name: %s", wRef.Namespace, wRef.Name))
			}

		case "DeploymentConfig":
			if err := kube.WaitForDeploymentConfigReady(ctx, r.osAppCli, r.kubeCli, wRef.Namespace, wRef.Name); err != nil {
				return errors.Wrap(err, fmt.Sprintf("DeploymentConfig not in ready state. Namespace: %s, Name: %s", wRef.Namespace, wRef.Name))
			}
		}
	}
	return nil
}

// ParseManifestYAML splits yaml manifests with multiple K8s object specs and returns list of workload object references
func ParseManifestYAML(manifest, namespace string) ([]corev1.ObjectReference, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	objRefs := []corev1.ObjectReference{}
	for _, spec := range yamlDelimiter.Split(manifest, -1) {
		if len(spec) == 0 {
			continue
		}
		obj, _, err := decode([]byte(spec), nil, nil)
		if err != nil {
			continue
		}

		// Set default namespace if empty
		if namespace == "" {
			namespace = "default"
		}

		switch w := obj.(type) {
		case *corev1.Pod:
			if w.GetNamespace() != "" {
				namespace = w.GetNamespace()
			}
			objRefs = append(objRefs, corev1.ObjectReference{Name: w.GetName(), Namespace: namespace, Kind: "Pod"})

		case *appsv1.Deployment:
			if w.GetNamespace() != "" {
				namespace = w.GetNamespace()
			}
			objRefs = append(objRefs, corev1.ObjectReference{Name: w.GetName(), Namespace: namespace, Kind: "Deployment"})

		case *appsv1.StatefulSet:
			if w.GetNamespace() != "" {
				namespace = w.GetNamespace()
			}
			objRefs = append(objRefs, corev1.ObjectReference{Name: w.GetName(), Namespace: namespace, Kind: "StatefulSet"})

		case *osappsv1.DeploymentConfig:
			if w.GetNamespace() != "" {
				namespace = w.GetNamespace()
			}
			objRefs = append(objRefs, corev1.ObjectReference{Name: w.GetName(), Namespace: namespace, Kind: "DeploymentConfig"})
		}

	}
	return objRefs, nil
}

func printList(app config.App) string {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tVERSION\tTYPE")
	fmt.Fprintf(w, "%s\t%s\t%s", app.Name, app.Version, app.Repository.Type)
	w.Flush()
	return b.String()
}

func patchManifest(manifest string, patches map[string]interface{}) (string, error) {
	e := yaml.NewEvaluator()
	patchedManifest := manifest
	var err error
	for _, expression := range createExpressions(patches) {
		patchedManifest, err = e.Eval(patchedManifest, expression)
		if err != nil {
			return "", err
		}
	}
	return patchedManifest, nil
}

func getManifest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "Error fetching from app URL")
	}

	defer resp.Body.Close()

	manifest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Error fetching from app URL")
	}
	return string(manifest), nil
}

func createExpressions(patches map[string]interface{}) []string {
	var expressions []string

	for k, v := range patches {
		// Type assertion is necessary for yq, strings without quotes result in error
		switch v.(type) {
		case string:
			v = fmt.Sprintf("\"%s\"", v)
		default:
		}

		keys := strings.Split(k, ".")
		// keys[0] - kind
		// keys[1] - metadata.name
		// keys[2:] - path of the field
		e := fmt.Sprintf(evalExpression, keys[0], keys[1], strings.Join(keys[2:], "."), v)
		expressions = append(expressions, e)
	}
	return expressions
}
