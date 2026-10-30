// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_web_app

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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"service_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateAppServicePlanID,
		},

		// Optional

		"app_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			ValidateFunc: validate.AppSettings,
		},

		"auth_settings": helpers.AuthSettingsSchema(),

		"auth_settings_v2": helpers.AuthV2SettingsSchema(),

		"backup": helpers.BackupSchema(),

		"client_affinity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_certificate_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"client_certificate_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "Required",
			ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForClientCertMode(), false),
		},

		"client_certificate_exclusion_paths": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Description: "Paths to exclude when using client certificates, separated by ;",
		},

		"connection_string": helpers.ConnectionStringSchema(),

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

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

		"end_to_end_tls_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"vnet_image_pull_enabled": {
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
		},

		"logs": helpers.LogsConfigSchema(),

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

		"site_config": helpers.SiteConfigSchemaLinux(),

		"sticky_settings": helpers.StickySettingsSchema(),

		"storage_account": helpers.StorageAccountSchema(),

		"zip_deploy_file": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The local path and filename of the Zip packaged application to deploy to this Linux Web App. **Note:** Using this value requires either `WEBSITE_RUN_FROM_PACKAGE=1` or `SCM_DO_BUILD_DURING_DEPLOYMENT=true` to be set on the App in `app_settings`.",
		},

		"tags": commonschema.Tags(),
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
				if planId.(string) == "" {
					return nil
				}
				_, newValue := rd.GetChange("vnet_image_pull_enabled")
				if newValue.(bool) {
					return nil
				}
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
				}
			}
			return nil
		},
	}
}
