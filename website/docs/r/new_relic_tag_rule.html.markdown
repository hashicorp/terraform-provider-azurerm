---
subcategory: "New Relic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_new_relic_tag_rule"
description: |-
  Manages an Azure Native New Relic Tag Rule.
---

# azurerm_new_relic_tag_rule

Manages an Azure Native New Relic Tag Rule.

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
}

resource "azurerm_new_relic_tag_rule" "example" {
  monitor_id                         = azurerm_new_relic_monitor.example.id
  azure_active_directory_log_enabled = true
  activity_log_enabled               = true
  metric_enabled                     = true
  subscription_log_enabled           = true

  log_tag_filter {
    name   = "key"
    action = "Include"
    value  = "value"
  }

  metric_tag_filter {
    name   = "key"
    action = "Exclude"
    value  = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `monitor_id` - (Required) Specifies the ID of the New Relic Monitor this Tag Rule should be created within. Changing this forces a new Azure Native New Relic Tag Rule to be created.

* `azure_active_directory_log_enabled` - (Optional) Whether Azure Active Directory logs should be sent for the Monitor resource. Defaults to `false`.

* `activity_log_enabled` - (Optional) Whether activity logs from Azure resources should be sent for the Monitor resource. Defaults to `false`.

* `log_tag_filter` - (Optional) A `log_tag_filter` block as defined below.

* `metric_enabled` - (Optional) Whether metrics should be sent for the Monitor resource. Defaults to `false`.

* `metric_tag_filter` - (Optional) A `metric_tag_filter` block as defined below.

* `subscription_log_enabled` - (Optional) Whether subscription logs should be sent for the Monitor resource. Defaults to `false`.

---

A `log_tag_filter` block supports the following:

* `name` - (Required) Specifies the name (also known as the key) of the tag.

* `action` - (Required) Valid actions for a filtering tag. Possible values are `Exclude` and `Include`. Exclusion takes priority over inclusion.

* `value` - (Required) Specifies the value of the tag.

---

A `metric_tag_filter` block supports the following:

* `name` - (Required) Specifies the name (also known as the key) of the tag.

* `action` - (Required) Valid actions for a filtering tag. Possible values are `Exclude` and `Include`. Exclusion takes priority over inclusion.

* `value` - (Required) Specifies the value of the tag.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Native New Relic Tag Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Native New Relic Tag Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Native New Relic Tag Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Native New Relic Tag Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Native New Relic Tag Rule.

## Import

Azure Native New Relic Tag Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_new_relic_tag_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/NewRelic.Observability/monitors/monitor1/tagRules/ruleSet1
```
