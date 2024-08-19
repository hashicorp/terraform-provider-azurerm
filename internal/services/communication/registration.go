// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/communication"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Communication"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Communication",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		CommunicationServiceDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		EmailCommunicationServiceDomainResource{},
		EmailCommunicationServiceResource{},
		EmailDomainAssociationResource{},
		CommunicationServiceResource{},
	}
}
