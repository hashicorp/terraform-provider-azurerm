package databoxedge

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Databox Edge"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Databox Edge",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_databox_edge_device": dataSourceDataboxEdgeDevice(),
		"azurerm_databox_edge_order":  dataSourceDataboxEdgeOrder(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_databox_edge_device": resourceDataboxEdgeDevice(),
		"azurerm_databox_edge_order":  resourceDataboxEdgeOrder(),
	}
}
