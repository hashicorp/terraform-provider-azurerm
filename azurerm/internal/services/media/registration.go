package media

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Media"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Media",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_media_services_account":   resourceMediaServicesAccount(),
		"azurerm_media_asset":              resourceMediaAsset(),
		"azurerm_media_job":                resourceMediaJob(),
		"azurerm_media_streaming_endpoint": resourceMediaStreamingEndpoint(),
		"azurerm_media_transform":          resourceMediaTransform(),
		"azurerm_media_streaming_locator":  resourceMediaStreamingLocator(),
		"azurerm_media_content_key_policy": resourceMediaContentKeyPolicy(),
		"azurerm_media_streaming_policy":   resourceMediaStreamingPolicy(),
		"azurerm_media_live_event":         resourceMediaLiveEvent(),
		"azurerm_media_live_event_output":  resourceMediaLiveOutput(),
		"azurerm_media_asset_filter":       resourceMediaAssetFilter(),
	}
}
