---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_cluster_maintenance_configuration"
description: |-
  Manages a Maintenance Configuration for a Kubernetes Cluster.
---

# azurerm_kubernetes_cluster_maintenance_configuration

Manages a Maintenance Configuration for a Kubernetes Cluster.

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

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }
}

resource "azurerm_kubernetes_cluster_maintenance_configuration" "example" {
  name                  = "example"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id

  maintenance_allowed {
    day        = "Monday"
    hour_slots = [1, 2]
  }

  maintenance_not_allowed_window {
    end   = "2021-11-30T12:00:00Z"
    start = "2021-11-26T03:00:00Z"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Maintenance Configuration for the Kubernetes Cluster. Changing this forces a new resource to be created.

* `kubernetes_cluster_id` - (Required) The ID of the Kubernetes Cluster where this Maintenance Configuration should applies to. Changing this forces a new resource to be created.

* `maintenance_allowed` - (Optional) One or more `maintenance_allowed` block as defined below.

* `maintenance_not_allowed_window` - (Optional) One or more `maintenance_not_allowed_window` block as defined below.

-> **Note:** At least one of `maintenance_not_allowed_window` and `maintenance_allowed` must be set.

---

An `maintenance_allowed` block exports the following:

* `day` - (Required) A day in a week. Possible values are `Sunday`, `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday` and `Saturday`.

* `hour_slots` - (Required) An array of hour slots in a day. Possible values are between `0` and `23`.

---

An `maintenance_not_allowed_window` block exports the following:

* `end` - (Required) The end of a time span, formatted as an RFC3339 string.

* `start` - (Required) The start of a time span, formatted as an RFC3339 string.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kubernetes Cluster Maintenance Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Cluster Maintenance Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Cluster Maintenance Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Cluster Maintenance Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Cluster Maintenance Configuration.

## Import

Kubernetes Cluster Maintenance Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_cluster_maintenance_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/maintenanceConfigurations/config1
```
