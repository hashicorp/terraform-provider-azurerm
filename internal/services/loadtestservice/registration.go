// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadtestservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct {
	autoRegistration
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/load-test"
}

func (r Registration) WebsiteCategories() []string {
	return r.autoRegistration.WebsiteCategories()
}

func (r Registration) Name() string {
	return r.autoRegistration.Name()
}

func (r Registration) DataSources() []sdk.DataSource {
	return r.autoRegistration.DataSources()
}

func (r Registration) Resources() []sdk.Resource {
	return r.autoRegistration.Resources()
}
