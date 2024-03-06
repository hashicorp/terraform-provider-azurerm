---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cost_anomaly_alert"
description: |-
  Manages a Cost Anomaly Alert.
---

# azurerm_cost_anomaly_alert

Manages a Cost Anomaly Alert.

~> **Note:** Anomaly alerts are sent based on the current access of the rule creator at the time that the email is sent. Learn more [here](https://learn.microsoft.com/en-us/azure/cost-management-billing/understand/analyze-unexpected-charges#create-an-anomaly-alert).

## Example Usage

```hcl
resource "azurerm_cost_anomaly_alert" "example" {
  name            = "alertname"
  display_name    = "Alert DisplayName"
  subscription_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
  email_subject   = "My Test Anomaly Alert"
  email_addresses = ["example@test.net"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cost Anomaly Alert. Changing this forces a new resource to be created. The name can contain only lowercase letters, numbers and hyphens.

* `display_name` - (Required) The display name which should be used for this Cost Anomaly Alert.

* `subscription_id` - (Optional) The ID of the Subscription this Cost Anomaly Alert is scoped to. Changing this forces a new resource to be created. When not supplied this defaults to the subscription configured in the provider.

* `email_addresses` - (Required) Specifies a list of email addresses which the Anomaly Alerts are send to.

* `email_subject` - (Required) The email subject of the Cost Anomaly Alerts. Maximum length of the subject is 70.



---

* `message` - (Optional) The message of the Cost Anomaly Alert. Maximum length of the message is 250.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cost Anomaly Alert.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cost Anomaly Alert.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cost Anomaly Alert.
* `update` - (Defaults to 30 minutes) Used when updating the Cost Anomaly Alert.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cost Anomaly Alert.

## Import

Cost Anomaly Alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cost_anomaly_alert.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/scheduledActions/dailyanomalybyresourcegroup
```
