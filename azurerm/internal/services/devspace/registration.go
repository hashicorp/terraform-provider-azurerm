package devspace

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

// TODO: this can be moved into Container

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "DevSpaces"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"DevSpace",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	if features.ThreePointOh() {
		return map[string]*pluginsdk.Resource{}
	}

	// TODO: remove this entire package in 3.0
	return map[string]*pluginsdk.Resource{
		"azurerm_devspace_controller": resourceDevSpaceController(),
	}
}
