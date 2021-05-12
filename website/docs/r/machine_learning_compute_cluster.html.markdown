---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_compute_cluster"
description: |-
  Manages a Machine Learning Compute Cluster.
---

# azurerm_machine_learning_compute_cluster

Manages a Machine Learning Compute Cluster.

## Example Usage

```hcl
resource "azurerm_machine_learning_compute_cluster" "example" {
  name = "example"
  location = "West Europe"
  vm_size = "Standard_DS2_v2"

  identity {
    type = "SystemAssigned"    
  }

  scale_settings {
    max_node_count = 1
    min_node_count = 0
    node_idle_time_before_scale_down = "PT30S"    
  }
  machine_learning_workspace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1"
  vm_priority = "LowPriority"
}
```

## Arguments Reference

The following arguments are supported:

* `identity` - (Required) A `identity` block as defined below. Changing this forces a new Machine Learning Compute Cluster to be created.

* `location` - (Required) The Azure Region where the Machine Learning Compute Cluster should exist. Changing this forces a new Machine Learning Compute Cluster to be created.

* `machine_learning_workspace_id` - (Required) The ID of the Machine Learning Workspace. Changing this forces a new Machine Learning Compute Cluster to be created.

* `name` - (Required) The name which should be used for this Machine Learning Compute Cluster. Changing this forces a new Machine Learning Compute Cluster to be created.

* `scale_settings` - (Required) A `scale_settings` block as defined below.

* `vm_priority` - (Required) The priority of the VM. Changing this forces a new Machine Learning Compute Cluster to be created.

* `vm_size` - (Required) The size of the VM. Changing this forces a new Machine Learning Compute Cluster to be created.

---

* `description` - (Optional) The description of the Machine Learning compute.

* `subnet_resource_id` - (Optional) The ID of the Subnet that the Compute Cluster should reside in. Changing this forces a new Machine Learning Compute Cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Compute Cluster.

---

A `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this Disk Encryption Set. At this time the only possible value is SystemAssigned.

---

A `scale_settings` block supports the following:

* `max_node_count` - (Required) Maximum node count.

* `min_node_count` - (Required) Minimal node count.

* `node_idle_time_before_scale_down` - (Required) Node Idle Time Before Scale Down: defines the time until the compute is shutdown when it has gone into Idle state.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Compute Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Compute Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Compute Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Compute Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Compute Cluster.

## Import

Machine Learning Compute Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_compute_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1
```