// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/orbital"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Orbital",
	}
}

func (r Registration) Name() string {
	return "Orbital"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		SpacecraftResource{},
		ContactProfileResource{},
		ContactResource{},
	}
}

var _ sdk.TypedServiceRegistration = Registration{}
