// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_function_app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var functionApp LinuxFunctionAppModel

			if err := metadata.Decode(&functionApp); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			resourcesClient := metadata.Client.AppService.ResourceProvidersClient
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := commonids.NewAppServiceID(subscriptionId, functionApp.ResourceGroup, functionApp.Name)

			servicePlanId, err := commonids.ParseAppServicePlanID(functionApp.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			var planSKU *string
			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: functionApp.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}
			if servicePlanModel := servicePlan.Model; servicePlanModel != nil {
				if sku := servicePlanModel.Sku; sku != nil && sku.Name != nil {
					planSKU = sku.Name
				}

				if ase := servicePlanModel.Properties.HostingEnvironmentProfile; ase != nil {
					// Attempt to check the ASE for the appropriate suffix for the name availability request.
					// This varies between internal and external ASE Types, and potentially has other names in other clouds
					// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
					nameSuffix := "appserviceenvironment.net"
					if ase.Id != nil {
						aseId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*ase.Id)
						nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
						if err != nil {
							metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
						} else {
							existingASE, err := aseClient.Get(ctx, *aseId)
							if err != nil || existingASE.Model == nil {
								metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionApp.Name, servicePlanId)
							} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
								nameSuffix = *props.DnsSuffix
							}
						}
					}

					availabilityRequest.Name = fmt.Sprintf("%s.%s", functionApp.Name, nameSuffix)
					availabilityRequest.IsFqdn = pointer.To(true)

					if !functionApp.VnetImagePullEnabled {
						return fmt.Errorf("`vnet_image_pull_enabled` cannot be disabled for app running in an app service environment")
					}
				}
			}
			// Only send for ElasticPremium and Consumption plan
			elasticOrConsumptionPlan := helpers.PlanIsElastic(planSKU) || helpers.PlanIsConsumption(planSKU)
			sendContentSettings := elasticOrConsumptionPlan && !functionApp.ForceDisableContentShare

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			subId := commonids.NewSubscriptionID(subscriptionId)

			checkName, err := resourcesClient.CheckNameAvailability(ctx, subId, availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for Linux %s: %+v", id, err)
			}
			if model := checkName.Model; model != nil && model.NameAvailable != nil && !*model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *model.Message)
			}

			storageString := functionApp.StorageAccountName
			if !functionApp.StorageUsesMSI {
				if functionApp.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, functionApp.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, functionApp.StorageAccountName, functionApp.StorageAccountKey, *storageDomainSuffix)
				}
			}

			siteConfig, err := helpers.ExpandSiteConfigLinuxFunctionApp(functionApp.SiteConfig, nil, metadata, functionApp.FunctionExtensionsVersion, storageString, functionApp.StorageUsesMSI)
			if err != nil {
				return fmt.Errorf("expanding site_config for Linux %s: %+v", id, err)
			}

			if functionApp.BuiltinLogging {
				if functionApp.AppSettings == nil {
					functionApp.AppSettings = make(map[string]string)
				}
				if !functionApp.StorageUsesMSI {
					functionApp.AppSettings["AzureWebJobsDashboard"] = storageString
				} else {
					functionApp.AppSettings["AzureWebJobsDashboard__accountName"] = functionApp.StorageAccountName
				}
			}
			// for function premium app with WEBSITE_CONTENTOVERVNET sets to 1, the user has to specify a predefined fire share.
			// https://docs.microsoft.com/azure/azure-functions/functions-app-settings#website_contentovervnet
			if sendContentSettings {
				if functionApp.AppSettings == nil {
					functionApp.AppSettings = make(map[string]string)
				}
				if !functionApp.StorageUsesMSI {
					suffix := uuid.New().String()[0:4]
					_, contentOverVnetEnabled := functionApp.AppSettings["WEBSITE_CONTENTOVERVNET"]
					_, contentSharePresent := functionApp.AppSettings["WEBSITE_CONTENTSHARE"]
					if _, contentShareConnectionStringPresent := functionApp.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"]; !contentShareConnectionStringPresent {
						functionApp.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = storageString
					}

					if !contentSharePresent {
						if contentOverVnetEnabled {
							return fmt.Errorf("the app_setting WEBSITE_CONTENTSHARE must be specified and set to a valid share when WEBSITE_CONTENTOVERVNET is specified")
						}
						functionApp.AppSettings["WEBSITE_CONTENTSHARE"] = fmt.Sprintf("%s-%s", strings.ToLower(functionApp.Name), suffix)
					}
				} else {
					if _, present := functionApp.AppSettings["AzureWebJobsStorage__accountName"]; !present {
						functionApp.AppSettings["AzureWebJobsStorage__accountName"] = storageString
					}
				}
			}

			siteConfig.LinuxFxVersion = helpers.EncodeFunctionAppLinuxFxVersion(functionApp.SiteConfig[0].ApplicationStack)
			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionApp.AppSettings)

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(functionApp.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(functionApp.Location),
				Tags:     pointer.To(functionApp.Tags),
				Kind:     pointer.To("functionapp,linux"),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					ServerFarmId:              pointer.To(functionApp.ServicePlanId),
					Enabled:                   pointer.To(functionApp.Enabled),
					HTTPSOnly:                 pointer.To(functionApp.HttpsOnly),
					SiteConfig:                siteConfig,
					ClientCertEnabled:         pointer.To(functionApp.ClientCertEnabled),
					ClientCertMode:            pointer.ToEnum[webapps.ClientCertMode](functionApp.ClientCertMode),
					DailyMemoryTimeQuota:      pointer.To(functionApp.DailyMemoryTimeQuota), // TODO - Investigate, setting appears silently ignored on Linux Function Apps?
					VnetBackupRestoreEnabled:  pointer.To(functionApp.VirtualNetworkBackupRestoreEnabled),
					VnetImagePullEnabled:      pointer.To(functionApp.VnetImagePullEnabled),
					VnetRouteAllEnabled:       siteConfig.VnetRouteAllEnabled,
					EndToEndEncryptionEnabled: pointer.To(functionApp.EndToEndTLSEncryptionEnabled),
				},
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !functionApp.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Values appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if functionApp.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(functionApp.KeyVaultReferenceIdentityID)
			}

			if functionApp.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(functionApp.VirtualNetworkSubnetID)
			}

			if functionApp.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(functionApp.ClientCertExclusionPaths)
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, siteEnvelope, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating Linux %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			// (@jackofallops) - updating the policy for publishing credentials resets the `Use32BitWorkerProcess` property
			if !functionApp.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateScmAllowed(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if !functionApp.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateFtpAllowed(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, siteEnvelope); err != nil {
				return fmt.Errorf("creating Linux %s: %+v", id, err)
			}

			stickySettings := helpers.ExpandStickySettings(functionApp.StickySettings)

			if stickySettings != nil {
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: stickySettings,
				}
				if _, err := client.UpdateSlotConfigurationNames(ctx, id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for Linux %s: %+v", id, err)
				}
			}

			backupConfig, err := helpers.ExpandBackupConfig(functionApp.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Linux %s: %+v", id, err)
			}
			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfiguration(ctx, id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Linux %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionApp.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettings(ctx, id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for Linux %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(functionApp.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2(ctx, id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Linux %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(functionApp.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccounts(ctx, id, *storageConfig); err != nil {
					return fmt.Errorf("setting Storage Accounts for Linux %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionApp.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
				}
			}

			if _, ok := metadata.ResourceData.GetOk("site_config.0.app_service_logs"); ok {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(functionApp.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			if functionApp.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, id, functionApp.ZipDeployFile); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
