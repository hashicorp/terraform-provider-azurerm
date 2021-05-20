---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_smart_detector_alert_rule"
description: |-
  Manages an Monitor Smart Detector Alert Rule.
---

# azurerm_monitor_smart_detector_alert_rule

Manages an Monitor Smart Detector Alert Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "example" {
  name                = "example-action-group"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "exampleactiongroup"
}

resource "azurerm_monitor_smart_detector_alert_rule" "example" {
  name                = "example-smart-detector-alert-rule"
  resource_group_name = azurerm_resource_group.example.name
  severity            = "Sev0"
  scope_resource_ids  = [azurerm_application_insights.example.id]
  frequency           = "PT1M"
  detector_type       = "FailureAnomaliesDetector"

  action_group {
    ids = [azurerm_monitor_action_group.test.id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Monitor Smart Detector Alert Rule. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group in which the Monitor Smart Detector Alert Rule should exist. Changing this forces a new resource to be created.

* `detector_type` - (Required) Specifies the Built-In Smart Detector type that this alert rule will use. Currently the only possible value is `FailureAnomaliesDetector`.

* `scope_resource_ids` - (Required) Specifies the scopes of this Smart Detector Alert Rule.

* `action_group` - (Required) An `action_group` block as defined below.

* `severity` - (Required) Specifies the severity of this Smart Detector Alert Rule. Possible values are `Sev0`, `Sev1`, `Sev2`, `Sev3` or `Sev4`.

* `frequency` - (Required) Specifies the frequency of this Smart Detector Alert Rule in ISO8601 format.

* `description` - (Optional) Specifies a description for the Smart Detector Alert Rule.

* `enabled` - (Optional) Is the Smart Detector Alert Rule enabled? Defaults to `true`.

* `throttling_duration` - (Optional) Specifies the duration (in ISO8601 format) to wait before notifying on the alert rule again.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `action_group` block supports the following:

* `ids` - (Required) Specifies the action group ids.

* `email_subject` - (Optional) Specifies a custom email subject if Email Receiver is specified in Monitor Action Group resource.

* `webhook_payload` - (Optional) A JSON String which Specifies the custom webhook payload if Webhook Receiver is specified in Monitor Action Group resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Monitor Smart Detector Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Monitor Smart Detector Alert Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Monitor Smart Detector Alert Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Monitor Smart Detector Alert Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Monitor Smart Detector Alert Rule.

## Import

Monitor Smart Detector Alert Rule can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_monitor_smart_detector_alert_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AlertsManagement/smartdetectoralertrules/rule1
```
