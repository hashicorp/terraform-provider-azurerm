---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_db_system"
description: |-
  Manages a DB System.
---

# azurerm_oracle_db_system

Manages a DB System.

## Example Usage

```hcl
resource "azurerm_oracle_db_system" "example" {
  name                            = "example-db-system"
  resource_group_name             = "example-resources"
  location                        = "West Europe"
  zones                           = ["1"]
  cpu_core_count                  = 4
  source                  		    = "None"
  database_edition      		      = "StandardEdition"
  db_version				              = "19.27.0.0"
  hostname                        = "hostname"
  network_anchor_id               = "/subscriptions/7a481e15-0e3c-420f-8dc7-4d183bd8d0f8/resourceGroups/wen_rg_eastus2euap/providers/Oracle.Database/networkAnchors/NetworkAnchorRegion1"
  resource_anchor_id              = "/subscriptions/7a481e15-0e3c-420f-8dc7-4d183bd8d0f8/resourceGroups/wen_rg_eastus2euap/providers/Oracle.Database/resourceAnchors/ResourceAnchorRegion1"
  shape                        	  = "VM.Standard.x86"
  ssh_public_keys                 = [file("~/.ssh/id_rsa.pub")]
}
```

## Arguments Reference

The following arguments are supported:

* `database_edition` - (Required) The Oracle Database Edition that applies to all the databases on the DB system. Exadata DB systems and 2-node RAC DB systems require EnterpriseEditionExtremePerformance.
Changing this forces a new DB System to be created.

* `db_version` - (Required) A valid Oracle Database version. For a list of supported versions, use the ListDbVersions operation. Changing this forces a new DB System to be created.

* `hostname` - (Required) The hostname for the DB system. Changing this forces a new DB system to be created.

* `location` - (Required) The Azure Region where the DB System should exist. Changing this forces a new DB system to be created.

* `name` - (Required) The name which should be used for this DB system. Changing this forces a new DB system to be created.

* `network_anchor_id` - (Required) The ID of the Azure Network Anchor. Changing this forces a new DB system to be created.

* `resource_anchor_id` - (Required) The ID of the Azure Resource Anchor. Changing this forces a new DB system to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the DB system should exist. Changing this forces a new DB system to be created.

* `shape` - (Required) The shape of the DB system. The shape determines resources to allocate to the DB system. For virtual machine shapes, the number of CPU cores and memory. For bare metal and Exadata shapes, the number of CPU cores, storage, and memory. Changing this forces a new DB system to be created.

* `source` - (Required) The source of the database: Use NONE for creating a new database. The default is `None`. Changing this forces a new Db System to be created.

* `ssh_public_keys` - (Required) Specifies a list of ssh public keys. The public key portion of one or more key pairs used for SSH access to the DB system. Changing this forces a new DB system to be created.

* `zones` - (Required) Specifies a list of DB System zones. Changing this forces a new DB System to be created.

---

* `admin_password` - (Optional) A strong password for SYS, SYSTEM, and PDB Admin. The password must be at least nine characters and contain at least two uppercase, two lowercase, two numbers, and two special characters. The special characters must be _, #, or -. Changing this forces a new DB system to be created.

* `cluster_name` - (Optional) The cluster name for Exadata and 2-node RAC virtual machine DB systems. The cluster name must begin with an alphabetic character, and may contain hyphens (-). Underscores (_) are not permitted. The cluster name can be no longer than 11 characters and is not case sensitive. Changing this forces a new DB system to be created.

* `compute_count` - (Optional) The number of compute servers for the DB system. Changing this forces a new DB system to be created.

* `compute_model` - (Optional) The compute model for Base Database Service. This is required if using the `computeCount` parameter. If using `cpuCoreCount` then it is an error to specify `computeModel` to a non-null value. The ECPU compute model is the recommended model, and the OCPU compute model is legacy. Changing this forces a new DB system to be created.

* `db_system_options` - (Optional) One or more `db_system_options` blocks as defined below.

* `disk_redundancy` - (Optional) The type of redundancy configured for the DB system. NORMAL is 2-way redundancy. HIGH is 3-way redundancy. Changing this forces a new DB system to be created.

* `display_name` - (Optional) The user-friendly name for the DB system. Changing this forces a new DB system to be created. The name does not need to be unique.

* `domain` - (Optional) The domain name for the DB system. Changing this forces a new DB system to be created.

* `initial_data_storage_size_in_gb` - (Optional) "Size in GB of the initial data volume that will be created and attached to a virtual machine DB system. You can scale up storage after provisioning, as needed. Note that the total storage size attached will be more than the amount you specify to allow for REDO/RECO space and software volume. Changing this forces a new DB system to be created.

* `license_model` - (Optional) The Oracle license model that applies to all the databases on the DB system. The default is LicenseIncluded. Changing this forces a new DB system to be created.

* `node_count` - (Optional) The number of nodes in the DB system. For RAC DB systems, the value is greater than 1. Changing this forces a new DB system to be created.

* `pluggable_database_name` - (Optional) The name of the pluggable database. The name must begin with an alphabetic character and can contain a maximum of thirty alphanumeric characters. Special characters are not permitted. Pluggable database name should not be same as database name. Changing this forces a new DB system to be created.

* `storage_volume_performance_mode` - (Optional) The block storage volume performance level. Valid values are Balanced and HighPerformance. See [Block Volume Performance](/Content/Block/Concepts/blockvolumeperformance.htm) for more information. Changing this forces a new DB system to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the DB System.

* `time_zone` - (Optional) The time zone of the DB system, e.g., UTC, to set the timeZone as UTC. Changing this forces a new DB System to be created.

---

A `db_system_options` block supports the following:

* `storage_management` - (Optional) The storage option used in DB system. ASM - Automatic storage management, LVM - Logical Volume management.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the DB System.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the DB System.
* `read` - (Defaults to 5 minutes) Used when retrieving the DB System.
* `update` - (Defaults to 30 minutes) Used when updating the DB System.
* `delete` - (Defaults to 30 minutes) Used when deleting the DB System.

## Import

DB Systems can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_db_system.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/dbsystems/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
