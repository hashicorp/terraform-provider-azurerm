---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
description: |-
  Manages an Cross Region Disaster Recovery Autonomous Database.
---

# azurerm_oracle_autonomous_database_cross_region_disaster_recovery

Manages an Cross Region Disaster Recovery Autonomous Database.
Cross Region Disaster Recovery Autonomous Database is an Autonomous Database with specific Cross Region Disaster Recovery role. It must be an exact copy of Autonomous Database for which you want to create a Disaster Recovery instance. Cross Region Disaster Recovery Autonomous Database must reside in a region that is different from region of main Autonomous Database. You must create a separate virtual network and subnet in this second region for Cross Region Disaster Recovery Autonomous Database to be able to communicate with it's original database. All parameters, except "name" and "display_name" must be exactly the same as for original database or creation of Cross Region Disaster Recovery Autonomous Database will fail.

## Example Usage

```hcl
resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "adbs_secondary_crdr" {
  name                             = "example"
  display_name                     = "example_display_name"
  location                         = "westus"
  database_type                    = "CrossRegionDisasterRecovery"
  source                           = "CrossRegionDisasterRecovery"
  source_id      				   = "/subscriptions/xxxxx-xxxxxxxxxxxxx-xxxxxxxx/resourceGroups/resource_group_name/providers/Oracle.Database/autonomousDatabases/autonomous_database_name"
  source_ocid                      = "ocid1.autonomousdatabase.xx.xx.xxxxxxxxxxxxxxx"
  remote_disaster_recovery_type    = "Adg"
  replicate_automatic_backups	   = true
  subnet_id                        = "/subscriptions/xxxxx-xxxxxxxxxxxxx-xxxxxxxx/resourceGroups/resource_group_name/providers/Microsoft.Network/virtualNetworks/virtual_network_name/subnets/subnet_name",
  virtual_network_id               = "/subscriptions/xxxxx-xxxxxxxxxxxxx-xxxxxxxx/resourceGroups/resource_group_name/providers/Microsoft.Network/virtualNetworks/virtual_network_name"
  resource_group_name              = "example"
  source_location				   = "eastus"
  license_model                    = "example"
  backup_retention_period_in_days  = 42
  auto_scaling_enabled             = true
  auto_scaling_for_storage_enabled = true
  mtls_connection_required         = false
  compute_model                    = "example"
  compute_count                    = 1.23456
  db_workload                      = "example"
  db_version                       = "example"
  admin_password                   = "example"
  character_set                    = "example"
  national_character_set           = "example"
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cross Region Disaster Recovery Autonomous Database. Changing this forces a new Cross Region Disaster Recovery Autonomous Database to be created.

* `display_name` - (Required) The user-friendly name for the Autonomous Database. The name does not have to be unique.

* `location` - (Required) The Azure Region where the Cross Region Disaster Recovery Autonomous Database should exist. Azure Region must be different from region where your original Autonomous Database resides. Changing this forces a new Cross Region Disaster Recovery Autonomous Database to be created.

* `database_type` - (Required) The type of Autonomous Database. To create Cross Region Disaster Recovery Autonomous Database set this value to "CrossRegionDisasterRecovery".

* `source` - (Required) Source of Autonomous Database. To create Cross Region Disaster Recovery Autonomous Database set this value to "CrossRegionDisasterRecovery".

* `remote_disaster_recovery_type` - (Required) Type of recovery. Value can be either `Adg` (Autonomous Data Guard) or `BackupBased`.

* `replicate_automatic_backups` - (Required) If true, 7 days worth of backups are replicated across regions for Cross-Region ADB or Backup-Based Disaster Recovery between Primary and Standby. If false, the backups taken on the Primary are not replicated to the Standby database.

* `subnet_id` - (Required) The Immutable Azure Resource ID of the subnet the resource is associated with. Must be associated with virtual network for this Cross Region Disaster Recovery Autonomous Database.

* `virtual_network_id` - (Required) The Immutable Azure Resource ID of the virtual network the resource is associated with. Must be located in a region that is different from the region of Autonomous Database for which the Cross Region Disaster Recovery Autonomous Database is created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cross Region Disaster Recovery Autonomous Database should exist.

* `source_location` - The Azure Region where source autonomous database for which cross region disaster recovery autonomous database is located.

* `license_model` - (Required) The Oracle license model that applies to the Oracle Autonomous Database. Bring your own license (BYOL) allows you to apply your current on-premises Oracle software licenses to equivalent, highly automated Oracle services in the cloud. License Included allows you to subscribe to new Oracle Database software licenses and the Oracle Database service. When provisioning an [Autonomous Database Serverless] (https://docs.oracle.com/en/cloud/paas/autonomous-database/index.html) database, if a value is not specified, the system defaults the value to `BRING_YOUR_OWN_LICENSE`.

* `backup_retention_period_in_days` - (Required) Retention period, in days, for backups.

* `auto_scaling_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database storage.

