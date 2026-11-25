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
			Description:  "The name of the service instance. Used for service endpoint, must be unique within the audience.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"kind": {
			Description:  "The type of the service. Values at time of publication are: `fhir`, `fhir-Stu3` and `fhir-R4`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(resource.KindFhir),
			ValidateFunc: validation.StringInSlice(resource.PossibleValuesForKind(), false),
		},

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"cosmosdb_throughput": {
			Description:  "The provisioned throughput for the backing database. Range of `400`-`100000`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1000,
			ValidateFunc: validation.IntBetween(1, 100000),
		},

		"cosmosdb_key_vault_key_versionless_id": {
			Description:  "A versionless Key Vault Key ID for CMK encryption of the backing database.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeKey),
		},

		"access_policy_object_ids": {
			Description: "A set of Azure object IDs that are allowed to access the Service.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
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
						Description:  "The Azure Active Directory (tenant) that serves as the authentication authority to access the service. Authority must be registered to Azure AD and in the following format: <https://{Azure-AD-endpoint}/{tenant-id}.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
					},
					"audience": {
						Description:  "The intended audience to receive authentication tokens for the service.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"authentication_configuration.0.authority", "authentication_configuration.0.audience", "authentication_configuration.0.smart_proxy_enabled"},
					},
					"smart_proxy_enabled": {
						Description:  "Enables the 'SMART on FHIR' option for mobile and web implementations.",
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
						Description: "A set of origins to be allowed via CORS.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						MaxItems:    64,
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
						Description: "A set of headers to be allowed via CORS.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						MaxItems:    64,
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
						Description: "The methods to be allowed via CORS.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    64,
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
						Description:  "The max age to be allowed via CORS.",
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
						Description: "If credentials are allowed via CORS.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
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
			Description:  "The name of the Storage Account which the operation configuration information is exported to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"public_network_access_enabled": {
			Description: "Whether public network access is enabled or disabled for this service instance.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tags": commonschema.Tags(),
	}
}
