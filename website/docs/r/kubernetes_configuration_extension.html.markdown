---
subcategory: "KubernetesConfiguration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_configuration_extension"
description: |-
  Manages a Kubernetes Configuration Extension.
---

# azurerm_kubernetes_configuration_extension

Manages a Kubernetes Configuration Extension.

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

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_configuration_extension" "example" {
  name                  = "example-akce1"
  resource_group_name   = azurerm_resource_group.example.name
  cluster_name          = azurerm_kubernetes_cluster.example.name
  cluster_resource_name = "managedClusters"
  extension_type        = "microsoft.flux"

  configuration_protected_settings = {
    "omsagent.secret.key" = "secretKeyValue01"
  }

  configuration_settings = {
    "omsagent.secret.wsid"     = "a38cef99-5a89-52ed-b6db-22095c23664b",
    "omsagent.env.clusterName" = "clusterName1"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Kubernetes Configuration Extension which must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number. Changing this forces a new Kubernetes Configuration Extension to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Kubernetes Configuration Extension should exist. Changing this forces a new Kubernetes Configuration Extension to be created.

* `cluster_name` - (Required) The name of the kubernetes cluster. Changing this forces a new Kubernetes Configuration Extension to be created.

* `cluster_resource_name` - (Required) The kubernetes cluster resource name. Current only possible value is `managedClusters`. Changing this forces a new Kubernetes Configuration Extension to be created.

* `extension_type` - (Required) Type of the Extension, of which this resource is an instance of.  It must be one of the Extension Types registered with Microsoft.KubernetesConfiguration by the Extension publisher. Possible Values are `microsoft.flux`, `microsoft.dapr`, `microsoft.azureml.kubernetes`, etc. Changing this forces a new Kubernetes Configuration Extension to be created.

* `auto_upgrade_minor_version` - (Optional) Flag to note if this extension participates in auto upgrade of minor version, or not. Defaults to `true`.

* `configuration_protected_settings` - (Optional) Configuration settings that are sensitive, as name-value pairs for configuring this extension.

* `configuration_settings` - (Optional) Configuration settings, as name-value pairs for configuring this extension.

* `release_train` - (Optional) ReleaseTrain this extension participates in for auto-upgrade. Possible Values are `Stable`, `Preview`, etc. It should be set only when `auto_upgrade_minor_version` is `true`.

* `version` - (Optional) Version of the extension for this extension, if it is 'pinned' to a specific version. It should be set when and only when `auto_upgrade_minor_version` is `false`.

* `release_namespace` - (Optional) Namespace where the extension Release must be placed, for a Cluster scoped extension.  If this namespace does not exist, it will be created. Changing this forces a new Kubernetes Configuration Extension to be created.

* `target_namespace` - (Optional) Namespace where the extension will be created for a namespace scoped extension.  If this namespace does not exist, it will be created. Changing this forces a new Kubernetes Configuration Extension to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Configuration Extension.

* `aks_assigned_identity` - An `aks_assigned_identity` block as defined below.

---

An `aks_assigned_identity` block exports the following:

* `type` - The type of Managed Service Identity.

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Configuration Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Configuration Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Configuration Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Configuration Extension.

## Import

Kubernetes Configuration Extensions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_configuration_extension.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/clusterName1/providers/Microsoft.KubernetesConfiguration/extensions/extension1
```
