---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cost_anomaly_alert"
description: |-
  Manages a Cost Management Anomaly Alert.
---

# azurerm_cost_anomaly_alert

Manages a Cost Management Anomaly Alert.

## Example Usage

```hcl
resource "azurerm_cost_anomaly_alert" "example" {
  name            = "alertname"
  display_name    = "Alert DisplayName"
  email_subject   = "My Test Anomaly Alert"
  email_addresses = ["example@test.net"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cost Management Anomaly Alert. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name which should be used for this Cost Management Anomaly Alert.

* `email_addresses` - (Required) Specifies a list of email addresses which the Anomaly Alerts are send to.

* `email_subject` - (Required) The email subject of the Anomaly Alerts.



---

* `message` - (Optional) The message of the Anomaly Alert.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cost Management Anomaly Alert.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cost Management Anomaly Alert.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cost Management Anomaly Alert.
* `update` - (Defaults to 30 minutes) Used when updating the Cost Management Anomaly Alert.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cost Management Anomaly Alert.

## Import

Cost Management Anomaly Alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cost_anomaly_alert.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/scheduledActions/dailyanomalybyresourcegroup
```
