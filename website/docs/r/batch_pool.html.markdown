---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_pool"
sidebar_current: "docs-azurerm-resource-batch-pool"
description: |-
  Manages an Azure Batch pool.

---

# azurerm_batch_pool

Manages an Azure Batch pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "testaccsa"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
  }

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
	name                   		    = "testaccpool"
	resource_group_name           = "${azurerm_resource_group.test.name}"
	account_name 		   		        = "${azurerm_batch_account.test.name}"
	display_name		   		        = "Test Acc Pool Auto"
	vm_size				   		          = "Standard_A1"
	node_agent_sku_id			        = "batch.node.ubuntu 16.04"
	scale_mode			   		        = "Auto"
	autoscale_evaluation_interval = "PT15M"
	autoscale_formula			        = <<EOF
	startingNumberOfVMs = 1;
	maxNumberofVMs = 25;
	pendingTaskSamplePercent = $PendingTasks.GetSamplePercent(180 * TimeInterval_Second);
	pendingTaskSamples = pendingTaskSamplePercent < 70 ? startingNumberOfVMs : avg($PendingTasks.GetSample(180 * TimeInterval_Second));
	$TargetDedicatedNodes=min(maxNumberofVMs, pendingTaskSamples);
	EOF

	storage_image_reference {
        publisher = "Canonical"
        offer     = "UbuntuServer"
        sku       = "16.04.0-LTS"
        version   = "latest"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Batch pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch pool. Changing this forces a new resource to be created.

* `account_name` - (Required) Specifies the name of the Batch account in which the pool will be created. Changing this forces a new resource to be created.

* `node_agent_sku_id` - (Required) Specifies the Sku of the node agents that will be created in the Batch pool.

* `storage_image_reference` - (Required) A `storage_image_reference` for the virtual machines that will compose the Batch pool.

* `display_name` - (Optional) Specifies the display name of the Batch pool.

* `vm_size` - (Optional) Specifies the size of the VM created in the Batch pool. Defaults to `Standard_A1`.

* `scale_mode` - (Optional) Specifies the mode to use for scaling of the Batch pool. Possible values are `Fixed` or `Auto`. Defaults to `Fixed`.

* `target_dedicated_nodes` - (Optional) Specifies the number of nodes to be created in the Batch pool when using scale mode `Fixed`. Defaults to 1.

* `target_low_priority_nodes` - (Optional) Specifies the target number of low priority nodes wanted in the Batch pool when using scale mode `Fixed`. Defaults to 0.

* `resize_timeout` - (Optional) Specifies the time in minutes for the resize operation to timeout. Defaults to `PT15M` (15 minutes).

* `autoscale_evaluation_interval` - (Optional) Specifies a time interval at which to automatically adjust the pool size according to the autoscale formula when using scale mode `Auto`. Defaults to `PT15M` (15 minutes).

* `autoscale_formula` - (Optional) Specifies the formula used to autoscale the pool. The formula is required when using scale mode `Auto`.

## Attributes Reference

The following attributes are exported:

* `id` - The Batch pool ID.
