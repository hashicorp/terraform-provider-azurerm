---
subcategory: "Digital Twins"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_digital_twins_instance"
description: |-
  Manages a Digital Twins instance.
---

# azurerm__digital_twins_instance

Manages a Digital Twins instance.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example_resources"
  location = "West Europe"
}

resource "azurerm_digital_twins_instance" "example" {
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

* `name` - (Required) The name which should be used for this Digital Twins instance. Changing this forces a new Digital Twins instance to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Digital Twins instance should exist. Changing this forces a new Digital Twins instance to be created.

* `location` - (Required) The Azure Region where the Digital Twins instance should exist. Changing this forces a new Digital Twins instance to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Digital Twins instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Digital Twins instance.

* `host_name` - The Api endpoint to work with this Digital Twins instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Digital Twins instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins instance.
* `update` - (Defaults to 30 minutes) Used when updating the Digital Twins instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Digital Twins instance.

## Import

Digital Twins instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_digital_twins_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt1
```
