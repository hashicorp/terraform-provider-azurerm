// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_flex_consumption

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Function App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"service_plan_id": {
			Description:  "The ID of the App Service Plan within which to create this Function App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"backend_storage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The backend Storage Account name which will be used by this Function App.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: storageValidate.StorageAccountName,
						ExactlyOneOf: []string{
							"backend_storage.0.name",
							"backend_storage.0.key_vault_secret_id",
						},
					},
					"access_key": {
						Description:  "The access key which will be used to access the Storage Account for the Function App.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.NoZeroValues,
						ConflictsWith: []string{
							"backend_storage.0.key_vault_secret_id",
						},
					},
					"key_vault_secret_id": {
						Description:  "The Key Vault Secret ID, including version, that contains the Connection String to connect to the Storage Account for this Function App.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
						ExactlyOneOf: []string{
							"backend_storage.0.name",
							"backend_storage.0.key_vault_secret_id",
						},
					},
				},
			},
		},

		"deployment_storage": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_endpoint": {
						Description:  "The endpoint of the Storage Container where the Function App's code is hosted.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ConflictsWith: []string{
							"deployment_storage.0.user_assigned_identity_id",
						},
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						ConflictsWith: []string{
							"deployment_storage.0.access_key",
						},
					},
				},
			},
		},

		"runtime_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForRuntimeName(), false),
		},

		"runtime_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"instance_memory_in_mb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  2048,
		},

		"maximum_instance_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"http_concurrency": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		// the name is always being lower-cased by the api: https://github.com/Azure/azure-rest-api-specs/issues/33095
		"always_ready": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringIsNotEmpty,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"instance_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
				},
			},
		},

		"site_config": helpers.SiteConfigSchemaFunctionAppFlexConsumption(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"app_settings": {
			Description: "A map of key-value pairs for [App Settings](https://learn.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"auth_settings_v2": helpers.AuthV2SettingsSchema(),

		"client_certificate_enabled": {
			Description: "Whether the Function App should use Client Certificates.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"client_certificate_mode": {
			Description:  "The mode of the Function App's client certificates requirement for incoming requests.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"client_certificate_exclusion_paths": {
			Description: "Paths to exclude when using client certificates, separated by `;`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"enabled": {
			Description: "Whether the Function App is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"key_vault_reference_identity_id": {
			Description:  "The User Assigned Identity to use for Key Vault access.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"https_only": {
			Description: "Whether the Function App can only be accessed via HTTPS.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"zip_deploy_file": {
			Description:  "The local path and filename of the Zip packaged application to deploy to this Function App. **Note:** Using this value requires either `WEBSITE_RUN_FROM_PACKAGE=1` or `SCM_DO_BUILD_DURING_DEPLOYMENT=true` to be set on the App in `app_settings`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		schema["backend_storage"].Computed = true
		schema["backend_storage"].Optional = true
		schema["backend_storage"].Required = false
		schema["backend_storage"].RequiredWith = []string{"deployment_storage"}

		schema["deployment_storage"].Computed = true
		schema["deployment_storage"].Optional = true
		schema["deployment_storage"].Required = false
		schema["deployment_storage"].RequiredWith = []string{"backend_storage"}

		schema["storage_container_type"] = &pluginsdk.Schema{
			Description: "The type of the Storage Container where the Function App's code is hosted.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because these are deprecated storage settings
			Computed:     true,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForFunctionsDeploymentStorageType(), false),
			ExactlyOneOf: []string{"storage_container_type", "deployment_storage"},
			Deprecated:   "`storage_container_type` has been deprecated in favour of the `deployment_storage` property and will be removed in v6.0 of the AzureRM Provider.",
		}

		schema["storage_container_endpoint"] = &pluginsdk.Schema{
			Description: "The endpoint of the Storage Container where the Function App's code is hosted.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because these are deprecated storage settings
			Computed: true,
			RequiredWith: []string{
				"storage_container_type",
			},
			Deprecated: "`storage_container_endpoint` has been deprecated in favour of the `deployment_storage` property and will be removed in v6.0 of the AzureRM Provider.",
		}

		schema["storage_authentication_type"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because these are deprecated storage settings
			Computed:     true,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAuthenticationType(), false),
			RequiredWith: []string{
				"storage_container_type",
			},
			Deprecated: "`storage_authentication_type` has been deprecated in favour of the `deployment_storage` property and will be removed in v6.0 of the AzureRM Provider.",
		}

		schema["storage_access_key"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because these are deprecated storage settings
			Computed:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{
				"storage_container_type",
			},
			Deprecated: "`storage_access_key` has been deprecated in favour of the `deployment_storage` property and will be removed in v6.0 of the AzureRM Provider.",
		}

		schema["storage_user_assigned_identity_id"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because these are deprecated storage settings
			Computed:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			RequiredWith: []string{
				"storage_container_type",
			},
			Deprecated: "`storage_user_assigned_identity_id` has been deprecated in favour of the `deployment_storage` property and will be removed in v6.0 of the AzureRM Provider.",
		}
	}

	return schema
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"default_hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"possible_outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"site_credential": helpers.SiteCredentialSchema(),
	}
}
