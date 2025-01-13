---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_db_nodes"
description: |-
  This data source provides the list of DB Nodes.
---

# Data Source: azurerm_oracle_db_nodes

Lists the database nodes for the specified Cloud VM Cluster.

## Example Usage

```hcl
data "azurerm_oracle_db_nodes" "example" {
  cloud_vm_cluster_id = "existing"
}

output "example" {
  value = data.azurerm_oracle_db_nodes.example
}
```

## Arguments Reference

The following arguments are supported:

* `cloud_vm_cluster_id` - (Required) The id of the Cloud VM cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `db_nodes` - A `db_nodes` block as defined below.

---

A `db_nodes` block exports the following:

* `additional_details` - Additional information about the planned maintenance.

* `backup_ip_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the backup IP address associated with the database node. Use this OCID with either the [GetPrivateIp](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PrivateIp/GetPrivateIp) or the [GetPublicIpByPrivateIpId](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PublicIp/GetPublicIpByPrivateIpId) API to get the IP address needed to make a database connection.

* `backup_vnic_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the backup VNIC.

* `backup_vnic2id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the second backup VNIC.

* `cpu_core_count` - The number of CPU cores enabled on the DB node.

* `db_node_storage_size_in_gbs` - The allocated local node storage in GBs on the DB node.

* `db_server_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the ExaCC DB server associated with the database node.

* `db_system_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the DB system.

* `fault_domain` - The name of the Fault Domain the instance is contained in.

* `host_ip_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the host IP address associated with the database node. Use this OCID with either the [GetPrivateIp](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PrivateIp/GetPrivateIp) or the [GetPublicIpByPrivateIpId](https://docs.cloud.oracle.com/iaas/api/#/en/iaas/20160918/PublicIp/GetPublicIpByPrivateIpId) API to get the IP address needed to make a database connection.

* `lifecycle_details` - Information about the current lifecycle details.

* `lifecycle_state` - Information about the current lifecycle state.

* `maintenance_type` - The type of database node maintenance.

* `memory_size_in_gbs` - The allocated memory in GBs on the DB Node.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DB node.

* `software_storage_size_in_gb` - The size (in GB) of the block storage volume allocation for the DB system. This attribute applies only for virtual machine DB systems.

* `time_created` - The date and time that the DB node was created.

* `time_maintenance_window_end` - End date and time of maintenance window.

* `time_maintenance_window_start` - Start date and time of maintenance window.

* `vnic2id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the second VNIC.

* `vnic_id` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the VNIC.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DB Node.
