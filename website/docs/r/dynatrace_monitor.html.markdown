---
subcategory: "Dynatrace"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dynatrace_monitor"
description: |-
  Manages Dynatrace monitors.
---

# azurerm_dynatrace_monitor

Manages a Dynatrace monitor.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dynatrace_monitor" "example" {
  name                            = "exmpledynatracemonitor"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.test.location
  monitoring_enabled              = true
  marketplace_subscription_status = "Active"

  identity {
    type = "SystemAssigned"
  }

  user {
    first_name   = "Alice"
    last_name    = "Bobab"
    email        = "alice@microsoft.com"
    phone_number = "123456"
    country      = "westus"
  }

  plan {
    usage_type     = "COMMITTED"
    billing_cycle  = "MONTHLY"
    plan           = "azureportalintegration_privatepreview@TIDhjdtn7tfnxcy"
    effective_date = "2019-08-30T15:14:33Z"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Dynatrace monitor. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Dynatrace monitor should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Dynatrace monitor should exist. Changing this forces a new resource to be created.

* `identity` - (Required) The kind of managed identity assigned to this resource.  A `identity` block as defined below.

* `marketplace_subscription` - (Required) Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state. Possible values are `Active` and `Suspended`.

* `plan` - (Required) Billing plan information. A `plan` block as defined below. Changing this forces a new resource to be created.

* `user` - (Required) User's information. A `user` block as defined below. Chainging this forces a new resource to be created.

* `monitoring_enabled` - (Optional) Flag specifying if the resource monitoring is enabled or disabled. Default is `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `identity` block supports the following:

* `type` - (Required) The type of identity used for the resource. Only possible value is `SystemAssigned`.

---

A `plan` block supports the following:

* `billing_cycle` - (Optional) Different billing cycles. Possible values are `MONTHLY` and `WEEKLY`.

* `effective_date` - (Required) Date when plan was applied.

* `plan` - (Required) Plan id as published by Dynatrace.

* `usage_type` - (Optional) Different usage type. Possible values are `PAYG` and `COMMITTED`.

---

A `user` block supports the following:

* `country` - (Required) Country of the user.

* `email` - (Required) Email of the user used by Dynatrace for contacting them if needed.

* `first_name` - (Required) First name of the user.

* `last_name` - (Required) Last name of the user.

* `phone_number` - (Required) phone number of the user by Dynatrace for contacting them if needed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dynatrace monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dynatrace monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dynatrace monitor.
* `update` - (Defaults to 30 minutes) Used when updating the Dynatrace monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dynatrace monitor.

## Import

Dynatrace monitor can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dynatrace_monitor.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Dynatrace.Observability/monitors/monitor1
```
