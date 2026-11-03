// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_web_app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient

			id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state WindowsWebAppModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading Windows %s: %v", *id, err)
			}

			model := *existing.Model

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteConfig, err := client.GetConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
			}
			model.Properties.SiteConfig = webAppSiteConfig.Model.Properties

			var serviceFarmId string
			servicePlanChange := false
			if model.Properties.ServerFarmId != nil {
				serviceFarmId = *model.Properties.ServerFarmId
			}

			if metadata.ResourceData.HasChange("service_plan_id") {
				serviceFarmId = state.ServicePlanId
				model.Properties.ServerFarmId = pointer.To(serviceFarmId)
				servicePlanChange = true
			}
			servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(serviceFarmId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading App %s: %+v", servicePlanId, err)
			}

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = pointer.To(state.Enabled)
			}
			if metadata.ResourceData.HasChange("https_only") {
				model.Properties.HTTPSOnly = pointer.To(state.HttpsOnly)
			}
			if metadata.ResourceData.HasChange("client_affinity_enabled") {
				model.Properties.ClientAffinityEnabled = pointer.To(state.ClientAffinityEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_enabled") {
				model.Properties.ClientCertEnabled = pointer.To(state.ClientCertEnabled)
			}
			if metadata.ResourceData.HasChange("client_certificate_mode") {
				model.Properties.ClientCertMode = pointer.ToEnum[webapps.ClientCertMode](state.ClientCertMode)
			}
			if metadata.ResourceData.HasChange("client_certificate_exclusion_paths") {
				model.Properties.ClientCertExclusionPaths = pointer.To(state.ClientCertExclusionPaths)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("key_vault_reference_identity_id") {
				model.Properties.KeyVaultReferenceIdentity = pointer.To(state.KeyVaultReferenceIdentityID)
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("virtual_network_backup_restore_enabled") {
				model.Properties.VnetBackupRestoreEnabled = pointer.To(state.VirtualNetworkBackupRestoreEnabled)
			}

			if metadata.ResourceData.HasChange("virtual_network_image_pull_enabled") {
				model.Properties.VnetImagePullEnabled = pointer.To(state.VirtualNetworkImagePullEnabled)
			}

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetwork(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					model.Properties.VirtualNetworkSubnetId = empty
				} else {
					model.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			if metadata.ResourceData.HasChange("end_to_end_tls_encryption_enabled") {
				model.Properties.EndToEndEncryptionEnabled = pointer.To(state.EndToEndTLSEncryptionEnabled)
			}

			currentStack := ""
			sc := state.SiteConfig[0]

			if servicePlan.Model != nil && servicePlan.Model.Sku != nil && servicePlan.Model.Sku.Name != nil {
				if helpers.IsFreeOrSharedServicePlan(*servicePlan.Model.Sku.Name) {
					if sc.AlwaysOn {
						return fmt.Errorf("always_on feature has to be turned off before switching to a free/shared Sku")
					}
				}
			}

			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			if metadata.ResourceData.HasChanges("site_config", "app_settings") || servicePlanChange {
				model.Properties.SiteConfig, err = sc.ExpandForUpdate(metadata, model.Properties.SiteConfig, state.AppSettings)
				if err != nil {
					return err
				}
				model.Properties.VnetRouteAllEnabled = existing.Model.Properties.SiteConfig.VnetRouteAllEnabled
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				pna := helpers.PublicNetworkAccessEnabled
				if !state.PublicNetworkAccess {
					pna = helpers.PublicNetworkAccessDisabled
				}

				// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
				model.Properties.PublicNetworkAccess = pointer.To(pna)
				model.Properties.SiteConfig.PublicNetworkAccess = model.Properties.PublicNetworkAccess
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}

			siteMetadata := webapps.StringDictionary{Properties: &map[string]string{
				"CURRENT_STACK": currentStack,
			}}
			if _, err := client.UpdateMetadata(ctx, *id, siteMetadata); err != nil {
				return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
			}

			updateLogs := false

			// sending App Settings updates can clobber logs configuration so must be updated before we send any Log updates
			if metadata.ResourceData.HasChanges("app_settings", "site_config") {
				appSettingsUpdate := helpers.ExpandAppSettingsForUpdate(model.Properties.SiteConfig.AppSettings)
				appSettingsProps := *appSettingsUpdate.Properties
				if state.SiteConfig[0].HealthCheckEvictionTime != 0 {
					appSettingsProps["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = strconv.FormatInt(state.SiteConfig[0].HealthCheckEvictionTime, 10)
					appSettingsUpdate.Properties = &appSettingsProps
				} else {
					delete(appSettingsProps, "WEBSITE_HEALTHCHECK_MAXPINGFAILURES")
					appSettingsUpdate.Properties = &appSettingsProps
				}
				if _, err := client.UpdateApplicationSettings(ctx, *id, *appSettingsUpdate); err != nil {
					return fmt.Errorf("updating App Settings for Linux %s: %+v", *id, err)
				}

				updateLogs = true
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = &map[string]webapps.ConnStringValueTypePair{}
				}
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("sticky_settings") {
				emptySlice := make([]string, 0)
				stickySettings := helpers.ExpandStickySettings(state.StickySettings)
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: &webapps.SlotConfigNames{
						AppSettingNames:       &emptySlice,
						ConnectionStringNames: &emptySlice,
					},
				}

				if stickySettings != nil {
					if stickySettings.AppSettingNames != nil {
						stickySettingsUpdate.Properties.AppSettingNames = stickySettings.AppSettingNames
					}
					if stickySettings.ConnectionStringNames != nil {
						stickySettingsUpdate.Properties.ConnectionStringNames = stickySettings.ConnectionStringNames
					}
				}

				if _, err := client.UpdateSlotConfigurationNames(ctx, *id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				if authUpdate.Properties == nil {
					authUpdate.Properties = helpers.DefaultAuthSettingsProperties()
					updateLogs = true
				}
				if _, err := client.UpdateAuthSettings(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				// (@toddgiguere) - in the case of a removal of this block, we need to zero these settings
				if authV2Update.Properties == nil {
					authV2Update.Properties = helpers.DefaultAuthV2SettingsProperties()
				}
				if _, err := client.UpdateAuthSettingsV2(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Windows %s: %+v", id, err)
				}
				updateLogs = true
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
				}

				if backupUpdate.Properties == nil {
					if _, err = client.DeleteBackupConfiguration(ctx, *id); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", id, err)
					}
				} else {
					if _, err = client.UpdateBackupConfiguration(ctx, *id, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("logs") || updateLogs {
				logsUpdate := helpers.ExpandLogsConfig(state.LogsConfig)
				if logsUpdate.Properties == nil {
					logsUpdate = helpers.DisabledLogsConfig() // The API is update only, so we need to send an update with everything switched of when a user removes the "logs" block
				}
				if _, err = client.UpdateDiagnosticLogsConfig(ctx, *id, *logsUpdate); err != nil {
					return fmt.Errorf("updating Logs Config for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err := client.UpdateAzureStorageAccounts(ctx, *id, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("zip_deploy_file") {
				if err = helpers.GetCredentialsAndPublish(ctx, client, *id, state.ZipDeployFile); err != nil {
					return err
				}
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingFTPBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateFtpAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("webdeploy_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingDeployBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateScmAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
