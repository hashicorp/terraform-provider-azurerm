---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance_key"
description: |-
  Manages a Managed instance keys for transport data encryption.
---

# azurerm_mssql_managed_instance_key

Manages a Managed instance keys for transport data encryption.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "sql_mi-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "miDelegation"

    service_delegation {
      name = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name                = "example-routetable"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  route {
    name                   = "example"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_mssql_managed_instance" "example" {
  name                         = "sql-mi"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  administrator_login          = "demoReadUser"
  administrator_login_password = "ReadUser@123456"
  subnet_id                    = "${azurerm_subnet.example.id}"
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = 8
    family   = "Gen5"
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
  }
  license_type                 = "LicenseIncluded"
  collation                    = "SQL_Latin1_General_CP1_CI_AS"
  proxy_override               = "Redirect"
  storage_size_gb              = 64
  vcores                       = 8
  public_data_endpoint_enabled = false
  timezone_id                  = "Central America Standard Time"
  minimal_tls_version          = "1.2"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                        = "test-encrypt2"
  location                    = "westus2"
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_enabled         = true
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "get",
      "wrapKey",
      "unwrapKey",
      "create",
      "list",
      "delete",
    ]

    secret_permissions = [
      "get",
      "delete",
    ]

    storage_permissions = [
      "get",
      "delete",
    ]
  }

  access_policy {
    tenant_id = azurerm_mssql_managed_instance.example.identity[0].tenant_id
    object_id = azurerm_mssql_managed_instance.example.identity[0].principal_id

    key_permissions = [
      "get",
      "wrapKey",
      "unwrapKey",
      "create",
      "list",
      "delete",
    ]

    secret_permissions = [
      "get",
      "delete",
    ]

    storage_permissions = [
      "get",
      "delete",
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "example-key1"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_mssql_managed_instance_key" "test" {
  key_name              = "${azurerm_key_vault.example.name}_${azurerm_key_vault_key.example.name}_${azurerm_key_vault_key.example.version}"
  managed_instance_name = azurerm_mssql_managed_instance.example.name
  resource_group_name   = azurerm_mssql_managed_instance.example.resource_group_name
  uri                   = azurerm_key_vault_key.example.id
}

```

## Argument Reference

The following arguments are supported:

* `managed_instance_name` - (Required) The resource id of the MS SQL Managed Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The resource group of the managed instance. Changing this forces a new resource to be created.

* `key_name` - (Required) The managed instance encryption key name. The format for this should alwaye be keyvaultname_keyname_keyversion. Changing this forces a new resource to be created.

* `uri` - (Optional) The uri of the azure keyvault key along with its version.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed instance key.

* `thumbprint` - The managed instance key thumbprint

* `creation_date` - The managed instance creation date

* `name` - The managed instance key name.

* `type` - The managed instance key resource type.

* `server_key_type` - The key type. Defaults to `AzureKeyVault`.

* `kind` - The key kind.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating Managed instance key. 
* `update` - (Defaults to 30 minutes) Used when updating Managed instance key.
* `read` - (Defaults to 5 minutes) Used when retrieving Managed instance key.
* `delete` - (Defaults to 5 minutes) Used when deleting Managed instance key.

## Import

SQL managed instance key details can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/sql-mi/keys/<keyname>
```
