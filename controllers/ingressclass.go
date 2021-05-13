package controllers

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *NginxIngressControllerReconciler) ingressClassForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *networking.IngressClass {
	ic := &networking.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.Spec.IngressClass,
		},
		Spec: networking.IngressClassSpec{
			Controller: "nginx.org/ingress-controller",
		},
	}
	return ic
}
