---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account"
sidebar_current: "docs-azurerm-resource-storage-account"
description: |-
  Manages a Azure Storage Account.
---

# azurerm_storage_account

Manages an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "storageaccountname"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}
```

## Example Usage with Network Rules

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "test" {
  name                = "virtnetname"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnetname"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Sql", "Microsoft.Storage"]
}

resource "azurerm_storage_account" "testsa" {
  name                = "storageaccountname"
  resource_group_name = "${azurerm_resource_group.test.name}"

  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["100.0.0.1"]
    virtual_network_subnet_ids = ["${azurerm_subnet.test.id}"]
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

* `account_kind` - (Optional) Defines the Kind of account. Valid options are `BlobStorage`, `BlockBlobStorage`, `FileStorage`, `Storage` and `StorageV2`. Changing this forces a new resource to be created. Defaults to `Storage`.

* `account_tier` - (Required) Defines the Tier to use for this storage account. Valid options are `Standard` and `Premium`. For `FileStorage` accounts only `Premium` is valid. Changing this forces a new resource to be created.

* `account_replication_type` - (Required) Defines the type of replication to use for this storage account. Valid options are `LRS`, `GRS`, `RAGRS` and `ZRS`.

* `access_tier` - (Optional) Defines the access tier for `BlobStorage`, `FileStorage` and `StorageV2` accounts. Valid options are `Hot` and `Cool`, defaults to `Hot`.

* `enable_blob_encryption` - (Optional) Boolean flag which controls if Encryption Services are enabled for Blob storage, see [here](https://azure.microsoft.com/en-us/documentation/articles/storage-service-encryption/) for more information. Defaults to `true`.

* `enable_file_encryption` - (Optional) Boolean flag which controls if Encryption Services are enabled for File storage, see [here](https://azure.microsoft.com/en-us/documentation/articles/storage-service-encryption/) for more information. Defaults to `true`.

* `enable_https_traffic_only` - (Optional) Boolean flag which forces HTTPS if enabled, see [here](https://docs.microsoft.com/en-us/azure/storage/storage-require-secure-transfer/)
    for more information.
    
* `is_hns_enabled` - (Optional) Is Hierarchical Namespace enabled? This can be used with Azure Data Lake Storage Gen 2 ([see here for more information](https://docs.microsoft.com/en-us/azure/storage/blobs/data-lake-storage-quickstart-create-account/)). Changing this forces a new resource to be created.

* `account_encryption_source` - (Optional) The Encryption Source for this Storage Account. Possible values are `Microsoft.Keyvault` and `Microsoft.Storage`. Defaults to `Microsoft.Storage`.

* `custom_domain` - (Optional) A `custom_domain` block as documented below.

* `enable_advanced_threat_protection` (Optional) Boolean flag which controls if advanced threat protection is enabled, see [here](https://docs.microsoft.com/en-us/azure/storage/common/storage-advanced-threat-protection) for more information. Defaults to `false`.

~> **Note:** `enable_advanced_threat_protection` is not supported in all regions.

* `identity` - (Optional) A `identity` block as defined below.

* `queue_properties` - (Optional) A `queue_properties` block as defined below.

~> **NOTE:** `queue_properties` cannot be set when the `access_tier` is set to `BlobStorage`

* `network_rules` - (Optional) A `network_rules` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `cors_rule` block supports the following:

* `allowed_headers` - (Required) A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - (Required) A list of http headers that are allowed to be executed by the origin. Valid options are
`DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS` or `PUT`.

* `allowed_origins` - (Required) A list of origin domains that will be allowed by CORS. 

* `exposed_headers` - (Required) A list of response headers that are exposed to CORS clients. 

* `max_age_in_seconds` - (Required) The number of seconds the client should cache a preflight response.

---

A `custom_domain` block supports the following:

* `name` - (Optional) The Custom Domain Name to use for the Storage Account, which will be validated by Azure.
* `use_subdomain` - (Optional) Should the Custom Domain Name be validated by using indirect CNAME validation?

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

~> **Note:** [More information on Validation is available here](https://docs.microsoft.com/en-gb/azure/storage/blobs/storage-custom-domain-name)

---

A `queue_properties` block supports the following:

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `logging` - (Optional) A `logging` block as defined below.

* `minute_metrics` - (Optional) A `minute_metrics` block as defined below.

* `hour_metrics` - (Optional) A `hour_metrics` block as defined below.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The storage account Resource ID.

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

-> You can access the Principal ID via `${azurerm_storage_account.test.identity.0.principal_id}` and the Tenant ID via `${azurerm_storage_account.test.identity.0.tenant_id}`

## Import

Storage Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account.storageAcc1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
