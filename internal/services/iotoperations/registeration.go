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
		"azurerm_iotoperations_instance":              resourceInstance(),
		"azurerm_iotoperations_broker":                resourceBroker(),
		"azurerm_iotoperations_broker_authentication": resourceBrokerAuthentication(),
		"azurerm_iotoperations_broker_authorization":  resourceBrokerAuthorization(),
		"azurerm_iotoperations_broker_listener":       resourceBrokerListener(),
		"azurerm_iotoperations_dataflow":              resourceDataflow(),
		"azurerm_iotoperations_dataflow_endpoint":     resourceDataflowEndpoint(),
		"azurerm_iotoperations_dataflow_profile":      resourceDataflowProfile(),
	}
}
