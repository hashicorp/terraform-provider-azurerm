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
  name = "cluster-name"
  resource_group_name = "cluster-rg"
  location = "West Europe"
  workspace_name = "aml-ws"

  identity {
    type = "SystemAssigned"    
  }
  kubernetes_cluster_name = "k8s-cluster-name"
  node_pool_name = "default"
  kubernetes_cluster_rg = "k8s-cluster-rg"
}
```

## Arguments Reference

The following arguments are supported:

* `identity` - (Required) A `identity` block as defined below. Changing this forces a new Machine Learning Inference Cluster to be created.

* `kubernetes_cluster_name` - (Required) The name of the Kubernetes Cluster resource to which to attach the inference cluster to. Changing this forces a new Machine Learning Inference Cluster to be created.

* `kubernetes_cluster_rg` - (Required) The name of the resource group in which the Kubernetes Cluster resides. Changing this forces a new Machine Learning Inference Cluster to be created.

* `location` - (Required) The Azure Region where the Machine Learning Inference Cluster should exist. Changing this forces a new Machine Learning Inference Cluster to be created.

* `name` - (Required) The name which should be used for this Machine Learning Inference Cluster. Changing this forces a new Machine Learning Inference Cluster to be created.

* `node_pool_name` - (Required) The name of the Kubernetes Cluster's node pool. Changing this forces a new Machine Learning Inference Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning Inference Cluster should exist. Changing this forces a new Machine Learning Inference Cluster to be created.

* `workspace_name` - (Required) The name of the Azure Machine Learning Workspace where the Machine Learning Inference Cluster should exist. Changing this forces a new Machine Learning Inference Cluster to be created.

---

* `cluster_purpose` - (Optional) The purpose of the Inference Cluster. If used for Development or Testing, use "Dev" or "Test" here. If using for production use "Prod" here. Changing this forces a new Machine Learning Inference Cluster to be created.

* `description` - (Optional) The description of the Machine Learning compute.

* `ssl_certificate_custom` - (Optional) One or more `ssl_certificate_custom` blocks as defined below. Changing this forces a new Machine Learning Inference Cluster to be created.

* `ssl_enabled` - (Optional) Should the SSL Configuration be enabled?

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Inference Cluster.

---

A `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this Disk Encryption Set. At this time the only possible value is `SystemAssigned`.

---

A `ssl_certificate_custom` block (maximally *one*) supports the following:

* `cert` - (Optional) The content of the custom SSL certificate.

* `cname` - (Optional) The Cname of the custom SSL certificate.

* `key` - (Optional) The content of the key file.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Inference Cluster.

* `resource_id` - The ID of the Machine Learning Inference Cluster.

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