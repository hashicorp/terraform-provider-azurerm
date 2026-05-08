---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
description: |-
  Manages a Cross Region Disaster Recovery Autonomous Database.
---

# azurerm_oracle_autonomous_database_cross_region_disaster_recovery

Manages a Cross Region Disaster Recovery Autonomous Database.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}

resource "azurerm_virtual_network" "example_primary_vnet" {
  name                = "example-primary-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example_primary_subnet" {
  name                 = "example-primary-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_primary_vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_virtual_network" "example_dr_vnet" {
  name                = "example-dr-vnet"
  location            = "westus"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "example_dr_subnet" {
  name                 = "example-dr-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_dr_vnet.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_oracle_autonomous_database" "example_primary" {
  name                             = "example-primary"
  display_name                     = "example-primary"
  resource_group_name              = azurerm_resource_group.example.name
  location                         = azurerm_resource_group.example.location
  subnet_id                        = azurerm_subnet.example_primary_subnet.id
  virtual_network_id               = azurerm_virtual_network.example_primary_vnet.id
  admin_password                   = "SomeP@ssw0rd123"
  auto_scaling_enabled             = false
  auto_scaling_for_storage_enabled = false
  backup_retention_period_in_days  = 7
  character_set                    = "AL32UTF8"
  compute_count                    = 2
  compute_model                    = "ECPU"
  customer_contacts                = ["test@example.com"]
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "DW"
  license_model                    = "LicenseIncluded"
  mtls_connection_required         = false
  national_character_set           = "AL16UTF16"
}

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "example" {
  name                                = "example-dr"
  resource_group_name                 = azurerm_resource_group.example.name
  location                            = "westus"
  display_name                        = "example-dr"
  source_autonomous_database_id       = azurerm_oracle_autonomous_database.example_primary.id
  subnet_id                           = azurerm_subnet.example_dr_subnet.id
  virtual_network_id                  = azurerm_virtual_network.example_dr_vnet.id
  replicate_automatic_backups_enabled = true

  tags = {
    environment = "production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name for this Cross Region Disaster Recovery Autonomous Database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the resource should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Cross Region Disaster Recovery Autonomous Database should exist. Changing this forces a new resource to be created.

* `display_name` - (Required) The user-friendly name for the Autonomous Database. Changing this forces a new resource to be created.

* `source_autonomous_database_id` - (Required) The ID of the source (primary) Autonomous Database. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet in the target region. Changing this forces a new resource to be created.

* `virtual_network_id` - (Required) The ID of the virtual network in the target region. Changing this forces a new resource to be created.

* `replicate_automatic_backups_enabled` - (Optional) If true, 7 days of backups are replicated across regions.

* `tags` - (Optional) A mapping of tags assigned to the resource. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above the following Attributes are exported:

* `id` - The ID of the Cross Region Disaster Recovery Autonomous Database.

* `auto_scaling_enabled` - Whether auto-scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Whether auto-scaling is enabled for the Autonomous Database storage.

* `backup_retention_period_in_days` - The backup retention period in days.

* `character_set` - The character set for the Autonomous Database.

* `compute_count` - The compute amount available to the database.

* `compute_model` - The compute model of the Autonomous Database.

* `customer_contacts` - A list of customer contact email addresses.

* `data_storage_size_in_tb` - The maximum storage that can be allocated for the database in terabytes.

* `database_version` - The Oracle Database version for the Autonomous Database.

* `database_workload` - The Autonomous Database workload type.

* `license_model` - The Oracle license model that applies to the Autonomous Database.

* `mtls_connection_required` - Specifies if the Autonomous Database requires mTLS connections.

* `national_character_set` - The national character set for the Autonomous Database.

* `remote_disaster_recovery_type` - Type of recovery.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Cross Region Disaster Recovery Autonomous Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cross Region Disaster Recovery Autonomous Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cross Region Disaster Recovery Autonomous Database.

## Import

Cross Region Disaster Recovery Autonomous Database can be imported using the `resource id`.

```shell
terraform import azurerm_oracle_autonomous_database_cross_region_disaster_recovery.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/autonomousDatabases/autonomousDatabases1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
