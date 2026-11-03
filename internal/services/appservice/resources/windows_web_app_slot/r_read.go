// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_web_app_slot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_web_app_slot/migration"
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

			webAppSlot, err := client.GetSlot(ctx, *id)
			if err != nil {
				if response.WasNotFound(webAppSlot.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Windows %s: %+v", *id, err)
			}

			// Despite being part of the defined `Get` response model, site_config is always nil so we get it explicitly
			webAppSiteSlotConfig, err := client.GetConfigurationSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Config for Windows %s: %+v", *id, err)
			}

			auth, err := client.GetAuthSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Auth Settings for Windows %s: %+v", *id, err)
			}

			var authV2 webapps.SiteAuthSettingsV2
			if strings.EqualFold(pointer.From(auth.Model.Properties.ConfigVersion), "v2") {
				authV2Resp, err := client.GetAuthSettingsV2Slot(ctx, *id)
				if err != nil {
					return fmt.Errorf("reading authV2 settings for Linux %s: %+v", *id, err)
				}
				authV2 = *authV2Resp.Model
			}

			backup, err := client.GetBackupConfigurationSlot(ctx, *id)
			if err != nil {
				if !response.WasNotFound(backup.HttpResponse) {
					return fmt.Errorf("reading Backup Settings for Windows %s: %+v", *id, err)
				}
			}

			logsConfig, err := client.GetDiagnosticLogsConfigurationSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Diagnostic Logs information for Windows %s: %+v", *id, err)
			}

			appSettings, err := client.ListApplicationSettingsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading App Settings for Windows %s: %+v", *id, err)
			}

			storageAccounts, err := client.ListAzureStorageAccountsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Storage Account information for Windows %s: %+v", *id, err)
			}

			connectionStrings, err := client.ListConnectionStringsSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Connection String information for Windows %s: %+v", *id, err)
			}

			siteCredentials, err := helpers.ListPublishingCredentialsSlot(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing Site Publishing Credential information for %s: %+v", id, err)
			}

			siteMetadata, err := client.ListMetadataSlot(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Site Metadata for Windows %s: %+v", *id, err)
			}

			appId := commonids.NewAppServiceID(metadata.Client.Account.SubscriptionId, id.ResourceGroupName, id.SiteName)

			webApp, err := client.Get(ctx, appId)
			if err != nil {
				return fmt.Errorf("reading parent Web App for Linux %s: %+v", *id, err)
			}
			if webApp.Model.Properties == nil || webApp.Model.Properties.ServerFarmId == nil {
				return fmt.Errorf("reading parent Function App Service Plan information for Linux %s: %+v", *id, err)
			}

			basicAuthFTP := true
			if basicAuthFTPResp, err := client.GetFtpAllowedSlot(ctx, *id); err != nil || basicAuthFTPResp.Model == nil {
				return fmt.Errorf("retrieving state of FTP Basic Auth for %s: %+v", *id, err)
			} else if csmProps := basicAuthFTPResp.Model.Properties; csmProps != nil {
				basicAuthFTP = csmProps.Allow
			}

			basicAuthWebDeploy := true
			if basicAuthWebDeployResp, err := client.GetScmAllowedSlot(ctx, *id); err != nil || basicAuthWebDeployResp.Model == nil {
				return fmt.Errorf("retrieving state of WebDeploy Basic Auth for %s: %+v", *id, err)
			} else if csmProps := basicAuthWebDeployResp.Model.Properties; csmProps != nil {
				basicAuthWebDeploy = csmProps.Allow
			}

			state := WindowsWebAppSlotModel{
				Name:              id.SlotName,
				AppServiceId:      appId.ID(),
				AuthSettings:      helpers.FlattenAuthSettings(auth.Model),
				AuthV2Settings:    helpers.FlattenAuthV2Settings(authV2),
				Backup:            helpers.FlattenBackupConfig(backup.Model),
				ConnectionStrings: helpers.FlattenConnectionStrings(connectionStrings.Model),
				LogsConfig:        helpers.FlattenLogsConfig(logsConfig.Model),
				SiteCredentials:   helpers.FlattenSiteCredentials(siteCredentials),
				StorageAccounts:   helpers.FlattenStorageAccounts(storageAccounts.Model),
			}

			if model := webAppSlot.Model; model != nil {
				state.Kind = pointer.From(model.Kind)
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.ClientAffinityEnabled = pointer.From(props.ClientAffinityEnabled)
					state.ClientCertEnabled = pointer.From(props.ClientCertEnabled)
					state.ClientCertMode = pointer.FromEnum(props.ClientCertMode)
					state.ClientCertExclusionPaths = pointer.From(props.ClientCertExclusionPaths)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.Enabled = pointer.From(props.Enabled)
					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.KeyVaultReferenceIdentityID = pointer.From(props.KeyVaultReferenceIdentity)
					state.OutboundIPAddresses = pointer.From(props.OutboundIPAddresses)
					state.OutboundIPAddressList = strings.Split(pointer.From(props.OutboundIPAddresses), ",")
					state.PossibleOutboundIPAddresses = pointer.From(props.PossibleOutboundIPAddresses)
					state.PossibleOutboundIPAddressList = strings.Split(pointer.From(props.PossibleOutboundIPAddresses), ",")
					state.PublicNetworkAccess = !strings.EqualFold(pointer.From(props.PublicNetworkAccess), helpers.PublicNetworkAccessDisabled)
					state.VirtualNetworkBackupRestoreEnabled = pointer.From(props.VnetBackupRestoreEnabled)
					state.VirtualNetworkImagePullEnabled = pointer.From(props.VnetImagePullEnabled)
					state.EndToEndTLSEncryptionEnabled = pointer.From(props.EndToEndEncryptionEnabled)

					if hostingEnv := props.HostingEnvironmentProfile; hostingEnv != nil {
						hostingEnvId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*hostingEnv.Id)
						if err != nil {
							return err
						}
						state.HostingEnvId = hostingEnvId.ID()
					}

					parentAppFarmId, err := commonids.ParseAppServicePlanIDInsensitively(*webApp.Model.Properties.ServerFarmId)
					if err != nil {
						return fmt.Errorf("reading parent Service Plan ID: %+v", err)
					}
					if slotPlanIdRaw := props.ServerFarmId; slotPlanIdRaw != nil && *slotPlanIdRaw != "" && !strings.EqualFold(parentAppFarmId.ID(), *slotPlanIdRaw) {
						slotPlanId, err := commonids.ParseAppServicePlanIDInsensitively(pointer.From(slotPlanIdRaw))
						if err != nil {
							return fmt.Errorf("reading Slot Service Plan ID: %+v", err)
						}
						state.ServicePlanID = slotPlanId.ID()
					}

					if subnetId := pointer.From(props.VirtualNetworkSubnetId); subnetId != "" {
						state.VirtualNetworkSubnetID = subnetId
					}
				}
				state.PublishingFTPBasicAuthEnabled = basicAuthFTP
				state.PublishingDeployBasicAuthEnabled = basicAuthWebDeploy

				state.AppSettings = helpers.FlattenWebStringDictionary(appSettings.Model)
				currentStack := ""
				if m := siteMetadata.Model; m != nil && m.Properties != nil {
					p := *m.Properties
					if v, ok := p["CURRENT_STACK"]; ok {
						currentStack = v
					}
				}

				siteConfig := helpers.SiteConfigWindowsWebAppSlot{}
				siteConfig.Flatten(webAppSiteSlotConfig.Model.Properties, currentStack)
				siteConfig.SetHealthCheckEvictionTime(state.AppSettings)
				state.AppSettings = siteConfig.ParseNodeVersion(state.AppSettings)

				if helpers.FxStringHasPrefix(siteConfig.WindowsFxVersion, helpers.FxStringPrefixDocker) {
					siteConfig.DecodeDockerAppStack(state.AppSettings)
				}

				state.SiteConfig = []helpers.SiteConfigWindowsWebAppSlot{siteConfig}

				// Filter out all settings we've consumed above
				state.AppSettings = helpers.FilterManagedAppSettings(state.AppSettings)

				// Zip Deploys are not retrievable, so attempt to get from config. This doesn't matter for imports as an unexpected value here could break the deployment.
				if deployFile, ok := metadata.ResourceData.Get("zip_deploy_file").(string); ok {
					state.ZipDeployFile = deployFile
				}

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				state.Identity = pointer.From(flattenedIdentity)

				if err := metadata.Encode(&state); err != nil {
					return fmt.Errorf("encoding: %+v", err)
				}
			}

			return nil
		},
	}
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsWebAppSlotV0toV1{},
		},
	}
}
