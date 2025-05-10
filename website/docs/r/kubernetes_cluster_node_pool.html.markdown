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

~> **Note:** Multiple Node Pools are only supported when the Kubernetes Cluster is using Virtual Machine Scale Sets.

-> **Note:** Changing certain properties is done by cycling the node pool. When cycling it, it doesn’t perform cordon and drain, and it will disrupt rescheduling pods currently running on the previous node pool. `temporary_name_for_rotation` must be specified when changing any of the following properties: `fips_enabled`, `host_encryption_enabled`, `kubelet_config`, `kubelet_disk_type`, `linux_os_config`, `max_pods`, `node_public_ip_enabled`, `os_disk_size_gb`, `os_disk_type`, `pod_subnet_id`, `snapshot_id`, `ultra_ssd_enabled`, `vm_size`, `vnet_subnet_id`, `zones`.

## Example Usage

This example provisions a basic Kubernetes Node Pool. Other examples of the `azurerm_kubernetes_cluster_node_pool` resource can be found in [the `./examples/kubernetes` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/kubernetes)

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

~> **Note:** A Windows Node Pool cannot have a `name` longer than 6 characters.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster where this Node Pool should exist. Changing this forces a new resource to be created.

~> **Note:** The type of Default Node Pool for the Kubernetes Cluster must be `VirtualMachineScaleSets` to attach multiple node pools.

* `vm_size` - (Required) The SKU which should be used for the Virtual Machines used in this Node Pool. Changing this property requires specifying `temporary_name_for_rotation`.

---

* `capacity_reservation_group_id` - (Optional) Specifies the ID of the Capacity Reservation Group where this Node Pool should exist. Changing this forces a new resource to be created.

