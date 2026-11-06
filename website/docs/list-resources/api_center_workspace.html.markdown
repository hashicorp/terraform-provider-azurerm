---
subcategory: "Api"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_workspace"
description: |-
    Lists Api Center Workspace resources.
---

# List resource: azurerm_api_center_workspace

Lists Api Center Workspace resources.

## Example Usage

### List API Center Workspaces in an API Center Service

```hcl
list "azurerm_api_center_workspace" "example" {
  provider = azurerm
  config {
    api_center_service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.ApiCenter/services/example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_center_service_id` - (Required) The ID of the Service to query.
