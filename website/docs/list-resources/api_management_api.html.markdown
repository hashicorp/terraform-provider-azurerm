---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api"
description: |-
    Lists API Management API resources.
---

# List resource: azurerm_api_management_api

Lists API Management API resources.

## Example Usage

### List API Management APIs in an API Management Service

```hcl
list "azurerm_api_management_api" "example" {
  provider = azurerm
  config {
    api_management_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiManagement/service/example-apim"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_management_id` - (Required) The ID of the API Management Service to query.
