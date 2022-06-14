---
subcategory: "Load Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_load_test"
description: |-
  Manages a Load Test.
---

# azurerm_load_test

Manages a Load Test.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_load_test" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Load Test should exist. Changing this forces a new Load Test to be created.

* `name` - (Required) The name which should be used for this Load Test. Changing this forces a new Load Test to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Load Test should exist. Changing this forces a new Load Test to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Load Test.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Load Test.

* `dataplane_uri` - Public URI of the Data Plane.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Load Test.
* `read` - (Defaults to 5 minutes) Used when retrieving the Load Test.
* `update` - (Defaults to 30 minutes) Used when updating the Load Test.
* `delete` - (Defaults to 30 minutes) Used when deleting the Load Test.

## Import

Load tests can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_load_test.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.LoadTestService/loadtests/example
```
