---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_script"
description: |-
  Manages a Kusto Script.
---

# azurerm_kusto_script

Manages a Kusto Script.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "setup-files"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "example" {
  name                   = "script.txt"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Block"
  source_content         = ".create table MyTable (Level:string, Timestamp:datetime, UserId:string, TraceId:string, Message:string, ProcessId:int32)"
}

data "azurerm_storage_account_blob_container_sas" "example" {
  connection_string = azurerm_storage_account.example.primary_connection_string
  container_name    = azurerm_storage_container.example.name
  https_only        = true

  start  = "2017-03-21"
  expiry = "2022-03-21"

  permissions {
    read   = true
    add    = false
    create = false
    write  = true
    delete = false
    list   = true
  }
}

resource "azurerm_kusto_script" "example" {
  name                               = "example"
  database_id                        = azurerm_kusto_database.example.id
  url                                = azurerm_storage_blob.example.id
  sas_token                          = data.azurerm_storage_account_blob_container_sas.example.sas
  continue_on_errors_enabled         = true
  force_an_update_when_value_changed = "first"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Kusto Script. Changing this forces a new Kusto Script to be created.
  
* `database_id` - (Required) The ID of the Kusto Database. Changing this forces a new Kusto Script to be created.

---

* `continue_on_errors_enabled` - (Optional) Flag that indicates whether to continue if one of the command fails.

* `force_an_update_when_value_changed` - (Optional) A unique string. If changed the script will be applied again.

* `script_content` - (Optional) The script content. This property should be used when the script is provide inline and not through file in a SA. Must not be used together with `url` and `sas_token` properties. Changing this forces a new resource to be created.

* `sas_token` - (Optional) The SAS token used to access the script. Must be provided when using scriptUrl property. Changing this forces a new resource to be created.

* `url` - (Optional) The url to the KQL script blob file. Must not be used together with scriptContent property. Please reference [this documentation](https://docs.microsoft.com/azure/data-explorer/database-script) that describes the commands that are allowed in the script.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kusto Script.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kusto Script.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Script.
* `update` - (Defaults to 30 minutes) Used when updating the Kusto Script.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kusto Script.

## Import

Kusto Scripts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_script.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/scripts/script1
```
