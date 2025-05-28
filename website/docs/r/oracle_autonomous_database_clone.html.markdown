---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_clone"
description: |-
  Manages an Oracle Autonomous Database Clone.
---

# azurerm_oracle_autonomous_database_clone

Manages an Oracle Autonomous Database Clone.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "Oracle.Database.networkAttachments"

    service_delegation {
      name    = "Oracle.Database/networkAttachments"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_oracle_autonomous_database" "source" {
  name                = "example-source-db"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 7
  character_set                     = "AL32UTF8"
  compute_count                     = 2.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "Example Source Database"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.example.id
  virtual_network_id               = azurerm_virtual_network.example.id
}

resource "azurerm_oracle_autonomous_database_clone" "example" {
  name                = "example-clone-db"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  # Clone-specific configuration
  source_id  = azurerm_oracle_autonomous_database.source.id
  clone_type = "Full"
  source     = "Database"

  # Database configuration
  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 7
  character_set                     = "AL32UTF8"
  compute_count                     = 2.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "Example Clone Database"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.example.id
  virtual_network_id               = azurerm_virtual_network.example.id

  tags = {
    Environment = "Development"
    Purpose     = "Clone"
  }
}
```

## Example Usage - Refreshable Clone

```hcl
resource "azurerm_oracle_autonomous_database_clone" "refreshable" {
  name                = "example-refreshable-clone"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  # Clone-specific configuration
  source_id            = azurerm_oracle_autonomous_database.source.id
  clone_type           = "Full"
  source               = "Database"
  is_refreshable_clone = true
  refreshable_model    = "Manual"

  # Database configuration
  admin_password                    = "BEstrO0ng_#11"
  backup_retention_period_in_days   = 7
  character_set                     = "AL32UTF8"
  compute_count                     = 2.0
  compute_model                     = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "Example Refreshable Clone"
  license_model                    = "LicenseIncluded"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.example.id
  virtual_network_id               = azurerm_virtual_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Autonomous Database Clone. Changing this forces a new Autonomous Database Clone to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Autonomous Database Clone should exist. Changing this forces a new Autonomous Database Clone to be created.

* `location` - (Required) The Azure Region where the Autonomous Database Clone should exist. Changing this forces a new Autonomous Database Clone to be created.

### Clone-Specific Arguments

* `source_id` - (Required) The ID of the source Autonomous Database to clone from. Changing this forces a new Autonomous Database Clone to be created.

* `clone_type` - (Required) The type of clone to create. Possible values are `Full` and `Metadata`. Changing this forces a new Autonomous Database Clone to be created.

* `source` - (Required) The source of the clone. Possible values are `None`, `Database`, `BackupFromId`, `BackupFromTimestamp`, `CloneToRefreshable`, `CrossRegionDataguard`, and `CrossRegionDisasterRecovery`. Changing this forces a new Autonomous Database Clone to be created.

* `is_refreshable_clone` - (Optional) Indicates whether the clone is a refreshable clone. Changing this forces a new Autonomous Database Clone to be created.

* `refreshable_model` - (Optional) The refreshable model for the clone. Possible values are `Automatic` and `Manual`. Changing this forces a new Autonomous Database Clone to be created.

* `is_reconnect_clone_enabled` - (Optional) Indicates whether reconnect clone is enabled. Changing this forces a new Autonomous Database Clone to be created.

* `time_until_reconnect_clone_enabled` - (Optional) The time until reconnect clone is enabled. Must be in RFC3339 format. Changing this forces a new Autonomous Database Clone to be created.

### Database Configuration Arguments

* `admin_password` - (Required) The password for the SYS, SYSTEM, and PDB Admin users. The password must be at least 12 characters long, and contain at least 1 uppercase, 1 lowercase, and 1 numeric character. It cannot contain the double quote symbol (") or the username "admin", regardless of casing. Changing this forces a new Autonomous Database Clone to be created.

* `backup_retention_period_in_days` - (Required) Retention period, in days, for backups. Possible values are between `1` and `60`. Changing this forces a new Autonomous Database Clone to be created.

* `character_set` - (Required) The character set for the autonomous database. Changing this forces a new Autonomous Database Clone to be created.

* `compute_count` - (Required) The compute amount (CPUs) available to the database. Possible values are between `2.0` and `512.0`.

* `compute_model` - (Required) The compute model of the Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `data_storage_size_in_tbs` - (Required) The maximum storage that can be allocated for the database, in terabytes. Possible values are between `1` and `384`.

* `db_version` - (Required) A valid Oracle Database version for Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `db_workload` - (Required) The Autonomous Database workload type. Possible values are `OLTP` and `DW`. Changing this forces a new Autonomous Database Clone to be created.

* `display_name` - (Required) The user-friendly name for the Autonomous Database. Changing this forces a new Autonomous Database Clone to be created.

* `license_model` - (Required) The Oracle license model that applies to the Oracle Autonomous Database. Possible values are `LicenseIncluded` and `BringYourOwnLicense`. Changing this forces a new Autonomous Database Clone to be created.

* `auto_scaling_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - (Required) Indicates if auto scaling is enabled for the Autonomous Database storage.

* `mtls_connection_required` - (Required) Specifies if the Autonomous Database requires mTLS connections. Changing this forces a new Autonomous Database Clone to be created.

* `national_character_set` - (Required) The national character set for the autonomous database. Changing this forces a new Autonomous Database Clone to be created.

* `subnet_id` - (Required) The ID of the subnet the resource is associated with. Changing this forces a new Autonomous Database Clone to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network this Autonomous Database Clone should be created in. Changing this forces a new Autonomous Database Clone to be created.

* `customer_contacts` - (Optional) Customer contact email addresses. Changing this forces a new Autonomous Database Clone to be created.

* `tags` - (Optional) A mapping of tags to assign to the Autonomous Database Clone.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Autonomous Database Clone.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 120 minutes) Used when creating the Autonomous Database Clone.
* `read` - (Defaults to 5 minutes) Used when retrieving the Autonomous Database Clone.
* `update` - (Defaults to 30 minutes) Used when updating the Autonomous Database Clone.
* `delete` - (Defaults to 30 minutes) Used when deleting the Autonomous Database Clone.

## Import

Autonomous Database Clones can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_autonomous_database_clone.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/autonomousDatabases/adb1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database`: 2024-06-01