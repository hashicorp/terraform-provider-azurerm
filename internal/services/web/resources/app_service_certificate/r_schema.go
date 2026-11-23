// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func appServiceCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"pfx_blob": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			RequiredWith: []string{"key_vault_secret_id"},
		},

		"key_vault_secret_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.DiffSuppressIgnoreKeyVaultKeyVersion,
			ValidateFunc:     keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny),
			ConflictsWith:    []string{"pfx_blob", "password"},
			ExactlyOneOf:     []string{"key_vault_secret_id", "pfx_blob"},
		},

		"app_service_plan_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subject_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"host_names": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"issuer": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"issue_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"expiration_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_profile_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
