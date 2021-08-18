package controllers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestServiceAccountForNginxIngressController(t *testing.T) {
	boolPointer := func(b bool) *bool { return &b }
	s := scheme.Scheme

	if err := k8sv1alpha1.AddToScheme(s); err != nil {
		t.Fatalf("Unable to add k8sv1alpha1 scheme: (%v)", err)
	}
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
			OwnerReferences: []v1.OwnerReference{
				{
					APIVersion:         "k8s.nginx.org/v1alpha1",
					Name:               instance.Name,
					Kind:               "NginxIngressController",
					Controller:         boolPointer(true),
					BlockOwnerDeletion: boolPointer(true),
				},
			},
		},
	}

	result, _ := serviceAccountForNginxIngressController(instance, s)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("serviceAccountForNginxIngressController() mismatch (-want +got):\n%s", diff)
	}
}
