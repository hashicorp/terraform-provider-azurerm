// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package digitaltwins

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/digital-twins"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Digital Twins"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Digital Twins",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_digital_twins_instance": dataSourceDigitalTwinsInstance(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_digital_twins_instance":            resourceDigitalTwinsInstance(),
		"azurerm_digital_twins_endpoint_eventgrid":  resourceDigitalTwinsEndpointEventGrid(),
		"azurerm_digital_twins_endpoint_eventhub":   resourceDigitalTwinsEndpointEventHub(),
		"azurerm_digital_twins_endpoint_servicebus": resourceDigitalTwinsEndpointServiceBus(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		TimeSeriesDatabaseConnectionResource{},
	}
}
