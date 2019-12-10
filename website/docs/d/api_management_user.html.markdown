---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_user"
sidebar_current: "docs-azurerm-datasource-api-management-user"
description: |-
  Gets information about an existing API Management User.
---

# Data Source: azurerm_api_management_user

Use this data source to access information about an existing API Management User.

## Example Usage

```hcl
data "azurerm_api_management_user" "example" {
  user_id             = "my-user"
  api_management_name = "example-apim"
  resource_group_name = "search-service"
}

output "notes" {
  value = "${data.azurerm_api_management_user.example.notes}"
}
```

## Argument Reference

* `api_management_name` - (Required) The Name of the API Management Service in which this User exists.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management Service exists.

* `user_id` - (Required) The Identifier for the User.

## Attributes Reference

* `id` - The ID of the API Management User.

* `first_name` - The First Name for the User.

* `last_name` - The Last Name for the User.

* `email` - The Email Address used for this User.

* `note` - Any notes about this User.

* `state` - The current state of this User, for example `active`, `blocked` or `pending`.
