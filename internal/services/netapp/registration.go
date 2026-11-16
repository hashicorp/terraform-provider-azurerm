// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/actions/netapp_volume_bucket_credentials"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_account_encryption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_backup_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_backup_vault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_snapshot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_snapshot_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_volume"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_volume_bucket"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_volume_bucket_with_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_volume_group_oracle"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_volume_group_sap_hana"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/resources/netapp_volume_quota_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/netapp"
}

func (r Registration) Name() string {
	return "NetApp"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"NetApp",
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_netapp_account":         netapp_account.RegisterDataSource(),
		"azurerm_netapp_pool":            netapp_pool.RegisterDataSource(),
		"azurerm_netapp_snapshot":        netapp_snapshot.RegisterDataSource(),
		"azurerm_netapp_snapshot_policy": netapp_snapshot_policy.RegisterDataSource(),
		"azurerm_netapp_volume":          netapp_volume.RegisterDataSource(),
	}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_netapp_account":         netapp_account.RegisterResource(),
		"azurerm_netapp_pool":            netapp_pool.RegisterResource(),
		"azurerm_netapp_snapshot":        netapp_snapshot.RegisterResource(),
		"azurerm_netapp_snapshot_policy": netapp_snapshot_policy.RegisterResource(),
		"azurerm_netapp_volume":          netapp_volume.RegisterResource(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		netapp_account_encryption.DataSource{},
		netapp_backup_policy.DataSource{},
		netapp_backup_vault.DataSource{},
		netapp_volume_bucket.DataSource{},
		netapp_volume_bucket_with_server.DataSource{},
		netapp_volume_group_oracle.DataSource{},
		netapp_volume_group_sap_hana.DataSource{},
		netapp_volume_quota_rule.DataSource{},
	}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		netapp_account_encryption.Resource{},
		netapp_backup_policy.Resource{},
		netapp_backup_vault.Resource{},
		netapp_volume_bucket.Resource{},
		netapp_volume_bucket_with_server.Resource{},
		netapp_volume_group_oracle.Resource{},
		netapp_volume_group_sap_hana.Resource{},
		netapp_volume_quota_rule.Resource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		netapp_volume_bucket_credentials.RegisterAction,
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
	return []sdk.FrameworkListWrappedResource{
		netapp_volume_bucket.ListResource{},
	}
}
