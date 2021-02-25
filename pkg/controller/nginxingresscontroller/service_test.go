package nginxingresscontroller

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestServiceForNginxIngressController(t *testing.T) {
	name := "my-service"
	namespace := "my-nginx-ingress"
	extraLabels := map[string]string{"app": "my-nginx-ingress"}

	instance := &k8sv1alpha1.NginxIngressController{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    extraLabels,
		},
		Spec: k8sv1alpha1.NginxIngressControllerSpec{
			ServiceType: string(corev1.ServiceTypeLoadBalancer),
			Service: &k8sv1alpha1.Service{
				ExtraLabels: extraLabels,
			},
		},
	}
	expected := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    extraLabels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: "TCP",
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: 80,
					},
				},
				{
					Name:     "https",
					Protocol: "TCP",
					Port:     443,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: 443,
					},
				},
			},
			Selector: map[string]string{"app": instance.Name},
			Type:     corev1.ServiceTypeLoadBalancer,
		},
	}

	result := serviceForNginxIngressController(instance)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("serviceForNginxIngressController() mismatch (-want +got):\n%s", diff)
	}
}

func TestHasServiceChanged(t *testing.T) {
	name := "my-service"
	namespace := "my-nginx-ingress"
	extraLabels := map[string]string{"app": "my-nginx-ingress"}

	tests := []struct {
		svc      *corev1.Service
		instance *k8sv1alpha1.NginxIngressController
		expected bool
		msg      string
	}{
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: string(corev1.ServiceTypeLoadBalancer),
				},
			},
			false,
			"no changes",
		},
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: "NodePort",
				},
			},
			true,
			"different service type",
		},
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    extraLabels,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: string(corev1.ServiceTypeLoadBalancer),
					Service: &k8sv1alpha1.Service{
						ExtraLabels: map[string]string{"new": "label"},
					},
				},
			},
			true,
			"different label",
		},
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    extraLabels,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: string(corev1.ServiceTypeLoadBalancer),
					Service: &k8sv1alpha1.Service{
						ExtraLabels: nil,
					},
				},
			},
			true,
			"remove label",
		},
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    extraLabels,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: string(corev1.ServiceTypeLoadBalancer),
				},
			},
			true,
			"remove service parameters",
		},
	}
	for _, test := range tests {
		result := hasServiceChanged(test.svc, test.instance)
		if result != test.expected {
			t.Errorf("hasServiceChanged() = %v, want %v for the case of %v", result, test.expected, test.msg)
		}
	}
}

func TestUpdateService(t *testing.T) {
	name := "my-service"
	namespace := "my-nginx-ingress"
	extraLabels := map[string]string{"app": "my-nginx-ingress"}

	tests := []struct {
		svc      *corev1.Service
		instance *k8sv1alpha1.NginxIngressController
		expected *corev1.Service
		msg      string
	}{
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: string(corev1.ServiceTypeNodePort),
				},
			},
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeNodePort,
				},
			},
			"override service type",
		},
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    map[string]string{"my": "labelToBeOverridden"},
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: string(corev1.ServiceTypeLoadBalancer),
					Service: &k8sv1alpha1.Service{
						ExtraLabels: extraLabels,
					},
				},
			},
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    extraLabels,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			"override labels",
		},
		{
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    map[string]string{"my": "labelToBeRemoved"},
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			&k8sv1alpha1.NginxIngressController{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					ServiceType: string(corev1.ServiceTypeLoadBalancer),
				},
			},
			&corev1.Service{
				ObjectMeta: v1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
				Spec: corev1.ServiceSpec{
					Type: corev1.ServiceTypeLoadBalancer,
				},
			},
			"remove labels",
		},
	}
	for _, test := range tests {
		result := updateService(test.svc, test.instance)
		if diff := cmp.Diff(test.expected, result); diff != "" {
			t.Errorf("updateService() mismatch for the case of %v (-want +got):\n%s", test.msg, diff)
		}
	}
}
