// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_recovery_plan

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9-]{1,63}[a-zA-Z0-9]$`), "The name can contain only letters, numbers, and hyphens. It should start with a letter and end with a letter or a number."),
		},

		"recovery_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateVaultID,
		},

		"source_recovery_fabric_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		"target_recovery_fabric_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		// lintignore:S013
		"shutdown_recovery_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pre_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},

					"post_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},
				},
			},
		},

		// lintignore:S013
		"failover_recovery_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pre_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},

					"post_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},
				},
			},
		},

		// lintignore:S013
		"boot_recovery_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"replicated_protected_items": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},

					"pre_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},

					"post_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},
				},
			},
		},

		"azure_to_azure_settings": replicationRecoveryPlanA2ASchema(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
