// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_function_app_slot

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_function_app_slot/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient
			id, err := webapps.ParseSlotID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			functionAppSlot, err := client.GetSlot(ctx, *id)
			if err != nil {
				if response.WasNotFound(functionAppSlot.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			functionAppId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)

			appSettingsResp, err := client.ListApplicationSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", id, err)
			}

			connectionStrings, err := client.ListConnectionStringsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccountsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentialsSlot(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if auth.Model != nil && auth.Model.Properties != nil && strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2Slot(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfigurationSlot(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logs, err := client.GetDiagnosticLogsConfigurationSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowedSlot(ctx, *id); err != nil || basicAuthFTPResp.Model == nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowedSlot(ctx, *id); err != nil || basicAuthWebDeployResp.Model == nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}
			state := WindowsFunctionAppSlotModel{
				Name:                             id.SlotName,
				FunctionAppID:                    functionAppId.ID(),
				PublishingFTPBasicAuthEnabled:    basicAuthFTP,
				PublishingDeployBasicAuthEnabled: basicAuthWebDeploy,
				ConnectionStrings:                helpers.FlattenConnectionStrings(connectionStrings.Model),
				SiteCredentials:                  helpers.FlattenSiteCredentials(siteCredentials),
				AuthSettings:                     helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:                   helpers.FlattenAuthV2Settings(authV2),
				Backup:                           helpers.FlattenBackupConfig(backup.Model),
				StorageAccounts:                  helpers.FlattenStorageAccounts(storageAccounts.Model),
			}
			if model := functionAppSlot.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				state.Kind = pointer.From(model.Kind)
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				state.Identity = pointer.From(flattenedIdentity)

				if props := model.Properties; props != nil {
					state.Enabled = pointer.From(props.Enabled)
					state.ClientCertMode = pointer.FromEnum(props.ClientCertMode)
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.DailyMemoryTimeQuota = pointer.From(props.DailyMemoryTimeQuota)
					state.KeyVaultReferenceIdentityID = pointer.From(props.KeyVaultReferenceIdentity)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
					state.VirtualNetworkBackupRestoreEnabled = pointer.From(props.VnetBackupRestoreEnabled)
					state.VnetImagePullEnabled = pointer.From(props.VnetImagePullEnabled)
					state.EndToEndTLSEncryptionEnabled = pointer.From(props.EndToEndEncryptionEnabled)

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						state.HostingEnvId = pointer.From(hostingEnv.Id)
					}

					functionApp, err := client.Get(ctx, functionAppId)
					if err != nil {
						return fmt.Errorf("reading parent Function App for Windows %s: %+v", *id, err)
					}
					if functionApp.Model == nil || functionApp.Model.Properties == nil || functionApp.Model.Properties.ServerFarmId == nil {
						return fmt.Errorf("reading parent Function App Service Plan information for Windows %s: %+v", *id, err)
					}
					parentAppFarmId, err := commonids.ParseAppServicePlanIDInsensitively(*functionApp.Model.Properties.ServerFarmId)
					if err != nil {
						return fmt.Errorf("parsing parent Service Plan ID for %s: %+v", *id, err)
					}

					if slotPlanIdRaw := props.ServerFarmId; slotPlanIdRaw != nil && !strings.EqualFold(parentAppFarmId.ID(), *slotPlanIdRaw) {
						slotPlanId, err := commonids.ParseAppServicePlanIDInsensitively(*slotPlanIdRaw)
						if err != nil {
							return fmt.Errorf("parsing Service Plan ID for %s: %+v", *id, err)
						}
						state.ServicePlanID = slotPlanId.ID()
					}

					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)

					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						state.VirtualNetworkSubnetID = subnetId
					}
				}

				configResp, err := client.GetConfigurationSlot(ctx, *id)
				if err != nil {
					return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionAppSlot(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
				}
				state.SiteConfig = []helpers.SiteConfigWindowsFunctionAppSlot{*siteConfig}

				state.unpackWindowsFunctionAppSettings(appSettingsResp.Model, metadata)

				state.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs.Model)

				if err := metadata.Encode(&state); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}
			}

			return nil
		},
	}
}

func (m *WindowsFunctionAppSlotModel) unpackWindowsFunctionAppSettings(input *webapps.StringDictionary, metadata sdk.ResourceMetaData) {
	if input == nil || input.Properties == nil {
		return
	}

	appSettings := make(map[string]string)
	var dockerSettings helpers.ApplicationStackDocker
	m.BuiltinLogging = false

	for k, v := range *input.Properties {
		switch k {
		case "FUNCTIONS_EXTENSION_VERSION":
			m.FunctionExtensionsVersion = v

		case "WEBSITE_NODE_DEFAULT_VERSION":
			if len(m.SiteConfig[0].ApplicationStack) == 0 {
				m.SiteConfig[0].ApplicationStack = []helpers.ApplicationStackWindowsFunctionApp{{}}
			}
			m.SiteConfig[0].ApplicationStack[0].NodeVersion = v

		case "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"); ok {
				appSettings[k] = v
			}

		case "WEBSITE_CONTENTSHARE":
			if _, ok := metadata.ResourceData.GetOk("app_settings.WEBSITE_CONTENTSHARE"); ok {
				appSettings[k] = v
			}

		case "WEBSITE_HTTPLOGGING_RETENTION_DAYS":
		case "FUNCTIONS_WORKER_RUNTIME":
			if _, ok := metadata.ResourceData.GetOk("app_settings.FUNCTIONS_WORKER_RUNTIME"); ok {
				appSettings[k] = v
			}
			switch v {
			case "dotnet-isolated":
				m.SiteConfig[0].ApplicationStack[0].DotNetIsolated = true
			case "custom":
				m.SiteConfig[0].ApplicationStack[0].CustomHandler = true
			}
		case "DOCKER_REGISTRY_SERVER_URL":
			dockerSettings.RegistryURL = v

		case "DOCKER_REGISTRY_SERVER_USERNAME":
			dockerSettings.RegistryUsername = v

		case "DOCKER_REGISTRY_SERVER_PASSWORD":
			dockerSettings.RegistryPassword = v

		case "APPINSIGHTS_INSTRUMENTATIONKEY":
			m.SiteConfig[0].AppInsightsInstrumentationKey = v

		case "APPLICATIONINSIGHTS_CONNECTION_STRING":
			m.SiteConfig[0].AppInsightsConnectionString = v

		case "AzureWebJobsStorage":
			if strings.HasPrefix(v, "@Microsoft.KeyVault") {
				m.StorageKeyVaultSecretID = strings.TrimPrefix(strings.TrimSuffix(v, ")"), "@Microsoft.KeyVault(SecretUri=")
			} else {
				m.StorageAccountName, m.StorageAccountKey = helpers.ParseWebJobsStorageString(v)
			}
		case "AzureWebJobsDashboard":
			m.BuiltinLogging = true

		case "WEBSITE_HEALTHCHECK_MAXPINGFAILURES":
			i, _ := strconv.Atoi(v)
			m.SiteConfig[0].HealthCheckEvictionTime = int64(i)

		case "AzureWebJobsStorage__accountName":
			m.StorageUsesMSI = true
			m.StorageAccountName = v

		case "AzureWebJobsDashboard__accountName":
			m.BuiltinLogging = true

		case "WEBSITE_VNET_ROUTE_ALL":
			// Filter out - handled by site_config setting `vnet_route_all_enabled`

		default:
			appSettings[k] = v
		}
	}

	m.AppSettings = appSettings
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsFunctionAppSlotV0toV1{},
		},
	}
}
