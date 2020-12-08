---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share"
description: |-
  Manages a Data Share.
---

# azurerm_data_share

Manages a Data Share.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_share_account" "example" {
  name                = "example-dsa"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    foo = "bar"
  }
}

resource "azurerm_data_share" "example" {
  name        = "example_dss"
  account_id  = azurerm_data_share_account.example.id
  kind        = "CopyBased"
  description = "example desc"
  terms       = "example terms"

  snapshot_schedule {
    name       = "example-ss"
    recurrence = "Day"
    start_time = "2020-04-17T04:47:52.9614956Z"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Share. Changing this forces a new Data Share to be created.

* `account_id` - (Required) The ID of the Data Share account in which the Data Share is created. Changing this forces a new Data Share to be created.

* `kind` - (Required) The kind of the Data Share. Possible values are `CopyBased` and `InPlace`. Changing this forces a new Data Share to be created.

* `description` - (Optional) The Data Share's description.

* `snapshot_schedule` - (Optional) A `snapshot_schedule` block as defined below.

* `terms` - (Optional) The terms of the Data Share.

---

A `snapshot_schedule` block supports the following:

* `name` - The name of the snapshot schedule.

* `recurrence` - (Required) The interval of the synchronization with the source data. Possible values are `Hour` and `Day`.

* `start_time` - (Required) The synchronization with the source data's start time.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share.
* `update` - (Defaults to 30 minutes) Used when updating the Data Share.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share.

## Import

Data Shares can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1
```
