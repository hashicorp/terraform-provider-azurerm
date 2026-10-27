// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_standard

import (
	"context"
	"fmt"
	"log"
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

			var state LogicAppResourceModel

			id, err := commonids.ParseLogicAppId(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading Linux %s: %+v", id, err)
			}

			state.Name = id.SiteName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Kind = pointer.From(model.Kind)
				state.Location = location.Normalize(model.Location)
				ident, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				state.Identity = pointer.From(ident)
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
					if err != nil {
						return err
					}

					state.AppServicePlanId = servicePlanId.ID()
					state.Enabled = pointer.From(props.Enabled)
					state.DefaultHostname = pointer.From(props.DefaultHostName)
					state.HTTPSOnly = pointer.From(props.HTTPSOnly)
					state.OutboundIpAddresses = pointer.From(props.OutboundIPAddresses)
					state.PossibleOutboundIpAddresses = pointer.From(props.PossibleOutboundIPAddresses)
					state.ClientAffinityEnabled = pointer.From(props.ClientAffinityEnabled)
					state.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					state.VirtualNetworkSubnetId = pointer.From(props.VirtualNetworkSubnetId)
					state.VNETContentShareEnabled = pointer.From(props.VnetContentShareEnabled)
					state.PublicNetworkAccess = pointer.From(props.PublicNetworkAccess)
					if kvRefId := pointer.From(props.KeyVaultReferenceIdentity); !strings.EqualFold(kvRefId, "SystemAssigned") {
						state.KeyvaultReferenceIdentityId = kvRefId
					}
					state.ClientCertificateMode = pointer.FromEnum(props.ClientCertMode)
				}
			}

			appSettingsResp, err := client.ListApplicationSettings(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing application settings for %s: %+v", *id, err)
			}

			if model := appSettingsResp.Model; model != nil {
				appSettings := pointer.From(model.Properties)

				connectionString := appSettings[storageAppSettingName]

				if strings.HasPrefix(connectionString, "@Microsoft.KeyVault") {
					state.StorageKeyVaultSecretID = strings.TrimPrefix(strings.TrimSuffix(connectionString, ")"), "@Microsoft.KeyVault(SecretUri=")
				} else {
					state.StorageAccountName, state.StorageAccountAccessKey = helpers.ParseWebJobsStorageString(connectionString)
				}

				if v, ok := appSettings[functionVersionAppSettingName]; ok {
					state.Version = v
				}

				if _, ok := appSettings["AzureFunctionsJobHost__extensionBundle__id"]; ok {
					state.UseExtensionBundle = true

					if val, ok := appSettings["AzureFunctionsJobHost__extensionBundle__version"]; ok {
						state.BundleVersion = val
					}
				} else {
					state.UseExtensionBundle = false
					state.BundleVersion = "[1.*, 2.0.0)"
				}

				state.StorageAccountShareName = appSettings[contentShareAppSettingName]
				delete(appSettings, contentFileConnStringAppSettingName)
				delete(appSettings, "APP_KIND")
				delete(appSettings, "AzureFunctionsJobHost__extensionBundle__id")
				delete(appSettings, "AzureFunctionsJobHost__extensionBundle__version")
				delete(appSettings, "AzureWebJobsDashboard")
				delete(appSettings, storageAppSettingName)
				delete(appSettings, functionVersionAppSettingName)
				delete(appSettings, contentShareAppSettingName)

				state.AppSettings = appSettings
			}

			connectionStringsResp, err := client.ListConnectionStrings(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing connection strings for %s: %+v", *id, err)
			}

			if model := connectionStringsResp.Model; model != nil {
				state.ConnectionStrings = helpers.FlattenConnectionStrings(model)
			}

			ftpBasicAuth, err := client.GetFtpAllowed(ctx, *id)
			if err != nil || ftpBasicAuth.Model == nil {
				return fmt.Errorf("retrieving FTP publish basic authentication policy for %s: %+v", id, err)
			}

			if props := ftpBasicAuth.Model.Properties; props != nil {
				state.FtpPublishBasicAuthEnabled = props.Allow
			}

			scmBasicAuth, err := client.GetScmAllowed(ctx, *id)
			if err != nil || scmBasicAuth.Model == nil {
				return fmt.Errorf("retrieving SCM publish basic authentication policy for %s: %+v", id, err)
			}

			if props := scmBasicAuth.Model.Properties; props != nil {
				state.SCMPublishBasicAuthEnabled = props.Allow
			}

			siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, *id)
			if err != nil {
				return fmt.Errorf("listing publishing credentials for %s: %+v", *id, err)
			}

			state.SiteCredential = helpers.FlattenSiteCredentialsLogicApp(siteCredentials)

			configResp, err := client.GetConfiguration(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving the configuration for %s: %+v", *id, err)
			}

			if model := configResp.Model; model != nil {
				state.SiteConfig = flattenLogicAppStandardSiteConfig(model.Properties)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenLogicAppStandardSiteConfig(input *webapps.SiteConfig) []helpers.LogicAppSiteConfig {
	results := make([]helpers.LogicAppSiteConfig, 0)
	result := helpers.LogicAppSiteConfig{}

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	result.AlwaysOn = pointer.From(input.AlwaysOn)
	result.Use32BitWorkerProcess = pointer.From(input.Use32BitWorkerProcess)
	result.WebSocketsEnabled = pointer.From(input.WebSocketsEnabled)
	result.LinuxFxVersion = pointer.From(input.LinuxFxVersion)
	result.HTTP2Enabled = pointer.From(input.HTTP20Enabled)
	result.PreWarmedInstanceCount = pointer.From(input.PreWarmedInstanceCount)
	result.IpRestriction = helpers.FlattenIpRestrictions(input.IPSecurityRestrictions)
	result.SCMIPRestriction = helpers.FlattenIpRestrictions(input.ScmIPSecurityRestrictions)
	result.SCMIpRestrictionDefaultAction = pointer.FromEnum(input.ScmIPSecurityRestrictionsDefaultAction)

	result.SCMUseMainIpRestriction = pointer.From(input.ScmIPSecurityRestrictionsUseMain)

	result.SCMType = pointer.FromEnum(input.ScmType)
	result.SCMMinTLSVersion = pointer.FromEnum(input.ScmMinTlsVersion)

	result.MinTLSVersion = pointer.FromEnum(input.MinTlsVersion)
	result.FTPSState = pointer.FromEnum(input.FtpsState)

	result.Cors = helpers.FlattenCorsSettings(input.Cors)

	result.AutoSwapSlotName = pointer.From(input.AutoSwapSlotName)

	result.HealthCheckPath = pointer.From(input.HealthCheckPath)

	result.ElasticInstanceMinimum = pointer.From(input.MinimumElasticInstanceCount)

	result.AppScaleLimit = pointer.From(input.FunctionAppScaleLimit)

	result.RuntimeScaleMonitoringEnabled = pointer.From(input.FunctionsRuntimeScaleMonitoringEnabled)

	result.DotnetFrameworkVersion = pointer.From(input.NetFrameworkVersion)

	result.VNETRouteAllEnabled = pointer.From(input.VnetRouteAllEnabled)

	result.IpRestrictionDefaultAction = pointer.FromEnum(input.IPSecurityRestrictionsDefaultAction)

	results = append(results, result)
	return results
}

func flattenLogicAppStandardIpRestriction(input *[]webapps.IPSecurityRestriction) []any {
	restrictions := make([]any, 0)

	if input == nil {
		return []any{}
	}

	for _, v := range *input {
		restriction := make(map[string]any)
		if ip := v.IPAddress; ip != nil {
			if *ip == "Any" {
				continue
			} else {
				switch pointer.From(v.Tag) {
				case webapps.IPFilterTagServiceTag:
					restriction["service_tag"] = *ip
				default:
					restriction["ip_address"] = *ip
				}
			}
		}

		restriction["virtual_network_subnet_id"] = pointer.From(v.VnetSubnetResourceId)

		restriction["name"] = pointer.From(v.Name)

		priority := 0
		if priorityRaw := v.Priority; priorityRaw != nil {
			priority = int(*priorityRaw)
		}
		restriction["priority"] = priority

		restriction["action"] = pointer.From(v.Action)

		if headers := v.Headers; headers != nil {
			restriction["headers"] = flattenHeaders(*headers)
		}

		restrictions = append(restrictions, restriction)
	}

	return restrictions
}

var (
	logicAppStdKind   = "functionapp,workflowapp"
	logicAppLinuxKind = "functionapp,linux,container,workflowapp"

	storageAppSettingName                = "AzureWebJobsStorage"
	contentShareAppSettingName           = "WEBSITE_CONTENTSHARE"
	contentFileConnStringAppSettingName  = "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"
	functionVersionAppSettingName        = "FUNCTIONS_EXTENSION_VERSION"
	extensionBundleAppSettingName        = "AzureFunctionsJobHost__extensionBundle__id"
	extensionBundleAppSettingValue       = "Microsoft.Azure.Functions.ExtensionBundle.Workflows"
	extensionBundleVersionAppSettingName = "AzureFunctionsJobHost__extensionBundle__version"
)

func getBasicLogicAppSettings(d LogicAppResourceModel, endpointSuffix string) ([]webapps.NameValuePair, error) {
	var storageConnection string
	if d.StorageKeyVaultSecretID != "" {
		storageConnection = fmt.Sprintf(helpers.StorageStringFmtKV, d.StorageKeyVaultSecretID)
	} else {
		storageConnection = fmt.Sprintf(helpers.StorageStringFmt, d.StorageAccountName, d.StorageAccountAccessKey, endpointSuffix)
	}

	contentShare := strings.ToLower(d.Name) + "-content"
	if d.StorageAccountShareName != "" {
		contentShare = d.StorageAccountShareName
	}

	basicSettings := []webapps.NameValuePair{
		{Name: &storageAppSettingName, Value: &storageConnection},
		{Name: &functionVersionAppSettingName, Value: pointer.To(d.Version)},
		{Name: pointer.To("APP_KIND"), Value: pointer.To("workflowApp")},
		{Name: &contentShareAppSettingName, Value: &contentShare},
		{Name: &contentFileConnStringAppSettingName, Value: &storageConnection},
	}

	if d.UseExtensionBundle {
		extensionBundleVersion := d.BundleVersion

		if extensionBundleVersion == "" {
			return nil, fmt.Errorf(
				"when `use_extension_bundle` is true, `bundle_version` must be specified",
			)
		}

		bundleSettings := []webapps.NameValuePair{
			{Name: pointer.To("AzureFunctionsJobHost__extensionBundle__id"), Value: pointer.To("Microsoft.Azure.Functions.ExtensionBundle.Workflows")},
			{Name: pointer.To("AzureFunctionsJobHost__extensionBundle__version"), Value: &extensionBundleVersion},
		}

		return append(basicSettings, bundleSettings...), nil
	}

	return basicSettings, nil
}
