package controllers

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ingressClassForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *networking.IngressClass {
	ingressClassName := "nginx"
	if instance.Spec.IngressClass != "" {
		ingressClassName = instance.Spec.IngressClass
	}
	ic := &networking.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressClassName,
		},
		Spec: networking.IngressClassSpec{
			Controller: "nginx.org/ingress-controller",
		},
	}
	return ic
}
