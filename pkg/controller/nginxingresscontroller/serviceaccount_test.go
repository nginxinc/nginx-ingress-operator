package nginxingresscontroller

import (
	"reflect"
	"testing"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestServiceAccountForNginxIngressController(t *testing.T) {
	namespace := "my-nginx-ingress"
	name := "my-sa"
	instance := &k8sv1alpha1.NginxIngressController{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	expected := &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	result := serviceAccountForNginxIngressController(instance)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("serviceAccountForNginxIngressController(%v) returned %+v but expected %+v", instance, result, expected)
	}
}
