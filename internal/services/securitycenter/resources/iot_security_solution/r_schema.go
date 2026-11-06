// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_solution

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func iotSecuritySolutionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotSecuritySolutionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"iothub_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: iothubValidate.IotHubID,
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
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInEnumSlice(security.PossibleAdditionalWorkspaceDataTypeValues(), false),
						},
					},

					"workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
				},
			},
		},

		"disabled_data_sources": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInEnumSlice(security.PossibleDataSourceValues(), false),
			},
		},

		"log_analytics_workspace_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     workspaces.ValidateWorkspaceID,
			DiffSuppressFunc: suppress.CaseDifference,
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
				ValidateFunc: validation.StringInSlice([]string{
					string(security.RawEvents),
				}, false),
			},
		},

		"recommendations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
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
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"query_subscription_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"tags": commonschema.Tags(),
	}
}
