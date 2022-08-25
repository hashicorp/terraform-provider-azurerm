---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_group"
description: |-
  Gets information about an existing API Management Group.
---

# Data Source: azurerm_api_management_group

Use this data source to access information about an existing API Management Group.

## Example Usage

```hcl
data "azurerm_api_management_group" "example" {
  name                = "my-group"
  api_management_name = "example-apim"
  resource_group_name = "search-service"
}

output "group_type" {
  value = data.azurerm_api_management_group.example.type
}
```

## Argument Reference

* `api_management_name` - The Name of the API Management Service in which this Group exists.

* `name` - The Name of the API Management Group.

* `resource_group_name` - The Name of the Resource Group in which the API Management Service exists.

## Attributes Reference

* `id` - The ID of the API Management Group.

* `display_name` - The display name of this API Management Group.

* `description` - The description of this API Management Group.

* `external_id` - The identifier of the external Group.

* `type` - The type of this API Management Group, such as `custom` or `external`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Group.
