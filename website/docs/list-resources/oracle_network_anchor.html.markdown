---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_network_anchor"
description: |-
  Lists Oracle Network Anchor resources.
---

# List resource: azurerm_oracle_network_anchor

Lists Oracle Network Anchor resources.

## Example Usage

### List all Oracle Network Anchors in the subscription

```hcl
list "azurerm_oracle_network_anchor" "example" {
  provider = azurerm
  config {}
}
```

### List all Oracle Network Anchors in a specific resource group

```hcl
list "azurerm_oracle_network_anchor" "example" {
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
