// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_function_app_slot

import (
	"context"
	"fmt"
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
			Description:  "Specifies the name of the Windows Function App Slot.",
		},

		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
			Description:  "The ID of the Windows Function App this Slot is a member of.",
		},

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			Description:  "The backend storage account name which will be used by this Function App Slot.",
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
			Description: "The access key which will be used to access the storage account for the Function App Slot.",
		},

		"storage_uses_managed_identity": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ConflictsWith: []string{
				"storage_account_access_key",
				"storage_key_vault_secret_id",
			},
			Description: "Should the Function App Slot use its Managed Identity to access storage?",
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
			Description: "Should built in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting.",
		},

		"client_certificate_enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the Function App Slot use Client Certificates.",
		},

		"client_certificate_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      webapps.ClientCertModeOptional,
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
			Description:  "The mode of the Function App Slot's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser`.",
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
			Description: "Is the Windows Function App Slot enabled.",
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
			Description: "The runtime version associated with the Function App Slot.",
		},

		"https_only": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Can the Function App Slot only be accessed via HTTPS?",
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

		"site_config": helpers.SiteConfigSchemaWindowsFunctionAppSlot(),

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
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The identifier used by App Service to perform domain ownership verification via DNS TXT record.",
		},

		"default_hostname": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The default hostname of the Windows Function App Slot.",
		},

		"hosting_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kind": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Kind value for this Windows Function App Slot.",
		},

		"outbound_ip_addresses": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.",
		},

		"outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of outbound IP addresses. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
		},

		"possible_outbound_ip_addresses": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`. For example `[\"52.23.25.3\", \"52.143.43.12\",\"52.143.43.17\"]`.",
		},

		"possible_outbound_ip_address_list": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
		},

		"site_credential": helpers.SiteCredentialSchema(),
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			appClient := metadata.Client.AppService.WebAppsClient
			rd := metadata.ResourceDiff
			if rd.HasChange("vnet_image_pull_enabled") {
				appId := rd.Get("function_app_id")
				if appId.(string) == "" {
					return nil
				}
				_, newValue := rd.GetChange("vnet_image_pull_enabled")
				functionAppId, err := commonids.ParseAppServiceID(appId.(string))
				if err != nil {
					return err
				}

				functionApp, err := appClient.Get(ctx, *functionAppId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", functionAppId, err)
				}
				if functionAppModel := functionApp.Model; functionAppModel != nil && functionAppModel.Properties != nil {
					if ase := functionAppModel.Properties.HostingEnvironmentProfile; ase != nil && ase.Id != nil && *ase.Id != "" && !newValue.(bool) {
						return fmt.Errorf("`vnet_image_pull_enabled` cannot be disabled for app slot running in an app service environment")
					}
				}
			}
			return nil
		},
	}
}
