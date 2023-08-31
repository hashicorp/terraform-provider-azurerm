// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/redis"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Redis"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Redis",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_redis_cache": dataSourceRedisCache(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_redis_cache":         resourceRedisCache(),
		"azurerm_redis_firewall_rule": resourceRedisFirewallRule(),
		"azurerm_redis_linked_server": resourceRedisLinkedServer(),
	}
}
