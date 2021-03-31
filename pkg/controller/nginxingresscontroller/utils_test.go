package nginxingresscontroller

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGeneratePodArgs(t *testing.T) {
	var promPort, statusPort uint16
	promPort = 9114
	statusPort = 9090
	name := "my-nginx-ingress"
	namespace := "my-nginx-ingress"
	enable := true
	disable := false
	tests := []struct {
		instance *k8sv1alpha1.NginxIngressController
		expected []string
	}{
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-nginx-ingress",
				"-leader-election-lock-name=my-nginx-ingress-lock",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					DefaultSecret: "my-nginx-ingress/my-secret",
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-secret",
				"-leader-election-lock-name=my-nginx-ingress-lock",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					NginxPlus: true,
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-nginx-ingress",
				"-nginx-plus",
				"-leader-election-lock-name=my-nginx-ingress-lock",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					EnableCRDs: &disable,
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-nginx-ingress",
				"-leader-election-lock-name=my-nginx-ingress-lock",
				"-enable-custom-resources=false",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					NginxPlus:     true,
					EnableCRDs:    &disable,
					DefaultSecret: "my-nginx-ingress/my-secret",
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-secret",
				"-nginx-plus",
				"-leader-election-lock-name=my-nginx-ingress-lock",
				"-enable-custom-resources=false",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					DefaultSecret: "my-nginx-ingress/my-secret",
					ServiceType:   "NodePort",
					ReportIngressStatus: &k8sv1alpha1.ReportIngressStatus{
						Enable:      true,
						IngressLink: "my-ingresslink",
					},
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-secret",
				"-report-ingress-status",
				"-ingresslink=my-ingresslink",
				"-leader-election-lock-name=my-nginx-ingress-lock",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					DefaultSecret: "my-nginx-ingress/my-secret",
					ServiceType:   "LoadBalancer",
					ReportIngressStatus: &k8sv1alpha1.ReportIngressStatus{
						Enable:      true,
						IngressLink: "my-invalid-ingresslink",
					},
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-secret",
				"-report-ingress-status",
				fmt.Sprintf("-external-service=%v", name),
				"-leader-election-lock-name=my-nginx-ingress-lock",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					EnableCRDs:            &enable,
					EnableSnippets:        true,
					EnablePreviewPolicies: true,
					EnableTLSPassthrough:  true,
					GlobalConfiguration:   "my-nginx-ingress/globalconfiguration",
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-nginx-ingress",
				"-leader-election-lock-name=my-nginx-ingress-lock",
				"-enable-tls-passthrough",
				"-global-configuration=my-nginx-ingress/globalconfiguration",
				"-enable-snippets",
				"-enable-preview-policies",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					EnableLeaderElection: &disable,
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-nginx-ingress",
				"-enable-leader-election=false",
			},
		},
		{
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					NginxPlus:           true,
					DefaultSecret:       "my-nginx-ingress/my-secret",
					IngressClass:        "ingressClass",
					UseIngressClassOnly: true,
					WatchNamespace:      "default",
					HealthStatus: &k8sv1alpha1.HealthStatus{
						Enable: true,
						URI:    "/healthz",
					},
					NginxDebug: true,
					LogLevel:   3,
					NginxStatus: &k8sv1alpha1.NginxStatus{
						Enable:     true,
						Port:       &statusPort,
						AllowCidrs: "127.0.0.1",
					},
					ReportIngressStatus: &k8sv1alpha1.ReportIngressStatus{
						Enable:          true,
						ExternalService: "external",
						IngressLink:     "my-invalid-ingressLink",
					},
					EnableLeaderElection: &enable,
					WildcardTLS:          "my-nginx-ingress/wildcard-secret",
					Prometheus: &k8sv1alpha1.Prometheus{
						Enable: true,
						Port:   &promPort,
					},
					EnableLatencyMetrics: true,
					GlobalConfiguration:  "my-nginx-ingress/globalconfiguration",
					EnableTLSPassthrough: true,
					AppProtect: &k8sv1alpha1.AppProtect{
						Enable: true,
					},
					NginxReloadTimeout:    5000,
					EnableCRDs:            &disable,
					EnableSnippets:        true,
					EnablePreviewPolicies: true,
				},
			},
			expected: []string{
				"-nginx-configmaps=my-nginx-ingress/my-nginx-ingress",
				"-default-server-tls-secret=my-nginx-ingress/my-secret",
				"-nginx-plus",
				"-enable-app-protect",
				"-ingress-class=ingressClass",
				"-use-ingress-class-only",
				"-watch-namespace=default",
				"-health-status",
				"-health-status-uri=/healthz",
				"-nginx-debug",
				"-v=3",
				"-nginx-status",
				"-nginx-status-port=9090",
				"-nginx-status-allow-cidrs=127.0.0.1",
				"-report-ingress-status",
				"-external-service=external",
				"-leader-election-lock-name=my-nginx-ingress-lock",
				"-wildcard-tls-secret=my-nginx-ingress/wildcard-secret",
				"-enable-prometheus-metrics",
				"-prometheus-metrics-listen-port=9114",
				"-enable-latency-metrics",
				"-enable-custom-resources=false",
				"-nginx-reload-timeout=5000",
			},
		},
	}

	for _, test := range tests {
		result := generatePodArgs(test.instance)
		if diff := cmp.Diff(test.expected, result); diff != "" {
			t.Errorf("generatePodArgs(%+v) mismatch (-want +got):\n%s", test.instance, diff)
		}
	}
}

