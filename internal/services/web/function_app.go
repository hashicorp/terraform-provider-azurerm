// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func schemaAppServiceFunctionAppSiteConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"cors": SchemaWebCorsSettings(),

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.FtpsStateAllAllowed),
						string(web.FtpsStateDisabled),
						string(web.FtpsStateFtpsOnly),
					}, false),
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": schemaAppServiceIpRestriction(),

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},

				"min_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"pre_warmed_instance_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
				},

				"scm_ip_restriction": schemaAppServiceIpRestriction(),

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ScmTypeBitbucketGit),
						string(web.ScmTypeBitbucketHg),
						string(web.ScmTypeCodePlexGit),
						string(web.ScmTypeCodePlexHg),
						string(web.ScmTypeDropbox),
						string(web.ScmTypeExternalGit),
						string(web.ScmTypeExternalHg),
						string(web.ScmTypeGitHub),
						string(web.ScmTypeLocalGit),
						string(web.ScmTypeNone),
						string(web.ScmTypeOneDrive),
						string(web.ScmTypeTfs),
						string(web.ScmTypeVSO),
						string(web.ScmTypeVSTSRM),
					}, false),
				},

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"use_32_bit_worker_process": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				// The following is only used for "slots"
				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"java_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice([]string{"1.8", "11", "17"}, false),
				},

				"elastic_instance_minimum": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					Computed: true,
					ValidateFunc: func() pluginsdk.SchemaValidateFunc {
						// 0 is no longer a valid value
						if features.FourPointOhBeta() {
							return validation.IntBetween(1, 20)
						}
						return validation.IntBetween(0, 20)
					}(),
				},

				"app_scale_limit": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"runtime_scale_monitoring_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"dotnet_framework_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "v4.0",
					ValidateFunc: validation.StringInSlice([]string{
						"v4.0",
						"v5.0",
						"v6.0",
					}, false),
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func schemaFunctionAppDataSourceSiteConfig() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"cors": SchemaWebCorsSettings(),

				"use_32_bit_worker_process": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": schemaAppServiceDataSourceIpRestriction(),

				"min_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"pre_warmed_instance_count": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				// The following is only used for "slots"
				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_ip_restriction": schemaAppServiceDataSourceIpRestriction(),

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"elastic_instance_minimum": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"app_scale_limit": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"runtime_scale_monitoring_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"dotnet_framework_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func getBasicFunctionAppAppSettings(d *pluginsdk.ResourceData, appServiceTier, endpointSuffix string, existingSettings map[string]*string) ([]web.NameValuePair, error) {
	// TODO: This is a workaround since there are no public Functions API
	// You may track the API request here: https://github.com/Azure/azure-rest-api-specs/issues/3750
	dashboardPropName := "AzureWebJobsDashboard"
	storagePropName := "AzureWebJobsStorage"
	functionVersionPropName := "FUNCTIONS_EXTENSION_VERSION"

	contentSharePropName := "WEBSITE_CONTENTSHARE"
	contentFileConnStringPropName := "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING"

	var storageConnection string

	storageAccountName := ""
	if v, ok := d.GetOk("storage_account_name"); ok {
		storageAccountName = v.(string)
	}

	storageAccountKey := ""
	if v, ok := d.GetOk("storage_account_access_key"); ok {
		storageAccountKey = v.(string)
	}

	if storageAccountName == "" && storageAccountKey == "" {
		return nil, fmt.Errorf("Both `storage_account_name` and `storage_account_access_key` must be specified")
	}

	if (storageAccountName == "" && storageAccountKey != "") || (storageAccountName != "" && storageAccountKey == "") {
		return nil, fmt.Errorf("both `storage_account_name` and `storage_account_access_key` must be specified")
	}

	if storageAccountKey != "" && storageAccountName != "" {
		storageConnection = fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", storageAccountName, storageAccountKey, endpointSuffix)
	}

	functionVersion := d.Get("version").(string)
	var contentShare string
	contentSharePreviouslySet := false
	if currentContentShare, ok := existingSettings[contentSharePropName]; ok {
		contentShare = *currentContentShare
		contentSharePreviouslySet = true
	} else {
		suffix := uuid.New().String()[0:4]
		contentShare = strings.ToLower(d.Get("name").(string)) + suffix
	}

	basicSettings := []web.NameValuePair{
		{Name: &storagePropName, Value: &storageConnection},
		{Name: &functionVersionPropName, Value: &functionVersion},
	}

	if d.Get("enable_builtin_logging").(bool) {
		basicSettings = append(basicSettings, web.NameValuePair{
			Name:  &dashboardPropName,
			Value: &storageConnection,
		})
	}

	consumptionSettings := []web.NameValuePair{
		{Name: &contentSharePropName, Value: &contentShare},
		{Name: &contentFileConnStringPropName, Value: &storageConnection},
	}

	// If there's an existing value for content, we need to send it. This can be the case for PremiumV2/PremiumV3 plans where the value has been previously configured.
	if contentSharePreviouslySet {
		return append(basicSettings, consumptionSettings...), nil
	}

	// On consumption and premium plans include WEBSITE_CONTENT components, unless it's a Linux consumption plan
	// Note: The docs on this are misleading. Premium here refers explicitly to `ElasticPremium`, and not `PremiumV2` / `PremiumV3` etc.
	if !(strings.EqualFold(appServiceTier, "dynamic") && strings.EqualFold(d.Get("os_type").(string), "linux")) &&
		(strings.EqualFold(appServiceTier, "dynamic") || strings.HasPrefix(strings.ToLower(appServiceTier), "elastic")) {
		return append(basicSettings, consumptionSettings...), nil
	}

	return basicSettings, nil
}

