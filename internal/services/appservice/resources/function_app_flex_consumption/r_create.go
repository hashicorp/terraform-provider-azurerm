// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_flex_consumption

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			resourcesClient := metadata.Client.AppService.ResourceProvidersClient
			servicePlanClient := metadata.Client.AppService.ServicePlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			var functionAppFlexConsumption FunctionAppFlexConsumptionModel
			if err := metadata.Decode(&functionAppFlexConsumption); err != nil {
				return err
			}

			id := commonids.NewAppServiceID(subscriptionId, functionAppFlexConsumption.ResourceGroup, functionAppFlexConsumption.Name)

			servicePlanId, err := commonids.ParseAppServicePlanID(functionAppFlexConsumption.ServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", servicePlanId, err)
			}

			var planSKU *string
			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: functionAppFlexConsumption.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}

			if servicePlanModel := servicePlan.Model; servicePlanModel != nil {
				if sku := servicePlanModel.Sku; sku != nil && sku.Name != nil {
					planSKU = sku.Name
				}
			}

			isFlexConsumptionSku := helpers.PlanIsFlexConsumption(planSKU)
			if !isFlexConsumptionSku {
				return fmt.Errorf("the sku name is %s which is not valid for a flex consumption function app", *planSKU)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			checkName, err := resourcesClient.CheckNameAvailability(ctx, commonids.NewSubscriptionID(subscriptionId), availabilityRequest)
			if err != nil {
				return fmt.Errorf("checking name availability for %s: %+v", id, err)
			}
			if model := checkName.Model; model != nil && model.NameAvailable != nil && !*model.NameAvailable {
				return fmt.Errorf("the Site Name %q failed the availability check: %+v", id.SiteName, *model.Message)
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(functionAppFlexConsumption.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			storageDeployment := &webapps.FunctionsDeployment{
				Storage: &webapps.FunctionsDeploymentStorage{
					Type:  pointer.ToEnum[webapps.FunctionsDeploymentStorageType](functionAppFlexConsumption.StorageContainerType),
					Value: &functionAppFlexConsumption.StorageContainerEndpoint,
				},
			}
			storageConnStringForFCApp := "DEPLOYMENT_STORAGE_CONNECTION_STRING"
			var storageString string
			if storageName, _, ok := strings.Cut(strings.TrimPrefix(functionAppFlexConsumption.StorageContainerEndpoint, "https://"), "."); ok {
				storageString = fmt.Sprintf(helpers.StorageStringFmt, storageName, functionAppFlexConsumption.StorageAccessKey, *storageDomainSuffix)
			} else {
				return fmt.Errorf("retrieving storage container endpoint error, the expected format is https://storagename.blob.core.windows.net/containername, the received value is %s", functionAppFlexConsumption.StorageContainerEndpoint)
			}
			storageAuth := webapps.FunctionsDeploymentStorageAuthentication{
				Type: pointer.ToEnum[webapps.AuthenticationType](functionAppFlexConsumption.StorageAuthType),
			}

			if functionAppFlexConsumption.StorageAuthType == string(webapps.AuthenticationTypeStorageAccountConnectionString) {
				if functionAppFlexConsumption.StorageAccessKey == "" {
					return fmt.Errorf("the storage account access key must be specified when using the storage key based access")
				}
			} else {
				storageConnStringForFCApp = ""
				if functionAppFlexConsumption.StorageAuthType == string(webapps.AuthenticationTypeUserAssignedIdentity) {
					if functionAppFlexConsumption.StorageUserAssignedIdentityID == "" {
						return fmt.Errorf("the user assigned identity id must be specified when using the user assigned identity to access the storage account")
					}
					storageAuth.UserAssignedIdentityResourceId = &functionAppFlexConsumption.StorageUserAssignedIdentityID
				}
			}

			storageAuth.StorageAccountConnectionStringName = &storageConnStringForFCApp
			storageDeployment.Storage.Authentication = &storageAuth
			runtime := webapps.FunctionsRuntime{
				Name:    pointer.ToEnum[webapps.RuntimeName](functionAppFlexConsumption.RuntimeName),
				Version: &functionAppFlexConsumption.RuntimeVersion,
			}

			alwaysReady, err := ExpandAlwaysReadyConfiguration(functionAppFlexConsumption.AlwaysReady, functionAppFlexConsumption.MaximumInstanceCount)
			if err != nil {
				return fmt.Errorf("expanding `always_ready` for %s: %+v", id, err)
			}

			scaleAndConcurrencyConfig := webapps.FunctionsScaleAndConcurrency{
				AlwaysReady:          alwaysReady,
				InstanceMemoryMB:     &functionAppFlexConsumption.InstanceMemoryInMB,
				MaximumInstanceCount: &functionAppFlexConsumption.MaximumInstanceCount,
			}

			if functionAppFlexConsumption.HttpConcurrency >= 1 {
				scaleAndConcurrencyConfig.Triggers = &webapps.FunctionsScaleAndConcurrencyTriggers{
					HTTP: &webapps.FunctionsScaleAndConcurrencyTriggersHTTP{
						PerInstanceConcurrency: &functionAppFlexConsumption.HttpConcurrency,
					},
				}
			}

			flexFunctionAppConfig := &webapps.FunctionAppConfig{
				Deployment:          storageDeployment,
				Runtime:             &runtime,
				ScaleAndConcurrency: &scaleAndConcurrencyConfig,
			}

			siteConfig, err := helpers.ExpandSiteConfigFunctionFlexConsumptionApp(functionAppFlexConsumption.SiteConfig, nil, metadata, false, storageString, storageConnStringForFCApp)
			if err != nil {
				return fmt.Errorf("expanding `site_config` for %s: %+v", id, err)
			}

			siteConfig.AppSettings = helpers.MergeUserAppSettings(siteConfig.AppSettings, functionAppFlexConsumption.AppSettings)

			siteEnvelope := webapps.Site{
				Location: location.Normalize(functionAppFlexConsumption.Location),
				Tags:     pointer.To(functionAppFlexConsumption.Tags),
				Kind:     pointer.To("functionapp,linux"),
				Identity: expandedIdentity,
				Properties: &webapps.SiteProperties{
					ServerFarmId:      pointer.To(functionAppFlexConsumption.ServicePlanId),
					Enabled:           pointer.To(functionAppFlexConsumption.Enabled),
					SiteConfig:        siteConfig,
					HTTPSOnly:         pointer.To(functionAppFlexConsumption.HttpsOnly),
					FunctionAppConfig: flexFunctionAppConfig,
					ClientCertEnabled: pointer.To(functionAppFlexConsumption.ClientCertEnabled),
					ClientCertMode:    pointer.ToEnum[webapps.ClientCertMode](functionAppFlexConsumption.ClientCertMode),
				},
			}

			pna := helpers.PublicNetworkAccessEnabled
			if !functionAppFlexConsumption.PublicNetworkAccess {
				pna = helpers.PublicNetworkAccessDisabled
			}

			siteEnvelope.Properties.PublicNetworkAccess = pointer.To(pna)
			siteEnvelope.Properties.SiteConfig.PublicNetworkAccess = siteEnvelope.Properties.PublicNetworkAccess

			if functionAppFlexConsumption.VirtualNetworkSubnetID != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(functionAppFlexConsumption.VirtualNetworkSubnetID)
				locks.ByID(functionAppFlexConsumption.VirtualNetworkSubnetID)
				defer locks.UnlockByID(functionAppFlexConsumption.VirtualNetworkSubnetID)
			}

			if functionAppFlexConsumption.ClientCertExclusionPaths != "" {
				siteEnvelope.Properties.ClientCertExclusionPaths = pointer.To(functionAppFlexConsumption.ClientCertExclusionPaths)
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, siteEnvelope, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if !functionAppFlexConsumption.PublishingDeployBasicAuthEnabled {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{},
				}
				if _, err := client.UpdateScmAllowed(ctx, id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			stickySettings := helpers.ExpandStickySettings(functionAppFlexConsumption.StickySettings)

			if stickySettings != nil {
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: stickySettings,
				}
				if _, err := client.UpdateSlotConfigurationNames(ctx, id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for %s: %+v", id, err)
				}
			}

			auth := helpers.ExpandAuthSettings(functionAppFlexConsumption.AuthSettings)
			if auth.Properties != nil {
				if _, err := client.UpdateAuthSettings(ctx, id, *auth); err != nil {
					return fmt.Errorf("setting Authorisation Settings for %s: %+v", id, err)
				}
			}

			authv2 := helpers.ExpandAuthV2Settings(functionAppFlexConsumption.AuthV2Settings)
			if authv2.Properties != nil {
				if _, err = client.UpdateAuthSettingsV2(ctx, id, *authv2); err != nil {
					return fmt.Errorf("updating AuthV2 settings for %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(functionAppFlexConsumption.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for %s: %+v", id, err)
				}
			}

			if functionAppFlexConsumption.ZipDeployFile != "" {
				if err = helpers.GetCredentialsAndPublish(ctx, client, id, functionAppFlexConsumption.ZipDeployFile); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
