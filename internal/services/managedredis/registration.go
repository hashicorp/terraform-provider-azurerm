// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

//_ sdk.FrameworkServiceRegistration             = Registration{}
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
		ManagedRedisDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ManagedRedisGeoReplicationResource{},
		ManagedRedisResource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		newManagedRedisFlushDatabasesAction,
	}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
