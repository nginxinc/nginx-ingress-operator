package nginxingresscontroller

import (
	"context"
	commonerrors "errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/version"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	apixv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nginxingresscontroller")

const (
	clusterRoleName = "nginx-ingress-role"
	sccName         = "nginx-ingress-scc"
	finalizer       = "finalizer.nginxingresscontroller.k8s.nginx.org"
)

// Add creates a new NginxIngressController Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	sccAPIExists, err := VerifySCCAPIExists()
	if err != nil {
		return err
	}

	return add(mgr, newReconciler(mgr, sccAPIExists), sccAPIExists)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, sccAPIExists bool) reconcile.Reconciler {
	return &ReconcileNginxIngressController{client: mgr.GetClient(), scheme: mgr.GetScheme(), sccAPIExists: sccAPIExists}
}

// isLocal returns true if the Operator is running locally and false if running inside a cluster
func isLocal() bool {
	_, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		if commonerrors.Is(err, k8sutil.ErrRunLocal) {
			return true
		}
	}

	return false
}

func createKICCustomResourceDefinitions(mgr manager.Manager) error {
	reqLogger := log.WithValues("Request.Namespace", "", "Request.Name", "nginxingresscontroller-controller")

	if isLocal() {
		reqLogger.Info("Skipping KIC CRDs creation; not running in a cluster")
		return nil
	}

	// Create CRDs with a different client (apiextensions)
	apixClient, err := apixv1client.NewForConfig(mgr.GetConfig())
	if err != nil {
		reqLogger.Error(err, "unable to create client for CRD registration")
		return err
	}

	crds, err := kicCRDs()

	if err != nil {
		return err
	}

	crdsClient := apixClient.CustomResourceDefinitions()
	for _, crd := range crds {
		oldCRD, err := crdsClient.Get(context.TODO(), crd.Name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				reqLogger.Info(fmt.Sprintf("no previous CRD %v found, creating a new one.", crd.Name))
				_, err = crdsClient.Create(context.TODO(), crd, metav1.CreateOptions{})
				if err != nil {
					return fmt.Errorf("error creating CustomResourceDefinition %v: %v", crd.Name, err)
				}
			} else {
				return fmt.Errorf("error getting CustomResourceDefinition %v: %v", crd.Name, err)
			}
		} else {
			// Update CRDs if they already exist
			reqLogger.Info(fmt.Sprintf("previous CustomResourceDefinition %v found, updating.", crd.Name))
			oldCRD.Spec = crd.Spec
			_, err = crdsClient.Update(context.TODO(), oldCRD, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("error updating CustomResourceDefinition %v: %v", crd.Name, err)
			}
		}
	}

	return nil
}

