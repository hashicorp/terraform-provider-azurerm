---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_machine_automanage_configuration_assignment"
description: |-
  Manages an Arc Machine Automanage Configuration Profile Assignment.
---

# azurerm_arc_machine_automanage_configuration_assignment

Manages an Arc Machine Automanage Configuration Profile Assignment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

variable "arc_machine_name" {
  description = "The name of the Arc Machine."
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_arc_machine" "example" {
  name                = var.arc_machine_name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_automanage_configuration" "example" {
  name                = "example-configuration"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_arc_machine_automanage_configuration_assignment" "example" {
  arc_machine_id   = data.azurerm_arc_machine.example.id
  configuration_id = azurerm_automanage_configuration.example.id
}

```

## Arguments Reference

The following arguments are supported:

* `arc_machine_id` - (Required) The ARM resource ID of the Arc Machine to assign the Automanage Configuration to. Changing this forces a new resource to be created.

* `configuration_id` - (Required) The ARM resource ID of the Automanage Configuration to assign to the Virtual Machine. Changing this forces a new resource to be created.

~> **Note:** For a successful creation of this resource, locate "Automanage API Access" app within your Entra ID tenant. Make sure it's granted access to the scope that includes the arc server.

---
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Arc Machine Automanage Configuration Profile Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automanage Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automanage Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automanage Configuration.

## Import

Virtual Machine Automanage Configuration Profile Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_arc_machine_automanage_configuration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AutoManage/configurationProfileAssignments/default
```
