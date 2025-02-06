---
subcategory: "Service Networking"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_load_balancer_frontend"
description: |-
  Manages an Application Gateway for Containers Frontend.
---

# azurerm_application_load_balancer_frontend

Manages an Application Gateway for Containers Frontend.

## Example Usage

```hcl
resource "azurerm_application_load_balancer" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
}

resource "azurerm_application_load_balancer_frontend" "example" {
  name                         = "example"
  application_load_balancer_id = azurerm_application_load_balancer.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Application Gateway for Containers Frontend. Changing this forces a new resource to be created.

* `application_load_balancer_id` - (Required) The ID of the Application Gateway for Containers. Changing this forces a new resource to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Gateway for Containers Frontend.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Gateway for Containers Frontend.

* `fully_qualified_domain_name` - The Fully Qualified Domain Name of the DNS record associated to an Application Gateway for Containers Frontend.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Gateway for Containers Frontend.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Gateway for Containers Frontend.
* `update` - (Defaults to 30 minutes) Used when updating the Application Gateway for Containers Frontend.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Gateway for Containers Frontend.

## Import

Application Gateway for Containers Frontend can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_load_balancer_frontend.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ServiceNetworking/trafficControllers/alb1/frontends/frontend1
```
