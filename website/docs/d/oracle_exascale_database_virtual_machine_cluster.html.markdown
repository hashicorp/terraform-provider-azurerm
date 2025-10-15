---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_exascale_database_virtual_machine_cluster"
description: |-
  Gets information about an existing Exadata VM Cluster.
---

# Data Source: azurerm_oracle_exascale_database_virtual_machine_cluster

Use this data source to access information about an existing Exadata VM Cluster.

## Example Usage

```hcl
data "azurerm_oracle_exascale_database_virtual_machine_cluster" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_exascale_database_virtual_machine_cluster.example.id
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

* `data_collection_option` - A `data_collection_option` block as defined below.

* `display_name` - The user-friendly name for the Exadata VM Cluster. The name does not need to be unique.

* `domain` - The domain name for the Exadata VM Cluster.

* `enabled_ecpu_count` - The number of ECPUs to enable for an Exadata VM cluster on Exascale Infrastructure.

* `exascale_database_storage_vault_id` - The OCID of the Exadata Database Storage Vault.

* `grid_infrastructure_version` - A valid Oracle Grid Infrastructure (GI) software version.

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

* `network_security_group_cidrs` - A `network_security_group_cidr` block as defined below.

* `network_security_group_url` - The link to OCI Network Security Group exposed to Azure Customer via the Azure Interface.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exadata VM Cluster.

* `private_zone_ocid` - The private zone ID in which you want DNS records to be created.

* `scan_dns_name` - The FQDN of the DNS record for the SCAN IP addresses that are associated with the Exadata VM Cluster.

* `scan_dns_record_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DNS record for the SCAN IP addresses that are associated with the Exadata VM Cluster.

* `scan_ip_ids` - A `scan_ip_ids` block as defined below.

* `single_client_access_name_listener_port_tcp` -  The TCP Single Client Access Name (SCAN) port. The default port is 1521.

* `single_client_access_name_listener_port_tcp_ssl` - The TCPS Single Client Access Name (SCAN) port. The default port is 2484.

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

* `virtual_ip_ids` - The virtual IP (VIP) addresses associated with the Exadata VM cluster on Exascale Infrastructure.

* `virtual_machine_file_system_storage` - A `file_system_storage_details` block as defined below.

* `virtual_network_id` - The ID to an Azure Resource Manager Virtual Network resource.

* `zone_ocid` - The OCID of the zone the Exadata VM cluster on Exascale Infrastructure is associated with.

---

A `data_collection_option` block exports the following:

* `diagnostics_events_enabled` - Indicates whether diagnostic collection is enabled for the VM cluster/Cloud VM cluster/VMBM DBCS. Changing this forces a new Cloud VM Cluster to be created.

* `health_monitoring_enabled` -  Indicates whether health monitoring is enabled for the VM cluster / Cloud VM cluster / VMBM DBCS. Changing this forces a new Cloud VM Cluster to be created.

* `incident_logs_enabled` - Indicates whether incident logs and trace collection are enabled for the VM cluster / Cloud VM cluster / VMBM DBCS. Changing this forces a new Cloud VM Cluster to be created.

---

A `network_security_group_cidr` block exports the following:

* `destination_port_range` - A `destination_port_range` block as defined below.

* `source` - It is a range of IP addresses that a packet coming into the instance can come from.

---

A `destination_port_range` block exports the following:

* `max` - The maximum port number, which must not be less than the minimum port number.

* `min` - The minimum port number, which must not be greater than the maximum port number.

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

* `total_size_in_gb` - Total Capacity 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata VM Cluster.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
