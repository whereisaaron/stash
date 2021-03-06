package controller

import (
	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	rbac_util "github.com/appscode/kutil/rbac/v1beta1"
	api "github.com/appscode/stash/apis/stash/v1alpha1"
	"github.com/appscode/stash/pkg/util"
	apps "k8s.io/api/apps/v1beta1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	rbac "k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SidecarClusterRole = "stash-sidecar"
	KubectlRole        = "stash-kubectl"
	RecoveryRole       = "stash-recovery"
)

func (c *StashController) getSidecarRoleBindingName(name string) string {
	return name + "-" + SidecarClusterRole
}

func (c *StashController) ensureSidecarRoleBinding(resource *core.ObjectReference, sa string) error {
	meta := metav1.ObjectMeta{
		Namespace: resource.Namespace,
		Name:      c.getSidecarRoleBindingName(resource.Name),
	}
	_, _, err := rbac_util.CreateOrPatchRoleBinding(c.k8sClient, meta, func(in *rbac.RoleBinding) *rbac.RoleBinding {
		in.ObjectMeta = util.EnsureOwnerReference(in.ObjectMeta, resource)

		if in.Annotations == nil {
			in.Annotations = map[string]string{}
		}

		in.RoleRef = rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     SidecarClusterRole,
		}
		in.Subjects = []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      sa,
				Namespace: resource.Namespace,
			},
		}
		return in
	})
	return err
}

func (c *StashController) ensureSidecarRoleBindingDeleted(resource metav1.ObjectMeta) error {
	log.Infof("Deleting RoleBinding %s/%s", resource.Namespace, c.getSidecarRoleBindingName(resource.Name))
	return c.k8sClient.RbacV1beta1().
		RoleBindings(resource.Namespace).
		Delete(c.getSidecarRoleBindingName(resource.Name), &metav1.DeleteOptions{})
}

func (c *StashController) ensureSidecarClusterRole() error {
	meta := metav1.ObjectMeta{Name: SidecarClusterRole}
	_, _, err := rbac_util.CreateOrPatchClusterRole(c.k8sClient, meta, func(in *rbac.ClusterRole) *rbac.ClusterRole {
		if in.Labels == nil {
			in.Labels = map[string]string{}
		}
		in.Labels["app"] = "stash"

		in.Rules = []rbac.PolicyRule{
			{
				APIGroups: []string{api.SchemeGroupVersion.Group},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{apps.GroupName},
				Resources: []string{"deployments", "statefulsets"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{extensions.GroupName},
				Resources: []string{"daemonsets", "replicasets"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"replicationcontrollers", "secrets"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"configmaps"},
				Verbs:     []string{"create", "update", "get"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"events"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{batch.GroupName},
				Resources: []string{"jobs"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{rbac.GroupName},
				Resources: []string{"clusterroles", "roles", "rolebindings"},
				Verbs:     []string{"get", "create"},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"serviceaccounts"},
				Verbs:     []string{"get", "create"},
			},
		}
		return in
	})
	return err
}

// use kubectl-role, service-account and role-binding name same as job name
// set job as owner of role, service-account and role-binding
func (c *StashController) ensureKubectlRBAC(resource *core.ObjectReference) error {
	// ensure roles
	meta := metav1.ObjectMeta{
		Name:      KubectlRole,
		Namespace: resource.Namespace,
	}
	_, _, err := rbac_util.CreateOrPatchRole(c.k8sClient, meta, func(in *rbac.Role) *rbac.Role {
		in.ObjectMeta = util.EnsureOwnerReference(in.ObjectMeta, resource)

		if in.Labels == nil {
			in.Labels = map[string]string{}
		}
		in.Labels["app"] = "stash"

		in.Rules = []rbac.PolicyRule{
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "delete", "deletecollection"},
			},
		}
		return in
	})
	if err != nil {
		return err
	}

	// ensure service account
	meta = metav1.ObjectMeta{
		Name:      resource.Name,
		Namespace: resource.Namespace,
	}
	_, _, err = core_util.CreateOrPatchServiceAccount(c.k8sClient, meta, func(in *core.ServiceAccount) *core.ServiceAccount {
		in.ObjectMeta = util.EnsureOwnerReference(in.ObjectMeta, resource)
		if in.Labels == nil {
			in.Labels = map[string]string{}
		}
		in.Labels["app"] = "stash"
		return in
	})
	if err != nil {
		return err
	}

	// ensure role binding
	_, _, err = rbac_util.CreateOrPatchRoleBinding(c.k8sClient, meta, func(in *rbac.RoleBinding) *rbac.RoleBinding {
		in.ObjectMeta = util.EnsureOwnerReference(in.ObjectMeta, resource)

		if in.Labels == nil {
			in.Labels = map[string]string{}
		}
		in.Labels["app"] = "stash"

		in.RoleRef = rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "Role",
			Name:     KubectlRole,
		}
		in.Subjects = []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      meta.Name,
				Namespace: meta.Namespace,
			},
		}
		return in
	})
	return err
}

// use sidecar-cluster-role, service-account and role-binding name same as job name
// set job as owner of service-account and role-binding
func (c *StashController) ensureRecoveryRBAC(resource *core.ObjectReference) error {
	// ensure service account
	meta := metav1.ObjectMeta{
		Name:      resource.Name,
		Namespace: resource.Namespace,
	}
	_, _, err := core_util.CreateOrPatchServiceAccount(c.k8sClient, meta, func(in *core.ServiceAccount) *core.ServiceAccount {
		in.ObjectMeta = util.EnsureOwnerReference(in.ObjectMeta, resource)
		if in.Labels == nil {
			in.Labels = map[string]string{}
		}
		in.Labels["app"] = "stash"
		return in
	})
	if err != nil {
		return err
	}

	// ensure role binding
	_, _, err = rbac_util.CreateOrPatchRoleBinding(c.k8sClient, meta, func(in *rbac.RoleBinding) *rbac.RoleBinding {
		in.ObjectMeta = util.EnsureOwnerReference(in.ObjectMeta, resource)

		if in.Labels == nil {
			in.Labels = map[string]string{}
		}
		in.Labels["app"] = "stash"

		in.RoleRef = rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     SidecarClusterRole,
		}
		in.Subjects = []rbac.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      meta.Name,
				Namespace: meta.Namespace,
			},
		}
		return in
	})
	return err
}
