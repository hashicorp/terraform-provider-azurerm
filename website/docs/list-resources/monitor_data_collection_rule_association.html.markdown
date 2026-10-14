---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_data_collection_rule_association"
description: |-
  Lists Data Collection Rule Association resources.
---

# List resource: azurerm_monitor_data_collection_rule_association

Lists Data Collection Rule Association resources.

## Example Usage

### List all Data Collection Rule Associations for a Target Azure Resource

```hcl
list "azurerm_monitor_data_collection_rule_association" "example" {
  provider = azurerm
  config {
    target_resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/virtualMachines/example-vm"
  }
}
```

### List all Data Collection Rule Associations for a Data Collection Rule

```hcl
list "azurerm_monitor_data_collection_rule_association" "example" {
  provider = azurerm
  config {
    data_collection_rule_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Insights/dataCollectionRules/example-dcr"
  }
}
```

### List all Data Collection Rule Associations for a Data Collection Endpoint

```hcl
list "azurerm_monitor_data_collection_rule_association" "example" {
  provider = azurerm
  config {
    data_collection_endpoint_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Insights/dataCollectionEndpoints/example-dce"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `data_collection_endpoint_id` - (Optional) The ID of the Data Collection Endpoint to query associations for. Exactly one of `data_collection_endpoint_id`, `data_collection_rule_id`, or `target_resource_id` must be specified.

* `data_collection_rule_id` - (Optional) The ID of the Data Collection Rule to query associations for. Exactly one of `data_collection_endpoint_id`, `data_collection_rule_id`, or `target_resource_id` must be specified.

* `target_resource_id` - (Optional) The ID of the Target Azure Resource to query associations for. Exactly one of `data_collection_endpoint_id`, `data_collection_rule_id`, or `target_resource_id` must be specified.
