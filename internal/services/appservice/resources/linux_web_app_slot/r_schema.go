// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_web_app_slot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Linux Web App Slot.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"app_service_id": {
			Description:  "The ID of the Linux Web App this Deployment Slot will be part of.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateWebAppID,
		},

		// Optional

		"service_plan_id": {
			Description:  "The ID of the Service Plan in which to run this slot. If not specified the same Service Plan as the Linux Web App will be used.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		"app_settings": {
			Description:  "A map of key-value pairs of App Settings.",
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ValidateFunc: validate.AppSettings,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"auth_settings_v2": helpers.AuthV2SettingsSchema(),

		"backup": helpers.BackupSchema(),

		"client_affinity_enabled": {
			Description: "Should Client Affinity be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"client_certificate_enabled": {
			Description: "Should Client Certificates be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"client_certificate_mode": {
			Description:  "The Client Certificate mode. This property has no effect when `client_cert_enabled` is `false`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Required",
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"client_certificate_exclusion_paths": {
			Description: "Paths to exclude when using client certificates, separated by `;`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"enabled": {
			Description: "Should the Linux Web App be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"https_only": {
			Description: "Should the Linux Web App require HTTPS connections?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_backup_restore_enabled": {
			Description: "Whether backup and restore operations over the linked Virtual Network are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_subnet_id": {
			Description:  "The subnet ID which will be used by this Web App Slot for [regional Virtual Network integration](https://learn.microsoft.com/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"vnet_image_pull_enabled": {
			Description: "Should the traffic for the image pull be routed over Virtual Network enabled?",
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
			Description: "Should public network access be enabled for the Web App?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"webdeploy_publish_basic_authentication_enabled": {
			Description: "Should the default WebDeploy Basic Authentication publishing credentials enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"ftp_publish_basic_authentication_enabled": {
			Description: "Should the default FTP Basic Authentication publishing profile be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"logs": helpers.LogsConfigSchema(),

		"site_config": SiteConfigSchemaLinuxWebAppSlot(),

		"storage_account": helpers.StorageAccountSchema(),

		"zip_deploy_file": {
			Description:  "The local path and filename of the Zip packaged application to deploy to this Linux Web App.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"end_to_end_tls_encryption_enabled": {
			Description: "Whether to enable end-to-end TLS encryption.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"tags": commonschema.Tags(),
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
			Description: "The ID of the App Service Environment used by App Service Slot.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The Kind value for this Linux Web App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"app_metadata": {
			Description: "A `app_metadata`.",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"outbound_ip_addresses": {
			Description: "A comma-separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_address_list": {
			Description: "A list of outbound IP addresses - such as `[\"52.23.25.3\", \"52.143.43.12\"]`.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"possible_outbound_ip_addresses": {
			Description: "A comma-separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"possible_outbound_ip_address_list": {
			Description: "A `possible_outbound_ip_address_list`.",
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
				appId := rd.Get("app_service_id")
				if appId.(string) == "" {
					return nil
				}
				_, newValue := rd.GetChange("vnet_image_pull_enabled")
				if newValue.(bool) {
					return nil
				}
				functionAppId, err := commonids.ParseAppServiceID(appId.(string))
				if err != nil {
					return err
				}

				webApp, err := appClient.Get(ctx, *functionAppId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", functionAppId, err)
				}
				if webAppModel := webApp.Model; webAppModel != nil && webAppModel.Properties != nil {
					if ase := webAppModel.Properties.HostingEnvironmentProfile; ase != nil && ase.Id != nil && *ase.Id != "" && !newValue.(bool) {
						return fmt.Errorf("`vnet_image_pull_enabled` cannot be disabled for app slot running in an app service environment")
					}
				}
			}
			return nil
		},
	}
}
