---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_access_policy"
sidebar_current: "docs-azurerm-resource-key-vault-access-policy"
description: |-
  Manages a Key Vault Access Policy.
---

# azurerm_key_vault_access_policy

Manages a Key Vault Access Policy.

~> **NOTE:** It's possible to define Key Vault Access Policies both within [the `azurerm_key_vault` resource](key_vault.html) via the `access_policy` block and by using [the `azurerm_key_vault_access_policy` resource](key_vault_access_policy.html). However it's not possible to use both methods to manage Access Policies within a KeyVault, since there'll be conflicts.

-> **NOTE:** Azure permits a maximum of 1024 Access Policies per Key Vault - [more information can be found in this document](https://docs.microsoft.com/en-us/azure/key-vault/key-vault-secure-your-key-vault#data-plane-access-control).

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

  tenant_id = "22222222-2222-2222-2222-222222222222"

  enabled_for_disk_encryption = true

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_access_policy" "test" {
  vault_name          = "${azurerm_key_vault.test.name}"
  resource_group_name = "${azurerm_key_vault.test.resource_group_name}"

  tenant_id = "00000000-0000-0000-0000-000000000000"
  object_id = "11111111-1111-1111-1111-111111111111"

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
    the following: `backup`, `create`, `delete`, `deleteissuers`, `get`, `getissuers`, `import`, `list`, `listissuers`, 
    `managecontacts`, `manageissuers`, `purge`, `recover`, `restore`, `setissuers` and `update`.

* `key_permissions` - (Required) List of key permissions, must be one or more from
    the following: `backup`, `create`, `decrypt`, `delete`, `encrypt`, `get`, `import`, `list`, `purge`, 
    `recover`, `restore`, `sign`, `unwrapKey`, `update`, `verify` and `wrapKey`.

* `secret_permissions` - (Required) List of secret permissions, must be one or more
    from the following: `backup`, `delete`, `get`, `list`, `purge`, `recover`, `restore` and `set`.

* `storage_permissions` - (Optional) List of storage permissions, must be one or more from the following: `backup`, `delete`, `deletesas`, `get`, `getsas`, `list`, `listsas`, `purge`, `recover`, `regeneratekey`, `restore`, `set`, `setsas` and `update`.

## Attributes Reference

The following attributes are exported:

* `id` - Key Vault Access Policy ID.

-> **NOTE:** This Identifier is unique to Terraform and doesn't map to an existing object within Azure.

## Import

Key Vault Access Policies can be imported using the Resource ID of the Key Vault, plus some additional metadata.

If both an `object_id` and `application_id` are specified, then the Access Policy can be imported using the following code:

```shell
terraform import azurerm_key_vault_access_policy.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/test-vault/objectId/11111111-1111-1111-1111-111111111111/applicationId/22222222-2222-2222-2222-222222222222
```

where `11111111-1111-1111-1111-111111111111` is the `object_id` and `22222222-2222-2222-2222-222222222222` is the `application_id`.

---

Access Policies with an `object_id` but no `application_id` can be imported using the following command:

```shell
terraform import azurerm_key_vault_access_policy.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/vaults/test-vault/objectId/11111111-1111-1111-1111-111111111111
```

where `11111111-1111-1111-1111-111111111111` is the `object_id`.

-> **NOTE:** Both Identifiers are unique to Terraform and don't map to an existing object within Azure.
