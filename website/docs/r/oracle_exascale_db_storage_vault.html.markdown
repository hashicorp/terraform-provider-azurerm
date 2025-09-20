---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_exascale_db_storage_vault"
description: |-
  Manages a Exadata Database Storage Vault.
---

# azurerm_oracle_exascale_db_storage_vault

Manages a Exadata Database Storage Vault.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_oracle_exascale_db_storage_vault" "example" {
  name                              = "example-exascale-db-storage-vault"
  resource_group_name               = azurerm_resource_group.example.name
  location                          = azurerm_resource_group.example.location
  zones                             = ["1"]
  display_name                      = "example-exascale-db-storage-vault"
  description                       = "description"
  additional_flash_cache_in_percent = 100
  high_capacity_database_storage {
    total_size_in_gb = 300
  }
  time_zone = "UTC"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Exadata Database Storage Vault should exist. Changing this forces a new Exadata Database Storage Vault to be created.

* `name` - (Required) The name which should be used for this Exadata Database Storage Vault. Changing this forces a new Exadata Database Storage Vault to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Exadata Database Storage Vault should exist. Changing this forces a new Exadata Database Storage Vault to be created.

* `additional_flash_cache_in_percent` - (Required) The size of additional Flash Cache in percentage of High Capacity database storage. Changing this forces a new Exadata Database Storage Vault to be created.

* `description` - (Required) Exadata Database Storage Vault description. Changing this forces a new Exadata Database Storage Vault to be created.

* `display_name` - (Required) The user-friendly name for the Exadata Database Storage Vault resource. The name does not need to be unique. Changing this forces a new Exadata Database Storage Vault to be created.

* `high_capacity_database_storage` - (Required) A `high_capacity_database_storage` block as defined below. Changing this forces a new Exadata Database Storage Vault to be created.

* `zones` - (Required) Exadata Database Storage Vault zones. Changing this forces a new Exadata Database Storage Vault to be created.

* `time_zone` - (Optional) The time zone that you want to use for the Exadata Database Storage Vault. Changing this forces a new Exadata Database Storage Vault to be created. For details, see [Time Zones](https://docs.oracle.com/en/cloud/paas/base-database/time-zone/).

* `tags` - (Optional) A mapping of tags which should be assigned to the Exadata Database Storage Vault.

---

A `high_capacity_database_storage` block supports the following:

* `total_size_in_gb` - (Required) Total Capacity. Changing this forces a new Exadata Database Storage Vault to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Exadata Database Storage Vault.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Exadata Database Storage Vault.
* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata Database Storage Vault.
* `update` - (Defaults to 30 minutes) Used when updating the Exadata Database Storage Vault.
* `delete` - (Defaults to 30 minutes) Used when deleting the Exadata Database Storage Vault.

## Import

Exadata Database Storage Vaults can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_exascale_db_storage_vault.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/exascaleDbStorageVaults/exascaleDbStorageVaults1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-03-01
