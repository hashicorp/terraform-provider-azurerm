// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_web_app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var webApp WindowsWebAppModel
			if err := metadata.Decode(&webApp); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			resourceProvidersClient := metadata.Client.AppService.ResourceProvidersClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			baseId := commonids.NewAppServiceID(subscriptionId, webApp.ResourceGroup, webApp.Name)
			id, err := commonids.ParseWebAppID(baseId.ID())
			if err != nil {
				return err
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, *id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Windows %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			sc := webApp.SiteConfig[0]

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: webApp.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			servicePlanId, err := commonids.ParseAppServicePlanID(webApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading App %s: %+v", servicePlanId, err)
			}
			if servicePlan.Model != nil && servicePlan.Model.Properties != nil {
				if ase := servicePlan.Model.Properties.HostingEnvironmentProfile; ase != nil {
					// Attempt to check the ASE for the appropriate suffix for the name availability request. Not convinced
					// the `DNSSuffix` field is still valid and possibly should have been deprecated / removed as is legacy
					// setting from ASEv1? Hence the non-fatal approach here.
					nameSuffix := "appserviceenvironment.net"
					if ase.Id != nil {
						aseId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*ase.Id)
						nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
						if err != nil {
							metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", webApp.Name, servicePlanId)
						} else {
							existingASE, err := aseClient.Get(ctx, *aseId)
							if err != nil || existingASE.Model == nil {
								metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", webApp.Name, servicePlanId)
							} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
								nameSuffix = *props.DnsSuffix
							}
						}
					}
					availabilityRequest.Name = fmt.Sprintf("%s.%s", webApp.Name, nameSuffix)
					availabilityRequest.IsFqdn = pointer.To(true)
				}
				if servicePlan.Model.Sku != nil && servicePlan.Model.Sku.Name != nil {
					if helpers.IsFreeOrSharedServicePlan(*servicePlan.Model.Sku.Name) {
						if sc.AlwaysOn {
							return fmt.Errorf("always_on cannot be set to true when using Free, F1, D1 Sku")
						}
					}
				}
			}

			subscriptionID := commonids.NewSubscriptionID(subscriptionId)
			checkName, err := resourceProvidersClient.CheckNameAvailability(ctx, subscriptionID, availabilityRequest)
			if err != nil || checkName.Model == nil {
				return fmt.Errorf("checking name availability for %s: %+v", id, err)
			}
			if !*checkName.Model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Model.Message)
			}

			siteConfig, err := sc.ExpandForCreate(webApp.AppSettings)
			if err != nil {
				return err
			}

			currentStack := ""
			if len(sc.ApplicationStack) == 1 {
				currentStack = sc.ApplicationStack[0].CurrentStack
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(webApp.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(webApp.Location),
				Tags:     pointer.To(webApp.Tags),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					ServerFarmId:              pointer.To(webApp.ServicePlanId),
					Enabled:                   pointer.To(webApp.Enabled),
					HTTPSOnly:                 pointer.To(webApp.HttpsOnly),
					SiteConfig:                siteConfig,
					ClientAffinityEnabled:     pointer.To(webApp.ClientAffinityEnabled),
					ClientCertEnabled:         pointer.To(webApp.ClientCertEnabled),
					ClientCertMode:            pointer.ToEnum[webapps.ClientCertMode](webApp.ClientCertMode),
					VnetBackupRestoreEnabled:  pointer.To(webApp.VirtualNetworkBackupRestoreEnabled),
					VnetRouteAllEnabled:       siteConfig.VnetRouteAllEnabled,
					EndToEndEncryptionEnabled: pointer.To(webApp.EndToEndTLSEncryptionEnabled),
				},
			}

			siteEnvelope.Properties.VnetImagePullEnabled = pointer.To(webApp.VirtualNetworkImagePullEnabled)

			pna := helpers.PublicNetworkAccessEnabled
			if !webApp.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if webApp.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(webApp.KeyVaultReferenceIdentityID)
			}

			if webApp.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(webApp.VirtualNetworkSubnetID)
			}

			if webApp.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(webApp.ClientCertExclusionPaths)
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, *id, siteEnvelope, metadata.SetIDCallback(id)); err != nil {
				return fmt.Errorf("creating Windows %s: %+v", id, err)
			}

			metadata.SetID(id)

			if currentStack != "" {
				siteMetadata := webapps.StringDictionary{Properties: &map[string]string{
					"CURRENT_STACK": currentStack,
				}}
				if _, err := client.UpdateMetadata(ctx, *id, siteMetadata); err != nil {
					return fmt.Errorf("setting Site Metadata for Current Stack on Windows %s: %+v", id, err)
				}
			}

			appSettings := helpers.ExpandAppSettingsForUpdate(siteConfig.AppSettings)
			appSettingsProps := *appSettings.Properties
			if metadata.ResourceData.HasChange("site_config.0.health_check_eviction_time_in_min") {
				appSettingsProps["WEBSITE_HEALTHCHECK_MAXPINGFAILURES"] = strconv.FormatInt(webApp.SiteConfig[0].HealthCheckEvictionTime, 10)
				appSettings.Properties = &appSettingsProps
			}

			if appSettings != nil {
				if _, err := client.UpdateApplicationSettings(ctx, *id, *appSettings); err != nil {
					return fmt.Errorf("setting App Settings for Windows %s: %+v", id, err)
				}
			}

			stickySettings := helpers.ExpandStickySettings(webApp.StickySettings)

			if stickySettings != nil {
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: stickySettings,
				}
				if _, err := client.UpdateSlotConfigurationNames(ctx, *id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Windows %s: %+v", *id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(webApp.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettings(ctx, *id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", *id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(webApp.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2(ctx, *id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Windows %s: %+v", *id, err)
				}
			}

			if metadata.ResourceData.HasChange("logs") {
				logsConfig := helpers.ExpandLogsConfig(webApp.LogsConfig)
				if logsConfig.Properties != nil {
					if _, err := client.UpdateDiagnosticLogsConfig(ctx, *id, *logsConfig); err != nil {
						return fmt.Errorf("setting Diagnostic Logs Configuration for Windows %s: %+v", id, err)
					}
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(webApp.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Windows %s: %+v", *id, err)
			}

			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, *id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Windows %s: %+v", *id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(webApp.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccounts(ctx, *id, *storageConfig); err != nil {
					return fmt.Errorf("setting Storage Accounts for Windows %s: %+v", *id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(webApp.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Windows %s: %+v", id, err)
				}
			}

			if webApp.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, *id, webApp.ZipDeployFile); err != nil {
					return err
				}
			}

			if !webApp.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateScmAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", *id, err)
				}
			}

			if !webApp.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateFtpAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", *id, err)
				}
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}
