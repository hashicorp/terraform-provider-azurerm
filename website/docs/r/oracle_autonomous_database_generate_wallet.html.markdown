---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_autonomous_database_wallet"
description: |-
  Manages an Autonomous Database Wallet.
---

# azurerm_oracle_autonomous_database_wallet

Generates and manages a wallet for an Oracle Autonomous Database. The wallet contains connection information and credentials required to connect to the Autonomous Database.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

# Create an Autonomous Database first
resource "azurerm_oracle_autonomous_database" "example" {
  name                = "example-adb"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  # Required properties for demonstration purposes
  admin_password                   = "StrongPassword12#$"
  backup_retention_period_in_days  = 60
  character_set                    = "AL32UTF8"
  compute_count                    = 2
  compute_model                    = "ECPU"
  data_storage_size_in_tbs         = 1
  db_version                       = "19c"
  db_workload                      = "OLTP"
  display_name                     = "example-adb"
  license_model                    = "BringYourOwnLicense"
  national_character_set           = "AL16UTF16"
  mtls_connection_required         = true
  auto_scaling_enabled             = true
  auto_scaling_for_storage_enabled = false
}

# Generate a wallet for the Autonomous Database
resource "azurerm_oracle_autonomous_database_wallet" "example" {
  autonomous_database_id = azurerm_oracle_autonomous_database.example.id
  password               = "WalletPassword123!"

  # Optional configurations
  generate_type = "SINGLE"
  is_regional   = false
}
``` 

## Arguments Reference
The following arguments are supported:

* `autonomous_database_id` - (Required) The ID of the Autonomous Database for which to generate a wallet. Changing this forces a new Wallet to be created.

* `password` - (Required) The password to protect the wallet. Must be between 8 and 30 characters. Changing this forces a new Wallet to be created.

* `generate_type` - (Optional) The type of wallet to generate. Valid values are SINGLE and ALL. Default is SINGLE. Changing this forces a new Wallet to be created.

    * `SINGLE` - Generate a wallet for a single database only.
    * `ALL - Generate` a wallet containing connection credentials for all Autonomous Databases in the region.
  
*`is_regional` - (Optional) Whether to create a regional wallet, which will include credentials for all database instances in the region. Default is false. Changing this forces a new Wallet to be created.

## Attributes Reference
In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Autonomous Database Wallet resource.

* `wallet_files` - The base64-encoded wallet file content. This attribute is sensitive.

## Timeouts
The timeouts block allows you to specify timeouts for certain actions:

* `create` - (Defaults to 10 minutes) Used when generating the Wallet.
* `read` - (Defaults to 5 minutes) Used when retrieving the Wallet.
