---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_named_value"
description: |-
  Manages an API Management Named Value.
---

# azurerm_api_management_named_value

Manages an API Management Named Value.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_named_value" "example" {
  name                = "example-apimg"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  display_name        = "ExampleProperty"
  value               = "Example Value"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Named Value. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Named Value should exist. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the [API Management Service](api_management.html) in which the API Management Named Value should exist. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of this API Management Named Value.

* `value` - (Required) The value of this API Management Named Value.

* `secret` - (Optional) Specifies whether the API Management Named Value is secret. Valid values are `true` or `false`. The default value is `false`.

~> **NOTE:** setting the field `secret` to `true` doesn't make this field sensitive in Terraform, instead it marks the value as secret and encrypts the value in Azure.

* `tags` - (Optional) A list of tags to be applied to the API Management Named Value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Named Value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Named Value.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Named Value.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Named Value.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Named Value.

## Import

API Management Properties can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_named_value.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.ApiManagement/service/example-apim/namedValues/example-apimp
```
