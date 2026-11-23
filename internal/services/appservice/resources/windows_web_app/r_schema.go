// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package windows_web_app

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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
			Description:  "The name which should be used for this Windows Web App.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"service_plan_id": {
			Description:  "The ID of the Service Plan that this Windows App Service will be created in.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		// Optional

		"app_settings": {
			Description: "A map of key-value pairs of App Settings.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			ValidateFunc: validate.AppSettings,
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
			Description:  "The Client Certificate mode. This property has no effect when `client_certificate_enabled` is `false`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(webapps.ClientCertModeRequired),
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Paths to exclude when using client certificates, separated by `;`.",
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"enabled": {
			Description: "Should the Windows Web App be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"https_only": {
			Description: "Should the Windows Web App require HTTPS connections?",
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

		"logs": helpers.LogsConfigSchema(),

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

		"site_config": SiteConfigSchemaWindows(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"storage_account": helpers.StorageAccountSchemaWindows(),

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Windows Web App.",
		},

		"tags": commonschema.Tags(),

		"virtual_network_backup_restore_enabled": {
			Description: "Whether backup and restore operations over the linked Virtual Network are enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_image_pull_enabled": {
			Description: "Whether traffic for the image pull should be routed over the Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_network_subnet_id": {
			Description:  "The subnet ID which will be used by this Web App for [regional Virtual Network integration](https://learn.microsoft.com/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"end_to_end_tls_encryption_enabled": {
			Description: "Whether to enable end-to-end TLS encryption.",
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hosting_environment_id": {
			Description: "The ID of the App Service Environment used by App Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kind": {
			Description: "The Kind value for this Windows Web App.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"outbound_ip_addresses": {
			Description: "A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`.",
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
			Description: "A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"possible_outbound_ip_address_list": {
			Description: "A list of possible outbound IP address.",
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

			model := WindowsWebAppModel{}
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			if metadata.ResourceDiff.HasChange("virtual_network_image_pull_enabled") {
				planId := model.ServicePlanId
				if planId == "" {
					return nil
				}

				_, newValue := metadata.ResourceDiff.GetChange("virtual_network_image_pull_enabled")
				if newValue.(bool) {
					return nil
				}

				servicePlanID, err := commonids.ParseAppServicePlanID(planId)
				if err != nil {
					return err
				}

				resp, err := client.Get(ctx, *servicePlanID)
				if err != nil {
					return fmt.Errorf("retrieving %s: %w", *servicePlanID, err)
				}

				if aspModel := resp.Model; aspModel != nil {
					if aspModel.Properties != nil && aspModel.Properties.HostingEnvironmentProfile != nil &&
						pointer.From(aspModel.Properties.HostingEnvironmentProfile.Id) != "" && !newValue.(bool) {
						return fmt.Errorf("`virtual_network_image_pull_enabled` cannot be disabled for app running in an app service environment")
					}
				}
			}

			return nil
		},
	}
}
