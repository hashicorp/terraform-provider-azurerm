// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"ftp_publish_basic_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"scm_publish_basic_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"site_config": schemaLogicAppStandardSiteConfigDataSource(),

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

			"tags": commonschema.TagsDataSource(),

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

			"public_network_access": {
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
	client := meta.(*clients.Client).AppService.WebAppsClient
	subscriptionId := meta.(*clients.Client).Web.AppServicesClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewAppServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("kind", pointer.From(model.Kind))
		d.Set("location", location.Normalize(model.Location))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %s", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("app_service_plan_id", pointer.From(props.ServerFarmId))
			d.Set("enabled", pointer.From(props.Enabled))
			d.Set("default_hostname", pointer.From(props.DefaultHostName))
			d.Set("https_only", pointer.From(props.HTTPSOnly))
			d.Set("outbound_ip_addresses", pointer.From(props.OutboundIPAddresses))
			d.Set("possible_outbound_ip_addresses", pointer.From(props.PossibleOutboundIPAddresses))
			d.Set("client_affinity_enabled", pointer.From(props.ClientAffinityEnabled))
			d.Set("custom_domain_verification_id", pointer.From(props.CustomDomainVerificationId))
			d.Set("public_network_access", pointer.From(props.PublicNetworkAccess))

			clientCertMode := ""
			if props.ClientCertEnabled != nil && *props.ClientCertEnabled {
				clientCertMode = string(pointer.From(props.ClientCertMode))
			}
			d.Set("client_certificate_mode", clientCertMode)

			d.Set("virtual_network_subnet_id", props.VirtualNetworkSubnetId)
		}
	}

	appSettingsResp, err := client.ListApplicationSettings(ctx, id)
	if err != nil {
		return fmt.Errorf("listing application settings for %s: %+v", id, err)
	}
	if model := appSettingsResp.Model; model != nil {
		appSettings := pointer.From(model.Properties)

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
	}

	connectionStringsResp, err := client.ListConnectionStrings(ctx, id)
	if err != nil {
		return fmt.Errorf("listing connection strings for %s: %+v", id, err)
	}

	if model := connectionStringsResp.Model; model != nil {
		if err = d.Set("connection_string", flattenLogicAppStandardDataSourceConnectionStrings(model.Properties)); err != nil {
			return err
		}
	}

	ftpBasicAuth, err := client.GetFtpAllowed(ctx, id)
	if err != nil || ftpBasicAuth.Model == nil {
		return fmt.Errorf("retrieving FTP publish basic authentication policy for %s: %+v", id, err)
	}

	if props := ftpBasicAuth.Model.Properties; props != nil {
		d.Set("ftp_publish_basic_authentication_enabled", props.Allow)
	}

	scmBasicAuth, err := client.GetScmAllowed(ctx, id)
	if err != nil || scmBasicAuth.Model == nil {
		return fmt.Errorf("retrieving SCM publish basic authentication policy for %s: %+v", id, err)
	}

	if props := scmBasicAuth.Model.Properties; props != nil {
		d.Set("scm_publish_basic_authentication_enabled", props.Allow)
	}

	configResp, err := client.GetConfiguration(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving the configuration for %s: %+v", id, err)
	}

	if model := configResp.Model; model != nil {
		siteConfig := flattenLogicAppStandardDataSourceSiteConfig(model.Properties)
		if err = d.Set("site_config", siteConfig); err != nil {
			return err
		}
	}

	siteCredentials, err := helpers.ListPublishingCredentials(ctx, client, id)
	if err != nil {
		return fmt.Errorf("listing publishing credentials for %s: %+v", id, err)
	}

	if err = d.Set("site_credential", flattenLogicAppStandardSiteCredential(siteCredentials)); err != nil {
		return err
	}

	return nil
}

func flattenLogicAppStandardDataSourceConnectionStrings(input *map[string]webapps.ConnStringValueTypePair) interface{} {
	results := make([]interface{}, 0)

	if input == nil || len(*input) == 0 {
		return results
	}

	for k, v := range *input {
		result := make(map[string]interface{})
		result["name"] = k
		result["type"] = string(v.Type)
		result["value"] = v.Value
		results = append(results, result)
	}

	return results
}

func flattenLogicAppStandardDataSourceSiteConfig(input *webapps.SiteConfig) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil {
		log.Printf("[DEBUG] SiteConfig is nil")
		return results
	}

	result["always_on"] = pointer.From(input.AlwaysOn)
	result["use_32_bit_worker_process"] = pointer.From(input.Use32BitWorkerProcess)
	result["websockets_enabled"] = pointer.From(input.WebSocketsEnabled)
	result["linux_fx_version"] = pointer.From(input.LinuxFxVersion)
	result["http2_enabled"] = pointer.From(input.HTTP20Enabled)
	result["pre_warmed_instance_count"] = pointer.From(input.PreWarmedInstanceCount)

	result["ip_restriction"] = flattenLogicAppStandardIpRestriction(input.IPSecurityRestrictions)

	result["scm_type"] = string(pointer.From(input.ScmType))
	result["scm_min_tls_version"] = string(pointer.From(input.ScmMinTlsVersion))
	result["scm_ip_restriction"] = flattenLogicAppStandardIpRestriction(input.ScmIPSecurityRestrictions)

	result["scm_use_main_ip_restriction"] = pointer.From(input.ScmIPSecurityRestrictionsUseMain)

	result["min_tls_version"] = string(pointer.From(input.MinTlsVersion))
	result["ftps_state"] = string(pointer.From(input.FtpsState))

	result["cors"] = flattenLogicAppStandardCorsSettings(input.Cors)

	result["auto_swap_slot_name"] = pointer.From(input.AutoSwapSlotName)
	result["health_check_path"] = pointer.From(input.HealthCheckPath)
	result["elastic_instance_minimum"] = pointer.From(input.MinimumElasticInstanceCount)
	result["app_scale_limit"] = pointer.From(input.FunctionAppScaleLimit)
	result["runtime_scale_monitoring_enabled"] = pointer.From(input.FunctionsRuntimeScaleMonitoringEnabled)

	result["dotnet_framework_version"] = pointer.From(input.NetFrameworkVersion)

	result["vnet_route_all_enabled"] = pointer.From(input.VnetRouteAllEnabled)

	results = append(results, result)
	return results
}

