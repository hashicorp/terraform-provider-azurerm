// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint_custom_domain

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnEndpointCustomDomainSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this CDN Endpoint Custom Domain.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CdnEndpointCustomDomainName(),
		},

		"cdn_endpoint_id": {
			Description:  "The ID of the CDN Endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(endpoints.ParseEndpointIDInsensitively),
		},

		"host_name": {
			Description:  "The host name of the custom domain.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"cdn_managed_https": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MinItems:      1,
			MaxItems:      1,
			ConflictsWith: []string{"user_managed_https"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_type": {
						Description:  "The type of HTTPS certificate.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInEnumSlice(cdn.PossibleCertificateTypeValues(), false),
					},
					"protocol_type": {
						Description:  "The type of protocol.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInEnumSlice(cdn.PossibleProtocolTypeValues(), false),
					},
					"tls_version": {
						Description: "The minimum TLS protocol version that is used for HTTPS.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(cdn.MinimumTLSVersionTLS12),
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.MinimumTLSVersionTLS12),
						}, false),
					},
				},
			},
		},

		"user_managed_https": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MinItems:      1,
			MaxItems:      1,
			ConflictsWith: []string{"cdn_managed_https"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tls_version": {
						Description: "The minimum TLS protocol version that is used for HTTPS.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(cdn.MinimumTLSVersionTLS12),
						ValidateFunc: validation.StringInSlice([]string{
							string(cdn.MinimumTLSVersionTLS12),
						}, false),
					},

					"key_vault_secret_id": {
						Description:  "The ID of the Key Vault Secret that contains the HTTPS certificate.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny),
					},
				},
			},
		},
	}
}
