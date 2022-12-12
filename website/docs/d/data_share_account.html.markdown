---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_account"
description: |-
  Gets information about an existing Data Share Account.
---

# Data Source: azurerm_data_share_account

Use this data source to access information about an existing Data Share Account.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_data_share_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share Account.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Share Account exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Share Account.

* `identity` - An `identity` block as defined below.

* `tags` - A mapping of tags assigned to the Data Share Account.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

* `type` - The identity type of this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Account.
