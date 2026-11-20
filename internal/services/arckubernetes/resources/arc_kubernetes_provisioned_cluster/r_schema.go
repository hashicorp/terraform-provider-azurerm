// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_provisioned_cluster

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Arc Kubernetes Provisioned Cluster.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-_a-zA-Z0-9]{1,260}$"),
				"The name of Arc Kubernetes Provisioned Cluster can only include alphanumeric characters, underscores, hyphens, has a maximum length of 260 characters, and must be unique.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"azure_active_directory": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_group_object_ids": {
						Description: "A list of IDs of Microsoft Entra ID Groups. All members of the specified Microsoft Entra ID Groups have the cluster administrator access to the Kubernetes cluster.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},

					"azure_rbac_enabled": {
						Description: "Whether to enable Azure RBAC for Kubernetes authorization.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"tenant_id": {
						Description:  "The Tenant ID to use for authentication. If not specified, the Tenant of the Arc Kubernetes Cluster will be used.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"arc_agent_desired_version": {
			Description:  "The version of the Arc agents to be installed on the cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"arc_agent_auto_upgrade_enabled": {
			Description: "Whether the Arc agents will be upgraded automatically to the latest version.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"agent_version": {
			Description: "The version of the agent running on the cluster resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"distribution": {
			Description: "The distribution running on this Arc Kubernetes Provisioned Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"infrastructure": {
			Description: "The infrastructure on which the Arc Kubernetes Provisioned Cluster is running on.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kubernetes_version": {
			Description: "The Kubernetes version of the cluster resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"offering": {
			Description: "The cluster offering.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"total_core_count": {
			Description: "The number of CPU cores present in the cluster resource.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"total_node_count": {
			Description: "The number of nodes present in the cluster resource.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},
	}
}
