---
subcategory: "Api"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_api_definition"
description: |-
    Lists Api Center Api Definition resources.
---

# List resource: azurerm_api_center_api_definition

Lists Api Center Api Definition resources.

## Example Usage

### List Api Center Api Definitions in a Api Center Api Version

```hcl
list "azurerm_api_center_api_definition" "example" {
  provider = azurerm
  config {
    api_center_api_version_id = "..."
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_center_api_version_id` - (Required) The ID of the Api Center Api Version to query.
