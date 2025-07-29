---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_gi_versions"
description: |-
  Provides the list of GI (Grid Infrastructure) Versions.
---

# Data Source: azurerm_oracle_gi_versions

This data source provides the list of GI Versions in Oracle Cloud Infrastructure Database service.

Gets a list of supported GI versions.

## Example Usage

```hcl
data "azurerm_oracle_gi_versions" "example" {
  location = "West Europe"
}

output "example" {
  value = data.azurerm_oracle_gi_versions.example
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region to query for the GI Versions in.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `versions` - A list of valid GI software versions.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the GI Versions.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database`: 2025-03-01
