---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_storage_account_table_sas"
description: |-
  Gets information about an existing Storage Account Table Sas.
---

# Data Source: azurerm_storage_account_table_sas

Use this data source to obtain a Shared Access Signature (SAS Token) for an existing Storage Account Table.

Shared access signatures allow fine-grained, ephemeral access control to various aspects of an Azure Storage Account Table.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "resourceGroupName"
  location = "westus"
}

resource "azurerm_storage_account" "storage" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "table" {
  name                  = "newtable"
  storage_account_name  = azurerm_storage_account.storage.name
}

data "azurerm_storage_account_table_sas" "example" {
  connection_string = azurerm_storage_account.storage.primary_connection_string
  table_name        = azurerm_storage_table.table.name
  https_only        = true

  ip_address = "168.1.5.65"

  start  = "2018-03-21"
  expiry = "2018-03-21"

  permissions {
    read   = true
    add    = true
    update = false
    delete = false    
  }

  start_partition_key = "Coho Winery"
  end_partition_key   = "Auburn"
  start_row_key       = "Coho Winery"
  end_row_key         = "Seattle"
}

output "sas_url_query_string" {
  value = data.azurerm_storage_account_table_sas.example.sas
}
```

## Arguments Reference

The following arguments are supported:

* `connection_string` - The connection string for the storage account to which this SAS applies. Typically directly from the primary_connection_string attribute of a terraform created azurerm_storage_account resource.

* `table_name` - Name of the table.

* `https_only` - (Optional) Only permit https access. If false, both http and https are permitted. Defaults to true.

* `ip_address` - (Optional) Single ipv4 address or range (connected with a dash) of ipv4 addresses.

* `start` - The starting time and date of validity of this SAS. Must be a valid ISO-8601 format time/date string.

* `expiry` - The expiration time and date of this SAS. Must be a valid ISO-8601 format time/date string.

* `permissions` - A `permissions` block as defined below.

* `start_partition_key` - (Optional) First Partition Key in a range of Partition Keys to allow the SAS to access.

* `end_partition_key` - (Optional) Last Partition Key in a range of Partition Keys to allow the SAS to access.

* `start_row_key` - (Optional) First Row Key in a range of Row Keys to allow the SAS to access.

* `end_row_key` - (Optional) Last Row Key in a range of Row Keys to allow the SAS to access.

---

A `permissions` block supports the following:

* `read` - Should Read permissions be enabled for this SAS?

* `add` - Should Add permissions be enabled for this SAS?

* `update` - Should Update permissions be enabled for this SAS?

* `delete` - Should Delete permissions be enabled for this SAS?

Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas)
for additional details on the fields above.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `sas` - The computed Table Shared Access Signature (SAS).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Table.