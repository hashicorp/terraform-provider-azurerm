// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func healthcareServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"kind": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(resource.KindFhir),
			ValidateFunc: validation.StringInSlice(resource.PossibleValuesForKind(), false),
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"cosmosdb_throughput": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1000,
			ValidateFunc: validation.IntBetween(1, 100000),
		},

		"cosmosdb_key_vault_key_versionless_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeKey),
		},

		"access_policy_object_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsUUID,
			},
		},

		"authentication_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authority": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
					},
					"audience": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
					},
					"smart_proxy_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
					},
				},
			},
		},

		"cors_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{
							"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
							"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
							"cors_configuration.0.allow_credentials",
						},
					},
					"allowed_headers": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						AtLeastOneOf: []string{
							"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
							"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
							"cors_configuration.0.allow_credentials",
						},
					},
					"allowed_methods": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 64,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"DELETE",
								"GET",
								"HEAD",
								"MERGE",
								"POST",
								"OPTIONS",
								"PUT",
								"PATCH",
							}, false),
						},
						AtLeastOneOf: []string{
							"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
							"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
							"cors_configuration.0.allow_credentials",
						},
					},
					"max_age_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 2000000000),
						AtLeastOneOf: []string{
							"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
							"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
							"cors_configuration.0.allow_credentials",
						},
					},
					"allow_credentials": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						AtLeastOneOf: []string{
							"cors_configuration.0.allowed_origins", "cors_configuration.0.allowed_headers",
							"cors_configuration.0.allowed_methods", "cors_configuration.0.max_age_in_seconds",
							"cors_configuration.0.allow_credentials",
						},
					},
				},
			},
		},

		"configuration_export_storage_account_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}
