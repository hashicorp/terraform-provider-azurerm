---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_kubernetes_cluster_node_pool"
description: |-
  Gets information about an existing Kubernetes Cluster Node Pool.
---

# Data Source: azurerm_kubernetes_cluster_node_pool

Use this data source to access information about an existing Kubernetes Cluster Node Pool.

## Example Usage

```hcl
data "azurerm_kubernetes_cluster_node_pool" "example" {
  name                    = "existing"
  kubernetes_cluster_name = "existing-cluster"
  resource_group_name     = "existing-resource-group"
}

output "id" {
  value = data.azurerm_kubernetes_cluster_node_pool.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `kubernetes_cluster_name` - (Required) The Name of the Kubernetes Cluster where this Node Pool is located.

* `name` - (Required) The name of this Kubernetes Cluster Node Pool.

* `resource_group_name` - (Required) The name of the Resource Group where the Kubernetes Cluster exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Kubernetes Cluster Node Pool.

* `availability_zones` - A list of Availability Zones in which the Nodes in this Node Pool exists.

* `enable_auto_scaling` - Does this Node Pool have Auto-Scaling enabled?

* `enable_node_public_ip` - Do nodes in this Node Pool have a Public IP Address?

* `eviction_policy` - The eviction policy used for Virtual Machines in the Virtual Machine Scale Set, when `priority` is set to `Spot`.

* `max_count` - The maximum number of Nodes allowed when auto-scaling is enabled.

* `max_pods` - The maximum number of Pods allowed on each Node in this Node Pool.

* `min_count` - The minimum number of Nodes allowed when auto-scaling is enabled.

* `mode` - The Mode for this Node Pool, specifying how these Nodes should be used (for either System or User resources).

* `node_count` - The current number of Nodes in the Node Pool.

* `node_labels` - A map of Kubernetes Labels applied to each Node in this Node Pool.

* `node_public_ip_prefix_id` - Resource ID for the Public IP Addresses Prefix for the nodes in this Agent Pool.

* `node_taints` - A map of Kubernetes Taints applied to each Node in this Node Pool.

* `orchestrator_version` - The version of Kubernetes configured on each Node in this Node Pool.

* `os_disk_size_gb` - The size of the OS Disk on each Node in this Node Pool.

* `os_disk_type` - The type of the OS Disk on each Node in this Node Pool.

* `os_type` - The operating system used on each Node in this Node Pool.

* `priority` - The priority of the Virtual Machines in the Virtual Machine Scale Set backing this Node Pool.

* `proximity_placement_group_id` - The ID of the Proximity Placement Group where the Virtual Machine Scale Set backing this Node Pool will be placed.

* `spot_max_price` - The maximum price being paid for Virtual Machines in this Scale Set. `-1` means the current on-demand price for a Virtual Machine.

* `tags` - A mapping of tags assigned to the Kubernetes Cluster Node Pool.

* `upgrade_settings` - A `upgrade_settings` block as documented below.

* `vm_size` - The size of the Virtual Machines used in the Virtual Machine Scale Set backing this Node Pool.

* `vnet_subnet_id` - The ID of the Subnet in which this Node Pool exists.

---

A `upgrade_settings` block exports the following:

* `max_surge` - The maximum number or percentage of nodes which will be added to the Node Pool size during an upgrade.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Node Pool.
