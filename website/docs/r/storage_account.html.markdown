---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account"
description: |-
  Manages a Azure Storage Account.
---

# azurerm_storage_account

Manages an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
```

## Example Usage with Network Rules

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "virtnetname"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "subnetname"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Sql", "Microsoft.Storage"]
}

resource "azurerm_storage_account" "example" {
  name                = "storageaccountname"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["100.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.example.id]
  }

  tags = {
    environment = "staging"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the storage account. Only lowercase Alphanumeric characters allowed. Changing this forces a new resource to be created. This must be unique across the entire Azure service, not just within the resource group.

* `resource_group_name` - (Required) The name of the resource group in which to create the storage account. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_kind` - (Optional) Defines the Kind of account. Valid options are `BlobStorage`, `BlockBlobStorage`, `FileStorage`, `Storage` and `StorageV2`. Defaults to `StorageV2`.

-> **Note:** Changing the `account_kind` value from `Storage` to `StorageV2` will not trigger a force new on the storage account, it will only upgrade the existing storage account from `Storage` to `StorageV2` keeping the existing storage account in place.

* `account_tier` - (Required) Defines the Tier to use for this storage account. Valid options are `Standard` and `Premium`. For `BlockBlobStorage` and `FileStorage` accounts only `Premium` is valid. Changing this forces a new resource to be created.

-> **Note:** Blobs with a tier of `Premium` are of account kind `StorageV2`.

* `account_replication_type` - (Required) Defines the type of replication to use for this storage account. Valid options are `LRS`, `GRS`, `RAGRS`, `ZRS`, `GZRS` and `RAGZRS`. Changing this forces a new resource to be created when types `LRS`, `GRS` and `RAGRS` are changed to `ZRS`, `GZRS` or `RAGZRS` and vice versa.

* `cross_tenant_replication_enabled` - (Optional) Should cross Tenant replication be enabled? Defaults to `false`.

* `access_tier` - (Optional) Defines the access tier for `BlobStorage`, `FileStorage` and `StorageV2` accounts. Valid options are `Hot`, `Cool`, `Cold` and `Premium`. Defaults to `Hot`.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Storage Account should exist. Changing this forces a new Storage Account to be created.

