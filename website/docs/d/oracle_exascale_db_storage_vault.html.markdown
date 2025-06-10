---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_exascale_db_storage_vault"
description: |-
  Gets information about an existing Exadata Database Storage Vault.
---

# Data Source: azurerm_oracle_exascale_db_storage_vault

Use this data source to access information about an existing Exadata Database Storage Vault

## Example Usage

```hcl
data "azurerm_oracle_exascale_db_storage_vault" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_exascale_db_storage_vault.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this EExadata Database Storage Vault.

* `resource_group_name` - (Required) The name of the Resource Group where the Exadata Database Storage Vault exists.

* `zones` - The Exadata Database Storage Vault Azure zones.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Exadata Database Storage Vault.

* `additional_flash_cache_in_percent` - The size of additional Flash Cache in percentage of High Capacity database storage.

* `description` - Exadata Database Storage Vault description.

* `display_name` - The user-friendly name for the Exadata Database Storage Vault. The name does not need to be unique.

* `high_capacity_database_storage` - A `high_capacity_database_storage` block as defined below.

* `lifecycle_details` - Additional information about the current lifecycle state.

* `lifecycle_state` - Exadata Database Storage Vault lifecycle state enum.

* `location` - The Azure Region where the Exadata Database Storage Vault should exist.

* `oci_url` - The URL of the resource in the OCI console.

* `ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the Exadata Database Storage Vault.

* `time_zone` - The time zone that you want to use for the Exadata Database Storage Vault.

* `vm_cluster_count` - The number of Exadata VM clusters used the Exadata Database Storage Vault.

---

A `high_capacity_database_storage` block exports the following:

* `available_size_in_gbs` - Available Capacity

* `total_size_in_gbs` - Total Capacity

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata Database Storage Vault.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database`: 2025-03-01
