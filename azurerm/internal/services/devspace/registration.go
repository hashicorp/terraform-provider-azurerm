package devspace

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	if features.ThreePointOh() {
		return map[string]*schema.Resource{}
	}

	// TODO: remove this entire package in 3.0
	return map[string]*schema.Resource{
		"azurerm_devspace_controller": resourceDevSpaceController(),
	}
}
