---
subcategory: "Custom Providers"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_custom_provider"
description: |-
  Manages an Azure Custom Provider.
---

# azurerm_custom_provider

Manages an Azure Custom Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_custom_provider" "example" {
  name                = "example_provider"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  resource_type {
    name     = "dEf1"
    endpoint = "https://testendpoint.com/"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Custom Provider. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Custom Provider.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_type` - (Optional) Any number of `resource_type` block as defined below. One of `resource_type` or `action` must be specified.

* `action` - (Optional) Any number of `action` block as defined below. One of `resource_type` or `action` must be specified.

* `validation` - (Optional) Any number of `validation` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `resource_type` block supports the following:

* `name` - (Required) Specifies the name of the route definition. 

* `endpoint` - (Required) Specifies the endpoint of the route definition. 

* `routing_type` - (Optional) The routing type that is supported for the resource request. Valid values are `ResourceTypeRoutingProxy` or `ResourceTypeRoutingProxyCache`. This value defaults to `ResourceTypeRoutingProxy`. 

---

A `action` block supports the following:

* `name` - (Required) Specifies the name of the action. 

* `endpoint` - (Required) Specifies the endpoint of the action. 

---

A `validation` block supports the following:

* `specification` - (Required) The endpoint where the validation specification is located. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Custom Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `read`   - (Defaults to 5 minutes) Used when retrieving the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

Custom Provider can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_custom_provider.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.CustomProviders/resourceProviders/example
```
