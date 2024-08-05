// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arcresourcebridge

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/arc-resource-bridge"
}

func (r Registration) Name() string {
	return "Arc Resource Bridge"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ArcResourceBridgeApplianceDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ArcResourceBridgeApplianceResource{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Arc Resource Bridge",
	}
}
