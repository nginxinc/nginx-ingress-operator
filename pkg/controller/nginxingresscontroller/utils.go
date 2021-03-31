package nginxingresscontroller

import (
	"fmt"
	"reflect"
	"strings"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
