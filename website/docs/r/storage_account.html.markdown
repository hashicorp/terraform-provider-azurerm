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
  name               = "subnetname"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes   = ["10.0.2.0/24"]
  service_endpoints  = ["Microsoft.Sql", "Microsoft.Storage"]
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

* `name` - (Required) Specifies the name of the storage account. Changing this forces a new resource to be created. This must be unique across the entire Azure service, not just within the resource group.

* `resource_group_name` - (Required) The name of the resource group in which to create the storage account. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `account_kind` - (Optional) Defines the Kind of account. Valid options are `BlobStorage`, `BlockBlobStorage`, `FileStorage`, `Storage` and `StorageV2`. Changing this forces a new resource to be created. Defaults to `StorageV2`.

-> **NOTE:** Changing the `account_kind` value from `Storage` to `StorageV2` will not trigger a force new on the storage account, it will only upgrade the existing storage account from `Storage` to `StorageV2` keeping the existing storage account in place.

* `account_tier` - (Required) Defines the Tier to use for this storage account. Valid options are `Standard` and `Premium`. For `BlockBlobStorage` and `FileStorage` accounts only `Premium` is valid. Changing this forces a new resource to be created.

* `account_replication_type` - (Required) Defines the type of replication to use for this storage account. Valid options are `LRS`, `GRS`, `RAGRS`, `ZRS`, `GZRS` and `RAGZRS`.

* `access_tier` - (Optional) Defines the access tier for `BlobStorage`, `FileStorage` and `StorageV2` accounts. Valid options are `Hot` and `Cool`, defaults to `Hot`.

