package controllers

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestDaemonSetForNginxIngressController(t *testing.T) {
	boolPointer := func(b bool) *bool { return &b }
	s := scheme.Scheme

	if err := k8sv1alpha1.AddToScheme(s); err != nil {
		t.Fatalf("Unable to add k8sv1alpha1 scheme: (%v)", err)
	}
	runAsUser := new(int64)
	allowPrivilegeEscalation := new(bool)
	*runAsUser = 101
	*allowPrivilegeEscalation = true

	instance := &k8sv1alpha1.NginxIngressController{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller",
		},
		Spec: k8sv1alpha1.NginxIngressControllerSpec{
			Image: k8sv1alpha1.Image{
				Repository: "nginx-ingress",
				Tag:        "edge",
			},
		},
	}
	expected := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller",
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
		Spec: appsv1.DaemonSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{"app": instance.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Name:      "my-nginx-ingress-controller",
					Namespace: "my-nginx-ingress-controller",
					Labels:    map[string]string{"app": instance.Name},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "my-nginx-ingress-controller",
					Containers: []corev1.Container{
						{
							Name:  "my-nginx-ingress-controller",
							Image: "nginx-ingress:edge",
							Args:  generatePodArgs(instance),
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
								},
								{
									Name:          "https",
									ContainerPort: 443,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
									Add:  []corev1.Capability{"NET_BIND_SERVICE"},
								},
								RunAsUser:                runAsUser,
								AllowPrivilegeEscalation: allowPrivilegeEscalation,
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	result, _ := daemonSetForNginxIngressController(instance, s)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("daemonSetForNginxIngressController() mismatch (-want +got):\n%s", diff)
	}
}
