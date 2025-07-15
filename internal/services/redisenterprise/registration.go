// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redisenterprise

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
	return "Redis Enterprise"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Redis Enterprise",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_redis_enterprise_database": dataSourceRedisEnterpriseDatabase(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_redis_enterprise_cluster":  resourceRedisEnterpriseCluster(),
		"azurerm_redis_enterprise_database": resourceRedisEnterpriseDatabase(),
	}
}