* `enable_https_traffic_only` - (Optional) Boolean flag which forces HTTPS if enabled, see [here](https://docs.microsoft.com/en-us/azure/storage/storage-require-secure-transfer/)
    for more information. Defaults to `true`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the storage account. Possible values are `TLS1_0`, `TLS1_1`, and `TLS1_2`. Defaults to `TLS1_0` for new storage accounts.

-> **NOTE:** At this time `min_tls_version` is only supported in the Public Cloud.

* `allow_blob_public_access` - Allow or disallow public access to all blobs or containers in the storage account. Defaults to `false`.

-> **NOTE:** At this time `allow_blob_public_access` is only supported in the Public Cloud.

* `is_hns_enabled` - (Optional) Is Hierarchical Namespace enabled? This can be used with Azure Data Lake Storage Gen 2 ([see here for more information](https://docs.microsoft.com/en-us/azure/storage/blobs/data-lake-storage-quickstart-create-account/)). Changing this forces a new resource to be created.

-> **NOTE:** When this is set to `true` the `account_tier` argument must be set to `Standard`

* `custom_domain` - (Optional) A `custom_domain` block as documented below.

* `identity` - (Optional) A `identity` block as defined below.

* `blob_properties` - (Optional) A `blob_properties` block as defined below.

* `queue_properties` - (Optional) A `queue_properties` block as defined below.

~> **NOTE:** `queue_properties` cannot be set when the `access_tier` is set to `BlobStorage`

* `static_website` - (Optional) A `static_website` block as defined below.

~> **NOTE:** `static_website` can only be set when the `account_kind` is set to `StorageV2` or `BlockBlobStorage`.

* `network_rules` - (Optional) A `network_rules` block as documented below.

* `large_file_share_enabled` - (Optional) Is Large File Share Enabled?

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `blob_properties` block supports the following:

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `delete_retention_policy` - (Optional) A `delete_retention_policy` block as defined below.

---

A `cors_rule` block supports the following:

* `allowed_headers` - (Required) A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - (Required) A list of http headers that are allowed to be executed by the origin. Valid options are
`DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS`, `PUT` or `PATCH`.

* `allowed_origins` - (Required) A list of origin domains that will be allowed by CORS.

* `exposed_headers` - (Required) A list of response headers that are exposed to CORS clients.

* `max_age_in_seconds` - (Required) The number of seconds the client should cache a preflight response.

---

A `custom_domain` block supports the following:

* `name` - (Required) The Custom Domain Name to use for the Storage Account, which will be validated by Azure.
* `use_subdomain` - (Optional) Should the Custom Domain Name be validated by using indirect CNAME validation?

---

A `delete_retention_policy` block supports the following:

* `days` - (Optional) Specifies the number of days that the blob should be retained, between `1` and `365` days. Defaults to `7`.

---

A `hour_metrics` block supports the following:

* `enabled` - (Required) Indicates whether hour metrics are enabled for the Queue service. Changing this forces a new resource.

* `version` - (Required) The version of storage analytics to configure. Changing this forces a new resource.

* `include_apis` - (Optional) Indicates whether metrics should generate summary statistics for called API operations.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained. Changing this forces a new resource.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Storage Account. At this time the only allowed value is `SystemAssigned`.

~> The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned`  and Storage Account has been created. More details are available below.

---

A `logging` block supports the following:

* `delete` - (Required) Indicates whether all delete requests should be logged. Changing this forces a new resource.

* `read` - (Required) Indicates whether all read requests should be logged. Changing this forces a new resource.

* `version` - (Required) The version of storage analytics to configure. Changing this forces a new resource.

* `write` - (Required) Indicates whether all write requests should be logged. Changing this forces a new resource.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained. Changing this forces a new resource.

---

A `minute_metrics` block supports the following:

* `enabled` - (Required) Indicates whether minute metrics are enabled for the Queue service. Changing this forces a new resource.

* `version` - (Required) The version of storage analytics to configure. Changing this forces a new resource.

* `include_apis` - (Optional) Indicates whether metrics should generate summary statistics for called API operations.

* `retention_policy_days` - (Optional) Specifies the number of days that logs will be retained. Changing this forces a new resource.

---

A `network_rules` block supports the following:

* `default_action` - (Required) Specifies the default action of allow or deny when no other rules match. Valid options are `Deny` or `Allow`.
* `bypass` - (Optional)  Specifies whether traffic is bypassed for Logging/Metrics/AzureServices. Valid options are
any combination of `Logging`, `Metrics`, `AzureServices`, or `None`.
* `ip_rules` - (Optional) List of public IP or IP ranges in CIDR Format. Only IPV4 addresses are allowed. Private IP address ranges (as defined in [RFC 1918](https://tools.ietf.org/html/rfc1918#section-3)) are not allowed.
* `virtual_network_subnet_ids` - (Optional) A list of resource ids for subnets.

~> **Note:** If specifying `network_rules`, one of either `ip_rules` or `virtual_network_subnet_ids` must be specified and `default_action` must be set to `Deny`.

~> **NOTE:** Network Rules can be defined either directly on the `azurerm_storage_account` resource, or using the `azurerm_storage_account_network_rules` resource - but the two cannot be used together. If both are used against the same Storage Account, spurious changes will occur. When managing Network Rules using this resource, to change from a `default_action` of `Deny` to `Allow` requires defining, rather than removing, the block.

~> **Note:** The prefix of `ip_rules` must be between 0 and 30 and only supports public IP addresses.

~> **Note:** [More information on Validation is available here](https://docs.microsoft.com/en-gb/azure/storage/blobs/storage-custom-domain-name)

---

A `queue_properties` block supports the following:

* `cors_rule` - (Optional) A `cors_rule` block as defined above.

* `logging` - (Optional) A `logging` block as defined below.

* `minute_metrics` - (Optional) A `minute_metrics` block as defined below.

* `hour_metrics` - (Optional) A `hour_metrics` block as defined below.

---

A `static_website` block supports the following:

* `index_document` - (Optional) The webpage that Azure Storage serves for requests to the root of a website or any subfolder. For example, index.html. The value is case-sensitive.

* `error_404_document` - (Optional) The absolute path to a custom webpage that should be used when a request is made which does not correspond to an existing file.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Account.

* `primary_location` - The primary location of the storage account.

* `secondary_location` - The secondary location of the storage account.

* `primary_blob_endpoint` - The endpoint URL for blob storage in the primary location.

* `primary_blob_host` - The hostname with port if applicable for blob storage in the primary location.

* `secondary_blob_endpoint` - The endpoint URL for blob storage in the secondary location.

* `secondary_blob_host` - The hostname with port if applicable for blob storage in the secondary location.

* `primary_queue_endpoint` - The endpoint URL for queue storage in the primary location.

* `primary_queue_host` - The hostname with port if applicable for queue storage in the primary location.

* `secondary_queue_endpoint` - The endpoint URL for queue storage in the secondary location.

* `secondary_queue_host` - The hostname with port if applicable for queue storage in the secondary location.

* `primary_table_endpoint` - The endpoint URL for table storage in the primary location.

* `primary_table_host` - The hostname with port if applicable for table storage in the primary location.

* `secondary_table_endpoint` - The endpoint URL for table storage in the secondary location.

* `secondary_table_host` - The hostname with port if applicable for table storage in the secondary location.

* `primary_file_endpoint` - The endpoint URL for file storage in the primary location.

* `primary_file_host` - The hostname with port if applicable for file storage in the primary location.

* `secondary_file_endpoint` - The endpoint URL for file storage in the secondary location.

* `secondary_file_host` - The hostname with port if applicable for file storage in the secondary location.

* `primary_dfs_endpoint` - The endpoint URL for DFS storage in the primary location.

* `primary_dfs_host` - The hostname with port if applicable for DFS storage in the primary location.

* `secondary_dfs_endpoint` - The endpoint URL for DFS storage in the secondary location.

* `secondary_dfs_host` - The hostname with port if applicable for DFS storage in the secondary location.

* `primary_web_endpoint` - The endpoint URL for web storage in the primary location.

* `primary_web_host` - The hostname with port if applicable for web storage in the primary location.

* `secondary_web_endpoint` - The endpoint URL for web storage in the secondary location.

* `secondary_web_host` - The hostname with port if applicable for web storage in the secondary location.

* `primary_access_key` - The primary access key for the storage account.

* `secondary_access_key` - The secondary access key for the storage account.

* `primary_connection_string` - The connection string associated with the primary location.

* `secondary_connection_string` - The connection string associated with the secondary location.

* `primary_blob_connection_string` - The connection string associated with the primary blob location.

* `secondary_blob_connection_string` - The connection string associated with the secondary blob location.

~> **NOTE:** If there's a Write Lock on the Storage Account, or the account doesn't have permission then these fields will have an empty value [due to a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/6363)

* `identity` - An `identity` block as defined below, which contains the Identity information for this Storage Account.

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Storage Account.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Storage Account.

-> You can access the Principal ID via `${azurerm_storage_account.example.identity.0.principal_id}` and the Tenant ID via `${azurerm_storage_account.example.identity.0.tenant_id}`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Storage Account.
* `update` - (Defaults to 60 minutes) Used when updating the Storage Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account.
* `delete` - (Defaults to 60 minutes) Used when deleting the Storage Account.

## Import

Storage Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account.storageAcc1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
