---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_node_pool"
description: |-
  Manages a Node Pool within a Kubernetes Cluster
---

# azurerm_kubernetes_cluster_node_pool

Manages a Node Pool within a Kubernetes Cluster

~> **NOTE:** Multiple Node Pools are only supported when the Kubernetes Cluster is using Virtual Machine Scale Sets.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "exampleaks1"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }
}

resource "azurerm_kubernetes_cluster_node_pool" "example" {
  name                  = "internal"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  vm_size               = "Standard_DS2_v2"
  node_count            = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Node Pool which should be created within the Kubernetes Cluster. Changing this forces a new resource to be created.

-> **NOTE:** A Windows Node Pool cannot have a `name` longer than 6 characters.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster where this Node Pool should exist. Changing this forces a new resource to be created.

* `vm_size` - (Required) The SKU which should be used for the Virtual Machines used in this Node Pool. Changing this forces a new resource to be created.

---

* `availability_zones` - (Optional) A list of Availability Zones where the Nodes in this Node Pool should be created in.

* `enable_auto_scaling` - (Optional) Whether to enable [auto-scaler](https://docs.microsoft.com/en-us/azure/aks/cluster-autoscaler). Defaults to `false`.

-> **NOTE:** Additional fields must be configured depending on the value of this field - see below.

* `enable_node_public_ip` - (Optional) Should each node have a Public IP Address? Defaults to `false`.

* `max_pods` - (Optional) The maximum number of pods that can run on each agent. Changing this forces a new resource to be created.

* `node_taints` - (Optional) A list of Kubernetes taints which should be applied to nodes in the agent pool (e.g `key=value:NoSchedule`).

* `os_disk_size_gb` - (Optional) The Agent Operating System disk size in GB. Changing this forces a new resource to be created.

* `os_type` - (Optional) The Operating System which should be used for this Node Pool. Changing this forces a new resource to be created. Possible values are `Linux` and `Windows`. Defaults to `Linux`.

* `vnet_subnet_id` - (Optional) The ID of the Subnet where this Node Pool should exist.

-> **NOTE:** At this time the `vnet_subnet_id` must be the same for all node pools in the cluster

~> **NOTE:** A route table must be configured on this Subnet.

---

When `enable_auto_scaling` is set to `true` the following fields are applicable:

* `max_count` - (Required) The maximum number of nodes which should exist within this Node Pool. Valid values are between `1` and `100` and must be greater than or equal to `min_count`.

* `min_count` - (Required) The minimum number of nodes which should exist within this Node Pool. Valid values are between `1` and `100` and must be less than or equal to `max_count`.

* `node_count` - (Optional) The initial number of nodes which should exist within this Node Pool. Valid values are between `1` and `100` and must be a value in the range `min_count` - `max_count`.

-> **NOTE:** If you're specifying an initial number of nodes you may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) to ignore changes to this field.

When `enable_auto_scaling` is set to `false` the following fields are applicable:

* `node_count` - (Required) The number of nodes which should exist within this Node Pool. Valid values are between `1` and `100`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes Cluster Node Pool.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Kubernetes Cluster Node Pool.
* `update` - (Defaults to 60 minutes) Used when updating the Kubernetes Cluster Node Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Node Pool.
* `delete` - (Defaults to 60 minutes) Used when deleting the Kubernetes Cluster Node Pool.

## Import

Kubernetes Cluster Node Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_node_pool.pool1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1
```
