package nginxingresscontroller

import (
	"reflect"
	"testing"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestVsForNginxIngressController(t *testing.T) {
	expected := &v1beta1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: "virtualservers.k8s.nginx.org",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group: "k8s.nginx.org",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural:     "virtualservers",
				Singular:   "virtualserver",
				ShortNames: []string{"vs"},
				Kind:       "VirtualServer",
			},
			Scope: "Namespaced",
			Versions: []v1beta1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
				},
			},
		},
	}

	result := vsForNginxIngressController()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("vsForNginxIngressController() returned %+v but expected %+v", result, expected)
	}
}

func TestVsrForNginxIngressController(t *testing.T) {
	expected := &v1beta1.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: "virtualserverroutes.k8s.nginx.org",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group: "k8s.nginx.org",
			Names: v1beta1.CustomResourceDefinitionNames{
				Plural:     "virtualserverroutes",
				Singular:   "virtualserveroute",
				ShortNames: []string{"vsr"},
				Kind:       "VirtualServerRoute",
			},
			Scope: "Namespaced",
			Versions: []v1beta1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
				},
			},
		},
	}

	result := vsrForNginxIngressController()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("vsrForNginxIngressController() returned %+v but expected %+v", result, expected)
	}
}
