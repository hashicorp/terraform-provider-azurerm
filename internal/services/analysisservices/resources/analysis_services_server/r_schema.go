// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package analysis_services_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func analysisServicesServerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Analysis Services Server. Only lowercase Alphanumeric characters allowed, starting with a letter.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Description: "SKU for the Analysis Services Server.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"D1",
				"B1",
				"B2",
				"S0",
				"S1",
				"S2",
				"S4",
				"S8",
				"S9",
				"S8v2",
				"S9v2",
			}, false),
		},

		"admin_users": {
			Description: "List of email addresses of admin users.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"power_bi_service_enabled": {
			Description: "Whether the Power BI service is allowed to access or not.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"ipv4_firewall_rule": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			RequiredWith: []string{"power_bi_service_enabled"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the firewall rule.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"range_start": {
						Description:  "Start of the firewall rule range as IPv4 address.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
					"range_end": {
						Description:  "End of the firewall rule range as IPv4 address.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsIPv4Address,
					},
				},
			},
			Set: hashAnalysisServicesServerIPv4FirewallRule,
		},

		"querypool_connection_mode": {
			Description:  "Controls how the read-write server is used in the query pool. If this value is set to `All` then read-write servers are also used for queries. Otherwise with `ReadOnly` these servers do not participate in query operations.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(servers.ConnectionModeAll),
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForConnectionMode(), false),
		},

		"backup_blob_container_uri": {
			Description:  "URI and SAS token for a blob container to store backups.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"server_full_name": {
			Description: "The full name of the Analysis Services Server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
