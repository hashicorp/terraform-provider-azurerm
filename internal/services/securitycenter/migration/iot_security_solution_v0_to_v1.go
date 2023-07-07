// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
)

type SecurityCenterIotSecuritySolutionV0ToV1 struct{}

func (s SecurityCenterIotSecuritySolutionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"iothub_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: set.HashStringIgnoreCase,
		},

		"additional_workspace": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_types": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"workspace_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"disabled_data_sources": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"log_analytics_workspace_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"log_unmasked_ips_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"events_to_export": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"recommendations_enabled": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"acr_authentication": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"agent_send_unutilized_msg": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"baseline": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"edge_hub_mem_optimize": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"edge_logging_option": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"inconsistent_module_settings": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"install_agent": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"ip_filter_deny_all": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"ip_filter_permissive_rule": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"open_ports": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"permissive_firewall_policy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"permissive_input_firewall_rules": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"permissive_output_firewall_rules": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"privileged_docker_options": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"shared_credentials": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"vulnerable_tls_cipher_suite": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"query_for_resources": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"query_subscription_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": tags.Schema(),
	}
}

func (s SecurityCenterIotSecuritySolutionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.IotSecuritySolutionIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
