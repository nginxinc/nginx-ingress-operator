package nginxingresscontroller

import (
	"reflect"
	"testing"

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
				APIGroups: []string{"extensions"},
				Resources: []string{"ingresses"},
			},
			{
				Verbs:     []string{"update"},
				APIGroups: []string{"extensions"},
				Resources: []string{"ingresses/status"},
			},
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{"k8s.nginx.org"},
				Resources: []string{"virtualservers", "virtualserverroutes", "globalconfigurations", "transportservers", "policies"},
			},
			{
				Verbs:     []string{"update"},
				APIGroups: []string{"k8s.nginx.org"},
				Resources: []string{"virtualservers/status", "virtualserverroutes/status"},
			},
			{
				Verbs:     []string{"get", "list", "watch"},
				APIGroups: []string{"appprotect.f5.com"},
				Resources: []string{"aplogconfs", "appolicies"},
			},
		},
	}

	result := clusterRoleForNginxIngressController(name)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("clusterRoleForNginxIngressController(%v) returned %+v but expected %+v", name, result, expected)
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
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("subjectForServiceAccount(%v, %v) returned %+v but expected %+v", namespace, name, result, expected)
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
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("clusterRoleBindingForNginxIngressController(%v) returned %+v but expected %+v", name, result, expected)
	}
}
