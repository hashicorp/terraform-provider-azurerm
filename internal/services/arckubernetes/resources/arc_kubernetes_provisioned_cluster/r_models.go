// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_provisioned_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ArcKubernetesProvisionedClusterModel struct {
	AgentVersion               string                         `tfschema:"agent_version"`
	ArcAgentAutoUpgradeEnabled bool                           `tfschema:"arc_agent_auto_upgrade_enabled"`
	ArcAgentDesiredVersion     string                         `tfschema:"arc_agent_desired_version"`
	AzureActiveDirectory       []AzureActiveDirectoryModel    `tfschema:"azure_active_directory"`
	Distribution               string                         `tfschema:"distribution"`
	Identity                   []identity.ModelSystemAssigned `tfschema:"identity"`
	Infrastructure             string                         `tfschema:"infrastructure"`
	KubernetesVersion          string                         `tfschema:"kubernetes_version"`
	Location                   string                         `tfschema:"location"`
	Name                       string                         `tfschema:"name"`
	Offering                   string                         `tfschema:"offering"`
	ResourceGroupName          string                         `tfschema:"resource_group_name"`
	Tags                       map[string]string              `tfschema:"tags"`
	TotalCoreCount             int64                          `tfschema:"total_core_count"`
	TotalNodeCount             int64                          `tfschema:"total_node_count"`
}

type AzureActiveDirectoryModel struct {
	AdminGroupObjectIds []string `tfschema:"admin_group_object_ids"`
	AzureRbacEnabled    bool     `tfschema:"azure_rbac_enabled"`
	TenantId            string   `tfschema:"tenant_id"`
}

func (r Resource) ModelObject() any {
	return &ArcKubernetesProvisionedClusterModel{}
}
