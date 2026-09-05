---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_grid_infrastructure_versions"
description: |-
  Gets the list of Oracle Grid Infrastructure versions.
---

# Data Source: azurerm_oracle_grid_infrastructure_versions

Use this data source to access the Oracle Grid Infrastructure versions available in a location.

## Example Usage

```hcl
data "azurerm_oracle_grid_infrastructure_versions" "example" {
  location = "West Europe"
  shape    = "Exadata.X9M"
  zone     = "2"
}

output "versions" {
  value = data.azurerm_oracle_grid_infrastructure_versions.example.versions
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region to query for Oracle Grid Infrastructure versions.

---

* `shape` - (Optional) The system shape used to filter the available Oracle Grid Infrastructure versions. Possible values are `ExaDbXS`, `Exadata.X9M`, and `Exadata.X11M`.

* `zone` - (Optional) The Azure availability zone used to filter the available Oracle Grid Infrastructure versions.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Oracle Grid Infrastructure versions data source.

* `versions` - A `versions` block as defined below.

---

A `versions` block exports the following:

* `id` - The ID of the Oracle Grid Infrastructure version.

* `name` - The name of the Oracle Grid Infrastructure version.

* `version` - The Oracle Grid Infrastructure version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle Grid Infrastructure versions.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
