// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func privateEndpointSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The Name of the Private Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PrivateLinkName,
		},

		"location": commonschema.Location(),

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"subnet_id": {
			Description:  "The ID of the Subnet from which Private IP Addresses will be allocated for this Private Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the network interface associated with the `private_endpoint`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"name": {
						Description: "The name of the network interface associated with the `private_endpoint`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"custom_network_interface_name": {
			Description: "The custom name of the network interface attached to the Private Endpoint.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"ip_version_type": {
			Description:  "The IP version type for the Private Endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(privateendpoints.PrivateEndpointIPVersionTypeIPvFour),
			ValidateFunc: validation.StringInSlice(privateendpoints.PossibleValuesForPrivateEndpointIPVersionType(), false),
		},

		"private_dns_zone_group": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Private DNS Zone Group.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"name": {
						Description:  "The Name of the Private DNS Zone Group.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.PrivateLinkName,
					},
					"private_dns_zone_ids": {
						Description: "The list of Private DNS Zones to include within the `private_dns_zone_group`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
						},
					},
				},
			},
		},

		"private_service_connection": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The Name of the Private Service Connection.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.PrivateLinkName,
					},
					"manual_connection": {
						Description: "Whether the Private Endpoint requires manual approval from the remote resource owner.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
						ForceNew:    true,
					},
					"private_connection_resource_id": {
						Description:  "The ID of the Private Link Enabled Remote Resource which this Private Endpoint should be connected to. For a Web App or Function App slot, the parent Web App should be used in this field instead of a reference to the slot itself.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
						ExactlyOneOf: []string{"private_service_connection.0.private_connection_resource_alias", "private_service_connection.0.private_connection_resource_id"},
					},
					"private_connection_resource_alias": {
						Description:  "The Service Alias of the Private Link Enabled Remote Resource which this Private Endpoint should be connected to.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.PrivateConnectionResourceAlias,
						ExactlyOneOf: []string{"private_service_connection.0.private_connection_resource_alias", "private_service_connection.0.private_connection_resource_id"},
					},
					"subresource_names": {
						Description: "A list of subresource names which the Private Endpoint can connect to. `subresource_names` corresponds to `group_id`.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.PrivateLinkSubResourceName,
						},
					},
					"request_message": {
						Description:  "A message passed to the owner of the remote resource when the Private Endpoint attempts to establish the connection to the remote resource. The provider allows a maximum request message length of `140` characters, however the request message maximum length is dependent on the service the Private Endpoint is connected to. Only valid if `manual_connection` is set to `true`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 140),
					},
					"private_ip_address": {
						Description: "The private IP address associated with the Private Endpoint, you will have a private IP address assigned to the Private Endpoint even if the connection request was `Rejected`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"action_required": {
						Description: "Whether changes on the service provider require any updates on the consumer.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"description": {
						Description: "The reason for approval/rejection of the connection.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"status": {
						Description: "The status of the connection between the service consumer and provider.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"ip_configuration": {
			Description: "This allows a static IP address to be set for this Private Endpoint, otherwise an address is dynamically allocated from the Subnet.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The Name of the IP Configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.PrivateLinkName,
					},
					"private_ip_address": {
						Description:  "The static IP address within the Private Endpoint's subnet to be used.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subresource_name": {
						Description:  "The subresource this IP address applies to. `subresource_names` corresponds to `group_id`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					// lintignore:S013
					"member_name": {
						Description: "The member name this IP address applies to. If it is not specified, it will use the value of `subresource_name`.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// NOTE: O+C This value should remain optional computed as there are certain cases where Azure will error if you pass in a member id when it isn't expecting one.
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"custom_dns_configs": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"fqdn": {
						Description: "The fully qualified domain name to the `private_endpoint`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"ip_addresses": {
						Description: "A list of all IP Addresses that map to the `private_endpoint` FQDN.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"private_dns_zone_configs": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the Private DNS Zone that the config belongs to.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"id": {
						Description: "The ID of the Private DNS Zone Config.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"private_dns_zone_id": {
						Description: "The ID of the Private DNS Zone that the config belongs to.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"record_sets": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name of the Private DNS Zone that the config belongs to.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"type": {
									Description: "The type of DNS record.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"fqdn": {
									Description: "The fully qualified domain name to the Private DNS Zone.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"ttl": {
									Description: "The time to live for each connection to the Private DNS Zone.",
									Type:        pluginsdk.TypeInt,
									Computed:    true,
								},
								"ip_addresses": {
									Description: "A list of all IP Addresses that map to the Private DNS Zone FQDN.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		connection := s["private_service_connection"].Elem.(*pluginsdk.Resource).Schema

		connection["is_manual_connection"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Description:  "Whether the Private Endpoint requires manual approval from the remote resource owner.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"private_service_connection.0.is_manual_connection", "private_service_connection.0.manual_connection"},
			Deprecated:   "`private_service_connection.is_manual_connection` has been deprecated in favour of `private_service_connection.manual_connection` and will be removed in v6.0 of the AzureRM Provider",
		}

		connection["manual_connection"] = &pluginsdk.Schema{
			Description:  "Whether the Private Endpoint requires manual approval from the remote resource owner.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"private_service_connection.0.is_manual_connection", "private_service_connection.0.manual_connection"},
		}
	}

	return s
}
