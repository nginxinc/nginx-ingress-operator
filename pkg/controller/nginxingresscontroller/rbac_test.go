package nginxingresscontroller

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClusterRoleForNginxIngressController(t *testing.T) {
	name := "my-rcluster-role"
	expected := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Rules: []rbacv1.PolicyRule{
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{""},
				Resources: []string{"services", "endpoints"},
			},
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{""},
				Resources: []string{"secrets"},
			},
			{
				Verbs:     []string{"get", "list", "watch", "update", "create"},
				APIGroups: []string{""},
				Resources: []string{"configmaps"},
			},
			{
				Verbs:     []string{"list", "watch"},
				APIGroups: []string{""},
				Resources: []string{"pods"},
			},
			{
				Verbs:     []string{"create", "patch"},
				APIGroups: []string{""},
				Resources: []string{"events"},
			},
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses"},
			},
			{
				Verbs:     []string{"update"},
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses/status"},
			},
			{
				Verbs:     []string{"get", "create"},
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingressclasses"},
			},
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{"k8s.nginx.org"},
				Resources: []string{"virtualservers", "virtualserverroutes", "globalconfigurations", "transportservers", "policies"},
			},
			{
				Verbs:     []string{"update"},
				APIGroups: []string{"k8s.nginx.org"},
				Resources: []string{"virtualservers/status", "virtualserverroutes/status", "policies/status", "transportservers/status"},
			},
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{"appprotect.f5.com"},
				Resources: []string{"aplogconfs", "appolicies", "apusersigs"},
			},
		},
	}

	result := clusterRoleForNginxIngressController(name)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("clusterRoleForNginxIngressController(%v) mismatch (-want +got):\n%s", name, diff)
	}
}

func TestSubjectForServiceAccount(t *testing.T) {
	name := "my-sa"
	namespace := "my-nginx-ingress"
	expected := rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      name,
		Namespace: namespace,
	}

	result := subjectForServiceAccount(namespace, name)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("subjectForServiceAccount(%v, %v) mismatch (-want +got):\n%s", namespace, name, diff)
	}
}

func TestClusterRoleBindingForNginxIngressController(t *testing.T) {
	name := "my-cluster-role-binding"
	expected := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     name,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	result := clusterRoleBindingForNginxIngressController(name)
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Errorf("clusterRoleBindingForNginxIngressController(%v) mismatch (-want +got):\n%s", name, diff)
	}
}
