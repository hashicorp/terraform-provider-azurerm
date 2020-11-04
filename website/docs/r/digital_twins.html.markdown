---
subcategory: "DigitalTwins"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_digital_twins"
description: |-
  Manages a Digital Twins.
---

# azurerm_digital_twins

Manages a Digital Twins.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example_resources"
  location = "West Europe"
}

resource "azurerm_digital_twins" "example" {
  name                = "example-DT"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Digital Twins. Changing this forces a new Digital Twins to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Digital Twins should exist. Changing this forces a new Digital Twins to be created.

* `location` - (Required) The Azure Region where the Digital Twins should exist. Changing this forces a new Digital Twins to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Digital Twins.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Digital Twins.

* `host_name` - The Api endpoint to work with this Digital Twins.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Digital Twins.
* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins.
* `update` - (Defaults to 30 minutes) Used when updating the Digital Twins.
* `delete` - (Defaults to 30 minutes) Used when deleting the Digital Twins.

## Import

Digital Twinss can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_digital_twins.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt1
```
