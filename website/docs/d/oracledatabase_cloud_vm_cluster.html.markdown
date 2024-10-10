---
subcategory: "App Service"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracledatabase_cloud_vm_cluster"
description: |-
  Gets information about an existing Cloud VM Cluster.
---

# Data Source: azurerm_oracledatabase_cloud_vm_cluster

Use this data source to access information about an existing Cloud VM Cluster.

## Example Usage

```hcl
data "azurerm_oracledatabase_cloud_vm_cluster" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracledatabase_cloud_vm_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Cloud VM Cluster.

* `resource_group_name` - (Required) The name of the Resource Group where the Cloud VM Cluster exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cloud VM Cluster.

* `backup_subnet_cidr` - Client OCI backup subnet CIDR, default is 192.168.252.0/22.

* `cloud_exadata_infrastructure_id` - Cloud Exadata Infrastructure ID.

* `cluster_name` - The cluster name for cloud VM cluster.

* `compartment_id` - The OCID of the compartment.

* `compute_nodes` - A `compute_nodes` block as defined below.

* `cpu_core_count` - The number of CPU cores enabled on the cloud VM cluster.

* `data_collection_options` - A `data_collection_options` block as defined below.

* `data_storage_percentage` - The percentage assigned to DATA storage (user data and database files). The remaining percentage is assigned to RECO storage (database redo logs, archive logs, and recovery manager backups). Accepted values are 35, 40, 60 and 80. The default is 80 percent assigned to DATA storage. See [Storage Configuration](https://docs.oracle.com/en-us/iaas/exadatacloud/index.html#Exadata) in the Exadata documentation for details on the impact of the configuration settings on storage.

* `data_storage_size_in_tbs` - The data disk group size to be allocated in TBs.

* `db_node_storage_size_in_gbs` - The local node storage to be allocated in GBs.

* `db_servers` - A `db_servers` block as defined below.

* `disk_redundancy` - The type of redundancy configured for the cloud Vm cluster. NORMAL is 2-way redundancy. HIGH is 3-way redundancy.

* `display_name` - The user-friendly name for the cloud VM cluster. The name does not need to be unique.

* `domain` - The domain name for the cloud VM cluster.

* `gi_version` - A valid Oracle Grid Infrastructure (GI) software version.

* `hostname` - The hostname for the cloud VM cluster.

* `iorm_config_cache` - A `iorm_config_cache` block as defined below.

* `is_local_backup_enabled` - If true, database backup on local Exadata storage is configured for the cloud VM cluster. If false, database backup on local Exadata storage is not available in the cloud VM cluster.

* `is_sparse_diskgroup_enabled` - If true, sparse disk group is configured for the cloud VM cluster. If false, sparse disk group is not created.

* `last_update_history_entry_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the last maintenance update history entry. This value is updated when a maintenance update starts.

* `license_model` - The Oracle license model that applies to the cloud VM cluster.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - Cloud VM Cluster lifecycle state enum.

* `listener_port` - The port number configured for the listener on the cloud VM cluster.

* `location` - The Azure Region where the Cloud VM Cluster exists.

* `memory_size_in_gbs` - The memory to be allocated in GBs.

* `node_count` - The number of nodes in the cloud VM cluster.

