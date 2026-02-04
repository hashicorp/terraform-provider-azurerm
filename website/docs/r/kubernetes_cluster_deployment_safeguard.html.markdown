---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_deployment_safeguard"
description: |-
  Manages a Deployment Safeguard for a Kubernetes Cluster.
---

# azurerm_kubernetes_cluster_deployment_safeguard

Manages a Deployment Safeguard for a Kubernetes Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "exampleaks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
    upgrade_settings {
      max_surge = "10%"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  azure_policy_enabled = true
}

resource "azurerm_kubernetes_cluster_deployment_safeguard" "example" {
  kubernetes_cluster_id        = azurerm_kubernetes_cluster.example.id
  level                        = "Enforce"
  excluded_namespaces          = ["my-app-namespace", "legacy-app"]
  pod_security_standards_level = "Restricted"
}
```

## Arguments Reference

The following arguments are supported:

* `kubernetes_cluster_id` - (Required) Specifies the Kubernetes Cluster ID for which Deployment Safeguards should be configured. Changing this forces a new resource to be created.

* `level` - (Required) The level of Deployment Safeguards enforcement. Possible values are `Warn` and `Enforce`.

---

* `excluded_namespaces` - (Optional) A set of Kubernetes namespace names that should be excluded from Deployment Safeguards enforcement. This allows certain namespaces to bypass the configured policies.

* `pod_security_standards_level` - (Optional) The Pod Security Standards level to enforce. Possible values are `Baseline`, `Privileged`, and `Restricted`. Defaults to `Privileged` if not specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Cluster (same as `kubernetes_cluster_id`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Deployment Safeguard.
* `read` - (Defaults to 5 minutes) Used when retrieving the Deployment Safeguard.
* `update` - (Defaults to 30 minutes) Used when updating the Deployment Safeguard.
* `delete` - (Defaults to 30 minutes) Used when deleting the Deployment Safeguard.

## Import

A Deployment Safeguard for a Kubernetes Cluster can be imported using the Kubernetes Cluster `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_deployment_safeguard.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ContainerService` - 2025-07-01
