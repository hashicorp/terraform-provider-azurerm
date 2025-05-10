---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_trusted_access_role_binding"
description: |-
  Manages a Kubernetes Cluster Trusted Access Role Binding.
---

# azurerm_kubernetes_cluster_trusted_access_role_binding

## Example Usage

```hcl
resource "azurerm_application_insights" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "example-value"
}
data "azurerm_client_config" "test" {}
resource "azurerm_key_vault" "example" {
  name                       = "example"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.example.tenant_id
  sku_name                   = "example-value"
  soft_delete_retention_days = "example-value"
}
resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id    = azurerm_key_vault.example.id
  tenant_id       = data.azurerm_client_config.example.tenant_id
  object_id       = data.azurerm_client_config.example.object_id
  key_permissions = "example-value"
}
resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "acctestaksexample"
  default_node_pool {
    name       = "example-value"
    node_count = "example-value"
    vm_size    = "example-value"
    upgrade_settings {
      max_surge = "example-value"
    }
  }
  identity {
    type = "example-value"
  }
}
resource "azurerm_machine_learning_workspace" "example" {
  name                    = "example"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  key_vault_id            = azurerm_key_vault.example.id
  storage_account_id      = azurerm_storage_account.example.id
  application_insights_id = azurerm_application_insights.example.id
  identity {
    type = "example-value"
  }
}
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_storage_account" "example" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "example-value"
  account_replication_type = "example-value"
}
resource "azurerm_kubernetes_cluster_trusted_access_role_binding" "example" {
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  name                  = "example"
  roles                 = "example-value"
  source_resource_id    = azurerm_machine_learning_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `kubernetes_cluster_id` - (Required) Specifies the Kubernetes Cluster Id within which this Kubernetes Cluster Trusted Access Role Binding should exist. Changing this forces a new Kubernetes Cluster Trusted Access Role Binding to be created.

* `name` - (Required) Specifies the name of this Kubernetes Cluster Trusted Access Role Binding. Changing this forces a new Kubernetes Cluster Trusted Access Role Binding to be created.

* `roles` - (Required) A list of roles to bind, each item is a resource type qualified role name.

* `source_resource_id` - (Required) The ARM resource ID of source resource that trusted access is configured for. Changing this forces a new Kubernetes Cluster Trusted Access Role Binding to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Cluster Trusted Access Role Binding.

---



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Cluster Trusted Access Role Binding.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Trusted Access Role Binding.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Cluster Trusted Access Role Binding.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Cluster Trusted Access Role Binding.

## Import

An existing Kubernetes Cluster Trusted Access Role Binding can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_trusted_access_role_binding.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{managedClusterName}/trustedAccessRoleBindings/{trustedAccessRoleBindingName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Kubernetes Cluster Trusted Access Role Binding exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Kubernetes Cluster Trusted Access Role Binding exists. For example `example-resource-group`.
* Where `{managedClusterName}` is the name of the Managed Cluster. For example `managedClusterValue`.
* Where `{trustedAccessRoleBindingName}` is the name of the Trusted Access Role Binding. For example `trustedAccessRoleBindingValue`.
