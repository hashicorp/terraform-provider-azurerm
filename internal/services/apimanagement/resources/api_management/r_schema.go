// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
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
			Description:  "The name of publisher/company.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApiManagementServicePublisherName,
		},

		"publisher_email": {
			Description:  "The email of publisher/company.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApiManagementServicePublisherEmail,
		},

		"sku_name": {
			Description:  "`sku_name` is a string consisting of two parts separated by an underscore(\\_). The first part is the `name`. The second part is the `capacity` (e.g. the number of deployed units of the `sku`), which must be a positive `integer` (e.g. `Developer_1`).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ApimSkuName(),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"virtual_network_type": {
			Description:  "The type of virtual network you want to use.",
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
						Description:  "The id of the subnet that will be used for the API Management.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"client_certificate_enabled": {
			Description: "Enforce a client certificate to be presented on each request to the gateway? This is only supported when SKU type is `Consumption`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"gateway_disabled": {
			Description: "Disable the gateway in main region? This is only supported when `additional_location` is set.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"min_api_version": {
			Description:  "The version which the control plane API calls to API Management service are limited with version equal to or newer than.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"notification_sender_email": {
			Description: "Email address from which the notification will be sent.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"additional_location": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"gateway_disabled": {
						Description: "Only valid for an Api Management service deployed in multiple locations. This can be used to disable the gateway in this additional location.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"location": commonschema.LocationWithoutForceNew(),

					"virtual_network_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Description:  "The id of the subnet that will be used for the API Management.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateSubnetID,
								},
							},
						},
					},

					"capacity": {
						Description:  "The number of compute units in this region.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(0, 50),
					},

					"zones": commonschema.ZonesMultipleOptional(),

					"gateway_regional_url": {
						Description: "The URL of the Regional Gateway for the API Management Service in the specified region.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"public_ip_addresses": {
						Description: "Public Static Load Balanced IP addresses of the API Management service in the additional location. Available only for Basic, Standard and Premium SKU.",
						Type:        pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Computed: true,
					},

					"public_ip_address_id": {
						Description:  "ID of a standard SKU IPv4 Public IP.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},

					"private_ip_addresses": {
						Description: "The Private IP addresses of the API Management Service. Available only when the API Manager instance is using Virtual Network mode.",
						Type:        pluginsdk.TypeList,
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
						Description: "The Base64 Encoded PFX or Base64 Encoded X.509 Certificate.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						Sensitive:   true,
					},

					"certificate_password": {
						Description: "The password for the certificate.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Sensitive:   true,
					},

					"store_name": {
						Description:  "The name of the Certificate Store where this certificate should be stored.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(apimanagementservice.PossibleValuesForStoreName(), false),
					},

					"expiry": {
						Description: "The expiration date of the certificate in RFC3339 format: `2000-01-02T03:04:05Z`.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"subject": {
						Description: "The subject of the certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"thumbprint": {
						Description: "The thumbprint of the certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
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
						Description: "Should HTTP/2 be supported by the API Management Service?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
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
						Description: "Should SSL 3.0 be enabled on the backend of the gateway?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"backend_tls10_enabled": {
						Description: "Should TLS 1.0 be enabled on the backend of the gateway?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"backend_tls11_enabled": {
						Description: "Should TLS 1.1 be enabled on the backend of the gateway?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"frontend_ssl30_enabled": {
						Description: "Should SSL 3.0 be enabled on the frontend of the gateway?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"frontend_tls10_enabled": {
						Description: "Should TLS 1.0 be enabled on the frontend of the gateway?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"frontend_tls11_enabled": {
						Description: "Should TLS 1.1 be enabled on the frontend of the gateway?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"triple_des_ciphers_enabled": {
						Description: "Should the `TLS_RSA_WITH_3DES_EDE_CBC_SHA` cipher be enabled for alL TLS versions (1.0, 1.1 and 1.2)?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled": {
						Description: "Should the `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled": {
						Description: "Should the `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled": {
						Description: "Should the `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled": {
						Description: "Should the `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_rsa_with_aes128_gcm_sha256_ciphers_enabled": {
						Description: "Should the `TLS_RSA_WITH_AES_128_GCM_SHA256` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_rsa_with_aes256_cbc_sha256_ciphers_enabled": {
						Description: "Should the `TLS_RSA_WITH_AES_256_CBC_SHA256` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_rsa_with_aes128_cbc_sha256_ciphers_enabled": {
						Description: "Should the `TLS_RSA_WITH_AES_128_CBC_SHA256` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_rsa_with_aes256_gcm_sha384_ciphers_enabled": {
						Description: "Should the `TLS_RSA_WITH_AES_256_GCM_SHA384` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_rsa_with_aes256_cbc_sha_ciphers_enabled": {
						Description: "Should the `TLS_RSA_WITH_AES_256_CBC_SHA` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"tls_rsa_with_aes128_cbc_sha_ciphers_enabled": {
						Description: "Should the `TLS_RSA_WITH_AES_128_CBC_SHA` cipher be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
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
						Description: "Should anonymous users be redirected to the sign in page?",
						Type:        pluginsdk.TypeBool,
						Required:    true,
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
						Description: "Should subscription requests be delegated to an external url?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"user_registration_enabled": {
						Description: "Should user registration requests be delegated to an external url?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"url": {
						Description:  "The delegation URL.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
					"validation_key": {
						Description:  "A base64-encoded validation key to validate, that a request is coming from Azure API Management.",
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
						Description: "Can users sign up on the development portal?",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},

					"terms_of_service": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Description: "Should Terms of Service be displayed during sign up?.",
									Type:        pluginsdk.TypeBool,
									Required:    true,
								},
								"consent_required": {
									Description: "Should the user be asked for consent during sign up?",
									Type:        pluginsdk.TypeBool,
									Required:    true,
								},
								"text": {
									Description: "The Terms of Service which users are required to agree to in order to sign up.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
							},
						},
					},
				},
			},
		},

		"zones": commonschema.ZonesMultipleOptional(),

		"gateway_url": {
			Description: "The URL of the Gateway for the API Management Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"management_api_url": {
			Description: "The URL for the Management API associated with this API Management service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"gateway_regional_url": {
			Description: "The Region URL for the Gateway of the API Management Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"public_ip_addresses": {
			Description: "The Public IP addresses of the API Management Service.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"public_ip_address_id": {
			Description:  "ID of a standard SKU IPv4 Public IP.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidatePublicIPAddressID,
		},

		"public_network_access_enabled": {
			Description: "Is public access to the service allowed?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"private_ip_addresses": {
			Description: "The Private IP addresses of the API Management Service.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"portal_url": {
			Description: "The URL for the Publisher Portal associated with this API Management service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"developer_portal_url": {
			Description: "The URL for the Developer Portal associated with this API Management service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"scm_url": {
			Description: "The URL for the SCM (Source Code Management) Endpoint associated with this API Management service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tenant_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Should the access to the management API be enabled?",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"tenant_id": {
						Description: "The identifier for the tenant access information contract.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"primary_key": {
						Description: "Primary access key for the tenant access information contract.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"secondary_key": {
						Description: "Secondary access key for the tenant access information contract.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
