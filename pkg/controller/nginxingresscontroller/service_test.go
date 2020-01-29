package nginxingresscontroller

import (
	"reflect"
	"testing"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestServiceForNginxIngressController(t *testing.T) {
	name := "my-service"
	namespace := "my-nginx-ingress"
	serviceType := "LoadBalancer"

	instance := &k8sv1alpha1.NginxIngressController{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: k8sv1alpha1.NginxIngressControllerSpec{
			ServiceType: serviceType,
		},
	}
	expected := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
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
			Type:     corev1.ServiceType(serviceType),
		},
	}

	result := serviceForNginxIngressController(instance)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("serviceForNginxIngressController(%v) returned %+v but expected %+v", instance, result, expected)
	}
}
