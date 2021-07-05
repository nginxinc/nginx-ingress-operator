package scc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSccForNginxIngressController(t *testing.T) {
	var uid int64 = 101
	allowPrivilegeEscalation := true

	expected := &secv1.SecurityContextConstraints{
		ObjectMeta: v1.ObjectMeta{
			Name: "nginx-ingress-scc",
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

	result := sccConfigTemplate()
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("sccConfigTemplate() mismatch (-want +got):\n%s", diff)
	}
}

func TestServiceAccountName(t *testing.T) {
	namespace := "my-nginx-ingress"
	name := "my-nginx-ingress-controller"
	expected := fmt.Sprintf("system:serviceaccount:%v:%v", namespace, name)

	result := serviceAccountName(namespace, name)
	if expected != result {
		t.Errorf("serviceAccountName(%v, %v) returned %v but expected %v", namespace, name, result, expected)
	}
}
