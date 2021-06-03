package controllers

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func serviceAccountForNginxIngressController(instance *k8sv1alpha1.NginxIngressController, scheme *runtime.Scheme) (*corev1.ServiceAccount, error) {
	svca := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	if err := ctrl.SetControllerReference(instance, svca, scheme); err != nil {
		return nil, err
	}

	return svca, nil
}
