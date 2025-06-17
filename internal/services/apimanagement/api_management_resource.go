// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/delegationsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/deletedservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signinsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/signupsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/tenantaccess"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var (
	apimBackendProtocolSsl3                  = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30"
	apimBackendProtocolTls10                 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10"
	apimBackendProtocolTls11                 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11"
	apimFrontendProtocolSsl3                 = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Ssl30"
	apimFrontendProtocolTls10                = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10"
	apimFrontendProtocolTls11                = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11"
	apimTripleDesCiphers                     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TripleDes168"
	apimHttp2Protocol                        = "Microsoft.WindowsAzure.ApiManagement.Gateway.Protocols.Server.Http2"
	apimTlsEcdheEcdsaWithAes256CbcShaCiphers = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
	apimTlsEcdheEcdsaWithAes128CbcShaCiphers = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
	apimTlsEcdheRsaWithAes256CbcShaCiphers   = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	apimTlsEcdheRsaWithAes128CbcShaCiphers   = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	apimTlsRsaWithAes128GcmSha256Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_GCM_SHA256"
	apimTlsRsaWithAes256CbcSha256Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA256"
	apimTlsRsaWithAes256GcmSha384Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_GCM_SHA384"
	apimTlsRsaWithAes128CbcSha256Ciphers     = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA256"
	apimTlsRsaWithAes256CbcShaCiphers        = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_256_CBC_SHA"
	apimTlsRsaWithAes128CbcShaCiphers        = "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Ciphers.TLS_RSA_WITH_AES_128_CBC_SHA"
)

func resourceApiManagementService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementServiceCreate,
		Read:   resourceApiManagementServiceRead,
		Update: resourceApiManagementServiceUpdate,
		Delete: resourceApiManagementServiceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apimanagementservice.ParseServiceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(3 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(3 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(3 * time.Hour),
		},

		Schema: resourceApiManagementSchema(),

		// we can only change `virtual_network_type` from None to Internal Or External, Else the subnet can not be destroyed cause “InUseSubnetCannotBeDeleted” for 3 hours
		// we can not change the subnet from subnet1 to subnet2 either, Else the subnet1 can not be destroyed cause “InUseSubnetCannotBeDeleted” for 3 hours
		// Issue: https://github.com/Azure/azure-rest-api-specs/issues/10395
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("virtual_network_type", func(ctx context.Context, old, new, meta interface{}) bool {
				return !(old.(string) == string(apimanagementservice.VirtualNetworkTypeNone) &&
					(new.(string) == string(apimanagementservice.VirtualNetworkTypeInternal) ||
						new.(string) == string(apimanagementservice.VirtualNetworkTypeExternal)))
			}),

			pluginsdk.ForceNewIfChange("virtual_network_configuration", func(ctx context.Context, old, new, meta interface{}) bool {
				return !(len(old.([]interface{})) == 0 && len(new.([]interface{})) > 0)
			}),
		),
	}
}

func resourceApiManagementSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementName(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"publisher_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: apimValidate.ApiManagementServicePublisherName,
		},

		"publisher_email": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: apimValidate.ApiManagementServicePublisherEmail,
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: apimValidate.ApimSkuName(),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"virtual_network_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(apimanagementservice.VirtualNetworkTypeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(apimanagementservice.VirtualNetworkTypeNone),
				string(apimanagementservice.VirtualNetworkTypeExternal),
				string(apimanagementservice.VirtualNetworkTypeInternal),
			}, false),
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
			Computed: true,
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
						Computed:     true,
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
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(apimanagementservice.StoreNameCertificateAuthority),
							string(apimanagementservice.StoreNameRoot),
						}, false),
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
			Computed: true,
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
			Computed: true,
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
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"management": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: apiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"portal": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: apiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"developer_portal": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: apiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"proxy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: apiManagementResourceHostnameProxySchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
					"scm": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: apiManagementResourceHostnameSchema(),
						},
						AtLeastOneOf: []string{"hostname_configuration.0.management", "hostname_configuration.0.portal", "hostname_configuration.0.developer_portal", "hostname_configuration.0.proxy", "hostname_configuration.0.scm"},
					},
				},
			},
		},

		"sign_in": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
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
			Computed: true,
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
						ValidateFunc: validate.Base64EncodedString,
						Sensitive:    true,
					},
				},
			},
		},

		"sign_up": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
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
			Computed: true,
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

	if !features.FivePointOh() {
		s["protocols"].Elem.(*pluginsdk.Resource).Schema["enable_http2"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"protocols.0.http2_enabled"},
			Deprecated:    "`protocols.enable_http2` has been deprecated in favour of the `protocols.http2_enabled` property and will be removed in v5.0 of the AzureRM Provider",
		}
		s["protocols"].Elem.(*pluginsdk.Resource).Schema["http2_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"protocols.0.enable_http2"},
		}

		s["security"].Elem.(*pluginsdk.Resource).Schema["enable_backend_ssl30"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.backend_ssl30_enabled"},
			Deprecated:    "`security.enable_backend_ssl30` has been deprecated in favour of the `security.backend_ssl30_enabled` property and will be removed in v5.0 of the AzureRM Provider",
		}
		s["security"].Elem.(*pluginsdk.Resource).Schema["backend_ssl30_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.enable_backend_ssl30"},
		}

		s["security"].Elem.(*pluginsdk.Resource).Schema["enable_backend_tls10"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.backend_tls10_enabled"},
			Deprecated:    "`security.enable_backend_tls10` has been deprecated in favour of the `security.backend_tls10_enabled` property and will be removed in v5.0 of the AzureRM Provider",
		}
		s["security"].Elem.(*pluginsdk.Resource).Schema["backend_tls10_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.enable_backend_tls10"},
		}

		s["security"].Elem.(*pluginsdk.Resource).Schema["enable_backend_tls11"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.backend_tls11_enabled"},
			Deprecated:    "`security.enable_backend_tls11` has been deprecated in favour of the `security.backend_tls11_enabled` property and will be removed in v5.0 of the AzureRM Provider",
		}
		s["security"].Elem.(*pluginsdk.Resource).Schema["backend_tls11_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.enable_backend_tls11"},
		}

		s["security"].Elem.(*pluginsdk.Resource).Schema["enable_frontend_ssl30"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.frontend_ssl30_enabled"},
			Deprecated:    "`security.enable_frontend_ssl30` has been deprecated in favour of the `security.frontend_ssl30_enabled` property and will be removed in v5.0 of the AzureRM Provider",
		}
		s["security"].Elem.(*pluginsdk.Resource).Schema["frontend_ssl30_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.enable_frontend_ssl30"},
		}

		s["security"].Elem.(*pluginsdk.Resource).Schema["enable_frontend_tls10"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.frontend_tls10_enabled"},
			Deprecated:    "`security.enable_frontend_tls10` has been deprecated in favour of the `security.frontend_tls10_enabled` property and will be removed in v5.0 of the AzureRM Provider",
		}
		s["security"].Elem.(*pluginsdk.Resource).Schema["frontend_tls10_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.enable_frontend_tls10"},
		}

		s["security"].Elem.(*pluginsdk.Resource).Schema["enable_frontend_tls11"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.frontend_tls11_enabled"},
			Deprecated:    "`security.enable_frontend_tls11` has been deprecated in favour of the `security.frontend_tls11_enabled` property and will be removed in v5.0 of the AzureRM Provider",
		}
		s["security"].Elem.(*pluginsdk.Resource).Schema["frontend_tls11_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"security.0.enable_frontend_tls11"},
		}
	}

	return s
}

func resourceApiManagementServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	apiClient := meta.(*clients.Client).ApiManagement.ApiClient
	deletedServicesClient := meta.(*clients.Client).ApiManagement.DeletedServicesClient
	productsClient := meta.(*clients.Client).ApiManagement.ProductsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sku := expandAzureRmApiManagementSkuName(d.Get("sku_name").(string))

	log.Printf("[INFO] preparing arguments for API Management Service creation.")

	id := apimanagementservice.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_api_management", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	publicIpAddressId := d.Get("public_ip_address_id").(string)
	notificationSenderEmail := d.Get("notification_sender_email").(string)
	virtualNetworkType := d.Get("virtual_network_type").(string)

	customProperties, err := expandApiManagementCustomProperties(d, sku.Name == apimanagementservice.SkuTypeConsumption)
	if err != nil {
		return err
	}
	certificates := expandAzureRmApiManagementCertificates(d)

	publicNetworkAccess := apimanagementservice.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = apimanagementservice.PublicNetworkAccessDisabled
	}

	// before creating check to see if the resource exists in the soft delete state
	deletedServiceId := deletedservice.NewDeletedServiceID(id.SubscriptionId, location, id.ServiceName)
	softDeleted, err := deletedServicesClient.GetByName(ctx, deletedServiceId)
	if err != nil {
		// If Terraform lacks permission to read at the Subscription we'll get 403, not 404
		if !response.WasNotFound(softDeleted.HttpResponse) && !response.WasForbidden(softDeleted.HttpResponse) {
			return fmt.Errorf("checking for the presence of an existing Soft-Deleted API Management %q (Location %q): %+v", id.ServiceName, location, err)
		}
	}

	// if so, does the user want us to recover it?
	if !response.WasNotFound(softDeleted.HttpResponse) && !response.WasForbidden(softDeleted.HttpResponse) {
		if !meta.(*clients.Client).Features.ApiManagement.RecoverSoftDeleted {
			// this exists but the users opted out, so they must import this it out-of-band
			return errors.New(optedOutOfRecoveringSoftDeletedApiManagementErrorFmt(id.ServiceName, location))
		}

		// First recover the deleted API Management, since all other properties are ignored during a restore operation
		// (don't set the ID just yet to avoid tainting on failure)
		params := apimanagementservice.ApiManagementServiceResource{
			Location: location,
			Properties: apimanagementservice.ApiManagementServiceProperties{
				Restore: pointer.To(true),
			},
			Sku: sku,
		}

		// retry to restore service since there is an API issue : https://github.com/Azure/azure-rest-api-specs/issues/25262
		err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
			resp, err := client.CreateOrUpdate(ctx, id, params)
			if err != nil {
				if response.WasBadRequest(resp.HttpResponse) {
					return pluginsdk.RetryableError(err)
				}
				return pluginsdk.NonRetryableError(err)
			}
			if err := resp.Poller.PollUntilDone(ctx); err != nil {
				return pluginsdk.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("recovering %s: %+v", id, err)
		}
	}

	properties := apimanagementservice.ApiManagementServiceResource{
		Location: location,
		Properties: apimanagementservice.ApiManagementServiceProperties{
			PublisherName:       d.Get("publisher_name").(string),
			PublisherEmail:      d.Get("publisher_email").(string),
			PublicNetworkAccess: pointer.To(publicNetworkAccess),
			CustomProperties:    pointer.To(customProperties),
			Certificates:        certificates,
		},
		Sku:  sku,
		Tags: tags.Expand(t),
	}

	if _, ok := d.GetOk("hostname_configuration"); ok {
		properties.Properties.HostnameConfigurations = expandAzureRmApiManagementHostnameConfigurations(d)
	}

	// intentionally not gated since we specify a default value (of None) in the expand, which we need on updates
	identityRaw := d.Get("identity").([]interface{})
	identity, err := identity.ExpandSystemAndUserAssignedMap(identityRaw)
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	properties.Identity = identity

	if _, ok := d.GetOk("additional_location"); ok {
		var err error
		properties.Properties.AdditionalLocations, err = expandAzureRmApiManagementAdditionalLocations(d, sku)
		if err != nil {
			return err
		}
	}

	if notificationSenderEmail != "" {
		properties.Properties.NotificationSenderEmail = pointer.To(notificationSenderEmail)
	}

	if virtualNetworkType != "" {
		properties.Properties.VirtualNetworkType = pointer.To(apimanagementservice.VirtualNetworkType(virtualNetworkType))

		if virtualNetworkType != string(apimanagementservice.VirtualNetworkTypeNone) {
			virtualNetworkConfiguration := expandAzureRmApiManagementVirtualNetworkConfigurations(d)
			if virtualNetworkConfiguration == nil {
				return fmt.Errorf("you must specify 'virtual_network_configuration' when 'virtual_network_type' is %q", virtualNetworkType)
			}
			properties.Properties.VirtualNetworkConfiguration = virtualNetworkConfiguration
		}
	}

	if publicIpAddressId != "" {
		if sku.Name != apimanagementservice.SkuTypePremium && sku.Name != apimanagementservice.SkuTypeDeveloper {
			if virtualNetworkType == string(apimanagementservice.VirtualNetworkTypeNone) {
				return errors.New("`public_ip_address_id` is only supported when sku type is `Developer` or `Premium`, and the APIM instance is deployed in a virtual network")
			}
		}
		properties.Properties.PublicIPAddressId = pointer.To(publicIpAddressId)
	}

	if d.HasChange("client_certificate_enabled") {
		enableClientCertificate := d.Get("client_certificate_enabled").(bool)
		if enableClientCertificate && sku.Name != apimanagementservice.SkuTypeConsumption {
			return errors.New("`client_certificate_enabled` is only supported when sku type is `Consumption`")
		}
		properties.Properties.EnableClientCertificate = pointer.To(enableClientCertificate)
	}

	gateWayDisabled := d.Get("gateway_disabled").(bool)
	if gateWayDisabled && len(*properties.Properties.AdditionalLocations) == 0 {
		return errors.New("`gateway_disabled` is only supported when `additional_location` is set")
	}
	properties.Properties.DisableGateway = pointer.To(gateWayDisabled)

	if v, ok := d.GetOk("min_api_version"); ok {
		properties.Properties.ApiVersionConstraint = &apimanagementservice.ApiVersionConstraint{
			MinApiVersion: pointer.To(v.(string)),
		}
	}

	if v := d.Get("zones").(*schema.Set).List(); len(v) > 0 {
		if sku.Name != apimanagementservice.SkuTypePremium {
			return errors.New("`zones` is only supported when sku type is `Premium`")
		}

		zones := zones.ExpandUntyped(v)
		properties.Zones = &zones
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, properties); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// Remove sample products and APIs after creating (v3.0 behaviour)
	apiServiceId := api.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)

	listResp, err := apiClient.ListByService(ctx, apiServiceId, api.ListByServiceOperationOptions{})
	if err != nil {
		return fmt.Errorf("listing APIs after creation of %s: %+v", id, err)
	}
	if model := listResp.Model; model != nil {
		for _, contract := range *model {
			if contract.Id == nil {
				continue
			}
			apiId, err := api.ParseApiID(pointer.From(contract.Id))
			if err != nil {
				return fmt.Errorf("parsing API ID: %+v", err)
			}
			log.Printf("[DEBUG] Deleting %s", apiId)
			if delResp, err := apiClient.Delete(ctx, *apiId, api.DeleteOperationOptions{DeleteRevisions: pointer.To(true)}); err != nil {
				if !response.WasNotFound(delResp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *apiId, err)
				}
			}
		}
	}

	produceServiceId := product.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
	proListResp, err := productsClient.ListByService(ctx, produceServiceId, product.ListByServiceOperationOptions{})
	if err != nil {
		return fmt.Errorf("listing products after creation of %s: %+v", id, err)
	}
	if model := proListResp.Model; model != nil {
		for _, contract := range *model {
			if contract.Id == nil {
				continue
			}
			productId, err := product.ParseProductID(pointer.From(contract.Id))
			if err != nil {
				return fmt.Errorf("parsing product ID: %+v", err)
			}
			log.Printf("[DEBUG] Deleting %s", productId)
			if delResp, err := productsClient.Delete(ctx, *productId, product.DeleteOperationOptions{DeleteSubscriptions: pointer.To(true)}); err != nil {
				if !response.WasNotFound(delResp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *productId, err)
				}
			}
		}
	}

	signInSettingsRaw := d.Get("sign_in").([]interface{})
	if sku.Name == apimanagementservice.SkuTypeConsumption && len(signInSettingsRaw) > 0 {
		return errors.New("`sign_in` is not supported for sku tier `Consumption`")
	}
	if sku.Name != apimanagementservice.SkuTypeConsumption {
		signInSettingServiceId := signinsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
		signInSettings := expandApiManagementSignInSettings(signInSettingsRaw)
		signInClient := meta.(*clients.Client).ApiManagement.SignInClient
		if _, err := signInClient.CreateOrUpdate(ctx, signInSettingServiceId, signInSettings, signinsettings.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf(" setting Sign In settings for %s: %+v", id, err)
		}
	}

	signUpSettingsRaw := d.Get("sign_up").([]interface{})
	if sku.Name == apimanagementservice.SkuTypeConsumption && len(signUpSettingsRaw) > 0 {
		return fmt.Errorf("`sign_up` is not supported for sku tier `Consumption`")
	}
	if sku.Name != apimanagementservice.SkuTypeConsumption {
		signUpSettingServiceId := signupsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
		signUpSettings := expandApiManagementSignUpSettings(signUpSettingsRaw)
		signUpClient := meta.(*clients.Client).ApiManagement.SignUpClient
		if _, err := signUpClient.CreateOrUpdate(ctx, signUpSettingServiceId, signUpSettings, signupsettings.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf(" setting Sign Up settings for %s: %+v", id, err)
		}
	}

	delegationSettingsRaw := d.Get("delegation").([]interface{})
	if sku.Name == apimanagementservice.SkuTypeConsumption && len(delegationSettingsRaw) > 0 {
		return fmt.Errorf("`delegation` is not supported for sku tier `Consumption`")
	}
	if sku.Name != apimanagementservice.SkuTypeConsumption && len(delegationSettingsRaw) > 0 {
		delegationSettingServiceId := delegationsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
		delegationSettings := expandApiManagementDelegationSettings(delegationSettingsRaw)
		delegationClient := meta.(*clients.Client).ApiManagement.DelegationSettingsClient
		if _, err := delegationClient.CreateOrUpdate(ctx, delegationSettingServiceId, delegationSettings, delegationsettings.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf(" setting Delegation settings for %s: %+v", id, err)
		}
	}

	tenantAccessRaw := d.Get("tenant_access").([]interface{})
	if sku.Name == apimanagementservice.SkuTypeConsumption && len(tenantAccessRaw) > 0 {
		return fmt.Errorf("`tenant_access` is not supported for sku tier `Consumption`")
	}
	if sku.Name != apimanagementservice.SkuTypeConsumption && d.HasChange("tenant_access") {
		tenantAccessServiceId := tenantaccess.NewAccessID(subscriptionId, id.ResourceGroupName, id.ServiceName, "access")
		tenantAccessInformationParametersRaw := d.Get("tenant_access").([]interface{})
		tenantAccessInformationParameters := expandApiManagementTenantAccessSettings(tenantAccessInformationParametersRaw)
		tenantAccessClient := meta.(*clients.Client).ApiManagement.TenantAccessClient
		if _, err := tenantAccessClient.Update(ctx, tenantAccessServiceId, tenantAccessInformationParameters, tenantaccess.UpdateOperationOptions{}); err != nil {
			return fmt.Errorf(" updating tenant access settings for %s: %+v", id, err)
		}
	}

	return resourceApiManagementServiceRead(d, meta)
}

func resourceApiManagementServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sku := expandAzureRmApiManagementSkuName(d.Get("sku_name").(string))
	virtualNetworkType := d.Get("virtual_network_type").(string)
	virtualNetworkConfiguration := expandAzureRmApiManagementVirtualNetworkConfigurations(d)

	log.Printf("[INFO] preparing arguments for API Management Service creation.")

	id := apimanagementservice.NewServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	_, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
	}

	props := apimanagementservice.ApiManagementServiceUpdateProperties{}
	payload := apimanagementservice.ApiManagementServiceUpdateParameters{}

	if d.HasChange("sku_name") {
		payload.Sku = pointer.To(sku)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("public_ip_address_id") {
		publicIpAddressId := d.Get("public_ip_address_id").(string)
		if publicIpAddressId != "" {
			if sku.Name != apimanagementservice.SkuTypePremium && sku.Name != apimanagementservice.SkuTypeDeveloper {
				if d.Get("virtual_network_type").(string) == string(apimanagementservice.VirtualNetworkTypeNone) {
					return fmt.Errorf("`public_ip_address_id` is only supported when sku type is `Developer` or `Premium`, and the APIM instance is deployed in a virtual network.")
				}
			}
			props.PublicIPAddressId = pointer.To(publicIpAddressId)
		}
	}

	if d.HasChange("notification_sender_email") {
		props.NotificationSenderEmail = pointer.To(d.Get("notification_sender_email").(string))
	}

	if d.HasChange("virtual_network_type") {
		props.VirtualNetworkType = pointer.To(apimanagementservice.VirtualNetworkType(virtualNetworkType))
		if virtualNetworkType != string(apimanagementservice.VirtualNetworkTypeNone) {
			if virtualNetworkConfiguration == nil {
				return fmt.Errorf("You must specify 'virtual_network_configuration' when 'virtual_network_type' is %q", virtualNetworkType)
			}
			props.VirtualNetworkConfiguration = virtualNetworkConfiguration
		}
	}

	if d.HasChange("virtual_network_configuration") {
		props.VirtualNetworkConfiguration = virtualNetworkConfiguration
		if virtualNetworkType == string(apimanagementservice.VirtualNetworkTypeNone) {
			if virtualNetworkConfiguration != nil {
				return fmt.Errorf("You must specify 'virtual_network_type' when specifying 'virtual_network_configuration'")
			}
		}
	}

	if d.HasChanges("security", "protocols") {
		customProperties, err := expandApiManagementCustomProperties(d, sku.Name == apimanagementservice.SkuTypeConsumption)
		if err != nil {
			return err
		}
		props.CustomProperties = pointer.To(customProperties)
	}

	if d.HasChange("certificate") {
		props.Certificates = expandAzureRmApiManagementCertificates(d)
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := apimanagementservice.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = apimanagementservice.PublicNetworkAccessDisabled
		}

		props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("publisher_name") {
		props.PublisherName = pointer.To(d.Get("publisher_name").(string))
	}

	if d.HasChange("publisher_email") {
		props.PublisherEmail = pointer.To(d.Get("publisher_email").(string))
	}

	if d.HasChange("hostname_configuration") {
		props.HostnameConfigurations = expandAzureRmApiManagementHostnameConfigurations(d)
	}

	// intentionally not gated since we specify a default value (of None) in the expand, which we need on updates
	if d.HasChange("identity") {
		identityRaw := d.Get("identity").([]interface{})
		identity, err := identity.ExpandSystemAndUserAssignedMap(identityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = identity
	}

	if d.HasChange("additional_location") {
		props.AdditionalLocations, err = expandAzureRmApiManagementAdditionalLocations(d, sku)
		if err != nil {
			return err
		}
	}

	if d.HasChange("client_certificate_enabled") {
		enableClientCertificate := d.Get("client_certificate_enabled").(bool)
		if enableClientCertificate && sku.Name != apimanagementservice.SkuTypeConsumption {
			return errors.New("`client_certificate_enabled` is only supported when sku type is `Consumption`")
		}
		props.EnableClientCertificate = pointer.To(enableClientCertificate)
	}

	if d.HasChange("gateway_disabled") {
		gateWayDisabled := d.Get("gateway_disabled").(bool)
		if gateWayDisabled && props.AdditionalLocations != nil && len(*props.AdditionalLocations) == 0 {
			return errors.New("`gateway_disabled` is only supported when `additional_location` is set")
		}
		props.DisableGateway = pointer.To(gateWayDisabled)
	}

	if d.HasChange("min_api_version") {
		props.ApiVersionConstraint = &apimanagementservice.ApiVersionConstraint{
			MinApiVersion: nil,
		}

		if v, ok := d.GetOk("min_api_version"); ok {
			props.ApiVersionConstraint.MinApiVersion = pointer.To(v.(string))
		}
	}

	if d.HasChange("zones") {
		if v := d.Get("zones").(*schema.Set).List(); len(v) > 0 {
			if sku.Name != apimanagementservice.SkuTypePremium {
				return errors.New("`zones` is only supported when sku type is `Premium`")
			}

			zones := zones.ExpandUntyped(v)
			payload.Zones = &zones
		}
	}

	payload.Properties = pointer.To(props)

	if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if d.HasChange("sign_in") {
		signInSettingsRaw := d.Get("sign_in").([]interface{})
		if sku.Name == apimanagementservice.SkuTypeConsumption && len(signInSettingsRaw) > 0 {
			return errors.New("`sign_in` is not supported for sku tier `Consumption`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption {
			signInSettingServiceId := signinsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
			signInSettings := expandApiManagementSignInSettings(signInSettingsRaw)
			signInClient := meta.(*clients.Client).ApiManagement.SignInClient
			if _, err := signInClient.CreateOrUpdate(ctx, signInSettingServiceId, signInSettings, signinsettings.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" setting Sign In settings for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("sign_up") {
		signUpSettingsRaw := d.Get("sign_up").([]interface{})
		if sku.Name == apimanagementservice.SkuTypeConsumption && len(signUpSettingsRaw) > 0 {
			return errors.New("`sign_up` is not supported for sku tier `Consumption`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption {
			signUpSettingServiceId := signupsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
			signUpSettings := expandApiManagementSignUpSettings(signUpSettingsRaw)
			signUpClient := meta.(*clients.Client).ApiManagement.SignUpClient
			if _, err := signUpClient.CreateOrUpdate(ctx, signUpSettingServiceId, signUpSettings, signupsettings.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" setting Sign Up settings for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("delegation") {
		delegationSettingsRaw := d.Get("delegation").([]interface{})
		if sku.Name == apimanagementservice.SkuTypeConsumption && len(delegationSettingsRaw) > 0 {
			return errors.New("`delegation` is not supported for sku tier `Consumption`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption && len(delegationSettingsRaw) > 0 {
			delegationSettingServiceId := delegationsettings.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)
			delegationSettings := expandApiManagementDelegationSettings(delegationSettingsRaw)
			delegationClient := meta.(*clients.Client).ApiManagement.DelegationSettingsClient
			if _, err := delegationClient.CreateOrUpdate(ctx, delegationSettingServiceId, delegationSettings, delegationsettings.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" setting Delegation settings for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("tenant_access") {
		tenantAccessRaw := d.Get("tenant_access").([]interface{})
		if sku.Name == apimanagementservice.SkuTypeConsumption && len(tenantAccessRaw) > 0 {
			return fmt.Errorf("`tenant_access` is not supported for sku tier `Consumption`")
		}
		if sku.Name != apimanagementservice.SkuTypeConsumption && d.HasChange("tenant_access") {
			tenantAccessServiceId := tenantaccess.NewAccessID(subscriptionId, id.ResourceGroupName, id.ServiceName, "access")
			tenantAccessInformationParametersRaw := d.Get("tenant_access").([]interface{})
			tenantAccessInformationParameters := expandApiManagementTenantAccessSettings(tenantAccessInformationParametersRaw)
			tenantAccessClient := meta.(*clients.Client).ApiManagement.TenantAccessClient
			if _, err := tenantAccessClient.Update(ctx, tenantAccessServiceId, tenantAccessInformationParameters, tenantaccess.UpdateOperationOptions{}); err != nil {
				return fmt.Errorf(" updating tenant access settings for %s: %+v", id, err)
			}
		}
	}

	return resourceApiManagementServiceRead(d, meta)
}

func resourceApiManagementServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	signInClient := meta.(*clients.Client).ApiManagement.SignInClient
	signUpClient := meta.(*clients.Client).ApiManagement.SignUpClient
	delegationClient := meta.(*clients.Client).ApiManagement.DelegationSettingsClient
	tenantAccessClient := meta.(*clients.Client).ApiManagement.TenantAccessClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apimanagementservice.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	policyServiceId := policy.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
	policyClient := meta.(*clients.Client).ApiManagement.PolicyClient
	policy, err := policyClient.Get(ctx, policyServiceId, policy.GetOperationOptions{Format: pointer.To(policy.PolicyExportFormatXml)})
	if err != nil {
		if !response.WasNotFound(policy.HttpResponse) {
			return fmt.Errorf("retrieving Policy for %s: %+v", *id, err)
		}
	}

	d.Set("name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))
		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("publisher_email", model.Properties.PublisherEmail)
		d.Set("publisher_name", model.Properties.PublisherName)
		d.Set("notification_sender_email", pointer.From(model.Properties.NotificationSenderEmail))
		d.Set("gateway_url", pointer.From(model.Properties.GatewayURL))
		d.Set("gateway_regional_url", pointer.From(model.Properties.GatewayRegionalURL))
		d.Set("portal_url", pointer.From(model.Properties.PortalURL))
		d.Set("developer_portal_url", pointer.From(model.Properties.DeveloperPortalURL))
		d.Set("management_api_url", pointer.From(model.Properties.ManagementApiURL))
		d.Set("scm_url", pointer.From(model.Properties.ScmURL))
		d.Set("public_ip_addresses", pointer.From(model.Properties.PublicIPAddresses))
		d.Set("public_ip_address_id", pointer.From(model.Properties.PublicIPAddressId))
		d.Set("public_network_access_enabled", pointer.From(model.Properties.PublicNetworkAccess) == apimanagementservice.PublicNetworkAccessEnabled)
		d.Set("private_ip_addresses", pointer.From(model.Properties.PrivateIPAddresses))
		d.Set("virtual_network_type", pointer.From(model.Properties.VirtualNetworkType))
		d.Set("client_certificate_enabled", pointer.From(model.Properties.EnableClientCertificate))
		d.Set("gateway_disabled", pointer.From(model.Properties.DisableGateway))

		d.Set("certificate", flattenAPIManagementCertificates(d, model.Properties.Certificates))

		if model.Sku.Name != "" {
			if err := d.Set("security", flattenApiManagementSecurityCustomProperties(*model.Properties.CustomProperties, model.Sku.Name == apimanagementservice.SkuTypeConsumption)); err != nil {
				return fmt.Errorf("setting `security`: %+v", err)
			}
		}

		if err := d.Set("protocols", flattenApiManagementProtocolsCustomProperties(*model.Properties.CustomProperties)); err != nil {
			return fmt.Errorf("setting `protocols`: %+v", err)
		}

		hostnameConfigs := flattenApiManagementHostnameConfigurations(model.Properties.HostnameConfigurations, d)
		if err := d.Set("hostname_configuration", hostnameConfigs); err != nil {
			return fmt.Errorf("setting `hostname_configuration`: %+v", err)
		}
		additionalLocation, err := flattenApiManagementAdditionalLocations(model.Properties.AdditionalLocations)
		if err != nil {
			return err
		}
		if err := d.Set("additional_location", additionalLocation); err != nil {
			return fmt.Errorf("setting `additional_location`: %+v", err)
		}

		virtualNetworkConfiguration, err := flattenApiManagementVirtualNetworkConfiguration(model.Properties.VirtualNetworkConfiguration)
		if err != nil {
			return err
		}
		if err := d.Set("virtual_network_configuration", virtualNetworkConfiguration); err != nil {
			return fmt.Errorf("setting `virtual_network_configuration`: %+v", err)
		}

		var minApiVersion string
		if model.Properties.ApiVersionConstraint != nil {
			minApiVersion = pointer.From(model.Properties.ApiVersionConstraint.MinApiVersion)
		}
		d.Set("min_api_version", minApiVersion)

		if err := d.Set("sku_name", flattenApiManagementServiceSkuName(&model.Sku)); err != nil {
			return fmt.Errorf("setting `sku_name`: %+v", err)
		}
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		if model.Sku.Name != apimanagementservice.SkuTypeConsumption {
			signInSettingServiceId := signinsettings.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
			signInSettings, err := signInClient.Get(ctx, signInSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Sign In Settings for %s: %+v", *id, err)
			}
			if err := d.Set("sign_in", flattenApiManagementSignInSettings(*signInSettings.Model)); err != nil {
				return fmt.Errorf("setting `sign_in`: %+v", err)
			}

			signUpSettingServiceId := signupsettings.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
			signUpSettings, err := signUpClient.Get(ctx, signUpSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Sign Up Settings for %s: %+v", *id, err)
			}

			if err := d.Set("sign_up", flattenApiManagementSignUpSettings(*signUpSettings.Model)); err != nil {
				return fmt.Errorf("setting `sign_up`: %+v", err)
			}

			delegationSettingServiceId := delegationsettings.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
			delegationSettings, err := delegationClient.Get(ctx, delegationSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Delegation Settings for %s: %+v", *id, err)
			}

			delegationValidationKeyContract, err := delegationClient.ListSecrets(ctx, delegationSettingServiceId)
			if err != nil {
				return fmt.Errorf("retrieving Delegation Validation Key for %s: %+v", *id, err)
			}

			if err := d.Set("delegation", flattenApiManagementDelegationSettings(*delegationSettings.Model, *delegationValidationKeyContract.Model)); err != nil {
				return fmt.Errorf("setting `delegation`: %+v", err)
			}

			tenantAccessServiceId := tenantaccess.NewAccessID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, "access")
			tenantAccessInformationContract, err := tenantAccessClient.ListSecrets(ctx, tenantAccessServiceId)
			if err != nil {
				return fmt.Errorf("retrieving tenant access properties for %s: %+v", *id, err)
			}
			if err := d.Set("tenant_access", flattenApiManagementTenantAccessSettings(*tenantAccessInformationContract.Model)); err != nil {
				return fmt.Errorf("setting `tenant_access`: %+v", err)
			}
		} else {
			d.Set("sign_in", []interface{}{})
			d.Set("sign_up", []interface{}{})
			d.Set("delegation", []interface{}{})
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceApiManagementServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ServiceClient
	deletedServicesClient := meta.(*clients.Client).ApiManagement.DeletedServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apimanagementservice.ParseServiceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if model := existing.Model; model != nil {
		locationName := location.NormalizeNilable(pointer.To(model.Location))

		// Purge the soft deleted Api Management permanently if the feature flag is enabled
		if meta.(*clients.Client).Features.ApiManagement.PurgeSoftDeleteOnDestroy {
			log.Printf("[DEBUG] %s marked for purge - executing purge", *id)
			deletedServiceId := deletedservice.NewDeletedServiceID(id.SubscriptionId, locationName, id.ServiceName)
			if _, err := deletedServicesClient.GetByName(ctx, deletedServiceId); err != nil {
				return fmt.Errorf("retrieving the deleted %s to be able to purge it: %+v", *id, err)
			}
			resp, err := deletedServicesClient.Purge(ctx, deletedServiceId)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("purging the deleted %s: %+v", *id, err)
			}

			if !response.WasNotFound(resp.HttpResponse) {
				if err := resp.Poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("purging the deleted %s: %+v", *id, err)
				}
			}

			log.Printf("[DEBUG] Purged %s.", *id)
			return nil
		}
	}

	return nil
}

func apiManagementRefreshFunc(ctx context.Context, client *apimanagementservice.ApiManagementServiceClient, id apimanagementservice.ServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if API Management Service %q (Resource Group: %q) is available..", id.ServiceName, id.ResourceGroupName)

		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				log.Printf("[DEBUG] Retrieving API Management %q (Resource Group: %q) returned 404.", id.ServiceName, id.ResourceGroupName)
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the state of the API Management Service %q (Resource Group: %q): %+v", id.ServiceName, id.ResourceGroupName, err)
		}

		state := ""
		if model := resp.Model; model != nil {
			if provisioningState := model.Properties.ProvisioningState; provisioningState != nil {
				state = pointer.From(provisioningState)
			}
		}

		return resp, state, nil
	}
}

func expandAzureRmApiManagementHostnameConfigurations(d *pluginsdk.ResourceData) *[]apimanagementservice.HostnameConfiguration {
	results := make([]apimanagementservice.HostnameConfiguration, 0)
	vs := d.Get("hostname_configuration")
	if vs == nil {
		return &results
	}
	hostnameVs := vs.([]interface{})

	for _, hostnameRawVal := range hostnameVs {
		// hostnameRawVal is guaranteed to be non-nil as there is AtLeastOneOf constraint on its containing properties.
		hostnameV := hostnameRawVal.(map[string]interface{})

		managementVs := hostnameV["management"].([]interface{})
		for _, managementV := range managementVs {
			v := managementV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeManagement)
			results = append(results, output)
		}

		portalVs := hostnameV["portal"].([]interface{})
		for _, portalV := range portalVs {
			v := portalV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypePortal)
			results = append(results, output)
		}

		developerPortalVs := hostnameV["developer_portal"].([]interface{})
		for _, developerPortalV := range developerPortalVs {
			v := developerPortalV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeDeveloperPortal)
			results = append(results, output)
		}

		proxyVs := hostnameV["proxy"].([]interface{})
		for _, proxyV := range proxyVs {
			v := proxyV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeProxy)
			if value, ok := v["default_ssl_binding"]; ok {
				output.DefaultSslBinding = pointer.To(value.(bool))
			}
			results = append(results, output)
		}

		scmVs := hostnameV["scm"].([]interface{})
		for _, scmV := range scmVs {
			v := scmV.(map[string]interface{})
			output := expandApiManagementCommonHostnameConfiguration(v, apimanagementservice.HostnameTypeScm)
			results = append(results, output)
		}
	}

	return &results
}

func expandApiManagementCommonHostnameConfiguration(input map[string]interface{}, hostnameType apimanagementservice.HostnameType) apimanagementservice.HostnameConfiguration {
	output := apimanagementservice.HostnameConfiguration{
		Type: hostnameType,
	}
	if v, ok := input["certificate"]; ok && v.(string) != "" {
		output.EncodedCertificate = pointer.To(v.(string))
	}
	if v, ok := input["certificate_password"]; ok && v.(string) != "" {
		output.CertificatePassword = pointer.To(v.(string))
	}
	if v, ok := input["host_name"]; ok && v.(string) != "" {
		output.HostName = v.(string)
	}
	if v, ok := input["key_vault_certificate_id"]; ok && v.(string) != "" {
		output.KeyVaultId = pointer.To(v.(string))
	}
	if !features.FivePointOh() {
		if v, ok := input["key_vault_id"]; ok && v.(string) != "" {
			output.KeyVaultId = pointer.To(v.(string))
		}
	}

	if v, ok := input["negotiate_client_certificate"]; ok {
		output.NegotiateClientCertificate = pointer.To(v.(bool))
	}

	if v, ok := input["ssl_keyvault_identity_client_id"]; ok && v.(string) != "" {
		output.IdentityClientId = pointer.To(v.(string))
	}

	return output
}

func flattenApiManagementHostnameConfigurations(input *[]apimanagementservice.HostnameConfiguration, d *pluginsdk.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	managementResults := make([]interface{}, 0)
	portalResults := make([]interface{}, 0)
	developerPortalResults := make([]interface{}, 0)
	proxyResults := make([]interface{}, 0)
	scmResults := make([]interface{}, 0)

	for _, config := range *input {
		output := make(map[string]interface{})

		output["host_name"] = config.HostName
		output["negotiate_client_certificate"] = pointer.From(config.NegotiateClientCertificate)
		output["key_vault_certificate_id"] = pointer.From(config.KeyVaultId)
		output["ssl_keyvault_identity_client_id"] = pointer.From(config.IdentityClientId)

		if !features.FivePointOh() {
			output["key_vault_id"] = pointer.From(config.KeyVaultId)
		}

		if config.Certificate != nil {
			if config.Certificate.Expiry != "" {
				output["expiry"] = config.Certificate.Expiry
			}
			output["thumbprint"] = config.Certificate.Thumbprint
			output["subject"] = config.Certificate.Subject
		}

		output["certificate_source"] = pointer.From(config.CertificateSource)
		output["certificate_status"] = pointer.From(config.CertificateStatus)

		var configType string
		switch strings.ToLower(string(config.Type)) {
		case strings.ToLower(string(apimanagementservice.HostnameTypeProxy)):
			// only set SSL binding for proxy types
			output["default_ssl_binding"] = pointer.From(config.DefaultSslBinding)
			proxyResults = append(proxyResults, output)
			configType = "proxy"

		case strings.ToLower(string(apimanagementservice.HostnameTypeManagement)):
			managementResults = append(managementResults, output)
			configType = "management"

		case strings.ToLower(string(apimanagementservice.HostnameTypePortal)):
			portalResults = append(portalResults, output)
			configType = "portal"

		case strings.ToLower(string(apimanagementservice.HostnameTypeDeveloperPortal)):
			developerPortalResults = append(developerPortalResults, output)
			configType = "developer_portal"

		case strings.ToLower(string(apimanagementservice.HostnameTypeScm)):
			scmResults = append(scmResults, output)
			configType = "scm"
		}

		existingHostnames := d.Get("hostname_configuration").([]interface{})
		if len(existingHostnames) > 0 && configType != "" {
			v := existingHostnames[0].(map[string]interface{})

			if valsRaw, ok := v[configType]; ok {
				vals := valsRaw.([]interface{})
				schemaz.CopyCertificateAndPassword(vals, config.HostName, output)
			}
		}
	}

	if len(managementResults) == 0 && len(portalResults) == 0 && len(developerPortalResults) == 0 && len(proxyResults) == 0 && len(scmResults) == 0 {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"management":       managementResults,
			"portal":           portalResults,
			"developer_portal": developerPortalResults,
			"proxy":            proxyResults,
			"scm":              scmResults,
		},
	}
}

