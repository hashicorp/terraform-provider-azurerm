// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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
		"azurerm_storage_account_blob_container_sas": dataSourceStorageAccountBlobContainerSharedAccessSignature(),
		"azurerm_storage_account_sas":                dataSourceStorageAccountSharedAccessSignature(),
		"azurerm_storage_account":                    dataSourceStorageAccount(),
		"azurerm_storage_blob":                       dataSourceStorageBlob(),
		"azurerm_storage_container":                  dataSourceStorageContainer(),
		"azurerm_storage_encryption_scope":           dataSourceStorageEncryptionScope(),
		"azurerm_storage_management_policy":          dataSourceStorageManagementPolicy(),
		"azurerm_storage_queue":                      dataSourceStorageQueue(),
		"azurerm_storage_share":                      dataSourceStorageShare(),
		"azurerm_storage_sync":                       dataSourceStorageSync(),
		"azurerm_storage_sync_group":                 dataSourceStorageSyncGroup(),
		"azurerm_storage_table_entity":               dataSourceStorageTableEntity(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_storage_account":                      resourceStorageAccount(),
		"azurerm_storage_account_customer_managed_key": resourceStorageAccountCustomerManagedKey(),
		"azurerm_storage_account_network_rules":        resourceStorageAccountNetworkRules(),
		"azurerm_storage_blob":                         resourceStorageBlob(),
		"azurerm_storage_blob_inventory_policy":        resourceStorageBlobInventoryPolicy(),
		"azurerm_storage_container":                    resourceStorageContainer(),
		"azurerm_storage_encryption_scope":             resourceStorageEncryptionScope(),
		"azurerm_storage_data_lake_gen2_filesystem":    resourceStorageDataLakeGen2FileSystem(),
		"azurerm_storage_data_lake_gen2_path":          resourceStorageDataLakeGen2Path(),
		"azurerm_storage_management_policy":            resourceStorageManagementPolicy(),
		"azurerm_storage_object_replication":           resourceStorageObjectReplication(),
		"azurerm_storage_queue":                        resourceStorageQueue(),
		"azurerm_storage_share":                        resourceStorageShare(),
		"azurerm_storage_share_file":                   resourceStorageShareFile(),
		"azurerm_storage_share_directory":              resourceStorageShareDirectory(),
		"azurerm_storage_table":                        resourceStorageTable(),
		"azurerm_storage_table_entity":                 resourceStorageTableEntity(),
		"azurerm_storage_sync":                         resourceStorageSync(),
		"azurerm_storage_sync_cloud_endpoint":          resourceStorageSyncCloudEndpoint(),
		"azurerm_storage_sync_group":                   resourceStorageSyncGroup(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		storageTableDataSource{},
		storageTableEntitiesDataSource{},
		storageContainersDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		LocalUserResource{},
		StorageContainerImmutabilityPolicyResource{},
		SyncServerEndpointResource{},
	}
}
