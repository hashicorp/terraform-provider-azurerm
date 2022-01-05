---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_tag_rule"
description: |-
  Manages an elastic Tag Rule.
---

# azurerm_elastic_tag_rule

Manages a elastic Tag Rule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "example-rg"
  location = "%s"
}
resource "azurerm_elastic_monitor" "test" {
 name = "example_elastic_monitor"
 resource_group_name = azurerm_resource_group.test.name
 location = azurerm_resource_group.test.location
 user_info {
  email_address = "abc@microsoft.com"
 }
 sku {
  name = "staging_Monthly"
 }
}
resource "azurerm_elastic_tag_rule" "test" {
  name = azurerm_elastic_monitor.test.name
  resource_group_name = azurerm_elastic_monitor.test.resource_group_name
  log_rules{
   send_subscription_logs = true
   send_activity_logs = true
   filtering_tag {
    name = "Test"
    value = "Terraform"
    action = "Exclude"
   }
  }
 }
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Elastic Monitor. Changing this forces a new elastic Tag Rule to be created.

* `resource_group_name` - (Required) The resource group of of the Elastic Monitor assosiated with the rule. Changing this forces a new elastic Tag Rule to be created.

* `rule_set_name` - (Optional)Default value is `default`. The name of the rule. Changing this forces a new elastic Tag Rule to be created.

---

* `filtering_tag` - (Optional) One or more `filtering_tag` blocks as defined below.This only takes effect if `send_activity_logs` flag is enabled. If empty, all resources will be captured. If only Exclude action is specified, the rules will apply to the list of all available resources. If Include actions are specified, the rules will only include resources with the associated tags.

* `send_aad_logs` - (Optional) Whether AAD logs should be sent to the Monitor resource?

* `send_activity_logs` - (Optional) Whether activity logs from Azure resources should be sent to the Monitor resource?

* `send_subscription_logs` - (Optional) Whether subscription logs should be sent to the Monitor resource?

---

An `filtering_tag` block exports the following:

* `name` - (Required) The name of this `filtering_tag`.

* `action` - (Required) The action for a filtering tag. Possible values are "Include" and "Exclude" is allowed. Note that the `Exclude` takes priority over the `Include`.

* `value` - (Optional) The value of this `filtering_tag`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the elastic Tag Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the elastic Tag Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the elastic Tag Rule.
* `update` - (Defaults to 30 minutes) Used when updating the elastic Tag Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the elastic Tag Rule.

## Import

elastic Tag Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_tag_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Elastic/monitors/monitor1/tagRules/ruleSet1
```
