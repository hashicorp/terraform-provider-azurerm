package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type Registry struct {
	Server   string `tfschema:"registry_server_url"`
	UserName string `tfschema:"registry_username"`
	Password string `tfschema:"registry_password"`
}

func RegistrySchemaLinuxFunctionAppOnContainer() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"registry_server_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The login endpoint of the container Registry url",
				},

				"registry_username": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					RequiredWith: []string{"registry.0.registry_password"},
					Description:  "The username to use for this Container Registry.",
				},

				"registry_password": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					RequiredWith: []string{"registry.0.registry_username"},
					Description:  "The name of the Secret Reference containing the password value for this user on the Container Registry.",
				},
			},
		},
	}
}

type SiteConfigLinuxFunctionAppOnContainer struct {
	AppInsightsInstrumentationKey string `tfschema:"application_insights_key"` // App Insights Instrumentation Key
	AppInsightsConnectionString   string `tfschema:"application_insights_connection_string"`
	AppScaleLimit                 int64  `tfschema:"app_scale_limit"`
	ElasticInstanceMinimum        int64  `tfschema:"elastic_instance_minimum"`
	LinuxFxVersion                string `tfschema:"linux_fx_version"`
}

func SiteConfigSchemaLinuxFunctionAppOnContainer() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"application_insights_key": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Instrumentation Key for connecting the Linux Function App to Application Insights.",
				},

				"application_insights_connection_string": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Connection String for linking the Linux Function App to Application Insights.",
				},

				"app_scale_limit": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Default:     10,
					Description: "The number of workers this function app can scale out to. Only applicable to apps on the Consumption and Premium plan.",
				},

				"elastic_instance_minimum": {
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					Description: "The number of minimum instances for this Linux Function App. Only affects apps on Elastic Premium plans.",
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

func ExpandSiteConfigLinuxFunctionAppOnContainer(siteConfig []SiteConfigLinuxFunctionAppOnContainer, existing *webapps.SiteConfig, metadata sdk.ResourceMetaData, registry Registry, version string, storageString string) *webapps.SiteConfig {
	if len(siteConfig) == 0 {
		return nil
	}

	expanded := &webapps.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	appSettings := make([]webapps.NameValuePair, 0)

	appSettings = updateOrAppendAppSettings(appSettings, "FUNCTIONS_EXTENSION_VERSION", version, false)
	appSettings = updateOrAppendAppSettings(appSettings, "AzureWebJobsStorage", storageString, false)

	linuxFunctionOnContainerSiteConfig := siteConfig[0]

	if metadata.ResourceData.HasChange("site_config.0.elastic_instance_minimum") {
		expanded.MinimumElasticInstanceCount = pointer.To(linuxFunctionOnContainerSiteConfig.ElasticInstanceMinimum)
	}

	if metadata.ResourceData.HasChange("site_config.0.app_scale_limit") {
		expanded.FunctionAppScaleLimit = pointer.To(linuxFunctionOnContainerSiteConfig.AppScaleLimit)
	}

	if linuxFunctionOnContainerSiteConfig.AppInsightsConnectionString == "" {
		appSettings = updateOrAppendAppSettings(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", linuxFunctionOnContainerSiteConfig.AppInsightsConnectionString, true)
	} else {
		appSettings = updateOrAppendAppSettings(appSettings, "APPLICATIONINSIGHTS_CONNECTION_STRING", linuxFunctionOnContainerSiteConfig.AppInsightsConnectionString, false)
	}

	// update docker related settings
	appSettings = updateOrAppendAppSettings(appSettings, "DOCKER_REGISTRY_SERVER_URL", registry.Server, false)
	if registry.UserName != "" {
		appSettings = updateOrAppendAppSettings(appSettings, "DOCKER_REGISTRY_SERVER_USERNAME", registry.UserName, false)
	}
	if registry.Password != "" {
		appSettings = updateOrAppendAppSettings(appSettings, "DOCKER_REGISTRY_SERVER_PASSWORD", registry.Password, false)
	}

	expanded.AppSettings = &appSettings
	return expanded
}

func FlattenSiteConfigLinuxFunctionAppOnContainer(functionAppOnContainer *webapps.SiteConfig) *SiteConfigLinuxFunctionAppOnContainer {
	result := &SiteConfigLinuxFunctionAppOnContainer{
		ElasticInstanceMinimum: pointer.From(functionAppOnContainer.MinimumElasticInstanceCount),
		AppScaleLimit:          pointer.From(functionAppOnContainer.FunctionAppScaleLimit),
		LinuxFxVersion:         pointer.From(functionAppOnContainer.LinuxFxVersion),
	}

	return result
}

func EncodeLinuxFunctionAppOnContainerRegistryImage(input []Registry, image string) *string {
	if len(input) == 0 {
		return utils.String("")
	}
	dockerUrl := input[0].Server
	httpPrefixes := []string{"https://", "http://"}
	for _, prefix := range httpPrefixes {
		dockerUrl = strings.TrimPrefix(dockerUrl, prefix)
	}
	return utils.String(fmt.Sprintf("DOCKER|%s/%s", dockerUrl, image))
}

func DecodeLinuxFunctionAppOnContainerRegistryImage(input *string, partial *webapps.StringDictionary) (string, []Registry, error) {
	containerRegistryList := make([]Registry, 0)
	if input == nil {
		return "", containerRegistryList, nil
	}
	parts := strings.Split(*input, "|")
	value := parts[1]
	if len(parts) != 2 || parts[0] != "DOCKER" {
		return "", []Registry{}, fmt.Errorf("unrecognised LinuxFxVersion format received, got %s", *input)
	}

	// mcr.microsoft.com/dockerimage:tag
	registryUrl := value[:strings.Index(parts[1], "/")]

	containerRegistry := Registry{
		Server: registryUrl,
	}
	if partial.Properties != nil {
		for k, v := range *partial.Properties {
			if k == "DOCKER_REGISTRY_SERVER_URL" {
				containerRegistry.Server = v
			}
			if k == "DOCKER_REGISTRY_SERVER_USERNAME" {
				containerRegistry.UserName = v
			}
			if k == "DOCKER_REGISTRY_SERVER_PASSWORD" {

				containerRegistry.Password = v
			}
		}
		containerRegistryList = append(containerRegistryList, containerRegistry)
	}

	return value[strings.Index(parts[1], "/")+1:], containerRegistryList, nil
}
