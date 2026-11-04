// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/actions/data_protection_backup_instance_protect"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/actions/data_protection_resource_guard_unlock_delete"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_instance_blob_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_instance_data_lake_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_instance_disk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_instance_kubernetes_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_instance_mysql_flexible_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_instance_postgresql_flexible_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_policy_blob_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_policy_cosmosdb_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_policy_data_lake_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_policy_disk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_policy_kubernetes_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_policy_mysql_flexible_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_policy_postgresql_flexible_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_vault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_backup_vault_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/resources/data_protection_resource_guard"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		data_protection_backup_policy_cosmosdb_account.ListResource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		data_protection_backup_instance_protect.RegisterAction,
		data_protection_resource_guard_unlock_delete.RegisterAction,
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

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/data-protection"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "DataProtection"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"DataProtection",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_protection_backup_vault": data_protection_backup_vault.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_protection_backup_instance_blob_storage": data_protection_backup_instance_blob_storage.RegisterResource(),
		"azurerm_data_protection_backup_instance_disk":         data_protection_backup_instance_disk.RegisterResource(),
		"azurerm_data_protection_backup_policy_blob_storage":   data_protection_backup_policy_blob_storage.RegisterResource(),
		"azurerm_data_protection_backup_policy_disk":           data_protection_backup_policy_disk.RegisterResource(),
		"azurerm_data_protection_backup_vault":                 data_protection_backup_vault.RegisterResource(),
		"azurerm_data_protection_resource_guard":               data_protection_resource_guard.RegisterResource(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		data_protection_backup_instance_data_lake_storage.Resource{},
		data_protection_backup_instance_kubernetes_cluster.Resource{},
		data_protection_backup_instance_mysql_flexible_server.Resource{},
		data_protection_backup_instance_postgresql_flexible_server.Resource{},
		data_protection_backup_policy_cosmosdb_account.Resource{},
		data_protection_backup_policy_data_lake_storage.Resource{},
		data_protection_backup_policy_kubernetes_cluster.Resource{},
		data_protection_backup_policy_mysql_flexible_server.Resource{},
		data_protection_backup_policy_postgresql_flexible_server.Resource{},
		data_protection_backup_vault_customer_managed_key.Resource{},
	}
}
