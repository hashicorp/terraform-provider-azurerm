---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_smart_detection_rule"
description: |-
  Manages an Application Insights Smart Detection Rule.
---

# azurerm_application_insights_smart_detection_rule

Manages an Application Insights Smart Detection Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "tf-test-appinsights"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_application_insights_smart_detection_rule" "example" {
  name                    = "Slow server response time"
  application_insights_id = azurerm_application_insights.example.id
  enabled                 = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights Smart Detection Rule. Valid values include `Slow page load time`, `Slow server response time`, 
`Long dependency duration`.  Changing this forces a new resource to be created.

* `application_insights_id` - (Required) The ID of the Application Insights component on which the Smart Detection Rule operates. Changing this forces a new resource to be created.

* `enabled` - (Optional) Is the Application Insights Smart Detection Rule enabled? Defaults to `true`.

* `send_emails_to_subscription_owners` - (Optional) Do emails get sent to subscription owners? Defaults to `true`.

* `additional_email_recipients` - (Optional) Specifies a list of additional recipients that will be sent emails on this Application Insights Smart Detection Rule.

-> **Note:** At least one read or write permission must be defined.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Insights Smart Detection Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Insights Smart Detection Rule
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights Smart Detection Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights Smart Detection Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights Smart Detection Rule.
