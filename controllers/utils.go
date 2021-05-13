package controllers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-logr/logr"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const apiVersionUnsupportedError = "server does not support API version"

// RunningK8sVersion contains the version of k8s
var RunningK8sVersion *version.Version

// generatePodArgs generate a list of arguments for the Ingress Controller pods based on the CRD.
func generatePodArgs(instance *k8sv1alpha1.NginxIngressController) []string {
	var args []string

	args = append(args, fmt.Sprintf("-nginx-configmaps=%v/%v", instance.Namespace, instance.Name))

	defaultSecretName := instance.Spec.DefaultSecret
	if instance.Spec.DefaultSecret == "" {
		defaultSecretName = fmt.Sprintf("%v/%v", instance.Namespace, instance.Name)
	}
	args = append(args, fmt.Sprintf("-default-server-tls-secret=%v", defaultSecretName))

	if instance.Spec.NginxPlus {
		args = append(args, "-nginx-plus")

		if instance.Spec.AppProtect != nil && instance.Spec.AppProtect.Enable {
			args = append(args, "-enable-app-protect")
		}
	}

	if instance.Spec.IngressClass != "" {
		args = append(args, fmt.Sprintf("-ingress-class=%v", instance.Spec.IngressClass))
	}

	if instance.Spec.UseIngressClassOnly {
		args = append(args, "-use-ingress-class-only")
	}

	if instance.Spec.WatchNamespace != "" {
		args = append(args, fmt.Sprintf("-watch-namespace=%v", instance.Spec.WatchNamespace))
	}

	if instance.Spec.HealthStatus != nil && instance.Spec.HealthStatus.Enable {
		args = append(args, "-health-status")
		if instance.Spec.HealthStatus.URI != "" {
			args = append(args, fmt.Sprintf("-health-status-uri=%v", instance.Spec.HealthStatus.URI))
		}
	}

	if instance.Spec.NginxDebug {
		args = append(args, "-nginx-debug")
	}

	if instance.Spec.LogLevel != 0 {
		args = append(args, fmt.Sprintf("-v=%v", instance.Spec.LogLevel))
	}

	if instance.Spec.NginxStatus != nil && instance.Spec.NginxStatus.Enable {
		args = append(args, "-nginx-status")

		if instance.Spec.NginxStatus.Port != nil {
			args = append(args, fmt.Sprintf("-nginx-status-port=%v", *instance.Spec.NginxStatus.Port))
		}

		if instance.Spec.NginxStatus.AllowCidrs != "" {
			args = append(args, fmt.Sprintf("-nginx-status-allow-cidrs=%v", instance.Spec.NginxStatus.AllowCidrs))
		}
	}

	if instance.Spec.ReportIngressStatus != nil && instance.Spec.ReportIngressStatus.Enable {
		args = append(args, "-report-ingress-status")

		if instance.Spec.ReportIngressStatus.ExternalService != "" {
			args = append(args, fmt.Sprintf("-external-service=%v", instance.Spec.ReportIngressStatus.ExternalService))
		} else if instance.Spec.ServiceType == "LoadBalancer" {
			args = append(args, fmt.Sprintf("-external-service=%v", instance.Name))
		} else if instance.Spec.ReportIngressStatus.IngressLink != "" {
			args = append(args, fmt.Sprintf("-ingresslink=%v", instance.Spec.ReportIngressStatus.IngressLink))
		}
	}

	if instance.Spec.EnableLeaderElection == nil || *instance.Spec.EnableLeaderElection {
		args = append(args, fmt.Sprintf("-leader-election-lock-name=%v-lock", instance.Name))
	} else {
		args = append(args, "-enable-leader-election=false")
	}

	if instance.Spec.WildcardTLS != "" {
		args = append(args, fmt.Sprintf("-wildcard-tls-secret=%v", instance.Spec.WildcardTLS))
	}

	if instance.Spec.Prometheus != nil && instance.Spec.Prometheus.Enable {
		args = append(args, "-enable-prometheus-metrics")

		if instance.Spec.Prometheus.Port != nil {
			args = append(args, fmt.Sprintf("-prometheus-metrics-listen-port=%v", *instance.Spec.Prometheus.Port))
		}

		if instance.Spec.EnableLatencyMetrics {
			args = append(args, "-enable-latency-metrics")
		}
	}

	if instance.Spec.EnableCRDs != nil && !*instance.Spec.EnableCRDs {
		args = append(args, "-enable-custom-resources=false")
	} else {
		if instance.Spec.EnableTLSPassthrough {
			args = append(args, "-enable-tls-passthrough")
		}

		if instance.Spec.GlobalConfiguration != "" {
			args = append(args, fmt.Sprintf("-global-configuration=%v", instance.Spec.GlobalConfiguration))
		}

		if instance.Spec.EnableSnippets {
			args = append(args, "-enable-snippets")
		}

		if instance.Spec.EnablePreviewPolicies {
			args = append(args, "-enable-preview-policies")
		}
	}

	if instance.Spec.NginxReloadTimeout != 0 {
		args = append(args, fmt.Sprintf("-nginx-reload-timeout=%v", instance.Spec.NginxReloadTimeout))
	}

	return args
}

