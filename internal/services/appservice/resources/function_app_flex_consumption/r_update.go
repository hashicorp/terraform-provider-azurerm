// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_flex_consumption

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageDomainSuffix, ok := metadata.Client.Account.Environment.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine Storage domain suffix for environment %q", metadata.Client.Account.Environment.Name)
			}

			client := metadata.Client.AppService.WebAppsClient
			id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state FunctionAppFlexConsumptionModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			model := *existing.Model

			if metadata.ResourceData.HasChange("enabled") {
				model.Properties.Enabled = pointer.To(state.Enabled)
			}

			if metadata.ResourceData.HasChange("virtual_network_subnet_id") {
				subnetId := metadata.ResourceData.Get("virtual_network_subnet_id").(string)
				if subnetId == "" {
					if _, err := client.DeleteSwiftVirtualNetwork(ctx, *id); err != nil {
						return fmt.Errorf("removing `virtual_network_subnet_id` association for %s: %+v", *id, err)
					}
					var empty *string
					model.Properties.VirtualNetworkSubnetId = empty
				} else {
					model.Properties.VirtualNetworkSubnetId = pointer.To(subnetId)

					locks.ByID(subnetId)
					defer locks.UnlockByID(subnetId)
				}
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

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = pointer.To(state.Tags)
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving App Settings for %s: %+v", id, err)
			}

			backendSaConStr, backendStorageUseMsi, deploymentSaConStrVal := FlattenStorageConnectionStrings(*appSettingsResp.Model)
			_, deploymentStorageKey := helpers.ParseWebJobsStorageString(deploymentSaConStrVal)
			deploymentStorage := model.Properties.FunctionAppConfig.Deployment.Storage

			if !features.SixPointOh() {
				endpoint, err := url.Parse(state.StorageContainerEndpoint)
				if err != nil {
					return fmt.Errorf("parsing storage container endpoint error, the expected format is https://storagename.blob.core.windows.net/containername, the received value is %s", state.StorageContainerEndpoint)
				}
				deploymentStorageName, _, ok := strings.Cut(endpoint.Hostname(), ".")
				if !ok {
					return fmt.Errorf("extracting storage account name from endpoint `%s`", state.StorageContainerEndpoint)
				}

				if metadata.ResourceData.HasChange("storage_container_endpoint") {
					deploymentStorage.Value = &state.StorageContainerEndpoint
				}

				if metadata.ResourceData.HasChange("storage_container_type") {
					deploymentStorage.Type = pointer.ToEnum[webapps.FunctionsDeploymentStorageType](state.StorageContainerType)
				}

				if metadata.ResourceData.HasChange("storage_access_key") {
					deploymentStorageKey = state.StorageAccessKey
					state.AppSettings[DeploymentStorageConnStr] = fmt.Sprintf(helpers.StorageStringFmt, deploymentStorageName, deploymentStorageKey, *storageDomainSuffix)
				}

				if metadata.ResourceData.HasChange("storage_authentication_type") {
					deploymentStorage.Authentication = &webapps.FunctionsDeploymentStorageAuthentication{}
					storageAuthType := webapps.AuthenticationType(state.StorageAuthType)
					switch storageAuthType {
					case webapps.AuthenticationTypeStorageAccountConnectionString:
						if deploymentStorageKey == "" {
							return fmt.Errorf("the storage account access key must be specified when using the storage key based access")
						}
						backendStorageUseMsi = false
						backendSaConStr = fmt.Sprintf(StorageStringFmt, deploymentStorageName, deploymentStorageKey, *storageDomainSuffix)
						deploymentStorage.Authentication.StorageAccountConnectionStringName = pointer.To(string(DeploymentStorageConnStr))
					case webapps.AuthenticationTypeUserAssignedIdentity:
						if state.StorageUserAssignedIdentityID == "" {
							return fmt.Errorf("the user assigned identity id must be specified when using the user assigned identity to access the storage account")
						}
						backendStorageUseMsi = true
						backendSaConStr = deploymentStorageName
						deploymentSaConStrVal = ""
						deploymentStorage.Authentication.UserAssignedIdentityResourceId = pointer.To(state.StorageUserAssignedIdentityID)
					case webapps.AuthenticationTypeSystemAssignedIdentity:
						backendStorageUseMsi = true
						backendSaConStr = deploymentStorageName
						deploymentSaConStrVal = ""
					}
					deploymentStorage.Authentication.Type = &storageAuthType
				}

				if metadata.ResourceData.HasChange("storage_user_assigned_identity_id") {
					deploymentStorage.Authentication.UserAssignedIdentityResourceId = pointer.To(state.StorageUserAssignedIdentityID)
				}
			}

			if metadata.ResourceData.HasChange("backend_storage") {
				backendSaConStr, backendStorageUseMsi = ExpandBackendStorage(state.BackendStorage, storageDomainSuffix)
			}

			if metadata.ResourceData.HasChange("deployment_storage") {
				deploymentStorage, deploymentSaConStrVal = ExpandDeploymentStorage(state.DeploymentStorage, DeploymentStorageConnStr, storageDomainSuffix)
				model.Properties.FunctionAppConfig.Deployment.Storage = deploymentStorage
				if deploymentSaConStrVal != "" {
					state.AppSettings[DeploymentStorageConnStr] = deploymentSaConStrVal
				}
			}

			if metadata.ResourceData.HasChange("key_vault_reference_identity_id") {
				model.Properties.KeyVaultReferenceIdentity = pointer.To(state.KeyVaultReferenceIdentityId)
			}

			// Note: We process this regardless to give us a "clean" view of service-side app_settings, so we can reconcile the user-defined entries later
			siteConfig, err := helpers.ExpandSiteConfigFunctionFlexConsumptionApp(state.SiteConfig, model.Properties.SiteConfig, metadata, backendStorageUseMsi, backendSaConStr, deploymentSaConStrVal)
			if err != nil {
				return fmt.Errorf("expanding Site Config for %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("site_config") {
				model.Properties.SiteConfig = siteConfig
			}

			if metadata.ResourceData.HasChange("maximum_instance_count") {
				model.Properties.FunctionAppConfig.ScaleAndConcurrency.MaximumInstanceCount = pointer.To(state.MaximumInstanceCount)
			}

			if metadata.ResourceData.HasChange("instance_memory_in_mb") {
				model.Properties.FunctionAppConfig.ScaleAndConcurrency.InstanceMemoryMB = pointer.To(state.InstanceMemoryInMB)
			}

			if metadata.ResourceData.HasChange("always_ready") {
				arc, err := ExpandAlwaysReadyConfiguration(state.AlwaysReady, state.MaximumInstanceCount)
				if err != nil {
					return fmt.Errorf("expanding `always_ready` for %s: %+v", id, err)
				}
				model.Properties.FunctionAppConfig.ScaleAndConcurrency.AlwaysReady = arc
			}

			if metadata.ResourceData.HasChange("maximum_instance_count") {
				model.Properties.FunctionAppConfig.ScaleAndConcurrency.MaximumInstanceCount = &state.MaximumInstanceCount
			}

			if metadata.ResourceData.HasChange("http_concurrency") {
				if state.HttpConcurrency > 0 {
					model.Properties.FunctionAppConfig.ScaleAndConcurrency.Triggers = &webapps.FunctionsScaleAndConcurrencyTriggers{
						HTTP: &webapps.FunctionsScaleAndConcurrencyTriggersHTTP{
							PerInstanceConcurrency: &state.HttpConcurrency,
						},
					}
				} else {
					model.Properties.FunctionAppConfig.ScaleAndConcurrency.Triggers = &webapps.FunctionsScaleAndConcurrencyTriggers{}
				}
			}

			if metadata.ResourceData.HasChange("runtime_name") {
				runtimeName := webapps.RuntimeName(state.RuntimeName)
				model.Properties.FunctionAppConfig.Runtime.Name = pointer.To(runtimeName)
			}

			if metadata.ResourceData.HasChange("runtime_version") {
				model.Properties.FunctionAppConfig.Runtime.Version = pointer.To(state.RuntimeVersion)
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

			if metadata.ResourceData.HasChange("https_only") {
				model.Properties.HTTPSOnly = pointer.To(state.HttpsOnly)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("webdeploy_publish_basic_authentication_enabled") {
				sitePolicy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: state.PublishingDeployBasicAuthEnabled,
					},
				}
				if _, err := client.UpdateScmAllowed(ctx, *id, sitePolicy); err != nil {
					return fmt.Errorf("setting basic auth for deploy publishing credentials for %s: %+v", id, err)
				}
			}

			if _, err := client.UpdateConfiguration(ctx, *id, webapps.SiteConfigResource{Properties: model.Properties.SiteConfig}); err != nil {
				return fmt.Errorf("updating Site Config for %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				connectionStringUpdate := helpers.ExpandConnectionStrings(state.ConnectionStrings)
				if connectionStringUpdate.Properties == nil {
					connectionStringUpdate.Properties = pointer.To(map[string]webapps.ConnStringValueTypePair{})
				}
				if _, err := client.UpdateConnectionStrings(ctx, *id, *connectionStringUpdate); err != nil {
					return fmt.Errorf("updating Connection Strings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("sticky_settings") {
				emptySlice := make([]string, 0)
				stickySettings := helpers.ExpandStickySettings(state.StickySettings)
				stickySettingsUpdate := webapps.SlotConfigNamesResource{
					Properties: &webapps.SlotConfigNames{
						AppSettingNames:       &emptySlice,
						ConnectionStringNames: &emptySlice,
					},
				}

				if stickySettings != nil {
					if stickySettings.AppSettingNames != nil {
						stickySettingsUpdate.Properties.AppSettingNames = stickySettings.AppSettingNames
					}
					if stickySettings.ConnectionStringNames != nil {
						stickySettingsUpdate.Properties.ConnectionStringNames = stickySettings.ConnectionStringNames
					}
				}

				if _, err := client.UpdateSlotConfigurationNames(ctx, *id, stickySettingsUpdate); err != nil {
					return fmt.Errorf("updating Sticky Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings") {
				authUpdate := helpers.ExpandAuthSettings(state.AuthSettings)
				// (@jackofallops) - in the case of a removal of this block, we need to zero these settings
				if authUpdate.Properties == nil {
					authUpdate.Properties = helpers.DefaultAuthSettingsProperties()
				}
				if _, err := client.UpdateAuthSettings(ctx, *id, *authUpdate); err != nil {
					return fmt.Errorf("updating Auth Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("auth_settings_v2") {
				authV2Update := helpers.ExpandAuthV2Settings(state.AuthV2Settings)
				// (@toddgiguere) - in the case of a removal of this block, we need to zero these settings
				if authV2Update.Properties == nil {
					authV2Update.Properties = helpers.DefaultAuthV2SettingsProperties()
				}
				if _, err := client.UpdateAuthSettingsV2(ctx, *id, *authV2Update); err != nil {
					return fmt.Errorf("updating AuthV2 Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("site_config.0.app_service_logs") {
				appServiceLogs := helpers.ExpandFunctionAppAppServiceLogs(state.SiteConfig[0].AppServiceLogs)
				if _, err := client.UpdateDiagnosticLogsConfig(ctx, *id, appServiceLogs); err != nil {
					return fmt.Errorf("updating App Service Log Settings for %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("zip_deploy_file") {
				if err = helpers.GetCredentialsAndPublish(ctx, client, *id, state.ZipDeployFile); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
