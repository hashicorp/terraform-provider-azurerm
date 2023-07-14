---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_alb"
description: |-
  Manages an Application Load Balancer.
---

# azurerm_alb

Manages an Application Load Balancer.

## Example Usage

```hcl
resource "azurerm_alb" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Application Load Balancer. Changing this forces a new Application Load Balancer to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Application Load Balancer should exist. Changing this forces a new Application Load Balancer to be created.

* `location` - (Required) The Azure Region where the Application Load Balancer should exist. Changing this forces a new Application Load Balancer to be created.

**Note:** The available values of `location` are `northeurope` and `north central us`.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Load Balancer.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Load Balancer.

* `configuration_endpoint` - The list of configuration endpoints for the Application Load Balancer. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Load Balancer.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Load Balancer.
* `update` - (Defaults to 30 minutes) Used when updating the Application Load Balancer.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Load Balancer.

## Import

Application Load Balancers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_alb.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceNetworking/trafficControllers/alb1
```
