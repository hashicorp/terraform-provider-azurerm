// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redismanaged

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) WebsiteCategories() []string {
	return nil
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
	}
}
