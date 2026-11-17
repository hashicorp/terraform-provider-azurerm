---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_operation"
description: |-
    Lists API Management API Operation resources.
---

# List resource: azurerm_api_management_api_operation

Lists API Management API Operation resources.

## Example Usage

### List all API Management API Operations in an API

```hcl
list "azurerm_api_management_api_operation" "example" {
  provider = azurerm
  config {
    api_management_api_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiManagement/service/example-apim/apis/example-api"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_management_api_id` - (Required) The ID of the API Management API to query.
