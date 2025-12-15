---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_perimeter_profile"
description: |-
  Manages a Network Security Perimeter Profile.
---

# azurerm_network_security_perimeter_profile

Manages a Network Security Perimeter Profile.

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

resource "azurerm_network_security_perimeter_profile" "example" {
  name = "example"
  perimeter_id = azurerm_network_security_perimeter.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Security Perimeter Profile. Changing this forces a new Network Security Perimeter Profile to be created.

* `perimeter_id` - (Required) The ID of the Network Security Perimeter within this Profile is created. Changing this forces a new Network Security Perimeter Profile to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Security Perimeter Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Security Perimeter Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Perimeter Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Security Perimeter Profile.

## Import

Network Security Perimeter Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_security_perimeter_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/networkSecurityPerimeters/example-nsp/profiles/defaultProfile
```