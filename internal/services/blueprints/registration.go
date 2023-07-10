// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package blueprints

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/blueprints"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Blueprints"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Blueprints",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_blueprint_definition":        dataSourceBlueprintDefinition(),
		"azurerm_blueprint_published_version": dataSourceBlueprintPublishedVersion(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_blueprint_assignment": resourceBlueprintAssignment(),
	}
}