func expandAzureRmApiManagementCertificates(d *pluginsdk.ResourceData) *[]apimanagementservice.CertificateConfiguration {
	vs := d.Get("certificate").([]interface{})

	results := make([]apimanagementservice.CertificateConfiguration, 0)

	for _, v := range vs {
		config := v.(map[string]interface{})

		certBase64 := config["encoded_certificate"].(string)
		storeName := apimanagementservice.StoreName(config["store_name"].(string))

		cert := apimanagementservice.CertificateConfiguration{
			EncodedCertificate: pointer.To(certBase64),
			StoreName:          storeName,
		}

		cert.CertificatePassword = pointer.To(config["certificate_password"].(string))

		results = append(results, cert)
	}

	return &results
}

func expandAzureRmApiManagementAdditionalLocations(d *pluginsdk.ResourceData, sku apimanagementservice.ApiManagementServiceSkuProperties) (*[]apimanagementservice.AdditionalLocation, error) {
	inputLocations := d.Get("additional_location").([]interface{})
	parentVnetConfig := d.Get("virtual_network_configuration").([]interface{})

	additionalLocations := make([]apimanagementservice.AdditionalLocation, 0)

	for _, v := range inputLocations {
		config := v.(map[string]interface{})
		location := azure.NormalizeLocation(config["location"].(string))

		if config["capacity"].(int) > 0 {
			sku.Capacity = int64(config["capacity"].(int))
		}

		additionalLocation := apimanagementservice.AdditionalLocation{
			Location:       location,
			Sku:            sku,
			DisableGateway: pointer.To(config["gateway_disabled"].(bool)),
		}

		childVnetConfig := config["virtual_network_configuration"].([]interface{})
		switch {
		case len(childVnetConfig) == 0 && len(parentVnetConfig) > 0:
			return nil, errors.New("`virtual_network_configuration` must be specified in any `additional_location` block when top-level `virtual_network_configuration` is supplied")
		case len(childVnetConfig) > 0 && len(parentVnetConfig) == 0:
			return nil, errors.New("`virtual_network_configuration` must be empty in all `additional_location` blocks when top-level `virtual_network_configuration` is not supplied")
		case len(childVnetConfig) > 0 && len(parentVnetConfig) > 0:
			v := childVnetConfig[0].(map[string]interface{})
			subnetResourceId := v["subnet_id"].(string)
			additionalLocation.VirtualNetworkConfiguration = &apimanagementservice.VirtualNetworkConfiguration{
				SubnetResourceId: pointer.To(subnetResourceId),
			}
		}

		publicIPAddressID := config["public_ip_address_id"].(string)
		if publicIPAddressID != "" {
			if sku.Name != apimanagementservice.SkuTypePremium {
				if len(childVnetConfig) == 0 {
					return nil, errors.New("`public_ip_address_id` for an additional location is only supported when sku type is `Premium`, and the APIM instance is deployed in a virtual network.")
				}
			}
			additionalLocation.PublicIPAddressId = &publicIPAddressID
		}

		zones := zones.ExpandUntyped(config["zones"].(*schema.Set).List())
		if len(zones) > 0 {
			additionalLocation.Zones = &zones
		}

		additionalLocations = append(additionalLocations, additionalLocation)
	}

	return &additionalLocations, nil
}

