package iotoperations

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.UntypedServiceRegistration = Registration{}
)

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the typed resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		InstanceResource{},
		BrokerResource{},
		BrokerAuthenticationResource{},
		BrokerAuthorizationResource{},
		BrokerListenerResource{},
		DataflowResource{},
		DataflowEndpointResource{},
		DataflowProfileResource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "IoT Operations"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"IoT Operations",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		// All IoT Operations resources are now typed and auto-registered through Resources() method
	}
}