* `auto_scaling_enabled` - (Optional) Whether to enable [auto-scaler](https://docs.microsoft.com/azure/aks/cluster-autoscaler).

* `host_encryption_enabled` - (Optional) Should the nodes in this Node Pool have host encryption enabled? Changing this property requires specifying `temporary_name_for_rotation`.

~> **Note:** Additional fields must be configured depending on the value of this field - see below.

* `node_public_ip_enabled` - (Optional) Should each node have a Public IP Address? Changing this property requires specifying `temporary_name_for_rotation`.

* `eviction_policy` - (Optional) The Eviction Policy which should be used for Virtual Machines within the Virtual Machine Scale Set powering this Node Pool. Possible values are `Deallocate` and `Delete`. Changing this forces a new resource to be created.

~> **Note:** An Eviction Policy can only be configured when `priority` is set to `Spot` and will default to `Delete` unless otherwise specified.

* `host_group_id` - (Optional) The fully qualified resource ID of the Dedicated Host Group to provision virtual machines from. Changing this forces a new resource to be created.

* `kubelet_config` - (Optional) A `kubelet_config` block as defined below. Changing this requires specifying `temporary_name_for_rotation`.

* `linux_os_config` - (Optional) A `linux_os_config` block as defined below. Changing this requires specifying `temporary_name_for_rotation`.

* `fips_enabled` - (Optional) Should the nodes in this Node Pool have Federal Information Processing Standard enabled? Changing this property requires specifying `temporary_name_for_rotation`.

~> **Note:** FIPS support is in Public Preview - more information and details on how to opt into the Preview can be found in [this article](https://docs.microsoft.com/azure/aks/use-multiple-node-pools#add-a-fips-enabled-node-pool-preview).

* `gpu_instance` - (Optional) Specifies the GPU MIG instance profile for supported GPU VM SKU. The allowed values are `MIG1g`, `MIG2g`, `MIG3g`, `MIG4g` and `MIG7g`. Changing this forces a new resource to be created.

* `kubelet_disk_type` - (Optional) The type of disk used by kubelet. Possible values are `OS` and `Temporary`. Changing this property requires specifying `temporary_name_for_rotation`.

* `max_pods` - (Optional) The maximum number of pods that can run on each agent. Changing this property requires specifying `temporary_name_for_rotation`.

* `mode` - (Optional) Should this Node Pool be used for System or User resources? Possible values are `System` and `User`. Defaults to `User`.

* `node_network_profile` - (Optional) A `node_network_profile` block as documented below.

* `node_labels` - (Optional) A map of Kubernetes labels which should be applied to nodes in this Node Pool.

* `node_public_ip_prefix_id` - (Optional) Resource ID for the Public IP Addresses Prefix for the nodes in this Node Pool. `node_public_ip_enabled` should be `true`. Changing this forces a new resource to be created.

* `node_taints` - (Optional) A list of Kubernetes taints which should be applied to nodes in the agent pool (e.g `key=value:NoSchedule`).

* `orchestrator_version` - (Optional) Version of Kubernetes used for the Agents. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade). AKS does not require an exact patch version to be specified, minor version aliases such as `1.22` are also supported. - The minor version's latest GA patch is automatically chosen in that case. More details can be found in [the documentation](https://docs.microsoft.com/en-us/azure/aks/supported-kubernetes-versions?tabs=azure-cli#alias-minor-version).

-> **Note:** This version must be supported by the Kubernetes Cluster - as such the version of Kubernetes used on the Cluster/Control Plane may need to be upgraded first.

* `os_disk_size_gb` - (Optional) The Agent Operating System disk size in GB. Changing this property requires specifying `temporary_name_for_rotation`.

* `os_disk_type` - (Optional) The type of disk which should be used for the Operating System. Possible values are `Ephemeral` and `Managed`. Defaults to `Managed`. Changing this property requires specifying `temporary_name_for_rotation`.

* `pod_subnet_id` - (Optional) The ID of the Subnet where the pods in the Node Pool should exist. Changing this property requires specifying `temporary_name_for_rotation`.

* `os_sku` - (Optional) Specifies the OS SKU used by the agent pool. Possible values are `AzureLinux`, `Ubuntu`, `Windows2019` and `Windows2022`. If not specified, the default is `Ubuntu` if OSType=Linux or `Windows2019` if OSType=Windows. And the default Windows OSSKU will be changed to `Windows2022` after Windows2019 is deprecated. Changing this from `AzureLinux` or `Ubuntu` to `AzureLinux` or `Ubuntu` will not replace the resource, otherwise it forces a new resource to be created.

* `os_type` - (Optional) The Operating System which should be used for this Node Pool. Changing this forces a new resource to be created. Possible values are `Linux` and `Windows`. Defaults to `Linux`.

* `priority` - (Optional) The Priority for Virtual Machines within the Virtual Machine Scale Set that powers this Node Pool. Possible values are `Regular` and `Spot`. Defaults to `Regular`. Changing this forces a new resource to be created.

* `proximity_placement_group_id` - (Optional) The ID of the Proximity Placement Group where the Virtual Machine Scale Set that powers this Node Pool will be placed. Changing this forces a new resource to be created.

-> **Note:** When setting `priority` to Spot - you must configure an `eviction_policy`, `spot_max_price` and add the applicable `node_labels` and `node_taints` [as per the Azure Documentation](https://docs.microsoft.com/azure/aks/spot-node-pool).

* `spot_max_price` - (Optional) The maximum price you're willing to pay in USD per Virtual Machine. Valid values are `-1` (the current on-demand price for a Virtual Machine) or a positive value with up to five decimal places. Changing this forces a new resource to be created.

~> **Note:** This field can only be configured when `priority` is set to `Spot`.

* `snapshot_id` - (Optional) The ID of the Snapshot which should be used to create this Node Pool. Changing this property requires specifying `temporary_name_for_rotation`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> **Note:** At this time there's a bug in the AKS API where Tags for a Node Pool are not stored in the correct case - you [may wish to use Terraform's `ignore_changes` functionality to ignore changes to the casing](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changess) until this is fixed in the AKS API.

* `scale_down_mode` - (Optional) Specifies how the node pool should deal with scaled-down nodes. Allowed values are `Delete` and `Deallocate`. Defaults to `Delete`.

* `temporary_name_for_rotation` - (Optional) Specifies the name of the temporary node pool used to cycle the node pool when one of the relevant properties are updated.

* `ultra_ssd_enabled` - (Optional) Used to specify whether the UltraSSD is enabled in the Node Pool. Defaults to `false`. See [the documentation](https://docs.microsoft.com/azure/aks/use-ultra-disks) for more information. Changing this property requires specifying `temporary_name_for_rotation`.

* `upgrade_settings` - (Optional) A `upgrade_settings` block as documented below.

* `vnet_subnet_id` - (Optional) The ID of the Subnet where this Node Pool should exist. Changing this property requires specifying `temporary_name_for_rotation`.

~> **Note:** A route table must be configured on this Subnet.

* `windows_profile` - (Optional) A `windows_profile` block as documented below. Changing this forces a new resource to be created.

* `workload_runtime` - (Optional) Used to specify the workload runtime. Allowed values are `OCIContainer` and `WasmWasi`.

~> **Note:** WebAssembly System Interface node pools are in Public Preview - more information and details on how to opt into the preview can be found in [this article](https://docs.microsoft.com/azure/aks/use-wasi-node-pools)

* `zones` - (Optional) Specifies a list of Availability Zones in which this Kubernetes Cluster Node Pool should be located. Changing this property requires specifying `temporary_name_for_rotation`.

---

If `auto_scaling_enabled` is set to `true`, then the following fields can also be configured:

* `max_count` - (Optional) The maximum number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000` and must be greater than or equal to `min_count`.

* `min_count` - (Optional) The minimum number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000` and must be less than or equal to `max_count`.

* `node_count` - (Optional) The initial number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000` (inclusive) for user pools and between `1` and `1000` (inclusive) for system pools and must be a value in the range `min_count` - `max_count`.

-> **Note:** If you're specifying an initial number of nodes you may wish to use [Terraform's `ignore_changes` functionality](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changess) to ignore changes to this field.

If `auto_scaling_enabled` is set to `false`, then the following fields can also be configured:

* `node_count` - (Optional) The number of nodes which should exist within this Node Pool. Valid values are between `0` and `1000` (inclusive) for user pools and between `1` and `1000` (inclusive) for system pools.

---

A `kubelet_config` block supports the following:

* `allowed_unsafe_sysctls` - (Optional) Specifies the allow list of unsafe sysctls command or patterns (ending in `*`). 

* `container_log_max_line` - (Optional) Specifies the maximum number of container log files that can be present for a container. must be at least 2. 

* `container_log_max_size_mb` - (Optional) Specifies the maximum size (e.g. 10MB) of container log file before it is rotated. 

* `cpu_cfs_quota_enabled` - (Optional) Is CPU CFS quota enforcement for containers enabled? Defaults to `true`.

* `cpu_cfs_quota_period` - (Optional) Specifies the CPU CFS quota period value. 

* `cpu_manager_policy` - (Optional) Specifies the CPU Manager policy to use. Possible values are `none` and `static`, 

* `image_gc_high_threshold` - (Optional) Specifies the percent of disk usage above which image garbage collection is always run. Must be between `0` and `100`. 

* `image_gc_low_threshold` - (Optional) Specifies the percent of disk usage lower than which image garbage collection is never run. Must be between `0` and `100`. 

* `pod_max_pid` - (Optional) Specifies the maximum number of processes per pod. 

* `topology_manager_policy` - (Optional) Specifies the Topology Manager policy to use. Possible values are `none`, `best-effort`, `restricted` or `single-numa-node`. 

---

A `linux_os_config` block supports the following:

* `swap_file_size_mb` - (Optional) Specifies the size of swap file on each node in MB.

* `sysctl_config` - (Optional) A `sysctl_config` block as defined below.

* `transparent_huge_page_defrag` - (Optional) specifies the defrag configuration for Transparent Huge Page. Possible values are `always`, `defer`, `defer+madvise`, `madvise` and `never`. 

* `transparent_huge_page_enabled` - (Optional) Specifies the Transparent Huge Page enabled configuration. Possible values are `always`, `madvise` and `never`.

---

A `node_network_profile` block supports the following:

* `allowed_host_ports` - (Optional) One or more `allowed_host_ports` blocks as defined below.

* `application_security_group_ids` - (Optional) A list of Application Security Group IDs which should be associated with this Node Pool.

* `node_public_ip_tags` - (Optional) Specifies a mapping of tags to the instance-level public IPs. Changing this forces a new resource to be created.

-> **Note:** To set the application security group, you must allow at least one host port. Without this, the configuration will fail silently. [Learn More](https://learn.microsoft.com/en-us/azure/aks/use-node-public-ips#allow-host-port-connections-and-add-node-pools-to-application-security-groups).

---

An `allowed_host_ports` block supports the following:

* `port_start` - (Optional) Specifies the start of the port range.

* `port_end` - (Optional) Specifies the end of the port range.

* `protocol` - (Optional) Specifies the protocol of the port range. Possible values are `TCP` and `UDP`.

---

A `sysctl_config` block supports the following:

~> **Note:** For more information, please refer to [Linux Kernel Doc](https://www.kernel.org/doc/html/latest/admin-guide/sysctl/index.html).

* `fs_aio_max_nr` - (Optional) The sysctl setting fs.aio-max-nr. Must be between `65536` and `6553500`.

* `fs_file_max` - (Optional) The sysctl setting fs.file-max. Must be between `8192` and `12000500`.

* `fs_inotify_max_user_watches` - (Optional) The sysctl setting fs.inotify.max_user_watches. Must be between `781250` and `2097152`.

* `fs_nr_open` - (Optional) The sysctl setting fs.nr_open. Must be between `8192` and `20000500`.

* `kernel_threads_max` - (Optional) The sysctl setting kernel.threads-max. Must be between `20` and `513785`.

* `net_core_netdev_max_backlog` - (Optional) The sysctl setting net.core.netdev_max_backlog. Must be between `1000` and `3240000`.

* `net_core_optmem_max` - (Optional) The sysctl setting net.core.optmem_max. Must be between `20480` and `4194304`.

* `net_core_rmem_default` - (Optional) The sysctl setting net.core.rmem_default. Must be between `212992` and `134217728`.

* `net_core_rmem_max` - (Optional) The sysctl setting net.core.rmem_max. Must be between `212992` and `134217728`.

* `net_core_somaxconn` - (Optional) The sysctl setting net.core.somaxconn. Must be between `4096` and `3240000`.

* `net_core_wmem_default` - (Optional) The sysctl setting net.core.wmem_default. Must be between `212992` and `134217728`.

* `net_core_wmem_max` - (Optional) The sysctl setting net.core.wmem_max. Must be between `212992` and `134217728`.

* `net_ipv4_ip_local_port_range_max` - (Optional) The sysctl setting net.ipv4.ip_local_port_range max value. Must be between `32768` and `65535`.

* `net_ipv4_ip_local_port_range_min` - (Optional) The sysctl setting net.ipv4.ip_local_port_range min value. Must be between `1024` and `60999`.

* `net_ipv4_neigh_default_gc_thresh1` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh1. Must be between `128` and `80000`.

* `net_ipv4_neigh_default_gc_thresh2` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh2. Must be between `512` and `90000`.

* `net_ipv4_neigh_default_gc_thresh3` - (Optional) The sysctl setting net.ipv4.neigh.default.gc_thresh3. Must be between `1024` and `100000`.

* `net_ipv4_tcp_fin_timeout` - (Optional) The sysctl setting net.ipv4.tcp_fin_timeout. Must be between `5` and `120`.

* `net_ipv4_tcp_keepalive_intvl` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_intvl. Must be between `10` and `90`.

* `net_ipv4_tcp_keepalive_probes` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_probes. Must be between `1` and `15`.

* `net_ipv4_tcp_keepalive_time` - (Optional) The sysctl setting net.ipv4.tcp_keepalive_time. Must be between `30` and `432000`.

* `net_ipv4_tcp_max_syn_backlog` - (Optional) The sysctl setting net.ipv4.tcp_max_syn_backlog. Must be between `128` and `3240000`.

* `net_ipv4_tcp_max_tw_buckets` - (Optional) The sysctl setting net.ipv4.tcp_max_tw_buckets. Must be between `8000` and `1440000`.

* `net_ipv4_tcp_tw_reuse` - (Optional) Is sysctl setting net.ipv4.tcp_tw_reuse enabled?

* `net_netfilter_nf_conntrack_buckets` - (Optional) The sysctl setting net.netfilter.nf_conntrack_buckets. Must be between `65536` and `524288`.

* `net_netfilter_nf_conntrack_max` - (Optional) The sysctl setting net.netfilter.nf_conntrack_max. Must be between `131072` and `2097152`.

* `vm_max_map_count` - (Optional) The sysctl setting vm.max_map_count. Must be between `65530` and `262144`.

* `vm_swappiness` - (Optional) The sysctl setting vm.swappiness. Must be between `0` and `100`.

* `vm_vfs_cache_pressure` - (Optional) The sysctl setting vm.vfs_cache_pressure. Must be between `0` and `100`.

---

A `upgrade_settings` block supports the following:

* `drain_timeout_in_minutes` - (Optional) The amount of time in minutes to wait on eviction of pods and graceful termination per node. This eviction wait time honors waiting on pod disruption budgets. If this time is exceeded, the upgrade fails. Unsetting this after configuring it will force a new resource to be created.

* `node_soak_duration_in_minutes` - (Optional) The amount of time in minutes to wait after draining a node and before reimaging and moving on to next node.

* `max_surge` - (Required) The maximum number or percentage of nodes which will be added to the Node Pool size during an upgrade.

---

A `windows_profile` block supports the following:

* `outbound_nat_enabled` - (Optional) Should the Windows nodes in this Node Pool have outbound NAT enabled? Defaults to `true`. Changing this forces a new resource to be created.

-> **Note:** If a percentage is provided, the number of surge nodes is calculated from the current node count on the cluster. Node surge can allow a cluster to have more nodes than `max_count` during an upgrade. Ensure that your cluster has enough [IP space](https://docs.microsoft.com/azure/aks/upgrade-cluster#customize-node-surge-upgrade) during an upgrade.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Cluster Node Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kubernetes Cluster Node Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Node Pool.
* `update` - (Defaults to 1 hour) Used when updating the Kubernetes Cluster Node Pool.
* `delete` - (Defaults to 1 hour) Used when deleting the Kubernetes Cluster Node Pool.

## Import

Kubernetes Cluster Node Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_node_pool.pool1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/agentPools/pool1
```
