---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
description: |-
  Manages an Cross Region Disaster Recovery Autonomous Database.
---

# azurerm_oracle_autonomous_database_cross_region_disaster_recovery

Manages Cross Region Disaster Recovery Autonomous Database.
Cross Region Disaster Recovery Autonomous Database is an Autonomous Database with a specific Cross-Region Disaster Recovery role. It must be an exact copy of Autonomous Database for which you want to create a Disaster Recovery instance. Cross Region Disaster Recovery Autonomous Database must reside in a region that is different from region of main Autonomous Database. You must create a separate virtual network and subnet in this second region for Cross Region Disaster Recovery Autonomous Database to be able to communicate with it's original database.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}

resource "azurerm_virtual_network" "example_vnet" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example_subnet" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}


resource "azurerm_oracle_autonomous_database" "primary" {
  name                             = "example-primary"
  display_name                     = "Example Primary"
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
  name                = "dr-vnet"
  location            = "westus"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "dr_subnet" {
  name                 = "dr-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.dr_vnet.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "azurerm_oracle_autonomous_database_cross_region_disaster_recovery" "dr_example" {
  name                          = "example-dr"
  display_name                  = "Example DR"
  location                      = "westus"
  resource_group_name           = azurerm_resource_group.example.name
  source_autonomous_database_id = azurerm_oracle_autonomous_database.primary.id
  subnet_id                     = azurerm_subnet.dr_subnet.id
  virtual_network_id            = azurerm_virtual_network.dr_vnet.id

  // Optional attributes
  replicate_automatic_backups_enabled = true

  tags = {
    Environment = "production"
    Purpose     = "disaster-recovery"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` (Required) - The name of the Resource Group where the resource should exist. Changing this forces creation of a new resource.

* `name` (Required) - The name for this Cross Region Disaster Recovery Autonomous Database. Changing this forces creation of a new resource.

* `location` (Required) - The Azure Region where the Cross Region Disaster Recovery Autonomous Database should exist. Must differ from the primary region. Changing this forces creation of a new resource.

* `display_name` (Required) - The user-friendly name for the Autonomous Database. Changing this forces creation of a new resource.

* `source_autonomous_database_id` (Required) - The Azure Resource ID of the source (primary) Autonomous Database. Changing this forces creation of a new resource.

* `subnet_id` (Required) - The Azure Resource ID of the subnet in the target region. Changing this forces creation of a new resource.

* `virtual_network_id` (Required) - The Azure Resource ID of the virtual network in the target region. Changing this forces creation of a new resource.

* `replicate_automatic_backups_enabled` (Optional) - If true, 7 days of backups are replicated across regions. Changing this forces creation of a new resource.

* `tags` (Optional) - Map of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Autonomous Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Cross Region Disaster Recovery Autonomous Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cross Region Disaster Recovery Autonomous Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cross Region Disaster Recovery Autonomous Database.

## Import

Cross Region Disaster Recovery Autonomous Database can be imported using the `resource id`.

```shell
terraform import azurerm_oracle_autonomous_database_cross_region_disaster_recovery.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/autonomousDatabases/autonomousDatabases1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
