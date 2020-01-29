package nginxingresscontroller

import (
	"reflect"
	"testing"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigMapForNginxIngressController(t *testing.T) {
	instance := &k8sv1alpha1.NginxIngressController{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller",
		},
	}
	expected := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller",
		},
	}

	result := configMapForNginxIngressController(instance)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("configMapForNginxIngressController(%+v) returned %+v but expected %+v", instance, result, expected)
	}
}

func TestHasConfigMapChanged(t *testing.T) {
	cm1 := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller",
		},
	}

	instance := &k8sv1alpha1.NginxIngressController{Spec: k8sv1alpha1.NginxIngressControllerSpec{}}
	result := hasConfigMapChanged(cm1, instance)
	if result != false {
		t.Errorf("hasConfigMapChanged(%v, %v) returned %v but expected %v for the case of same configmaps.", cm1, instance, result, false)
	}

	instance.Spec.ConfigMapData = map[string]string{"key": "value"}
	result = hasConfigMapChanged(cm1, instance)
	if result != true {
		t.Errorf("hasConfigMapChanged(%v, %v) returned %v but expected %v for the case of different configmaps.", cm1, instance, result, true)
	}
}
