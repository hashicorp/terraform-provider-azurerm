package storage

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_storage_account_blob_container_sas": dataSourceStorageAccountBlobContainerSharedAccessSignature(),
		"azurerm_storage_account_sas":                dataSourceStorageAccountSharedAccessSignature(),
		"azurerm_storage_account":                    dataSourceStorageAccount(),
		"azurerm_storage_container":                  dataSourceStorageContainer(),
		"azurerm_storage_encryption_scope":           dataSourceStorageEncryptionScope(),
		"azurerm_storage_management_policy":          dataSourceStorageManagementPolicy(),
		"azurerm_storage_sync":                       dataSourceStorageSync(),
		"azurerm_storage_sync_group":                 dataSourceStorageSyncGroup(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_storage_account":                      resourceStorageAccount(),
		"azurerm_storage_account_customer_managed_key": resourceStorageAccountCustomerManagedKey(),
		"azurerm_storage_account_network_rules":        resourceStorageAccountNetworkRules(),
		"azurerm_storage_blob":                         resourceStorageBlob(),
		"azurerm_storage_container":                    resourceStorageContainer(),
		"azurerm_storage_encryption_scope":             resourceStorageEncryptionScope(),
		"azurerm_storage_data_lake_gen2_filesystem":    resourceStorageDataLakeGen2FileSystem(),
		"azurerm_storage_data_lake_gen2_path":          resourceStorageDataLakeGen2Path(),
		"azurerm_storage_management_policy":            resourceStorageManagementPolicy(),
		"azurerm_storage_queue":                        resourceStorageQueue(),
		"azurerm_storage_share":                        resourceStorageShare(),
		"azurerm_storage_share_directory":              resourceStorageShareDirectory(),
		"azurerm_storage_table":                        resourceStorageTable(),
		"azurerm_storage_table_entity":                 resourceStorageTableEntity(),
		"azurerm_storage_sync":                         resourceStorageSync(),
		"azurerm_storage_sync_group":                   resourceStorageSyncGroup(),
	}
}
