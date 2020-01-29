package nginxingresscontroller

import (
	"reflect"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"
)

func configMapForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Data: instance.Spec.ConfigMapData,
	}
}

func hasConfigMapChanged(cm *v1.ConfigMap, instance *k8sv1alpha1.NginxIngressController) bool {
	return !reflect.DeepEqual(cm.Data, instance.Spec.ConfigMapData)
}
