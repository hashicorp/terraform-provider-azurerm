---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_agent"
description: |-
  Manages a Storage Mover Agent.
---

# azurerm_storage_mover_agent

Manages a Storage Mover Agent.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_storage_mover_agent" "example" {
  name                     = "example-sa"
  storage_mover_id         = azurerm_storage_mover.example.id
  arc_virtual_machine_id   = "${azurerm_resource_group.example.id}/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"
  arc_virtual_machine_uuid = "3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"
  description              = "Example Agent Description"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover Agent. Changing this forces a new resource to be created.

* `arc_virtual_machine_id` - (Required) Specifies the fully qualified ID of the Hybrid Compute resource for the Storage Mover Agent. Changing this forces a new resource to be created.

* `arc_virtual_machine_uuid` - (Required) Specifies the Hybrid Compute resource's unique SMBIOS ID. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the Storage Mover that this Agent should be connected to. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies a description for this Storage Mover Agent.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover Agent.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover Agent.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover Agent.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover Agent.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover Agent.

## Import

Storage Mover Agent can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_agent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/agents/agent1
```
