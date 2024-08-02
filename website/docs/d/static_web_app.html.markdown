---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_static_web_app"
description: |-
  Gets information about an existing Static Web App.
---

# Data Source: azurerm_static_web_app

Use this data source to access information about an existing Static Web App.

## Example Usage

```hcl
data "azurerm_static_web_app" "example" {
  name                = "existing"
  resource_group_name = "existing"
}


```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Static Web App.

* `resource_group_name` - (Required) The name of the Resource Group where the Static Web App exists.

* `location` - The Azure region in which this Static Web App exists.

* `api_key` - The API key of this Static Web App, which is used for later interacting with this Static Web App from other clients, e.g. GitHub Action.

* `app_settings` - The map of key-value pairs of App Settings for the Static Web App.

* `basic_auth` - A `basic_auth` block as defined below.

* `configuration_file_changes_enabled` - Are changes to the configuration file permitted. 

* `default_host_name` - The default host name of the Static Web App.

* `preview_environments_enabled` - Are Preview (Staging) environments enabled. 

* `public_network_access_enabled` - (Optional) Should public network access be enabled for the Static Web App. Defaults to `true`.

* `sku_tier` - The SKU tier of the Static Web App.

* `sku_size` - The SKU size of the Static Web App.

* `identity` - An `identity` block as defined below.

* `tags` - The mapping of tags assigned to the resource.

--- 

An `identity` block exports the following:

* `type` - The Type of Managed Identity assigned to this Static Web App resource.

* `identity_ids` - The list of Managed Identity IDs which are assigned to this Static Web App resource.

---

A `basic_auth` block exports the following:

* `environments` - The Environment types which are configured to use Basic Auth access.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Static Web App
