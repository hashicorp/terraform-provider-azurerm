// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/access_policy_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/flush_databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/geo_replication"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/managed_redis"
)

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

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
		access_policy_assignment.DataSource{},
		managed_redis.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		access_policy_assignment.Resource{},
		geo_replication.Resource{},
		managed_redis.Resource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		flush_databases.Action,
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
