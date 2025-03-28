// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devcenter

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) Name() string {
	return "Dev Center"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/dev-center"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Dev Center",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		DevCenterAttachedNetworkResource{},
		DevCenterGalleryResource{},
		DevCenterCatalogsResource{},
		DevCenterDevBoxDefinitionResource{},
		DevCenterEnvironmentTypeResource{},
		DevCenterNetworkConnectionResource{},
		DevCenterProjectPoolResource{},
		DevCenterProjectResource{},
		DevCenterProjectEnvironmentTypeResource{},
		DevCenterResource{},
	}
}
