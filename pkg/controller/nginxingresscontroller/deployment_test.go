package nginxingresscontroller

import (
	"reflect"
	"testing"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDeploymentForNginxIngressController(t *testing.T) {
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
	expected := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller",
		},
		Spec: appsv1.DeploymentSpec{
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

	result := deploymentForNginxIngressController(instance)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("deploymentForNginxIngressController(%+v) returned %+v but expected %+v", instance, result, expected)
	}
}

func TestHasDeploymentChanged(t *testing.T) {
	runAsUser := new(int64)
	allowPrivilegeEscalation := new(bool)
	*runAsUser = 101
	*allowPrivilegeEscalation = true
	replicas := new(int32)
	*replicas = 1

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
			Replicas: replicas,
		},
	}

	defaultDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-nginx-ingress-controller",
			Namespace: "my-nginx-ingress-controller",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Name:      "my-nginx-ingress-controller",
					Namespace: "my-nginx-ingress-controller",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "my-nginx-ingress-controller",
							Image: "nginx-ingress:edge",
							Args:  generatePodArgs(instance),
						},
					},
				},
			},
		},
	}

	tenReplicas := int32(10)

	tests := []struct {
		deployment *appsv1.Deployment
		instance   *k8sv1alpha1.NginxIngressController
		expected   bool
		msg        string
	}{
		{
			deployment: defaultDeployment,
			instance:   instance,
			expected:   false,
			msg:        "no changes",
		},
		{
			deployment: defaultDeployment,
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-nginx-ingress-controller",
					Namespace: "my-nginx-ingress-controller",
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					Image: k8sv1alpha1.Image{
						Repository: "nginx-ingress",
						Tag:        "edge",
					},
					Replicas: &tenReplicas,
				},
			},
			expected: true,
			msg:      "replicas increased",
		},
		{
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-nginx-ingress-controller",
					Namespace: "my-nginx-ingress-controller",
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &tenReplicas, // Deployment with 10 replicas
					Template: corev1.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Name:      "my-nginx-ingress-controller",
							Namespace: "my-nginx-ingress-controller",
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "my-nginx-ingress-controller",
									Image: "nginx-ingress:edge",
									Args:  generatePodArgs(instance),
								},
							},
						},
					},
				},
			},
			instance: &k8sv1alpha1.NginxIngressController{
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
			},
			expected: true,
			msg:      "replicas field removed",
		},
		{
			deployment: defaultDeployment,
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-nginx-ingress-controller",
					Namespace: "my-nginx-ingress-controller",
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					Image: k8sv1alpha1.Image{
						Repository: "nginx-plus-ingress",
						Tag:        "edge",
					},
				},
			},
			expected: true,
			msg:      "image repository update",
		},
		{
			deployment: defaultDeployment,
			instance: &k8sv1alpha1.NginxIngressController{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-nginx-ingress-controller",
					Namespace: "my-nginx-ingress-controller",
				},
				Spec: k8sv1alpha1.NginxIngressControllerSpec{
					Image: k8sv1alpha1.Image{
						Repository: "nginx-ingress",
						Tag:        "edge",
						PullPolicy: "Always",
					},
				},
			},
			expected: true,
			msg:      "pull policy update",
		},
	}
	for _, test := range tests {
		result := hasDeploymentChanged(test.deployment, test.instance)
		if result != test.expected {
			t.Errorf("hasDeploymentChanged() returned %v but expected %v for the case of %v", result, test.expected, test.msg)
		}
	}
}