* `mtls_connection_required` - (Required) Specifies if the Autonomous Database requires mTLS connections.

* `compute_model` - (Required) The compute model of the Autonomous Database. This is required if using the `computeCount` parameter. If using `cpuCoreCount` then it is an error to specify `computeModel` to a non-null value. ECPU compute model is the recommended model and OCPU compute model is legacy.

* `compute_count` - (Required) The compute amount (CPUs) available to the database. For an Autonomous Database Serverless instance, the `ECPU` compute model requires a minimum value of one. Required when using the `computeModel` parameter.

* `db_workload` - (Required) The Autonomous Database workload type. The following values are valid:
  * OLTP - indicates an Autonomous Transaction Processing database
  * DW - indicates an Autonomous Data Warehouse database
  * AJD - indicates an Autonomous JSON Database
  * APEX - indicates an Autonomous Database with the Oracle APEX Application Development workload type.

* `db_version` - (Required) A valid Oracle Database version for Autonomous Database.

* `admin_password` - (Required) The password must be between `12` and `30 `characters long, and must contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing.

* `character_set` - (Required) The character set for the autonomous database.  The default is `AL32UTF8`. Allowed values are:  `AL32UTF8`, `AR8ADOS710`, `AR8ADOS720`, `AR8APTEC715`, `AR8ARABICMACS`, `AR8ASMO8X`, `AR8ISO8859P6`, `AR8MSWIN1256`, `AR8MUSSAD768`, `AR8NAFITHA711`, `AR8NAFITHA721`, `AR8SAKHR706`, `AR8SAKHR707`, `AZ8ISO8859P9E`, `BG8MSWIN`, `BG8PC437S`, `BLT8CP921`, `BLT8ISO8859P13`, `BLT8MSWIN1257`, `BLT8PC775`, `BN8BSCII`, `CDN8PC863`, `CEL8ISO8859P14`, `CL8ISO8859P5`, `CL8ISOIR111`, `CL8KOI8R`, `CL8KOI8U`, `CL8MACCYRILLICS`, `CL8MSWIN1251`, `EE8ISO8859P2`, `EE8MACCES`, `EE8MACCROATIANS`, `EE8MSWIN1250`, `EE8PC852`, `EL8DEC`, `EL8ISO8859P7`, `EL8MACGREEKS`, `EL8MSWIN1253`, `EL8PC437S`, `EL8PC851`, `EL8PC869`, `ET8MSWIN923`, `HU8ABMOD`, `HU8CWI2`, `IN8ISCII`, `IS8PC861`, `IW8ISO8859P8`, `IW8MACHEBREWS`, `IW8MSWIN1255`, `IW8PC1507`, `JA16EUC`, `JA16EUCTILDE`, `JA16SJIS`, `JA16SJISTILDE`, `JA16VMS`, `KO16KSC5601`, `KO16KSCCS`, `KO16MSWIN949`, `LA8ISO6937`, `LA8PASSPORT`, `LT8MSWIN921`, `LT8PC772`, `LT8PC774`, `LV8PC1117`, `LV8PC8LR`, `LV8RST104090`, `N8PC865`, `NE8ISO8859P10`, `NEE8ISO8859P4`, `RU8BESTA`, `RU8PC855`, `RU8PC866`, `SE8ISO8859P3`, `TH8MACTHAIS`, `TH8TISASCII`, `TR8DEC`, `TR8MACTURKISHS`, `TR8MSWIN1254`, `TR8PC857`, `US7ASCII`, `US8PC437`, `UTF8`, `VN8MSWIN1258`, `VN8VN3`, `WE8DEC`, `WE8DG`, `WE8ISO8859P1`, `WE8ISO8859P15`, `WE8ISO8859P9`, `WE8MACROMAN8S`, `WE8MSWIN1252`, `WE8NCR4970`, `WE8NEXTSTEP`, `WE8PC850`, `WE8PC858`, `WE8PC860`, `WE8ROMAN8`, `ZHS16CGB231280`, `ZHS16GBK`, `ZHT16BIG5`, `ZHT16CCDC`, `ZHT16DBT`, `ZHT16HKSCS`, `ZHT16MSWIN950`, `ZHT32EUC`, `ZHT32SOPS`, `ZHT32TRIS`.

* `national_character_set` - (Required) The national character set for the autonomous database. The default is AL16UTF16. Allowed values are: AL16UTF16 or UTF8.

## Import

Cross Region Disaster Recovery Autonomous Database cannot be imported. `admin_password` parameter is required, but this parameter cannot be imported, therefore the entire import procedure cannot be executed.

## Update

Cross Region Disaster Recovery Autonomous Database cannot be updated as is must have exactly same parameters as peered original Autonomous Database. Updating original Autonomous Database will trigger an update for peered Cross Region Disaster Recovery Autonomous Database.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database`: 2025-03-01
