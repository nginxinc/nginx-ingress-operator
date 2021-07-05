package scc

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultName = "nginx-ingress-scc"

func sccConfigTemplate() *secv1.SecurityContextConstraints {
	var uid int64 = 101
	allowPrivilegeEscalation := true

	return &secv1.SecurityContextConstraints{
		ObjectMeta: v1.ObjectMeta{
			Name: defaultName,
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

func serviceAccountName(namespace string, name string) string {
	return fmt.Sprintf("system:serviceaccount:%v:%v", namespace, name)
}

func Create(client client.Client, log logr.Logger) error {
	scc := sccConfigTemplate()
	err := client.Get(context.TODO(), types.NamespacedName{Name: defaultName, Namespace: v1.NamespaceAll}, scc)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("no previous SecurityContextConstraints found, creating a new one.")
			err = client.Create(context.TODO(), scc)
			if err != nil {
				return fmt.Errorf("error creating SecurityContextConstraints: %w", err)
			}
		}
		return fmt.Errorf("error getting scc: %w", err)
	}

	return nil
}

func AddServiceAccount(client client.Client, namespace string, name string) error {
	scc := sccConfigTemplate()
	err := client.Get(context.TODO(), types.NamespacedName{Name: defaultName, Namespace: v1.NamespaceAll}, scc)
	if err != nil {
		return fmt.Errorf("failed to get scc: %w", err)
	}

	saName := serviceAccountName(namespace, name)
	for _, u := range scc.Users {
		if u == saName {
			// scc already has the service account name
			return nil
		}
	}

	scc.Users = append(scc.Users, saName)
	err = client.Update(context.TODO(), scc)
	if err != nil {
		return fmt.Errorf("failed to update scc: %w", err)
	}

	return nil
}

func RemoveServiceAccount(client client.Client, namespace string, name string) error {
	scc := sccConfigTemplate()
	err := client.Get(context.TODO(), types.NamespacedName{Name: defaultName, Namespace: v1.NamespaceAll}, scc)
	if err != nil {
		return fmt.Errorf("failed to get scc: %w", err)
	}

	scc.Users = removeStringValue(scc.Users, serviceAccountName(namespace, name))

	err = client.Update(context.TODO(), scc)
	if err != nil {
		return fmt.Errorf("failed to update scc: %w", err)
	}
	return nil
}

func removeStringValue(values []string, value string) []string {
	var filtered []string
	for _, v := range values {
		if v != value {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
