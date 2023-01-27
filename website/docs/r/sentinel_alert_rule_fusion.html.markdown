---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_alert_rule_fusion"
description: |-
  Manages a Sentinel Fusion Alert Rule.
---

# azurerm_sentinel_alert_rule_fusion

Manages a Sentinel Fusion Alert Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}

resource "azurerm_sentinel_alert_rule_fusion" "example" {
  name                       = "example-fusion-alert-rule"
  log_analytics_workspace_id = azurerm_log_analytics_solution.example.workspace_resource_id
  alert_rule_template_guid   = "f71aba3d-28fb-450b-b192-4e76a83015c8"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel Fusion Alert Rule. Changing this forces a new Sentinel Fusion Alert Rule to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this Sentinel Fusion Alert Rule belongs to. Changing this forces a new Sentinel Fusion Alert Rule to be created.

* `alert_rule_template_guid` - (Required) The GUID of the alert rule template which is used for this Sentinel Fusion Alert Rule. Changing this forces a new Sentinel Fusion Alert Rule to be created.

* `enabled` - (Optional) Should this Sentinel Fusion Alert Rule be enabled? Defaults to `true`.

* `source` - (Optional) One or more `source` blocks as defined below.

---

A `source` block supports the following:

* `name` - (Required) The name of the Fusion source signal. Refer to Fusion alert rule template for supported values.

* `enabled` - (Optional) Whether this source signal is enabled or disabled in Fusion detection? Defaults to `true`.

* `sub_type` - (Optional) One or more `sub_type` blocks as defined below.

---

A `sub_type` block supports the following:

* `name` - (Required) The Name of the source subtype under a given source signal in Fusion detection. Refer to Fusion alert rule template for supported values.

* `enabled` - (Optional) Whether this source subtype under source signal is enabled or disabled in Fusion detection. Defaults to `true`.

* `severities_allowed` - (Required) A list of severities that are enabled for this source subtype consumed in Fusion detection. Possible values for each element are `High`, `Medium`, `Low`, `Informational`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Sentinel Fusion Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Fusion Alert Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Fusion Alert Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Fusion Alert Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Fusion Alert Rule.

## Import

Sentinel Fusion Alert Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_alert_rule_fusion.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1
```
