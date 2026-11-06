---
subcategory: "Api"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_deployment"
description: |-
    Lists Api Center Deployment resources.
---

# List resource: azurerm_api_center_deployment

Lists Api Center Deployment resources.

## Example Usage

### List Api Center Deployments in a Api Center Api

```hcl
list "azurerm_api_center_deployment" "example" {
  provider = azurerm
  config {
    api_center_api_id = "..."
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_center_api_id` - (Required) The ID of the Api Center Api to query.
