---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_named_value"
description: |-
    Lists Api Management Named Value resources.
---

# List resource: azurerm_api_management_named_value

Lists Api Management Named Value resources.

## Example Usage

### List Api Management Named Values in an API Management Service

```hcl
list "azurerm_api_management_named_value" "example" {
  provider = azurerm
  config {
    service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiManagement/service/example-service"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `service_id` - (Required) The ID of the API Management Service to query.
