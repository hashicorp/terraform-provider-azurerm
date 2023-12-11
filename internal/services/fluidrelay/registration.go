// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fluidrelay

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/fluid-relay"
}

func (r Registration) Name() string {
	return "Fluid Relay"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		Server{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Fluid Relay",
	}
}

var (
	_ sdk.TypedServiceRegistration = (*Registration)(nil)
)
