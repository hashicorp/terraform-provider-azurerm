---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_job_definition"
description: |-
  Manages a Storage Mover Job Definition.
---

# azurerm_storage_mover_job_definition

Manages a Storage Mover Job Definition.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_storage_mover_agent" "example" {
  name                     = "example-agent"
  storage_mover_id         = azurerm_storage_mover.example.id
  arc_virtual_machine_id   = "${azurerm_resource_group.example.id}/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"
  arc_virtual_machine_uuid = "3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"
}

resource "azurerm_storage_account" "example" {
  name                            = "examplesa"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "example" {
  name                  = "acccontainer"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "blob"
}

resource "azurerm_storage_mover_target_endpoint" "example" {
  name                   = "example-smte"
  storage_mover_id       = azurerm_storage_mover.example.id
  storage_account_id     = azurerm_storage_account.example.id
  storage_container_name = azurerm_storage_container.example.name
}

resource "azurerm_storage_mover_source_endpoint" "example" {
  name             = "example-smse"
  storage_mover_id = azurerm_storage_mover.example.id
  host             = "192.168.0.1"
}

resource "azurerm_storage_mover_project" "example" {
  name             = "example-sp"
  storage_mover_id = azurerm_storage_mover.example.id
}

resource "azurerm_storage_mover_job_definition" "example" {
  name                     = "example-sjd"
  storage_mover_project_id = azurerm_storage_mover_project.example.id
  agent_name               = azurerm_storage_mover_agent.example.name
  copy_mode                = "Additive"
  source_name              = azurerm_storage_mover_source_endpoint.example.name
  source_sub_path          = "/"
  target_name              = azurerm_storage_mover_target_endpoint.example.name
  target_sub_path          = "/"
  description              = "Example Job Definition Description"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Job Definition. Changing this forces a new resource to be created.

* `storage_mover_project_id` - (Required) Specifies the ID of the Storage Mover Project. Changing this forces a new resource to be created.

* `source_name` - (Required) Specifies the name of the Storage Mover Source Endpoint. Changing this forces a new resource to be created.

* `target_name` - (Required) Specifies the name of the Storage Mover target Endpoint. Changing this forces a new resource to be created.

* `copy_mode` - (Required) Specifies the strategy to use for copy. Possible values are `Additive` and `Mirror`.

* `source_sub_path` - (Optional) Specifies the sub path to use when reading from the Storage Mover Source Endpoint. Changing this forces a new resource to be created.

* `target_sub_path` - (Optional) Specifies the sub path to use when writing to the Storage Mover Target Endpoint. Changing this forces a new resource to be created.

* `agent_name` - (Optional) Specifies the name of the Storage Mover Agent to assign for new Job Runs of this Storage Mover Job Definition.

* `description` - (Optional) Specifies a description for this Storage Mover Job Definition.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Job Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Job Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Job Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Job Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Job Definition.

## Import

Storage Mover Job Definition can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_job_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/projects/project1/jobDefinitions/jobDefinition1
```
