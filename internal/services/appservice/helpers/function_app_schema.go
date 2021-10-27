package helpers

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	apimValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	StorageStringFmt = "DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s"
)

type SiteConfigLinuxFunctionApp struct {
	AlwaysOn                      bool                               `tfschema:"always_on"`
	AppCommandLine                string                             `tfschema:"app_command_line"`
	ApiDefinition                 string                             `tfschema:"api_definition_url"`
	ApiManagementConfigId         string                             `tfschema:"api_management_api_id"`
	AppInsightsInstrumentationKey string                             `tfschema:"application_insights_key"` // App Insights Instrumentation Key
	AppInsightsConnectionString   string                             `tfschema:"application_insights_connection_string"`
	AppScaleLimit                 int                                `tfschema:"app_scale_limit"`
	UseManagedIdentityACR         bool                               `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI          string                             `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments              []string                           `tfschema:"default_documents"`
	ElasticInstanceMinimum        int                                `tfschema:"elastic_instance_minimum"`
	Http2Enabled                  bool                               `tfschema:"http2_enabled"`
	IpRestriction                 []IpRestriction                    `tfschema:"ip_restriction"`
	LoadBalancing                 string                             `tfschema:"load_balancing_mode"` // TODO - Valid for FunctionApps?
	LocalMysql                    bool                               `tfschema:"local_mysql"`
	ManagedPipelineMode           string                             `tfschema:"managed_pipeline_mode"`
	PreWarmedInstanceCount        int                                `tfschema:"pre_warmed_instance_count"`
	RemoteDebugging               bool                               `tfschema:"remote_debugging"`
	RemoteDebuggingVersion        string                             `tfschema:"remote_debugging_version"`
	RuntimeScaleMonitoring        bool                               `tfschema:"runtime_scale_monitoring_enabled"`
	ScmIpRestriction              []IpRestriction                    `tfschema:"scm_ip_restriction"`
	ScmType                       string                             `tfschema:"scm_type"` // Computed?
	ScmUseMainIpRestriction       bool                               `tfschema:"scm_use_main_ip_restriction"`
	Use32BitWorker                bool                               `tfschema:"use_32_bit_worker"`
	WebSockets                    bool                               `tfschema:"websockets_enabled"`
	FtpsState                     string                             `tfschema:"ftps_state"`
	HealthCheckPath               string                             `tfschema:"health_check_path"`
	NumberOfWorkers               int                                `tfschema:"number_of_workers"`
	ApplicationStack              []ApplicationStackLinuxFunctionApp `tfschema:"application_stack"`
	MinTlsVersion                 string                             `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion              string                             `tfschema:"scm_minimum_tls_version"`
	AutoSwapSlotName              string                             `tfschema:"auto_swap_slot_name"`
	Cors                          []CorsSetting                      `tfschema:"cors"`
	DetailedErrorLogging          bool                               `tfschema:"detailed_error_logging"`
	LinuxFxVersion                string                             `tfschema:"linux_fx_version"`
	VnetRouteAllEnabled           bool                               `tfschema:"vnet_route_all_enabled"` // Not supported in Dynamic plans
}

func SiteConfigSchemaLinuxFunctionApp() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Computed:    true, // Note - several factors change the default for this, so needs to be computed.
					Description: "If this Linux Web App is Always On enabled. Defaults to `false`.",
				},

				"api_management_api_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: apimValidate.ApiID,
					Description:  "The ID of the API Management API for this Linux Function App.",
				},

				"api_definition_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					Description:  "The URL of the API definition that describes this Linux Function App.",
				},

				"app_command_line": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The App command line to launch.",
				},

				"app_scale_limit": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "The number of workers this function app can scale out to. Only applicable to apps on the Consumption and Premium plan.",
					// TODO Validation?
				},

				"application_insights_key": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
					RequiredWith: []string{
						"site_config.0.application_insights_connection_string",
					},
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Instrumentation Key for connecting the Linux Function App to Application Insights.",
				},

				"application_insights_connection_string": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
					RequiredWith: []string{
						"site_config.0.application_insights_key",
					},
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Connection String for linking the Linux Function App to Application Insights.",
				},

				"application_stack": linuxFunctionAppStackSchema(),

				"container_registry_use_managed_identity": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should connections for Azure Container Registry use Managed Identity.",
				},

				"container_registry_managed_identity_client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
					Description:  "The Client ID of the Managed Service Identity to use for connections to the Azure Container Registry.",
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
					Description: "Specifies a list of Default Documents for the Linux Web App.",
				},

				"elastic_instance_minimum": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "The number of minimum instances for this Linux Function App. Only affects apps on Elastic Premium plans.",
				},

				"http2_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Specifies if the http2 protocol should be enabled. Defaults to `false`.",
				},

				"ip_restriction": IpRestrictionSchema(),

				"scm_use_main_ip_restriction": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Linux Function App `ip_restriction` configuration be used for the SCM also.",
				},

				"scm_ip_restriction": IpRestrictionSchema(),

				"local_mysql": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Use Local MySQL. Defaults to `false`.",
				},

				"load_balancing_mode": { // Supported on Function Apps?
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "LeastRequests",
					ValidateFunc: validation.StringInSlice([]string{
						"LeastRequests", // Service default
						"WeightedRoundRobin",
						"LeastResponseTime",
						"WeightedTotalTraffic",
						"RequestHash",
						"PerSiteRoundRobin",
					}, false),
					Description: "The Site load balancing mode. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.",
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
					}, false),
					Description: "The Managed Pipeline mode. Possible values include: `Integrated`, `Classic`. Defaults to `Integrated`.",
				},

				"pre_warmed_instance_count": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Computed:    true, // Variable defaults depending on plan etc
					Description: "The number of pre-warmed instances for this function app. Only affects apps on an Elastic Premium plan.",
				},

				"remote_debugging": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should Remote Debugging be enabled. Defaults to `false`.",
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"VS2017",
						"VS2019",
					}, false),
					Description: "The Remote Debugging Version. Possible values include `VS2017` and `VS2019`",
				},

				"runtime_scale_monitoring_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Description: "Should Functions Runtime Scale Monitoring be enabled.",
				},

				"scm_type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The SCM Type in use by the Linux Function App.",
				},

				"use_32_bit_worker": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should the Linux Web App use a 32-bit worker.",
				},

				"websockets_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should Web Sockets be enabled. Defaults to `false`.",
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.FtpsStateDisabled),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.FtpsStateAllAllowed),
						string(web.FtpsStateDisabled),
						string(web.FtpsStateFtpsOnly),
					}, false),
					Description: "State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `Disabled`.",
				},

				"health_check_path": {
					Type:        pluginsdk.TypeString,
					Optional:    true,
					Description: "The path to be checked for this function app health.",
				},

				"number_of_workers": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
					Description:  "The number of Workers for this Linux Function App.",
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
					Description: "The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
					Description: "Configures the minimum version of TLS required for SSL requests to the SCM site Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.",
				},

				"cors": CorsSettingsSchema(),

				"auto_swap_slot_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					// TODO - Add slot name validation here when the resource is added
					Description: "The Linux Function App Slot Name to automatically swap to when deployment to that slot is successfully completed.",
				},

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.",
				},

				"detailed_error_logging": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is detailed error logging enabled",
				},

				"linux_fx_version": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The Linux FX Version",
				},
			},
		},
	}
}

type ApplicationStackLinuxFunctionApp struct {
	// Note - Function Apps differ to Web Apps here. They do not use the named properties in the SiteConfig block and exclusively use the app_settings map
	DotNetVersion string                   `tfschema:"dotnet_version"` // Supported values `3.1`. Version 6 is in preview on Windows Only
	NodeVersion   string                   `tfschema:"node_version"`   // Supported values `12LTS`, `14LTS`
	PythonVersion string                   `tfschema:"python_version"` // Supported values `3.9`, `3.8`, `3.7`, `3.6`
	JavaVersion   string                   `tfschema:"java_version"`   // Supported values `8`, `11`
	CustomHandler bool                     `tfschema:"use_custom"`     // Supported values `true`
	Docker        []ApplicationStackDocker `tfschema:"docker"`         // Needs ElasticPremium or Basic (B1) Standard (S 1-3) or Premium(PxV2 or PxV3) LINUX Service Plan
}

type ApplicationStackDocker struct {
	RegistryURL      string `tfschema:"registry_url"`
	RegistryUsername string `tfschema:"registry_username"`
	RegistryPassword string `tfschema:"registry_password"`
	ImageName        string `tfschema:"image_name"`
	ImageTag         string `tfschema:"image_tag"`
}

func linuxFunctionAppStackSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"3.1",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom",
					},
					Description: "The version of .Net. Possible values are `3.1`",
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"3.9",
						"3.8",
						"3.7",
						"3.6", // EOL Soon, just remove it now?
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom",
					},
					Description: "The version of Python to use. Possible values include `3.9`, `3.8`, `3.7`, and `3.6`, ",
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"12",
						"14",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom",
					},
					Description: "The version of Node to use. Possible values include `12`, and `14`",
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"8",
						"11",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom",
					},
					Description: "The version of Java to use. Possible values are `8`, and `11`",
				},

				"docker": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*schema.Schema{
							"registry_url": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The URL of the docker registry.",
							},

							"registry_username": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								Sensitive:    true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The username to use for connections to the registry.",
							},

							"registry_password": {
								Type:        pluginsdk.TypeString,
								Required:    true,
								Sensitive:   true, // Note: whilst it's not a good idea, this _can_ be blank...
								Description: "The password for the account to use to connect to the registry.",
							},

							"image_name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The name of the Docker image to use.",
							},

							"image_tag": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
								Description:  "The image tag of the image to use.",
							},
						},
					},
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom",
					},
					Description: "A docker block",
				},

				"use_custom": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.docker",
						"site_config.0.application_stack.0.use_custom",
					},
				},
			},
		},
	}
}

func ExpandSiteConfigLinuxFunctionApp(siteConfig []SiteConfigLinuxFunctionApp, existing *web.SiteConfig, metadata sdk.ResourceMetaData, version string, storageString string, storageUsesMSI bool) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}
	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
		// need to zero fxversion to re-calculate based on changes below or removing app_stack doesn't apply
		expanded.LinuxFxVersion = utils.String("")
	}

	appSettings := make([]web.NameValuePair, 0)

	appSettings = append(appSettings, web.NameValuePair{
		Name:  utils.String("FUNCTIONS_EXTENSION_VERSION"),
		Value: utils.String(version),
	})

	if storageUsesMSI {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("AzureWebJobsStorage__accountName"),
			Value: utils.String(storageString),
		})
	} else {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("AzureWebJobsStorage"),
			Value: utils.String(storageString),
		})
	}

	linuxSiteConfig := siteConfig[0]

	expanded.AlwaysOn = utils.Bool(linuxSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.app_scale_limit") {
		expanded.FunctionAppScaleLimit = utils.Int32(int32(linuxSiteConfig.AppScaleLimit))
	}

	if linuxSiteConfig.AppInsightsConnectionString != "" {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("APPLICATIONINSIGHTS_CONNECTION_STRING"),
			Value: utils.String(linuxSiteConfig.AppInsightsConnectionString),
		})
	}

	if linuxSiteConfig.AppInsightsInstrumentationKey != "" {
		appSettings = append(appSettings, web.NameValuePair{
			Name:  utils.String("APPINSIGHTS_INSTRUMENTATIONKEY"),
			Value: utils.String(linuxSiteConfig.AppInsightsInstrumentationKey),
		})
	}

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: utils.String(linuxSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: utils.String(linuxSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = utils.String(linuxSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") && len(linuxSiteConfig.ApplicationStack) > 0 {
		if len(linuxSiteConfig.ApplicationStack) > 0 {
			linuxAppStack := linuxSiteConfig.ApplicationStack[0]
			if linuxAppStack.DotNetVersion != "" {
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
					Value: utils.String("dotnet"),
				})
				linuxSiteConfig.LinuxFxVersion = fmt.Sprintf("DOTNET|%s", linuxAppStack.DotNetVersion)
			}

			if linuxAppStack.NodeVersion != "" {
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
					Value: utils.String("node"),
				})
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("WEBSITE_NODE_DEFAULT_VERSION"),
					Value: utils.String(linuxAppStack.NodeVersion),
				})
				linuxSiteConfig.LinuxFxVersion = fmt.Sprintf("Node|%s", linuxAppStack.NodeVersion)
			}

			if linuxAppStack.PythonVersion != "" {
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
					Value: utils.String("python"),
				})
				linuxSiteConfig.LinuxFxVersion = fmt.Sprintf("Python|%s", linuxAppStack.PythonVersion)
			}

			if linuxAppStack.JavaVersion != "" {
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
					Value: utils.String("java"),
				})
				linuxSiteConfig.LinuxFxVersion = fmt.Sprintf("Java|%s", linuxAppStack.JavaVersion)
			}

			if linuxAppStack.CustomHandler {
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
					Value: utils.String("custom"),
				})
				linuxSiteConfig.LinuxFxVersion = "" // Custom needs an explicit empty string here
			}

			if linuxAppStack.Docker != nil && len(linuxAppStack.Docker) == 1 {
				dockerConfig := linuxAppStack.Docker[0]
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("DOCKER_REGISTRY_SERVER_URL"),
					Value: utils.String(dockerConfig.RegistryURL),
				})
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("DOCKER_REGISTRY_SERVER_USERNAME"),
					Value: utils.String(dockerConfig.RegistryUsername),
				})
				appSettings = append(appSettings, web.NameValuePair{
					Name:  utils.String("DOCKER_REGISTRY_SERVER_PASSWORD"),
					Value: utils.String(dockerConfig.RegistryPassword),
				})
				linuxSiteConfig.LinuxFxVersion = fmt.Sprintf("DOCKER|%s/%s:%s", dockerConfig.RegistryURL, dockerConfig.ImageName, dockerConfig.ImageTag)
			}
		} else {
			appSettings = append(appSettings, web.NameValuePair{
				Name:  utils.String("FUNCTIONS_WORKER_RUNTIME"),
				Value: utils.String(""),
			})
			linuxSiteConfig.LinuxFxVersion = ""
		}
	}

	expanded.AcrUseManagedIdentityCreds = utils.Bool(linuxSiteConfig.UseManagedIdentityACR)

	expanded.VnetRouteAllEnabled = utils.Bool(linuxSiteConfig.VnetRouteAllEnabled)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = utils.String(linuxSiteConfig.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &linuxSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = utils.Bool(linuxSiteConfig.Http2Enabled)

	expanded.LocalMySQLEnabled = utils.Bool(linuxSiteConfig.LocalMysql)

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = utils.Bool(linuxSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	// TODO - Supported for Function Apps?
	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(linuxSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(linuxSiteConfig.ManagedPipelineMode)
	}

	expanded.RemoteDebuggingEnabled = utils.Bool(linuxSiteConfig.RemoteDebugging)

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = utils.String(linuxSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = utils.Bool(linuxSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = utils.Bool(linuxSiteConfig.WebSockets)

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(linuxSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = utils.String(linuxSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.number_of_workers") {
		expanded.NumberOfWorkers = utils.Int32(int32(linuxSiteConfig.NumberOfWorkers))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_swap_slot_name") {
		expanded.AutoSwapSlotName = utils.String(linuxSiteConfig.AutoSwapSlotName)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(linuxSiteConfig.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.pre_warmed_instance_count") {
		expanded.PreWarmedInstanceCount = utils.Int32(int32(linuxSiteConfig.PreWarmedInstanceCount))
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = utils.Bool(linuxSiteConfig.VnetRouteAllEnabled)
	}

	expanded.AppSettings = &appSettings

	return expanded, nil
}

func FlattenSiteConfigLinuxFunctionApp(functionAppSiteConfig *web.SiteConfig) (*SiteConfigLinuxFunctionApp, error) {
	if functionAppSiteConfig == nil {
		return nil, fmt.Errorf("flattening site config: SiteConfig was nil")
	}

	result := &SiteConfigLinuxFunctionApp{
		AppCommandLine:         utils.NormalizeNilableString(functionAppSiteConfig.AppCommandLine),
		AppScaleLimit:          int(utils.NormaliseNilableInt32(functionAppSiteConfig.FunctionAppScaleLimit)),
		AutoSwapSlotName:       utils.NormalizeNilableString(functionAppSiteConfig.AutoSwapSlotName),
		ContainerRegistryMSI:   utils.NormalizeNilableString(functionAppSiteConfig.AcrUserManagedIdentityID),
		HealthCheckPath:        utils.NormalizeNilableString(functionAppSiteConfig.HealthCheckPath),
		LinuxFxVersion:         utils.NormalizeNilableString(functionAppSiteConfig.LinuxFxVersion),
		LoadBalancing:          string(functionAppSiteConfig.LoadBalancing),
		ManagedPipelineMode:    string(functionAppSiteConfig.ManagedPipelineMode),
		NumberOfWorkers:        int(utils.NormaliseNilableInt32(functionAppSiteConfig.NumberOfWorkers)),
		ScmType:                string(functionAppSiteConfig.ScmType),
		FtpsState:              string(functionAppSiteConfig.FtpsState),
		MinTlsVersion:          string(functionAppSiteConfig.MinTLSVersion),
		ScmMinTlsVersion:       string(functionAppSiteConfig.ScmMinTLSVersion),
		PreWarmedInstanceCount: int(utils.NormaliseNilableInt32(functionAppSiteConfig.PreWarmedInstanceCount)),
		ElasticInstanceMinimum: int(utils.NormaliseNilableInt32(functionAppSiteConfig.MinimumElasticInstanceCount)),
	}

	if functionAppSiteConfig.AlwaysOn != nil {
		result.AlwaysOn = *functionAppSiteConfig.AlwaysOn
	}

	if v := functionAppSiteConfig.APIDefinition; v != nil && v.URL != nil {
		result.ApiDefinition = *v.URL
	}

	if v := functionAppSiteConfig.APIManagementConfig; v != nil && v.ID != nil {
		result.ApiManagementConfigId = *v.ID
	}

	if v := functionAppSiteConfig.Use32BitWorkerProcess; v != nil {
		result.Use32BitWorker = *v
	}

	if v := functionAppSiteConfig.WebSocketsEnabled; v != nil {
		result.WebSockets = *v
	}

	if v := functionAppSiteConfig.HTTP20Enabled; v != nil {
		result.Http2Enabled = *v
	}

	if functionAppSiteConfig.IPSecurityRestrictions != nil {
		result.IpRestriction = FlattenIpRestrictions(functionAppSiteConfig.IPSecurityRestrictions)
	}

	if v := functionAppSiteConfig.ScmIPSecurityRestrictionsUseMain; v != nil {
		result.ScmUseMainIpRestriction = *v
	}

	if functionAppSiteConfig.ScmIPSecurityRestrictions != nil {
		result.ScmIpRestriction = FlattenIpRestrictions(functionAppSiteConfig.ScmIPSecurityRestrictions)
	}

	if v := functionAppSiteConfig.AcrUseManagedIdentityCreds; v != nil {
		result.UseManagedIdentityACR = *v
	}

	if v := functionAppSiteConfig.DefaultDocuments; v != nil {
		result.DefaultDocuments = *v
	}

	if v := functionAppSiteConfig.DetailedErrorLoggingEnabled; v != nil {
		result.DetailedErrorLogging = *v
	}

	if functionAppSiteConfig.RemoteDebuggingEnabled != nil {
		result.RemoteDebugging = *functionAppSiteConfig.RemoteDebuggingEnabled
	}

	if functionAppSiteConfig.RemoteDebuggingVersion != nil {
		// Note - This is sometimes returned in lower case, so we ToUpper it to avoid the need for a diff suppression
		result.RemoteDebuggingVersion = strings.ToUpper(*functionAppSiteConfig.RemoteDebuggingVersion)
	}

	if v := functionAppSiteConfig.FunctionsRuntimeScaleMonitoringEnabled; v != nil {
		result.RuntimeScaleMonitoring = *v
	}

	if functionAppSiteConfig.Cors != nil {
		corsSettings := functionAppSiteConfig.Cors
		cors := CorsSetting{}
		if corsSettings.SupportCredentials != nil {
			cors.SupportCredentials = *corsSettings.SupportCredentials
		}

		if corsSettings.AllowedOrigins != nil && len(*corsSettings.AllowedOrigins) != 0 {
			cors.AllowedOrigins = *corsSettings.AllowedOrigins
			result.Cors = []CorsSetting{cors}
		}
	}

	var appStack []ApplicationStackLinuxFunctionApp
	if functionAppSiteConfig.LinuxFxVersion != nil {
		decoded, err := DecodeFunctionAppLinuxFxVersion(*functionAppSiteConfig.LinuxFxVersion)
		if err != nil {
			return nil, fmt.Errorf("flattening site config: %s", err)
		}
		appStack = decoded
	}
	result.ApplicationStack = appStack

	if functionAppSiteConfig.LocalMySQLEnabled != nil {
		result.LocalMysql = *functionAppSiteConfig.LocalMySQLEnabled
	}

	if functionAppSiteConfig.VnetRouteAllEnabled != nil {
		result.VnetRouteAllEnabled = *functionAppSiteConfig.VnetRouteAllEnabled
	}

	return result, nil
}

func ParseWebJobsStorageString(input *string) (name, key string) {
	if input == nil {
		return
	}

	parts := strings.Split(*input, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "AccountName") {
			name = strings.TrimPrefix(part, "AccountName=")
		}
		if strings.HasPrefix(part, "AccountKey") {
			key = strings.TrimPrefix(part, "AccountKey=")
		}
	}

	return
}

func MergeUserAppSettings(systemSettings *[]web.NameValuePair, userSettings map[string]string) *[]web.NameValuePair {
	if len(userSettings) == 0 {
		return systemSettings
	}
	combined := *systemSettings
	for k, v := range userSettings {
		// Dedupe, explicit user settings take priority over enumerated, e.g. specifying KeyVault for `AzureWebJobsStorage`
		for i, x := range combined {
			if *x.Name == v {
				copy(combined[i:], combined[i+1:])
				combined[len(combined)-1] = web.NameValuePair{}
				combined = combined[:len(combined)-1]
			}
		}
		combined = append(combined, web.NameValuePair{
			Name:  utils.String(k),
			Value: utils.String(v),
		})
	}
	return &combined
}
