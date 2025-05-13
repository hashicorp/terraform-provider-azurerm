---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module_key"
description: |-
  Manages a Key Vault Managed Hardware Security Module Key.
---

# azurerm_key_vault_managed_hardware_security_module_key

Manages a Key Vault Managed Hardware Security Module Key.

~> **Note:** The Azure Provider includes a Feature Toggle which will purge a Key Vault Managed Hardware Security Module Key resource on destroy, rather than the default soft-delete. See [`purge_soft_deleted_hardware_security_modules_on_destroy`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block#purge_soft_deleted_hardware_security_module_keys_on_destroy) for more information.

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault_managed_hardware_security_module" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false

  active_config {
    security_domain_certificate = [
      azurerm_key_vault_certificate.cert[0].id,
      azurerm_key_vault_certificate.cert[1].id,
      azurerm_key_vault_certificate.cert[2].id,
    ]
    security_domain_quorum = 2
  }
}

// this gives your service principal the HSM Crypto User role which lets you create and destroy hsm keys
resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "hsm-crypto-user" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "1e243909-064c-6ac3-84e9-1c8bf8d6ad22"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = data.azurerm_client_config.current.object_id
}

// this gives your service principal the HSM Crypto Officer role which lets you purge hsm keys
resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "hsm-crypto-officer" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "1e243909-064c-6ac3-84e9-1c8bf8d6ad23"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "example" {
  name           = "example"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "EC-HSM"
  curve          = "P-521"
  key_opts       = ["sign"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test1
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Managed Hardware Security Module Key. Changing this forces a new resource to be created.

* `managed_hsm_id` - (Required) Specifies the ID of the Key Vault Managed Hardware Security Module that they key will be owned by. Changing this forces a new resource to be created.

* `key_type` - (Required) Specifies the Key Type to use for this Key Vault Managed Hardware Security Module Key. Possible values are `EC-HSM`, `oct-HSM` and `RSA-HSM`. More details see [HSM-protected keys](https://learn.microsoft.com/en-us/azure/key-vault/keys/about-keys#hsm-protected-keys). Changing this forces a new resource to be created.

* `key_size` - (Optional) Specifies the Size of the RSA key to create in bytes. For example, 1024 or 2048. *Note*: This field is required if `key_type` is `RSA-HSM` or `oct-HSM`. Changing this forces a new resource to be created.

* `curve` - (Optional) Specifies the curve to use when creating an `EC-HSM` key. Possible values are `P-256`, `P-256K`, `P-384`, and `P-521`. This field is required if `key_type` is `EC-HSM`. Changing this forces a new resource to be created.

* `key_opts` - (Required) A list of JSON web key operations. Possible values include: `decrypt`, `encrypt`, `sign`, `unwrapKey`, `verify`, `wrapKey` and `import`. Please note these values are case-sensitive.

* `not_before_date` - (Optional) Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').

~> **Note:** Once `expiration_date` is set, it's not possible to unset the key even if it is deleted & recreated as underlying Azure API uses the restore of the purged key.

* `expiration_date` - (Optional) Expiration UTC datetime (Y-m-d'T'H:M:S'Z'). When this parameter gets changed on reruns, if newer date is ahead of current date, an update is performed. If the newer date is before the current date, resource will be force created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Key Vault Secret Managed Hardware Security Module Key ID.

* `versioned_id` - The versioned Key Vault Secret Managed Hardware Security Module Key ID.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Managed Hardware Security Module Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Hardware Security Module Key.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Managed Hardware Security Module Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Managed Hardware Security Module Key.

## Import

Key Vault Managed Hardware Security Module Key can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module_key.example https://exampleHSM.managedhsm.azure.net/keys/exampleKey
```
