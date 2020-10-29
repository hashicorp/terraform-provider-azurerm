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
		"azurerm_storage_account_blob_container_sas": dataSourceArmStorageAccountBlobContainerSharedAccessSignature(),
		"azurerm_storage_account_sas":                dataSourceArmStorageAccountSharedAccessSignature(),
		"azurerm_storage_account":                    dataSourceArmStorageAccount(),
		"azurerm_storage_container":                  dataSourceArmStorageContainer(),
		"azurerm_storage_management_policy":          dataSourceArmStorageManagementPolicy(),
		"azurerm_storage_sync":                       dataSourceArmStorageSync(),
		"azurerm_storage_sync_group":                 dataSourceArmStorageSyncGroup(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_storage_account":                      resourceArmStorageAccount(),
		"azurerm_storage_account_customer_managed_key": resourceArmStorageAccountCustomerManagedKey(),
		"azurerm_storage_account_network_rules":        resourceArmStorageAccountNetworkRules(),
		"azurerm_storage_blob":                         resourceArmStorageBlob(),
		"azurerm_storage_container":                    resourceArmStorageContainer(),
		"azurerm_storage_data_lake_gen2_filesystem":    resourceArmStorageDataLakeGen2FileSystem(),
		"azurerm_storage_management_policy":            resourceArmStorageManagementPolicy(),
		"azurerm_storage_queue":                        resourceArmStorageQueue(),
		"azurerm_storage_share":                        resourceArmStorageShare(),
		"azurerm_storage_share_directory":              resourceArmStorageShareDirectory(),
		"azurerm_storage_table":                        resourceArmStorageTable(),
		"azurerm_storage_table_entity":                 resourceArmStorageTableEntity(),
		"azurerm_storage_sync":                         resourceArmStorageSync(),
		"azurerm_storage_sync_cloud_endpoint":          resourceArmStorageSyncCloudEndpoint(),
		"azurerm_storage_sync_group":                   resourceArmStorageSyncGroup(),
	}
}