func flattenLogicAppStandardSiteCredential(input *webapps.User) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if input == nil || input.Properties == nil {
		log.Printf("[DEBUG] UserProperties is nil")
		return results
	}

	result["username"] = input.Properties.PublishingUserName

	result["password"] = pointer.From(input.Properties.PublishingPassword)

	return append(results, result)
}

func flattenLogicAppStandardCorsSettings(input *webapps.CorsSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	allowedOrigins := make([]interface{}, 0)
	if s := input.AllowedOrigins; s != nil {
		for _, v := range *s {
			allowedOrigins = append(allowedOrigins, v)
		}
	}
	result["allowed_origins"] = pluginsdk.NewSet(pluginsdk.HashString, allowedOrigins)

	if input.SupportCredentials != nil {
		result["support_credentials"] = *input.SupportCredentials
	}

	return append(results, result)
}

func flattenHeaders(input map[string][]string) []interface{} {
	output := make([]interface{}, 0)
	headers := make(map[string]interface{})
	if input == nil {
		return output
	}

	if forwardedHost, ok := input["x-forwarded-host"]; ok && len(forwardedHost) > 0 {
		headers["x_forwarded_host"] = forwardedHost
	}
	if forwardedFor, ok := input["x-forwarded-for"]; ok && len(forwardedFor) > 0 {
		headers["x_forwarded_for"] = forwardedFor
	}
	if fdids, ok := input["x-azure-fdid"]; ok && len(fdids) > 0 {
		headers["x_azure_fdid"] = fdids
	}
	if healthProbe, ok := input["x-fd-healthprobe"]; ok && len(healthProbe) > 0 {
		headers["x_fd_health_probe"] = healthProbe
	}

	return append(output, headers)
}

func schemaLogicAppStandardSiteConfigDataSource() *pluginsdk.Schema {
	schema := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"cors": schemaLogicAppCorsSettingsDataSource(),

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": schemaLogicAppStandardIpRestrictionDataSource(),

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
						string(webapps.SupportedTlsVersionsOnePointTwo),
					}, false),
				},

				"pre_warmed_instance_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
				},

				"scm_ip_restriction": schemaLogicAppStandardIpRestrictionDataSource(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_min_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.SupportedTlsVersionsOnePointTwo),
					}, false),
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.ScmTypeBitbucketGit),
						string(webapps.ScmTypeBitbucketHg),
						string(webapps.ScmTypeCodePlexGit),
						string(webapps.ScmTypeCodePlexHg),
						string(webapps.ScmTypeDropbox),
						string(webapps.ScmTypeExternalGit),
						string(webapps.ScmTypeExternalHg),
						string(webapps.ScmTypeGitHub),
						string(webapps.ScmTypeLocalGit),
						string(webapps.ScmTypeNone),
						string(webapps.ScmTypeOneDrive),
						string(webapps.ScmTypeTfs),
						string(webapps.ScmTypeVSO),
						string(webapps.ScmTypeVSTSRM),
					}, false),
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

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"elastic_instance_minimum": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(0, 20),
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
						"v8.0",
					}, false),
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
				},

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}

	if !features.FivePointOh() {
		schema.Elem.(*pluginsdk.Resource).Schema["public_network_access_enabled"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Computed:   true,
			Deprecated: "the `site_config.public_network_access_enabled` property has been superseded by the `public_network_access` property and will be removed in v5.0 of the AzureRM Provider.",
		}
		schema.Elem.(*pluginsdk.Resource).Schema["scm_min_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.SupportedTlsVersionsOnePointZero),
				string(webapps.SupportedTlsVersionsOnePointOne),
				string(webapps.SupportedTlsVersionsOnePointTwo),
			}, false),
		}
		schema.Elem.(*pluginsdk.Resource).Schema["min_tls_version"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webapps.SupportedTlsVersionsOnePointZero),
				string(webapps.SupportedTlsVersionsOnePointOne),
				string(webapps.SupportedTlsVersionsOnePointTwo),
			}, false),
		}
	}

	return schema
}

func schemaLogicAppCorsSettingsDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
				"support_credentials": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func schemaLogicAppStandardIpRestrictionDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"ip_address": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"service_tag": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"virtual_network_subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"priority": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      65000,
					ValidateFunc: validation.IntBetween(1, math.MaxInt32),
				},

				"action": {
					Type:     pluginsdk.TypeString,
					Default:  "Allow",
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"Allow",
						"Deny",
					}, false),
				},

				// lintignore:XS003
				"headers": {
					Type:       pluginsdk.TypeList,
					Optional:   true,
					Computed:   true,
					MaxItems:   1,
					ConfigMode: pluginsdk.SchemaConfigModeAttr,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							// lintignore:S018
							"x_forwarded_host": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							// lintignore:S018
							"x_forwarded_for": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsCIDR,
								},
							},

							// lintignore:S018
							"x_azure_fdid": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 8,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.IsUUID,
								},
							},

							// lintignore:S018
							"x_fd_health_probe": {
								Type:     pluginsdk.TypeSet,
								Optional: true,
								MaxItems: 1,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										"1",
									}, false),
								},
							},
						},
					},
				},
			},
		},
	}
}
