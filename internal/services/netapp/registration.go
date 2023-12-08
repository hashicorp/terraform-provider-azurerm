// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/netapp"
}

func (r Registration) Name() string {
	return "NetApp"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"NetApp",
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_netapp_account":         dataSourceNetAppAccount(),
		"azurerm_netapp_pool":            dataSourceNetAppPool(),
		"azurerm_netapp_volume":          dataSourceNetAppVolume(),
		"azurerm_netapp_snapshot":        dataSourceNetAppSnapshot(),
		"azurerm_netapp_snapshot_policy": dataSourceNetAppSnapshotPolicy(),
	}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_netapp_account":         resourceNetAppAccount(),
		"azurerm_netapp_pool":            resourceNetAppPool(),
		"azurerm_netapp_volume":          resourceNetAppVolume(),
		"azurerm_netapp_snapshot":        resourceNetAppSnapshot(),
		"azurerm_netapp_snapshot_policy": resourceNetAppSnapshotPolicy(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		NetAppVolumeGroupSapHanaDataSource{},
		NetAppVolumeQuotaRuleDataSource{},
	}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		NetAppVolumeGroupSapHanaResource{},
		NetAppVolumeQuotaRuleResource{},
	}
}
