---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_database_system_versions"
description: |-
  Gets information about existing Oracle Database Systems Versions.
---

# Data Source: azurerm_oracle_database_system_versions

Use this data source to access information about existing Oracle Database Systems Versions.

## Example Usage

```hcl
data "azurerm_oracle_db_system_versions" "example" {
  location = "eastus"
}

output "id" {
  value = data.azurerm_oracle_db_system_versions.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region to query for the Oracle Database Systems Versions.

---

* `database_software_image_supported` - (Optional) Whether to filter the results to the set of Oracle Database versions that are supported for the database software images.

* `database_system_shape` - (Optional) If provided, filters the results to the set of database versions which are supported for the given shape. Possible value is `VM.Standard.x86`.

* `shape_family` - (Optional) If provided, filters the results to the set of database versions which are supported for the given shape family. Possible values are  `EXADATA`, `EXADB_XS`, `SINGLENODE`, `VIRTUALMACHINE` 

* `storage_management` - (Optional) The DB system storage management option. Used to list database versions available for that storage manager. Possible value is `LVM`.

* `upgrade_support_enabled` - (Optional) Whether to filter the results to the set of database versions which are supported for Upgrade.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the azurerm Oracle Database Systems Versions.

* `versions` - A `versions` block as defined below.

---

A `versions` block exports the following:

* `name` - The name of Oracle Database version. 

* `latest_version` - Indicates if this version of the Oracle Database software is the latest version for a release.

* `pluggable_database_supported` - Indicates if this version of the Oracle Database software supports pluggable databases.

* `version` - The value of the Oracle Database version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the azurerm Oracle Database Systems Versions.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
