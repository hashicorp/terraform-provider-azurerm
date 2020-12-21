package netapp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "NetApp"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"NetApp",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_netapp_account":  dataSourceNetAppAccount(),
		"azurerm_netapp_pool":     dataSourceNetAppPool(),
		"azurerm_netapp_volume":   dataSourceNetAppVolume(),
		"azurerm_netapp_snapshot": dataSourceNetAppSnapshot(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_netapp_account":  resourceNetAppAccount(),
		"azurerm_netapp_pool":     resourceNetAppPool(),
		"azurerm_netapp_volume":   resourceNetAppVolume(),
		"azurerm_netapp_snapshot": resourceNetAppSnapshot(),
	}
}
