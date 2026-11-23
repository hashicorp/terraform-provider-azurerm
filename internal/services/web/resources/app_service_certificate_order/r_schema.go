// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate_order

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func appServiceCertificateOrderSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"auto_renew": {
			Type:     pluginsdk.TypeBool,
			Default:  true,
			Optional: true,
		},

		"certificates": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"key_vault_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"key_vault_secret_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"provisioning_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"csr": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{"distinguished_name"},
		},

		"distinguished_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{"csr"},
		},

		"key_size": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      2048,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"product_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Standard",
			ValidateFunc: validation.StringInSlice([]string{
				"Standard",
				"WildCard",
			}, false),
		},

		"validity_in_years": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntBetween(1, 3),
		},

		"domain_verification_token": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"expiration_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"is_private_key_external": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"app_service_certificate_not_renewable_reasons": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"signed_certificate_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"root_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"intermediate_thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
