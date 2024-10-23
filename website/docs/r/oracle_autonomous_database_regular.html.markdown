---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database"
description: |-
  Manages an Autonomous Database.
---

# azurerm_oracle_autonomous_database

Manages an Autonomous Database.

## Example Usage

```hcl
resource "azurerm_oracle_autonomous_database" "example" {
  name                             = "example"
  resource_group_name              = "example"
  location                         = "West Europe"
  subnet_id                        = "example"
  display_name                     = "example"
  db_workload                      = "example"
  mtls_connection_required         = false
  backup_retention_period_in_days  = 42
  compute_model                    = "example"
  data_storage_size_in_gbs         = 42
  auto_scaling_for_storage_enabled = false
  virtual_network_id               = "example"
  admin_password                   = "example"
  auto_scaling_enabled             = "example"
  character_set                    = "example"
  compute_count                    = 1.23456
  national_character_set           = "example"
  license_model                    = false
  db_version                       = "example"
}
```

## Arguments Reference

The following arguments are supported:
 
* `name` - (Required) The name which should be used for this Autonomous Database.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database should exist.

* `location` - (Required) The Azure Region where the Autonomous Database should exist. Changing this forces a new Autonomous Database to be created.

* `admin_password` - (Required) The password must be between `12` and `30 `characters long, and must contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing.

* `backup_retention_period_in_days` - (Optional) (Updatable) Retention period, in days, for backups.

* `character_set` - (Required) The character set for the autonomous database.  The default is `AL32UTF8`. Allowed values are:  `AL32UTF8`, `AR8ADOS710`, `AR8ADOS720`, `AR8APTEC715`, `AR8ARABICMACS`, `AR8ASMO8X`, `AR8ISO8859P6`, `AR8MSWIN1256`, `AR8MUSSAD768`, `AR8NAFITHA711`, `AR8NAFITHA721`, `AR8SAKHR706`, `AR8SAKHR707`, `AZ8ISO8859P9E`, `BG8MSWIN`, `BG8PC437S`, `BLT8CP921`, `BLT8ISO8859P13`, `BLT8MSWIN1257`, `BLT8PC775`, `BN8BSCII`, `CDN8PC863`, `CEL8ISO8859P14`, `CL8ISO8859P5`, `CL8ISOIR111`, `CL8KOI8R`, `CL8KOI8U`, `CL8MACCYRILLICS`, `CL8MSWIN1251`, `EE8ISO8859P2`, `EE8MACCES`, `EE8MACCROATIANS`, `EE8MSWIN1250`, `EE8PC852`, `EL8DEC`, `EL8ISO8859P7`, `EL8MACGREEKS`, `EL8MSWIN1253`, `EL8PC437S`, `EL8PC851`, `EL8PC869`, `ET8MSWIN923`, `HU8ABMOD`, `HU8CWI2`, `IN8ISCII`, `IS8PC861`, `IW8ISO8859P8`, `IW8MACHEBREWS`, `IW8MSWIN1255`, `IW8PC1507`, `JA16EUC`, `JA16EUCTILDE`, `JA16SJIS`, `JA16SJISTILDE`, `JA16VMS`, `KO16KSC5601`, `KO16KSCCS`, `KO16MSWIN949`, `LA8ISO6937`, `LA8PASSPORT`, `LT8MSWIN921`, `LT8PC772`, `LT8PC774`, `LV8PC1117`, `LV8PC8LR`, `LV8RST104090`, `N8PC865`, `NE8ISO8859P10`, `NEE8ISO8859P4`, `RU8BESTA`, `RU8PC855`, `RU8PC866`, `SE8ISO8859P3`, `TH8MACTHAIS`, `TH8TISASCII`, `TR8DEC`, `TR8MACTURKISHS`, `TR8MSWIN1254`, `TR8PC857`, `US7ASCII`, `US8PC437`, `UTF8`, `VN8MSWIN1258`, `VN8VN3`, `WE8DEC`, `WE8DG`, `WE8ISO8859P1`, `WE8ISO8859P15`, `WE8ISO8859P9`, `WE8MACROMAN8S`, `WE8MSWIN1252`, `WE8NCR4970`, `WE8NEXTSTEP`, `WE8PC850`, `WE8PC858`, `WE8PC860`, `WE8ROMAN8`, `ZHS16CGB231280`, `ZHS16GBK`, `ZHT16BIG5`, `ZHT16CCDC`, `ZHT16DBT`, `ZHT16HKSCS`, `ZHT16MSWIN950`, `ZHT32EUC`, `ZHT32SOPS`, `ZHT32TRIS`

