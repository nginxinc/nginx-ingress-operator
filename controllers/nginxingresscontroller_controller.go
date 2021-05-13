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
	"strings"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
)

const (
	clusterRoleName = "nginx-ingress-role"
	sccName         = "nginx-ingress-scc"
	finalizer       = "nginxingresscontroller.k8s.nginx.org/finalizer"
)

// NginxIngressControllerReconciler reconciles a NginxIngressController object
type NginxIngressControllerReconciler struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	SccAPIExists bool
	Mgr          ctrl.Manager
}

//+kubebuilder:rbac:groups=k8s.nginx.org,resources=nginxingresscontrollers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.nginx.org,resources=nginxingresscontrollers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.nginx.org,resources=nginxingresscontrollers/finalizers,verbs=update
//+kubebuilder:rbac:groups=k8s.nginx.org;appprotect.f5.com,resources=*,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;replicasets;statefulsets,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=list;watch;get
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingressclasses,verbs=get;create;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=update
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=get;create;delete;update
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles;clusterrolebindings;roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups=security.openshift.io,resources=securitycontextconstraints,verbs=create;update;get;list;watch

//+kubebuilder:rbac:groups="",resources=pods;services;services/finalizers;endpoints;persistentvolumeclaims;events;configmaps;secrets;serviceaccounts;namespaces,verbs=create;update;get;list;watch;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NginxIngressController object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NginxIngressControllerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("nginxingresscontroller", req.NamespacedName)
	log.Info("Reconciling NginxIngressController")

	instance := &k8sv1alpha1.NginxIngressController{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		log.Info("NginxIngressController resource not found. Ignoring since object must be deleted")
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get NginxIngressController")
		return ctrl.Result{}, err
	}

	// Check if the NginxIngressController instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isNginxIngressControllerMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isNginxIngressControllerMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(instance, finalizer) {
			// Run finalization logic for nginxingresscontrollerFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeNginxIngressController(log, instance); err != nil {
				return ctrl.Result{}, err
			}

			// Remove nginxingresscontrollerFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(instance, finalizer)
			err := r.Update(ctx, instance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	if !controllerutil.ContainsFinalizer(instance, finalizer) {
		controllerutil.AddFinalizer(instance, finalizer)
		err = r.Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Namespace could have been deleted in the middle of the reconcile
	ns := &v1.Namespace{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Namespace, Namespace: v1.NamespaceAll}, ns)
	if (err != nil && errors.IsNotFound(err)) || (ns.Status.Phase == "Terminating") {
		log.Info(fmt.Sprintf("The namespace '%v' does not exist or is in Terminating status, canceling Reconciling", instance.Namespace))
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Failed to check if namespace exists")
		return ctrl.Result{}, err
	}

	if err := r.createCommonResources(log, instance); err != nil {
		return ctrl.Result{}, err
	}

	err = r.checkPrerequisites(log, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if strings.ToLower(instance.Spec.Type) == "deployment" {
		found := &appsv1.Deployment{}
		dep, err := r.deploymentForNginxIngressController(instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Deployment for NGINX Ingress Controller", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		} else if hasDeploymentChanged(dep, instance) {
			log.Info("NginxIngressController spec has changed, updating Deployment")
			updated := updateDeployment(found, instance)
			err = r.Update(ctx, updated)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Remove possible DaemonSet
		ds, err := r.daemonSetForNginxIngressController(instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Delete(ctx, ds); client.IgnoreNotFound(err) != nil {
			return reconcile.Result{}, err
		}
	} else if strings.ToLower(instance.Spec.Type) == "daemonset" {
		found := &appsv1.DaemonSet{}
		ds, err := r.daemonSetForNginxIngressController(instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new DaemonSet for NGINX Ingress Controller", "DaemonSet.Namespace", ds.Namespace, "DaemonSet.Name", ds.Name)

			err = r.Create(ctx, ds)
			if err != nil {
				log.Error(err, "Failed to create new DaemonSet", "DaemonSet.Namespace", ds.Namespace, "DaemonSet.Name", ds.Name)
				return ctrl.Result{}, err
			}
		} else if err != nil {
			return ctrl.Result{}, err
		} else if hasDaemonSetChanged(ds, instance) {
			log.Info("NginxIngressController spec has changed, updating DaemonSet")
			updated := updateDaemonSet(found, instance)
			err = r.Update(ctx, updated)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Remove possible Deployment
		dep, err := r.deploymentForNginxIngressController(instance)
		if err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Delete(ctx, dep); client.IgnoreNotFound(err) != nil {
			return reconcile.Result{}, err
		}

	}

	svc := r.serviceForNginxIngressController(instance)
	var extraLabels map[string]string
	if instance.Spec.Service != nil {
		extraLabels = instance.Spec.Service.ExtraLabels
	}
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, svc, serviceMutateFn(svc, instance.Spec.ServiceType, extraLabels))
	log.Info(fmt.Sprintf("Service %s %s", svc.Name, res))
	if err != nil {
		return ctrl.Result{}, err
	}

	cm, err := r.configMapForNginxIngressController(instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, cm, configMapMutateFn(cm, instance.Spec.ConfigMapData))
	log.Info(fmt.Sprintf("ConfigMap %s %s", svc.Name, res))
	if err != nil {
		return ctrl.Result{}, err
	}

	if !instance.Status.Deployed {
		instance.Status.Deployed = true
		err := r.Status().Update(ctx, instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// createIfNotExists creates a new object. If the object exists, does nothing. It returns whether the object existed before or not.
func (r *NginxIngressControllerReconciler) createIfNotExists(object client.Object) (error, bool) {
	err := r.Create(context.TODO(), object)
	if err != nil && errors.IsAlreadyExists(err) {
		return nil, true
	}

	return err, false
}

func (r *NginxIngressControllerReconciler) finalizeNginxIngressController(log logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	crb := r.clusterRoleBindingForNginxIngressController(clusterRoleName)

	err := r.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil {
		return err
	}

	var subjects []rbacv1.Subject
	for _, s := range crb.Subjects {
		if s.Name != instance.Name || s.Namespace != instance.Namespace {
			subjects = append(subjects, s)
		}
	}

	crb.Subjects = subjects

	err = r.Update(context.TODO(), crb)
	if err != nil {
		return err
	}

	if r.SccAPIExists {
		scc := r.sccForNginxIngressController(sccName)

		err := r.Get(context.TODO(), types.NamespacedName{Name: sccName, Namespace: v1.NamespaceAll}, scc)
		if err != nil {
			return err
		}

		var users []string
		for _, u := range scc.Users {
			if u != userForSCC(instance.Namespace, instance.Name) {
				users = append(users, u)
			}
		}

		scc.Users = users

		err = r.Update(context.TODO(), scc)
		if err != nil {
			return err
		}
	}

	ic := r.ingressClassForNginxIngressController(instance)
	if err := r.Delete(context.TODO(), ic); client.IgnoreNotFound(err) != nil {
		return err
	}

	log.Info("Successfully finalized NginxIngressController")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxIngressControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1alpha1.NginxIngressController{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&v1.ServiceAccount{}).
		Owns(&v1.Service{}).
		Owns(&v1.ConfigMap{}).
		Owns(&v1.Secret{}).
		Complete(r)
}
