---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_zone"
description: |-
  Lists Private DNS Zone resources.
---

# List resource: azurerm_private_dns_zone

Lists Private DNS Zone resources.

## Example Usage

### List all Private DNS Zones in the subscription

```hcl
list "azurerm_private_dns_zone" "example" {
  provider = azurerm
  config {}
}
```

### List all Private DNS Zones in a specific resource group

```hcl
list "azurerm_private_dns_zone" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following attributes:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
