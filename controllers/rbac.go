package controllers

import (
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *NginxIngressControllerReconciler) clusterRoleForNginxIngressController(name string) *rbacv1.ClusterRole {
	rules := []rbacv1.PolicyRule{
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
	}
	rbac := &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Rules: rules,
	}
	return rbac
}

func subjectForServiceAccount(namespace string, name string) rbacv1.Subject {
	sa := rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      name,
		Namespace: namespace,
	}
	return sa
}

func (r *NginxIngressControllerReconciler) clusterRoleBindingForNginxIngressController(name string) *rbacv1.ClusterRoleBinding {
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     name,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	return crb
}
