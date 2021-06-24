---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_node_pool"
description: |-
  Manages a Node Pool within a Kubernetes Cluster
---

# azurerm_kubernetes_cluster_node_pool

Manages a Node Pool within a Kubernetes Cluster

-> **Note:** Due to the fast-moving nature of AKS, we recommend using the latest version of the Azure Provider when using AKS - you can find [the latest version of the Azure Provider here](https://registry.terraform.io/providers/hashicorp/azurerm/latest).

~> **NOTE:** Multiple Node Pools are only supported when the Kubernetes Cluster is using Virtual Machine Scale Sets.

## Example Usage

This example provisions a basic Kubernetes Node Pool. Other examples of the `azurerm_kubernetes_cluster_node_pool` resource can be found in [the `./examples/kubernetes` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/kubernetes)


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

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Node Pool which should be created within the Kubernetes Cluster. Changing this forces a new resource to be created.

~> **NOTE:** A Windows Node Pool cannot have a `name` longer than 6 characters.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster where this Node Pool should exist. Changing this forces a new resource to be created.

~> **NOTE:** The type of Default Node Pool for the Kubernetes Cluster must be `VirtualMachineScaleSets` to attach multiple node pools.

* `vm_size` - (Required) The SKU which should be used for the Virtual Machines used in this Node Pool. Changing this forces a new resource to be created.

---

* `availability_zones` - (Optional) A list of Availability Zones where the Nodes in this Node Pool should be created in. Changing this forces a new resource to be created.

* `enable_auto_scaling` - (Optional) Whether to enable [auto-scaler](https://docs.microsoft.com/en-us/azure/aks/cluster-autoscaler). Defaults to `false`.

* `enable_host_encryption` - (Optional) Should the nodes in this Node Pool have host encryption enabled? Defaults to `false`.

~> **NOTE:** Additional fields must be configured depending on the value of this field - see below.

* `enable_node_public_ip` - (Optional) Should each node have a Public IP Address? Defaults to `false`.  Changing this forces a new resource to be created.

* `eviction_policy` - (Optional) The Eviction Policy which should be used for Virtual Machines within the Virtual Machine Scale Set powering this Node Pool. Possible values are `Deallocate` and `Delete`. Changing this forces a new resource to be created.

-> **Note:** An Eviction Policy can only be configured when `priority` is set to `Spot`.

* `kubelet_config` - (Optional) A `kubelet_config` block as defined below.

* `linux_os_config` - (Optional) A `linux_os_config` block as defined below.

* `max_pods` - (Optional) The maximum number of pods that can run on each agent. Changing this forces a new resource to be created.

* `mode` - (Optional) Should this Node Pool be used for System or User resources? Possible values are `System` and `User`. Defaults to `User`.

* `node_labels` - (Optional) A map of Kubernetes labels which should be applied to nodes in this Node Pool. Changing this forces a new resource to be created.

* `node_public_ip_prefix_id` - (Optional) Resource ID for the Public IP Addresses Prefix for the nodes in this Node Pool. `enable_node_public_ip` should be `true`. Changing this forces a new resource to be created.

* `node_taints` - (Optional) A list of Kubernetes taints which should be applied to nodes in the agent pool (e.g `key=value:NoSchedule`). Changing this forces a new resource to be created.

* `orchestrator_version` - (Optional) Version of Kubernetes used for the Agents. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade)

-> **Note:** This version must be supported by the Kubernetes Cluster - as such the version of Kubernetes used on the Cluster/Control Plane may need to be upgraded first.

* `os_disk_size_gb` - (Optional) The Agent Operating System disk size in GB. Changing this forces a new resource to be created.

* `os_disk_type` - (Optional) The type of disk which should be used for the Operating System. Possible values are `Ephemeral` and `Managed`. Defaults to `Managed`. Changing this forces a new resource to be created.

* `os_type` - (Optional) The Operating System which should be used for this Node Pool. Changing this forces a new resource to be created. Possible values are `Linux` and `Windows`. Defaults to `Linux`.

* `priority` - (Optional) The Priority for Virtual Machines within the Virtual Machine Scale Set that powers this Node Pool. Possible values are `Regular` and `Spot`. Defaults to `Regular`. Changing this forces a new resource to be created.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group where the Virtual Machine Scale Set that powers this Node Pool will be placed. Changing this forces a new resource to be created.

-> **Note:** When setting `priority` to Spot - you must configure an `eviction_policy`, `spot_max_price` and add the applicable `node_labels` and `node_taints` [as per the Azure Documentation](https://docs.microsoft.com/en-us/azure/aks/spot-node-pool).

~> **Note:** Spot Node Pools are in Preview and must be opted-into - [more information on how to opt into this Preview can be found in the AKS Documentation](https://docs.microsoft.com/en-us/azure/aks/spot-node-pool).

* `spot_max_price` - (Optional) The maximum price you're willing to pay in USD per Virtual Machine. Valid values are `-1` (the current on-demand price for a Virtual Machine) or a positive value with up to five decimal places. Changing this forces a new resource to be created.

~> **Note:** This field can only be configured when `priority` is set to `Spot`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> At this time there's a bug in the AKS API where Tags for a Node Pool are not stored in the correct case - you [may wish to use Terraform's `ignore_changes` functionality to ignore changes to the casing](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) until this is fixed in the AKS API.

* `upgrade_settings` - (Optional) A `upgrade_settings` block as documented below.

* `vnet_subnet_id` - (Optional) The ID of the Subnet where this Node Pool should exist.

-> **NOTE:** At this time the `vnet_subnet_id` must be the same for all node pools in the cluster

~> **NOTE:** A route table must be configured on this Subnet.

---

If `enable_auto_scaling` is set to `true`, then the following fields can also be configured:

* `max_count` - (Required) The maximum number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000` and must be greater than or equal to `min_count`.

* `min_count` - (Required) The minimum number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000` and must be less than or equal to `max_count`.

* `node_count` - (Optional) The initial number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000` and must be a value in the range `min_count` - `max_count`.

-> **NOTE:** If you're specifying an initial number of nodes you may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/docs/configuration/resources.html#ignore_changes) to ignore changes to this field.

If `enable_auto_scaling` is set to `false`, then the following fields can also be configured:

* `node_count` - (Required) The number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000`.

---

A `kubelet_config` block supports the following:

* `allowed_unsafe_sysctls` - (Optional) Specifies the allow list of unsafe sysctls command or patterns (ending in `*`). Changing this forces a new resource to be created.

* `container_log_max_line` - (Optional) Specifies the maximum number of container log files that can be present for a container. must be at least 2. Changing this forces a new resource to be created.

* `container_log_max_size_mb` - (Optional) Specifies the maximum size (e.g. 10MB) of container log file before it is rotated. Changing this forces a new resource to be created.

* `cpu_cfs_quota_enabled` - (Optional) Is CPU CFS quota enforcement for containers enabled? Changing this forces a new resource to be created.

* `cpu_cfs_quota_period` - (Optional) Specifies the CPU CFS quota period value. Changing this forces a new resource to be created.

* `cpu_manager_policy` - (Optional) Specifies the CPU Manager policy to use. Possible values are `none` and `static`, Changing this forces a new resource to be created.

* `image_gc_high_threshold` - (Optional) Specifies the percent of disk usage above which image garbage collection is always run. Must be between `0` and `100`. Changing this forces a new resource to be created.

* `image_gc_low_threshold` - (Optional) Specifies the percent of disk usage lower than which image garbage collection is never run. Must be between `0` and `100`. Changing this forces a new resource to be created.

* `pod_max_pids` - (Optional) Specifies the maximum number of processes per pod. Changing this forces a new resource to be created.

* `topology_manager_policy` - (Optional) Specifies the Topology Manager policy to use. Possible values are `none`, `best-effort`, `restricted` or `single-numa-node`. Changing this forces a new resource to be created.

---

A `linux_os_config` block supports the following:

* `swap_file_size_mb` - (Optional) Specifies the size of swap file on each node in MB. Changing this forces a new resource to be created.

* `sysctl_config` - (Optional) A `sysctl_config` block as defined below. Changing this forces a new resource to be created.

* `transparent_huge_page_defrag` - (Optional) specifies the defrag configuration for Transparent Huge Page. Possible values are `always`, `defer`, `defer+madvise`, `madvise` and `never`. Changing this forces a new resource to be created.

* `transparent_huge_page_enabled` - (Optional) Specifies the Transparent Huge Page enabled configuration. Possible values are `always`, `madvise` and `never`. Changing this forces a new resource to be created.

---

A `sysctl_config` block supports the following:

~> For more information, please refer to [Linux Kernel Doc](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/index.html).

* `fs_aio_max_nr` - (Optional) The sysctl setting fs.aio-max-nr. Must be between `65536` and `6553500`. Changing this forces a new resource to be created.

* `fs_file_max` - (Optional) The sysctl setting fs.file-max. Must be between `8192` and `12000500`. Changing this forces a new resource to be created.

* `fs_inotify_max_user_watches` - (Optional) The sysctl setting fs.inotify.max_user_watches. Must be between `781250` and `2097152`. Changing this forces a new resource to be created.

* `fs_nr_open` - (Optional) The sysctl setting fs.nr_open. Must be between `8192` and `20000500`. Changing this forces a new resource to be created.

* `kernel_threads_max` - (Optional) The sysctl setting kernel.threads-max. Must be between `20` and `513785`. Changing this forces a new resource to be created.

* `net_core_netdev_max_backlog` - (Optional) The sysctl setting net.core.netdev_max_backlog. Must be between `1000` and `3240000`. Changing this forces a new resource to be created.

* `net_core_optmem_max` - (Optional) The sysctl setting net.core.optmem_max. Must be between `20480` and `4194304`. Changing this forces a new resource to be created.

* `net_core_rmem_default` - (Optional) The sysctl setting net.core.rmem_default. Must be between `212992` and `134217728`. Changing this forces a new resource to be created.

* `net_core_rmem_max` - (Optional) The sysctl setting net.core.rmem_max. Must be between `212992` and `134217728`. Changing this forces a new resource to be created.

* `net_core_somaxconn` - (Optional) The sysctl setting net.core.somaxconn. Must be between `4096` and `3240000`. Changing this forces a new resource to be created.

* `net_core_wmem_default` - (Optional) The sysctl setting net.core.wmem_default. Must be between `212992` and `134217728`. Changing this forces a new resource to be created.

* `net_core_wmem_max` - (Optional) The sysctl setting net.core.wmem_max. Must be between `212992` and `134217728`. Changing this forces a new resource to be created.

* `net_ipv4_ip_local_port_range_max` - (Optional) The sysctl setting net.ipv4.ip_local_port_range max value. Must be between `1024` and `60999`. Changing this forces a new resource to be created.

* `net_ipv4_ip_local_port_range_min` - (Optional) The sysctl setting net.ipv4.ip_local_port_range min value. Must be between `1024` and `60999`. Changing this forces a new resource to be created.

* `net_ipv4_neigh_default_gc_thresh1` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh1. Must be between `128` and `80000`. Changing this forces a new resource to be created.

* `net_ipv4_neigh_default_gc_thresh2` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh2. Must be between `512` and `90000`. Changing this forces a new resource to be created.

* `net_ipv4_neigh_default_gc_thresh3` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh3. Must be between `1024` and `100000`. Changing this forces a new resource to be created.

* `net_ipv4_tcp_fin_timeout` - (Optional) The sysctl setting net.ipv4.tcp_fin_timeout. Must be between `5` and `120`. Changing this forces a new resource to be created.

* `net_ipv4_tcp_keepalive_intvl` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_intvl. Must be between `10` and `75`. Changing this forces a new resource to be created.

* `net_ipv4_tcp_keepalive_probes` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_probes. Must be between `1` and `15`. Changing this forces a new resource to be created.

* `net_ipv4_tcp_keepalive_time` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_time. Must be between `30` and `432000`. Changing this forces a new resource to be created.

* `net_ipv4_tcp_max_syn_backlog` - (Optional) The sysctl setting net.ipv4.tcp_max_syn_backlog. Must be between `128` and `3240000`. Changing this forces a new resource to be created.

* `net_ipv4_tcp_max_tw_buckets` - (Optional) The sysctl setting net.ipv4.tcp_max_tw_buckets. Must be between `8000` and `1440000`. Changing this forces a new resource to be created.

* `net_ipv4_tcp_tw_reuse` - (Optional) Is sysctl setting net.ipv4.tcp_tw_reuse enabled? Changing this forces a new resource to be created.

* `net_netfilter_nf_conntrack_buckets` - (Optional) The sysctl setting net.netfilter.nf_conntrack_buckets. Must be between `65536` and `147456`. Changing this forces a new resource to be created.

* `net_netfilter_nf_conntrack_max` - (Optional) The sysctl setting net.netfilter.nf_conntrack_max. Must be between `131072` and `589824`. Changing this forces a new resource to be created.

* `vm_max_map_count` - (Optional) The sysctl setting vm.max_map_count. Must be between `65530` and `262144`. Changing this forces a new resource to be created.

* `vm_swappiness` - (Optional) The sysctl setting vm.swappiness. Must be between `0` and `100`. Changing this forces a new resource to be created.

* `vm_vfs_cache_pressure` - (Optional) The sysctl setting vm.vfs_cache_pressure. Must be between `0` and `100`. Changing this forces a new resource to be created.

---

A `upgrade_settings` block supports the following:

* `max_surge` - (Required) The maximum number or percentage of nodes which will be added to the Node Pool size during an upgrade.

-> **Note:** If a percentage is provided, the number of surge nodes is calculated from the current node count on the cluster. Node surge can allow a cluster to have more nodes than `max_count` during an upgrade. Ensure that your cluster has enough [IP space](https://docs.microsoft.com/en-us/azure/aks/upgrade-cluster#customize-node-surge-upgrade) during an upgrade.

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
