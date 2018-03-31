---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_access_policy"
sidebar_current: "docs-azurerm-resource-key-vault-access-policy"
description: |-
  Create a Key Vault.
---

# azurerm\_key\_vault\_access\_policy

Assign access policies a Key Vault. This is useful for giving newly created VMs access to an existing Key Vault instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "vmss1"
  location            = "West US"
  resource_group_name = "resourceGroup1"
  upgrade_policy_mode = "Manual"
  overprovision       = true
  depends_on          = ["azurerm_lb.lb", "azurerm_virtual_network.vnet"]

  # Enable MSI authentication
  identity {
    type     = "systemAssigned"
  }

  extension {
    name                       = "MSILinuxExtension"
    publisher                  = "Microsoft.ManagedIdentity"
    type                       = "ManagedIdentityExtensionForLinux"
    type_handler_version       = "1.0"
    settings                   = "{\"port\": 50342}"
  }

  *# Remaining scaleset configuration omitted for brevity...*
}

resource "azurerm_key_vault_access_policy" "test" {
  tenant_id = "${var.tenantID}"

  # principal_id is the scale set's SPN
  object_id = "${lookup(azurerm_virtual_machine_scale_set.vmss1.identity[0], "principal_id")}"

  key_permissions = [
    "get",
  ]

  secret_permissions = [
    "get",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault resource. Changing this
    forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists.
    Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the namespace. Changing this forces a new resource to be created.

* `tenant_id` - (Required) The Azure Active Directory tenant ID that should be
    used for authenticating requests to the key vault.

* `certificate_permissions` - (Optional) List of certificate permissions, must be one or more from
    the following: `create`, `delete`, `deleteissuers`, `get`, `getissuers`, `import`, `list`, `listissuers`, `managecontacts`, `manageissuers`, `purge`, `recover`, `setissuers` and `update`.

* `key_permissions` - (Required) List of key permissions, must be one or more from
    the following: `backup`, `create`, `decrypt`, `delete`, `encrypt`, `get`, `import`, `list`, `purge`, `recover`, `restore`, `sign`, `unwrapKey`, `update`, `verify` and `wrapKey`.

* `secret_permissions` - (Required) List of secret permissions, must be one or more
    from the following: `backup`, `delete`, `get`, `list`, `purge`, `recover`, `restore` and `set`.

## Import

Access policies can be imported from any Key Vault using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/vault1
```
