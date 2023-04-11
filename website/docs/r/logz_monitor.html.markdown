---
subcategory: "Logz"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logz_monitor"
description: |-
  Manages a logz Monitor.
---

# azurerm_logz_monitor

Manages a logz Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-logz"
  location = "West Europe"
}

resource "azurerm_logz_monitor" "example" {
  name                = "example-monitor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "2022-06-06T00:00:00Z"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "user@example.com"
    first_name   = "Example"
    last_name    = "User"
    phone_number = "+12313803556"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this logz Monitor. Changing this forces a new logz Monitor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the logz Monitor should exist. Changing this forces a new logz Monitor to be created.

* `location` - (Required) The Azure Region where the logz Monitor should exist. Changing this forces a new logz Monitor to be created.

* `plan` - (Required) A `plan` block as defined below. Changing this forces a new resource to be created.

* `user` - (Required) A `user` block as defined below. Changing this forces a new resource to be created.

---

* `company_name` - (Optional) Name of the Logz organization. Changing this forces a new logz Monitor to be created.

* `enterprise_app_id` - (Optional) The ID of the Enterprise App. Changing this forces a new logz Monitor to be created.

~> **NOTE** Please follow [Set up Logz.io single sign-on](https://docs.microsoft.com/azure/partner-solutions/logzio/setup-sso) to create the ID of the Enterprise App.

* `enabled` - (Optional) Whether the resource monitoring is enabled? Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the logz Monitor.

---

An `plan` block exports the following:

* `billing_cycle` - (Required) Different billing cycles. Possible values are `MONTHLY` or `WEEKLY`. Changing this forces a new logz Monitor to be created.

* `effective_date` - (Required) Date when plan was applied. Changing this forces a new logz Monitor to be created.

* `usage_type` - (Required) Different usage types. Possible values are `PAYG` or `COMMITTED`. Changing this forces a new logz Monitor to be created.

* `plan_id` - (Optional) Plan id as published by Logz. The only possible value is `100gb14days`. Defaults to `100gb14days`. Changing this forces a new logz Monitor to be created.

---

An `user` block exports the following:

* `email` - (Required) Email of the user used by Logz for contacting them if needed. Changing this forces a new logz Monitor to be created.

~> **NOTE** If you use the Azure CLI to authenticate to Azure, the Email of your Azure account needs to be granted the admin permission in your Logz.io account. Otherwise, you may not be able to delete this resource. There is no such limitation for the Service Principal authentication.

* `first_name` - (Required) First Name of the user. Changing this forces a new logz Monitor to be created.

* `last_name` - (Required) Last Name of the user. Changing this forces a new logz Monitor to be created.

* `phone_number` - (Required) Phone number of the user used by Logz for contacting them if needed. Changing this forces a new logz Monitor to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the logz Monitor.

* `single_sign_on_url` - The single sign on url associated with the logz organization of this logz Monitor.

* `logz_organization_id` - The ID associated with the logz organization of this logz Monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the logz Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the logz Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the logz Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the logz Monitor.

## Import

logz Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logz_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logz/monitors/monitor1
```
