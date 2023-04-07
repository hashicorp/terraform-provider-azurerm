---
subcategory: "Kubernetes Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_kubernetes_cluster_extension"
description: |-
  Manages an Arc Kubernetes Cluster Extension.
---

# azurerm_arc_kubernetes_cluster_extension

Manages an Arc Kubernetes Cluster Extension.

-> **Note:** Installing and configuring the Azure Arc Agent on your Kubernetes Cluster to establish connectivity is outside the scope of this document. For more details refer to [Deploy agents to your cluster](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/conceptual-agent-overview#deploy-agents-to-your-cluster) and [Connect an existing Kubernetes Cluster](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/quickstart-connect-cluster?tabs=azure-cli#connect-an-existing-kubernetes-cluster). If you encounter issues connecting your Kubernetes Cluster to Azure Arc, we'd recommend opening a ticket with Microsoft Support.

-> **Note:** An extensive example on connecting the `azurerm_arc_kubernetes_cluster` to an external kubernetes cluster can be found in [the `./examples/arckubernetes` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/arckubernetes).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_kubernetes_cluster" "example" {
  name                         = "example-akcc"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = "West Europe"
  agent_public_key_certificate = filebase64("testdata/public.cer")

  identity {
    type = "SystemAssigned"
  }

  tags = {
    ENV = "Test"
  }
}

resource "azurerm_arc_kubernetes_cluster_extension" "example" {
  name                = "example-ext"
  resource_group_name = azurerm_resource_group.example.name
  cluster_name        = azurerm_arc_kubernetes_cluster.example.name
  extension_type      = "microsoft.flux"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Arc Kubernetes Cluster Extension. It must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Arc Kubernetes Cluster Extension should exist. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `cluster_name` - (Required) Specifies the Cluster Name. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `extension_type` - (Required) Type of the Extension, of which this resource is an instance of. It must be one of the Extension Types registered with Microsoft.KubernetesConfiguration by the Extension publisher. For more information, please refer to [Available Extensions for Arc-enabled Kubernetes clusters](https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/extensions-release). Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `configuration_protected_settings` - (Optional) Configuration settings that are sensitive, as name-value pairs for configuring this extension.

* `configuration_settings` - (Optional) Configuration settings, as name-value pairs for configuring this extension.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `release_train` - (Optional) The release train which is used by this extension. Possible values include but not limited to `Stable`, `Preview`. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `release_namespace` - (Optional) Namespace where the extension Release must be placed, for a Cluster scoped extension.  If this namespace does not exist, it will be created. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `target_namespace` - (Optional) Namespace where the extension will be created for a Namespace scoped extension.  If this namespace does not exist, it will be created. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

* `version` - (Optional) User-specified version of the extension for this extension to 'pin'. If it is not set, Azure will use the latest version and auto upgrade it. Changing this forces a new Arc Kubernetes Cluster Extension to be created.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Arc Kubernetes Cluster Extension.

* `current_version` - Currently installed version of the extension.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Arc Kubernetes Cluster Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Arc Kubernetes Cluster Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Arc Kubernetes Cluster Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Arc Kubernetes Cluster Extension.

## Import

Arc Kubernetes Cluster Extension can be imported using the `resource id` for different `cluster_resource_name`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_extension.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Kubernetes/connectedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/extensions/extension1
```
