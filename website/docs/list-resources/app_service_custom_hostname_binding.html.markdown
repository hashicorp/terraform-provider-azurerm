---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_custom_hostname_binding"
description: |-
    Lists App Service Custom Hostname Binding resources.
---

# List resource: azurerm_app_service_custom_hostname_binding

Lists App Service Custom Hostname Binding resources.

## Example Usage

### List all App Service Custom Hostname Bindings in an App Service

```hcl
list "azurerm_app_service_custom_hostname_binding" "example" {
  provider = azurerm
  config {
    app_service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Web/sites/example-app-service"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `app_service_id` - (Required) The ID of the App Service to query.
