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

* `backup_subnet_cidr` - Client OCI backup subnet CIDR.

* `cluster_name` - The cluster name for Exadata VM Cluster.

* `data_collection` - A `data_collection` block as defined below.

* `display_name` - The user-friendly name for the Exadata VM Cluster.

* `domain` - The domain name for the Exadata VM Cluster.

* `enabled_ecpu_count` - The number of ECPUs enabled for an Exadata VM cluster on Exascale Infrastructure.

* `exascale_database_storage_vault_id` - The OCID of the Exadata Database Storage Vault.

* `grid_infrastructure_version` - The Oracle Grid Infrastructure (GI) software version.

* `grid_image_ocid` - The Grid Image [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) used for Grid setup.

* `grid_image_type` - The type of Grid Image

* `hostname` - The hostname for the Exadata VM Cluster.

* `hostname_actual` - The hostname for the Exadata VM Cluster with suffix.

* `iorm_config_cache` - A `iorm_config_cache` block as defined below.

* `license_model` - The Oracle license model that applies to the Exadata VM Cluster.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - Exadata VM Cluster lifecycle state enum.

* `listener_port` - The port number configured for the listener on the Exadata VM Cluster.

* `location` - The Azure Region where the Exadata VM Cluster exists.

* `memory_size_in_gb` - The memory capacity allocated for the Exadata VM Cluster in GB.

* `node_count` - The number of nodes in the Exadata VM Cluster.

* `inbound_network_security_group_rule` - A `inbound_network_security_group_rule` block as defined below.

* `network_security_group_url` - The link to OCI Network Security Group exposed to Azure Customer via the Azure Interface.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exadata VM Cluster.

* `private_zone_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the zone the Exadata VM cluster on Exascale Infrastructure is associated with.

* `single_client_access_name_dns_name` - The FQDN of the DNS record for the SCAN IP addresses that are associated with the Exadata VM Cluster.

* `single_client_access_name_dns_record_id` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DNS record for the SCAN IP addresses that are associated with the Exadata VM Cluster.

* `single_client_access_name_ip_ids` - The Single Client Access Name (SCAN) IP addresses associated with the Exadata VM cluster on Exascale Infrastructure.

* `single_client_access_name_listener_port_tcp` -  The TCP Single Client Access Name (SCAN) port.

* `single_client_access_name_listener_port_tcp_ssl` - The TCPS Single Client Access Name (SCAN) port.

* `shape` - The model name of the Exadata hardware running the Exadata VM Cluster.

* `snapshot_file_system_storage` - A `snapshot_file_system_storage` block as defined below.

* `ssh_public_keys` - The public key portion of one or more key pairs used for SSH access to the Exadata VM Cluster.

* `subnet_id` - The ID of the Azure Resource Manager subnet resource.

* `subnet_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the subnet associated with the Exadata VM Cluster.

* `system_version` - Operating system version of the image.

* `tags` - A mapping of tags assigned to the Exadata VM Cluster.

* `time_zone` - The time zone of the Exadata VM Cluster. For details, see [Exadata Infrastructure Time Zones](https://docs.oracle.com/en-us/iaas/base-database/doc/manage-time-zone.html).

* `total_ecpu_count` - The number of Total ECPUs for an Exadata VM cluster on Exascale Infrastructure.

* `total_file_system_storage` - A `total_file_system_storage` block as defined below.

* `virtual_ip_ids` - The virtual IP (VIP) addresses associated with the Exadata VM cluster on Exascale Infrastructure.

* `virtual_machine_file_system_storage` - A `virtual_machine_file_system_storage` block as defined below.

* `virtual_network_id` - The ID to an Azure Resource Manager Virtual Network resource.

* `zone_ocid` - The OCID of the zone the Exadata VM cluster on Exascale Infrastructure is associated with.

* `zones` - A list of Availability Zones in which this Exadata VM cluster is located.

---

A `data_collection` block exports the following:

* `diagnostics_events_enabled` - Indicates whether diagnostic collection is enabled for the VM cluster, Cloud VM cluster or VMBM DBCS.

* `health_monitoring_enabled` -  Indicates whether health monitoring is enabled for the VM cluster, Cloud VM cluster or VMBM DBCS.

* `incident_logs_enabled` - Indicates whether incident logs and trace collection are enabled for the VM cluster, Cloud VM cluster or VMBM DBCS.

---

A `inbound_network_security_group_rule` block exports the following:

* `destination_port_range` - A `destination_port_range` block as defined below.

* `source_ip_range` - It is a range of IP addresses that a packet coming into the instance can come from.

---

A `destination_port_range` block exports the following:

* `maximum` - The maximum port number.

* `minimum` - The minimum port number.

---

A `database_plans` block exports the following:

* `database_name` - The database name.

* `flash_cache_limit` - The flash cache limit for this database. This value is internally configured based on the share value assigned to the database.

* `share` - The relative priority of this database.

---

A `iorm_config_cache` block exports the following:

* `database_plans` - A `database_plans` block as defined above.

* `lifecycle_details` - Additional information about the current `lifecycle_state`.

* `lifecycle_state` - The current state of IORM configuration for the Exadata DB system.

* `objective` - The current value for the IORM objective.

---

A `snapshot_file_system_storage` block exports the following:

* `total_size_in_gb` - Total Capacity in GB.

---

A `total_file_system_storage` block exports the following:

* `total_size_in_gb` - Total Capacity in GB.

---

A `virtual_machine_file_system_storage` block exports the following:

* `total_size_in_gb` - Total Capacity in GB.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata VM Cluster.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
