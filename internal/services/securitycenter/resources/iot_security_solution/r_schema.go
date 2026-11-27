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
			Description:  "The name of the Iot Security Solution.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotSecuritySolutionName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Description:  "The Display Name for this Iot Security Solution.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"iothub_ids": {
			Description: "The IoT Hub resource IDs to which this Iot Security Solution is applied.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Set:         set.HashStringIgnoreCase,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: iothubValidate.IotHubID,
			},
		},

		"additional_workspace": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_types": {
						Description: "A list of data types which are sent to the workspace.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInEnumSlice(security.PossibleAdditionalWorkspaceDataTypeValues(), false),
						},
					},

					"workspace_id": {
						Description:  "The resource ID of the Log Analytics Workspace.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
				},
			},
		},

		"disabled_data_sources": {
			Description: "A list of disabled data sources for the Iot Security Solution.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInEnumSlice(security.PossibleDataSourceValues(), false),
			},
		},

		"log_analytics_workspace_id": {
			Description:      "The Log Analytics Workspace ID to which the security data will be sent.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     workspaces.ValidateWorkspaceID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"enabled": {
			Description: "Whether the Iot Security Solution is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"log_unmasked_ips_enabled": {
			Description: "Whether IP addresses should be unmasked in the log.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"events_to_export": {
			Description: "A list of data which is to be exported to the analytic workspace.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
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
						Description: "Whether Principal Authentication is enabled for the ACR repository.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"agent_send_unutilized_msg": {
						Description: "Whether the agent sends underutilized messages is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"baseline": {
						Description: "Whether to identify security related system configuration issues.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"edge_hub_mem_optimize": {
						Description: "Whether to optimize IoT Edge Hub memory.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"edge_logging_option": {
						Description: "Whether to configure logging for the IoT Edge module.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"inconsistent_module_settings": {
						Description: "Whether inconsistent module settings are enabled for SecurityGroup.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"install_agent": {
						Description: "Whether to install the Azure IoT Security agent.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"ip_filter_deny_all": {
						Description: "Whether to add the default IP filter policy that denies all.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"ip_filter_permissive_rule": {
						Description: "Whether the IP filter rule's source allowable IP range is too large.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"open_ports": {
						Description: "Whether any ports are open on the device.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"permissive_firewall_policy": {
						Description: "Whether a Firewall Policy exists which allows necessary communication to/from the device.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"permissive_input_firewall_rules": {
						Description: "Whether only necessary addresses or ports are allowed in.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"permissive_output_firewall_rules": {
						Description: "Whether only necessary addresses or ports are allowed out.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"privileged_docker_options": {
						Description: "Whether high level permissions are necessary for the module.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"shared_credentials": {
						Description: "Whether any credentials are shared among devices.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"vulnerable_tls_cipher_suite": {
						Description: "Whether the TLS cypher suite needs to be updated.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		"query_for_resources": {
			Description: "An Azure Resource Graph query used to set the resources monitored.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"query_subscription_ids": {
			Description: "A list of subscription IDs on which the user defined resources query should be executed.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"tags": commonschema.Tags(),
	}
}
