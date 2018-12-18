---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_pool"
sidebar_current: "docs-azurerm-datasource-batch-pool"
description: |-
  Get information about an existing Azure Batch pool.

---

# Data source: azurerm_batch_pool

Use this data source to access information about an existing Batch pool

## Example Usage

```hcl
data "azurerm_batch_pool "test" {
  name                = "testbatchpool"
  account_name        = "testbatchaccount"
  resource_group_name = "test"

  fixed_scale {
    target_dedicated_nodes = 2
    resize_timeout         = "PT15M"
  }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The Batch pool ID.

* `name` - The name of the Batch pool.

* `account_name` - The name of the Batch account.

* `node_agent_sku_id` - The Sku of the node agents in the Batch pool.

* `vm_size` - The size of the VM created in the Batch pool.

* `fixed_scale` - A `fixed_scale` block that describes the scale settings when using fixed scale.

* `auto_scale` - A `auto_scale` block that describes the scale settings when using auto scale.

* `storage_image_reference` - The reference of the storage image used by the nodes in the Batch pool.

---

A `fixed_scale` block exports the following:

* `target_dedicated_nodes` - The number of nodes in the Batch pool

* `target_low_priority_nodes` - The number of low priority nodes in the Batch pool

* `resize_timeout` - The timeout for resize operations

--- 

A `auto_scale` block exports the following:

* `evaluation_interval` - The interval to wait before evaluating if the pool needs to be scaled

* `formula` - The autoscale formula that needs to be used for scaling the Batch pool