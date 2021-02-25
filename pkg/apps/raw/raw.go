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

	"github.com/infracloudio/kbrew/pkg/apps"
	"github.com/infracloudio/kbrew/pkg/config"
	"github.com/infracloudio/kbrew/pkg/kube"
)

type method string

const (
	install   method = "create"
	uninstall method = "delete"
	upgrade   method = "apply"
)

var yamlDelimiter = regexp.MustCompile(`(?m)^---$`)

type RawApp struct {
	apps.BaseApp
	KubeCli  kubernetes.Interface
	OSAppCli osversioned.Interface
}

func New(c config.App, namespace string) (*RawApp, error) {
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

	rApp := &RawApp{
		BaseApp: apps.BaseApp{
			App: c,
		},
		KubeCli:  cli,
		OSAppCli: osCli,
	}
	if namespace != "" {
		rApp.Namespace = namespace
	}
	return rApp, nil
}

func (r *RawApp) Install(ctx context.Context, name, version string, options map[string]string) error {
	// TODO(@prasad): Use go sdks
	if err := kubectlCommand(install, name, r.Namespace, r.App.Repository.URL); err != nil {
		return err
	}
	return r.waitForReady(ctx)
}

func (r *RawApp) Uninstall(ctx context.Context, name string) error {
	// TODO(@prasad): Use go sdks
	return kubectlCommand(uninstall, name, r.Namespace, r.App.Repository.URL)
}

func (r *RawApp) Search(ctx context.Context, name string) (string, error) {
	return printList(r.App), nil
}

func kubectlCommand(m method, name, namespace, url string) error {
	c := exec.Command("kubectl", string(m), "-f", url)
	if namespace != "" {
		c.Args = append(c.Args, "--namespace", namespace)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (r *RawApp) waitForReady(ctx context.Context) error {
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

		namespace := r.Namespace
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
