---
subcategory: "DataShare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_account"
description: |-
  Gets information about an existing DataShare Account
---

# Data Source: azurerm_data_share_account

Uses this data source to access information about an existing DataShare Account
---

## DataShare Account Usage

```hcl
data "azurerm_data_share_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resource-group"
}

output "data_share_account_id" {
  value = data.azurerm_data_share_account.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the DataShare Account.

* `resource_group_name` - (Required) Specifies the name of the resource group the DataShare Account is located in.

## Attributes Reference

The following attributes are exported:

* `id` - The Data Share Account ID.

* `tags` - A mapping of tags to assign to the resource.

## Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 & 1.44 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DataShare Account.
