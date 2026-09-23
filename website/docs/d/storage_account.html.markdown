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

## Arguments Reference

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

* `blob_properties` - A `blob_properties` block as documented below.

* `cross_tenant_replication_enabled` - Whether cross Tenant replication is enabled.

* `dns_endpoint_type` - Which DNS endpoint type is used - either `Standard` or `AzureDnsZone`.

* `https_traffic_only_enabled` - Is traffic only allowed via HTTPS? See [here](https://docs.microsoft.com/azure/storage/storage-require-secure-transfer/) for more information.

* `min_tls_version` - The minimum supported TLS version for this storage account.

* `allow_nested_items_to_be_public` - Can nested items in the storage account opt into allowing public access?

* `allowed_copy_scope` - The permitted scope for copy operations between storage accounts.

* `is_hns_enabled` - Is Hierarchical Namespace enabled?

* `large_file_share_enabled` - Whether Large File Shares are enabled.

* `local_user_enabled` - Whether Local User Authentication is enabled.

* `network_rules` - A `network_rules` block as documented below.

* `nfsv3_enabled` - Is NFSv3 protocol enabled?

* `custom_domain` - A `custom_domain` block as documented below.

* `customer_managed_key` - A `customer_managed_key` block as documented below.

* `default_to_oauth_authentication` - Whether Azure Active Directory authorization is used by default when accessing the Storage Account.

* `edge_zone` - The Edge Zone within the Azure Region where this Storage Account exists.

* `immutability_policy` - An `immutability_policy` block as documented below.

* `tags` - A mapping of tags to assigned to the resource.

* `primary_location` - The primary location of the Storage Account.

* `routing` - A `routing` block as documented below.

* `sas_policy` - A `sas_policy` block as documented below.

* `provisioned_billing_model_version` - The version of the `provisioned` billing model, for example when `account_kind` is `FileStorage`.

* `public_network_access` - The public network access setting of the Storage Account.

* `secondary_location` - The secondary location of the Storage Account.

* `sftp_enabled` - Whether SFTP is enabled for the storage account.

* `share_properties` - A `share_properties` block as documented below.

* `shared_access_key_enabled` - Whether the storage account permits requests to be authorized with the account access key via Shared Key.

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

A `blob_properties` block exports the following:

* `change_feed_enabled` - Whether the blob service properties for change feed events is enabled.

* `change_feed_retention_in_days` - The duration of change feed events retention in days.

* `container_delete_retention_policy` - A `container_delete_retention_policy` block as documented below.

* `cors_rule` - A `cors_rule` block as documented below.

* `default_service_version` - The API Version used by default for requests to the Data Plane API if an incoming request doesn't specify an API Version.

* `delete_retention_policy` - A `delete_retention_policy` block as documented below.

* `last_access_time_enabled` - Whether last access time based tracking is enabled.

* `restore_policy` - A `restore_policy` block as documented below.

* `versioning_enabled` - Whether versioning is enabled.

---

A `container_delete_retention_policy` block exports the following:

* `days` - The number of days that the container is retained.

---

A `cors_rule` block exports the following:

* `allowed_headers` - A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - A list of HTTP methods that are allowed to be executed by the origin.

* `allowed_origins` - A list of origin domains that are allowed by CORS.

* `exposed_headers` - A list of response headers that are exposed to CORS clients.

* `max_age_in_seconds` - The number of seconds the client should cache a preflight response.

---

* `custom_domain` supports the following:

* `name` - The Custom Domain Name used for the Storage Account.

---

A `customer_managed_key` block exports the following:

* `key_vault_key_id` - The ID of the Key Vault Key.

* `user_assigned_identity_id` - The ID of the user assigned identity.

---

A `delete_retention_policy` block exports the following:

* `days` - The number of days that the blob is retained.

* `permanent_delete_enabled` - Whether permanent deletion of the soft deleted blob versions and snapshots is allowed.

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

* `default_share_level_permission` - The default share level permissions applied to all users.

---

`active_directory` supports the following:

* `domain_name` - The primary domain that the AD DNS server is authoritative for.

* `netbios_domain_name` - The NetBIOS domain name.

* `forest_name` - The name of the Active Directory forest.

* `domain_guid` - The domain GUID.

* `domain_sid` - The domain security identifier.

* `storage_sid` - The security identifier for Azure Storage.

---

An `immutability_policy` block exports the following:

* `allow_protected_append_writes` - Whether new blocks can be written to an append blob while maintaining immutability protection and compliance.

* `period_since_creation_in_days` - The immutability period for the blobs in the container since the policy creation, in days.

* `state` - The mode of the policy.

---

A `network_rules` block exports the following:

* `bypass` - Whether traffic is bypassed for Logging/Metrics/AzureServices.

* `default_action` - The default action of allow or deny when no other rules match.

* `ip_rules` - The list of public IP or IP ranges in CIDR Format allowed to access the Storage Account.

* `private_link_access` - One or more `private_link_access` blocks as documented below.

* `virtual_network_subnet_ids` - The list of resource IDs for subnets allowed to access the Storage Account.

---

A `private_link_access` block exports the following:

* `endpoint_resource_id` - The ID of the Azure resource that is allowed access to the target storage account.

* `endpoint_tenant_id` - The tenant ID of the resource of the resource access rule granted access.

---

A `restore_policy` block exports the following:

* `days` - The number of days that the blob can be restored.

---

A `retention_policy` block exports the following:

* `days` - The number of days that the `azurerm_storage_share` is retained.

---

A `routing` block exports the following:

* `choice` - The kind of network routing opted by the user.

* `publish_internet_endpoints` - Whether internet routing storage endpoints are published.

* `publish_microsoft_endpoints` - Whether Microsoft routing storage endpoints are published.

---

A `sas_policy` block exports the following:

* `expiration_action` - The SAS expiration action.

* `expiration_period` - The SAS expiration period in the format of `DD.HH:MM:SS`.

---

A `share_properties` block exports the following:

* `cors_rule` - A `cors_rule` block as documented below.

* `retention_policy` - A `retention_policy` block as documented below.

* `smb` - A `smb` block as documented below.

---

A `smb` block exports the following:

* `authentication_types` - The set of SMB authentication methods.

* `channel_encryption_type` - The set of SMB channel encryption.

* `kerberos_ticket_encryption_type` - The set of Kerberos ticket encryption.

* `multichannel_enabled` - Whether multichannel is enabled.

* `versions` - The set of SMB protocol versions.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Storage` - 2025-08-01
