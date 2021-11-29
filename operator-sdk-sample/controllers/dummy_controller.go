/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrltypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dummyv1alpha1 "github.com/anynines/tmp-homework-mj/api/v1alpha1"
	"github.com/google/go-cmp/cmp"
)

// DummyReconciler reconciles a Dummy object
type DummyReconciler struct {
	client.Client
	recorder record.EventRecorder

	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=interview.com,resources=dummies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=interview.com,resources=dummies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=interview.com,resources=dummies/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DummyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Retrieve Dummy object
	var dummy dummyv1alpha1.Dummy
	if err := r.Get(ctx, req.NamespacedName, &dummy); err != nil {
		logger.Error(err, "unable to fetch Dummy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Process Dummy
	logMessage := fmt.Sprintf("Process Dummy Object: Name: %v, Namespace: %v, Message: %v", dummy.Name, dummy.Namespace, dummy.Spec.Message)
	logger.Info(logMessage)

	// ensurePod
	pod, err := r.ensurePod(ctx, &dummy)
	if err != nil {
		logger.Error(err, "unable to ensure Pod %v")
		r.recorder.Eventf(&dummy, corev1.EventTypeWarning, "FailedCreatePod", "error to creating pod, %v", err)
	}

	r.updateStatus(ctx, &dummy, pod)

	// Update Dummy status
	return ctrl.Result{}, nil
}

func (r *DummyReconciler) updateStatus(ctx context.Context, d *dummyv1alpha1.Dummy, pod *corev1.Pod) error {
	logger := log.FromContext(ctx)

	cp := d.Status.DeepCopy()
	logger.Info("Updating dummy status")
	d.Status.SpecEcho = d.Spec.Message
	d.Status.PodStatus = string(pod.Status.Phase)
	if !cmp.Equal(d.Status, cp) {
		return r.Status().Update(ctx, d)
	}

	return nil
}

// ensurePod ensures the Pod exists
func (r *DummyReconciler) ensurePod(ctx context.Context, d *dummyv1alpha1.Dummy) (*corev1.Pod, error) {
	logger := log.FromContext(ctx)

	// Retrieve Pod
	var desiredPod corev1.Pod
	nn := ctrltypes.NamespacedName{Namespace: d.ObjectMeta.Namespace, Name: d.ObjectMeta.Name}
	if err := r.Get(ctx, nn, &desiredPod); err != nil {
		logger.Error(err, "unable to get Pod")
		if apierrors.IsNotFound(err) {
			// Create Pod
			pod := corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: d.Name, Namespace: d.Namespace, Labels: d.DeepCopy().Labels},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "nginx", Image: "nginx"}}},
			}

			if err := ctrl.SetControllerReference(d, &pod, r.Scheme); err != nil {
				return nil, err
			}
			if err := r.Create(ctx, &pod); err != nil {
				return nil, err
			} else if err == nil {
				r.recorder.Eventf(d, corev1.EventTypeNormal, "PodCreated", "Pod %v is created", d.Name)
			}
		}
	}

	return &desiredPod, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("dummy-controller")

	return ctrl.NewControllerManagedBy(mgr).
		For(&dummyv1alpha1.Dummy{}).
		Owns(&corev1.Event{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
