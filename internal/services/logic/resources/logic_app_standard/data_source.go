// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package logic_app_standard

import (
	"fmt"
	"log"
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Logic App Standard instance.",
		Read:        dataSourceLogicAppStandardRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Logic App.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.LogicAppStandardName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"app_service_plan_id": {
				Description: "The ID of the App Service Plan.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"app_settings": {
				Description: "A map of key-value pairs for [App Settings](https://learn.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.",
				Type:        pluginsdk.TypeMap,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"use_extension_bundle": {
				Description: "Whether the Logic App should use the bundled extension package.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"bundle_version": {
				Description: "Controls the allowed range for bundle versions.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"client_affinity_enabled": {
				Description: "Should the Logic App send session affinity cookies, which route client requests in the same session to the same instance?",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"client_certificate_mode": {
				Description: "The mode of the Logic App's client certificates requirement for incoming requests.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"enabled": {
				Description: "Whether the Logic App is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"ftp_publish_basic_authentication_enabled": {
				Description: "Whether the default FTP basic authentication publishing profile is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"https_only": {
				Description: "Whether the Logic App can only be accessed via HTTPS.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"scm_publish_basic_authentication_enabled": {
				Description: "Whether the default SCM basic authentication publishing profile is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
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
				Description: "The backend Storage Account name which will be used by this Logic App (e.g. for Stateful workflows data).",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage_account_access_key": {
				Description: "The access key which will be used to access the backend Storage Account for the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Sensitive:   true,
			},

			"storage_key_vault_secret_id": {
				Type:        pluginsdk.TypeString,
				Computed:    true,
				Description: "The Key Vault Secret ID, optionally including version, that contains the connection string to the backend Storage Account for the Logic App.",
			},

			"storage_account_share_name": {
				Description: "The name of the share used by the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"version": {
				Description: "The runtime version associated with the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),

			"custom_domain_verification_id": {
				Description: "The custom domain verification of the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"default_hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kind": {
				Description: "The kind of the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"outbound_ip_addresses": {
				Description: "The outbound IP addresses of the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"possible_outbound_ip_addresses": {
				Description: "The possible outbound IP addresses of the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"public_network_access": {
				Description: "Whether Public Network Access should be enabled or not.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"virtual_network_subnet_id": {
				Description: "The subnet ID for the Logic App.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
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

func dataSourceLogicAppStandardRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppService.WebAppsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

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
				clientCertMode = pointer.FromEnum(props.ClientCertMode)
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

		if strings.HasPrefix(connectionString, "@Microsoft.KeyVault") {
			d.Set("storage_key_vault_secret_id", strings.TrimPrefix(strings.TrimSuffix(connectionString, ")"), "@Microsoft.KeyVault(SecretUri="))
		} else {
			name, key := helpers.ParseWebJobsStorageString(connectionString)
			d.Set("storage_account_name", name)
			d.Set("storage_account_access_key", key)
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
		if err = d.Set("site_config", flattenLogicAppStandardDataSourceSiteConfig(model.Properties)); err != nil {
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

func flattenLogicAppStandardDataSourceConnectionStrings(input *map[string]webapps.ConnStringValueTypePair) any {
	results := make([]any, 0)

	if input == nil || len(*input) == 0 {
		return results
	}

	for k, v := range *input {
		result := make(map[string]any)
		result["name"] = k
		result["type"] = string(v.Type)
		result["value"] = v.Value
		results = append(results, result)
	}

	return results
}

func flattenLogicAppStandardDataSourceSiteConfig(input *webapps.SiteConfig) []any {
	results := make([]any, 0)
	result := make(map[string]any)

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

	result["scm_type"] = pointer.FromEnum(input.ScmType)
	result["scm_min_tls_version"] = pointer.FromEnum(input.ScmMinTlsVersion)
	result["scm_ip_restriction"] = flattenLogicAppStandardIpRestriction(input.ScmIPSecurityRestrictions)
	result["scm_ip_restriction_default_action"] = pointer.FromEnum(input.ScmIPSecurityRestrictionsDefaultAction)
	result["scm_use_main_ip_restriction"] = pointer.From(input.ScmIPSecurityRestrictionsUseMain)

	result["min_tls_version"] = pointer.FromEnum(input.MinTlsVersion)
	result["ftps_state"] = pointer.FromEnum(input.FtpsState)

	result["cors"] = flattenLogicAppStandardCorsSettings(input.Cors)

	result["auto_swap_slot_name"] = pointer.From(input.AutoSwapSlotName)
	result["health_check_path"] = pointer.From(input.HealthCheckPath)
	result["elastic_instance_minimum"] = pointer.From(input.MinimumElasticInstanceCount)
	result["app_scale_limit"] = pointer.From(input.FunctionAppScaleLimit)
	result["runtime_scale_monitoring_enabled"] = pointer.From(input.FunctionsRuntimeScaleMonitoringEnabled)

	result["dotnet_framework_version"] = pointer.From(input.NetFrameworkVersion)

	result["vnet_route_all_enabled"] = pointer.From(input.VnetRouteAllEnabled)

	result["ip_restriction_default_action"] = pointer.FromEnum(input.IPSecurityRestrictionsDefaultAction)

	results = append(results, result)
	return results
}

func flattenLogicAppStandardSiteCredential(input *webapps.User) []any {
	results := make([]any, 0)
	result := make(map[string]any)

	if input == nil || input.Properties == nil {
		log.Printf("[DEBUG] UserProperties is nil")
		return results
	}

	result["username"] = input.Properties.PublishingUserName

	result["password"] = pointer.From(input.Properties.PublishingPassword)

	return append(results, result)
}

func flattenLogicAppStandardCorsSettings(input *webapps.CorsSettings) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	result := make(map[string]any)

	allowedOrigins := make([]any, 0)
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

func flattenHeaders(input map[string][]string) []any {
	output := make([]any, 0)
	headers := make(map[string]any)
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
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Description: "Should the Logic App be loaded at all times?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"cors": schemaLogicAppCorsSettingsDataSource(),

				"ftps_state": {
					Description: "The state of FTP / FTPS service for this Logic App.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"http2_enabled": {
					Description: "Whether the HTTP2 protocol should be enabled.",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"ip_restriction": schemaLogicAppStandardIpRestrictionDataSource(),

				"linux_fx_version": {
					Description: "Linux App Framework and version for the Logic App.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"min_tls_version": {
					Description: "The minimum supported TLS version for the Logic App.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"pre_warmed_instance_count": {
					Description: "The number of pre-warmed instances for this Logic App. Only affects apps on the Premium plan.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},

				"scm_ip_restriction": schemaLogicAppStandardIpRestrictionDataSource(),

				"scm_ip_restriction_default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_use_main_ip_restriction": {
					Description: "Should the Logic App `ip_restriction` configuration be used for the SCM too?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"scm_min_tls_version": {
					Description: "The minimum version of TLS required for SSL requests to the SCM site.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"scm_type": {
					Description: "The type of Source Control used by the Logic App in use by the Windows Function App.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"use_32_bit_worker_process": {
					Description: "Should the Logic App run in 32 bit mode, rather than 64 bit mode?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"websockets_enabled": {
					Description: "Should WebSockets be enabled?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"health_check_path": {
					Description: "Path which will be checked for this Logic App health.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"elastic_instance_minimum": {
					Description: "The number of minimum instances for this Logic App. Only affects apps on the Premium plan.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},

				"app_scale_limit": {
					Description: "The number of workers this Logic App can scale out to. Only applicable to apps on the Consumption and Premium plan.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},

				"runtime_scale_monitoring_enabled": {
					Description: "Should Runtime Scale Monitoring be enabled? Only applicable to apps on the Premium plan.",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"dotnet_framework_version": {
					Description: "The version of the .NET framework's CLR used in this Logic App.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"vnet_route_all_enabled": {
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
				},

				"auto_swap_slot_name": {
					Description: "The Auto-swap slot name.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"ip_restriction_default_action": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func schemaLogicAppCorsSettingsDataSource() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"allowed_origins": {
					Description: "A list of origins which should be able to make cross-origin calls.",
					Type:        pluginsdk.TypeSet,
					Computed:    true,
					Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
				"support_credentials": {
					Description: "Are credentials supported?",
					Type:        pluginsdk.TypeBool,
					Computed:    true,
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
					Description: "The IP Address used for this IP Restriction in CIDR notation.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"service_tag": {
					Description: "The Service Tag used for this IP Restriction.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"virtual_network_subnet_id": {
					Description: "The Virtual Network Subnet ID used for this IP Restriction.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"name": {
					Description: "The name for this IP Restriction.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"priority": {
					Description: "The priority for this IP Restriction. Restrictions are enforced in priority order.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},

				"action": {
					Description: "Does this restriction `Allow` or `Deny` access for this IP range?",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				// lintignore:XS003
				"headers": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							// lintignore:S018
							"x_forwarded_host": {
								Description: "A list of allowed 'X-Forwarded-Host' domains.",
								Type:        pluginsdk.TypeSet,
								Computed:    true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							// lintignore:S018
							"x_forwarded_for": {
								Description: "A list of allowed 'X-Forwarded-For' IPs in CIDR notation.",
								Type:        pluginsdk.TypeSet,
								Computed:    true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							// lintignore:S018
							"x_azure_fdid": {
								Description: "A list of allowed Azure FrontDoor IDs in UUID notation.",
								Type:        pluginsdk.TypeSet,
								Computed:    true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},

							// lintignore:S018
							"x_fd_health_probe": {
								Description: "A list to allow the Azure FrontDoor health probe header.",
								Type:        pluginsdk.TypeSet,
								Computed:    true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
				},
			},
		},
	}
}
