---
subcategory: "Datadog"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_datadog_monitor_sso_configuration"
description: |-
  Manages SingleSignOn on the datadog Monitor.
---

# azurerm_datadog_monitor_sso_configuration

Manages SingleSignOn on the datadog Monitor.

## Example Usage

### Enabling SSO on monitor
```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-datadog"
  location = "West US 2"
}
resource "azurerm_datadog_monitor_sso_configuration" "test" {
    datadog_monitor_id = azurerm_datadog_monitor.example.id
    singlesignon_state = "Enable"
    enterprise_application_id = "XXXX"
}
```

## Arguments Reference

The following arguments are supported:

* `datadog_monitor_id` - (Required) The Datadog Monitor Id which should be used for this.

* `singlesignon_state` - (Required) The state of SingleSignOn configuration.

* `enterprise_application_id` - (Required) The application Id to perform SSO operation.

--- 

* `name` - (Optional) The name of the SingleSignOn configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `singlesignon_url` - The SingleSignOn URL to login to Datadog org.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SingleSignOn on the datadog Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the SingleSignOn on the datadog Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the SingleSignOn on the datadog Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the SingleSignOn on the datadog Monitor.