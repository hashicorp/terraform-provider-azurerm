// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package vpn_server_configuration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func vpnServerConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"vpn_authentication_types": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVpnAuthenticationType(), false),
			},
		},

		// Optional
		"azure_active_directory_authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"audience": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"issuer": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"tenant": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"client_revoked_certificate": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"thumbprint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"client_root_certificate": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"public_cert_data": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"ipsec_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dh_group": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForDhGroup(), false),
					},

					"ike_encryption": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeEncryption(), false),
					},

					"ike_integrity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIkeIntegrity(), false),
					},

					"ipsec_encryption": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecEncryption(), false),
					},

					"ipsec_integrity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForIPsecIntegrity(), false),
					},

					"pfs_group": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForPfsGroup(), false),
					},

					"sa_lifetime_seconds": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"sa_data_size_kilobytes": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"radius": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"secret": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
									Sensitive:    true,
								},

								"score": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 30),
								},
							},
						},
					},

					"client_root_certificate": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"thumbprint": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"server_root_certificate": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"public_cert_data": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"vpn_protocols": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(virtualwans.PossibleValuesForVpnGatewayTunnelingProtocol(), false),
			},
		},

		"tags": commonschema.Tags(),
	}
}
