---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_trigger_blob_event"
description: |-
  Manages a Blob Event Trigger inside an Azure Data Factory.
---

# azurerm_data_factory_trigger_blob_event

Manages a Blob Event Trigger inside an Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_pipeline" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  data_factory_name   = azurerm_data_factory.example.name
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory_trigger_blob_event" "example" {
  name                = "example"
  data_factory_id     = azurerm_data_factory.example.id
  storage_account_id  = azurerm_storage_account.example.id
  events              = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobDeleted"]
  blob_path_ends_with = ".txt"
  ignore_empty_blobs  = true

  annotations = ["test1", "test2", "test3"]
  description = "example description"

  pipeline {
    name = azurerm_data_factory_pipeline.example.name
    parameters = {
      Env = "Prod"
    }
  }

  additional_properties = {
    foo = "foo1"
    bar = "bar2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Blob Event Trigger. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The ID of Data Factory in which to associate the Trigger with. Changing this forces a new resource.

* `storage_account_id` - (Required) The ID of Storage Account in which blob event will be listened. Changing this forces a new resource.

* `events` - (Required) List of events that will fire this trigger. Possible values are `Microsoft.Storage.BlobCreated` and `Microsoft.Storage.BlobDeleted`.

* `pipeline` - (Required) One or more `pipeline` blocks as defined below.

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Blob Event Trigger.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Blob Event Trigger.

* `blob_path_begins_with` - (Optional) The pattern that blob path starts with for trigger to fire.

* `blob_path_ends_with` - (Optional) The pattern that blob path ends with for trigger to fire.

~> **Note:** At least one of `blob_path_begins_with` and `blob_path_ends_with` must be set.

* `description` - (Optional) The description for the Data Factory Blob Event Trigger.

* `ignore_empty_blobs` - (Optional) are blobs with zero bytes ignored?

---

A `pipeline` block supports the following:

* `name` - (Required) The Data Factory Pipeline name that the trigger will act on.

* `parameters` - (Optional) The Data Factory Pipeline parameters that the trigger will act on.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Factory Blob Event Trigger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Blob Event Trigger.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Blob Event Trigger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Blob Event Trigger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Blob Event Trigger.

## Import

Data Factory Blob Event Trigger can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_trigger_blob_event.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/triggers/example
```
