/*
Copyright 2022.

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

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	trackerv1 "tod.cc4i.xyz/tod-operator/api/v1"
)

// TrackerTopReconciler reconciles a TrackerTop object
type TrackerTopReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tracker.tod.cc4i.xyz,resources=trackertops,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tracker.tod.cc4i.xyz,resources=trackertops/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tracker.tod.cc4i.xyz,resources=trackertops/finalizers,verbs=update
//+kubebuilder:rbac:groups=tracker.tod.cc4i.xyz,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tracker.tod.cc4i.xyz,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups=tracker.tod.cc4i.xyz,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TrackerTop object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *TrackerTopReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var tto trackerv1.TrackerTop
	if err := r.Get(ctx, req.NamespacedName, &tto); err != nil {
		l.Error(err, "unable to fetch TrackerTop")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, tracker := range tto.Spec.Trackers {

		td := appv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      tracker.Name,
				Namespace: tto.Spec.Where,
				// OwnerReferences: []metav1.OwnerReference{
				// 	{
				// 		APIVersion: tto.APIVersion,
				// 		Kind:       tto.Kind,
				// 		Name:       tto.Name,
				// 		UID:        tto.UID,
				// 	},
				// },
			},
			Spec: appv1.DeploymentSpec{
				Replicas: tracker.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "tracker",
						"svc": tracker.Name,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "tracker",
							"svc": tracker.Name,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "tracker",
								Image: tracker.Image,
								Env: []corev1.EnvVar{
									{
										Name:  "TRACKER_VERSION",
										Value: tracker.Version,
									},
									{
										Name: "POD_NAME",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "metadata.name",
											},
										},
									},
									{
										Name: "POD_NAMESPACE",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "metadata.namespace",
											},
										},
									},
									{
										Name: "POD_NODE_NAME",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "spec.nodeName",
											},
										},
									},
									{
										Name:  "REDIS_SERVER_ADDRESS",
										Value: "redis-cart.run-tracker.svc.cluster.local:6379",
									},
								},
								Ports: []corev1.ContainerPort{
									{
										Name:          "http",
										ContainerPort: 8000,
									},
									{
										Name:          "tcp",
										ContainerPort: 8008,
									},
								},
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										"cpu": resource.Quantity{
											Format: "200m",
										},
									},

									Limits: corev1.ResourceList{
										"cpu": resource.Quantity{
											Format: "500m",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		foundDeployment := &appv1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{Name: td.Name, Namespace: td.Namespace}, foundDeployment)
		if err != nil && errors.IsNotFound(err) {
			l.Info("Creating Deployment for Tracker", "deployment", td.Name)
			if err := ctrl.SetControllerReference(&tto, &td, r.Scheme); err != nil {
				l.Error(err, "Unable to Set OwnerReferences to Tracker's Deployment", "td", td)
				return ctrl.Result{}, err
			}
			if err = r.Create(ctx, &td); err != nil {
				l.Error(err, "Unable to create Deployment for Tracker", "td", td)
				return ctrl.Result{}, err
			}
		} else if err == nil {
			if foundDeployment.Spec.Replicas != td.Spec.Replicas {
				foundDeployment.Spec.Replicas = td.Spec.Replicas
				l.Info("Updating Deployment for Tracker", "deployment", foundDeployment.Name)
				if err = r.Update(ctx, foundDeployment); err != nil {
					l.Error(err, "Unable to update Deployment for Tracker", "td", td)
					return ctrl.Result{}, err
				}
			}
		}

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TrackerTopReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&trackerv1.TrackerTop{}).
		Owns(&appv1.Deployment{}).
		Complete(r)
}
