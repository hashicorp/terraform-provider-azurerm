// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_function_app_slot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state WindowsFunctionAppSlotModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.GetSlot(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("reading Windows %s: %v", id, err)
			}

			model := *existing.Model

			_, planSKU, err := helpers.ServicePlanInfoForAppSlot(ctx, metadata, *id)
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("service_plan_id") {
				o, n := metadata.ResourceData.GetChange("service_plan_id")
				oldPlan, err := commonids.ParseAppServicePlanID(o.(string))
				if err != nil {
					return err
				}

				newPlan, err := commonids.ParseAppServicePlanID(n.(string))
				if err != nil {
					return err
				}
				locks.ByID(oldPlan.ID())
				defer locks.UnlockByID(oldPlan.ID())
				locks.ByID(newPlan.ID())
				defer locks.UnlockByID(newPlan.ID())
				if model.Properties == nil {
					return fmt.Errorf("updating Service Plan for Windows %s: Slot SiteProperties was nil", *id)
				}
				model.Properties.ServerFarmId = pointer.To(newPlan.ID())
			}

			// Only send for Dynamic and ElasticPremium
			sendContentSettings := (helpers.PlanIsConsumption(planSKU) || helpers.PlanIsElastic(planSKU)) && !state.ForceDisableContentShare

			// Some service plan updates are allowed - see customiseDiff for exceptions
			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = pointer.To(state.Enabled)
			}

			if metadata.ResourceData.HasChange("https_only") {
				model.Properties.HTTPSOnly = pointer.To(state.HttpsOnly)
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

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err = client.DeleteSwiftVirtualNetworkSlot(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					model.Properties.VirtualNetworkSubnetId = empty
				} else {
					model.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)
				}
			}

			if metadata.ResourceData.HasChange("vnet_image_pull_enabled") {
				model.Properties.VnetImagePullEnabled = pointer.To(state.VnetImagePullEnabled)
			}

			if metadata.ResourceData.HasChange("end_to_end_tls_encryption_enabled") {
				model.Properties.EndToEndEncryptionEnabled = pointer.To(state.EndToEndTLSEncryptionEnabled)
			}

			if metadata.ResourceData.HasChange("storage_account") {
				storageAccountUpdate := helpers.ExpandStorageConfig(state.StorageAccounts)
				if _, err = client.UpdateAzureStorageAccountsSlot(ctx, *id, *storageAccountUpdate); err != nil {
					return fmt.Errorf("updating Storage Accounts for Windows %s: %+v", id, err)
				}
			}

			storageString := state.StorageAccountName
			if !state.StorageUsesMSI {
				if state.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, state.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, state.StorageAccountName, state.StorageAccountKey, *storageDomainSuffix)
				}
			}

			if sendContentSettings {
				appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
				}
				if state.AppSettings == nil {
					state.AppSettings = make(map[string]string)
				}
				state.AppSettings = helpers.ParseContentSettings(appSettingsResp.Model, state.AppSettings)
			}

			// Note: We process this regardless to give us a "clean" view of service-side app_settings, so we can reconcile the user-defined entries later
			siteConfig, err := helpers.ExpandSiteConfigWindowsFunctionAppSlot(state.SiteConfig, model.Properties.SiteConfig, metadata, state.FunctionExtensionsVersion, storageString, state.StorageUsesMSI)
			if state.BuiltinLogging {
				if state.AppSettings == nil && !state.StorageUsesMSI {
					state.AppSettings = make(map[string]string)
				}
				if !state.StorageUsesMSI {
					state.AppSettings["AzureWebJobsDashboard"] = storageString
				} else {
					state.AppSettings["AzureWebJobsDashboard__accountName"] = state.StorageAccountName
				}
			}

			if metadata.ResourceData.HasChange("site_config") {
				if err != nil {
					return fmt.Errorf("expanding Site Config for Windows %s: %+v", id, err)
				}
				model.Properties.SiteConfig = siteConfig
				model.Properties.VnetRouteAllEnabled = model.Properties.SiteConfig.VnetRouteAllEnabled
			}

			model.Properties.SiteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, state.AppSettings)

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				pna := helpers.PublicNetworkAccessEnabled
				if !state.PublicNetworkAccess {
					pna = helpers.PublicNetworkAccessDisabled
				}

				// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
				model.Properties.PublicNetworkAccess = pointer.To(pna)
				model.Properties.SiteConfig.PublicNetworkAccess = model.Properties.PublicNetworkAccess
			}

			if err = client.CreateOrUpdateSlotThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating Windows %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("ftp_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingFTPBasicAuthEnabled,
					},
				}
				if _, err = client.UpdateFtpAllowedSlot(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("webdeploy_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingDeployBasicAuthEnabled,
					},
				}
				if _, err = client.UpdateScmAllowedSlot(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if _, err = client.UpdateConfigurationSlot(ctx, *id, webapps.SiteConfigResource{Properties: siteConfig}); err != nil {
				return fmt.Errorf("updating Site Config for Windows %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = &map[string]webapps.ConnStringValueTypePair{}
				}
				if _, err = client.UpdateConnectionStringsSlot(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				// (@jackofallops) - in the case of a removal of this block, we need to zero these settings
				if authUpdate.Properties == nil {
					authUpdate.Properties = helpers.DefaultAuthSettingsProperties()
				}
				if _, err = client.UpdateAuthSettingsSlot(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				// (@toddgiguere) - in the case of a removal of this block, we need to zero these settings
				if authV2Update.Properties == nil {
					authV2Update.Properties = helpers.DefaultAuthV2SettingsProperties()
				}
				if _, err := client.UpdateAuthSettingsV2Slot(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for Windows %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("backup") {
				backupUpdate, err := helpers.ExpandBackupConfig(state.Backup)
				if err != nil {
					return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
				}

				if backupUpdate.Properties == nil {
					if _, err = client.DeleteBackupConfigurationSlot(ctx, *id); err != nil {
						return fmt.Errorf("removing Backup Settings for Windows %s: %+v", id, err)
					}
				} else {
					if _, err = client.UpdateBackupConfigurationSlot(ctx, *id, *backupUpdate); err != nil {
						return fmt.Errorf("updating Backup Settings for Windows %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("site_config.0.app_service_logs") {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(state.SiteConfig[0].AppServiceLogs)
				if _, err = client.UpdateDiagnosticLogsConfigSlot(ctx, *id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