// hasDifferentArguments returns whether the arguments of a container are different than the NginxIngressController spec.
func hasDifferentArguments(container corev1.Container, instance *k8sv1alpha1.NginxIngressController) bool {
	newArgs := generatePodArgs(instance)
	return !reflect.DeepEqual(newArgs, container.Args)
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func VerifySCCAPIExists() (bool, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return false, err
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return false, err
	}

	gv := schema.GroupVersion{
		Group:   secv1.GroupName,
		Version: secv1.GroupVersion.Version,
	}

	err = discovery.ServerSupportsVersion(clientSet, gv)
	if err != nil {
		// This error means the call could not find SCC in the API, but there was no API error.
		if strings.Contains(err.Error(), apiVersionUnsupportedError) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func generateImage(repository string, tag string) string {
	return fmt.Sprintf("%v:%v", repository, tag)
}

// GetK8sVersion returns the running version of k8s
func GetK8sVersion(client kubernetes.Interface) (v *version.Version, err error) {
	serverVersion, err := client.Discovery().ServerVersion()
	if err != nil {
		return nil, err
	}

	runningVersion, err := version.ParseGeneric(serverVersion.String())
	if err != nil {
		return nil, fmt.Errorf("unexpected error parsing running Kubernetes version: %v", err)
	}

	return runningVersion, nil
}

// checkPrerequisites creates all necessary objects before the deployment of a new Ingress Controller.
func (r *NginxIngressControllerReconciler) checkPrerequisites(log logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	sa := r.serviceAccountForNginxIngressController(instance)
	err, existed := r.createIfNotExists(sa)
	if err != nil {
		return err
	}

	if !existed {
		log.Info("ServiceAccount created", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
	}

	// Assign this new ServiceAccount to the ClusterRoleBinding (if is not present already)
	crb := r.clusterRoleBindingForNginxIngressController(clusterRoleName)

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

	// IngressClass is available from k8s 1.18+
	minVersion, _ := version.ParseGeneric("v1.18.0")
	if RunningK8sVersion.AtLeast(minVersion) {
		if instance.Spec.IngressClass == "" {
			instance.Spec.IngressClass = "nginx"
			log.Info("Warning! IngressClass not set, using default", "IngressClass.Name", instance.Spec.IngressClass)
		}
		ic := r.ingressClassForNginxIngressController(instance)

		err, existed = r.createIfNotExists(ic)
		if err != nil {
			return err
		}

		if !existed {
			log.Info("IngressClass created", "IngressClass.Name", ic.Name)
		}
	}

	if instance.Spec.DefaultSecret == "" {
		err = r.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, &v1.Secret{})

		if err != nil && errors.IsNotFound(err) {
			secret, err := r.defaultSecretForNginxIngressController(instance)
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

	if r.SccAPIExists {
		// Assign this new User to the SCC (if is not present already)
		scc := r.sccForNginxIngressController(sccName)

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

// create common resources shared by all the Ingress Controllers
func (r *NginxIngressControllerReconciler) createCommonResources(log logr.Logger, instance *k8sv1alpha1.NginxIngressController) error {
	// Create ClusterRole and ClusterRoleBinding for all the NginxIngressController resources.
	var err error

	cr := r.clusterRoleForNginxIngressController(clusterRoleName)

	err = r.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, cr)

	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("no previous ClusterRole found, creating a new one.")
			err = r.Create(context.TODO(), cr)
			if err != nil {
				return fmt.Errorf("error creating ClusterRole: %v", err)
			}
		} else {
			return fmt.Errorf("error getting ClusterRole: %v", err)
		}
	} else {
		// For updates in the ClusterRole permissions (eg new CRDs of the Ingress Controller).
		log.Info("previous ClusterRole found, updating.")
		cr := r.clusterRoleForNginxIngressController(clusterRoleName)
		err = r.Update(context.TODO(), cr)
		if err != nil {
			return fmt.Errorf("error updating ClusterRole: %v", err)
		}
	}

	crb := r.clusterRoleBindingForNginxIngressController(clusterRoleName)

	err = r.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName, Namespace: v1.NamespaceAll}, crb)
	if err != nil && errors.IsNotFound(err) {
		log.Info("no previous ClusterRoleBinding found, creating a new one.")
		err = r.Create(context.TODO(), crb)
	}

	if err != nil {
		return fmt.Errorf("error creating ClusterRoleBinding: %v", err)
	}

	err = createKICCustomResourceDefinitions(log, r.Mgr)
	if err != nil {
		return fmt.Errorf("error creating KIC CRDs: %v", err)
	}

	if r.SccAPIExists {
		log.Info("OpenShift detected as platform.")

		scc := r.sccForNginxIngressController(sccName)

		err = r.Get(context.TODO(), types.NamespacedName{Name: sccName, Namespace: v1.NamespaceAll}, scc)
		if err != nil && errors.IsNotFound(err) {
			log.Info("no previous SecurityContextConstraints found, creating a new one.")
			err = r.Create(context.TODO(), scc)
		}

		if err != nil {
			return fmt.Errorf("error creating SecurityContextConstraints: %v", err)
		}
	}

	return nil
}
