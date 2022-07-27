---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_batch_configuration"
description: |-
  Manages a Logic App Integration Account Batch Configuration.
---

# azurerm_logic_app_integration_account_batch_configuration

Manages a Logic App Integration Account Batch Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "example" {
  name                = "example-ia"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_logic_app_integration_account_batch_configuration" "example" {
  name                     = "exampleiabc"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name
  batch_group_name         = "TestBatchGroup"

  release_criteria {
    message_count = 80
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Batch Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Batch Configuration should exist. Changing this forces a new resource to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new resource to be created.

* `batch_group_name` - (Required) The batch group name of the Logic App Integration Batch Configuration. Changing this forces a new resource to be created.

* `release_criteria` - (Required) A `release_criteria` block as documented below, which is used to select the criteria to meet before processing each batch.

* `metadata` - (Optional) A JSON mapping of any Metadata for this Logic App Integration Account Batch Configuration.

---

A `release_criteria` block exports the following:

* `batch_size` - (Optional) The batch size in bytes for the Logic App Integration Batch Configuration.

* `message_count` - (Optional) The message count for the Logic App Integration Batch Configuration.

* `recurrence` - (Optional) A `recurrence` block as documented below.

---

A `recurrence` block exports the following:

* `frequency` - (Required) The frequency of the schedule. Possible values are `Day`, `Hour`, `Minute`, `Month`, `Second`, `Week` and `Year`.

* `interval` - (Required) The number of `frequency`s between runs.

* `end_time` - (Optional) The end time of the schedule, formatted as an RFC3339 string.

* `schedule` - (Optional) A `schedule` block as documented below.

* `start_time` - (Optional) The start time of the schedule, formatted as an RFC3339 string.

* `time_zone` - (Optional) The timezone of the start/end time.

---

A `schedule` block exports the following:

* `hours` - (Optional) A list containing a single item, which specifies the Hour interval at which this recurrence should be triggered.

* `minutes` - (Optional) A list containing a single item which specifies the Minute interval at which this recurrence should be triggered.

* `month_days` - (Optional) A list of days of the month that the job should execute on.

* `monthly` - (Optional) A `monthly` block as documented below.

* `week_days` - (Optional) A list of days of the week that the job should execute on. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and `Saturday`.

---

A `monthly` block exports the following:

* `weekday` - (Required) The day of the occurrence. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and `Saturday`.

* `week` - (Required) The occurrence of the week within the month.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Batch Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Batch Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Batch Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Batch Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Batch Configuration.

## Import

Logic App Integration Account Batch Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_batch_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/batchConfigurations/batchConfiguration1
```
