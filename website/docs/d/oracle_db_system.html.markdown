---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_db_system"
description: |-
  Gets information about an existing DB System.
---

# Data Source: azurerm_oracle_db_system

Use this data source to access information about an existing DB System.

## Example Usage

```hcl
data "azurerm_oracle_db_system" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_db_system.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this DB System.

* `resource_group_name` - (Required) The name of the Resource Group where the Db System exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the DB System.

* `cluster_name` - The cluster name for Exadata and 2-node RAC virtual machine DB systems. The cluster name must begin with an alphabetic character, and may contain hyphens (-). Underscores (_) are not permitted. The cluster name can be no longer than 11 characters and is not case sensitive.

* `compute_count` - The number of compute servers for the DB system.

* `compute_model` - The compute model for Base Database Service. This is required if using the `computeCount` parameter. If using `cpuCoreCount` then it is an error to specify `computeModel` to a non-null value. The ECPU compute model is the recommended model, and the OCPU compute model is legacy.

* `data_storage_size_in_gbs` - The data storage size, in gigabytes, that is currently available to the DB system. Applies only for virtual machine DB systems.

* `database_edition` - The Oracle Database Edition that applies to all the databases on the DB system. Exadata DB systems and 2-node RAC DB systems require EnterpriseEditionExtremePerformance.

* `db_system_options` - A `db_system_options` block as defined below.

* `db_version` - A valid Oracle Database version. For a list of supported versions, use the ListDbVersions operation.

* `disk_redundancy` - The type of redundancy configured for the DB system. NORMAL is 2-way redundancy. HIGH is 3-way redundancy.

* `display_name` - The user-friendly name for the DB system. The name does not have to be unique.

* `domain` - The domain name for the DB system.

* `grid_image_ocid` - The OCID of a grid infrastructure software image. This is a database software image of the type GRID_IMAGE.

* `hostname` - The hostname for the DB system.

* `license_model` - The Oracle license model that applies to all the databases on the DB system. The default is LicenseIncluded.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - The current state of the DB system.

* `listener_port` - The port number configured for the listener on the DB system.

* `location` - The Azure Region where the DB System exists.

* `memory_size_in_gbs` - Memory allocated to the DB system, in gigabytes.

* `network_anchor_id` - The ID of the Azure Network Anchor.

* `node_count` - The number of nodes in the DB system. For RAC DB systems, the value is greater than 1.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the DB system.

* `resource_anchor_id` - The ID of the Azure Resource Anchor.

* `scan_dns_name` - The FQDN of the DNS record for the SCAN IP addresses that are associated with the DB system.

* `scan_ips` - The list of Single Client Access Name (SCAN) IP addresses associated with the DB system. SCAN IP addresses are typically used for load balancing and are not assigned to any interface. Oracle Clusterware directs the requests to the appropriate nodes in the cluster. Note: For a single-node DB system, this list is empty.

* `shape` - The shape of the DB system. The shape determines resources to allocate to the DB system. For virtual machine shapes, the number of CPU cores and memory. For bare metal and Exadata shapes, the number of CPU cores, storage, and memory.

* `source` - The source of the database for creating a new database.

* `ssh_public_keys` - A `ssh_public_keys` block as defined below. The public key portion of one or more key pairs used for SSH access to the DB system.

* `storage_volume_performance_mode` - The block storage volume performance level. Valid values are Balanced and HighPerformance. See [Block Volume Performance](/Content/Block/Concepts/blockvolumeperformance.htm) for more information.

* `time_zone` - The time zone of the DB system, e.g., UTC, to set the timeZone as UTC.

* `version` - The Oracle Database version of the DB system.

* `zones` - The DB System Azure zones.

---

A `db_system_options` block supports the following:

* `storage_management` - (Optional) The storage option used in DB system. ASM - Automatic storage management, LVM - Logical Volume management.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
