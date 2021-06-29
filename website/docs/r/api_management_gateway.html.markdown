---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_gateway"
description: |-
  Manages a API Management Gateway.
---

# azurerm_api_management_gateway

Manages a API Management Gateway.

## Example Usage

```hcl
resource "azurerm_api_management_gateway" "example" {
  resource_group_name = "example"
  location            = "West Europe"
  gateway_id          = "my-gateway"
  api_management_name = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_name` - (Required) The Name of the API Management Service in which this Gateway exists. Changing this forces a new API Management Gateway to be created.

* `gateway_id` - (Required) The Identifier for the API Management Gateway. Changing this forces a new API Management Gateway to be created.

* `location` - (Required) The Azure Region where the API Management Gateway should exist.

* `resource_group_name` - (Required) The name of the Resource Group where the API Management Gateway should exist. Changing this forces a new API Management Gateway to be created.

---

* `description` - (Optional) Description of the API Management Gateway.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the API Management Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Gateway.
* `update` - (Defaults to 10 minutes) Used when updating the API Management Gateway.
* `delete` - (Defaults to 10 minutes) Used when deleting the API Management Gateway.

## Import

API Management Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1
```
