---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_network_acls"
description: |-
  Manages a Key Vault Network ACLs.
---

# azurerm_key_vault_network_acls

Manages a Key Vault Network ACLs.

~> **Note:** It's possible to define Key Vault Network ACLs both within [the `azurerm_key_vault` resource](key_vault.html) via the `network_acls` block and by using [the `azurerm_key_vault_network_acls` resource](key_vault_network_acls.html). However it's not possible to use both methods to manage Access Policies within a KeyVault, since there'll be conflicts.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "rg" {
  name     = "resourceGroup1"
  location = "East US"
}

resource "azurerm_key_vault" "kv" {
  name                = "testvault"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_virtual_network" "vnet" {
  name                = "virtual-network"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "subnet" {
  name                 = "kv-subnet"
  virtual_network_name = azurerm_virtual_network.vnet.name
  resource_group_name  = azurerm_resource_group.rg.name
  address_prefixes     = ["10.0.1.0/24"]
  service_endpoints    = ["Microsoft.KeyVault"]
}

resource "azurerm_key_vault_network_acls" "acls" {
  key_vault_name      = azurerm_key_vault.kv.name
  resource_group_name = azurerm_resource_group.rg.name
  network_acls {
    default_action             = "Deny"
    bypass                     = "None"
    ip_rules                   = ["43.0.0.0/24"]
    virtual_network_subnet_ids = [azurerm_subnet.subnet.id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `key_vault_name` - (Required) Specifies the name of the Key Vault. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Key Vault. Changing this forces a new resource to be created.

A `network_acls` block supports the following:

* `bypass` - (Required) Specifies which traffic can bypass the network rules. Possible values are `AzureServices` and `None`.

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_subnet_ids`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Key Vault.

* `virtual_network_subnet_ids` - (Optional) One or more Subnet ID's which should be able to access this Key Vault.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - Key Vault Network ACLs ID.

-> **NOTE:** This Identifier is unique to Terraform and doesn't map to an existing object within Azure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Network ACLs.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Network ACLs.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Network ACLs.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Network ACLs.

## Import

Key Vault Network ACLs can be imported using the Resource ID of the Key Vault.

```shell
terraform import azurerm_key_vault_network_acls.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/vault1
```
