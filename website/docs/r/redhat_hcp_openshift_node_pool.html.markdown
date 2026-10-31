---
subcategory: "Red Hat OpenShift"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_redhat_hcp_openshift_node_pool"
description: |-
  Manages an Azure Red Hat OpenShift with Hosted Control Planes Node Pool.
---

# azurerm_redhat_hcp_openshift_node_pool

Manages an Azure Red Hat OpenShift with Hosted Control Planes Node Pool.

The node pool inherits its location from the parent HCP cluster.

~> **Caution:** The ARO HCP service currently does not support deleting the last node pool in a cluster. When destroying a configuration containing an HCP cluster, its dependencies, and one or more node pools, Terraform cannot delete the last node pool, preventing it from completing the destruction of the cluster and its dependencies. Before running `terraform destroy`, remove all of the cluster's node pools from Terraform state using `terraform state rm`, or apply `removed` blocks for the node pools with `lifecycle { destroy = false }`. Both workarounds leave the node pools in Azure so they can be deleted as part of deleting the parent cluster.

## Example Usage

The following example adds a node pool to the ARO HCP cluster `azurerm_redhat_hcp_openshift_cluster.example`:

```hcl
resource "azurerm_redhat_hcp_openshift_node_pool" "example" {
  name                     = "workers"
  hcp_openshift_cluster_id = azurerm_redhat_hcp_openshift_cluster.example.id
  replicas                 = 2

  platform {
    vm_size = "Standard_D4s_v3"
  }

  version {
    id = azurerm_redhat_hcp_openshift_cluster.example.version[0].id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the node pool. Must be between 1 and 15 characters in length. Changing this forces a new resource to be created.
* `hcp_openshift_cluster_id` - (Required) The ID of the parent ARO HCP cluster. Changing this forces a new resource to be created.
* `platform` - (Required) A `platform` block as defined below.
* `auto_repair` - (Optional) Whether automatic node repair is enabled. Defaults to `true`. Changing this forces a new resource to be created.
* `auto_scaling` - (Optional) An `auto_scaling` block as defined below. Conflicts with `replicas`.
* `labels` - (Optional) A mapping of Kubernetes labels assigned to nodes.
* `node_drain_timeout_minutes` - (Optional) The node drain timeout in minutes, between `0` and `10080`. Defaults to `0`.
* `replicas` - (Optional) The number of nodes. Conflicts with `auto_scaling`. When omitted, the value is determined by the service.
* `taints` - (Optional) One or more `taints` blocks as defined below.
* `version` - (Optional) A `version` block as defined below. When omitted, the value is determined by the service.
* `tags` - (Optional) A mapping of tags assigned to the resource.

---

An `auto_scaling` block supports the following:

* `max` - (Required) The maximum number of nodes, at least `0` and not less than `min`. Must not exceed `200` when `platform.availability_zone` is not specified. There is no upper limit when `platform.availability_zone` is specified.
* `min` - (Required) The minimum number of nodes, at least `0`. Must not exceed `200` when `platform.availability_zone` is not specified. There is no upper limit when `platform.availability_zone` is specified.

---

An `os_disk` block supports the following:

* `disk_storage_account_type` - (Optional) The storage account type. Possible values are `Premium_LRS`, `Standard_LRS`, and `StandardSSD_LRS`.
* `disk_type` - (Optional) The disk type. Possible values are `Managed` and `Ephemeral`.
* `encryption_set_id` - (Optional) The disk encryption set ID.
* `size_gb` - (Optional) The OS disk size in GiB, at least `1`.

---

A `platform` block supports the following:

* `vm_size` - (Required) The size of the virtual machines. Changing this forces a new resource to be created.
* `availability_zone` - (Optional) The availability zone for the nodes. Changing this forces a new resource to be created.
* `enable_encryption_at_host` - (Optional) Whether host encryption is enabled. Changing this forces a new resource to be created.
* `os_disk` - (Optional) An `os_disk` block as defined above. Changing this forces a new resource to be created.
* `subnet_id` - (Optional) The subnet ID for the nodes. Changing this forces a new resource to be created.

---

A `version` block supports the following:

* `id` - (Required) The OpenShift version in `X.Y.Z` format with an optional pre-release suffix.
* `channel_group` - (Optional) The channel group. Possible values are `nightly`, `candidate`, `fast`, and `stable`. Defaults to `stable`. Changing this forces a new resource to be created.

---

A `taints` block supports the following:

* `effect` - (Required) The effect. Possible values are `NoExecute`, `NoSchedule`, and `PreferNoSchedule`.
* `key` - (Required) The taint key.
* `value` - (Optional) The taint value.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the node pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the node pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the node pool.
* `update` - (Defaults to 60 minutes) Used when updating the node pool.
* `delete` - (Defaults to 60 minutes) Used when deleting the node pool.

## Import

Node Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_redhat_hcp_openshift_node_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.RedHatOpenShift/hcpOpenShiftClusters/example-cluster/nodePools/workers
```
