---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_inference_cluster"
description: |-
  Manages a Machine Learning Inference Cluster.
---

# azurerm_machine_learning_inference_cluster

Manages a Machine Learning Inference Cluster.

## Example Usage

```hcl
resource "azurerm_machine_learning_inference_cluster" "example" {
  name                  = "example"
  location              = "West Europe"
  kubernetes_cluster_id = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1"

  identity {
    type = "SystemAssigned"
  }
  machine_learning_workspace_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1"
  node_pool_name                = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `identity` - (Required) A `identity` block as defined below. Changing this forces a new Machine Learning Inference Cluster to be created.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster. Changing this forces a new Machine Learning Inference Cluster to be created.

* `location` - (Required) The Azure Region where the Machine Learning Inference Cluster should exist. Changing this forces a new Machine Learning Inference Cluster to be created.

* `machine_learning_workspace_id` - (Required) The ID of the Machine Learning Workspace. Changing this forces a new Machine Learning Inference Cluster to be created.

* `name` - (Required) The name which should be used for this Machine Learning Inference Cluster. Changing this forces a new Machine Learning Inference Cluster to be created.

* `node_pool_name` - (Required) The name of the Kubernetes Cluster's node pool. Changing this forces a new Machine Learning Inference Cluster to be created.

---

* `cluster_purpose` - (Optional) The purpose of the Inference Cluster. If used for Development or Testing, use "Dev" or "Test" here. If using for production use "Prod" here. Changing this forces a new Machine Learning Inference Cluster to be created.

* `description` - (Optional) The description of the Machine Learning compute.

* `ssl` - (Optional) A `ssl` block as defined below. Changing this forces a new Machine Learning Inference Cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Inference Cluster.

---

A `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this Disk Encryption Set. At this time the only possible value is `SystemAssigned`.

---

A `ssl` block supports the following:

* `cert` - (Optional) The certificate for the ssl configuration.

* `cname` - (Optional) The cname of the ssl configuration.

* `key` - (Optional) The key content for the ssl configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Inference Cluster.

* `sku_name` - The type of SKU.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Inference Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Inference Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Inference Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Inference Cluster.

## Import

Machine Learning Inference Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_inference_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1
```
