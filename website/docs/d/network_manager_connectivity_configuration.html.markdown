---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_network_manager_connectivity_configuration"
description: |-
  Gets information about an existing Network Manager Connectivity Configuration.
---

# Data Source: azurerm_network_manager_connectivity_configuration

Use this data source to access information about an existing Network Manager Connectivity Configuration.

## Example Usage

```hcl
data "azurerm_network_manager_connectivity_configuration" "example" {
  name               = "existing"
  network_manager_id = "TODO"
}

output "id" {
  value = data.azurerm_network_manager_connectivity_configuration.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Network Manager Connectivity Configuration.

* `network_manager_id` - (Required) The ID of the Network Manager.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Manager Connectivity Configuration.

* `applies_to_group` - An `applies_to_group` block as defined below.

* `connectivity_topology` - The connectivity topology type.

* `delete_existing_peering_enabled` - Whether to current existing Virtual Network Peering in the Connectivity Configuration affected scope.

* `description` - The description of the Connectivity Configuration.

* `global_mesh_enabled` - Whether global mesh is supported.

* `hub` - A `hub` block as defined below.

---

An `applies_to_group` block exports the following:

* `global_mesh_enabled` - Whether global mesh is supported.

* `group_connectivity` - The group connectivity type.

* `network_group_id` - The ID of the Network Manager Network Group.

* `use_hub_gateway` - Whether hub gateway is used.

---

A `hub` block exports the following:

* `resource_id` - The resource ID used as hub in Hub and Spoke topology.

* `resource_type` - The resource type used as hub in Hub and Spoke topology.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Connectivity Configuration.
