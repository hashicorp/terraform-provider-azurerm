---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_policy"
description: |-
    Lists API Management API Policy resources.
---

# List resource: azurerm_api_management_api_policy

Lists API Management API Policy resources.

## Example Usage

### List API Management API Policies in an API Management API

```hcl
list "azurerm_api_management_api_policy" "example" {
  provider = azurerm
  config {
    api_management_api_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/apis/api1;rev=1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_management_api_id` - (Required) The ID of the API Management API to query.
