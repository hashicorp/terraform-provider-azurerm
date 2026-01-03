---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_perimeter"
description: |-
  Manages a Network Security Perimeter.
---

# azurerm_network_security_perimeter

Manages a Network Security Perimeter.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_perimeter" "example" {
  name = "example"
  resource_group_name = azurerm_resource_group.example.name
  location = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Network Security Perimeter should exist. Changing this forces a new Network Security Perimeter to be created.

* `name` - (Required) The name which should be used for this Network Security Perimeter. Changing this forces a new Network Security Perimeter to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Network Security Perimeter should exist.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Security Perimeter.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Security Perimeter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Security Perimeter.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Perimeter.
* `update` - (Defaults to 30 minutes) Used when updating the Network Security Perimeter.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Security Perimeter.

## Import

Network Security Perimeters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_security_perimeter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/networkSecurityPerimeters/example-nsp
```