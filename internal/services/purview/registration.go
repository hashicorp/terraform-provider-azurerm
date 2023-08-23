// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package purview

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/purview"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Purview"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Purview",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_purview_account": resourcePurviewAccount(),
	}
}
