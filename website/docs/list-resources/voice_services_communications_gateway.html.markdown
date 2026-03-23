---
subcategory: "Voice Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_voice_services_communications_gateway"
description: |-
  Lists Voice Services Communications Gateway resources.
---

# List resource: azurerm_voice_services_communications_gateway

Lists Voice Services Communications Gateway resources.

## Example Usage

### List all Communications Gateways in the subscription

```hcl
list "azurerm_voice_services_communications_gateway" "example" {
  provider = azurerm
  config {}
}
```

### List all Communications Gateways in a specific resource group

```hcl
list "azurerm_voice_services_communications_gateway" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
