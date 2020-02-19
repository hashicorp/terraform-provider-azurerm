---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app"
description: |-
  Gets information about an existing Function App.

---

# Data Source: azurerm_function_app

Use this data source to access information about a Function App.

## Example Usage

```hcl
data "azurerm_function_app" "example" {
  name                = "test-azure-functions"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Function App resource.

* `resource_group_name` - The name of the Resource Group where the Function App exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Function App

* `app_service_plan_id` - The ID of the App Service Plan within which to create this Function App.

* `app_settings` - A key-value pair of App Settings.

* `connection_string` - An `connection_string` block as defined below.

* `default_hostname` - The default hostname associated with the Function App.

* `enabled` - Is the Function App enabled?

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses.

* `possible_outbound_ip_addresses` - A comma separated list of outbound IP addresses, not all of which are necessarily in use. Superset of `outbound_ip_addresses`.

---

The `connection_string` supports the following:

* `name` - The name of the Connection String.
* `type` - The type of the Connection String. 
* `value` - The value for the Connection String.

---

The `site_credential` block exports the following:

* `username` - The username which can be used to publish to this App Service
* `password` - The password associated with the username, which can be used to publish to this App Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Function App.
