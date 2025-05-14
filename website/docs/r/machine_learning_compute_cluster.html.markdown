---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_compute_cluster"
description: |-
  Manages a Machine Learning Compute Cluster.
---

# azurerm_machine_learning_compute_cluster

Manages a Machine Learning Compute Cluster.
**NOTE:** At this point in time the resource cannot be updated (not supported by the backend Azure Go SDK). Therefore it can only be created and deleted, not updated. At the moment, there is also no possibility to specify ssh User Account Credentials to ssh into the compute cluster.

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
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_machine_learning_compute_cluster" "test" {
  name                          = "example"
  location                      = azurerm_resource_group.example.location
  vm_priority                   = "LowPriority"
  vm_size                       = "Standard_DS2_v2"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.example.id
  subnet_resource_id            = azurerm_subnet.example.id

  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 1
    scale_down_nodes_after_idle_duration = "PT30S" # 30 seconds
  }

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Machine Learning Compute Cluster. Changing this forces a new Machine Learning Compute Cluster to be created.

* `machine_learning_workspace_id` - (Required) The ID of the Machine Learning Workspace. Changing this forces a new Machine Learning Compute Cluster to be created.

* `location` - (Required) The Azure Region where the Machine Learning Compute Cluster should exist. Changing this forces a new Machine Learning Compute Cluster to be created.

* `vm_priority` - (Required) The priority of the VM. Changing this forces a new Machine Learning Compute Cluster to be created. Accepted values are `Dedicated` and `LowPriority`.

* `vm_size` - (Required) The size of the VM. Changing this forces a new Machine Learning Compute Cluster to be created.

* `scale_settings` - (Required) A `scale_settings` block as defined below. Changing this forces a new Machine Learning Compute Cluster to be created.

---

* `ssh` - (Optional) Credentials for an administrator user account that will be created on each compute node. A `ssh` block as defined below. Changing this forces a new Machine Learning Compute Cluster to be created.

* `description` - (Optional) The description of the Machine Learning compute. Changing this forces a new Machine Learning Compute Cluster to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `local_auth_enabled` - (Optional) Whether local authentication methods is enabled. Defaults to `true`. Changing this forces a new Machine Learning Compute Cluster to be created.

* `node_public_ip_enabled` - (Optional) Whether the compute cluster will have a public ip. Defaults to `true`. Changing this forces a new Machine Learning Compute Cluster to be created.

* `ssh_public_access_enabled` - (Optional) A boolean value indicating whether enable the public SSH port. Defaults to `false`. Changing this forces a new Machine Learning Compute Cluster to be created.

* `subnet_resource_id` - (Optional) The ID of the Subnet that the Compute Cluster should reside in. Changing this forces a new Machine Learning Compute Cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Compute Cluster. Changing this forces a new Machine Learning Compute Cluster to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Machine Learning Compute Cluster. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both). Changing this forces a new resource to be created.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Machine Learning Compute Cluster. Changing this forces a new resource to be created.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---
A `ssh` block supports the following:

* `admin_username` - (Required) Name of the administrator user account which can be used to SSH to nodes. Changing this forces a new Machine Learning Compute Cluster to be created.

* `admin_password` - (Optional) Password of the administrator user account. Changing this forces a new Machine Learning Compute Cluster to be created.

* `key_value` - (Optional) SSH public key of the administrator user account. Changing this forces a new Machine Learning Compute Cluster to be created.

~> **Note:** At least one of `admin_password` and `key_value` shoud be specified.

---

A `scale_settings` block supports the following:

* `max_node_count` - (Required) Maximum node count. Changing this forces a new Machine Learning Compute Cluster to be created.

* `min_node_count` - (Required) Minimal node count. Changing this forces a new Machine Learning Compute Cluster to be created.

* `scale_down_nodes_after_idle_duration` - (Required) Node Idle Time Before Scale Down: defines the time until the compute is shutdown when it has gone into Idle state. Is defined according to W3C XML schema standard for duration. Changing this forces a new Machine Learning Compute Cluster to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Machine Learning Compute Cluster.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Machine Learning Compute Cluster.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Machine Learning Compute Cluster.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Machine Learning Compute Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Compute Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Compute Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Compute Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Compute Cluster.

## Import

Machine Learning Compute Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_compute_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1
```
