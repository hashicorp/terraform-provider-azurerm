---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_inference_cluster"
description: |-
  Manages a Machine Learning Inference Cluster.
---

# azurerm_machine_learning_inference_cluster

Manages a Machine Learning Inference Cluster.

~> **NOTE:** The Machine Learning Inference Cluster resource is used to attach an existing AKS cluster to the Machine Learning Workspace, it doesn't create the AKS cluster itself. Therefore it can only be created and deleted, not updated. Any change to the configuration will recreate the resource.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
  tags = {
    "stage" = "example"
  }
}

resource "azurerm_application_insights" "example" {
  name                = "example-ai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_key_vault" "example" {
  name                = "example-kv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "standard"

  purge_protection_enabled = true
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example-mlw"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  default_node_pool {
    name           = "default"
    node_count     = 3
    vm_size        = "Standard_D3_v2"
    vnet_subnet_id = azurerm_subnet.example.id
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_machine_learning_inference_cluster" "example" {
  name                  = "example"
  location              = azurerm_resource_group.example.location
  cluster_purpose       = "FastProd"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  description           = "This is an example cluster used with Terraform"

  machine_learning_workspace_id = azurerm_machine_learning_workspace.example.id

  tags = {
    "stage" = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Machine Learning Inference Cluster. Changing this forces a new Machine Learning Inference Cluster to be created.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster. Changing this forces a new Machine Learning Inference Cluster to be created.

* `location` - (Required) The Azure Region where the Machine Learning Inference Cluster should exist. Changing this forces a new Machine Learning Inference Cluster to be created.

* `machine_learning_workspace_id` - (Required) The ID of the Machine Learning Workspace. Changing this forces a new Machine Learning Inference Cluster to be created.

---

* `cluster_purpose` - (Optional) The purpose of the Inference Cluster. Options are `DevTest`, `DenseProd` and `FastProd`. If used for Development or Testing, use `DevTest` here. Default purpose is `FastProd`, which is recommended for production workloads. Changing this forces a new Machine Learning Inference Cluster to be created.

~> **NOTE:** When creating or attaching a cluster, if the cluster will be used for production (`cluster_purpose = "FastProd"`), then it must contain at least 12 virtual CPUs. The number of virtual CPUs can be calculated by multiplying the number of nodes in the cluster by the number of cores provided by the VM size selected. For example, if you use a VM size of "Standard_D3_v2", which has 4 virtual cores, then you should select 3 or greater as the number of nodes.

* `description` - (Optional) The description of the Machine Learning compute. Changing this forces a new Machine Learning Inference Cluster to be created.

* `ssl` - (Optional) A `ssl` block as defined below. Changing this forces a new Machine Learning Inference Cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Inference Cluster. Changing this forces a new Machine Learning Inference Cluster to be created.


---

A `ssl` block supports the following:

* `cert` - (Optional) The certificate for the ssl configuration.Conflicts with `ssl.0.leaf_domain_label`,`ssl.0.overwrite_existing_domain`. Changing this forces a new Machine Learning Inference Cluster to be created.

* `cname` - (Optional) The cname of the ssl configuration.Conflicts with `ssl.0.leaf_domain_label`,`ssl.0.overwrite_existing_domain`. Changing this forces a new Machine Learning Inference Cluster to be created.

* `key` - (Optional) The key content for the ssl configuration.Conflicts with `ssl.0.leaf_domain_label`,`ssl.0.overwrite_existing_domain`. Changing this forces a new Machine Learning Inference Cluster to be created.

* `leaf_domain_label` - (Optional) The leaf domain label for the ssl configuration. Conflicts with `ssl.0.cert`,`ssl.0.key`,`ssl.0.cname`. Changing this forces a new Machine Learning Inference Cluster to be created.

* `overwrite_existing_domain` - (Optional) Whether or not to overwrite existing leaf domain. Conflicts with `ssl.0.cert`,`ssl.0.key`,`ssl.0.cname` Changing this forces a new Machine Learning Inference Cluster to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Inference Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Inference Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Inference Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Inference Cluster.

## Import

Machine Learning Inference Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_inference_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1
```
