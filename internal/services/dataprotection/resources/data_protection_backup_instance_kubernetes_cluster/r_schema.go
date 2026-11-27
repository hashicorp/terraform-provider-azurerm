// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_kubernetes_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Backup Instance Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"vault_id": {
			Description:  "The ID of the Backup Vault within which the Backup Instance Kubernetes Cluster should exist.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: backupvaultresources.ValidateBackupVaultID,
		},

		"backup_policy_id": {
			Description:  "The ID of the Backup Policy.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupPolicyID,
		},

		"kubernetes_cluster_id": {
			Description:  "The ID of the Kubernetes Cluster.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKubernetesClusterID,
		},

		"snapshot_resource_group_name": commonschema.ResourceGroupName(),

		"user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		// TODO: in v6.0, make this block required because the API requires cluster_scoped_resources_enabled
		// and volume_snapshot_enabled to be present. Currently, expand sends defaults when the block is
		// omitted, and read suppresses the default block during flattening to avoid introducing it into state.
		"backup_datasource_parameters": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"excluded_namespaces": {
						Description: "The namespaces to be excluded during backup.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"excluded_resource_types": {
						Description: "The resource types to be excluded during backup.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"cluster_scoped_resources_enabled": {
						Description: "Whether to include cluster scope resources during backup.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},
					"included_namespaces": {
						Description: "The namespaces to be included during backup.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"included_resource_types": {
						Description: "The resource types to be included during backup.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"label_selectors": {
						Description: "The resources with such label selectors to be included during backup.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"volume_snapshot_enabled": {
						Description: "Whether to take volume snapshots during backup.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Description: "The protection state of the Backup Instance Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
