// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct {
}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/authorization"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Authorization"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Authorization",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_client_config":   dataSourceArmClientConfig(),
		"azurerm_role_definition": dataSourceArmRoleDefinition(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_role_assignment": resourceArmRoleAssignment(),
		"azurerm_role_definition": resourceArmRoleDefinition(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	resources := []sdk.Resource{
		PimActiveRoleAssignmentResource{},
		PimEligibleRoleAssignmentResource{},
		RoleAssignmentMarketplaceResource{},
	}
	return resources
}
