---
subcategory: "ApiCenter"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_api"
description: |-
    Lists Api Center Api resources in the default workspace of an API Center Service.
---

# List resource: azurerm_api_center_api

Lists Api Center Api resources.

## Example Usage

### List Api Center Apis in an API Center Service

```hcl
list "azurerm_api_center_api" "example" {
  provider = azurerm
  config {
    api_center_service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiCenter/services/example-service"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_center_service_id` - (Required) The ID of the API Center Service whose default workspace to query.
