---
subcategory: "ApiCenter"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_metadata_schema"
description: |-
    Lists Api Center Metadata Schema resources.
---

# List resource: azurerm_api_center_metadata_schema

Lists Api Center Metadata Schema resources.

## Example Usage

### List Api Center Metadata Schemas in a Service

```hcl
list "azurerm_api_center_metadata_schema" "example" {
  provider = azurerm
  config {
    api_center_service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiCenter/services/example-service"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_center_service_id` - (Required) The ID of the API Center Service to query.
