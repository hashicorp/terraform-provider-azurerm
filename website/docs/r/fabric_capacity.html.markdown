---
subcategory: "Fabric"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_fabric_capacity"
description: |-
  Manages a Fabric Capacity.
---

# azurerm_fabric_capacity

Manages a Fabric Capacity.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_fabric_capacity" "example" {
  name                = "exampleffc"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"

  administration_members = [data.azurerm_client_config.current.object_id]

  sku {
    name = "F32"
    tier = "Fabric"
  }

  tags = {
    environment = "test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the Fabric Capacity. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the Fabric Capacity exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the Fabric Capacity. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `administration_members` - (Optional) An array of administrator user identities. The member must be an Entra user or a service principal.

~> **Note:** If the member is an Entra user, use user principal name (UPN) format. If the user is a service principal, use object ID.

* `tags` - (Optional) A mapping of tags to assign to the Fabric Capacity.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU to use for the Fabric Capacity. Possible values are `F2`, `F4`, `F8`, `F16`, `F32`, `F64`, `F128`, `F256`, `F512`, `F1024`, `F2048`.

* `tier` - (Required) The tier of the SKU to use for the Fabric Capacity. The only possible value is `Fabric`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Fabric Capacity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Fabric Capacity.
* `read` - (Defaults to 5 minutes) Used when retrieving the Fabric Capacity.
* `update` - (Defaults to 30 minutes) Used when updating the Fabric Capacity.
* `delete` - (Defaults to 30 minutes) Used when deleting the Fabric Capacity.

## Import

Fabric Capacities can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_fabric_capacity.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Fabric/capacities/capacity1
```
