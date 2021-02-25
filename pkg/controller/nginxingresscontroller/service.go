package nginxingresscontroller

import (
	"reflect"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func serviceForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *corev1.Service {
	extraLabels := map[string]string{}
	if instance.Spec.Service != nil {
		extraLabels = instance.Spec.Service.ExtraLabels
	}

	return &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
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
			Type:     corev1.ServiceType(instance.Spec.ServiceType),
		},
	}
}

func hasServiceChanged(svc *corev1.Service, instance *k8sv1alpha1.NginxIngressController) bool {
	if svc.Spec.Type != corev1.ServiceType(instance.Spec.ServiceType) {
		return true
	}
	if instance.Spec.Service != nil && !reflect.DeepEqual(svc.Labels, instance.Spec.Service.ExtraLabels) {
		return true
	}
	return svc.Labels != nil && instance.Spec.Service == nil
}

func updateService(svc *corev1.Service, instance *k8sv1alpha1.NginxIngressController) *corev1.Service {
	svc.Spec.Type = corev1.ServiceType(instance.Spec.ServiceType)
	if instance.Spec.Service != nil {
		svc.Labels = instance.Spec.Service.ExtraLabels
	} else {
		svc.Labels = nil
	}
	return svc
}
