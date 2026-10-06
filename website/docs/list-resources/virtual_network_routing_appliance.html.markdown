---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network_routing_appliance"
description: |-
  Lists Virtual Network Routing Appliance resources.
---

# List resource: azurerm_virtual_network_routing_appliance

Lists Virtual Network Routing Appliance resources.

## Example Usage

```hcl
list "azurerm_virtual_network_routing_appliance" "example" {
  provider         = azurerm
  include_resource = true

  config {
    resource_group_name = "example-resource-group"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query. When omitted, all Resource Groups in the subscription are queried.

* `subscription_id` - (Optional) The ID of the subscription to query. Defaults to the subscription configured in the provider.
