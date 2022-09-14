---
subcategory: "Dynatrace"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dynatrace_monitors"
description: |-
  Manages Dynatrace monitors.
---

# azurerm_dynatrace_monitors

Manages Dynatrace monitors.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dynatrace_monitors" "example" {
  name                            = "exmpledynatracemonitor"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.test.location
  identity_type                   = "SystemAssigned"
  monitoring_status               = "Enabled"
  marketplace_subscription_status = "Active"

  user_info {
    first_name    = "Alice"
    last_name     = "Bobab"
    email_address = "alice@microsoft.com"
    phone_number  = "123456"
    country       = "westus"
  }

  plan_data {
    usage_type     = "COMMITTED"
    billing_cycle  = "MONTHLY"
    plan_details   = "azureportalintegration_privatepreview@TIDhjdtn7tfnxcy"
    effective_date = "2019-08-30T15:14:33Z"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Dynatrace monitor. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Dynatrace monitor should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Dynatrace monitor should exist. Changing this forces a new resource to be created.

* `identity_type` - (Optional) The kind of managed identity assigned to this resource. Possible values are `SystemAssigned`, `UserAssigned`. Changing this forces a new resource to be created.

* `monitoring_status` - (Optional) Flag specifying if the resource monitoring is enabled or disabled. Possible values aree `Enabled`, `Disabled`.

* `marketplace_subscription_status` - (Optional) Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state. Possible values are `Active`, `Suspended`.

* `plan_data` - (Optional) Billing plan information. A `plan_data` block as defined below.

* `user_info` - (Optional) User's information. A `user_info` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `plan_data` block supports the following:

* `billing_cycle` - (Optional) Different billing cycles. Possible values are `MONTHLY`, `WEEKLY`.

* `effective_date` - (Optional) Date when plan was applied.

* `plan_details` - (Optional) Plan id as published by Dynatrace.

* `usage_type` - (Optional) Different usage type. Possible values are `PAYG`, `COMMITTED`.

---

A `user_info` block supports the following:

* `country` - (Optional) Country of the user.

* `email_address` - (Optional) Email of the user used by Dynatrace for contacting them if needed.

* `first_name` - (Optional) First name of the user.

* `last_name` - (Optional) Last name of the user.

* `phone_number` - (Optional) phone number of the user by Dynatrace for contacting them if needed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dynatrace monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Dynatrace monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dynatrace monitor.
* `update` - (Defaults to 1 hour) Used when updating the Dynatrace monitor.
* `delete` - (Defaults to 1 hour) Used when deleting the Dynatrace monitor.

## Import

Dynatrace monitor can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dynatrace_monitors.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Dynatrace.Observability/monitors/monitor1
```
