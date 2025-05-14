---
subcategory: "Hybrid Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_machine"
description: |-
  Manages a Hybrid Compute Machine.
---

# azurerm_arc_machine

Manages a Hybrid Compute Machine.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_machine" "example" {
  name                = "example-arcmachine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  kind                = "SCVMM"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Arc machine. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Arc Machine should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Arc Machine should exist. Changing this forces a new resource to be created.

* `kind` - (Required) The kind of the Arc Machine. Possible values are `AVS`, `AWS`, `EPS`, `GCP`, `HCI`, `SCVMM` and `VMware`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Arc Machine.

---

* An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity assigned to this Arc Machine. At this time the only possible value is `SystemAssigned`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Arc Machine.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Arc Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Arc Machine.
* `update` - (Defaults to 30 minutes) Used when updating the Arc Machine.
* `delete` - (Defaults to 30 minutes) Used when deleting the Arc Machine.

## Import

Arc Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_arc_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1
```
