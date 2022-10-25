---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_source_control_token"
description: |-
  Gets information about an existing App Service GitHub Token.
---

# Data Source: azurerm_source_control_token

Use this data source to access information about an existing App Service Source Control Token.

~> **Note:** This value can only be queried for the user or service principal that is executing Terraform. It is not possible to retrieve for another user.

## Example Usage

```hcl
data "azurerm_source_control_token" "example" {
  type = "GitHub"
}

output "id" {
  value = data.azurerm_app_service_github_token.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `type` - (Required) The Token type. Possible values include `Bitbucket`, `Dropbox`, `Github`, and `OneDrive`.

## Attributes Reference

The following Attributes are exported:

* `id` - The ID of the App Service Source GitHub Token.

* `token` - The GitHub Token value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Source GitHub Token.