func TestHasDifferentArguments(t *testing.T) {
	name := "my-nginx-ingress"
	namespace := "my-nginx-ingress"
	tests := []struct {
		container corev1.Container
		instance  *k8sv1alpha1.NginxIngressController
		expected  bool
	}{
		{
			container: corev1.Container{
				Args: []string{
					fmt.Sprintf("-nginx-configmaps=%v/%v", namespace, name),
					fmt.Sprintf("-default-server-tls-secret=%v/%v", namespace, name),
					"-nginx-plus",
					"-leader-election-lock-name=my-nginx-ingress-lock",
				},
			},
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					NginxPlus: true,
				},
			},
			expected: false,
		},
		{
			container: corev1.Container{
				Args: []string{
					fmt.Sprintf("-nginx-configmaps=%v/%v", namespace, name),
					fmt.Sprintf("-default-server-tls-secret=%v/%v", namespace, name),
					"-nginx-plus=false",
					"-leader-election-lock-name=my-nginx-ingress-lock",
				},
			},
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					NginxPlus: true,
				},
			},
			expected: true,
		},
		{
			container: corev1.Container{
				Args: []string{
					fmt.Sprintf("-nginx-configmaps=%v/%v", namespace, name),
					"-default-server-tls-secret=default/mysecret",
					"-nginx-plus",
					"-leader-election-lock-name=my-nginx-ingress-lock",
				},
			},
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					NginxPlus:     true,
					DefaultSecret: "my-namespace/my-secret",
				},
			},
			expected: true,
		},
		{
			container: corev1.Container{
				Args: []string{fmt.Sprintf(
					"-nginx-configmaps=%v/%v", namespace, name),
					"-default-server-tls-secret=default/mysecret",
					"-nginx-plus",
					"-leader-election-lock-name=my-nginx-ingress-lock",
					"-enable-custom-resources=false",
				},
			},
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					NginxPlus:     true,
					DefaultSecret: "my-namespace/my-secret",
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		result := hasDifferentArguments(test.container, test.instance)
		if diff := cmp.Diff(test.expected, result); diff != "" {
			t.Errorf("hasDifferentArguments(%+v, %+v) mismatch (-want +got):\n%s", test.container, test.instance, diff)
		}
	}
}

func TestGenerateImage(t *testing.T) {
	rep := "repository/image"
	version := "version"
	expected := "repository/image:version"
	result := generateImage(rep, version)

	if expected != result {
		t.Errorf("generateImage(%v, %v) returned %v but expected %v", rep, version, result, expected)
	}
}
