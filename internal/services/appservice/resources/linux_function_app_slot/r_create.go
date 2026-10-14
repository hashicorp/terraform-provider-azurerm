// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_function_app_slot

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
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var functionAppSlot LinuxFunctionAppSlotModel

			if err := metadata.Decode(&functionAppSlot); err != nil {
				return err
			}

			client := metadata.Client.AppService.WebAppsClient
			resourceProvidersClient := metadata.Client.AppService.ResourceProvidersClient
			functionAppId, err := commonids.ParseFunctionAppID(functionAppSlot.FunctionAppID)
			if err != nil {
				return err
			}

			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := webapps.NewSlotID(subscriptionId, functionAppId.ResourceGroupName, functionAppId.SiteName, functionAppSlot.Name)

			functionApp, err := client.Get(ctx, *functionAppId)
			if err != nil {
				return fmt.Errorf("retrieving parent Linux %s: %+v", *functionAppId, err)
			}

			var servicePlanId *commonids.AppServicePlanId
			servicePlanId, err = commonids.ParseAppServicePlanIDInsensitively(*functionApp.Model.Properties.ServerFarmId)
			if err != nil {
				return err
			}

			if functionAppSlot.ServicePlanID != "" {
				newServicePlanId, err := commonids.ParseAppServicePlanID(functionAppSlot.ServicePlanID)
				if err != nil {
					return err
				}
				// we only set `service_plan_id` when it differs from the parent `service_plan_id` which is causing issues
				// https://github.com/hashicorp/terraform-provider-azurerm/issues/21024
				// we'll error here if the `service_plan_id` equals the parent `service_plan_id`
				if strings.EqualFold(newServicePlanId.ID(), servicePlanId.ID()) {
					return fmt.Errorf("`service_plan_id` should only be specified when it differs from the `service_plan_id` of the associated Web App")
				}

				servicePlanId = newServicePlanId
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: fmt.Sprintf("%s-%s", id.SiteName, id.SlotName),
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			var planSKU *string
			if model := servicePlan.Model; model != nil {
				if sku := model.Sku; sku != nil && sku.Name != nil {
					planSKU = sku.Name
				}
				if ase := model.Properties.HostingEnvironmentProfile; ase != nil {
					// Attempt to check the ASE for the appropriate suffix for the name availability request.
					// This varies between internal and external ASE Types, and potentially has other names in other clouds
					// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
					nameSuffix := "appserviceenvironment.net"
					if ase.Id != nil {
						aseId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*ase.Id)
						nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
						if err != nil {
							metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionAppSlot.Name, servicePlanId)
						} else {
							existingASE, err := aseClient.Get(ctx, *aseId)
							if err != nil || existingASE.Model == nil {
								metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", functionAppSlot.Name, servicePlanId)
							} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
								nameSuffix = *props.DnsSuffix
							}
						}
					}

					availabilityRequest.Name = fmt.Sprintf("%s.%s", functionAppSlot.Name, nameSuffix)
					availabilityRequest.IsFqdn = pointer.To(true)

					if !functionAppSlot.VnetImagePullEnabled {
						return fmt.Errorf("`vnet_image_pull_enabled` cannot be disabled for app running in an app service environment")
					}
				}
			}

			// Only send for ElasticPremium and Consumption plan
			elasticOrConsumptionPlan := helpers.PlanIsElastic(planSKU) || helpers.PlanIsConsumption(planSKU)
			sendContentSettings := elasticOrConsumptionPlan && !functionAppSlot.ForceDisableContentShare

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetSlot(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			subscriptionID := commonids.NewSubscriptionID(subscriptionId)

			checkName, err := resourceProvidersClient.CheckNameAvailability(ctx, subscriptionID, availabilityRequest)
			if err != nil || checkName.Model == nil {
				return fmt.Errorf("checking name availability for Linux %s: %+v", id, err)
			}
			if checkName.Model.NameAvailable != nil && !*checkName.Model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *checkName.Model.Message)
			}

			storageString := functionAppSlot.StorageAccountName
			if !functionAppSlot.StorageUsesMSI {
				if functionAppSlot.StorageKeyVaultSecretID != "" {
					storageString = fmt.Sprintf(helpers.StorageStringFmtKV, functionAppSlot.StorageKeyVaultSecretID)
				} else {
					storageString = fmt.Sprintf(helpers.StorageStringFmt, functionAppSlot.StorageAccountName, functionAppSlot.StorageAccountKey, *storageDomainSuffix)
				}
			}
			siteConfig, err := helpers.ExpandSiteConfigLinuxFunctionAppSlot(functionAppSlot.SiteConfig, nil, metadata, functionAppSlot.FunctionExtensionsVersion, storageString, functionAppSlot.StorageUsesMSI)
			if err != nil {
				return fmt.Errorf("expanding site_config for Linux %s: %+v", id, err)
			}

			if functionAppSlot.BuiltinLogging {
				if functionAppSlot.AppSettings == nil {
					functionAppSlot.AppSettings = make(map[string]string)
				}
				if !functionAppSlot.StorageUsesMSI {
					functionAppSlot.AppSettings["AzureWebJobsDashboard"] = storageString
				} else {
					functionAppSlot.AppSettings["AzureWebJobsDashboard__accountName"] = functionAppSlot.StorageAccountName
				}
			}

			if sendContentSettings {
				if functionAppSlot.AppSettings == nil {
					functionAppSlot.AppSettings = make(map[string]string)
				}
				if !functionAppSlot.StorageUsesMSI {
					suffix := uuid.New().String()[0:4]
					if _, present := functionAppSlot.AppSettings["WEBSITE_CONTENTSHARE"]; !present {
						functionAppSlot.AppSettings["WEBSITE_CONTENTSHARE"] = fmt.Sprintf("%s-%s", strings.ToLower(functionAppSlot.Name), suffix)
					}
					if _, present := functionAppSlot.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"]; !present {
						functionAppSlot.AppSettings["WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"] = storageString
					}
				} else {
					if _, present := functionAppSlot.AppSettings["AzureWebJobsStorage__accountName"]; !present {
						functionAppSlot.AppSettings["AzureWebJobsStorage__accountName"] = storageString
					}
				}
			}

			siteConfig.LinuxFxVersion = helpers.EncodeFunctionAppLinuxFxVersion(functionAppSlot.SiteConfig[0].ApplicationStack)
			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionAppSlot.AppSettings)

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(functionAppSlot.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Location: location.Normalize(functionApp.Model.Location),
				Tags:     pointer.To(functionAppSlot.Tags),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					ServerFarmId:              pointer.To(servicePlanId.ID()),
					Enabled:                   pointer.To(functionAppSlot.Enabled),
					HTTPSOnly:                 pointer.To(functionAppSlot.HttpsOnly),
					SiteConfig:                siteConfig,
					ClientCertEnabled:         pointer.To(functionAppSlot.ClientCertEnabled),
					ClientCertMode:            pointer.ToEnum[webapps.ClientCertMode](functionAppSlot.ClientCertMode),
					DailyMemoryTimeQuota:      pointer.To(functionAppSlot.DailyMemoryTimeQuota),
					VnetBackupRestoreEnabled:  pointer.To(functionAppSlot.VirtualNetworkBackupRestoreEnabled),
					VnetImagePullEnabled:      pointer.To(functionAppSlot.VnetImagePullEnabled),
					VnetRouteAllEnabled:       siteConfig.VnetRouteAllEnabled, // (@jackofallops) - Value appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
					EndToEndEncryptionEnabled: pointer.To(functionAppSlot.EndToEndTLSEncryptionEnabled),
				},
			}

			pan := helpers.PublicNetworkAccessEnabled
			if !functionAppSlot.PublicNetworkAccess {
				pan = helpers.PublicNetworkAccessDisabled
			}

			// (@jackofallops) - Value appear to need to be set in both SiteProperties and SiteConfig for now? https://github.com/Azure/azure-rest-api-specs/issues/24681
			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pan)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if functionAppSlot.KeyVaultReferenceIdentityID != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(functionAppSlot.KeyVaultReferenceIdentityID)
			}

			if functionAppSlot.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(functionAppSlot.VirtualNetworkSubnetID)
			}

			if functionAppSlot.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(functionAppSlot.ClientCertExclusionPaths)
			}

			if err := client.CreateOrUpdateSlotCallbackThenPoll(ctx, id, siteEnvelope, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating Linux %s: %+v", id, err)
			}

			metadata.SetID(id)

			if !functionAppSlot.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateScmAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if !functionAppSlot.PublishingFTPBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateFtpAllowedSlot(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for ftp publishing credentials for %s: %+v", id, err)
				}
			}

			if err := client.CreateOrUpdateSlotThenPoll(ctx, id, siteEnvelope); err != nil {
				return fmt.Errorf("updating properties of Linux %s: %+v", id, err)
			}

			backupConfig, err := helpers.ExpandBackupConfig(functionAppSlot.Backup)
			if err != nil {
				return fmt.Errorf("expanding backup configuration for Linux %s: %+v", id, err)
			}

			if backupConfig.Properties != nil {
				if _, err := client.UpdateBackupConfigurationSlot(ctx, id, *backupConfig); err != nil {
					return fmt.Errorf("adding Backup Settings for Linux %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionAppSlot.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettingsSlot(ctx, id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for Linux %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(functionAppSlot.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2Slot(ctx, id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for Linux %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionAppSlot.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStringsSlot(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
				}
			}

			storageConfig := helpers.ExpandStorageConfig(functionAppSlot.StorageAccounts)
			if storageConfig.Properties != nil {
				if _, err := client.UpdateAzureStorageAccountsSlot(ctx, id, *storageConfig); err != nil {
					return fmt.Errorf("setting Storage Accounts for Linux %s: %+v", id, err)
				}
			}

			if _, ok := metadata.ResourceData.GetOk("site_config.0.app_service_logs"); ok {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(functionAppSlot.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfigSlot(ctx, id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
