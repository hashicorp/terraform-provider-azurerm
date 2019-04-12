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

* `start_task` - A `start_task` block that describes the start task settings for the Batch pool.

* `max_tasks_per_node` - The maximum number of tasks that can run concurrently on a single compute node in the pool.

---

A `fixed_scale` block exports the following:

* `target_dedicated_nodes` - The number of nodes in the Batch pool.

* `target_low_priority_nodes` - The number of low priority nodes in the Batch pool.

* `resize_timeout` - The timeout for resize operations.

--- 

A `auto_scale` block exports the following:

* `evaluation_interval` - The interval to wait before evaluating if the pool needs to be scaled.

* `formula` - The autoscale formula that needs to be used for scaling the Batch pool.

---

A `start_task` block exports the following:

* `command_line` - The command line executed by the start task.

* `max_task_retry_count` - The number of retry count.

* `wait_for_success` - A flag that indicates if the Batch pool should wait for the start task to be completed.

* `environment` - A map of strings (key,value) that represents the environment variables to set in the start task.

* `user_identity` - A `user_identity` block that describes the user identity under which the start task runs.

---

A `user_identity` block exports the following:

* `user_name` - The username to be used by the Batch pool start task.

* `auto_user` - A `auto_user` block that describes the user identity under which the start task runs.

---

A `auto_user` block exports the following:

* `elevation_level` - The elevation level of the user identity under which the start task runs.

* `scope` - The scope of the user identity under which the start task runs.
