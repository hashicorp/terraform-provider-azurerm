---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_key"
sidebar_current: "docs-azurerm-resource-key-vault-key"
description: |-
  Manages a Key Vault Key.

---

# azurerm_key_vault_key

Manages a Key Vault Key.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "my-resource-group"
  location = "West US"
}

resource "azurerm_key_vault" "test" {
  name                = "my-key-vault"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "premium"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    key_permissions = [
      "create",
      "get",
    ]

    secret_permissions = [
      "create",
      "set",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "generated" {
  name      = "generated-certificate"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"
  key_type  = "RSA"
  key_size  = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Key. Changing this forces a new resource to be created.

* `vault_uri` - (Required) Specifies the URI used to access the Key Vault instance, available on the `azurerm_key_vault` resource.

* `key_type` - (Required) Specifies the Key Type to use for this Key Vault Key. Possible values are `EC` (Elliptic Curve), `Oct` (Octet), `RSA` and `RSA-HSM`. Changing this forces a new resource to be created.

* `key_size` - (Required) Specifies the Size of the Key to create in bytes. For example, 1024 or 2048. Changing this forces a new resource to be created.

* `key_opts` - (Required) A list of JSON web key operations. Possible values include: `decrypt`, `encrypt`, `sign`, `unwrapKey`, `verify` and `wrapKey`. Please note these values are case sensitive.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Key ID.
* `version` - The current version of the Key Vault Key.
* `n` - The RSA modulus of this Key Vault Key.
* `e` - The RSA public exponent of this Key Vault Key.


## Import

Key Vault Key which is Enabled can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_key.test https://example-keyvault.vault.azure.net/keys/example/fdf067c93bbb4b22bff4d8b7a9a56217
```
