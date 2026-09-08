---
subcategory: "Analysis Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_analysis_services_server"
description: |-
    Lists Analysis Services Server resources.
---

# List resource: azurerm_analysis_services_server

Lists Analysis Services Server resources.

## Example Usage

### List all Analysis Services Servers in the subscription

```hcl
list "azurerm_analysis_services_server" "example" {
  provider = azurerm
  config {}
}
```

### List all Analysis Services Servers in a specific resource group

```hcl
list "azurerm_analysis_services_server" "example" {
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
