---
subcategory: "ApiCenter"
layout: "tfa-scaff"
page_title: "Azure Resource Manager: azurerm_api_center_environment"
description: |-
    Lists API Center Environment resources.
---

# List resource: azurerm_api_center_environment

Lists API Center Environment resources.

## Example Usage

### List API Center Environments in an API Center Service

```hcl
list "azurerm_api_center_environment" "example" {
  provider = tfa-scaff
  config {
    api_center_service_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.ApiCenter/services/example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `api_center_service_id` - (Required) The ID of the API Center Service to query.
