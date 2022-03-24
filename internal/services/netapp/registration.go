package netapp

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/netapp"
}

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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_netapp_account":         dataSourceNetAppAccount(),
		"azurerm_netapp_pool":            dataSourceNetAppPool(),
		"azurerm_netapp_volume":          dataSourceNetAppVolume(),
		"azurerm_netapp_snapshot":        dataSourceNetAppSnapshot(),
		"azurerm_netapp_snapshot_policy": dataSourceNetAppSnapshotPolicy(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_netapp_account":         resourceNetAppAccount(),
		"azurerm_netapp_pool":            resourceNetAppPool(),
		"azurerm_netapp_volume":          resourceNetAppVolume(),
		"azurerm_netapp_snapshot":        resourceNetAppSnapshot(),
		"azurerm_netapp_snapshot_policy": resourceNetAppSnapshotPolicy(),
	}
}
