// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_function_app

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
			Description:  "The name which should be used for this Linux Function App. Limit the function name to 32 characters to avoid naming collisions. For more information about [Function App naming rule](https://learn.microsoft.com/azure/azure-resource-manager/management/resource-name-rules#microsoftweb) and [Host ID Collisions](https://github.com/Azure/azure-functions-host/wiki/Host-IDs#host-id-collisions).",
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
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"storage_account_name": {
			Description:  "The backend Storage Account name which will be used by this Function App.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: storageValidate.StorageAccountName,
			ExactlyOneOf: []string{
				"storage_account_name",
				"storage_key_vault_secret_id",
			},
		},

		"storage_account_access_key": {
			Description:  "The access key which will be used to access the backend Storage Account for the Function App.",
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
			Description: "Whether the Function App should use its Managed Identity to access storage.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			ConflictsWith: []string{
				"storage_account_access_key",
				"storage_key_vault_secret_id",
			},
		},

		"storage_key_vault_secret_id": {
			Description:  "The Key Vault Secret ID, optionally including version, that contains the Connection String to connect to the Storage Account for this Function App.",
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

		"daily_memory_time_quota": {
			Description:  "The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects Function Apps under the consumption plan.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"enabled": {
			Description: "Whether the Linux Function App is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"content_share_force_disabled": {
			Description: "Whether the settings for linking the Function App to storage should be suppressed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"functions_extension_version": {
			Description: "The runtime version associated with the Function App.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "~4",
		},

		"https_only": {
			Description: "Whether the Function App can only be accessed via HTTPS.",
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

		"site_config": helpers.SiteConfigSchemaLinuxFunctionApp(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"storage_account": helpers.StorageAccountSchema(),

		"tags": commonschema.Tags(),

		"virtual_network_backup_restore_enabled": {
			Description: "Whether backup and restore operations over the linked Virtual Network are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_subnet_id": {
			Description:  "The subnet ID which will be used by this Function App for [regional Virtual Network integration](https://learn.microsoft.com/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).",
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

		"end_to_end_tls_encryption_enabled": {
			Description: "Whether to enable end-to-end TLS encryption.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"zip_deploy_file": {
			Description:  "The local path and filename of the Zip packaged application to deploy to this Linux Function App.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_id": {
			Description: "The ID of the App Service Environment used by Function App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The Kind value for this Linux Function App.",
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
			Description: "A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`.",
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
				if err != nil {
					return fmt.Errorf("could not read new Service Plan to check tier %s: %+v ", newPlanId, err)
				}
				if newPlan.Model != nil {
					if planSku := newPlan.Model.Sku; planSku != nil {
						if tier := planSku.Tier; tier != nil {
							newTierIsDynamic = strings.EqualFold(*tier, "dynamic")
							newTierIsBasic = strings.EqualFold(*tier, "basic")
						}
					}
				}
				if _, ok := rd.GetOk("backup"); ok && newTierIsDynamic {
					return fmt.Errorf("cannot specify backup configuration for Dynamic tier Service Plans, Standard or higher is required")
				}
				if _, ok := rd.GetOk("backup"); ok && newTierIsBasic {
					return fmt.Errorf("cannot specify backup configuration for Basic tier Service Plans, Standard or higher is required")
				}

				if strings.EqualFold(currentPlanIdRaw.(string), newPlanIdRaw.(string)) || currentPlanIdRaw.(string) == "" {
					// State migration escape for correcting case in serverFarms
					// change of case here will not move the app to a new Service Plan
					// also if the current Service Plan is empty, this is a new resource, so can skip this
					return nil
				}

				// Service Plans can only be updated in place when both New and Existing are not Dynamic
				if currentPlanIdRaw.(string) != "" {
					currentPlanId, err := commonids.ParseAppServicePlanIDInsensitively(currentPlanIdRaw.(string))
					if err != nil {
						return fmt.Errorf("reading existing plan id %+v", err)
					}

					currentPlan, err := client.Get(ctx, *currentPlanId)
					if err != nil || currentPlan.Model == nil {
						return fmt.Errorf("could not read old Service Plan to check tier %s: %+v", currentPlanId, err)
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