* `compute_count` - (Required) The compute amount (CPUs) available to the database. Minimum and maximum values depend on the compute model and whether the database is an Autonomous Database Serverless instance or an Autonomous Database on Dedicated Exadata Infrastructure.  For an Autonomous Database Serverless instance, the `ECPU` compute model requires a minimum value of one, for databases in the elastic resource pool and minimum value of two, otherwise. Required when using the `computeModel` parameter. When using `cpuCoreCount` parameter, it is an error to specify computeCount to a non-null value. Providing `computeModel` and `computeCount` is the preferred method for both OCPU and ECPU.

* `compute_model` - (Required) The compute model of the Autonomous Database. This is required if using the `computeCount` parameter. If using `cpuCoreCount` then it is an error to specify `computeModel` to a non-null value. ECPU compute model is the recommended model and OCPU compute model is legacy.

* `data_storage_size_in_tbs` - (Required) The maximum storage that can be allocated for the database, in terabytes.

* `db_version` - (Required) A valid Oracle Database version for Autonomous Database.

* `db_workload` - (Required) The Autonomous Database workload type. The following values are valid:
    * OLTP - indicates an Autonomous Transaction Processing database
    * DW - indicates an Autonomous Data Warehouse database
    * AJD - indicates an Autonomous JSON Database
    * APEX - indicates an Autonomous Database with the Oracle APEX Application Development workload type.

* `display_name` - (Required) The user-friendly name for the Autonomous Database. The name does not have to be unique.

* `auto_scaling_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database CPU core count. The default value is `true`.

* `auto_scaling_for_storage_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database storage. The default value is `false`.

* `mtls_connection_required` - (Required) Specifies if the Autonomous Database requires mTLS connections.

* `license_model` - (Required) The Oracle license model that applies to the Oracle Autonomous Database. Bring your own license (BYOL) allows you to apply your current on-premises Oracle software licenses to equivalent, highly automated Oracle services in the cloud. License Included allows you to subscribe to new Oracle Database software licenses and the Oracle Database service. Note that when provisioning an [Autonomous Database on dedicated Exadata infrastructure](https://docs.oracle.com/en/cloud/paas/autonomous-database/index.html), this attribute must be null. It is already set at the Autonomous Exadata Infrastructure level. When provisioning an [Autonomous Database Serverless] (https://docs.oracle.com/en/cloud/paas/autonomous-database/index.html) database, if a value is not specified, the system defaults the value to `BRING_YOUR_OWN_LICENSE`. Bring your own license (BYOL) also allows you to select the DB edition using the optional parameter.

* `national_character_set` - (Required) The national character set for the autonomous database.  The default is AL16UTF16. Allowed values are: AL16UTF16 or UTF8.

* `subnet_id` - (Required) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the subnet the resource is associated with.

* `virtual_network_id` - (Required) The ID of the vnet associated with the cloud VM cluster.

---

* `customer_contacts` - (Optional) Specifies a list of customer contacts as email addresses.

* `tags` - (Optional) A mapping of tags which should be assigned to the Autonomous Database.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Autonomous Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Autonomous Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database.
* `update` - (Defaults to 30 minutes) Used when updating the Autonomous Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Autonomous Database.

## Import

Autonomous Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_autonomous_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/autonomousDatabases/autonomousDatabases1
```
