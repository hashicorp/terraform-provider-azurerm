// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_deployment_safeguard

type KubernetesClusterDeploymentSafeguardResourceModel struct {
	KubernetesClusterId       string   `tfschema:"kubernetes_cluster_id"`
	Level                     string   `tfschema:"level"`
	ExcludedNamespaces        []string `tfschema:"excluded_namespaces"`
	PodSecurityStandardsLevel string   `tfschema:"pod_security_standards_level"`
}

func (r Resource) ModelObject() any {
	return &KubernetesClusterDeploymentSafeguardResourceModel{}
}
