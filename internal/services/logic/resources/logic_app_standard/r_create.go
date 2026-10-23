// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_standard

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
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
			aseClient := metadata.Client.AppService.AppServiceEnvironmentClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			env := metadata.Client.Account.Environment
			storageAccountDomainSuffix, ok := env.Storage.DomainSuffix()
			if !ok {
				return fmt.Errorf("could not determine the domain suffix for storage accounts in environment %q: %+v", env.Name, env.Storage)
			}

			data := LogicAppResourceModel{}

			if err := metadata.Decode(&data); err != nil {
				return err
			}

			id := commonids.NewAppServiceID(subscriptionId, data.ResourceGroupName, data.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_logic_app_standard", id.ID())
				}
			}

			servicePlanId, err := commonids.ParseAppServicePlanID(data.AppServicePlanId)
			if err != nil {
				return err
			}

			servicePlan, err := servicePlanClient.Get(ctx, *servicePlanId)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", servicePlanId, err)
			}

			availabilityRequest := resourceproviders.ResourceNameAvailabilityRequest{
				Name: data.Name,
				Type: resourceproviders.CheckNameResourceTypesMicrosoftPointWebSites,
			}
			if servicePlanModel := servicePlan.Model; servicePlanModel != nil {
				if ase := servicePlanModel.Properties.HostingEnvironmentProfile; ase != nil {
					// Attempt to check the ASE for the appropriate suffix for the name availability request.
					// This varies between internal and external ASE Types, and potentially has other names in other clouds
					// We use the "internal" as the fallback here, if we can read the ASE, we'll get the full one
					nameSuffix := "appserviceenvironment.net"
					if ase.Id != nil {
						aseId, err := commonids.ParseAppServiceEnvironmentIDInsensitively(*ase.Id)
						nameSuffix = fmt.Sprintf("%s.%s", aseId.HostingEnvironmentName, nameSuffix)
						if err != nil {
							metadata.Logger.Warnf("could not parse App Service Environment ID determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", data.Name, servicePlanId)
						} else {
							existingASE, err := aseClient.Get(ctx, *aseId)
							if err != nil || existingASE.Model == nil {
								metadata.Logger.Warnf("could not read App Service Environment to determine FQDN for name availability check, defaulting to `%s.%s.appserviceenvironment.net`", data.Name, servicePlanId)
							} else if props := existingASE.Model.Properties; props != nil && props.DnsSuffix != nil && *props.DnsSuffix != "" {
								nameSuffix = *props.DnsSuffix
							}
						}
					}

					availabilityRequest.Name = fmt.Sprintf("%s.%s", data.Name, nameSuffix)
					availabilityRequest.IsFqdn = pointer.To(true)
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

			basicAppSettings, err := getBasicLogicAppSettings(data, *storageAccountDomainSuffix)
			if err != nil {
				return err
			}

			siteConfig, err := expandLogicAppStandardSiteConfigForCreate(data.SiteConfig, metadata)
			if err != nil {
				return fmt.Errorf("expanding `site_config`: %+v", err)
			}

			kind := logicAppStdKind
			if siteConfig.LinuxFxVersion != nil && len(*siteConfig.LinuxFxVersion) > 0 {
				kind = logicAppLinuxKind
			}

			appSettings := expandAppSettings(data.AppSettings)
			appSettings = append(appSettings, basicAppSettings...)

			siteConfig.AppSettings = pointer.To(appSettings)

			if v, ok := data.AppSettings["WEBSITE_VNET_ROUTE_ALL"]; ok {
				// For compatibility between app_settings and site_config, we need to set the API property based on the presence of the app_setting map value if present.
				// a replacement of this resource should consider deprecating support for this.
				vnetRouteAll, _ := strconv.ParseBool(v)
				siteConfig.VnetRouteAllEnabled = pointer.To(vnetRouteAll)
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(data.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			siteEnvelope := webapps.Site{
				Identity: expandedIdentity,
				Kind:     pointer.To(kind),
				Location: location.Normalize(data.Location),
				Properties: &webapps.SiteProperties{
					ServerFarmId:            pointer.To(data.AppServicePlanId),
					Enabled:                 pointer.To(data.Enabled),
					ClientAffinityEnabled:   pointer.To(data.ClientAffinityEnabled),
					ClientCertEnabled:       pointer.To(data.ClientCertificateMode != ""),
					HTTPSOnly:               pointer.To(data.HTTPSOnly),
					SiteConfig:              siteConfig,
					VnetContentShareEnabled: pointer.To(data.VNETContentShareEnabled),
					PublicNetworkAccess:     pointer.To(data.PublicNetworkAccess),
				},
				Tags: pointer.To(data.Tags),
			}

			siteEnvelope.Properties.ClientCertMode = pointer.ToEnum[webapps.ClientCertMode](data.ClientCertificateMode)

			if data.VirtualNetworkSubnetId != "" {
				siteEnvelope.Properties.VirtualNetworkSubnetId = pointer.To(data.VirtualNetworkSubnetId)
			}

			if data.KeyvaultReferenceIdentityId != "" {
				siteEnvelope.Properties.KeyVaultReferenceIdentity = pointer.To(data.KeyvaultReferenceIdentityId)
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, siteEnvelope, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			if !data.FtpPublishBasicAuthEnabled {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.FtpPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateFtpAllowed(ctx, id, policy); err != nil {
					return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
				}
			}

			if !data.SCMPublishBasicAuthEnabled {
				policy := webapps.CsmPublishingCredentialsPoliciesEntity{
					Properties: &webapps.CsmPublishingCredentialsPoliciesEntityProperties{
						Allow: data.SCMPublishBasicAuthEnabled,
					},
				}

				if _, err := client.UpdateScmAllowed(ctx, id, policy); err != nil {
					return fmt.Errorf("updating FTP publish basic authentication policy for %s: %+v", id, err)
				}
			}

			connectionStrings := helpers.ExpandConnectionStrings(data.ConnectionStrings)
			if connectionStrings.Properties != nil {
				if _, err := client.UpdateConnectionStrings(ctx, id, *connectionStrings); err != nil {
					return fmt.Errorf("setting Connection Strings for Linux %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func expandLogicAppStandardSiteConfigForCreate(d []helpers.LogicAppSiteConfig, metadata sdk.ResourceMetaData) (*webapps.SiteConfig, error) {
	siteConfig := &webapps.SiteConfig{}
	if len(d) == 0 {
		return siteConfig, nil
	}

	config := d[0]

	siteConfig.AlwaysOn = pointer.To(config.AlwaysOn)
	siteConfig.HTTP20Enabled = pointer.To(config.HTTP2Enabled)
	siteConfig.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(config.RuntimeScaleMonitoringEnabled)
	siteConfig.Use32BitWorkerProcess = pointer.To(config.Use32BitWorkerProcess)
	siteConfig.WebSocketsEnabled = pointer.To(config.WebSocketsEnabled)
	siteConfig.ScmIPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](config.SCMIpRestrictionDefaultAction)

	if config.LinuxFxVersion != "" {
		siteConfig.LinuxFxVersion = pointer.To(config.LinuxFxVersion)
	}

	if len(config.Cors) > 0 {
		siteConfig.Cors = helpers.ExpandCorsSettings(config.Cors)
	}

	ipr, err := helpers.ExpandIpRestrictions(config.IpRestriction)
	if err != nil {
		return nil, err
	}
	siteConfig.IPSecurityRestrictions = ipr

	ipr, err = helpers.ExpandIpRestrictions(config.SCMIPRestriction)
	if err != nil {
		return nil, err
	}
	siteConfig.ScmIPSecurityRestrictions = ipr

	siteConfig.ScmIPSecurityRestrictionsUseMain = pointer.To(config.SCMUseMainIpRestriction)

	siteConfig.ScmMinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.SCMMinTLSVersion)
	siteConfig.ScmType = pointer.ToEnum[webapps.ScmType](config.SCMType)
	siteConfig.MinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.MinTLSVersion)

	siteConfig.FtpsState = pointer.ToEnum[webapps.FtpsState](config.FTPSState)

	if metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["pre_warmed_instance_count"].IsKnown() {
		siteConfig.PreWarmedInstanceCount = pointer.To(config.PreWarmedInstanceCount)
	}
	siteConfig.HealthCheckPath = pointer.To(config.HealthCheckPath)

	if metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["elastic_instance_minimum"].IsKnown() && config.ElasticInstanceMinimum > 0 {
		siteConfig.MinimumElasticInstanceCount = pointer.To(config.ElasticInstanceMinimum)
	}

	if metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["app_scale_limit"].IsKnown() {
		siteConfig.FunctionAppScaleLimit = pointer.To(config.AppScaleLimit)
	}

	siteConfig.NetFrameworkVersion = pointer.To(config.DotnetFrameworkVersion)
	siteConfig.VnetRouteAllEnabled = pointer.To(config.VNETRouteAllEnabled)

	siteConfig.PublicNetworkAccess = pointer.To(metadata.ResourceData.Get("public_network_access").(string))

	siteConfig.IPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](config.IpRestrictionDefaultAction)

	return siteConfig, nil
}

func expandLogicAppStandardSiteConfigForUpdate(d []helpers.LogicAppSiteConfig, metadata sdk.ResourceMetaData, existing *webapps.SiteConfig) (*webapps.SiteConfig, error) {
	siteConfig := &webapps.SiteConfig{}
	if len(d) == 0 {
		siteConfig.Cors = &webapps.CorsSettings{
			AllowedOrigins:     pointer.To(make([]string, 0)),
			SupportCredentials: pointer.To(false),
		}
		return siteConfig, nil
	}

	if existing != nil {
		siteConfig = existing
	}

	config := d[0]

	siteConfig.AlwaysOn = pointer.To(config.AlwaysOn)
	siteConfig.HTTP20Enabled = pointer.To(config.HTTP2Enabled)
	siteConfig.FunctionsRuntimeScaleMonitoringEnabled = pointer.To(config.RuntimeScaleMonitoringEnabled)
	siteConfig.Use32BitWorkerProcess = pointer.To(config.Use32BitWorkerProcess)
	siteConfig.ScmIPSecurityRestrictionsUseMain = pointer.To(config.SCMUseMainIpRestriction)
	siteConfig.WebSocketsEnabled = pointer.To(config.WebSocketsEnabled)
	siteConfig.VnetRouteAllEnabled = pointer.To(config.VNETRouteAllEnabled)

	if metadata.ResourceData.HasChange("site_config.0.linux_fx_version") {
		siteConfig.LinuxFxVersion = pointer.To(config.LinuxFxVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		siteConfig.Cors = helpers.ExpandCorsSettings(config.Cors)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipr, err := helpers.ExpandIpRestrictions(config.IpRestriction)
		if err != nil {
			return nil, err
		}
		siteConfig.IPSecurityRestrictions = ipr
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		ipr, err := helpers.ExpandIpRestrictions(config.SCMIPRestriction)
		if err != nil {
			return nil, err
		}
		siteConfig.ScmIPSecurityRestrictions = ipr
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction_default_action") {
		siteConfig.ScmIPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](config.SCMIpRestrictionDefaultAction)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_min_tls_version") {
		siteConfig.ScmMinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.SCMMinTLSVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_type") {
		siteConfig.ScmType = pointer.ToEnum[webapps.ScmType](config.SCMType)
	}

	if metadata.ResourceData.HasChange("site_config.0.min_tls_version") {
		siteConfig.MinTlsVersion = pointer.ToEnum[webapps.SupportedTlsVersions](config.MinTLSVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		siteConfig.FtpsState = pointer.ToEnum[webapps.FtpsState](config.FTPSState)
	}

	if len(metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()) > 0 && metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["pre_warmed_instance_count"].IsKnown() {
		siteConfig.PreWarmedInstanceCount = pointer.To(config.PreWarmedInstanceCount)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		siteConfig.HealthCheckPath = pointer.To(config.HealthCheckPath)
	}

	if len(metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()) > 0 && metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["elastic_instance_minimum"].IsKnown() {
		siteConfig.MinimumElasticInstanceCount = pointer.To(config.ElasticInstanceMinimum)
	}

	if len(metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()) > 0 && metadata.ResourceData.GetRawConfig().AsValueMap()["site_config"].AsValueSlice()[0].AsValueMap()["app_scale_limit"].IsKnown() {
		siteConfig.FunctionAppScaleLimit = pointer.To(config.AppScaleLimit)
	}

	if metadata.ResourceData.HasChange("site_config.0.dotnet_framework_version") {
		siteConfig.NetFrameworkVersion = pointer.To(config.DotnetFrameworkVersion)
	}

	// lintignore:R019 // deliberate subset: only the fields that feed the app_settings merge; the remaining attributes do not affect it
	if metadata.ResourceData.HasChanges(
		"app_settings",
		"storage_account_name",
		"storage_account_share_name",
		"storage_account_access_key",
		"version",
		"storage_key_vault_secret_id",
	) {
		o, n := metadata.ResourceData.GetChange("app_settings")

		appSettings := make([]webapps.NameValuePair, 0)
		if existing != nil {
			appSettings = *existing.AppSettings
		}

		siteConfig.AppSettings = mergeAppSettings(appSettings, o.(map[string]any), n.(map[string]any), metadata)
	}

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction_default_action") {
		siteConfig.IPSecurityRestrictionsDefaultAction = pointer.ToEnum[webapps.DefaultAction](config.IpRestrictionDefaultAction)
	}

	return siteConfig, nil
}

func expandAppSettings(input map[string]string) []webapps.NameValuePair {
	output := make([]webapps.NameValuePair, 0)

	for k, v := range input {
		nameValue := webapps.NameValuePair{
			Name:  pointer.To(k),
			Value: pointer.To(v),
		}
		output = append(output, nameValue)
	}

	return output
}

func mergeAppSettings(existing []webapps.NameValuePair, old, new map[string]any, metadata sdk.ResourceMetaData) *[]webapps.NameValuePair {
	f := func(input map[string]any) (result map[string]string) {
		result = make(map[string]string)
		for k, v := range input {
			result[k] = v.(string)
		}

		return
	}

	eMap := make(map[string]string)
	for _, i := range existing {
		n, v := pointer.From(i.Name), pointer.From(i.Value)
		eMap[n] = v
	}

	oMap := f(old)
	cMap := f(new)

	if metadata.ResourceData.HasChange("storage_key_vault_secret_id") && metadata.ResourceData.Get("storage_key_vault_secret_id").(string) != "" {
		kvRef := fmt.Sprintf(helpers.StorageStringFmtKV, metadata.ResourceData.Get("storage_key_vault_secret_id").(string))
		eMap[storageAppSettingName] = kvRef
		eMap[contentFileConnStringAppSettingName] = kvRef
	} else if metadata.ResourceData.HasChanges("storage_account_name", "storage_account_access_key") {
		accountName := metadata.ResourceData.Get("storage_account_name").(string)
		accountAccessKey := metadata.ResourceData.Get("storage_account_access_key").(string)
		suffix, _ := metadata.Client.Account.Environment.Storage.DomainSuffix()

		conn := fmt.Sprintf(helpers.StorageStringFmt, accountName, accountAccessKey, *suffix)
		eMap[storageAppSettingName] = conn
		eMap[contentFileConnStringAppSettingName] = conn
	}

	if metadata.ResourceData.HasChange("storage_account_share_name") {
		n := metadata.ResourceData.Get("storage_account_share_name").(string)

		if n != "" {
			eMap[contentShareAppSettingName] = n
		} else {
			name := metadata.ResourceData.Get("name").(string)
			eMap[contentShareAppSettingName] = strings.ToLower(name) + "-content"
		}
	}

	if metadata.ResourceData.HasChange("version") {
		eMap[functionVersionAppSettingName] = metadata.ResourceData.Get("version").(string)
	}

	if metadata.ResourceData.HasChanges("bundle_version", "use_extension_bundle") {
		if bundleVersion := metadata.ResourceData.Get("bundle_version").(string); bundleVersion != "" {
			eMap[extensionBundleAppSettingName] = extensionBundleAppSettingValue
			eMap[extensionBundleVersionAppSettingName] = bundleVersion
		} else {
			delete(eMap, extensionBundleAppSettingName)
		}
	}

	remove := map[string]string{}
	addOrUpdate := map[string]string{}

	for k, v := range oMap {
		if _, ok := cMap[k]; !ok {
			remove[k] = v
			break
		}

		addOrUpdate[k] = v
	}

	maps.Copy(addOrUpdate, cMap)

	for k := range remove {
		delete(eMap, k)
	}

	maps.Copy(eMap, addOrUpdate)

	return pointer.To(expandAppSettings(eMap))
}
