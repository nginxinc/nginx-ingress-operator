package controllers

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	v1 "k8s.io/api/core/v1"
)

func (r *NginxIngressControllerReconciler) configMapForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) (*v1.ConfigMap, error) {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Data: instance.Spec.ConfigMapData,
	}
	if err := ctrl.SetControllerReference(instance, cm, r.Scheme); err != nil {
		return nil, err
	}
	return cm, nil
}

func configMapMutateFn(cm *v1.ConfigMap, configMapData map[string]string) controllerutil.MutateFn {
	return func() error {
		cm.Data = configMapData
		return nil
	}
}
