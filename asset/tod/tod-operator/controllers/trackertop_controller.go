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
	"fmt"
	"reflect"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	// Reconcile Tracker' Deployment & Service
	tsa := buildServiceAccount(tto.Spec.Where)
	foundServiceAccount := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: tsa.Name, Namespace: tsa.Namespace}, foundServiceAccount)
	if err != nil && errors.IsNotFound(err) {
		//SA
		l.Info("Creating ServiceAccount for Tracker", "serviceaccount", tsa.Name)
		if err := ctrl.SetControllerReference(&tto, tsa, r.Scheme); err != nil {
			l.Error(err, "Unable to Set OwnerReferences to Tracker's ServiceAccount", "tsa", tsa)
			return ctrl.Result{}, err
		}
		if err = r.Create(ctx, tsa); err != nil {
			l.Error(err, "Unable to create ServiceAccount for Tracker", "tsa", tsa)
			return ctrl.Result{}, err
		}
		//Role
		tr := buildRole(tsa)
		l.Info("Creating Role for ServiceAccount", "role", tr.Name)
		if err := ctrl.SetControllerReference(&tto, tr, r.Scheme); err != nil {
			l.Error(err, "Unable to Set OwnerReferences to Role", "tr", tr)
			return ctrl.Result{}, err
		}
		if err = r.Create(ctx, tr); err != nil {
			l.Error(err, "Unable to create Role", "tr", tr)
			return ctrl.Result{}, err
		}
		//RoleBinding
		trb := buildRoleBinding(tr, tsa)
		l.Info("Creating RoleBinding for ServiceAccount", "rolebinding", trb.Name)
		if err := ctrl.SetControllerReference(&tto, trb, r.Scheme); err != nil {
			l.Error(err, "Unable to Set OwnerReferences to RoleBinding", "trb", trb)
			return ctrl.Result{}, err
		}
		if err = r.Create(ctx, trb); err != nil {
			l.Error(err, "Unable to create RoleBinding", "tr", tr)
			return ctrl.Result{}, err
		}
		//ClusterRole
		tcr := buildClusterRole(tsa)
		l.Info("Creating ClusterRole for ServiceAccount", "clusterrole", tcr.Name)
		if err := ctrl.SetControllerReference(&tto, tcr, r.Scheme); err != nil {
			l.Error(err, "Unable to Set OwnerReferences to ClusterRole", "tcr", tcr)
			return ctrl.Result{}, err
		}
		if err = r.Create(ctx, tcr); err != nil {
			l.Error(err, "Unable to create ClusterRole", "tcr", tcr)
			return ctrl.Result{}, err
		}
		//ClusterRoleBinding
		tcrb := buildClusterRoleBinding(tcr, tsa)
		l.Info("Creating ClusterRoleBinding for ServiceAccount", "clusterrolebinding", tcrb.Name)
		if err := ctrl.SetControllerReference(&tto, tcrb, r.Scheme); err != nil {
			l.Error(err, "Unable to Set OwnerReferences to ClusterRoleBinding", "tcrb", tcrb)
			return ctrl.Result{}, err
		}
		if err = r.Create(ctx, tcrb); err != nil {
			l.Error(err, "Unable to create RoleBinding", "tr", tr)
			return ctrl.Result{}, err
		}

	}

	for _, tracker := range tto.Spec.Trackers {

		td := buildTrackerDeploy(tracker, tsa, tto.Spec.Redis, tto.Spec.Where)
		foundDeployment := &appv1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{Name: td.Name, Namespace: td.Namespace}, foundDeployment)
		if err != nil && errors.IsNotFound(err) {
			l.Info("Creating Deployment for Tracker", "deployment", td.Name)
			if err := ctrl.SetControllerReference(&tto, td, r.Scheme); err != nil {
				l.Error(err, "Unable to Set OwnerReferences to Tracker's Deployment", "td", td)
				return ctrl.Result{}, err
			}
			if err = r.Create(ctx, td); err != nil {
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

		ts := buildTrackerService(tracker, tto.Spec.Where)
		foundService := &corev1.Service{}
		err = r.Get(ctx, types.NamespacedName{Name: ts.Name, Namespace: ts.Namespace}, foundService)
		if err != nil && errors.IsNotFound(err) {
			l.Info("Creating Service for Tracker", "service", ts.Name)
			if err := ctrl.SetControllerReference(&tto, ts, r.Scheme); err != nil {
				l.Error(err, "Unable to Set OwnerReferences to Tracker's Service", "ts", ts)
				return ctrl.Result{}, err
			}
			if err = r.Create(ctx, ts); err != nil {
				l.Error(err, "Unable to create Service for Tracker", "ts", ts)
				return ctrl.Result{}, err
			}
		} else if err == nil {
			if !reflect.DeepEqual(foundService.Spec.Selector, ts.Spec.Selector) ||
				foundService.Spec.Type != ts.Spec.Type ||
				!reflect.DeepEqual(foundService.Spec.Ports, ts.Spec.Ports) {

				foundService.Spec.Type = ts.Spec.Type
				foundService.Spec.Selector = ts.Spec.Selector
				foundService.Spec.Ports = ts.Spec.Ports
				l.Info("Updating Service for Tracker", "service", foundService.Name)
				if err = r.Update(ctx, foundService); err != nil {
					l.Error(err, "Unable to update Service for Tracker", "ts", ts)
					return ctrl.Result{}, err
				}
			}
		}

	}

	// Reconcile Redis' Deployment & Service
	if tto.Spec.Redis.Name != "" {
		rd := buildRedisDeploy(tto.Spec.Redis, tto.Spec.Where)
		foundDeployment := &appv1.Deployment{}
		err := r.Get(ctx, types.NamespacedName{Name: rd.Name, Namespace: rd.Namespace}, foundDeployment)
		if err != nil && errors.IsNotFound(err) {
			l.Info("Creating Deployment for Redis", "deployment", rd.Name)
			if err := ctrl.SetControllerReference(&tto, rd, r.Scheme); err != nil {
				l.Error(err, "Unable to Set OwnerReferences to Redis's Deployment", "rd", rd)
				return ctrl.Result{}, err
			}
			if err = r.Create(ctx, rd); err != nil {
				l.Error(err, "Unable to create Deployment for Redis", "td", rd)
				return ctrl.Result{}, err
			}
		}

		rs := buildRedisService(tto.Spec.Redis, tto.Spec.Where)
		foundService := &corev1.Service{}
		err = r.Get(ctx, types.NamespacedName{Name: rs.Name, Namespace: rs.Namespace}, foundService)
		if err != nil && errors.IsNotFound(err) {
			l.Info("Creating Service for Redis", "service", rs.Name)
			if err := ctrl.SetControllerReference(&tto, rs, r.Scheme); err != nil {
				l.Error(err, "Unable to Set OwnerReferences to Redis's Service", "rs", rs)
				return ctrl.Result{}, err
			}
			if err = r.Create(ctx, rs); err != nil {
				l.Error(err, "Unable to create Service for Redis", "rs", rs)
				return ctrl.Result{}, err
			}
		} else if err == nil {
			if !reflect.DeepEqual(foundService.Spec.Selector, rs.Spec.Selector) ||
				foundService.Spec.Type != rs.Spec.Type ||
				!reflect.DeepEqual(foundService.Spec.Ports, rs.Spec.Ports) {

				foundService.Spec.Type = rs.Spec.Type
				foundService.Spec.Selector = rs.Spec.Selector
				foundService.Spec.Ports = rs.Spec.Ports
				l.Info("Updating Service for Redis", "service", foundService.Name)
				if err = r.Update(ctx, foundService); err != nil {
					l.Error(err, "Unable to update Service for Redis", "rs", rs)
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

func buildTrackerDeploy(tracker trackerv1.Tracker, tsa *corev1.ServiceAccount, redis trackerv1.ThirdParty, ns string) *appv1.Deployment {
	terminate := int64(60)
	td := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tracker.Name,
			Namespace: ns,
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
					TerminationGracePeriodSeconds: &terminate,
					ServiceAccountName:            tsa.GetName(),
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
									Value: fmt.Sprintf("%s.%s.svc.cluster.local:%d", redis.Host, ns, redis.Port),
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
									corev1.ResourceCPU: resource.MustParse("200m"),
								},

								Limits: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse("500m"),
								},
							},
						},
					},
				},
			},
		},
	}
	return &td
}

