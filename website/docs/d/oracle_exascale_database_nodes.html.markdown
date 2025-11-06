---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_db_nodes"
description: |-
  This data source provides the list of DB Nodes.
---

# Data Source: azurerm_oracle_exascale_database_nodes

Lists the database nodes for the specified Exadb VM Cluster.

## Example Usage

```hcl
data "azurerm_oracle_exascale_database_nodes" "example" {
  exadb_vm_cluster_id = "existing"
}

output "example" {
  value = data.azurerm_oracle_exascale_database_nodes.example
}
```

## Arguments Reference

The following arguments are supported:

* `exascale_database_virtual_machine_cluster_id` - (Required) The id of the Exadb VM cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `exascale_database_nodes` - A `exascale_database_node` block as defined below.

---

A `exascale_database_node` block exports the following:

* `additional_details` - Additional information about the planned maintenance.

* `cpu_core_count` - The number of CPU cores enabled on the exascale DB node.

* `database_node_storage_size_in_gb` - The allocated local node storage in GBs on the exascale DB node.

* `fault_domain` - The name of the Fault Domain the instance is contained in.

* `hostname` - The hostname for the exascale db node.

* `lifecycle_state` - Information about the current lifecycle state.

* `maintenance_type` - The type of database node maintenance.

* `memory_size_in_gb` - The allocated memory in GBs on the DB Node.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exascale DB node.

* `software_storage_size_in_gb` - The size (in GB) of the block storage volume allocation for the DB system. This attribute applies only for virtual machine DB systems.

* `time_maintenance_window_end` - End date and time of maintenance window.

* `time_maintenance_window_start` - Start date and time of maintenance window.

* `total_cpu_core_count` - The total number of CPU cores reserved on the Db node.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DB Node.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-03-01
