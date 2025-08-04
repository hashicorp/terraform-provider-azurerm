---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_clone_from_database"
description: |-
  Manages an Oracle Autonomous Database Clone From Database Instance.
---

# azurerm_oracle_autonomous_database_clone_from_database

Manages an Oracle Autonomous Database Clone from database instance.

## Example Usage

```hcl

resource "azurerm_oracle_autonomous_database_clone_from_database" "database_clone" {
  name                = "Example%[2]dclone"
  resource_group_name = azurerm_oracle_autonomous_database.test.resource_group_name
  location            = azurerm_oracle_autonomous_database.test.location

  source_autonomous_database_id = azurerm_oracle_autonomous_database.test.id
  clone_type                    = "Full"
  source                        = "Database"

  admin_password                   = "BEstrO0ng_#11"
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2.0
  compute_model                    = "ECPU"
  data_storage_size_in_tb          = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "Example%[2]dclone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = true
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_oracle_autonomous_database.test.subnet_id
  virtual_network_id               = azurerm_oracle_autonomous_database.test.virtual_network_id
}
```

## Arguments Reference

The following arguments are supported:

* `admin_password` - (Required) The password for the SYS, SYSTEM, and PDB Admin users. The password must be at least 12 characters long, and contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing. Changing this forces a new Autonomous Database Clone to be created.

* `auto_scaling_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database storage.

* `backup_retention_period_in_days` - (Required) Retention period, in days, for backups. Possible values are between `1` and `60`. Changing this forces a new Autonomous Database Clone to be created.

* `character_set` - (Required) The character set for the autonomous database. Changing this forces a new Autonomous Database Clone to be created.

* `clone_type` - (Required) The type of clone to create. Possible values are `Full` and `Metadata`. Changing this forces a new Autonomous Database Clone to be created.

* `compute_count` - (Required) The compute amount (CPUs) available to the database. Possible values are between `2.0` and `512.0`.

* `compute_model` - (Required) The compute model of the Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `customer_contacts` - (Optional) Customer contact email addresses. Changing this forces a new Autonomous Database Clone to be created.

* `database_version` - (Required) A valid Oracle Database version for Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `data_storage_size_in_tb` - (Required) The maximum storage that can be allocated for the database, in terabytes. Possible values are between `1` and `384`.

* `db_workload` - (Required) The Autonomous Database workload type. Possible values are `OLTP` and `DW`, `APEX`, `AJD`. Changing this forces a new Autonomous Database Clone to be created.
  * OLTP - indicates an Autonomous Transaction Processing database
  * DW - indicates an Autonomous Data Warehouse database
  * AJD - indicates an Autonomous JSON Database
  * APEX - indicates an Autonomous Database with the Oracle APEX Application Development workload type.

~> **Note:** To clone the database with a different `db_workload` type, please refer to the documentation [here](https://docs.public.oneportal.content.oci.oraclecloud.com/en-us/iaas/autonomous-database-serverless/doc/autonomous-clone-cross-workload-type.html#GUID-527A712D-FF82-498B-AB35-8A1623E36EDD) for correct configuration steps.

* `display_name` - (Required) The user-friendly name for the Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `license_model` - (Required) The Oracle license model that applies to the Oracle Autonomous Database. Possible values are `LicenseIncluded` and `BringYourOwnLicense`. Changing this forces a new Autonomous Database Clone to be created.

* `location` - (Required) The Azure Region where the Autonomous Database Clone should exist. Changing this forces a new Autonomous Database Clone to be created.

* `mtls_connection_required` - (Required) Specifies if the Autonomous Database requires mTLS connections. Changing this forces a new Autonomous Database Clone to be created.

* `name` - (Required) The name of this Autonomous Database Clone. Changing this forces a new Autonomous Database Clone to be created.

* `national_character_set` - (Required) The national character set for the autonomous database. Changing this forces a new Autonomous Database Clone to be created.

* `refreshable_model` - (Optional) The refreshable model for the clone. Possible values are `Automatic` and `Manual`.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database Clone should exist. Changing this forces a new Autonomous Database Clone to be created.

* `source_autonomouse_database_id` - (Required) The ID of the source Autonomous Database to clone from. Changing this forces a new Autonomous Database Clone to be created.

* `source` - (Required) The source of the clone. Possible values are  `Database` . Changing this forces a new Autonomous Database Clone to be created.

* `subnet_id` - (Optional) The ID of the subnet the resource is associated with. Changing this forces a new Autonomous Database Clone to be created.

* `time_until_reconnect` - (Optional) The time until reconnect clone is enabled. Must be in RFC3339 format.

* `virtual_network_id` - (Optional) The ID of the Virtual Network this Autonomous Database Clone should be created in. Changing this forces a new Autonomous Database Clone to be created.

* `tags` - (Optional) A mapping of tags to assign to the Autonomous Database Clone.

* `allowed_ips` - (Optional) (Optional) Defines the network access type for the Autonomous Database. If the property is explicitly set to an empty list, it allows secure public access to the database from any IP address. If specific ACL (Access Control List) values are provided, access will be restricted to only the specified IP addresses.

~> **Note:** `allowed_ips`  cannot be updated after provisioning the resource with an empty list (i.e., a publicly accessible Autonomous Database)
size: the maximum number of Ips provided shouldn't exceed 1024. At this time we only support IpV4.
---

## Attributes Reference

In addition to the Arguments listed above—the following Attributes are exported:

* `id` - The ID of the Autonomous Database Clone.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Autonomous Database Clone.
* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database Clone.
* `update` - (Defaults to 30 minutes) Used when updating the Autonomous Database Clone.
* `delete` - (Defaults to 30 minutes) Used when deleting the Autonomous Database Clone.

## Import

Autonomous Database Clones can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_autonomous_database_clone_from_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/autonomousDatabases/adb1
```

## API Providers

* `Oracle.Database`: 2025-03-01
