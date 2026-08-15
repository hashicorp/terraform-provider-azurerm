// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SynapseWorkspaceV0ToV1{}

type SynapseWorkspaceV0ToV1 struct{}

func (SynapseWorkspaceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_data_lake_gen2_filesystem_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"sql_administrator_login": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"sql_administrator_login_password": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"linking_allowed_for_aad_tenant_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"compute_subnet_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"data_exfiltration_protection_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"managed_virtual_network_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"connectivity_endpoints": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"managed_resource_group_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"azure_devops_repo": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"last_commit_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"project_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"github_repo": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"git_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"last_commit_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"purview_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"sql_identity_control_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_versionless_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"key_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "cmk",
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"azuread_authentication_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (SynapseWorkspaceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// IDs imported while this resource parsed them with the legacy resourceids.ParseAzureResourceID
		// can contain non-canonically cased static segments (e.g. `resourcegroups`, `microsoft.synapse`),
		// which the case-sensitive SDK parser rejects - normalise them to the canonical casing
		oldId := rawState["id"].(string)
		id, err := workspaces.ParseWorkspaceIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
