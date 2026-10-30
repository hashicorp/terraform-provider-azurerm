---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_linux_web_app_site_container"
description: |-
  Lists Linux Web App Site Container resources.
---

# List resource: azurerm_linux_web_app_site_container

Lists Linux Web App Site Container resources for a given Linux Web App.

## Example Usage

### List all Site Containers for a Linux Web App

```hcl
list "azurerm_linux_web_app_site_container" "example" {
  provider = azurerm
  config {
    linux_web_app_id = azurerm_linux_web_app.example.id
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `linux_web_app_id` - (Required) The ID of the Linux Web App whose Site Containers should be listed.
