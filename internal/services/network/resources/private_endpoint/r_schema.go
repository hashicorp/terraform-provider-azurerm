// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func privateEndpointSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.PrivateLinkName,
		},

		"location": commonschema.Location(),

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"subnet_id": {
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"custom_network_interface_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"private_dns_zone_group": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.PrivateLinkName,
					},
					"private_dns_zone_ids": {
						Type:     pluginsdk.TypeList,
						Required: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.PrivateLinkName,
					},
					"is_manual_connection": {
						Type:     pluginsdk.TypeBool,
						Required: true,
						ForceNew: true,
					},
					"private_connection_resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: azure.ValidateResourceID,
						ExactlyOneOf: []string{"private_service_connection.0.private_connection_resource_alias", "private_service_connection.0.private_connection_resource_id"},
					},
					"private_connection_resource_alias": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.PrivateConnectionResourceAlias,
						ExactlyOneOf: []string{"private_service_connection.0.private_connection_resource_alias", "private_service_connection.0.private_connection_resource_id"},
					},
					"subresource_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.PrivateLinkSubResourceName,
						},
					},
					"request_message": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringLenBetween(1, 140),
					},
					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.PrivateLinkName,
					},
					"private_ip_address": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subresource_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					// lintignore:S013
					"member_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
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
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"private_dns_zone_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"record_sets": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"fqdn": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"ttl": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
								"ip_addresses": {
									Type:     pluginsdk.TypeList,
									Computed: true,
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
}
