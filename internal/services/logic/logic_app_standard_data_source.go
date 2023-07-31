// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceLogicAppStandard() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceLogicAppStandardRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.LogicAppStandardName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"app_service_plan_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"app_settings": {
				Type:     pluginsdk.TypeMap,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"use_extension_bundle": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"bundle_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"client_affinity_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"client_certificate_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

			"site_config": schemaLogicAppStandardSiteConfig(),

			"connection_string": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"value": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},

			"storage_account_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"storage_account_share_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),

			"custom_domain_verification_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"outbound_ip_addresses": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"possible_outbound_ip_addresses": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"virtual_network_subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"site_credential": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"password": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLogicAppStandardRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	subscriptionId := meta.(*clients.Client).Web.AppServicesClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLogicAppStandardID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Logic App Standard %s was not found", id)
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Standard %s: %+v", id, err)
	}

	d.SetId(id.ID())

	appSettingsResp, err := client.ListApplicationSettings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("[ERROR] Listing application settings for %s: %+v", id, err)
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("[ERROR] Listing connection strings for %s: %+v", id, err)
	}

	siteCredFuture, err := client.ListPublishingCredentials(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("[ERROR] Listing publishing credentials for %s: %+v", id, err)
	}
	if err = siteCredFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("[ERROR] Waiting to list the publishing credentials for %s: %+v", id, err)
	}
	siteCredResp, err := siteCredFuture.Result(*client)
	if err != nil {
		return fmt.Errorf("[ERROR] Retrieving the publishing credentials for %s: %+v", id, err)
	}

	d.Set("kind", resp.Kind)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.SiteProperties; props != nil {
		d.Set("app_service_plan_id", props.ServerFarmID)
		d.Set("enabled", props.Enabled)
		d.Set("default_hostname", props.DefaultHostName)
		d.Set("https_only", props.HTTPSOnly)
		d.Set("outbound_ip_addresses", props.OutboundIPAddresses)
		d.Set("possible_outbound_ip_addresses", props.PossibleOutboundIPAddresses)
		d.Set("client_affinity_enabled", props.ClientAffinityEnabled)
		d.Set("custom_domain_verification_id", props.CustomDomainVerificationID)

		clientCertMode := ""
		if props.ClientCertEnabled != nil && *props.ClientCertEnabled {
			clientCertMode = string(props.ClientCertMode)
		}
		d.Set("client_certificate_mode", clientCertMode)
	}

	appSettings := flattenLogicAppStandardDataSourceAppSettings(appSettingsResp.Properties)

	if err = d.Set("virtual_network_subnet_id", resp.SiteProperties.VirtualNetworkSubnetID); err != nil {
		return err
	}

	if err = d.Set("connection_string", flattenLogicAppStandardDataSourceConnectionStrings(connectionStringsResp.Properties)); err != nil {
		return err
	}

	connectionString := appSettings["AzureWebJobsStorage"]

	// This teases out the necessary attributes from the storage connection string
	connectionStringParts := strings.Split(connectionString, ";")
	for _, part := range connectionStringParts {
		if strings.HasPrefix(part, "AccountName") {
			accountNameParts := strings.Split(part, "AccountName=")
			if len(accountNameParts) > 1 {
				d.Set("storage_account_name", accountNameParts[1])
			}
		}
		if strings.HasPrefix(part, "AccountKey") {
			accountKeyParts := strings.Split(part, "AccountKey=")
			if len(accountKeyParts) > 1 {
				d.Set("storage_account_access_key", accountKeyParts[1])
			}
		}
	}

	d.Set("version", appSettings["FUNCTIONS_EXTENSION_VERSION"])

	if _, ok := appSettings["AzureFunctionsJobHost__extensionBundle__id"]; ok {
		d.Set("use_extension_bundle", true)
		if val, ok := appSettings["AzureFunctionsJobHost__extensionBundle__version"]; ok {
			d.Set("bundle_version", val)
		}
	} else {
		d.Set("use_extension_bundle", false)
		d.Set("bundle_version", "[1.*, 2.0.0)")
	}

	d.Set("storage_account_share_name", appSettings["WEBSITE_CONTENTSHARE"])

	// Remove all the settings that are created by this resource so we don't to have to specify in app_settings
	// block whenever we use azurerm_logic_app_standard.
	delete(appSettings, "WEBSITE_CONTENTAZUREFILECONNECTIONSTRING")
	delete(appSettings, "APP_KIND")
	delete(appSettings, "AzureFunctionsJobHost__extensionBundle__id")
	delete(appSettings, "AzureFunctionsJobHost__extensionBundle__version")
	delete(appSettings, "AzureWebJobsDashboard")
	delete(appSettings, "AzureWebJobsStorage")
	delete(appSettings, "FUNCTIONS_EXTENSION_VERSION")
	delete(appSettings, "WEBSITE_CONTENTSHARE")

	if err = d.Set("app_settings", appSettings); err != nil {
		return err
	}

	identity := flattenLogicAppStandardDataSourceIdentity(resp.Identity)
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	configResp, err := client.GetConfiguration(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		return fmt.Errorf("retrieving the configuration for %s: %+v", id, err)
	}

	siteConfig := flattenLogicAppStandardDataSourceSiteConfig(configResp.SiteConfig)
	if err = d.Set("site_config", siteConfig); err != nil {
		return err
	}

	siteCred := flattenLogicAppStandardDataSourceSiteCredential(siteCredResp.UserProperties)
	if err = d.Set("site_credential", siteCred); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenLogicAppStandardDataSourceAppSettings(input map[string]*string) map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = *v
	}

	return output
}

func flattenLogicAppStandardDataSourceConnectionStrings(input map[string]*web.ConnStringValueTypePair) interface{} {
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

func flattenLogicAppStandardDataSourceIdentity(identity *web.ManagedServiceIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	principalId := ""
	if identity.PrincipalID != nil {
		principalId = *identity.PrincipalID
	}

	tenantId := ""
	if identity.TenantID != nil {
		tenantId = *identity.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"principal_id": principalId,
			"tenant_id":    tenantId,
			"type":         string(identity.Type),
		},
	}
}

func flattenLogicAppStandardDataSourceSiteConfig(input *web.SiteConfig) []interface{} {
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

	result["ip_restriction"] = flattenLogicAppStandardIpRestriction(input.IPSecurityRestrictions)

	result["scm_type"] = string(input.ScmType)
	result["scm_min_tls_version"] = string(input.ScmMinTLSVersion)
	result["scm_ip_restriction"] = flattenLogicAppStandardIpRestriction(input.ScmIPSecurityRestrictions)

	if input.ScmIPSecurityRestrictionsUseMain != nil {
		result["scm_use_main_ip_restriction"] = *input.ScmIPSecurityRestrictionsUseMain
	}

	result["min_tls_version"] = string(input.MinTLSVersion)
	result["ftps_state"] = string(input.FtpsState)

	result["cors"] = flattenLogicAppStandardCorsSettings(input.Cors)

	if input.AutoSwapSlotName != nil {
		result["auto_swap_slot_name"] = *input.AutoSwapSlotName
	}

	if input.HealthCheckPath != nil {
		result["health_check_path"] = *input.HealthCheckPath
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

func flattenLogicAppStandardDataSourceSiteCredential(input *web.UserProperties) []interface{} {
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
