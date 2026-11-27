// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_blob_container_sas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_local_user"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_network_rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_queue_properties"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_sas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_static_website"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_account_table_properties"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_blob"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_blob_inventory_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_container"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_container_immutability_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_containers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_data_lake_gen2_filesystem"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_data_lake_gen2_path"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_discovery_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_encryption_scope"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_management_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_object_replication"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_queue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_share_directory"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_share_file"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_sync"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_sync_cloud_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_sync_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_sync_server_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_table_entities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/resources/storage_table_entity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/storage"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Storage"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Storage",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_storage_account":                    storage_account.RegisterDataSource(),
		"azurerm_storage_account_blob_container_sas": storage_account_blob_container_sas.RegisterDataSource(),
		"azurerm_storage_account_sas":                storage_account_sas.RegisterDataSource(),
		"azurerm_storage_blob":                       storage_blob.RegisterDataSource(),
		"azurerm_storage_container":                  storage_container.RegisterDataSource(),
		"azurerm_storage_encryption_scope":           storage_encryption_scope.RegisterDataSource(),
		"azurerm_storage_management_policy":          storage_management_policy.RegisterDataSource(),
		"azurerm_storage_queue":                      storage_queue.RegisterDataSource(),
		"azurerm_storage_share":                      storage_share.RegisterDataSource(),
		"azurerm_storage_sync":                       storage_sync.RegisterDataSource(),
		"azurerm_storage_sync_group":                 storage_sync_group.RegisterDataSource(),
		"azurerm_storage_table_entity":               storage_table_entity.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_storage_account":                      storage_account.RegisterResource(),
		"azurerm_storage_account_customer_managed_key": storage_account_customer_managed_key.RegisterResource(),
		"azurerm_storage_account_network_rules":        storage_account_network_rules.RegisterResource(),
		"azurerm_storage_blob":                         storage_blob.RegisterResource(),
		"azurerm_storage_blob_inventory_policy":        storage_blob_inventory_policy.RegisterResource(),
		"azurerm_storage_container":                    storage_container.RegisterResource(),
		"azurerm_storage_data_lake_gen2_filesystem":    storage_data_lake_gen2_filesystem.RegisterResource(),
		"azurerm_storage_data_lake_gen2_path":          storage_data_lake_gen2_path.RegisterResource(),
		"azurerm_storage_encryption_scope":             storage_encryption_scope.RegisterResource(),
		"azurerm_storage_management_policy":            storage_management_policy.RegisterResource(),
		"azurerm_storage_object_replication":           storage_object_replication.RegisterResource(),
		"azurerm_storage_queue":                        storage_queue.RegisterResource(),
		"azurerm_storage_share":                        storage_share.RegisterResource(),
		"azurerm_storage_share_directory":              storage_share_directory.RegisterResource(),
		"azurerm_storage_share_file":                   storage_share_file.RegisterResource(),
		"azurerm_storage_sync":                         storage_sync.RegisterResource(),
		"azurerm_storage_sync_cloud_endpoint":          storage_sync_cloud_endpoint.RegisterResource(),
		"azurerm_storage_sync_group":                   storage_sync_group.RegisterResource(),
		"azurerm_storage_table":                        storage_table.RegisterResource(),
		"azurerm_storage_table_entity":                 storage_table_entity.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		storage_account_static_website.DataSource{},
		storage_containers.DataSource{},
		storage_table.DataSource{},
		storage_table_entities.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		storage_account_queue_properties.Resource{},
		storage_account_static_website.Resource{},
		storage_account_table_properties.Resource{},
		storage_account_local_user.Resource{},
		storage_container_immutability_policy.Resource{},
		storage_discovery_workspace.Resource{},
		storage_sync_server_endpoint.Resource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
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
		storage_account_customer_managed_key.ListResource{},
		storage_account.ListResource{},
		storage_account_network_rules.ListResource{},
		storage_container.ListResource{},
		storage_sync.ListResource{},
		storage_sync_server_endpoint.ListResource{},
	}
}
