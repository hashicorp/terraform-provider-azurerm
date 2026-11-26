---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_custom_policy"
description: |-
    Lists Network DDoS Custom Policy resources.
---

# List resource: azurerm_network_ddos_custom_policy

Lists Network DDoS Custom Policy resources.

## Example Usage

### List all Network DDoS Custom Policies in the subscription

```hcl
list "azurerm_network_ddos_custom_policy" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Network DDoS Custom Policies in a Resource Group

```hcl
list "azurerm_network_ddos_custom_policy" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query. Defaults to the value specified in the Provider Configuration.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
