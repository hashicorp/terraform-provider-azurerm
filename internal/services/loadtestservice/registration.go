// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadtestservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/load-test"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Load Test",
	}
}

func (r Registration) Name() string {
	return "LoadTestService"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		LoadTestDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		LoadTestResource{},
	}
}
