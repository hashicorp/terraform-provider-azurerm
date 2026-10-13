// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceApiManagementSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"publisher_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApiManagementServicePublisherName,
		},

		"publisher_email": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApiManagementServicePublisherEmail,
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApimSkuName(),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"virtual_network_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(apimanagementservice.VirtualNetworkTypeNone),
			ValidateFunc: validation.StringInSlice(apimanagementservice.PossibleValuesForVirtualNetworkType(), false),
		},

		"virtual_network_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"client_certificate_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"gateway_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"min_api_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"notification_sender_email": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"additional_location": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"gateway_disabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"location": commonschema.LocationWithoutForceNew(),

					"virtual_network_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateSubnetID,
								},
							},
						},
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(0, 50),
					},

					"zones": commonschema.ZonesMultipleOptional(),

					"gateway_regional_url": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"public_ip_addresses": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Computed: true,
					},

					"public_ip_address_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},

					"private_ip_addresses": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Computed: true,
					},
				},
			},
		},

		"certificate": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 10,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"encoded_certificate": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},

					"certificate_password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"store_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(apimanagementservice.PossibleValuesForStoreName(), false),
					},

					"expiry": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"subject": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"thumbprint": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"protocols": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"http2_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"security": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend_ssl30_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"backend_tls10_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"backend_tls11_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"frontend_ssl30_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"frontend_tls10_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"frontend_tls11_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"triple_des_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_rsa_with_aes128_gcm_sha256_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_rsa_with_aes256_cbc_sha256_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_rsa_with_aes128_cbc_sha256_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_rsa_with_aes256_gcm_sha384_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_rsa_with_aes256_cbc_sha_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"tls_rsa_with_aes128_cbc_sha_ciphers_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"hostname_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"management": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: helpers.ApiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"portal": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: helpers.ApiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"developer_portal": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: helpers.ApiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"proxy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: helpers.ApiManagementResourceHostnameProxySchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"scm": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: helpers.ApiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
				},
			},
		},

		"sign_in": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"delegation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subscriptions_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"user_registration_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
					"validation_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
					},
				},
			},
		},

		"sign_up": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"terms_of_service": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"consent_required": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"text": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},

		"zones": commonschema.ZonesMultipleOptional(),

		"gateway_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"management_api_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"gateway_regional_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"public_ip_address_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidatePublicIPAddressID,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"private_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"portal_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"developer_portal_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scm_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tenant_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"primary_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"secondary_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
