/*


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
	"os"
	"strings"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const clusterRoleName = "nginx-ingress-role"
const sccName = "nginx-ingress-scc"
const finalizer = "finalizer.nginxingresscontroller.k8s.nginx.org"

// NginxIngressControllerReconciler reconciles a NginxIngressController object
type NginxIngressControllerReconciler struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	sccAPIExists bool
}

// +kubebuilder:rbac:groups=k8s.nginx.org,resources=nginxingresscontrollers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.nginx.org,resources=nginxingresscontrollers/status,verbs=get;update;patch

// Reconcile watches kubernetes events against the custom resources
func (r *NginxIngressControllerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("nginxingresscontroller", req.NamespacedName)

	// your logic here
	instance := &k8sv1alpha1.NginxIngressController{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		return ctrl.Result{}, nil
	} else if err != nil {
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	err = r.handleFinalizers(log, instance)
	if err != nil {
		return ctrl.Result{}, err
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

	err = r.checkPrerequisites(ctx, instance, log)
	if err != nil {
		return ctrl.Result{}, err
	}

	if strings.ToLower(instance.Spec.Type) == "deployment" {
		dep := deploymentForNginxIngressController(instance)
		found := &appsv1.Deployment{}
		err = r.Get(ctx, types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new Deployment for NGINX Ingress Controller", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

			err = controllerutil.SetControllerReference(instance, dep, r.Scheme)
			if err != nil {
				log.Error(err, "Error setting controller reference")
				return ctrl.Result{}, err
			}

			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}
		} else if err != nil {
			return ctrl.Result{}, err
		} else if hasDeploymentChanged(found, instance) {
			log.Info("NginxIngressController spec has changed, updating Deployment")
			updated := updateDeployment(found, instance)
			err = r.Update(ctx, updated)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		// Remove possible DaemonSet
		ds := daemonSetForNginxIngressController(instance)
		err = r.deleteIfExists(ds.Name, ds.Namespace, ds)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if strings.ToLower(instance.Spec.Type) == "daemonset" {
		ds := daemonSetForNginxIngressController(instance)
		found := &appsv1.DaemonSet{}
		err = r.Get(ctx, types.NamespacedName{Name: ds.Name, Namespace: ds.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating a new DaemonSet for NGINX Ingress Controller", "DaemonSet.Namespace", ds.Namespace, "DaemonSet.Name", ds.Name)

			err = controllerutil.SetControllerReference(instance, ds, r.Scheme)
			if err != nil {
				log.Error(err, "Error setting controller reference")
				return ctrl.Result{}, err
			}

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
				return ctrl.Result{}, err
			}
		}

		// Remove possible Deployment
		dep := deploymentForNginxIngressController(instance)
		err = r.deleteIfExists(dep.Name, dep.Namespace, dep)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	svc := serviceForNginxIngressController(instance)
	svcFound := &v1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Service for NGINX Ingress Controller", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)

		err = controllerutil.SetControllerReference(instance, svc, r.Scheme)
		if err != nil {
			log.Error(err, "Error setting controller reference")
			return ctrl.Result{}, err
		}

		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else if hasServiceChanged(svcFound, instance) {
		log.Info("NginxIngressController spec has changed, updating Service")
		updated := updateService(svcFound, instance)
		err = r.Update(ctx, updated)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	cm := configMapForNginxIngressController(instance)
	cmFound := &v1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Namespace: cm.Namespace, Name: cm.Name}, cmFound)
	if err != nil {
		return ctrl.Result{}, err
	}

	if hasConfigMapChanged(cmFound, instance) {
		err = r.Update(ctx, cm)
		log.Info("NginxIngressController spec has changed, updating ConfigMap")
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if !instance.Status.Deployed {
		instance.Status.Deployed = true
		err := r.Status().Update(ctx, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager defines the child resources associated with the custom resource
func (r *NginxIngressControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1alpha1.NginxIngressController{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

// isLocal returns true if the Operator is running locally and false if running inside a cluster
func isLocal() bool {
	_, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if !ok {
		return true
	}

	return false
}

func createKICCustomResourceDefinitions(ctx context.Context, mgr manager.Manager, log logr.Logger) error {

	if isLocal() {
		log.Info("Skipping KIC CRDs creation; not running in a cluster")
		return nil
	}

	// Create CRDs with a different client (apiextensions)
	apixClient, err := apixv1beta1client.NewForConfig(mgr.GetConfig())
	if err != nil {
		log.Error(err, "unable to create client for CRD registration")
		return err
	}

	crds, err := kicCRDs()

	if err != nil {
		return err
	}

	crdsClient := apixClient.CustomResourceDefinitions()
	for _, crd := range crds {
		oldCRD, err := crdsClient.Get(ctx, crd.Name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info(fmt.Sprintf("no previous CRD %v found, creating a new one.", crd.Name))
				_, err = crdsClient.Create(ctx, crd, metav1.CreateOptions{})
				if err != nil {
					return fmt.Errorf("error creating CustomResourceDefinition %v: %v", crd.Name, err)
				}
			} else {
				return fmt.Errorf("error getting CustomResourceDefinition %v: %v", crd.Name, err)
			}
		} else {
			// Update CRDs if they already exist
			log.Info(fmt.Sprintf("previous CustomResourceDefinition %v found, updating.", crd.Name))
			oldCRD.Spec = crd.Spec
			_, err = crdsClient.Update(ctx, oldCRD, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("error updating CustomResourceDefinition %v: %v", crd.Name, err)
			}
		}
	}

	return nil
}

// create common resources shared by all the Ingress Controllers
func createCommonResources(ctx context.Context, mgr manager.Manager, sccAPIExists bool, log logr.Logger) error {

	// Create ClusterRole and ClusterRoleBinding for all the NginxIngressController resources.
	clientReader := mgr.GetAPIReader()
	clientWriter := mgr.GetClient()
	cr := clusterRoleForNginxIngressController(clusterRoleName)

	err := clientReader.Get(ctx, types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, cr)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("no previous ClusterRole found, creating a new one.")
			err = clientWriter.Create(ctx, cr)
			if err != nil {
				return fmt.Errorf("error creating ClusterRole: %v", err)
			}
		} else {
			return fmt.Errorf("error getting ClusterRole: %v", err)
		}
	} else {
		// For updates in the ClusterRole permissions (eg new CRDs of the Ingress Controller).
		log.Info("previous ClusterRole found, updating.")
		cr := clusterRoleForNginxIngressController(clusterRoleName)
		err = clientWriter.Update(ctx, cr)
		if err != nil {
			return fmt.Errorf("error updating ClusterRole: %v", err)
		}
	}

	crb := clusterRoleBindingForNginxIngressController(clusterRoleName)

	err = clientReader.Get(ctx, types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil && errors.IsNotFound(err) {
		log.Info("no previous ClusterRoleBinding found, creating a new one.")
		err = clientWriter.Create(ctx, crb)
	}

	if err != nil {
		return fmt.Errorf("error creating ClusterRoleBinding: %v", err)
	}

	err = createKICCustomResourceDefinitions(ctx, mgr, log)
	if err != nil {
		return fmt.Errorf("error creating KIC CRDs: %v", err)
	}

	if sccAPIExists {
		log.Info("OpenShift detected as platform.")

		scc := sccForNginxIngressController(sccName)

		err = clientReader.Get(ctx, types.NamespacedName{Name: sccName, Namespace: v1.NamespaceAll}, scc)
		if err != nil && errors.IsNotFound(err) {
			log.Info("no previous SecurityContextConstraints found, creating a new one.")
			err = clientWriter.Create(ctx, scc)
		}

		if err != nil {
			return fmt.Errorf("error creating SecurityContextConstraints: %v", err)
		}
	}

	return nil
}

// checkPrerequisites creates all necessary objects before the deployment of a new Ingress Controller.
func (r *NginxIngressControllerReconciler) checkPrerequisites(ctx context.Context, instance *k8sv1alpha1.NginxIngressController, log logr.Logger) error {
	sa := serviceAccountForNginxIngressController(instance)
	err := controllerutil.SetControllerReference(instance, sa, r.Scheme)
	if err != nil {
		return err
	}

	existed, err := r.createIfNotExists(sa)
	if err != nil {
		return err
	}

	if !existed {
		log.Info("ServiceAccount created", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
	}

	// Assign this new ServiceAccount to the ClusterRoleBinding (if is not present already)
	crb := clusterRoleBindingForNginxIngressController(clusterRoleName)

	err = r.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil {
		return err
	}

	subject := subjectForServiceAccount(sa.Namespace, sa.Name)
	found := false
	for _, s := range crb.Subjects {
		if s.Name == subject.Name && s.Namespace == subject.Namespace {
			found = true
			break
		}
	}

	if !found {
		crb.Subjects = append(crb.Subjects, subject)

		err = r.Update(context.TODO(), crb)
		if err != nil {
			return err
		}
	}

	cm := configMapForNginxIngressController(instance)

	err = controllerutil.SetControllerReference(instance, cm, r.Scheme)
	if err != nil {
		return err
	}

	existed, err = r.createIfNotExists(cm)
	if err != nil {
		return err
	}

	if !existed {
		log.Info("ConfigMap created", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
	}

	if instance.Spec.DefaultSecret == "" {
		err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, &v1.Secret{})

		if err != nil && errors.IsNotFound(err) {
			secret, err := defaultSecretForNginxIngressController(instance)
			if err != nil {
				return err
			}

			err = controllerutil.SetControllerReference(instance, secret, r.Scheme)
			if err != nil {
				return err
			}

			err = r.Create(context.TODO(), secret)
			if err != nil {
				return err
			}

			log.Info("Warning! A custom self-signed TLS Secret has been created for the default server. "+
				"Update your 'DefaultSecret' with your own Secret in Production",
				"Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)

		} else if err != nil {
			return err
		}
	}

	if r.sccAPIExists {
		// Assign this new User to the SCC (if is not present already)
		scc := sccForNginxIngressController(sccName)

		err = r.Get(context.TODO(), types.NamespacedName{Name: sccName, Namespace: v1.NamespaceAll}, scc)
		if err != nil {
			return err
		}

		user := userForSCC(sa.Namespace, sa.Name)
		found := false
		for _, u := range scc.Users {
			if u == user {
				found = true
				break
			}
		}

		if !found {
			scc.Users = append(scc.Users, user)

			err = r.Update(context.TODO(), scc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// createIfNotExists creates a new object. If the object exists, does nothing. It returns whether the object existed before or not.
func (r *NginxIngressControllerReconciler) createIfNotExists(object runtime.Object) (bool, error) {
	err := r.Create(context.TODO(), object)
	if err != nil && errors.IsAlreadyExists(err) {
		return true, nil
	}

	return false, err
}

// deleteIfExists deletes an object if it exists in the cluster.
func (r *NginxIngressControllerReconciler) deleteIfExists(name string, namespace string, object runtime.Object) error {
	err := r.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, object)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err == nil {
		err := r.Delete(context.TODO(), object)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *NginxIngressControllerReconciler) finalizeNginxIngressController(log logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	crb := clusterRoleBindingForNginxIngressController(clusterRoleName)

	err := r.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil {
		return err
	}

	var subjects []rbacv1.Subject
	for _, s := range crb.Subjects {
		if s.Name != instance.Name && s.Namespace != instance.Namespace {
			subjects = append(subjects, s)
		}
	}

	crb.Subjects = subjects

	err = r.Update(context.TODO(), crb)
	if err != nil {
		return err
	}

	if r.sccAPIExists {
		scc := sccForNginxIngressController(sccName)

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

	log.Info("Successfully finalized NginxIngressController")
	return nil
}

func (r *NginxIngressControllerReconciler) addFinalizer(log logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	log.Info("Adding Finalizer for the NginxIngressController")
	instance.SetFinalizers(append(instance.GetFinalizers(), finalizer))

	err := r.Update(context.TODO(), instance)
	if err != nil {
		log.Error(err, "Failed to update NginxIngressController with finalizer")
		return err
	}

	return nil
}

func (r *NginxIngressControllerReconciler) handleFinalizers(log logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	// If instance has been marked to be deleted
	if instance.GetDeletionTimestamp() != nil {
		if contains(instance.GetFinalizers(), finalizer) {
			err := r.finalizeNginxIngressController(log, instance)
			if err != nil {
				return err
			}

			instance.SetFinalizers(remove(instance.GetFinalizers(), finalizer))
			err = r.Update(context.TODO(), instance)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if !contains(instance.GetFinalizers(), finalizer) {
		err := r.addFinalizer(log, instance)
		if err != nil {
			return err
		}
	}

	return nil
}
