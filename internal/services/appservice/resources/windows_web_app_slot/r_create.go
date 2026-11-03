// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_web_app_slot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var webAppSlot WindowsWebAppSlotModel
			if err := metadata.Decode(&webAppSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			appId, err := commonids.ParseWebAppID(webAppSlot.AppServiceId)
			if err != nil {
				return err
			}

			id := webapps.NewSlotID(appId.SubscriptionId, appId.ResourceGroupName, appId.SiteName, webAppSlot.Name)

			webApp, err := client.Get(ctx, *appId)
			if err != nil {
				return fmt.Errorf("reading parent Windows Web App for %s: %+v", id, err)
			}

			var servicePlanId *commonids.AppServicePlanId
			differentServicePlanToParent := false
			if webApp.Model == nil || webApp.Model.Properties == nil || webApp.Model.Properties.ServerFarmId == nil {
				return fmt.Errorf("could not determine Service Plan ID for %s: %+v", id, err)
			}

			servicePlanId, err = commonids.ParseAppServicePlanIDInsensitively(*webApp.Model.Properties.ServerFarmId)
			if err != nil {
				return err
			}

			if webAppSlot.ServicePlanID != "" {
				newServicePlanId, err := commonids.ParseAppServicePlanID(webAppSlot.ServicePlanID)
				if err != nil {
					return fmt.Errorf("parsing service_plan_id for %s: %+v", id, err)
				}
				// we only set `service_plan_id` when it differs from the parent `service_plan_id` which is causing issues
				// https://github.com/hashicorp/terraform-provider-azurerm/issues/21024
				// we'll error here if the `service_plan_id` equals the parent `service_plan_id`
				if strings.EqualFold(newServicePlanId.ID(), servicePlanId.ID()) {
					return fmt.Errorf("`service_plan_id` should only be specified when it differs from the `service_plan_id` of the associated Web App")
				}
				differentServicePlanToParent = true
				servicePlanId = newServicePlanId
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetSlot(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			sc := webAppSlot.SiteConfig[0]
			siteConfig, err := sc.ExpandForCreate(webAppSlot.AppSettings)
			if err != nil {
				return err
			}

			currentStack := ""
			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(webAppSlot.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(webApp.Model.Location),
				Tags:     pointer.To(webAppSlot.Tags),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					Enabled:                   pointer.To(webAppSlot.Enabled),
					HTTPSOnly:                 pointer.To(webAppSlot.HttpsOnly),
					SiteConfig:                siteConfig,
					ClientAffinityEnabled:     pointer.To(webAppSlot.ClientAffinityEnabled),
					ClientCertEnabled:         pointer.To(webAppSlot.ClientCertEnabled),
					ClientCertMode:            pointer.ToEnum[webapps.ClientCertMode](webAppSlot.ClientCertMode),
					ClientCertExclusionPaths:  pointer.To(webAppSlot.ClientCertExclusionPaths),
					VnetBackupRestoreEnabled:  pointer.To(webAppSlot.VirtualNetworkBackupRestoreEnabled),
					VnetRouteAllEnabled:       siteConfig.VnetRouteAllEnabled,
					EndToEndEncryptionEnabled: pointer.To(webAppSlot.EndToEndTLSEncryptionEnabled),
				},
			}

			siteEnvelope.Properties.VnetImagePullEnabled = pointer.To(webAppSlot.VirtualNetworkImagePullEnabled)

			if differentServicePlanToParent {
				siteEnvelope.Properties.ServerFarmId = pointer.To(servicePlanId.ID())
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !webAppSlot.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if webAppSlot.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(webAppSlot.KeyVaultReferenceIdentityID)
			}

			if webAppSlot.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(webAppSlot.VirtualNetworkSubnetID)
				siteEnvelope.Properties.ServerFarmId = pointer.To(servicePlanId.ID())
			}

			if err := client.CreateOrUpdateSlotCallbackThenPoll(ctx, id, siteEnvelope, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			metadata.SetID(id)

			// (@jackofallops) - Windows Web App Slots need the siteConfig sending individually to actually accept the `windowsFxVersion` value or it's set as `DOCKER|` only.
			siteConfigUpdate := webapps.SiteConfigResource{
				Properties: siteConfig,
			}
			if _, err = client.UpdateConfigurationSlot(ctx, id, siteConfigUpdate); err != nil {
				return fmt.Errorf("updating %s site config: %+v", id, err)
			}

			if currentStack != "" {
				siteMetadata := webapps.StringDictionary{Properties: &map[string]string{
					"CURRENT_STACK": currentStack,
				}}
				if _, err := client.UpdateMetadataSlot(ctx, id, siteMetadata); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
				}
			}

			appSettings := helpers.ExpandAppSettingsForUpdate(siteConfig.AppSettings)
			appSettingsProps := *appSettings.Properties
			if metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettingsProps["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = strconv.FormatInt(webAppSlot.SiteConfig[0].HealthCheckEvictionTime, 10)
				appSettings.Properties = &appSettingsProps
			}
			if appSettings != nil {
				if _, err := client.UpdateApplicationSettingsSlot(ctx, id, *appSettings); err != nil {
					return fmt.Errorf("setting App Settings for Windows %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webAppSlot.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettingsSlot(ctx, id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(webAppSlot.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2Slot(ctx, id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Linux %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsConfig := helpers.ExpandLogsConfig(webAppSlot.LogsConfig)
				if logsConfig.Properties != nil {
					if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id, *logsConfig); err != nil {
						return fmt.Errorf("setting Diagnostic Logs Configuration for Windows %s: %+v", id, err)
					}
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(webAppSlot.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Windows %s: %+v", id, err)
			}

			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfigurationSlot(ctx, id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(webAppSlot.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id, *storageConfig); err != nil {
					return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(webAppSlot.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStringsSlot(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if webAppSlot.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublishSlot(ctx, client, id, webAppSlot.ZipDeployFile); err != nil {
					return err
				}
			}

			if !webAppSlot.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateScmAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if !webAppSlot.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateFtpAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}
