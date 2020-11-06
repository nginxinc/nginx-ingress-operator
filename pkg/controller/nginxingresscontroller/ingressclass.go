package nginxingresscontroller

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ingressClassForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *networking.IngressClass {
	return &networking.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.Spec.IngressClass,
		},
		Spec: networking.IngressClassSpec{
			Controller: "nginx.org/ingress-controller",
		},
	}
}
