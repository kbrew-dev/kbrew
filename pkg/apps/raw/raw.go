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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"

	// Load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kbrew-dev/kbrew/pkg/config"
	"github.com/kbrew-dev/kbrew/pkg/kube"
	"github.com/kbrew-dev/kbrew/pkg/yaml"
)

type method string

const (
	install   method = "apply"
	uninstall method = "delete"
	upgrade   method = "apply"

	evalExpression = `select(.kind  == "%s" and .metadata.name == "%s").%s |= %v`
)

var yamlDelimiter = regexp.MustCompile(`(?m)^---$`)

// App represents K8s app defined with plain YAML manifests
type App struct {
	App      config.App
	KubeCli  kubernetes.Interface
	OSAppCli osversioned.Interface
}

// New returns new instance of raw App
func New(c config.App) (*App, error) {
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
		App:      c,
		KubeCli:  cli,
		OSAppCli: osCli,
	}
	return rApp, nil
}

// Install installs the app specified by name, version and namespace.
func (r *App) Install(ctx context.Context, name, namespace, version string, options map[string]string) error {
	fmt.Printf("Installing raw app %s/%s\n", r.App.Repository.Name, name)

	manifest, err := getManifest(r.App.Repository.URL)
	if err != nil {
		return err
	}

	patchedManifest, err := patchManifest(manifest, r.App.Args)
	if err != nil {
		return err
	}

	// TODO(@prasad): Use go sdks
	if err := kubectlCommand(install, name, namespace, patchedManifest); err != nil {
		return err
	}
	return r.waitForReady(ctx, namespace)
}

// Uninstall uninstalls the app specified by name and namespace.
func (r *App) Uninstall(ctx context.Context, name, namespace string) error {
	fmt.Printf("Unistalling raw app %s\n", name)
	// TODO(@prasad): Use go sdks
	return kubectlCommand(uninstall, name, namespace, r.App.Repository.URL)
}

// Search searches the app specified by name.
func (r *App) Search(ctx context.Context, name string) (string, error) {
	return printList(r.App), nil
}

func kubectlCommand(m method, name, namespace, manifest string) error {
	var c *exec.Cmd
	switch m {
	case install:
		c = exec.Command("kubectl", string(m), "-f", "-")
		// Pass the manifest on STDIN
		c.Stdin = strings.NewReader(manifest)
	default:
		c = exec.Command("kubectl", string(m), "-f", manifest)
	}

	if namespace != "" {
		c.Args = append(c.Args, "--namespace", namespace)
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (r *App) waitForReady(ctx context.Context, namespace string) error {
	resp, err := http.Get(r.App.Repository.URL)
	if err != nil {
		return errors.Wrap(err, "Failed to read resource manifest from URL")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read resource manifest from URL")
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	for _, spec := range yamlDelimiter.Split(string(data), -1) {
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
			if err := kube.WaitForPodReady(ctx, r.KubeCli, namespace, w.GetName()); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Pod not in ready state. Namespace: %s, Name: %s", namespace, w.GetName()))
			}

		case *appsv1.Deployment:
			if w.GetNamespace() != "" {
				namespace = w.GetNamespace()
			}
			if err := kube.WaitForDeploymentReady(ctx, r.KubeCli, namespace, w.GetName()); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Deployment not in ready state. Namespace: %s, Name: %s", namespace, w.GetName()))
			}

		case *appsv1.StatefulSet:
			if w.GetNamespace() != "" {
				namespace = w.GetNamespace()
			}
			if err := kube.WaitForStatefulSetReady(ctx, r.KubeCli, namespace, w.GetName()); err != nil {
				return errors.Wrap(err, fmt.Sprintf("StatefulSet not in ready state. Namespace: %s, Name: %s", namespace, w.GetName()))
			}

		case *osappsv1.DeploymentConfig:
			if w.GetNamespace() != "" {
				namespace = w.GetNamespace()
			}
			if err := kube.WaitForDeploymentConfigReady(ctx, r.OSAppCli, r.KubeCli, namespace, w.GetName()); err != nil {
				return errors.Wrap(err, fmt.Sprintf("DeploymentConfig not in ready state. Namespace: %s, Name: %s", namespace, w.GetName()))
			}
		}
	}
	return nil

}

func printList(app config.App) string {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tVERSION\tTYPE")
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", app.Name, app.Version, app.Repository.Type))
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
