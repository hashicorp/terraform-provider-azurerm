---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
description: |-
  Manages a Cross Region Disaster Recovery Autonomous Database.
---

# azurerm_oracle_autonomous_database_cross_region_disaster_recovery

Manages Cross Region Disaster Recovery Autonomous Database.
Cross Region Disaster Recovery Autonomous Database is an Autonomous Database with a specific Cross-Region Disaster Recovery role. It must be an exact copy of Autonomous Database for which you want to create a Disaster Recovery instance. Cross Region Disaster Recovery Autonomous Database must reside in a region that is different from region of main Autonomous Database. You must create a separate virtual network and subnet in this second region for Cross Region Disaster Recovery Autonomous Database to be able to communicate with its original database.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}

resource "azurerm_virtual_network" "example_vnet" {
  name                = "exampleVnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example_subnet" {
  name                 = "exampleSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}


resource "azurerm_oracle_autonomous_database" "primary" {
  name                             = "examplePrimary"
  display_name                     = "ExamplePrimary"
  resource_group_name              = azurerm_resource_group.example.name
  location                         = "eastus"
  compute_model                    = "ECPU"
  compute_count                    = 2
  license_model                    = "LicenseIncluded"
  backup_retention_period_in_days  = 7
  auto_scaling_enabled             = true
  auto_scaling_for_storage_enabled = true
  mtls_connection_required         = false
  data_storage_size_in_tbs         = 1
  db_workload                      = "DW"
  admin_password                   = "SomeP@ssw0rd123"
  db_version                       = "19c"
  character_set                    = "AL32UTF8"
  national_character_set           = "AL16UTF16"
  subnet_id                        = azurerm_subnet.example_subnet.id
  virtual_network_id               = azurerm_virtual_network.example_vnet.id
  customer_contacts                = ["test@example.com"]
}

resource "azurerm_virtual_network" "dr_vnet" {
  name                = "drVnet"
  location            = "westus"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "dr_subnet" {
  name                 = "drSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.dr_vnet.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "dr_example" {
  name                                = "exampledr"
  display_name                        = "ExampleDR"
  location                            = "westus"
  resource_group_name                 = azurerm_resource_group.example.name
  source_autonomous_database_id       = azurerm_oracle_autonomous_database.primary.id
  subnet_id                           = azurerm_subnet.dr_subnet.id
  replicate_automatic_backups_enabled = true

  tags = {
    Environment = "production"
    Purpose     = "disaster-recovery"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the resource should exist. Changing this forces a new resource to be created.

* `name` - (Required) The name for this Cross Region Disaster Recovery Autonomous Database. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Cross Region Disaster Recovery Autonomous Database should exist. Must differ from the primary region. Changing this forces a new resource to be created.

* `display_name` - (Required) The user-friendly name for the Autonomous Database. Changing this forces a new resource to be created.

* `replicate_automatic_backups_enabled` - (Required) If true, 7 days of backups are replicated across regions. Changing this forces a new resource to be created.

* `source_autonomous_database_id` - (Required) The ID of the source (primary) Autonomous Database. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet in the target region. Changing this forces a new resource to be created.

* `tags` - (Optional) Map of tags to assign to the resource. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Autonomous Database.

* `auto_scaling_enabled` - Whether auto-scaling is enabled for the Autonomous Database CPU core count.

* `auto_scaling_for_storage_enabled` - Whether auto-scaling is enabled for the Autonomous Database storage.

* `backup_retention_period_in_days` - The backup retention period in days.

* `character_set` - The character set for the autonomous database.

* `compute_count` - The compute amount (CPUs) available to the database.

* `compute_model` - The compute model of the Autonomous Database.

* `customer_contacts` - A list of Customer's contact email addresses.

* `data_storage_size_in_tb` - The maximum storage that can be allocated for the database in terabytes.

* `database_version` - The Oracle Database version for Autonomous Database.

* `database_workload` - The Autonomous Database workload type.

* `license_model` - The Oracle license model that applied to the Oracle Autonomous Database.

* `mtls_connection_required` - Specifies if the Autonomous Database requires mTLS connections.

* `national_character_set` - The national character set for the autonomous database.

* `remote_disaster_recovery_type` - Type of recovery.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Cross Region Disaster Recovery Autonomous Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cross Region Disaster Recovery Autonomous Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cross Region Disaster Recovery Autonomous Database.

## Import

Cross Region Disaster Recovery Autonomous Database can be imported using the `resource id`.

```shell
terraform import azurerm_oracle_autonomous_database_cross_region_disaster_recovery.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/autonomousDatabases/database1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
