// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package extendedlocation

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/extended-location"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Extended Location",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		CustomLocationResource{},
	}
}

func (r Registration) Name() string {
	return "ExtendedLocation"
}