func getFunctionAppServiceTier(ctx context.Context, appServicePlanId string, meta interface{}) (string, error) {
	id, err := parse.AppServicePlanID(appServicePlanId)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to parse App Service Plan ID %q: %+v", appServicePlanId, err)
	}

	log.Printf("[DEBUG] Retrieving App Service Plan %q (Resource Group %q)", id.ServerFarmName, id.ResourceGroup)

	appServicePlansClient := meta.(*clients.Client).Web.AppServicePlansClient
	appServicePlan, err := appServicePlansClient.Get(ctx, id.ResourceGroup, id.ServerFarmName)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Could not retrieve App Service Plan ID %q: %+v", appServicePlanId, err)
	}

	if sku := appServicePlan.Sku; sku != nil {
		if tier := sku.Tier; tier != nil {
			return *tier, nil
		}
	}
	return "", fmt.Errorf("No `sku` block was returned for App Service Plan ID %q", appServicePlanId)
}

func expandFunctionAppAppSettings(d *pluginsdk.ResourceData, basicAppSettings []web.NameValuePair) map[string]*string {
	output := expandAppServiceAppSettings(d)

	for _, p := range basicAppSettings {
		output[*p.Name] = p.Value
	}

	return output
}

func expandFunctionAppSiteConfig(d *pluginsdk.ResourceData) (web.SiteConfig, error) {
	configs := d.Get("site_config").([]interface{})
	siteConfig := web.SiteConfig{}

	if len(configs) == 0 {
		return siteConfig, nil
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["always_on"]; ok {
		siteConfig.AlwaysOn = utils.Bool(v.(bool))
	}

	if v, ok := config["use_32_bit_worker_process"]; ok {
		siteConfig.Use32BitWorkerProcess = utils.Bool(v.(bool))
	}

	if v, ok := config["websockets_enabled"]; ok {
		siteConfig.WebSocketsEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["linux_fx_version"]; ok {
		siteConfig.LinuxFxVersion = utils.String(v.(string))
	}

	if v, ok := config["cors"]; ok {
		expand := ExpandWebCorsSettings(v)
		siteConfig.Cors = &expand
	}

	if v, ok := config["http2_enabled"]; ok {
		siteConfig.HTTP20Enabled = utils.Bool(v.(bool))
	}

	if v, ok := config["ip_restriction"]; ok {
		restrictions, err := expandAppServiceIpRestriction(v)
		if err != nil {
			return siteConfig, err
		}
		siteConfig.IPSecurityRestrictions = &restrictions
	}

	if v, ok := config["scm_use_main_ip_restriction"]; ok {
		siteConfig.ScmIPSecurityRestrictionsUseMain = utils.Bool(v.(bool))
	}

	if v, ok := config["scm_ip_restriction"]; ok {
		scmIPSecurityRestrictions := v.([]interface{})
		scmRestrictions, err := expandAppServiceIpRestriction(scmIPSecurityRestrictions)
		if err != nil {
			return siteConfig, err
		}
		siteConfig.ScmIPSecurityRestrictions = &scmRestrictions
	}

	if v, ok := config["min_tls_version"]; ok {
		siteConfig.MinTLSVersion = web.SupportedTLSVersions(v.(string))
	}

	if v, ok := config["ftps_state"]; ok {
		siteConfig.FtpsState = web.FtpsState(v.(string))
	}

	if v, ok := config["pre_warmed_instance_count"]; ok {
		siteConfig.PreWarmedInstanceCount = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["scm_type"]; ok {
		siteConfig.ScmType = web.ScmType(v.(string))
	}

	// This optional parameter can only present in "slot" resources
	if v, ok := config["auto_swap_slot_name"]; ok {
		siteConfig.AutoSwapSlotName = utils.String(v.(string))
	}

	if v, ok := config["health_check_path"]; ok {
		siteConfig.HealthCheckPath = utils.String(v.(string))
	}

	if v, ok := config["java_version"]; ok {
		siteConfig.JavaVersion = utils.String(v.(string))
	}

	if v, ok := config["elastic_instance_minimum"]; ok && v.(int) > 0 {
		siteConfig.MinimumElasticInstanceCount = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["app_scale_limit"]; ok {
		siteConfig.FunctionAppScaleLimit = utils.Int32(int32(v.(int)))
	}

	if v, ok := config["runtime_scale_monitoring_enabled"]; ok {
		siteConfig.FunctionsRuntimeScaleMonitoringEnabled = utils.Bool(v.(bool))
	}

	if v, ok := config["dotnet_framework_version"]; ok {
		siteConfig.NetFrameworkVersion = utils.String(v.(string))
	}

	if v, ok := config["vnet_route_all_enabled"]; ok {
		siteConfig.VnetRouteAllEnabled = utils.Bool(v.(bool))
	}

	return siteConfig, nil
}

func flattenFunctionAppSiteConfig(input *web.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	if input.AlwaysOn != nil {
		result["always_on"] = *input.AlwaysOn
	}

	if input.Use32BitWorkerProcess != nil {
		result["use_32_bit_worker_process"] = *input.Use32BitWorkerProcess
	}

	if input.WebSocketsEnabled != nil {
		result["websockets_enabled"] = *input.WebSocketsEnabled
	}

	if input.LinuxFxVersion != nil {
		result["linux_fx_version"] = *input.LinuxFxVersion
	}

	if input.HTTP20Enabled != nil {
		result["http2_enabled"] = *input.HTTP20Enabled
	}

	if input.PreWarmedInstanceCount != nil {
		result["pre_warmed_instance_count"] = *input.PreWarmedInstanceCount
	}

	result["ip_restriction"] = flattenAppServiceIpRestriction(input.IPSecurityRestrictions)

	if input.ScmIPSecurityRestrictionsUseMain != nil {
		result["scm_use_main_ip_restriction"] = *input.ScmIPSecurityRestrictionsUseMain
	}

	result["scm_ip_restriction"] = flattenAppServiceIpRestriction(input.ScmIPSecurityRestrictions)

	result["min_tls_version"] = string(input.MinTLSVersion)
	result["ftps_state"] = string(input.FtpsState)
	result["scm_type"] = string(input.ScmType)

	result["cors"] = FlattenWebCorsSettings(input.Cors)

	if input.AutoSwapSlotName != nil {
		result["auto_swap_slot_name"] = *input.AutoSwapSlotName
	}

	if input.HealthCheckPath != nil {
		result["health_check_path"] = *input.HealthCheckPath
	}

	if input.JavaVersion != nil {
		result["java_version"] = *input.JavaVersion
	}

	if input.MinimumElasticInstanceCount != nil {
		result["elastic_instance_minimum"] = *input.MinimumElasticInstanceCount
	}

	if input.FunctionAppScaleLimit != nil {
		result["app_scale_limit"] = *input.FunctionAppScaleLimit
	}

	if input.FunctionsRuntimeScaleMonitoringEnabled != nil {
		result["runtime_scale_monitoring_enabled"] = *input.FunctionsRuntimeScaleMonitoringEnabled
	}

	if input.NetFrameworkVersion != nil {
		result["dotnet_framework_version"] = *input.NetFrameworkVersion
	}

	vnetRouteAllEnabled := false
	if input.VnetRouteAllEnabled != nil {
		vnetRouteAllEnabled = *input.VnetRouteAllEnabled
	}
	result["vnet_route_all_enabled"] = vnetRouteAllEnabled

	results = append(results, result)
	return results
}

func expandFunctionAppConnectionStrings(d *pluginsdk.ResourceData) map[string]*web.ConnStringValueTypePair {
	input := d.Get("connection_string").(*pluginsdk.Set).List()
	output := make(map[string]*web.ConnStringValueTypePair, len(input))

	for _, v := range input {
		vals := v.(map[string]interface{})

		csName := vals["name"].(string)
		csType := vals["type"].(string)
		csValue := vals["value"].(string)

		output[csName] = &web.ConnStringValueTypePair{
			Value: utils.String(csValue),
			Type:  web.ConnectionStringType(csType),
		}
	}

	return output
}

func flattenFunctionAppConnectionStrings(input map[string]*web.ConnStringValueTypePair) interface{} {
	results := make([]interface{}, 0)

	for k, v := range input {
		result := make(map[string]interface{})
		result["name"] = k
		result["type"] = string(v.Type)
		result["value"] = *v.Value
		results = append(results, result)
	}

	return results
}

func flattenFunctionAppSiteCredential(input *web.UserProperties) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] UserProperties is nil")
		return results
	}

	if input.PublishingUserName != nil {
		result["username"] = *input.PublishingUserName
	}

	if input.PublishingPassword != nil {
		result["password"] = *input.PublishingPassword
	}

	return append(results, result)
}

func appSettingsMapToNameValuePair(input map[string]*string) *[]web.NameValuePair {
	if len(input) == 0 {
		return nil
	}
	result := make([]web.NameValuePair, 0)
	for k, v := range input {
		result = append(result, web.NameValuePair{
			Name:  utils.String(k),
			Value: v,
		})
	}

	return &result
}
