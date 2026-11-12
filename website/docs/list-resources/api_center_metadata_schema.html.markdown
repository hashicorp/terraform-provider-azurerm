---
subcategory: "Api"
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
    service_id = "..."
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `service_id` - (Required) The ID of the Service to query.
