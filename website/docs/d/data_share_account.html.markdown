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

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory.

* `tenant_id` - The ID of the Azure Active Directory Tenant.

* `type` - The identity type of the Data Share Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Account.
