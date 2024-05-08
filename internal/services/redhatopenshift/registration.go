// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Red Hat OpenShift"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/redhatopenshift"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Red Hat OpenShift",
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		RedHatOpenShiftCluster{},
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}
