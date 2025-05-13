---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance_transparent_data_encryption"
description: |-
  Manages the transparent data encryption configuration for a MSSQL Managed Instance
---

# azurerm_mssql_managed_instance_transparent_data_encryption

Manages the transparent data encryption configuration for a MSSQL Managed Instance

~> **Note:** Once transparent data encryption(TDE) is enabled on a MS SQL instance, it is not possible to remove TDE. You will be able to switch between 'ServiceManaged' and 'CustomerManaged' keys, but will not be able to remove encryption. For safety when this resource is deleted, the TDE mode will automatically be set to 'ServiceManaged'. See `key_vault_uri` for more information on how to specify the key types. As SQL Managed Instance only supports a single configuration for encryption settings, this resource will replace the current encryption settings on the server.

~> **Note:** See [documentation](https://docs.microsoft.com/azure/azure-sql/database/transparent-data-encryption-byok-overview) for important information on how handle lifecycle management of the keys to prevent data lockout.

## Example Usage with Service Managed Key

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "EastUs"
}

resource "azurerm_virtual_network" "example" {
  name                = "acctest-vnet1-mssql"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "example" {
  name                 = "subnet1-mssql"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_mssql_managed_instance" "example" {
  name                = "mssqlinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.example.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_mssql_managed_instance_transparent_data_encryption" "example" {
  managed_instance_id = azurerm_mssql_managed_instance.example.id
}
```

## Example Usage with Customer Managed Key

```hcl

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "EastUs"
}

resource "azurerm_virtual_network" "example" {
  name                = "acctest-vnet1-mssql"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_subnet" "example" {
  name                 = "subnet1-mssql"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_mssql_managed_instance" "example" {
  name                = "mssqlinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.example.id
  vcores             = 4

  administrator_login          = "missadministrator"
  administrator_login_password = "NCC-1701-D"

  identity {
    type = "SystemAssigned"
  }
}

# Create a key vault with policies for the deployer to create a key & SQL Managed Instance to wrap/unwrap/get key
resource "azurerm_key_vault" "example" {
  name                        = "example"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get", "List", "Create", "Delete", "Update", "Recover", "Purge", "GetRotationPolicy",
    ]
  }
  access_policy {
    tenant_id = azurerm_mssql_managed_instance.example.identity[0].tenant_id
    object_id = azurerm_mssql_managed_instance.example.identity[0].principal_id

    key_permissions = [
      "Get", "WrapKey", "UnwrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "byok"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "unwrapKey",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault.example
  ]
}

resource "azurerm_mssql_managed_instance_transparent_data_encryption" "example" {
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  key_vault_key_id    = azurerm_key_vault_key.example.id
}


```

## Arguments Reference

The following arguments are supported:

* `managed_instance_id` - (Required) Specifies the name of the MS SQL Managed Instance. Changing this forces a new resource to be created.

---

* `key_vault_key_id` - (Optional) To use customer managed keys from Azure Key Vault, provide the AKV Key ID. To use service managed keys, omit this field.

* `managed_hsm_key_id` -  (Optional)  To use customer managed keys from a managed HSM, provide the Managed HSM Key ID. To use service managed keys, omit this field.

~> **Note:** In order to use customer managed keys, the identity of the MSSQL Managed Instance must have the following permissions on the key vault: 'get', 'wrapKey' and 'unwrapKey'

~> **Note:** If `managed_instance_id` denotes a secondary instance deployed for disaster recovery purposes, then the `key_vault_key_id` should be the same key used for the primary instance's transparent data encryption. Both primary and secondary instances should be encrypted with same key material.

* `auto_rotation_enabled` - (Optional) When enabled, the SQL Managed Instance will continuously check the key vault for any new versions of the key being used as the TDE protector. If a new version of the key is detected, the TDE protector on the SQL Managed Instance will be automatically rotated to the latest key version within 60 minutes.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the MSSQL encryption protector

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MSSQL.
* `read` - (Defaults to 5 minutes) Used when retrieving the MSSQL.
* `update` - (Defaults to 30 minutes) Used when updating the MSSQL.
* `delete` - (Defaults to 30 minutes) Used when deleting the MSSQL.

## Import

~> **Note:** This resource does not need to be imported to manage it, however the import will work.

SQL Managed Instance Transparent Data Encryption can be imported using the resource id, e.g.

```shell
terraform import azurerm_mssql_managed_instance_transparent_data_encryption.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/instance1/encryptionProtector/current
```

