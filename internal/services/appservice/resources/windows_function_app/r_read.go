// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_function_app

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_function_app/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
			functionApp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(functionApp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", id, err)
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", *id, err)
			}

			connectionStrings, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", *id, err)
			}

			stickySettings, err := client.ListSlotConfigurationNames(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Sticky Settings for Windows %s: %+v", *id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccounts(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if auth.Model != nil && auth.Model.Properties != nil && strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfiguration(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", id, err)
				}
			}

			logs, err := client.GetDiagnosticLogsConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading logs configuration for Windows %s: %+v", id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowed(ctx, *id); err != nil || basicAuthFTPResp.Model == nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowed(ctx, *id); err != nil || basicAuthWebDeployResp.Model == nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			if model := functionApp.Model; model != nil {
				state := WindowsFunctionAppModel{
					Name:          id.SiteName,
					ResourceGroup: id.ResourceGroupName,
					Location:      location.Normalize(model.Location),
					Tags:          pointer.From(model.Tags),
					Kind:          pointer.From(model.Kind),
				}

				if props := model.Properties; props != nil {
					state.Enabled = pointer.From(props.Enabled)
					state.ClientCertMode = pointer.FromEnum(props.ClientCertMode)
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.DailyMemoryTimeQuota = pointer.From(props.DailyMemoryTimeQuota)
					state.StickySettings = helpers.FlattenStickySettings(stickySettings.Model.Properties)
					state.KeyVaultReferenceIdentityID = pointer.From(props.KeyVaultReferenceIdentity)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
					state.VirtualNetworkBackupRestoreEnabled = pointer.From(props.VnetBackupRestoreEnabled)
					state.VnetImagePullEnabled = pointer.From(props.VnetImagePullEnabled)
					state.EndToEndTLSEncryptionEnabled = pointer.From(props.EndToEndEncryptionEnabled)

					servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmId))
					if err != nil {
						return err
					}
					state.ServicePlanId = servicePlanId.ID()

					state.PublishingFTPBasicAuthEnabled = basicAuthFTP
					state.PublishingDeployBasicAuthEnabled = basicAuthWebDeploy

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						hostingEnvId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*hostingEnv.Id)
						if err != nil {
							return err
						}
						state.HostingEnvId = hostingEnvId.ID()
					}

					if v := props.OutboundIPAddresses; v != nil {
						state.OutboundIPAddresses = *v
						state.OutboundIPAddressList = strings.Split(*v, ",")
					}

					if v := props.PossibleOutboundIPAddresses; v != nil {
						state.PossibleOutboundIPAddresses = *v
						state.PossibleOutboundIPAddressList = strings.Split(*v, ",")
					}

					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)

					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						state.VirtualNetworkSubnetID = subnetId
					}
				}
				configResp, err := client.GetConfiguration(ctx, *id)
				if err != nil {
					return fmt.Errorf("making Read request on AzureRM Function App Configuration %q: %+v", id.SiteName, err)
				}

				siteConfig, err := helpers.FlattenSiteConfigWindowsFunctionApp(configResp.Model.Properties)
				if err != nil {
					return fmt.Errorf("reading Site Config for Windows %s: %+v", id, err)
				}

				state.SiteConfig = []helpers.SiteConfigWindowsFunctionApp{*siteConfig}

				state.unpackWindowsFunctionAppSettings(appSettingsResp.Model, metadata)

				state.ConnectionStrings = helpers.FlattenConnectionStrings(connectionStrings.Model)

				state.SiteCredentials = helpers.FlattenSiteCredentials(siteCredentials)

				state.AuthSettings = helpers.FlattenAuthSettings(auth.Model)

				state.AuthV2Settings = helpers.FlattenAuthV2Settings(authV2)

				state.Backup = helpers.FlattenBackupConfig(backup.Model)

				state.SiteConfig[0].AppServiceLogs = helpers.FlattenFunctionAppAppServiceLogs(logs.Model)

				state.StorageAccounts = helpers.FlattenStorageAccounts(storageAccounts.Model)

				// Zip Deploys are not retrievable, so attempt to get from config. This doesn't matter for imports as an unexpected value here could break the deployment.
				if deployFile, ok := metadata.ResourceData.Get("zip_deploy_file").(string); ok {
					state.ZipDeployFile = deployFile
				}
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				if err := metadata.ResourceData.Set("identity", flattenedIdentity); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if err := metadata.Encode(&state); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}
			}

			return nil
		},
	}
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.AppService.WebAppsClient
		servicePlanClient := metadata.Client.AppService.ServicePlanClient

		id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		site, err := client.Get(ctx, *id)
		if err != nil || site.Model == nil || site.Model.Properties == nil {
			return fmt.Errorf("reading Windows %s: %+v", id, err)
		}
		props := site.Model.Properties
		if props.ServerFarmId == nil {
			return fmt.Errorf("determining Service Plan ID for Windows %s: %+v", id, err)
		}
		servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(props.ServerFarmId))
		if err != nil {
			return err
		}

		sp, err := servicePlanClient.Get(ctx, *servicePlanId)
		if err != nil || sp.Model == nil || sp.Model.Kind == nil {
			return fmt.Errorf("reading Service Plan for Windows %s: %+v", id, err)
		}

		if strings.Contains(strings.ToLower(*sp.Model.Kind), "Windows") {
			return fmt.Errorf("specified Service Plan is not a Windows Functionapp plan")
		}

		return nil
	}
}

func (m *WindowsFunctionAppModel) unpackWindowsFunctionAppSettings(input *webapps.StringDictionary, metadata sdk.ResourceMetaData) {
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
			0: migration.WindowsFunctionAppV0toV1{},
		},
	}
}
