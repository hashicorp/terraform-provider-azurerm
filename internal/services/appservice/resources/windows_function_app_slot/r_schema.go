// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Windows Function App Slot.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"function_app_id": {
			Description:  "The name of the Windows Function App this Slot is a member of.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
		},

		"service_plan_id": {
			Description:  "The ID of the Service Plan in which to run this slot. If not specified the same Service Plan as the Windows Function App will be used.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"storage_account_name": {
			Description:  "The backend Storage Account name which will be used by this Function App Slot.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
		},

		"storage_account_access_key": {
			Description:  "The access key which will be used to access the Storage Account for the Function App Slot.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
			ConflictsWith: []string{
				"storage_uses_managed_identity",
				"storage_key_vault_secret_id",
			},
		},

		"storage_uses_managed_identity": {
			Description: "Whether the Function App Slot should use its Managed Identity to access storage.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			ConflictsWith: []string{
				"storage_account_access_key",
				"storage_key_vault_secret_id",
			},
		},

		"storage_key_vault_secret_id": {
			Description:  "The Key Vault Secret ID, optionally including version, that contains the Connection String to connect to the Storage Account for this Function App Slot.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
		},

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

		"backup": helpers.BackupSchema(),

		"builtin_logging_enabled": {
			Description: "Whether built in logging should be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"client_certificate_enabled": {
			Description: "Whether the Function App Slot should use Client Certificates.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"client_certificate_mode": {
			Description:  "The mode of the Function App Slot's client certificates requirement for incoming requests.",
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

		"daily_memory_time_quota": {
			Description:  "The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects Function Apps in Consumption Plans.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"enabled": {
			Description: "Whether the Windows Function App Slot is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"content_share_force_disabled": {
			Description: "Force disable the content share settings.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"functions_extension_version": {
			Description: "The runtime version associated with the Function App Slot.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "~4",
		},

		"https_only": {
			Description: "Whether the Function App Slot can only be accessed via HTTPS.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"end_to_end_tls_encryption_enabled": {
			Description: "Whether to enable end-to-end TLS encryption.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"key_vault_reference_identity_id": {
			Description:  "The User Assigned Identity ID used for accessing KeyVault secrets. The identity must be assigned to the application in the `identity` block. [For more information see - Access vaults with a user-assigned identity](https://learn.microsoft.com/azure/app-service/app-service-key-vault-references#access-vaults-with-a-user-assigned-identity).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"public_network_access_enabled": {
			Description: "Whether public network access should be enabled for the Function App.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Description: "Whether the default WebDeploy Basic Authentication publishing credentials are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Description: "Whether the default FTP Basic Authentication publishing profile should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"site_config": SiteConfigSchemaWindowsFunctionAppSlot(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"tags": commonschema.Tags(),

		"virtual_network_backup_restore_enabled": {
			Description: "Whether backup and restore operations over the linked Virtual Network are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_subnet_id": {
			Description:  "The subnet ID which will be used by this Function App Slot for [regional Virtual Network integration](https://learn.microsoft.com/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"vnet_image_pull_enabled": {
			Description: "Whether container image pull over Virtual Network is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"custom_domain_verification_id": {
			Description: "The identifier used by App Service to perform domain ownership verification via DNS TXT record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"default_hostname": {
			Description: "The default hostname of the Windows Function App Slot.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"hosting_environment_id": {
			Description: "The ID of the App Service Environment used by Function App Slot.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The Kind value for this Windows Function App Slot.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_addresses": {
			Description: "A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_address_list": {
			Description: "A list of outbound IP addresses. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Description: "A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`. For example `[\"52.23.25.3\", \"52.143.43.12\",\"52.143.43.17\"]`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"possible_outbound_ip_address_list": {
			Description: "A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `[\"52.23.25.3\", \"52.143.43.12\"]`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
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
