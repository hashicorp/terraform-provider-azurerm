---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_network_acl_rule"
description: |-
  Manages a Key Vault Network ACL Rule.
---

# azurerm_key_vault_network_acl_rule

Manages a Key Vault Network ACL Rule.

~> **NOTE:** It's possible to define Key Vault Network ACL Rules both within [the `azurerm_key_vault` resource](key_vault.html) via the `network_acls` block and by using [the `azurerm_key_vault_network_acl_rule` resource](key_vault_network_acl_rule.html). However it's not possible to use both methods to manage Network ACL Rules within a KeyVault, since there'll be conflicts.

-> **NOTE:** Azure permits a maximum of 200 virtual network rules and 1000 IPv4 rules - [more information can be found in this document](https://docs.microsoft.com/en-us/azure/key-vault/general/network-security).

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.KeyVault"]
}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_key_vault_network_acl_rule" "example_ip" {
  key_vault_id = azurerm_key_vault.example.id
  source       = "1.2.3.4"
}

resource "azurerm_key_vault_network_acl_rule" "example_cidr" {
  key_vault_id = azurerm_key_vault.example.id
  source       = "40.82.252.0/24"
}

resource "azurerm_key_vault_network_acl_rule" "example_subnet" {
  key_vault_id = azurerm_key_vault.example.id
  source       = azurerm_subnet.example.id
}
```

## Argument Reference

The following arguments are supported:

* `key_vault_id` - (Required) Specifies the id of the Key Vault resource. Changing this
    forces a new resource to be created.

* `source` - (Required) The IP address, CIDR Blocks, or Subnet ID which should be able
    to access this Key Vault. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - Key Vault Access Policy ID.

-> **NOTE:** This Identifier is unique to Terraform and doesn't map to an existing object within Azure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Access Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Access Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Access Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Access Policy.

## Import

Key Vault Network ACL Rules can be imported using the Resource ID of the Key Vault, plus some additional metadata.


If the `source` is an IP address, then the Access Policy can be imported using the following code:

```shell
terraform import azurerm_key_vault_access_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/test-vault|1.2.3.4
```

where `1.2.3.4` is a public, IPv4 address.

---

If the `source` is CIDR block, then the Access Policy can be imported using the following code:

```shell
terraform import azurerm_key_vault_access_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/test-vault|40.82.252.0/24
```

where `40.82.252.0/24` is a public, IPv4 CIDR Block.

---

If the `source` is a virtual network subnet, then the Access Policy can be imported using the following code:

```shell
terraform import azurerm_key_vault_access_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/test-vault|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```

where `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1` is a virtual network subnet ID.

-> **NOTE:** All Identifiers are unique to Terraform and don't map to an existing object within Azure.
