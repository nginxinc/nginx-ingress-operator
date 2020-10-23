package nginxingresscontroller

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	networking "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIngressClassForNginxIngressController(t *testing.T) {
	instance := &k8sv1alpha1.NginxIngressController{
		Spec: k8sv1alpha1.NginxIngressControllerSpec{
			IngressClass: "nginx",
		},
	}
	expected := &networking.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: networking.IngressClassSpec{
			Controller: "nginx.org/ingress-controller",
		},
	}

	result := ingressClassForNginxIngressController(instance)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("ingressClassForNginxIngressController() mismatch (-want +got):\n%s", diff)
	}
}
