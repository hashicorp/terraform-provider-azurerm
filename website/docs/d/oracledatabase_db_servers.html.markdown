---
subcategory: "App Service"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracledatabase_db_servers"
description: |-
  Gets information about an existing DB Server.
---

# Data Source: azurerm_oracledatabase_db_servers

Use this data source to access information about an existing DB Server.

## Example Usage

```hcl
data "azurerm_oracledatabase_db_servers" "example" {
  resource_group_name               = "existing"
  cloud_exadata_infrastructure_name = "existing"
}

output "id" {
  value = data.azurerm_oracledatabase_db_servers.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `cloud_exadata_infrastructure_name` - (Required) The name of the ExadataInfrastructure.

* `resource_group_name` - (Required) The name of the Resource Group where the DB Server exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exacc Db server.

* `db_servers` - A `db_servers` block as defined below.

---

A `db_servers` block exports the following:

* `autonomous_virtual_machine_ds` - The list of [OCIDs](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Autonomous Virtual Machines associated with the Db server.

* `autonomous_vm_cluster_ids` - The list of [OCIDs](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Autonomous VM Clusters associated with the Db server.

* `compartment_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the compartment.

* `cpu_core_count` - The number of CPU cores enabled on the Db server.

* `db_node_ids` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Db nodes associated with the Db server.

* `db_node_storage_size_in_gbs` - The allocated local node storage in GBs on the Db server.

* `display_name` - The user-friendly name for the Db server. The name does not need to be unique.

* `exadata_infrastructure_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exadata infrastructure.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - The current state of the DB server.

* `max_cpu_count` - The total number of CPU cores available.

* `max_db_node_storage_in_gbs` -The total local node storage available in GBs.

* `max_memory_in_gbs` - The total memory available in GBs.

* `memory_size_in_gbs` - The allocated memory in GBs on the Db server.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DB server.

* `provisioning_state` - DbServer provisioning state.

* `shape` - The shape of the Db server. The shape determines the amount of CPU, storage, and memory resources available.

* `time_created` - The date and time that the Db Server was created.

* `vm_cluster_ids` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the VM Clusters associated with the Db server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DB Server.
