package kube

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kanisterio/kanister/pkg/kube"
	osversioned "github.com/openshift/client-go/apps/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var k8sVersion string

// Client contains Kubernetes clients to call APIs
type Client struct {
	KubeCli      kubernetes.Interface
	OSCli        osversioned.Interface
	DiscoveryCli discovery.DiscoveryInterface
}

// WaitForPodReady waits till the pod gets ready
func WaitForPodReady(ctx context.Context, kubeCli kubernetes.Interface, namespace string, name string) error {
	return kube.WaitForPodReady(ctx, kubeCli, namespace, name)
}

// WaitForDeploymentReady waits till the deployment gets ready
func WaitForDeploymentReady(ctx context.Context, kubeCli kubernetes.Interface, namespace string, name string) error {
	return kube.WaitOnDeploymentReady(ctx, kubeCli, namespace, name)
}

// WaitForStatefulSetReady waits till the statefulset gets ready
func WaitForStatefulSetReady(ctx context.Context, kubeCli kubernetes.Interface, namespace string, name string) error {
	return kube.WaitOnStatefulSetReady(ctx, kubeCli, namespace, name)
}

// WaitForDeploymentConfigReady waits till the deployment config gets ready
func WaitForDeploymentConfigReady(ctx context.Context, osCli osversioned.Interface, kubeCli kubernetes.Interface, namespace string, name string) error {
	return kube.WaitOnDeploymentConfigReady(ctx, osCli, kubeCli, namespace, name)
}

// FetchNonRunningPods returns list of non running Pods owned by the workloads
func FetchNonRunningPods(ctx context.Context, workloads []corev1.ObjectReference) ([]corev1.Pod, error) {
	clis, err := NewClient()
	if err != nil {
		return nil, err
	}

	pods := []corev1.Pod{}
	for _, wRef := range workloads {
		fmt.Println("WORKLOAD", wRef)
		switch wRef.Kind {
		case "Pod":
			pod, err := clis.KubeCli.CoreV1().Pods(wRef.Namespace).Get(ctx, wRef.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
			if pod.Status.Phase != corev1.PodRunning {
				pods = append(pods, *pod)
			}

		case "Deployment":
			_, notRunningPods, err := kube.DeploymentPods(ctx, clis.KubeCli, wRef.Namespace, wRef.Name)
			if err != nil {
				return nil, err
			}
			pods = append(pods, notRunningPods...)
		case "StatefulSet":
			_, notRunningPods, err := kube.StatefulSetPods(ctx, clis.KubeCli, wRef.Namespace, wRef.Name)
			if err != nil {
				return nil, err
			}
			pods = append(pods, notRunningPods...)

		case "DeploymentConfig":
			_, notRunningPods, err := kube.DeploymentConfigPods(ctx, clis.OSCli, clis.KubeCli, wRef.Namespace, wRef.Name)
			if err != nil {
				return nil, err
			}
			pods = append(pods, notRunningPods...)
		}
	}
	return pods, nil
}

func newConfig() (*rest.Config, error) {
	kubeConfig, err := rest.InClusterConfig()
	if err == nil {
		return kubeConfig, nil
	}
	kubeconfig, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// NewClient initializes and returns client object
func NewClient() (*Client, error) {
	kubeConfig, err := newConfig()
	if err != nil {
		return nil, err
	}
	disClient, err := discovery.NewDiscoveryClientForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	kubeCli, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	osCli, err := osversioned.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &Client{
		KubeCli:      kubeCli,
		DiscoveryCli: disClient,
		OSCli:        osCli,
	}, nil
}

// GetK8sVersion returns Kubernetes server version
func GetK8sVersion() (string, error) {
	if k8sVersion != "" {
		return k8sVersion, nil
	}
	clis, err := NewClient()
	if err != nil {
		return "", err
	}
	versionInfo, err := clis.DiscoveryCli.ServerVersion()
	if err != nil {
		return "", err
	}
	// Store the version in a global var to avoid repeatitive API calls
	k8sVersion = versionInfo.String()
	return k8sVersion, nil
}
