package nginxingresscontroller

import (
	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func daemonSetForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) *appsv1.DaemonSet {
	runAsUser := new(int64)
	allowPrivilegeEscalation := new(bool)
	*runAsUser = 101
	*allowPrivilegeEscalation = true

	return &appsv1.DaemonSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{"app": instance.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Name:      instance.Name,
					Namespace: instance.Namespace,
					Labels:    map[string]string{"app": instance.Name},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: instance.Name,
					Containers: []corev1.Container{
						{
							Name:            instance.Name,
							Image:           generateImage(instance.Spec.Image.Repository, instance.Spec.Image.Tag),
							ImagePullPolicy: corev1.PullPolicy(instance.Spec.Image.PullPolicy),
							Args:            generatePodArgs(instance),
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
}

func hasDaemonSetChanged(ds *appsv1.DaemonSet, instance *k8sv1alpha1.NginxIngressController) bool {
	// There is only 1 container in our template
	container := ds.Spec.Template.Spec.Containers[0]
	if container.Image != generateImage(instance.Spec.Image.Repository, instance.Spec.Image.Tag) {
		return true
	}

	if container.ImagePullPolicy != corev1.PullPolicy(instance.Spec.Image.PullPolicy) {
		return true
	}

	return hasDifferentArguments(container, instance)
}

func updateDaemonSet(ds *appsv1.DaemonSet, instance *k8sv1alpha1.NginxIngressController) *appsv1.DaemonSet {
	ds.Spec.Template.Spec.Containers[0].Image = generateImage(instance.Spec.Image.Repository, instance.Spec.Image.Tag)
	ds.Spec.Template.Spec.Containers[0].Args = generatePodArgs(instance)
	return ds
}
