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

!> **Note:** This Data Source is coming in version 3.0 of the Azure Provider and is available **as an opt-in Beta** - more information can be found in [the upcoming version 3.0 of the Azure Provider](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/3.0-overview). 


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

* `id` - The ID of the App Service Source Github Token.

* `token` - The GitHub Token value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Source GitHub Token.
