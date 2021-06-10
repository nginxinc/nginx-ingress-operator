package controllers

import (
	"reflect"
	"testing"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestDefaultSecretForNginxIngressController(t *testing.T) {
	s := scheme.Scheme

	if err := k8sv1alpha1.AddToScheme(s); err != nil {
		t.Fatalf("Unable to add k8sv1alpha1 scheme: (%v)", err)
	}
	instance := &k8sv1alpha1.NginxIngressController{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller-ns",
		},
	}

	expectedObjectMeta := &metav1.ObjectMeta{
		Name:      "my-nginx-ingress-controller",
		Namespace: "my-nginx-ingress-controller-ns",
	}
	expectedType := corev1.SecretTypeTLS

	secret, err := defaultSecretForNginxIngressController(instance, s)
	if err != nil {
		t.Fatalf("defaultSecretForNginxIngressController() returned unexpected error %v", err)
	}

	if reflect.DeepEqual(expectedObjectMeta, secret.ObjectMeta) {
		t.Errorf("defaultSecretForNginxIngressController() returned %v but expected %v", secret.ObjectMeta, expectedObjectMeta)
	}
	if expectedType != secret.Type {
		t.Errorf("defaultSecretForNginxIngressController() returned %s but expected %s", secret.Type, expectedType)
	}
	if len(secret.Data[corev1.TLSCertKey]) == 0 {
		t.Errorf("defaultSecretForNginxIngressController() returned empty data key %s", corev1.TLSCertKey)
	}
	if len(secret.Data[corev1.TLSPrivateKeyKey]) == 0 {
		t.Errorf("defaultSecretForNginxIngressController() returned empty data key %s", corev1.TLSPrivateKeyKey)
	}
}
