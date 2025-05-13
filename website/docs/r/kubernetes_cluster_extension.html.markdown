---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_extension"
description: |-
  Manages a Kubernetes Cluster Extension.
---

# azurerm_kubernetes_cluster_extension

Manages a Kubernetes Cluster Extension.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "example-aks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster_extension" "example" {
  name           = "example-ext"
  cluster_id     = azurerm_kubernetes_cluster.example.id
  extension_type = "microsoft.flux"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Kubernetes Cluster Extension. Changing this forces a new Kubernetes Cluster Extension to be created.

* `cluster_id` - (Required) Specifies the Cluster ID. Changing this forces a new Kubernetes Cluster Extension to be created.

* `extension_type` - (Required) Specifies the type of extension. It must be one of the extension types registered with Microsoft.KubernetesConfiguration by the Extension publisher. For more information, please refer to [Available Extensions for AKS](https://learn.microsoft.com/en-us/azure/aks/cluster-extensions?tabs=azure-cli#currently-available-extensions). Changing this forces a new Kubernetes Cluster Extension to be created.

* `configuration_protected_settings` - (Optional) Configuration settings that are sensitive, as name-value pairs for configuring this extension.

* `configuration_settings` - (Optional) Configuration settings, as name-value pairs for configuring this extension.

* `plan` - (Optional) A `plan` block as defined below. Changing this forces a new resource to be created.

* `release_train` - (Optional) The release train used by this extension. Possible values include but are not limited to `Stable`, `Preview`. Changing this forces a new Kubernetes Cluster Extension to be created.

* `release_namespace` - (Optional) Namespace where the extension release must be placed for a cluster scoped extension. If this namespace does not exist, it will be created. Changing this forces a new Kubernetes Cluster Extension to be created.

* `target_namespace` - (Optional) Namespace where the extension will be created for a namespace scoped extension. If this namespace does not exist, it will be created. Changing this forces a new Kubernetes Cluster Extension to be created.

* `version` - (Optional) User-specified version that the extension should pin to. If it is not set, Azure will use the latest version and auto upgrade it. Changing this forces a new Kubernetes Cluster Extension to be created.

---

A `plan` block supports the following:

* `name` - (Required) Specifies the name of the plan from the marketplace. Changing this forces a new Kubernetes Cluster Extension to be created.

* `product` - (Required) Specifies the product of the plan from the marketplace. Changing this forces a new Kubernetes Cluster Extension to be created.

* `publisher` - (Required) Specifies the publisher of the plan. Changing this forces a new Kubernetes Cluster Extension to be created.

* `promotion_code` - (Optional) Specifies the promotion code to use with the plan. Changing this forces a new Kubernetes Cluster Extension to be created.

* `version` - (Optional) Specifies the version of the plan from the marketplace. Changing this forces a new Kubernetes Cluster Extension to be created.

~> **Note:** When `plan` is specified, legal terms must be accepted for this item on this subscription before creating the Kubernetes Cluster Extension. The `azurerm_marketplace_agreement` resource or AZ CLI tool can be used to do this.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Cluster Extension.

* `aks_assigned_identity` - An `aks_assigned_identity` block as defined below.

* `current_version` - The current version of the extension.

---

An `aks_assigned_identity` block exports the following:

* `type` - The identity type.

* `principal_id` - The principal ID of resource identity.

* `tenant_id` - The tenant ID of resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Cluster Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Cluster Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Cluster Extension.

## Import

Kubernetes Cluster Extension can be imported using the `resource id` for different `cluster_resource_name`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_extension.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/extensions/extension1
```
