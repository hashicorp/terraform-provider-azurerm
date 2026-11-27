---
subcategory: "SRE Agents"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sre_agent"
description: |-
  Lists SRE Agent resources.
---

# List resource: azurerm_sre_agent

Lists SRE Agent resources.

## Example Usage

### List all SRE Agents in the subscription

```hcl
list "azurerm_sre_agent" "example" {
  provider = azurerm
  config {}
}
```

### List all SRE Agents in a specific Resource Group

```hcl
list "azurerm_sre_agent" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-resource-group"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
