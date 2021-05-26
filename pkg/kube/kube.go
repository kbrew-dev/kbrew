package kube

import (
	"context"

	"github.com/kanisterio/kanister/pkg/kube"
	osversioned "github.com/openshift/client-go/apps/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

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
