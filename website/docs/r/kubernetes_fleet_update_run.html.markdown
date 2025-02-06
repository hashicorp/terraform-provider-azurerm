---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_fleet_update_run"
description: |-
  Manages a Kubernetes Fleet Update Run.
---

# azurerm_kubernetes_fleet_update_run

Manages a Kubernetes Fleet Update Run.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "westeurope"
}

resource "azurerm_kubernetes_fleet_manager" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "example"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_fleet_member" "example" {
  name                  = "example"
  kubernetes_fleet_id   = azurerm_kubernetes_fleet_manager.example.id
  kubernetes_cluster_id = azurerm_kubernetes_cluster.example.id
  group                 = "example-group"
}

resource "azurerm_kubernetes_fleet_update_run" "example" {
  name                        = "example"
  kubernetes_fleet_manager_id = azurerm_kubernetes_fleet_manager.example.id
  managed_cluster_update {
    upgrade {
      type               = "Full"
      kubernetes_version = "1.27"
    }
    node_image_selection {
      type = "Latest"
    }
  }
  stage {
    name = "example"
    group {
      name = "example-group"
    }
    after_stage_wait_in_seconds = 21
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Kubernetes Fleet Update Run. Changing this forces a new Kubernetes Fleet Update Run to be created.

* `kubernetes_fleet_manager_id` - (Required) The ID of the Fleet Manager. Changing this forces a new Kubernetes Fleet Update Run to be created.

* `managed_cluster_update` - (Required) A `managed_cluster_update` block as defined below.

* `fleet_update_strategy_id` - (Optional) The ID of the Fleet Update Strategy. Only one of `fleet_update_strategy_id` or `stage` can be specified. 

* `stage` - (Optional) One or more `stage` blocks as defined below. Only one of `stage` or `fleet_update_strategy_id` can be specified.

---

A `managed_cluster_update` block supports the following:

* `upgrade` - (Required) A `upgrade` block as defined below.

* `node_image_selection` - (Optional) A `node_image_selection` block as defined below.

---

A `upgrade` block supports the following:

* `type` - (Required) Specifies the type of upgrade to perform. Possible values are `Full` and `NodeImageOnly`.

* `kubernetes_version` - (Optional) Specifies the Kubernetes version to upgrade the member clusters to. This is required if `type` is set to `Full`.

---

A `node_image_selection` block supports the following:

* `type` - (Required) Specifies the node image upgrade type. Possible values are `Latest` and `Consistent`.

---

A `stage` block supports the following:

* `group` - (Required) One or more `group` blocks as defined below.

* `name` - (Required) The name which should be used for this stage.

* `after_stage_wait_in_seconds` - (Optional) Specifies the time in seconds to wait at the end of this stage before starting the next one.

---

A `group` block supports the following:

* `name` - (Required) The name which should be used for this group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Fleet Update Run.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Fleet Update Run.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Update Run.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Fleet Update Run.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Fleet Update Run.

## Import

Kubernetes Fleet Update Runs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_fleet_update_run.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/fleets/fleet1/updateRuns/updateRun1
```
