// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedidentity

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct {
	autoRegistration
}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/authorization"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return r.autoRegistration.Name()
}

func (r Registration) DataSources() []sdk.DataSource {
	dataSources := []sdk.DataSource{}
	dataSources = append(dataSources, r.autoRegistration.DataSources()...)
	return dataSources
}

func (r Registration) Resources() []sdk.Resource {
	resources := []sdk.Resource{
		FederatedIdentityCredentialResource{},
	}
	resources = append(resources, r.autoRegistration.Resources()...)
	return resources
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return r.autoRegistration.WebsiteCategories()
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_user_assigned_identity": dataSourceArmUserAssignedIdentity(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}
