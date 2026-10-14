// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_function_app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
			Description:  "Specifies the name of the Function App.",
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
			Description:  "The ID of the App Service Plan within which to create this Function App",
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			Description:  "The backend storage account name which will be used by this Function App.",
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
		},

		"storage_account_access_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
			ConflictsWith: []string{
				"storage_uses_managed_identity",
				"storage_key_vault_secret_id",
			},
			Description: "The access key which will be used to access the storage account for the Function App.",
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"storage_account_access_key",
				"storage_key_vault_secret_id",
			},
			Description: "Should the Function App use its Managed Identity to access storage?",
		},

		"storage_key_vault_secret_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
			Description: "The Key Vault Secret ID, including version, that contains the Connection String to connect to the storage account for this Function App.",
		},

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A map of key-value pairs for [App Settings](https://docs.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.",
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"auth_settings_v2": helpers.AuthV2SettingsSchema(),

		"backup": helpers.BackupSchema(),

		"builtin_logging_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Should built in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting",
		},

		"client_certificate_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the function app use Client Certificates",
		},

		"client_certificate_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
			Description:  "The mode of the Function App's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser` ",
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Paths to exclude when using client certificates, separated by ;",
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"daily_memory_time_quota": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
			Description:  "The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects function apps in Consumption Plans.",
		},

		"enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Is the Windows Function App enabled.",
		},

		"content_share_force_disabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force disable the content share settings.",
		},

		"functions_extension_version": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "~4",
			Description: "The runtime version associated with the Function App.",
		},

		"https_only": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Description: "Can the Function App only be accessed via HTTPS?",
		},

		"end_to_end_tls_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"key_vault_reference_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			Description:  "The User Assigned Identity to use for Key Vault access.",
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"site_config": helpers.SiteConfigSchemaWindowsFunctionApp(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"tags": commonschema.Tags(),

		"virtual_network_backup_restore_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"virtual_network_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"vnet_image_pull_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Is container image pull over virtual network enabled? Defaults to `false`.",
		},

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Windows Function App. **Note:** Using this value requires `WEBSITE_RUN_FROM_PACKAGE=1` to be set on the App in `app_settings`.",
		},
	}
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

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ServicePlanClient
			rd := metadata.ResourceDiff

			if rd.HasChange("vnet_image_pull_enabled") {
				planId := rd.Get("service_plan_id")
				// the plan id is known after apply during the initial creation
				if planId.(string) == "" {
					return nil
				}
				_, newValue := rd.GetChange("vnet_image_pull_enabled")
				servicePlanId, err := commonids.ParseAppServicePlanID(planId.(string))
				if err != nil {
					return err
				}

				asp, err := client.Get(ctx, *servicePlanId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", servicePlanId, err)
				}
				if aspModel := asp.Model; aspModel != nil {
					if aspModel.Properties != nil && aspModel.Properties.HostingEnvironmentProfile != nil &&
						aspModel.Properties.HostingEnvironmentProfile.Id != nil && *aspModel.Properties.HostingEnvironmentProfile.Id != "" && !newValue.(bool) {
						return fmt.Errorf("`vnet_image_pull_enabled` cannot be disabled for app running in an app service environment")
					}
					if sku := aspModel.Sku; sku != nil {
						if helpers.PlanIsConsumption(sku.Name) && newValue.(bool) {
							return fmt.Errorf("`vnet_image_pull_enabled` cannot be enabled on consumption plans")
						}
					}
				}
			}

			if rd.HasChange("service_plan_id") {
				currentPlanIdRaw, newPlanIdRaw := rd.GetChange("service_plan_id")
				if newPlanIdRaw.(string) == "" {
					// Plans creating a new service_plan inline will be empty as `Computed` known after apply
					return nil
				}
				newPlanId, err := commonids.ParseAppServicePlanID(newPlanIdRaw.(string))
				if err != nil {
					return fmt.Errorf("reading new plan id %+v", err)
				}

				var currentTierIsDynamic, newTierIsDynamic, newTierIsBasic bool

				newPlan, err := client.Get(ctx, *newPlanId)
				if err != nil || newPlan.Model == nil {
					return fmt.Errorf("could not read new Service Plan to check tier %s: %+v", newPlanId, err)
				}
				if planSku := newPlan.Model.Sku; planSku != nil {
					if tier := planSku.Tier; tier != nil {
						newTierIsDynamic = strings.EqualFold(*tier, "dynamic")
						newTierIsBasic = strings.EqualFold(*tier, "basic")
					}
				}

				if _, ok := rd.GetOk("backup"); ok && newTierIsDynamic {
					return fmt.Errorf("cannot specify backup configuration for Dynamic tier Service Plans, Standard or higher is required")
				}
				if _, ok := rd.GetOk("backup"); ok && newTierIsBasic {
					return fmt.Errorf("cannot specify backup configuration for Basic tier Service Plans, Standard or higher is required")
				}

				if strings.EqualFold(currentPlanIdRaw.(string), newPlanIdRaw.(string)) || currentPlanIdRaw == "" {
					// State migration escape for correcting case in serverFarms
					// change of case here will not move the app to a new Service Plan
					return nil
				}
				// Service Plans can only be updated in place when both New and Existing are not Dynamic
				if currentPlanIdRaw.(string) != "" {
					currentPlanId, err := commonids.ParseAppServicePlanID(currentPlanIdRaw.(string))
					if err != nil {
						return fmt.Errorf("reading existing plan id %+v", err)
					}

					currentPlan, err := client.Get(ctx, *currentPlanId)
					if err != nil || currentPlan.Model == nil {
						return fmt.Errorf("could not read current Service Plan to check tier %s: %+v", currentPlanId, err)
					}

					if planSku := currentPlan.Model.Sku; planSku != nil {
						if tier := planSku.Tier; tier != nil {
							currentTierIsDynamic = strings.EqualFold(*tier, "dynamic")
						}
					}

					if currentTierIsDynamic || newTierIsDynamic {
						if err := rd.ForceNew("service_plan_id"); err != nil {
							return err
						}
					}
				}
			}
			return nil
		},
	}
}
