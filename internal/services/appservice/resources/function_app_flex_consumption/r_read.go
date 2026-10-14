// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_flex_consumption

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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			functionAppFlexConsumption, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(functionAppFlexConsumption.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving App Settings for %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Connection String information for %s: %+v", id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, *id)
			if err != nil || stickySettings.Model == nil {
				return fmt.Errorf("retrieving Sticky Settings for %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Auth Settings for %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if auth.Model != nil && auth.Model.Properties != nil && strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2(ctx, *id)
				if err != nil {
					return fmt.Errorf("retrieving authV2 settings for %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowed(ctx, *id); err != nil && basicAuthWebDeployResp.Model != nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			model := functionAppFlexConsumption.Model
			if model == nil {
				return fmt.Errorf("function app %s : model is nil", id)
			}
			flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			state := FunctionAppFlexConsumptionModel{
				Name:                             id.SiteName,
				ResourceGroup:                    id.ResourceGroupName,
				Location:                         location.Normalize(model.Location),
				ConnectionStrings:                helpers.FlattenConnectionStrings(connectionStrings.Model),
				StickySettings:                   helpers.FlattenStickySettings(stickySettings.Model.Properties),
				SiteCredentials:                  helpers.FlattenSiteCredentials(siteCredentials),
				AuthSettings:                     helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:                   helpers.FlattenAuthV2Settings(authV2),
				PublishingDeployBasicAuthEnabled: basicAuthWebDeploy,
				Tags:                             pointer.From(model.Tags),
				Kind:                             pointer.From(model.Kind),
				Identity:                         pointer.From(flattenedIdentity),
			}

			if props := model.Properties; props != nil {
				state.Enabled = pointer.From(props.Enabled)
				state.ClientCertMode = pointer.FromEnum(props.ClientCertMode)
				state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
				state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
				state.DefaultHostname = pointer.From(props.DefaultHostName)
				state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)

				servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
				if err != nil {
					return err
				}
				state.ServicePlanId = servicePlanId.ID()

				if v := props.OutboundIPAddresses; v != nil {
					state.OutboundIPAddresses = *v
					state.OutboundIPAddressList = strings.Split(*v, ",")
				}

				if v := props.PossibleOutboundIPAddresses; v != nil {
					state.PossibleOutboundIPAddresses = *v
					state.PossibleOutboundIPAddressList = strings.Split(*v, ",")
				}

				configResp, err := client.GetConfiguration(ctx, *id)
				if err != nil {
					return fmt.Errorf("retrieving Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigFunctionAppFlexConsumption(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("retrieving Site Config for %s: %+v", id, err)
				}
				state.SiteConfig = []helpers.SiteConfigFunctionAppFlexConsumption{*siteConfig}

				if functionAppConfig := props.FunctionAppConfig; functionAppConfig != nil {
					if faConfigDeployment := functionAppConfig.Deployment; faConfigDeployment != nil && faConfigDeployment.Storage != nil {
						storageConfig := *faConfigDeployment.Storage
						state.StorageContainerType = pointer.FromEnum(storageConfig.Type)
						state.StorageContainerEndpoint = pointer.From(storageConfig.Value)
						if storageConfig.Authentication != nil && storageConfig.Authentication.Type != nil {
							state.StorageAuthType = pointer.FromEnum(storageConfig.Authentication.Type)
							if storageConfig.Authentication.UserAssignedIdentityResourceId != nil {
								state.StorageUserAssignedIdentityID = pointer.From(storageConfig.Authentication.UserAssignedIdentityResourceId)
							}
						}
					}

					if faConfigRuntime := functionAppConfig.Runtime; faConfigRuntime != nil {
						state.RuntimeName = pointer.FromEnum(faConfigRuntime.Name)
						state.RuntimeVersion = pointer.From(faConfigRuntime.Version)
					}

					if faConfigScale := functionAppConfig.ScaleAndConcurrency; faConfigScale != nil {
						state.AlwaysReady = FlattenAlwaysReadyConfiguration(faConfigScale.AlwaysReady)
						state.InstanceMemoryInMB = pointer.From(faConfigScale.InstanceMemoryMB)
						state.MaximumInstanceCount = pointer.From(faConfigScale.MaximumInstanceCount)
						if faConfigScale.Triggers != nil && faConfigScale.Triggers.HTTP != nil {
							state.HttpConcurrency = pointer.From(faConfigScale.Triggers.HTTP.PerInstanceConcurrency)
						}
					}
				}

				state.unpackFunctionAppFlexConsumptionSettings(*appSettingsResp.Model)

				state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)

				state.HttpsOnly = pointer.From(props.HTTPSOnly)

				if props.VirtualNetworkSubnetId != nil && pointer.From(props.VirtualNetworkSubnetId) != "" {
					subnetId, err := commonids.ParseSubnetID(*props.VirtualNetworkSubnetId)
					if err != nil {
						return err
					}
					state.VirtualNetworkSubnetID = subnetId.ID()
				}

				// Zip Deploys are not retrievable, so attempt to get from config. This doesn't matter for imports as an unexpected value here could break the deployment.
				if deployFile, ok := metadata.ResourceData.Get("zip_deploy_file").(string); ok {
					state.ZipDeployFile = deployFile
				}
			}
			return metadata.Encode(&state)
		},
	}
}

const (
	StorageStringFmt = "DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s"
)

type FunctionAppAlwaysReady struct {
	Name          string `tfschema:"name"`
	InstanceCount int64  `tfschema:"instance_count"`
}

func (m *FunctionAppFlexConsumptionModel) unpackFunctionAppFlexConsumptionSettings(input webapps.StringDictionary) {
	if input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)

	for k, v := range *input.Properties {
		switch k {
		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = v

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = v

		case "AzureWebJobsStorage":
			_, m.StorageAccessKey = helpers.ParseWebJobsStorageString(v)

		case "WEBSITE_HEALTHCHECK_MAXPINGFAILURES":
			i, _ := strconv.Atoi(v)
			m.SiteConfig[0].HealthCheckEvictionTime = int64(i)

		case "DEPLOYMENT_STORAGE_CONNECTION_STRING":
			// Filter out - not user faced

		default:
			appSettings[k] = v
		}
	}
	m.AppSettings = appSettings
}

func ExpandAlwaysReadyConfiguration(input []FunctionAppAlwaysReady, maximumInstanceCount int64) (*[]webapps.FunctionsAlwaysReadyConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}
	var totalInstanceCount int64
	arList := make([]webapps.FunctionsAlwaysReadyConfig, 0)
	for _, v := range input {
		totalInstanceCount += v.InstanceCount
		arList = append(arList, webapps.FunctionsAlwaysReadyConfig{
			Name:          pointer.To(v.Name),
			InstanceCount: pointer.To(v.InstanceCount),
		})
	}

	if totalInstanceCount > maximumInstanceCount {
		return nil, fmt.Errorf("the total number of always-ready instances should not exceed the maximum scale out limit")
	}

	return &arList, nil
}

func FlattenAlwaysReadyConfiguration(alwaysReady *[]webapps.FunctionsAlwaysReadyConfig) []FunctionAppAlwaysReady {
	if alwaysReady == nil || len(*alwaysReady) == 0 {
		return []FunctionAppAlwaysReady{}
	}

	alwaysReadyList := make([]FunctionAppAlwaysReady, 0)

	for _, v := range *alwaysReady {
		alwaysReadyList = append(alwaysReadyList, FunctionAppAlwaysReady{
			Name:          pointer.From(v.Name),
			InstanceCount: pointer.From(v.InstanceCount),
		})
	}

	return alwaysReadyList
}
