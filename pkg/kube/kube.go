package kube

import (
	"context"
	"time"

	"github.com/kanisterio/kanister/pkg/kube"
	osversioned "github.com/openshift/client-go/apps/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

const workloadReadyTimeout = 15 * time.Minute

func WaitForPodReady(ctx context.Context, kubeCli kubernetes.Interface, namespace string, name string) error {
	ctx, cancel := context.WithTimeout(ctx, workloadReadyTimeout)
	defer cancel()
	return kube.WaitForPodReady(ctx, kubeCli, namespace, name)
}

func WaitForDeploymentReady(ctx context.Context, kubeCli kubernetes.Interface, namespace string, name string) error {
	ctx, cancel := context.WithTimeout(ctx, workloadReadyTimeout)
	defer cancel()
	return kube.WaitOnDeploymentReady(ctx, kubeCli, namespace, name)
}

func WaitForStatefulSetReady(ctx context.Context, kubeCli kubernetes.Interface, namespace string, name string) error {
	ctx, cancel := context.WithTimeout(ctx, workloadReadyTimeout)
	defer cancel()
	return kube.WaitOnStatefulSetReady(ctx, kubeCli, namespace, name)
}

func WaitForDeploymentConfigReady(ctx context.Context, osCli osversioned.Interface, kubeCli kubernetes.Interface, namespace string, name string) error {
	ctx, cancel := context.WithTimeout(ctx, workloadReadyTimeout)
	defer cancel()
	return kube.WaitOnDeploymentConfigReady(ctx, osCli, kubeCli, namespace, name)
}
