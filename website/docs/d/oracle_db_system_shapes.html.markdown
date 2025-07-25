---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_db_system_shapes"
description: |-
  Provides the list of DB System Shapes.
---

# Data Source: azurerm_oracle_db_system_shapes

This data source provides the list of DB System Shapes in Oracle Cloud Infrastructure Database service.

Gets a list of the shapes that can be used to launch a new DB system. The shape determines resources to allocate to the DB system - CPU cores and memory for VM shapes; CPU cores, memory and storage for non-VM (or bare metal) shapes.

## Example Usage

```hcl
data "azurerm_oracle_db_system_shapes" "example" {
  location = "eastus"
  zone     = "2"
}

output "example" {
  value = data.azurerm_oracle_db_system_shapes.example
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region to query for the system shapes in.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `db_system_shapes` - A `db_system_shapes` block as defined below.

---

A `db_system_shapes` block exports the following:

* `are_server_types_supported` - Indicates if the shape supports database and storage server types.

* `available_core_count` - The maximum number of CPU cores that can be enabled on the DB system for this shape.

* `available_core_count_per_node` - The maximum number of CPU cores per database node that can be enabled for this shape. Only applicable to the flex Exadata shape, ExaCC Elastic shapes and VM Flex shapes.

* `available_data_storage_in_tbs` - The maximum data storage that can be enabled for this shape.

* `available_data_storage_per_server_in_tbs` - The maximum data storage available per storage server for this shape. Only applicable to ExaCC Elastic shapes.

* `available_db_node_per_node_in_gbs` - The maximum DB Node storage available per database node for this shape. Only applicable to ExaCC Elastic shapes.

* `available_db_node_storage_in_gbs` - The maximum DB Node storage that can be enabled for this shape.

* `available_memory_in_gbs` - The maximum memory that can be enabled for this shape.

* `available_memory_per_node_in_gbs` - The maximum memory available per database node for this shape. Only applicable to ExaCC Elastic shapes.

* `compute_model` - The compute model of the Exadata Infrastructure.

* `core_count_increment` - The discrete number by which the CPU core count for this shape can be increased or decreased.

* `display_name` - The display name of the shape used for the DB system.

* `maximum_storage_count` - The maximum number of Exadata storage servers available for the Exadata infrastructure.

* `maximum_node_count` - The maximum number of compute servers available for this shape.

* `minimum_core_count_per_node` - The minimum number of CPU cores that can be enabled per node for this shape.

* `minimum_data_storage_in_tbs` - The minimum data storage that need be allocated for this shape.

* `minimum_db_node_storage_per_node_in_gbs` - The minimum DB Node storage that need be allocated per node for this shape.

* `minimum_memory_per_node_in_gbs` - The minimum memory that need be allocated per node for this shape.

* `minimum_storage_count` - The minimum number of Exadata storage servers available for the Exadata infrastructure.

* `minimum_core_count` - The minimum number of CPU cores that can be enabled on the DB system for this shape.

* `minimum_node_count` - The minimum number of compute servers available for this shape.

* `runtime_minimum_core_count` - The runtime minimum number of compute servers available for this shape.

* `shape_family` - The family of the shape used for the DB system.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the System Shapes.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database`: 2025-03-01
