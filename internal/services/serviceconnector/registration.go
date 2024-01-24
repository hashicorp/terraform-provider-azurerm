// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serviceconnector

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/service-connector"
}

func (r Registration) WebsiteCategories() []string {
	return []string{}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AppServiceConnectorResource{},
		SpringCloudConnectorResource{},
		FunctionAppConnectorResource{},
	}
}

func (r Registration) Name() string {
	return "ServiceConnector"
}
