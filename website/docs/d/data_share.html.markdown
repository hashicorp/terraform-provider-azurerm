---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share"
description: |-
  Gets information about an existing Data Share.
---

# Data Source: azurerm_data_share

Use this data source to access information about an existing Data Share.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resource-group"
}

data "azurerm_data_share" "example" {
  name       = "existing"
  account_id = data.azurerm_data_share_account.exmaple.id
}

output "id" {
  value = data.azurerm_data_share.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share.

* `account_id` - (Required) The ID of the Data Share account from which the Data Share is created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Share.

* `share_kind` - The kind of the Data Share.

* `description` - The description of the Data Share.

* `snapshot_schedule` - A `snapshot_schedule` block as defined below.

* `terms` - The terms of the Data Share.

---

A `snapshot_schedule` block exports the following:

* `recurrence` - The recurrence interval of the synchronization of the source data. Possible values are 'Hour'and 'Day'.

* `start_time` - The start time of the synchronization of the source data.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share.
