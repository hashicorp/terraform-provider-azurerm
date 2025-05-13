---
subcategory: "New Relic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_new_relic_monitor"
description: |-
  Manages an Azure Native New Relic Monitor.
---

# azurerm_new_relic_monitor

Manages an Azure Native New Relic Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_new_relic_monitor" "example" {
  name                = "example-nrm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  plan {
    effective_date = "2023-06-06T00:00:00Z"
  }

  user {
    email        = "user@example.com"
    first_name   = "Example"
    last_name    = "User"
    phone_number = "+12313803556"
  }

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Azure Native New Relic Monitor. Changing this forces a new Azure Native New Relic Monitor to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Azure Native New Relic Monitor should exist. Changing this forces a new Azure Native New Relic Monitor to be created.

* `location` - (Required) Specifies the Azure Region where the Azure Native New Relic Monitor should exist. Changing this forces a new Azure Native New Relic Monitor to be created.

* `plan` - (Required) A `plan` block as defined below. Changing this forces a new Azure Native New Relic Monitor to be created.

* `user` - (Required) A `user` block as defined below. Changing this forces a new Azure Native New Relic Monitor to be created.

* `account_creation_source` - (Optional) Specifies the source of account creation. Possible values are `LIFTR` and `NEWRELIC`. Defaults to `LIFTR`. Changing this forces a new Azure Native New Relic Monitor to be created.

* `account_id` - (Optional) Specifies the account id. Changing this forces a new Azure Native New Relic Monitor to be created.

-> **Note:** The value of `account_id` must come from an Azure Native New Relic Monitor instance of another different subscription.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new Azure Native New Relic Monitor to be created.

* `ingestion_key` - (Optional) Specifies the ingestion key of account. Changing this forces a new Azure Native New Relic Monitor to be created.

* `organization_id` - (Optional) Specifies the organization id. Changing this forces a new Azure Native New Relic Monitor to be created.

-> **Note:** The value of `organization_id` must come from an Azure Native New Relic Monitor instance of another different subscription.

* `org_creation_source` - (Optional) Specifies the source of org creation. Possible values are `LIFTR` and `NEWRELIC`. Defaults to `LIFTR`. Changing this forces a new Azure Native New Relic Monitor to be created.

* `user_id` - (Optional) Specifies the user id. Changing this forces a new Azure Native New Relic Monitor to be created.

---

A `plan` block supports the following:

* `effective_date` - (Required) Specifies the date when plan was applied. Changing this forces a new Azure Native New Relic Monitor to be created.

* `billing_cycle` - (Optional) Specifies the billing cycles. Possible values are `MONTHLY`, `WEEKLY` and `YEARLY`. Defaults to `MONTHLY`. Changing this forces a new Azure Native New Relic Monitor to be created.

* `plan_id` - (Optional) Specifies the plan id published by NewRelic. The only possible value is `newrelic-pay-as-you-go-free-live`. Defaults to `newrelic-pay-as-you-go-free-live`. Changing this forces a new Azure Native New Relic Monitor to be created.

* `usage_type` - (Optional) Specifies the usage type. Possible values are `COMMITTED` and `PAYG`. Defaults to `PAYG`. Changing this forces a new Azure Native New Relic Monitor to be created.

---

A `user` block supports the following:

* `email` - (Required) Specifies the user Email. Changing this forces a new Azure Native New Relic Monitor to be created.

* `first_name` - (Required) Specifies the first name. Changing this forces a new Azure Native New Relic Monitor to be created.

* `last_name` - (Required) Specifies the last name. Changing this forces a new Azure Native New Relic Monitor to be created.

* `phone_number` - (Required) Specifies the contact phone number. Changing this forces a new Azure Native New Relic Monitor to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Azure Native New Relic Monitor. The only possible value is `SystemAssigned`. Changing this forces a new Azure Native New Relic Monitor to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Native New Relic Monitor.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Azure Native New Relic Monitor.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Azure Native New Relic Monitor.

## Role Assignment

To enable metrics flow, perform role assignment on the identity created above. `Monitoring reader(43d0d8ad-25c7-4714-9337-8ba259a9fe05)` role is required .

### Role assignment on the monitor created

```hcl
data "azurerm_subscription" "primary" {}

data "azurerm_role_definition" "monitoring_reader" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment" "example" {
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.monitoring_reader.id}"
  principal_id       = azurerm_new_relic_monitor.example.identity[0].principal_id
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Native New Relic Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Native New Relic Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Native New Relic Monitor.

## Import

Azure Native New Relic Monitor can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_new_relic_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/NewRelic.Observability/monitors/monitor1
```