func flattenApiManagementAdditionalLocations(input *[]apimanagementservice.AdditionalLocation) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, prop := range *input {
		virtualNetworkConfiguration, err := flattenApiManagementVirtualNetworkConfiguration(prop.VirtualNetworkConfiguration)
		if err != nil {
			return results, err
		}

		results = append(results, map[string]interface{}{
			"capacity":                      int32(prop.Sku.Capacity),
			"gateway_regional_url":          pointer.From(prop.GatewayRegionalURL),
			"location":                      location.NormalizeNilable(pointer.To(prop.Location)),
			"private_ip_addresses":          pointer.From(prop.PrivateIPAddresses),
			"public_ip_address_id":          pointer.From(prop.PublicIPAddressId),
			"public_ip_addresses":           pointer.From(prop.PublicIPAddresses),
			"virtual_network_configuration": virtualNetworkConfiguration,
			"zones":                         zones.FlattenUntyped(prop.Zones),
			"gateway_disabled":              pointer.From(prop.DisableGateway),
		})
	}

	return results, nil
}

func expandAzureRmApiManagementSkuName(input string) apimanagementservice.ApiManagementServiceSkuProperties {
	// "sku_name" is validated to be in this format above, and is required
	skuParts := strings.Split(input, "_")
	name := skuParts[0]
	capacity, _ := strconv.Atoi(skuParts[1])
	return apimanagementservice.ApiManagementServiceSkuProperties{
		Name:     apimanagementservice.SkuType(name),
		Capacity: int64(capacity),
	}
}

