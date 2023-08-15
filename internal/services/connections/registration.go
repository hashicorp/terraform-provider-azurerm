// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connections

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.UntypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/connections"
}

func (r Registration) Name() string {
	return "Connections"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Connections",
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_managed_api": dataSourceManagedApi(),
	}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_api_connection": resourceConnection(),
	}
}
