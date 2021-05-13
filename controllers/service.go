package controllers

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *NginxIngressControllerReconciler) serviceForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *corev1.Service {
	extraLabels := map[string]string{}
	if instance.Spec.Service != nil {
		extraLabels = instance.Spec.Service.ExtraLabels
	}

	svc := &corev1.Service{
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
	ctrl.SetControllerReference(instance, svc, r.Scheme)

	return svc
}

func serviceMutateFn(svc *corev1.Service, serviceType string, labels map[string]string) controllerutil.MutateFn {
	return func() error {
		svc.Spec.Type = corev1.ServiceType(serviceType)
		svc.Labels = labels
		return nil
	}
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
