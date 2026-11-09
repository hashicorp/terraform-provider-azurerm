---
subcategory: "Attestation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_attestation_provider"
description: |-
    Lists Attestation Provider resources.
---

# List resource: azurerm_attestation_provider

Lists Attestation Provider resources.

## Example Usage

### List all Attestation Providers in the subscription

```hcl
list "azurerm_attestation_provider" "example" {
  provider = azurerm
  config {}
}
```

### List all Attestation Providers in a specific resource group

```hcl
list "azurerm_attestation_provider" "example" {
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
