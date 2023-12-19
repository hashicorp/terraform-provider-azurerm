---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account"
description: |-
  Gets information about an existing Storage Account.

---

# Data Source: azurerm_storage_account

Use this data source to access information about an existing Storage Account.

## Example Usage

```hcl
data "azurerm_storage_account" "example" {
  name                = "packerimages"
  resource_group_name = "packer-storage"
}

output "storage_account_tier" {
  value = data.azurerm_storage_account.example.account_tier
}
```

## Argument Reference

* `name` - Specifies the name of the Storage Account
* `resource_group_name` - Specifies the name of the resource group the Storage Account is located in.

## Attributes Reference

* `id` - The ID of the Storage Account.

* `location` - The Azure location where the Storage Account exists

* `identity` - An `identity` block as documented below.

* `account_kind` - The Kind of account.

* `account_tier` - The Tier of this storage account.

* `account_replication_type` - The type of replication used for this storage account.

* `access_tier` - The access tier for `BlobStorage` accounts.

* `enable_https_traffic_only` - Is traffic only allowed via HTTPS? See [here](https://docs.microsoft.com/azure/storage/storage-require-secure-transfer/)
    for more information.

* `min_tls_version` - The minimum supported TLS version for this storage account.

* `allow_nested_items_to_be_public` - Can nested items in the storage account opt into allowing public access?

* `is_hns_enabled` - Is Hierarchical Namespace enabled?

* `nfsv3_enabled` - Is NFSv3 protocol enabled?

* `custom_domain` - A `custom_domain` block as documented below.

* `tags` - A mapping of tags to assigned to the resource.

* `primary_location` - The primary location of the Storage Account.

* `secondary_location` - The secondary location of the Storage Account.

* `primary_blob_endpoint` - The endpoint URL for blob storage in the primary location.

* `primary_blob_host` - The hostname with port if applicable for blob storage in the primary location.

* `primary_blob_internet_endpoint` - The internet routing endpoint URL for blob storage in the primary location.

* `primary_blob_internet_host` - The internet routing hostname with port if applicable for blob storage in the primary location.

* `primary_blob_microsoft_endpoint` - The microsoft routing endpoint URL for blob storage in the primary location.

* `primary_blob_microsoft_host` - The microsoft routing hostname with port if applicable for blob storage in the primary location.

* `secondary_blob_endpoint` - The endpoint URL for blob storage in the secondary location.

* `secondary_blob_host` - The hostname with port if applicable for blob storage in the secondary location.

* `secondary_blob_internet_endpoint` - The internet routing endpoint URL for blob storage in the secondary location.

* `secondary_blob_internet_host` - The internet routing hostname with port if applicable for blob storage in the secondary location.

* `secondary_blob_microsoft_endpoint` - The microsoft routing endpoint URL for blob storage in the secondary location.

* `secondary_blob_microsoft_host` - The microsoft routing hostname with port if applicable for blob storage in the secondary location.

* `primary_queue_endpoint` - The endpoint URL for queue storage in the primary location.

* `primary_queue_host` - The hostname with port if applicable for queue storage in the primary location.

* `primary_queue_microsoft_endpoint` - The microsoft routing endpoint URL for queue storage in the primary location.

* `primary_queue_microsoft_host` - The microsoft routing hostname with port if applicable for queue storage in the primary location.

* `secondary_queue_endpoint` - The endpoint URL for queue storage in the secondary location.

* `secondary_queue_host` - The hostname with port if applicable for queue storage in the secondary location.

* `secondary_queue_microsoft_endpoint` - The microsoft routing endpoint URL for queue storage in the secondary location.

* `secondary_queue_microsoft_host` - The microsoft routing hostname with port if applicable for queue storage in the secondary location.

* `primary_table_endpoint` - The endpoint URL for table storage in the primary location.

* `primary_table_host` - The hostname with port if applicable for table storage in the primary location.

* `primary_table_microsoft_endpoint` - The microsoft routing endpoint URL for table storage in the primary location.

* `primary_table_microsoft_host` - The microsoft routing hostname with port if applicable for table storage in the primary location.

* `secondary_table_endpoint` - The endpoint URL for table storage in the secondary location.

* `secondary_table_host` - The hostname with port if applicable for table storage in the secondary location.

* `secondary_table_microsoft_endpoint` - The microsoft routing endpoint URL for table storage in the secondary location.

* `secondary_table_microsoft_host` - The microsoft routing hostname with port if applicable for table storage in the secondary location.

* `primary_file_endpoint` - The endpoint URL for file storage in the primary location.

* `primary_file_host` - The hostname with port if applicable for file storage in the primary location.

* `primary_file_internet_endpoint` - The internet routing endpoint URL for file storage in the primary location.

* `primary_file_internet_host` - The internet routing hostname with port if applicable for file storage in the primary location.

* `primary_file_microsoft_endpoint` - The microsoft routing endpoint URL for file storage in the primary location.

* `primary_file_microsoft_host` - The microsoft routing hostname with port if applicable for file storage in the primary location.

* `secondary_file_endpoint` - The endpoint URL for file storage in the secondary location.

* `secondary_file_host` - The hostname with port if applicable for file storage in the secondary location.

* `secondary_file_internet_endpoint` - The internet routing endpoint URL for file storage in the secondary location.

* `secondary_file_internet_host` - The internet routing hostname with port if applicable for file storage in the secondary location.

* `secondary_file_microsoft_endpoint` - The microsoft routing endpoint URL for file storage in the secondary location.

* `secondary_file_microsoft_host` - The microsoft routing hostname with port if applicable for file storage in the secondary location.

* `primary_dfs_endpoint` - The endpoint URL for DFS storage in the primary location.

* `primary_dfs_host` - The hostname with port if applicable for DFS storage in the primary location.

* `primary_dfs_internet_endpoint` - The internet routing endpoint URL for DFS storage in the primary location.

* `primary_dfs_internet_host` - The internet routing hostname with port if applicable for DFS storage in the primary location.

* `primary_dfs_microsoft_endpoint` - The microsoft routing endpoint URL for DFS storage in the primary location.

* `primary_dfs_microsoft_host` - The microsoft routing hostname with port if applicable for DFS storage in the primary location.

* `secondary_dfs_endpoint` - The endpoint URL for DFS storage in the secondary location.

* `secondary_dfs_host` - The hostname with port if applicable for DFS storage in the secondary location.

* `secondary_dfs_internet_endpoint` - The internet routing endpoint URL for DFS storage in the secondary location.

* `secondary_dfs_internet_host` - The internet routing hostname with port if applicable for DFS storage in the secondary location.

* `secondary_dfs_microsoft_endpoint` - The microsoft routing endpoint URL for DFS storage in the secondary location.

* `secondary_dfs_microsoft_host` - The microsoft routing hostname with port if applicable for DFS storage in the secondary location.

* `primary_web_endpoint` - The endpoint URL for web storage in the primary location.

* `primary_web_host` - The hostname with port if applicable for web storage in the primary location.

* `primary_web_internet_endpoint` - The internet routing endpoint URL for web storage in the primary location.

* `primary_web_internet_host` - The internet routing hostname with port if applicable for web storage in the primary location.

* `primary_web_microsoft_endpoint` - The microsoft routing endpoint URL for web storage in the primary location.

* `primary_web_microsoft_host` - The microsoft routing hostname with port if applicable for web storage in the primary location.

* `secondary_web_endpoint` - The endpoint URL for web storage in the secondary location.

* `secondary_web_host` - The hostname with port if applicable for web storage in the secondary location.

* `secondary_web_internet_endpoint` - The internet routing endpoint URL for web storage in the secondary location.

* `secondary_web_internet_host` - The internet routing hostname with port if applicable for web storage in the secondary location.

* `secondary_web_microsoft_endpoint` - The microsoft routing endpoint URL for web storage in the secondary location.

* `secondary_web_microsoft_host` - The microsoft routing hostname with port if applicable for web storage in the secondary location.

* `primary_access_key` - The primary access key for the Storage Account.

* `secondary_access_key` - The secondary access key for the Storage Account.

* `primary_connection_string` - The connection string associated with the primary location

* `secondary_connection_string` - The connection string associated with the secondary location

* `primary_blob_connection_string` - The connection string associated with the primary blob location

* `secondary_blob_connection_string` - The connection string associated with the secondary blob location

~> **Note:** If there's a Write Lock on the Storage Account, or the account doesn't have permission then these fields will have an empty value [due to a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/6363)

* `queue_encryption_key_type` - The encryption key type of the queue.

* `table_encryption_key_type` - The encryption key type of the table.

* `infrastructure_encryption_enabled` - Is infrastructure encryption enabled? See [here](https://docs.microsoft.com/azure/storage/common/infrastructure-encryption-enable/)
    for more information.

* `azure_files_authentication` - A `azure_files_authentication` block as documented below.

---

* `custom_domain` supports the following:

* `name` - The Custom Domain Name used for the Storage Account.

---

`identity` supports the following:

* `type` - The type of Managed Service Identity that is configured on this Storage Account

* `identity_ids` - A list of User Assigned Managed Identity IDs assigned with the Identity of this Storage Account.

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Storage Account.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Storage Account.

---

`azure_files_authentication` supports the following:

* `directory_type` - The directory service used for this Storage Account.

* `active_directory` - An `active_directory` block as documented below.

---

`active_directory` supports the following:

* `domain_name` - The primary domain that the AD DNS server is authoritative for.

* `netbios_domain_name` - The NetBIOS domain name.

* `forest_name` - The name of the Active Directory forest.

* `domain_guid` - The domain GUID.

* `domain_sid` - The domain security identifier.

* `storage_sid` - The security identifier for Azure Storage.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account.
