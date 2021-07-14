---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_app_service_github_token"
description: |-
  Gets information about an existing App Service GitHub Token.
---

# Data Source: azurerm_app_service_github_token

Use this data source to access information about an existing App Service GitHub Token. 
~> **Note:** This value can only be queried for the user or service principal that is executing Terraform. It is not possible to retrieve for another user.

## Example Usage

```hcl
data "azurerm_app_service_github_token" "example" {

}

output "id" {
  value = data.azurerm_app_service_github_token.example.id
}
```

## Arguments Reference

The data source does not accept any arguments.


## Attributes Reference

The following Attributes are exported: 

* `id` - The ID of the App Service Source Github Token.

* `token` - The GitHub Token value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Source GitHub Token.