func flattenApiManagementServiceSkuName(input *apimanagementservice.ApiManagementServiceSkuProperties) string {
	if input == nil {
		return ""
	}

	return fmt.Sprintf("%s_%d", string(input.Name), input.Capacity)
}

func expandApiManagementCustomProperties(d *pluginsdk.ResourceData, skuIsConsumption bool) (map[string]string, error) {
	backendProtocolSsl3 := false
	backendProtocolTls10 := false
	backendProtocolTls11 := false
	frontendProtocolSsl3 := false
	frontendProtocolTls10 := false
	frontendProtocolTls11 := false
	tripleDesCiphers := false
	tlsEcdheEcdsaWithAes256CbcShaCiphers := false
	tlsEcdheEcdsaWithAes128CbcShaCiphers := false
	tlsEcdheRsaWithAes256CbcShaCiphers := false
	tlsEcdheRsaWithAes128CbcShaCiphers := false
	tlsRsaWithAes128GcmSha256Ciphers := false
	tlsRsaWithAes256GcmSha384Ciphers := false
	tlsRsaWithAes256CbcSha256Ciphers := false
	tlsRsaWithAes128CbcSha256Ciphers := false
	tlsRsaWithAes256CbcShaCiphers := false
	tlsRsaWithAes128CbcShaCiphers := false

	if vs := d.Get("security").([]interface{}); len(vs) > 0 {
		v := vs[0].(map[string]interface{})
		backendProtocolSsl3 = v["backend_ssl30_enabled"].(bool)
		backendProtocolTls10 = v["backend_tls10_enabled"].(bool)
		backendProtocolTls11 = v["backend_tls11_enabled"].(bool)
		frontendProtocolSsl3 = v["frontend_ssl30_enabled"].(bool)
		frontendProtocolTls10 = v["frontend_tls10_enabled"].(bool)
		frontendProtocolTls11 = v["frontend_tls11_enabled"].(bool)

		if !features.FivePointOh() {
			if val, ok := d.GetOk("security.0.enable_backend_ssl30"); ok {
				backendProtocolSsl3 = val.(bool)
			}
			if val, ok := d.GetOk("security.0.enable_backend_tls10"); ok {
				backendProtocolTls10 = val.(bool)
			}
			if val, ok := d.GetOk("security.0.enable_backend_tls11"); ok {
				backendProtocolTls11 = val.(bool)
			}
			if val, ok := d.GetOk("security.0.enable_frontend_ssl30"); ok {
				frontendProtocolSsl3 = val.(bool)
			}
			if val, ok := d.GetOk("security.0.enable_frontend_tls10"); ok {
				frontendProtocolTls10 = val.(bool)
			}
			if val, ok := d.GetOk("security.0.enable_frontend_tls11"); ok {
				frontendProtocolTls11 = val.(bool)
			}
		}

		if v, exists := v["triple_des_ciphers_enabled"]; exists {
			tripleDesCiphers = v.(bool)
		}

		tlsEcdheEcdsaWithAes256CbcShaCiphers = v["tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled"].(bool)
		tlsEcdheEcdsaWithAes128CbcShaCiphers = v["tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled"].(bool)
		tlsEcdheRsaWithAes256CbcShaCiphers = v["tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled"].(bool)
		tlsEcdheRsaWithAes128CbcShaCiphers = v["tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled"].(bool)
		tlsRsaWithAes128GcmSha256Ciphers = v["tls_rsa_with_aes128_gcm_sha256_ciphers_enabled"].(bool)
		tlsRsaWithAes256GcmSha384Ciphers = v["tls_rsa_with_aes256_gcm_sha384_ciphers_enabled"].(bool)
		tlsRsaWithAes256CbcSha256Ciphers = v["tls_rsa_with_aes256_cbc_sha256_ciphers_enabled"].(bool)
		tlsRsaWithAes128CbcSha256Ciphers = v["tls_rsa_with_aes128_cbc_sha256_ciphers_enabled"].(bool)
		tlsRsaWithAes256CbcShaCiphers = v["tls_rsa_with_aes256_cbc_sha_ciphers_enabled"].(bool)
		tlsRsaWithAes128CbcShaCiphers = v["tls_rsa_with_aes128_cbc_sha_ciphers_enabled"].(bool)

		if skuIsConsumption && frontendProtocolSsl3 {
			if !features.FivePointOh() {
				return nil, errors.New("`frontend_ssl30_enabled`/`enable_frontend_ssl30` are not supported for Sku Tier `Consumption`")
			}
			return nil, errors.New("`frontend_ssl30_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tripleDesCiphers {
			return nil, errors.New("`enable_triple_des_ciphers` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheEcdsaWithAes256CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheEcdsaWithAes128CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheRsaWithAes256CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsEcdheRsaWithAes128CbcShaCiphers {
			return nil, errors.New("`tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes128GcmSha256Ciphers {
			return nil, errors.New("`tls_rsa_with_aes128_gcm_sha256_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes256CbcSha256Ciphers {
			return nil, errors.New("`tls_rsa_with_aes256_cbc_sha256_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes128CbcSha256Ciphers {
			return nil, errors.New("`tls_rsa_with_aes128_cbc_sha256_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes256CbcShaCiphers {
			return nil, errors.New("`tls_rsa_with_aes256_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}

		if skuIsConsumption && tlsRsaWithAes128CbcShaCiphers {
			return nil, errors.New("`tls_rsa_with_aes128_cbc_sha_ciphers_enabled` is not supported for Sku Tier `Consumption`")
		}
	}

	customProperties := map[string]string{
		apimBackendProtocolSsl3:   strconv.FormatBool(backendProtocolSsl3),
		apimBackendProtocolTls10:  strconv.FormatBool(backendProtocolTls10),
		apimBackendProtocolTls11:  strconv.FormatBool(backendProtocolTls11),
		apimFrontendProtocolTls10: strconv.FormatBool(frontendProtocolTls10),
		apimFrontendProtocolTls11: strconv.FormatBool(frontendProtocolTls11),
	}

	if !skuIsConsumption {
		customProperties[apimFrontendProtocolSsl3] = strconv.FormatBool(frontendProtocolSsl3)
		customProperties[apimTripleDesCiphers] = strconv.FormatBool(tripleDesCiphers)
		customProperties[apimTlsEcdheEcdsaWithAes256CbcShaCiphers] = strconv.FormatBool(tlsEcdheEcdsaWithAes256CbcShaCiphers)
		customProperties[apimTlsEcdheEcdsaWithAes128CbcShaCiphers] = strconv.FormatBool(tlsEcdheEcdsaWithAes128CbcShaCiphers)
		customProperties[apimTlsEcdheRsaWithAes256CbcShaCiphers] = strconv.FormatBool(tlsEcdheRsaWithAes256CbcShaCiphers)
		customProperties[apimTlsEcdheRsaWithAes128CbcShaCiphers] = strconv.FormatBool(tlsEcdheRsaWithAes128CbcShaCiphers)
		customProperties[apimTlsRsaWithAes128GcmSha256Ciphers] = strconv.FormatBool(tlsRsaWithAes128GcmSha256Ciphers)
		customProperties[apimTlsRsaWithAes256GcmSha384Ciphers] = strconv.FormatBool(tlsRsaWithAes256GcmSha384Ciphers)
		customProperties[apimTlsRsaWithAes256CbcSha256Ciphers] = strconv.FormatBool(tlsRsaWithAes256CbcSha256Ciphers)
		customProperties[apimTlsRsaWithAes128CbcSha256Ciphers] = strconv.FormatBool(tlsRsaWithAes128CbcSha256Ciphers)
		customProperties[apimTlsRsaWithAes256CbcShaCiphers] = strconv.FormatBool(tlsRsaWithAes256CbcShaCiphers)
		customProperties[apimTlsRsaWithAes128CbcShaCiphers] = strconv.FormatBool(tlsRsaWithAes128CbcShaCiphers)
	}

	if vp := d.Get("protocols").([]interface{}); len(vp) > 0 {
		vpr := vp[0].(map[string]interface{})
		enableHttp2 := vpr["http2_enabled"].(bool)
		if !features.FivePointOh() {
			if v, ok := d.GetOk("protocols.0.enable_http2"); ok {
				enableHttp2 = v.(bool)
			}
		}
		customProperties[apimHttp2Protocol] = strconv.FormatBool(enableHttp2)
	}

	return customProperties, nil
}

func expandAzureRmApiManagementVirtualNetworkConfigurations(d *pluginsdk.ResourceData) *apimanagementservice.VirtualNetworkConfiguration {
	vs := d.Get("virtual_network_configuration").([]interface{})
	if len(vs) == 0 {
		return nil
	}

	v := vs[0].(map[string]interface{})

	return &apimanagementservice.VirtualNetworkConfiguration{
		SubnetResourceId: pointer.To(v["subnet_id"].(string)),
	}
}

func flattenApiManagementSecurityCustomProperties(input map[string]string, skuIsConsumption bool) []interface{} {
	output := make(map[string]interface{})

	output["backend_ssl30_enabled"] = parseApiManagementNilableDictionary(input, apimBackendProtocolSsl3)
	output["backend_tls10_enabled"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls10)
	output["backend_tls11_enabled"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls11)
	output["frontend_tls10_enabled"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls10)
	output["frontend_tls11_enabled"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls11)

	if !features.FivePointOh() {
		output["enable_backend_ssl30"] = parseApiManagementNilableDictionary(input, apimBackendProtocolSsl3)
		output["enable_backend_tls10"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls10)
		output["enable_backend_tls11"] = parseApiManagementNilableDictionary(input, apimBackendProtocolTls11)
		output["enable_frontend_tls10"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls10)
		output["enable_frontend_tls11"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolTls11)
	}

	if !skuIsConsumption {
		output["frontend_ssl30_enabled"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolSsl3)

		if !features.FivePointOh() {
			output["enable_frontend_ssl30"] = parseApiManagementNilableDictionary(input, apimFrontendProtocolSsl3)
		}

		output["triple_des_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTripleDesCiphers)
		output["tls_ecdhe_ecdsa_with_aes256_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheEcdsaWithAes256CbcShaCiphers)
		output["tls_ecdhe_ecdsa_with_aes128_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheEcdsaWithAes128CbcShaCiphers)
		output["tls_ecdhe_rsa_with_aes256_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheRsaWithAes256CbcShaCiphers)
		output["tls_ecdhe_rsa_with_aes128_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsEcdheRsaWithAes128CbcShaCiphers)
		output["tls_rsa_with_aes256_gcm_sha384_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes256GcmSha384Ciphers)
		output["tls_rsa_with_aes128_gcm_sha256_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes128GcmSha256Ciphers)
		output["tls_rsa_with_aes256_cbc_sha256_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes256CbcSha256Ciphers)
		output["tls_rsa_with_aes128_cbc_sha256_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes128CbcSha256Ciphers)
		output["tls_rsa_with_aes256_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes256CbcShaCiphers)
		output["tls_rsa_with_aes128_cbc_sha_ciphers_enabled"] = parseApiManagementNilableDictionary(input, apimTlsRsaWithAes128CbcShaCiphers)
	}

	return []interface{}{output}
}

func flattenApiManagementProtocolsCustomProperties(input map[string]string) []interface{} {
	output := make(map[string]interface{})

	output["http2_enabled"] = parseApiManagementNilableDictionary(input, apimHttp2Protocol)

	if !features.FivePointOh() {
		output["enable_http2"] = parseApiManagementNilableDictionary(input, apimHttp2Protocol)
	}

	return []interface{}{output}
}

func flattenApiManagementVirtualNetworkConfiguration(input *apimanagementservice.VirtualNetworkConfiguration) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	virtualNetworkConfiguration := make(map[string]interface{})

	if input.SubnetResourceId != nil {
		subnetId, err := commonids.ParseSubnetIDInsensitively(*input.SubnetResourceId)
		if err != nil {
			return []interface{}{}, err
		}
		virtualNetworkConfiguration["subnet_id"] = subnetId.ID()
	}

	return []interface{}{virtualNetworkConfiguration}, nil
}

func parseApiManagementNilableDictionary(input map[string]string, key string) bool {
	log.Printf("Parsing value for %q", key)

	v, ok := input[key]
	if !ok {
		log.Printf("%q was not found in the input - returning `false` as the default value", key)
		return false
	}

	val, err := strconv.ParseBool(v)
	if err != nil {
		log.Printf(" parsing %q (key %q) as bool: %+v - assuming false", key, v, err)
		return false
	}

	return val
}

func expandApiManagementSignInSettings(input []interface{}) signinsettings.PortalSigninSettings {
	enabled := false

	if len(input) > 0 {
		vs := input[0].(map[string]interface{})
		enabled = vs["enabled"].(bool)
	}

	return signinsettings.PortalSigninSettings{
		Properties: &signinsettings.PortalSigninSettingProperties{
			Enabled: pointer.To(enabled),
		},
	}
}

func flattenApiManagementSignInSettings(input signinsettings.PortalSigninSettings) []interface{} {
	enabled := false

	if props := input.Properties; props != nil {
		if props.Enabled != nil {
			enabled = pointer.From(props.Enabled)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled": enabled,
		},
	}
}

func expandApiManagementDelegationSettings(input []interface{}) delegationsettings.PortalDelegationSettings {
	if len(input) == 0 {
		return delegationsettings.PortalDelegationSettings{}
	}

	vs := input[0].(map[string]interface{})

	props := delegationsettings.PortalDelegationSettingsProperties{
		UserRegistration: &delegationsettings.RegistrationDelegationSettingsProperties{
			Enabled: pointer.To(vs["user_registration_enabled"].(bool)),
		},
		Subscriptions: &delegationsettings.SubscriptionsDelegationSettingsProperties{
			Enabled: pointer.To(vs["subscriptions_enabled"].(bool)),
		},
	}

	validationKey := vs["validation_key"].(string)
	if !vs["user_registration_enabled"].(bool) && !vs["subscriptions_enabled"].(bool) && validationKey == "" {
		// for some reason we cannot leave this empty
		props.ValidationKey = pointer.To("cGxhY2Vob2xkZXIxCg==")
	}
	if validationKey != "" {
		props.ValidationKey = pointer.To(validationKey)
	}

	url := vs["url"].(string)
	if !vs["user_registration_enabled"].(bool) && !vs["subscriptions_enabled"].(bool) && url == "" {
		// for some reason we cannot leave this empty
		props.Url = pointer.To("https://www.placeholder.com")
	}
	if url != "" {
		props.Url = pointer.To(url)
	}

	return delegationsettings.PortalDelegationSettings{
		Properties: &props,
	}
}

func flattenApiManagementDelegationSettings(input delegationsettings.PortalDelegationSettings, keyContract delegationsettings.PortalSettingValidationKeyContract) []interface{} {
	url := ""
	subscriptionsEnabled := false
	userRegistrationEnabled := false

	if props := input.Properties; props != nil {
		url = pointer.From(props.Url)
		if props.Subscriptions != nil {
			subscriptionsEnabled = pointer.From(props.Subscriptions.Enabled)
		}
		if props.UserRegistration != nil {
			userRegistrationEnabled = pointer.From(props.UserRegistration.Enabled)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"url":                       url,
			"subscriptions_enabled":     subscriptionsEnabled,
			"user_registration_enabled": userRegistrationEnabled,
			"validation_key":            pointer.From(keyContract.ValidationKey),
		},
	}
}

func expandApiManagementSignUpSettings(input []interface{}) signupsettings.PortalSignupSettings {
	if len(input) == 0 {
		return signupsettings.PortalSignupSettings{
			Properties: &signupsettings.PortalSignupSettingsProperties{
				Enabled: pointer.To(false),
				TermsOfService: &signupsettings.TermsOfServiceProperties{
					ConsentRequired: pointer.To(false),
					Enabled:         pointer.To(false),
					Text:            pointer.To(""),
				},
			},
		}
	}

	vs := input[0].(map[string]interface{})

	props := signupsettings.PortalSignupSettingsProperties{
		Enabled: pointer.To(vs["enabled"].(bool)),
	}

	termsOfServiceRaw := vs["terms_of_service"].([]interface{})
	if len(termsOfServiceRaw) > 0 {
		termsOfServiceVs := termsOfServiceRaw[0].(map[string]interface{})
		props.TermsOfService = &signupsettings.TermsOfServiceProperties{
			Enabled:         pointer.To(termsOfServiceVs["enabled"].(bool)),
			ConsentRequired: pointer.To(termsOfServiceVs["consent_required"].(bool)),
			Text:            pointer.To(termsOfServiceVs["text"].(string)),
		}
	}

	return signupsettings.PortalSignupSettings{
		Properties: &props,
	}
}

func flattenApiManagementSignUpSettings(input signupsettings.PortalSignupSettings) []interface{} {
	enabled := false
	termsOfService := make([]interface{}, 0)

	if props := input.Properties; props != nil {
		enabled = pointer.From(props.Enabled)

		if tos := props.TermsOfService; tos != nil {
			output := make(map[string]interface{})

			output["enabled"] = pointer.From(tos.Enabled)
			output["consent_required"] = pointer.From(tos.ConsentRequired)
			output["text"] = pointer.From(tos.Text)

			termsOfService = append(termsOfService, output)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":          enabled,
			"terms_of_service": termsOfService,
		},
	}
}

func expandApiManagementTenantAccessSettings(input []interface{}) tenantaccess.AccessInformationUpdateParameters {
	enabled := false

	if len(input) > 0 {
		vs := input[0].(map[string]interface{})
		enabled = vs["enabled"].(bool)
	}

	return tenantaccess.AccessInformationUpdateParameters{
		Properties: &tenantaccess.AccessInformationUpdateParameterProperties{
			Enabled: pointer.To(enabled),
		},
	}
}

func flattenApiManagementTenantAccessSettings(input tenantaccess.AccessInformationSecretsContract) []interface{} {
	result := make(map[string]interface{})

	result["enabled"] = pointer.From(input.Enabled)
	result["tenant_id"] = pointer.From(input.Id)
	result["primary_key"] = pointer.From(input.PrimaryKey)
	result["secondary_key"] = pointer.From(input.SecondaryKey)

	return []interface{}{result}
}

func flattenAPIManagementCertificates(d *pluginsdk.ResourceData, inputs *[]apimanagementservice.CertificateConfiguration) []interface{} {
	if inputs == nil || len(*inputs) == 0 {
		return []interface{}{}
	}

	outputs := []interface{}{}
	for i, input := range *inputs {
		var pwd, encodedCertificate string
		if v, ok := d.GetOk(fmt.Sprintf("certificate.%d.certificate_password", i)); ok {
			pwd = v.(string)
		}

		if v, ok := d.GetOk(fmt.Sprintf("certificate.%d.encoded_certificate", i)); ok {
			encodedCertificate = v.(string)
		}

		output := map[string]interface{}{
			"certificate_password": pwd,
			"encoded_certificate":  encodedCertificate,
			"store_name":           string(input.StoreName),
			"expiry":               input.Certificate.Expiry,
			"subject":              input.Certificate.Subject,
			"thumbprint":           input.Certificate.Thumbprint,
		}
		outputs = append(outputs, output)
	}
	return outputs
}

func optedOutOfRecoveringSoftDeletedApiManagementErrorFmt(name, location string) string {
	message := `
An existing soft-deleted API Management exists with the Name %q in the location %q, however
automatically recovering this API Management has been disabled via the "features" block.

Terraform can automatically recover the soft-deleted API Management when this behaviour is
enabled within the "features" block (located within the "provider" block) - more
information can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import", or pick a different name/location.
`
	return fmt.Sprintf(message, name, location)
}