func buildTrackerService(tracker trackerv1.Tracker, ns string) *corev1.Service {
	ts := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tracker.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(tracker.ServingType),
			Selector: map[string]string{
				"app": "tracker",
				"svc": tracker.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       8000,
					TargetPort: intstr.IntOrString{IntVal: 8000},
				},
				{
					Name:       "tcp",
					Port:       8008,
					TargetPort: intstr.IntOrString{IntVal: 8008},
				},
			},
		},
	}
	return &ts
}

func buildRedisDeploy(redis trackerv1.ThirdParty, ns string) *appv1.Deployment {
	rd := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: ns,
		},
		Spec: appv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "redis-cart",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "false",
					},
					Labels: map[string]string{
						"app": "redis-cart",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "redis",
							Image: redis.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "tcp",
									ContainerPort: redis.Port,
								},
							},
							ReadinessProbe: &corev1.Probe{
								PeriodSeconds: 5,
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{IntVal: redis.Port},
									},
								},
							},

							LivenessProbe: &corev1.Probe{
								PeriodSeconds: 5,
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{IntVal: 6379},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("1024Mi"),
								},

								Limits: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse("2000m"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "redis-data",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "redis-data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}
	return &rd
}

func buildRedisService(redis trackerv1.ThirdParty, ns string) *corev1.Service {
	rs := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name,
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType("ClusterIP"),
			Selector: map[string]string{
				"app": "redis-cart",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp",
					Port:       redis.Port,
					TargetPort: intstr.IntOrString{IntVal: redis.Port},
				},
			},
		},
	}
	return &rs
}

func buildServiceAccount(ns string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tracker-sa",
			Namespace: ns,
		},
	}
}
func buildRole(tsa *corev1.ServiceAccount) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsa.GetName() + "-role",
			Namespace: tsa.GetNamespace(),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "services", "pods", "endpoints"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"extensions"},
				Resources: []string{"deployments"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}
func buildRoleBinding(tr *rbacv1.Role, tsa *corev1.ServiceAccount) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tr.Name + "-rolebinding",
			Namespace: tr.Namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     tr.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: tsa.Name,
			},
		},
	}
}

func buildClusterRole(tsa *corev1.ServiceAccount) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsa.GetName() + "-clusterrole",
			Namespace: tsa.GetNamespace(),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "services", "pods", "endpoints"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"extensions"},
				Resources: []string{"deployments"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}
func buildClusterRoleBinding(tcr *rbacv1.ClusterRole, tsa *corev1.ServiceAccount) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tcr.Name + "-clusterrolebinding",
			Namespace: tcr.Namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     tcr.Name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      tsa.Name,
				Namespace: tsa.Namespace,
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *TrackerTopReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&trackerv1.TrackerTop{}).
		Owns(&appv1.Deployment{}).
		Complete(r)
}
