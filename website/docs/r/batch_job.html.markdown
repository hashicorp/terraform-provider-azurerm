---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_job"
description: |-
  Manages a Batch Job.
---

# azurerm_batch_job

Manages a Batch Job.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_batch_account" "example" {
  name                = "exampleaccount"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_batch_pool" "example" {
  name                = "examplepool"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_batch_account.example.name
  node_agent_sku_id   = "batch.node.ubuntu 16.04"
  vm_size             = "Standard_A1"

  fixed_scale {
    target_dedicated_nodes = 1
  }

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04.0-LTS"
    version   = "latest"
  }
}

resource "azurerm_batch_job" "example" {
  name          = "examplejob"
  batch_pool_id = azurerm_batch_pool.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `batch_pool_id` - (Required) The ID of the Batch Pool. Changing this forces a new Batch Job to be created.

* `name` - (Required) The name which should be used for this Batch Job. Changing this forces a new Batch Job to be created.

---

* `common_environment_properties` - (Optional) Specifies a map of common environment settings applied to this Batch Job. Changing this forces a new Batch Job to be created.

* `display_name` - (Optional) The display name of this Batch Job. Changing this forces a new Batch Job to be created.

* `task_retry_maximum` - (Optional) The number of retries to each Batch Task belongs to this Batch Job. If this is set to `0`, the Batch service does not retry Tasks. If this is set to `-1`, the Batch service retries Batch Tasks without limit. Default value is `0`.

* `priority` - (Optional) The priority of this Batch Job, possible values can range from -1000 (lowest) to 1000 (highest). Defaults to `0`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Batch Job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Batch Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Job.
* `update` - (Defaults to 5 minutes) Used when updating the Batch Job.
* `delete` - (Defaults to 5 minutes) Used when deleting the Batch Job.

## Import

Batch Jobs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_batch_job.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Batch/batchAccounts/account1/pools/pool1/jobs/job1
```
