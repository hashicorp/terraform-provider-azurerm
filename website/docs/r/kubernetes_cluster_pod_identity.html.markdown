---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_pod_identity"
description: |-
  Manages the Pod Identities within an Azure Kubernetes Cluster.
---

# azurerm_kubernetes_cluster_pod_identity

Manages the Pod Identities within an Azure Kubernetes Cluster.

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

  linux_profile {
    admin_username = "acctestuser1"

    ssh_key {
      key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
    }
  }

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    load_balancer_sku = "standard"
  }
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_user_assigned_identity.example.id
  role_definition_name = "Managed Identity Operator"
  principal_id         = azurerm_kubernetes_cluster.example.identity.0.principal_id
}

resource "azurerm_kubernetes_cluster_pod_identity" "example" {
  cluster_id = azurerm_kubernetes_cluster.example.id

  pod_identity {
    name        = "name"
    namespace   = "ns"
    identity_id = azurerm_user_assigned_identity.example.id
  }

  depends_on = [
    azurerm_role_assignment.example
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the Managed Kubernetes Cluster in which to create the Pod Identity. Changing this forces a new resource to be created.

* `pod_identity` - (Optional) One or more `pod_identity` blocks as defined below. At least one of `pod_identity` and `exception` should be specified.

* `exception` - (Optional) One or more `exception` blocks as defined below. At least one of `pod_identity` and `exception` should be specified.

---

A `pod_identity` block supports the following:

* `name` - (Required) The name of the Pod Identity.

* `namespace` - (Required) The namespace where the Pod Identity should be created.

* `identity_id` - (Required) The ID of the user assigned identity.

---

An `exception` block supports the following:

* `name` - (Required) The name of the Pod Identity Exception.

* `namespace` - (Required) The namespace where the Pod Identity Exception should be created.

* `pod_labels` - (Required) A map of Pod labels, the Pod Identity Exception will take effect only if Pod labels match.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Cluster Pod Identity.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Cluster Pod Identity.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Pod Identity.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Cluster Pod Identity.

## Import

Kubernetes Clusters Pod Identity can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_pod_identity.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1
```