* `https_traffic_only_enabled` - (Optional) Boolean flag which forces HTTPS if enabled, see [here](https://docs.microsoft.com/azure/storage/storage-require-secure-transfer/) for more information. Defaults to `true`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the storage account. Possible values are `TLS1_0`, `TLS1_1`, and `TLS1_2`. Defaults to `TLS1_2` for new storage accounts.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

-> **Note:** At this time `min_tls_version` is only supported in the Public Cloud, China Cloud, and US Government Cloud.

* `allow_nested_items_to_be_public` - (Optional) Allow or disallow nested items within this Account to opt into being public. Defaults to `true`.

-> **Note:** At this time `allow_nested_items_to_be_public` is only supported in the Public Cloud, China Cloud, and US Government Cloud.

* `shared_access_key_enabled` - (Optional) Indicates whether the storage account permits requests to be authorized with the account access key via Shared Key. If false, then all requests, including shared access signatures, must be authorized with Azure Active Directory (Azure AD). Defaults to `true`.

~> **Note:** Terraform uses Shared Key Authorisation to provision Storage Containers, Blobs and other items - when Shared Key Access is disabled, you will need to enable [the `storage_use_azuread` flag in the Provider block](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#storage_use_azuread) to use Azure AD for authentication, however not all Azure Storage services support Active Directory authentication.

* `public_network_access_enabled` - (Optional) Whether the public network access is enabled? Defaults to `true`.

* `default_to_oauth_authentication` - (Optional) Default to Azure Active Directory authorization in the Azure portal when accessing the Storage Account. The default value is `false`

* `is_hns_enabled` - (Optional) Is Hierarchical Namespace enabled? This can be used with Azure Data Lake Storage Gen 2 ([see here for more information](https://docs.microsoft.com/azure/storage/blobs/data-lake-storage-quickstart-create-account/)). Changing this forces a new resource to be created.

-> **Note:** This can only be `true` when `account_tier` is `Standard` or when `account_tier` is `Premium` *and* `account_kind` is `BlockBlobStorage`

* `nfsv3_enabled` - (Optional) Is NFSv3 protocol enabled? Changing this forces a new resource to be created. Defaults to `false`.

-> **Note:** This can only be `true` when `account_tier` is `Standard` and `account_kind` is `StorageV2`, or `account_tier` is `Premium` and `account_kind` is `BlockBlobStorage`. Additionally, the `is_hns_enabled` is `true` and `account_replication_type` must be `LRS` or `RAGRS`.

* `custom_domain` - (Optional) A `custom_domain` block as documented below.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as documented below.

~> **Note:** It's possible to define a Customer Managed Key both within either the `customer_managed_key` block or by using the [`azurerm_storage_account_customer_managed_key`](storage_account_customer_managed_key.html) resource. However, it's not possible to use both methods to manage a Customer Managed Key for a Storage Account, since these will conflict. When using the `azurerm_storage_account_customer_managed_key` resource, you will need to use `ignore_changes` on the `customer_managed_key` block.

* `identity` - (Optional) An `identity` block as defined below.

* `blob_properties` - (Optional) A `blob_properties` block as defined below.

* `queue_properties` - (Optional) A `queue_properties` block as defined below.

~> **Note:** `queue_properties` can only be configured when `account_tier` is set to `Standard` and `account_kind` is set to either `Storage` or `StorageV2`.

* `static_website` - (Optional) A `static_website` block as defined below.

~> **Note:** `static_website` can only be set when the `account_kind` is set to `StorageV2` or `BlockBlobStorage`.

~> **Note:** If `static_website` is specified, the service will automatically create a `azurerm_storage_container` named `$web`.

* `share_properties` - (Optional) A `share_properties` block as defined below.

~> **Note:** `share_properties` can only be configured when either `account_tier` is `Standard` and `account_kind` is either `Storage` or `StorageV2` - or when `account_tier` is `Premium` and `account_kind` is `FileStorage`.

* `network_rules` - (Optional) A `network_rules` block as documented below.

* `large_file_share_enabled` - (Optional) Are Large File Shares Enabled? Defaults to `false`.

-> **Note:** Large File Shares are enabled by default when using an `account_kind` of `FileStorage`.

* `local_user_enabled` - (Optional) Is Local User Enabled? Defaults to `true`.

* `azure_files_authentication` - (Optional) A `azure_files_authentication` block as defined below.

* `routing` - (Optional) A `routing` block as defined below.

* `queue_encryption_key_type` - (Optional) The encryption type of the queue service. Possible values are `Service` and `Account`. Changing this forces a new resource to be created. Default value is `Service`.

* `table_encryption_key_type` - (Optional) The encryption type of the table service. Possible values are `Service` and `Account`. Changing this forces a new resource to be created. Default value is `Service`.

~> **Note:** `queue_encryption_key_type` and `table_encryption_key_type` cannot be set to `Account` when `account_kind` is set `Storage`

* `infrastructure_encryption_enabled` - (Optional) Is infrastructure encryption enabled? Changing this forces a new resource to be created. Defaults to `false`.

-> **Note:** This can only be `true` when `account_kind` is `StorageV2` or when `account_tier` is `Premium` *and* `account_kind` is one of `BlockBlobStorage` or `FileStorage`.

* `immutability_policy` - (Optional) An `immutability_policy` block as defined below. Changing this forces a new resource to be created.

* `sas_policy` - (Optional) A `sas_policy` block as defined below.

* `allowed_copy_scope` - (Optional) Restrict copy to and from Storage Accounts within an AAD tenant or with Private Links to the same VNet. Possible values are `AAD` and `PrivateLink`.

* `sftp_enabled` - (Optional) Boolean, enable SFTP for the storage account

-> **Note:** SFTP support requires `is_hns_enabled` set to `true`. [More information on SFTP support can be found here](https://learn.microsoft.com/azure/storage/blobs/secure-file-transfer-protocol-support). Defaults to `false`

* `dns_endpoint_type` - (Optional) Specifies which DNS endpoint type to use. Possible values are `Standard` and `AzureDnsZone`. Defaults to `Standard`. Changing this forces a new resource to be created.

-> **Note:** Azure DNS zone support requires `PartitionedDns` feature to be enabled. To enable this feature for your subscription, use the following command: `az feature register --namespace "Microsoft.Storage" --name "PartitionedDns"`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `blob_properties` block supports the following:

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `delete_retention_policy` - (Optional) A `delete_retention_policy` block as defined below.

* `restore_policy` - (Optional) A `restore_policy` block as defined below. This must be used together with `delete_retention_policy` set, `versioning_enabled` and `change_feed_enabled` set to `true`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

-> **Note:** `restore_policy` can not be configured when `dns_endpoint_type` is `AzureDnsZone`.

* `versioning_enabled` - (Optional) Is versioning enabled? Default to `false`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `change_feed_enabled` - (Optional) Is the blob service properties for change feed events enabled? Default to `false`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `change_feed_retention_in_days` - (Optional) The duration of change feed events retention in days. The possible values are between 1 and 146000 days (400 years). Setting this to null (or omit this in the configuration file) indicates an infinite retention of the change feed.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `default_service_version` - (Optional) The API Version which should be used by default for requests to the Data Plane API if an incoming request doesn't specify an API Version.

* `last_access_time_enabled` - (Optional) Is the last access time based tracking enabled? Default to `false`.

-> **Note:** This field cannot be configured when `kind` is set to `Storage` (V1).

* `container_delete_retention_policy` - (Optional) A `container_delete_retention_policy` block as defined below.

---

A `cors_rule` block supports the following:

* `allowed_headers` - (Required) A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - (Required) A list of HTTP methods that are allowed to be executed by the origin. Valid options are
`DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS`, `PUT` or `PATCH`.

* `allowed_origins` - (Required) A list of origin domains that will be allowed by CORS.

* `exposed_headers` - (Required) A list of response headers that are exposed to CORS clients.

* `max_age_in_seconds` - (Required) The number of seconds the client should cache a preflight response.

---

A `custom_domain` block supports the following:

* `name` - (Required) The Custom Domain Name to use for the Storage Account, which will be validated by Azure.

* `use_subdomain` - (Optional) Should the Custom Domain Name be validated by using indirect CNAME validation?

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Optional) The ID of the Key Vault Key, supplying a version-less key ID will enable auto-rotation of this key. Exactly one of `key_vault_key_id` and `managed_hsm_key_id` may be specified.

* `managed_hsm_key_id` -  (Optional) The ID of the managed HSM Key. Exactly one of `key_vault_key_id` and `managed_hsm_key_id` may be specified.

* `user_assigned_identity_id` - (Required) The ID of a user assigned identity.

~> **Note:** `customer_managed_key` can only be set when the `account_kind` is set to `StorageV2` or `account_tier` set to `Premium`, and the identity type is `UserAssigned`.

---

A `delete_retention_policy` block supports the following:

* `days` - (Optional) Specifies the number of days that the blob should be retained, between `1` and `365` days. Defaults to `7`.

* `permanent_delete_enabled` - (Optional) Indicates whether permanent deletion of the soft deleted blob versions and snapshots is allowed. Defaults to `false`.

~> **Note:** `permanent_delete_enabled` cannot be set to true if a `restore_policy` block is defined.

---

A `restore_policy` block supports the following:

* `days` - (Required) Specifies the number of days that the blob can be restored, between `1` and `365` days. This must be less than the `days` specified for `delete_retention_policy`.

---

A `container_delete_retention_policy` block supports the following:

* `days` - (Optional) Specifies the number of days that the container should be retained, between `1` and `365` days. Defaults to `7`.

---

A `hour_metrics` block supports the following:

* `enabled` - (Required) Indicates whether hour metrics are enabled for the Queue service.

* `version` - (Required) The version of storage analytics to configure.

* `include_apis` - (Optional) Indicates whether metrics should generate summary statistics for called API operations.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Storage Account. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Storage Account.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

~> **Note:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned`  and Storage Account has been created. More details are available below.

---

An `immutability_policy` block supports the following:

~> **Note:** This argument specifies the default account-level immutability policy which is inherited and applied to objects that do not possess an explicit immutability policy at the object level. The object-level immutability policy has higher precedence than the container-level immutability policy, which has a higher precedence than the account-level immutability policy.

* `allow_protected_append_writes` - (Required) When enabled, new blocks can be written to an append blob while maintaining immutability protection and compliance. Only new blocks can be added and any existing blocks cannot be modified or deleted.

* `state` - (Required) Defines the mode of the policy. `Disabled` state disables the policy, `Unlocked` state allows increase and decrease of immutability retention time and also allows toggling allowProtectedAppendWrites property, `Locked` state only allows the increase of the immutability retention time. A policy can only be created in a Disabled or Unlocked state and can be toggled between the two states. Only a policy in an Unlocked state can transition to a Locked state which cannot be reverted. Changing from `Locked` forces a new resource to be created.

* `period_since_creation_in_days` - (Required) The immutability period for the blobs in the container since the policy creation, in days.

---

A `logging` block supports the following:

* `delete` - (Required) Indicates whether all delete requests should be logged.

* `read` - (Required) Indicates whether all read requests should be logged.

* `version` - (Required) The version of storage analytics to configure.

* `write` - (Required) Indicates whether all write requests should be logged.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained.

---

A `minute_metrics` block supports the following:

* `enabled` - (Required) Indicates whether minute metrics are enabled for the Queue service.

* `version` - (Required) The version of storage analytics to configure.

* `include_apis` - (Optional) Indicates whether metrics should generate summary statistics for called API operations.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained.

---

A `network_rules` block supports the following:

* `default_action` - (Required) Specifies the default action of allow or deny when no other rules match. Valid options are `Deny` or `Allow`.
* `bypass` - (Optional) Specifies whether traffic is bypassed for Logging/Metrics/AzureServices. Valid options are any combination of `Logging`, `Metrics`, `AzureServices`, or `None`.
* `ip_rules` - (Optional) List of public IP or IP ranges in CIDR Format. Only IPv4 addresses are allowed. /31 CIDRs, /32 CIDRs, and Private IP address ranges (as defined in [RFC 1918](https://tools.ietf.org/html/rfc1918#section-3)), are not allowed.

* `virtual_network_subnet_ids` - (Optional) A list of resource ids for subnets.

* `private_link_access` - (Optional) One or more `private_link_access` block as defined below.

~> **Note:** If specifying `network_rules`, one of either `ip_rules` or `virtual_network_subnet_ids` must be specified and `default_action` must be set to `Deny`.

~> **Note:** Network Rules can be defined either directly on the `azurerm_storage_account` resource, or using the `azurerm_storage_account_network_rules` resource - but the two cannot be used together. If both are used against the same Storage Account, spurious changes will occur. When managing Network Rules using this resource, to change from a `default_action` of `Deny` to `Allow` requires defining, rather than removing, the block.

~> **Note:** The prefix of `ip_rules` must be between 0 and 30 and only supports public IP addresses.

~> **Note:** [More information on Validation is available here](https://docs.microsoft.com/en-gb/azure/storage/blobs/storage-custom-domain-name)

---

A `private_link_access` block supports the following:

* `endpoint_resource_id` - (Required) The ID of the Azure resource that should be allowed access to the target storage account.

* `endpoint_tenant_id` - (Optional) The tenant id of the resource of the resource access rule to be granted access. Defaults to the current tenant id.

---

A `azure_files_authentication` block supports the following:

* `directory_type` - (Required) Specifies the directory service used. Possible values are `AADDS`, `AD` and `AADKERB`.

* `active_directory` - (Optional) A `active_directory` block as defined below. Required when `directory_type` is `AD`.

* `default_share_level_permission` - (Optional) Specifies the default share level permissions applied to all users. Possible values are `StorageFileDataSmbShareReader`, `StorageFileDataSmbShareContributor`, `StorageFileDataSmbShareElevatedContributor`, or `None`.

---

A `active_directory` block supports the following:

* `domain_name` - (Required) Specifies the primary domain that the AD DNS server is authoritative for.

* `domain_guid` - (Required) Specifies the domain GUID.

* `domain_sid` - (Optional) Specifies the security identifier (SID). This is required when `directory_type` is set to `AD`.

* `storage_sid` - (Optional) Specifies the security identifier (SID) for Azure Storage. This is required when `directory_type` is set to `AD`.

* `forest_name` - (Optional) Specifies the Active Directory forest. This is required when `directory_type` is set to `AD`.

* `netbios_domain_name` - (Optional) Specifies the NetBIOS domain name. This is required when `directory_type` is set to `AD`.

---

A `routing` block supports the following:

* `publish_internet_endpoints` - (Optional) Should internet routing storage endpoints be published? Defaults to `false`.

* `publish_microsoft_endpoints` - (Optional) Should Microsoft routing storage endpoints be published? Defaults to `false`.

* `choice` - (Optional) Specifies the kind of network routing opted by the user. Possible values are `InternetRouting` and `MicrosoftRouting`. Defaults to `MicrosoftRouting`.

---

A `queue_properties` block supports the following:

* `cors_rule` - (Optional) A `cors_rule` block as defined above.

* `logging` - (Optional) A `logging` block as defined below.

* `minute_metrics` - (Optional) A `minute_metrics` block as defined below.

* `hour_metrics` - (Optional) A `hour_metrics` block as defined below.

---

A `sas_policy` block supports the following:

* `expiration_period` - (Required) The SAS expiration period in format of `DD.HH:MM:SS`.

* `expiration_action` - (Optional) The SAS expiration action. The only possible value is `Log` at this moment. Defaults to `Log`.

---

A `static_website` block supports the following:

* `index_document` - (Optional) The webpage that Azure Storage serves for requests to the root of a website or any subfolder. For example, index.html. The value is case-sensitive.

* `error_404_document` - (Optional) The absolute path to a custom webpage that should be used when a request is made which does not correspond to an existing file.

---

A `share_properties` block supports the following:

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `retention_policy` - (Optional) A `retention_policy` block as defined below.

* `smb` - (Optional) A `smb` block as defined below.

---

A `retention_policy` block supports the following:

* `days` - (Optional) Specifies the number of days that the `azurerm_storage_share` should be retained, between `1` and `365` days. Defaults to `7`.

---

A `smb` block supports the following:

* `versions` - (Optional) A set of SMB protocol versions. Possible values are `SMB2.1`, `SMB3.0`, and `SMB3.1.1`.

* `authentication_types` - (Optional) A set of SMB authentication methods. Possible values are `NTLMv2`, and `Kerberos`.

* `kerberos_ticket_encryption_type` - (Optional) A set of Kerberos ticket encryption. Possible values are `RC4-HMAC`, and `AES-256`.

* `channel_encryption_type` - (Optional) A set of SMB channel encryption. Possible values are `AES-128-CCM`, `AES-128-GCM`, and `AES-256-GCM`.

* `multichannel_enabled` - (Optional) Indicates whether multichannel is enabled. Defaults to `false`. This is only supported on Premium storage accounts.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account.

* `primary_location` - The primary location of the storage account.

* `secondary_location` - The secondary location of the storage account.

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

* `primary_access_key` - The primary access key for the storage account.

* `secondary_access_key` - The secondary access key for the storage account.

* `primary_connection_string` - The connection string associated with the primary location.

* `secondary_connection_string` - The connection string associated with the secondary location.

* `primary_blob_connection_string` - The connection string associated with the primary blob location.

* `secondary_blob_connection_string` - The connection string associated with the secondary blob location.

~> **Note:** If there's a write-lock on the Storage Account, or the account doesn't have permission then these fields will have an empty value [due to a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/6363)

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Storage Account.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Storage Account.

-> **Note:** You can access the Principal ID via `${azurerm_storage_account.example.identity[0].principal_id}` and the Tenant ID via `${azurerm_storage_account.example.identity[0].tenant_id}`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Storage Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account.
* `update` - (Defaults to 1 hour) Used when updating the Storage Account.
* `delete` - (Defaults to 1 hour) Used when deleting the Storage Account.

## Import

Storage Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account.storageAcc1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
