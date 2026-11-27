// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_verifier_workspace_reachability_analysis_intent

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/reachabilityanalysisintents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Network Manager Verifier Workspace Reachability Analysis Intent.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"verifier_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&reachabilityanalysisintents.VerifierWorkspaceId{}, "Network Manager Verifier Workspace"),

		"source_resource_id": {
			Description: "The ID of the source resource.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				commonids.ValidatePublicIPAddressID,
				commonids.ValidateSubnetID,
				commonids.ValidateVirtualMachineID,
			),
		},

		"destination_resource_id": {
			Description: "The ID of the destination resource.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				commonids.ValidatePublicIPAddressID,
				commonids.ValidateSqlServerID,
				commonids.ValidateStorageAccountID,
				commonids.ValidateSubnetID,
				commonids.ValidateVirtualMachineID,
				cosmosdb.ValidateDatabaseAccountID,
			),
		},

		"ip_traffic": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destination_ips": {
						Description: "A list of IPv4 or IPv6 addresses or ranges using CIDR notation of the source you want to verify.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsCIDR,
								validation.IsIPAddress,
							),
						},
					},
					"destination_ports": {
						Description: "A list of ports or ranges of the destination you want to verify. To specify any port, use `[\"*\"]`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.IpTrafficPort,
						},
					},
					"source_ips": {
						Description: "A list of IPv4 or IPv6 addresses or ranges using CIDR notation of the source you want to verify.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsCIDR,
								validation.IsIPAddress,
							),
						},
					},
					"source_ports": {
						Description: "A list of ports or ranges of the source you want to verify. To specify any port, use `[\"*\"]`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.IpTrafficPort,
						},
					},
					"protocols": {
						Description: "A list of network protocols.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(reachabilityanalysisintents.PossibleValuesForNetworkProtocol(), false),
						},
					},
				},
			},
		},

		"description": {
			Description:  "The description of the resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
