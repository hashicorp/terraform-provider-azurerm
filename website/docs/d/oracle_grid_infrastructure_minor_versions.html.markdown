---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_grid_infrastructure_minor_versions"
description: |-
  Gets the list of Oracle Grid Infrastructure minor versions.
---

# Data Source: azurerm_oracle_grid_infrastructure_minor_versions

Use this data source to access the Oracle Grid Infrastructure minor versions available for a Grid Infrastructure version.

## Example Usage

```hcl
data "azurerm_oracle_grid_infrastructure_versions" "example" {
  location = "uksouth"
  shape    = "ExaDbXS"
  zone     = "1"
}

data "azurerm_oracle_grid_infrastructure_minor_versions" "example" {
  location                    = "uksouth"
  grid_infrastructure_version = data.azurerm_oracle_grid_infrastructure_versions.example.versions[0].name
  shape_family                = "EXADB_XS"
  zone                        = "1"
}

output "versions" {
  value = data.azurerm_oracle_grid_infrastructure_minor_versions.example.versions
}
```

## Arguments Reference

The following arguments are supported:

* `grid_infrastructure_version` - (Required) The name of the Oracle Grid Infrastructure version to query for minor versions.

* `location` - (Required) The Azure Region to query for Oracle Grid Infrastructure minor versions.

---

* `shape_family` - (Optional) The shape family used to filter the available Oracle Grid Infrastructure minor versions. Possible values are `EXADATA` and `EXADB_XS`.

* `zone` - (Optional) The Azure availability zone used to filter the available Oracle Grid Infrastructure minor versions.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Oracle Grid Infrastructure minor versions data source.

* `versions` - A `versions` block as defined below.

---

A `versions` block exports the following:

* `id` - The ID of the Oracle Grid Infrastructure minor version.

* `grid_image_ocid` - The Oracle Cloud Identifier (OCID) of the Grid Infrastructure image.

* `name` - The name of the Oracle Grid Infrastructure minor version.

* `version` - The Oracle Grid Infrastructure minor version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle Grid Infrastructure minor versions.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
