// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/redis"
}

func (r Registration) WebsiteCategories() []string {
	return []string{"Managed Redis"}
}

func (r Registration) Name() string {
	return "Managed Redis"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ManagedRedisDatabaseDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ManagedRedisClusterResource{},
		ManagedRedisDatabaseResource{},
		ManagedRedisDatabaseGeoReplicationResource{},
	}
}
