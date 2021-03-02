package nginxingresscontroller

import (
	"fmt"

	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func sccForNginxIngressController(name string) *secv1.SecurityContextConstraints {
	var uid int64 = 101

	allowPrivilegeEscalation := true

	return &secv1.SecurityContextConstraints{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		AllowHostPorts:           false,
		AllowPrivilegedContainer: false,
		RunAsUser: secv1.RunAsUserStrategyOptions{
			Type: "MustRunAs",
			UID:  &uid,
		},
		Users:                    nil,
		AllowHostDirVolumePlugin: false,
		AllowHostIPC:             false,
		SELinuxContext: secv1.SELinuxContextStrategyOptions{
			Type: "MustRunAs",
		},
		ReadOnlyRootFilesystem: false,
		FSGroup: secv1.FSGroupStrategyOptions{
			Type: "MustRunAs",
		},
		SupplementalGroups: secv1.SupplementalGroupsStrategyOptions{
			Type: "MustRunAs",
		},
		Volumes:                  []secv1.FSType{"secret"},
		AllowHostPID:             false,
		AllowHostNetwork:         false,
		AllowPrivilegeEscalation: &allowPrivilegeEscalation,
		RequiredDropCapabilities: []corev1.Capability{"ALL"},
		DefaultAddCapabilities:   []corev1.Capability{"NET_BIND_SERVICE"},
		AllowedCapabilities:      nil,
	}
}

func userForSCC(namespace string, name string) string {
	return fmt.Sprintf("system:serviceaccount:%v:%v", namespace, name)
}
