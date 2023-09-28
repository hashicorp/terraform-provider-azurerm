---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_load_balancer"
description: |-
  Manages an Application Gateway for Containers (ALB).
---

# azurerm_application_load_balancer

Manages an Application Gateway for Containers (ALB).

## Example Usage

```hcl
resource "azurerm_application_load_balancer" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the resource should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the resource should exist. Changing this forces a new resource to be created. Available regions can be found [here](https://learn.microsoft.com/en-us/azure/application-gateway/for-containers/overview#supported-regions)

**Note:** The available values of `location` are `northeurope` and `north central us`.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Gateway for Containers.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the resource.

* `configuration_endpoint` - The list of configuration endpoints for the resource. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Gateway for Containers resource
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Gateway for Containers resource.
* `update` - (Defaults to 30 minutes) Used when updating the Application Gateway for Containers resource
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Gateway for Containers resource.

## Import

Application Gateway for Containers (ALB) can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_load_balancer.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceNetworking/trafficControllers/alb1
```