// create common resources shared by all the Ingress Controllers
func createCommonResources(mgr manager.Manager, sccAPIExists bool) error {
	reqLogger := log.WithValues("Request.Namespace", "", "Request.Name", "nginxingresscontroller-controller")

	// Create ClusterRole and ClusterRoleBinding for all the NginxIngressController resources.
	clientReader := mgr.GetAPIReader()
	clientWriter := mgr.GetClient()
	cr := clusterRoleForNginxIngressController(clusterRoleName)

	err := clientReader.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, cr)

	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("no previous ClusterRole found, creating a new one.")
			err = clientWriter.Create(context.TODO(), cr)
			if err != nil {
				return fmt.Errorf("error creating ClusterRole: %v", err)
			}
		} else {
			return fmt.Errorf("error getting ClusterRole: %v", err)
		}
	} else {
		// For updates in the ClusterRole permissions (eg new CRDs of the Ingress Controller).
		reqLogger.Info("previous ClusterRole found, updating.")
		cr := clusterRoleForNginxIngressController(clusterRoleName)
		err = clientWriter.Update(context.TODO(), cr)
		if err != nil {
			return fmt.Errorf("error updating ClusterRole: %v", err)
		}
	}

	crb := clusterRoleBindingForNginxIngressController(clusterRoleName)

	err = clientReader.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("no previous ClusterRoleBinding found, creating a new one.")
		err = clientWriter.Create(context.TODO(), crb)
	}

	if err != nil {
		return fmt.Errorf("error creating ClusterRoleBinding: %v", err)
	}

	err = createKICCustomResourceDefinitions(mgr)
	if err != nil {
		return fmt.Errorf("error creating KIC CRDs: %v", err)
	}

	if sccAPIExists {
		reqLogger.Info("OpenShift detected as platform.")

		scc := sccForNginxIngressController(sccName)

		err = clientReader.Get(context.TODO(), types.NamespacedName{Name: sccName, Namespace: v1.NamespaceAll}, scc)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("no previous SecurityContextConstraints found, creating a new one.")
			err = clientWriter.Create(context.TODO(), scc)
		}

		if err != nil {
			return fmt.Errorf("error creating SecurityContextConstraints: %v", err)
		}
	}

	return nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler, sccAPIExists bool) error {
	// Create a new controller
	c, err := controller.New("nginxingresscontroller-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Create common resources
	err = createCommonResources(mgr, sccAPIExists)
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NginxIngressController
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.NginxIngressController{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to any of the following resources where NginxIngressController is their owner
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NginxIngressController{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NginxIngressController{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.ServiceAccount{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NginxIngressController{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NginxIngressController{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NginxIngressController{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NginxIngressController{},
	})
	if err != nil {
		return err
	}

	return nil
}

// checkPrerequisites creates all necessary objects before the deployment of a new Ingress Controller.
func (r *ReconcileNginxIngressController) checkPrerequisites(reqLogger logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	sa := serviceAccountForNginxIngressController(instance)
	err := controllerutil.SetControllerReference(instance, sa, r.scheme)
	if err != nil {
		return err
	}

	err, existed := r.createIfNotExists(sa)
	if err != nil {
		return err
	}

	if !existed {
		reqLogger.Info("ServiceAccount created", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
	}

	// Assign this new ServiceAccount to the ClusterRoleBinding (if is not present already)
	crb := clusterRoleBindingForNginxIngressController(clusterRoleName)

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
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

		err = r.client.Update(context.TODO(), crb)
		if err != nil {
			return err
		}
	}

	cm := configMapForNginxIngressController(instance)

	err = controllerutil.SetControllerReference(instance, cm, r.scheme)
	if err != nil {
		return err
	}

	err, existed = r.createIfNotExists(cm)
	if err != nil {
		return err
	}

	if !existed {
		reqLogger.Info("ConfigMap created", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
	}

	// IngressClass is available from k8s 1.18+
	minVersion, _ := version.ParseGeneric("v1.18.0")
	if RunningK8sVersion.AtLeast(minVersion) {
		if instance.Spec.IngressClass == "" {
			instance.Spec.IngressClass = "nginx"
			reqLogger.Info("Warning! IngressClass not set, using default", "IngressClass.Name", instance.Spec.IngressClass)
		}
		ic := ingressClassForNginxIngressController(instance)

		err, existed = r.createIfNotExists(ic)
		if err != nil {
			return err
		}

		if !existed {
			reqLogger.Info("IngressClass created", "IngressClass.Name", ic.Name)
		}
	}

	if instance.Spec.DefaultSecret == "" {
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, &v1.Secret{})

		if err != nil && errors.IsNotFound(err) {
			secret, err := defaultSecretForNginxIngressController(instance)
			if err != nil {
				return err
			}

			err = controllerutil.SetControllerReference(instance, secret, r.scheme)
			if err != nil {
				return err
			}

			err = r.client.Create(context.TODO(), secret)
			if err != nil {
				return err
			}

			reqLogger.Info("Warning! A custom self-signed TLS Secret has been created for the default server. "+
				"Update your 'DefaultSecret' with your own Secret in Production",
				"Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)

		} else if err != nil {
			return err
		}
	}

	if r.sccAPIExists {
		// Assign this new User to the SCC (if is not present already)
		scc := sccForNginxIngressController(sccName)

		err = r.client.Get(context.TODO(), types.NamespacedName{Name: sccName, Namespace: v1.NamespaceAll}, scc)
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

			err = r.client.Update(context.TODO(), scc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// blank assignment to verify that ReconcileNginxIngressController implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNginxIngressController{}

// ReconcileNginxIngressController reconciles a NginxIngressController object
type ReconcileNginxIngressController struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client       client.Client
	scheme       *runtime.Scheme
	sccAPIExists bool
}

// Reconcile reads that state of the cluster for a NginxIngressController object and makes changes based on the state read
// and what is in the NginxIngressController.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNginxIngressController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NginxIngressController")

	instance := &k8sv1alpha1.NginxIngressController{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil && errors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.GetDeletionTimestamp() != nil {
		err = r.handleDeletion(reqLogger, instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		reqLogger.Info("NginxIngressController was successfully deleted")
		return reconcile.Result{}, nil
	}

	err = r.ensureFinalizer(reqLogger, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Namespace could have been deleted in the middle of the reconcile
	ns := &v1.Namespace{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Namespace, Namespace: v1.NamespaceAll}, ns)
	if (err != nil && errors.IsNotFound(err)) || (ns.Status.Phase == "Terminating") {
		reqLogger.Info(fmt.Sprintf("The namespace '%v' does not exist or is in Terminating status, canceling Reconciling", instance.Namespace))
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to check if namespace exists")
		return reconcile.Result{}, err
	}

	err = r.checkPrerequisites(reqLogger, instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	if strings.ToLower(instance.Spec.Type) == "deployment" {
		dep := deploymentForNginxIngressController(instance)
		found := &appsv1.Deployment{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Deployment for NGINX Ingress Controller", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

			err = controllerutil.SetControllerReference(instance, dep, r.scheme)
			if err != nil {
				reqLogger.Error(err, "Error setting controller reference")
				return reconcile.Result{}, err
			}

			err = r.client.Create(context.TODO(), dep)
			if err != nil {
				reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return reconcile.Result{}, err
			}
		} else if err != nil {
			return reconcile.Result{}, err
		} else if hasDeploymentChanged(found, instance) {
			reqLogger.Info("NginxIngressController spec has changed, updating Deployment")
			updated := updateDeployment(found, instance)
			err = r.client.Update(context.TODO(), updated)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Remove possible DaemonSet
		ds := daemonSetForNginxIngressController(instance)
		err = r.deleteIfExists(ds.Name, ds.Namespace, ds)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if strings.ToLower(instance.Spec.Type) == "daemonset" {
		ds := daemonSetForNginxIngressController(instance)
		found := &appsv1.DaemonSet{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: ds.Name, Namespace: ds.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new DaemonSet for NGINX Ingress Controller", "DaemonSet.Namespace", ds.Namespace, "DaemonSet.Name", ds.Name)

			err = controllerutil.SetControllerReference(instance, ds, r.scheme)
			if err != nil {
				reqLogger.Error(err, "Error setting controller reference")
				return reconcile.Result{}, err
			}

			err = r.client.Create(context.TODO(), ds)
			if err != nil {
				reqLogger.Error(err, "Failed to create new DaemonSet", "DaemonSet.Namespace", ds.Namespace, "DaemonSet.Name", ds.Name)
				return reconcile.Result{}, err
			}
		} else if err != nil {
			return reconcile.Result{}, err
		} else if hasDaemonSetChanged(ds, instance) {
			reqLogger.Info("NginxIngressController spec has changed, updating DaemonSet")
			updated := updateDaemonSet(found, instance)
			err = r.client.Update(context.TODO(), updated)
			if err != nil {
				return reconcile.Result{}, err
			}
		}

		// Remove possible Deployment
		dep := deploymentForNginxIngressController(instance)
		err = r.deleteIfExists(dep.Name, dep.Namespace, dep)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	svc := serviceForNginxIngressController(instance)
	svcFound := &v1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service for NGINX Ingress Controller", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)

		err = controllerutil.SetControllerReference(instance, svc, r.scheme)
		if err != nil {
			reqLogger.Error(err, "Error setting controller reference")
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else if hasServiceChanged(svcFound, instance) {
		reqLogger.Info("NginxIngressController spec has changed, updating Service")
		updated := updateService(svcFound, instance)
		err = r.client.Update(context.TODO(), updated)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	cm := configMapForNginxIngressController(instance)
	cmFound := &v1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: cm.Namespace, Name: cm.Name}, cmFound)
	if err != nil {
		return reconcile.Result{}, err
	} else {
		if hasConfigMapChanged(cmFound, instance) {
			err = r.client.Update(context.TODO(), cm)
			reqLogger.Info("NginxIngressController spec has changed, updating ConfigMap")
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	if !instance.Status.Deployed {
		instance.Status.Deployed = true
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Finish reconcile for NginxIngressController")
	return reconcile.Result{}, nil
}

// createIfNotExists creates a new object. If the object exists, does nothing. It returns whether the object existed before or not.
func (r *ReconcileNginxIngressController) createIfNotExists(object runtime.Object) (error, bool) {
	err := r.client.Create(context.TODO(), object)
	if err != nil && errors.IsAlreadyExists(err) {
		return nil, true
	}

	return err, false
}

// deleteIfExists deletes an object if it exists in the cluster.
func (r *ReconcileNginxIngressController) deleteIfExists(name string, namespace string, object runtime.Object) error {
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, object)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err == nil {
		err := r.client.Delete(context.TODO(), object)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileNginxIngressController) finalizeNginxIngressController(reqLogger logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	crb := clusterRoleBindingForNginxIngressController(clusterRoleName)

	err := r.client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
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

	err = r.client.Update(context.TODO(), crb)
	if err != nil {
		return err
	}

	if r.sccAPIExists {
		scc := sccForNginxIngressController(sccName)

		err := r.client.Get(context.TODO(), types.NamespacedName{Name: sccName, Namespace: v1.NamespaceAll}, scc)
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

		err = r.client.Update(context.TODO(), scc)
		if err != nil {
			return err
		}
	}

	reqLogger.Info("Successfully finalized NginxIngressController")
	return nil
}

func (r *ReconcileNginxIngressController) addFinalizer(reqLogger logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	reqLogger.Info("Adding Finalizer for the NginxIngressController")
	instance.SetFinalizers(append(instance.GetFinalizers(), finalizer))

	err := r.client.Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Failed to update NginxIngressController with finalizer")
		return err
	}

	return nil
}

func (r *ReconcileNginxIngressController) handleDeletion(reqLogger logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	if !contains(instance.GetFinalizers(), finalizer) {
		return nil
	}

	err := r.finalizeNginxIngressController(reqLogger, instance)
	if err != nil {
		return err
	}

	instance.SetFinalizers(remove(instance.GetFinalizers(), finalizer))
	return r.client.Update(context.TODO(), instance)
}

func (r *ReconcileNginxIngressController) ensureFinalizer(reqLogger logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	if contains(instance.GetFinalizers(), finalizer) {
		return nil
	}

	return r.addFinalizer(reqLogger, instance)
}
