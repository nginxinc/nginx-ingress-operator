package controllers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestServiceForNginxIngressController(t *testing.T) {
	boolPointer := func(b bool) *bool { return &b }
	s := scheme.Scheme

	if err := k8sv1alpha1.AddToScheme(s); err != nil {
		t.Fatalf("Unable to add k8sv1alpha1 scheme: (%v)", err)
	}
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
			OwnerReferences: []v1.OwnerReference{
				{
					APIVersion:         "k8s.nginx.org/v1alpha1",
					Name:               instance.Name,
					Kind:               "NginxIngressController",
					Controller:         boolPointer(true),
					BlockOwnerDeletion: boolPointer(true),
				},
			},
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

	result, _ := serviceForNginxIngressController(instance, s)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("serviceForNginxIngressController() mismatch (-want +got):\n%s", diff)
	}
}
