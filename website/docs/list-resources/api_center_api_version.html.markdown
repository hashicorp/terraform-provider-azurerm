---
subcategory: "Api"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_api_version"
description: |-
    Lists Api Center Api Version resources.
---

# List resource: azurerm_api_center_api_version

Lists Api Center Api Version resources.

## Example Usage

### List Api Center Api Versions in a Api Center Api

```hcl
list "azurerm_api_center_api_version" "example" {
  provider = azurerm
  config {
    api_center_api_id = "..."
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_center_api_id` - (Required) The ID of the Api Center Api to query.
