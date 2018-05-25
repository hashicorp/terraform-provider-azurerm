---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_access_policy"
sidebar_current: "docs-azurerm-resource-key-vault-access-policy"
description: |-
  Manages a Key Vault Access Policy.
---

# azurerm_key_vault_access_policy

Manages a Key Vault Access Policy.

~> **NOTE on Key Vaults and Key Vault Policies:** Terraform currently
provides both a standalone [Key Vault Policy Resource](key_vault_policy.html), and allows for Key Vault Access Polcies to be defined in-line within the [Key Vault Resource](key_vault.html).
At this time you cannot define Key Vault Policy with in-line Key Vault in conjunction with any Key Vault Policy resources. Doing so may cause a conflict of Access Policies and will overwrite Access Policies.


## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "${azurerm_resource_group.test.location}"
}

resource "azurerm_key_vault" "test" {
  name                = "testvault"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "standard"
  }

  tenant_id = "d6e396d0-5584-41dc-9fc0-268df99bc610"

  enabled_for_disk_encryption = true

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_policy" "test" {
  vault_name           = "${azurerm_key_vault.test.name}"
  resource_group_name  = "${azurerm_key_vault.test.resource_group_name}"
  
  tenant_id = "d6e396d0-5584-41dc-9fc0-268df99bc610"
  object_id = "d746815a-0433-4a21-b95d-fc437d2d475b"

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

* `vault_name` - (Required) Specifies the name of the Key Vault resource. Changing this
    forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the namespace. Changing this forces a new resource to be created.

* `tenant_id` - (Required) The Azure Active Directory tenant ID that should be used
    for authenticating requests to the key vault. Changing this forces a new resource 
    to be created.

* `object_id` - (Required) The object ID of a user, service principal or security
    group in the Azure Active Directory tenant for the vault. The object ID must
    be unique for the list of access policies. Changing this forces a new resource 
    to be created.

* `application_id` - (Optional) The object ID of an Application in Azure Active Directory.

* `certificate_permissions` - (Optional) List of certificate permissions, must be one or more from
    the following: `create`, `delete`, `deleteissuers`, `get`, `getissuers`, `import`, `list`, `listissuers`, `managecontacts`, `manageissuers`, `purge`, `recover`, `setissuers` and `update`.

* `key_permissions` - (Required) List of key permissions, must be one or more from
    the following: `backup`, `create`, `decrypt`, `delete`, `encrypt`, `get`, `import`, `list`, `purge`, `recover`, `restore`, `sign`, `unwrapKey`, `update`, `verify` and `wrapKey`.

* `secret_permissions` - (Required) List of secret permissions, must be one or more
    from the following: `backup`, `delete`, `get`, `list`, `purge`, `recover`, `restore` and `set`.

## Attributes Reference

The following attributes are exported:

* `id` - Key Vault Access Policy ID.
