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
  account_name		  = "testbatchaccount"
  resource_group_name = "test"
}

output "autoscale_mode" {
  value = "${data.azurerm_batch_pool.test.resource_group_name}"
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The Batch pool ID.

* `name` - The name of the Batch pool.

* `account_name` - The name of the Batch account.

* `node_agent_sku_id` - The Sku of the node agents in the Batch pool.

* `vm_size` - The size of the VM created in the Batch pool.

* `scale_mode` - The mode to use for scaling of the Batch pool. Possible values are `Fixed` or `Auto`.

* `target_dedicated_nodes` - The number of nodes to be created in the Batch pool when using scale mode `Fixed`.

* `target_low_priority_nodes` - The target number of low priority nodes wanted in the Batch pool when using scale mode `Fixed`.

* `resize_timeout` - The time in minutes for the resize operation to timeout.

* `autoscale_evaluation_interval` - A time interval at which to automatically adjust the pool size according to the autoscale formula when using scale mode `Auto`.

* `autoscale_formula` - The formula used to autoscale the pool when using scale mode `Auto`.

* `storage_image_reference` - The reference of the storage image used by the nodes in the Batch pool.
