package controllers

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *NginxIngressControllerReconciler) serviceAccountForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *corev1.ServiceAccount {
	svca := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}
	ctrl.SetControllerReference(instance, svca, r.Scheme)
	return svca
}
