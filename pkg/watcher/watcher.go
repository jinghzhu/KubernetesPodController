package watcher

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jinghzhu/KubernetesPodOperator/pkg/types"
	"github.com/jinghzhu/goutils/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PendingPodsWatcher watches the Pods in Pending phase and perform some actions.
func PendingPodsWatcher(ctx context.Context, namespace string, watchInterval time.Duration) {
	defer func() { // Always keep it running.
		if err := recover(); err != nil {
			fmt.Printf("Catch panic in PendingPodsWatcher: %+v\n", err)
		}
		go PendingPodsWatcher(ctx, namespace, watchInterval)
	}()
	for {
		fmt.Println("Start to check Pending Pods")
		err := CheckPods(
			ctx,
			os.Getenv("KUBECONFIG"),
			namespace,
			metav1.ListOptions{
				FieldSelector: string(types.StatusPhasePending),
				// LabelSelector: "something=",
			},
			checkPendingPod,
		)
		if err != nil {
			fmt.Printf("Fail to check Pods in Pending status: %+v\n", err)
		} else {
			fmt.Println("Successfully finish checking Pods in Pending status")
		}
		time.Sleep(watchInterval)
	}
}

// checkPendingPod checks the Pending Pod.
func checkPendingPod(ctx context.Context, pod *corev1.Pod) error {
	defer utils.PanicHandler()
	sub := time.Now().Sub(pod.ObjectMeta.CreationTimestamp.Time)
	if sub.Seconds() < types.MaxPendingPeriod {
		return nil
	}

	return processBadPendingPod(ctx, pod)
}
