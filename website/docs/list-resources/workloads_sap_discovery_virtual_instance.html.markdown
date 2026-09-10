---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_discovery_virtual_instance"
description: |-
  Lists SAP Discovery Virtual Instance resources.
---

# List resource: azurerm_workloads_sap_discovery_virtual_instance

Lists SAP Discovery Virtual Instance resources.

## Example Usage

### List all SAP Discovery Virtual Instances in the subscription

```hcl
list "azurerm_workloads_sap_discovery_virtual_instance" "example" {
  provider = azurerm
  config {}
}
```

### List all SAP Discovery Virtual Instances in a specific resource group

```hcl
list "azurerm_workloads_sap_discovery_virtual_instance" "example" {
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
