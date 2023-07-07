// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package voiceservices

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration = Registration{}
)

// Name is the name of this Service
func (r Registration) Name() string {
	return "Voice Services"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Voice Services",
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		CommunicationsGatewayResource{},
		CommunicationsGatewayTestLineResource{},
	}
}
