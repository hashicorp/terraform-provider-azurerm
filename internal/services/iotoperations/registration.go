// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotoperations

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

// Change this to typed registration
var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/iot-operations"
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

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		// Add typed data sources here when implemented
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		InstanceResource{},
		// Comment out others until they compile
		BrokerResource{},
		BrokerAuthenticationResource{},
		BrokerAuthorizationResource{},
		BrokerListenerResource{},
		DataflowResource{},
	    DataflowEndpointResource{},
		DataflowProfileResource{},
	}
}
