---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_exa_db_vm_cluster"
description: |-
  Gets information about an existing Exadata VM Cluster.
---

# Data Source: azurerm_oracle_exa_db_vm_cluster

Use this data source to access information about an existing Exadata VM Cluster.

## Example Usage

```hcl
data "azurerm_oracle_exa_db_vm_cluster" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_exa_db_vm_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Exadata VM Cluster.

* `resource_group_name` - (Required) The name of the Resource Group where the Exadata VM Cluster exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Exadata VM Cluster.

* `backup_subnet_cidr` - Client OCI backup subnet CIDR, default is `192.168.252.0/22`.

* `cluster_name` - The cluster name for Exadata VM Cluster.

* `data_collection_options` - A `data_collection_options` block as defined below.

* `display_name` - The user-friendly name for the Exadata VM Cluster. The name does not need to be unique.

* `domain` - The domain name for the Exadata VM Cluster.

* `enabled_ecpu_count` - The number of ECPUs to enable for an Exadata VM cluster on Exascale Infrastructure.

* `exascale_db_storage_vault_id` - The OCID of the Exadata Database Storage Vault.

* `gi_version` - A valid Oracle Grid Infrastructure (GI) software version.

* `grid_image_ocid` - Grid Setup will be done using this grid image id.

* `grid_image_type` - The type of Grid Image

* `hostname` - The hostname for the Exadata VM Cluster without suffix.

* `hostname_actual` - The hostname for the Exadata VM Cluster with suffix.

* `iorm_config_cache` - A `iorm_config_cache` block as defined below.

* `license_model` - The Oracle license model that applies to the Exadata VM Cluster.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - Exadata VM Cluster lifecycle state enum.

* `listener_port` - The port number configured for the listener on the Exadata VM Cluster.

* `location` - The Azure Region where the Exadata VM Cluster exists.

* `memory_size_in_gbs` - The memory to be allocated in GBs.

* `node_count` - The number of nodes in the Exadata VM Cluster.

* `nsg_cidrs` - CIDR blocks for additional NSG ingress rules. The VNET CIDRs used to provision the VM Cluster will be added by default.

* `nsg_url` - The list of [OCIDs](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) for the network security groups (NSGs) to which this resource belongs. Setting this to an empty list removes all resources from all NSGs. For more information about NSGs, see [Security Rules](https://docs.oracle.com/en-us/iaas/Content/Network/Concepts/securityrules.htm). NsgIds restrictions:

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exadata VM Cluster.

* `private_zone_ocid` - The private zone ID in which you want DNS records to be created.

* `scan_dns_name` - The FQDN of the DNS record for the SCAN IP addresses that are associated with the Exadata VM Cluster.

* `scan_dns_record_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DNS record for the SCAN IP addresses that are associated with the Exadata VM Cluster.

* `scan_ip_ids` - A `scan_ip_ids` block as defined below.

* `scan_listener_port_tcp` -  The TCP Single Client Access Name (SCAN) port. The default port is 1521.

* `scan_listener_port_tcp_ssl` - The TCPS Single Client Access Name (SCAN) port. The default port is 2484.

* `shape` - The model name of the Exadata hardware running the Exadata VM Cluster.

* `snapshot_file_system_storage` - A `file_system_storage_details` block as defined below.

* `ssh_public_keys` - The public key portion of one or more key pairs used for SSH access to the Exadata VM Cluster.

* `subnet_id` - The ID of the Azure Resource Manager subnet resource.

* `subnet_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the subnet associated with the Exadata VM Cluster.

* `system_version` - Operating system version of the image.

* `tags` - A mapping of tags assigned to the Exadata VM Cluster.

* `time_zone` - The time zone of the Exadata VM Cluster. For details, see [Exadata Infrastructure Time Zones](https://docs.oracle.com/en-us/iaas/base-database/doc/manage-time-zone.html).

* `total_ecpu_count` - The number of Total ECPUs for an Exadata VM cluster on Exascale Infrastructure.

* `total_file_system_storage` - A `file_system_storage_details` block as defined below.

* `vip_ids` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the virtual IP (VIP) addresses associated with the Exadata VM Cluster. The Cluster Ready Services (CRS) creates and maintains one VIP address for each node in the Exadata Cloud Service instance to enable failover. If one node fails, the VIP is reassigned to another active node in the Cluster. 

* `vm_file_system_storage` - A `file_system_storage_details` block as defined below.

* `virtual_network_id` - The ID to an Azure Resource Manager Virtual Network resource.

* `zone_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the zone the Exadata VM Cluster is associated with.

---

A `data_collection_options` block exports the following:

* `diagnostics_events_enabled` - Indicates whether diagnostic collection is enabled for the VM Cluster/Exadata VM Cluster/VMBM DBCS. Enabling diagnostic collection allows you to receive Events service notifications for guest VM issues. Diagnostic collection also allows Oracle to provide enhanced service and proactive support for your Exadata system. You can enable diagnostic collection during VM Cluster/Exadata VM Cluster provisioning. You can also disable or enable it at any time using the `UpdateVmCluster` or `updateExadataVmCluster` API.

* `health_monitoring_enabled` -  Indicates whether health monitoring is enabled for the VM Cluster / Exadata VM Cluster / VMBM DBCS. Enabling health monitoring allows Oracle to collect diagnostic data and share it with its operations and support personnel. You may also receive notifications for some events. Collecting health diagnostics enables Oracle to provide proactive support and enhanced service for your system. Optionally enable health monitoring while provisioning a system. You can also disable or enable health monitoring anytime using the `UpdateVmCluster`, `UpdateExadataVmCluster` or `updateDbsystem` API.

* `incident_logs_enabled` - Indicates whether incident logs and trace collection are enabled for the VM Cluster / Exadata VM Cluster / VMBM DBCS. Enabling incident logs collection allows Oracle to receive Events service notifications for guest VM issues, collect incident logs and traces, and use them to diagnose issues and resolve them. Optionally enable incident logs collection while provisioning a system. You can also disable or enable incident logs collection anytime using the `UpdateVmCluster`, `updateExadataVmCluster` or `updateDbsystem` API.

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

---

A `file_system_storage_details` block exports the following:

* `total_size_in_gbs` - Total Capacity 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata VM Cluster.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database`: 2025-03-01