* `nsg_url` - The list of [OCIDs](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) for the network security groups (NSGs) to which this resource belongs. Setting this to an empty list removes all resources from all NSGs. For more information about NSGs, see [Security Rules](https://docs.oracle.com/en-us/iaas/Content/Network/Concepts/securityrules.htm). NsgIds restrictions:
  * A network security group (NSG) is optional for Autonomous Databases with private access. The nsgIds list can be empty.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the cloud VM cluster.

* `ocpu_count` - The number of OCPU cores to enable on the cloud VM cluster. Only 1 decimal place is allowed for the fractional part.

* `provisioning_state` - Azure resource provisioning state.

* `scan_dns_name` - The FQDN of the DNS record for the SCAN IP addresses that are associated with the cloud VM cluster.

* `scan_dns_record_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DNS record for the SCAN IP addresses that are associated with the cloud VM cluster.

* `scan_ip_ids` - A `scan_ip_ids` block as defined below.

* `scan_listener_port_tcp` -  The TCP Single Client Access Name (SCAN) port. The default port is 1521.

* `scan_listener_port_tcp_ssl` - The TCPS Single Client Access Name (SCAN) port. The default port is 2484.

* `shape` - The model name of the Exadata hardware running the cloud VM cluster.

* `ssh_public_keys` - The public key portion of one or more key pairs used for SSH access to the cloud VM cluster.

* `storage_size_in_gbs` - The storage allocation for the disk group, in gigabytes (GB).

* `subnet_id` - The ID of the Azure Resource Manager subnet resource.

* `subnet_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the subnet associated with the cloud VM cluster.

* `system_data` - A `system_data` block as defined below.

* `system_version` - Operating system version of the image.

* `tags` - A mapping of tags assigned to the Cloud VM Cluster.

* `time_created` - The date and time that the cloud VM cluster was created.

* `time_zone` - The time zone of the cloud VM cluster. For details, see [Exadata Infrastructure Time Zones](https://docs.oracle.com/en-us/iaas/base-database/doc/manage-time-zone.html).

* `vip_ods` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the virtual IP (VIP) addresses associated with the cloud VM cluster. The Cluster Ready Services (CRS) creates and maintains one VIP address for each node in the Exadata Cloud Service instance to enable failover. If one node fails, the VIP is reassigned to another active node in the cluster. **Note:** For a single-node DB system, this list is empty.

* `vnet_id` - The ID to an Azure Resource Manager vnet resource.

* `zone_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the zone the cloud VM cluster is associated with.

---

A `data_collection_options` block exports the following:

* `is_diagnostics_events_enabled` - Indicates whether diagnostic collection is enabled for the VM cluster/Cloud VM cluster/VMBM DBCS. Enabling diagnostic collection allows you to receive Events service notifications for guest VM issues. Diagnostic collection also allows Oracle to provide enhanced service and proactive support for your Exadata system. You can enable diagnostic collection during VM cluster/Cloud VM cluster provisioning. You can also disable or enable it at any time using the `UpdateVmCluster` or `updateCloudVmCluster` API.

* `is_health_monitoring_enabled` -  Indicates whether health monitoring is enabled for the VM cluster / Cloud VM cluster / VMBM DBCS. Enabling health monitoring allows Oracle to collect diagnostic data and share it with its operations and support personnel. You may also receive notifications for some events. Collecting health diagnostics enables Oracle to provide proactive support and enhanced service for your system. Optionally enable health monitoring while provisioning a system. You can also disable or enable health monitoring anytime using the `UpdateVmCluster`, `UpdateCloudVmCluster` or `updateDbsystem` API.

* `is_incident_logs_enabled` - Indicates whether incident logs and trace collection are enabled for the VM cluster / Cloud VM cluster / VMBM DBCS. Enabling incident logs collection allows Oracle to receive Events service notifications for guest VM issues, collect incident logs and traces, and use them to diagnose issues and resolve them. Optionally enable incident logs collection while provisioning a system. You can also disable or enable incident logs collection anytime using the `UpdateVmCluster`, `updateCloudVmCluster` or `updateDbsystem` API.

---

A `db_plans` block exports the following:

* `db_name` - The database name. For the default `DbPlan`, the `dbName` is `default`.

* `flash_cache_limit` - The flash cache limit for this database. This value is internally configured based on the share value assigned to the database.

* `share` - The relative priority of this database.

---

A `iorm_config_cache` block exports the following:

* `db_plans` - A `db_plans` block as defined above.

* `lifecycle_details` - Additional information about the current `lifecycleState`.

* `lifecycle_state` - The current state of IORM configuration for the Exadata DB system.

* `objective` - The current value for the IORM objective. The default is `AUTO`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cloud VM Cluster.
