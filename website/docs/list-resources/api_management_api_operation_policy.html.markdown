---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_operation_policy"
description: |-
    Lists API Management API Operation Policy resources.
---

# List resource: azurerm_api_management_api_operation_policy

Lists API Management API Operation Policy resources.

## Example Usage

### List API Management API Operation Policies in an API Operation

```hcl
list "azurerm_api_management_api_operation_policy" "example" {
  provider = azurerm
  config {
    api_management_api_operation_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiManagement/service/example-apim/apis/example-api/operations/example-operation"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_management_api_operation_id` - (Required) The ID of the API Management API Operation to query.
