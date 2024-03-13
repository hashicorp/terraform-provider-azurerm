---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module_key"
description: |-
  Manages a Managed Hardware Security Module Key.

---

# azurerm_key_vault_managed_hardware_security_module_key

Manages a Managed Hardware Security Module Key.

## Example Usage

~> **Note:** To use this resource, your client should have RBAC roles with permissions like `Managed HSM Crypto Officer` or `Managed HSM Administrator` See [built-in-roles](https://learn.microsoft.com/en-us/azure/key-vault/managed-hsm/built-in-roles).

~> **Note:** The Azure Provider includes a Feature Toggle which will purge a Managed Hardware Security Module Key resource on destroy, rather than the default soft-delete. See [`purge_soft_deleted_keys_on_destroy`](https://registry.terraform.io/providers/hashicorp/azurerm/l.example.docs/guides/features-block#purge_soft_deleted_keys_on_destroy) for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_keys_on_destroy = true
      recover_soft_deleted_keys          = true
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}


resource "azurerm_key_vault" "example" {
  name                       = "acc240226165403061820"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
    ]
    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
}

resource "azurerm_key_vault_certificate" "cert" {
  count        = 3
  name         = "hsmcertexample${count.index}"
  key_vault_id = azurerm_key_vault.example.id
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }
    lifetime_action {
      action {
        action_type = "AutoRenew"
      }
      trigger {
        days_before_expiry = 30
      }
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      extended_key_usage = []
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_managed_hardware_security_module" example {
  name                     = "hsmexample"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false

  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 2
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "role_user" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.example.hsm_uri
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "role_user" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.example.hsm_uri
  name               = "706c03c7-69ad-33e5-2796-b3380d3a6e1a"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.role_user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "role_officer" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.example.hsm_uri
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "role_officer" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.example.hsm_uri
  name               = "d1a3242a-d521-11ee-9880-00155d316070"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.role_officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}


resource "azurerm_key_vault_managed_hardware_security_module_key" "example" {
  name           = "key-example"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.example.id
  key_type       = "EC-HSM"

  key_opts = [
    "sign",
    "verify",
  ]
  rotation_policy {
    expire_after_duration = "P66D"

    automatic {
      duration_before_expiry = "P30D"
    }
  }

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_user,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_officer
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Hardware Security Module Key. Changing this forces a new resource to be created.

* `managed_hsm_id` - (Required) The ID of the Managed Hardware Security Module where the Key should be created. Changing this forces a new resource to be created.

* `key_type` - (Required) Specifies the Key Type to use for this Managed Hardware Security Module Key. Possible values are `EC-HSM`, `RSA-HSM` and `oct-HSM`. Changing this forces a new resource to be created.

* `key_size` - (Optional) Specifies the Size of the RSA key to create in bytes. For example, `1024` or `2048`. *Note*: This field is required if `key_type` is `RSA-HSM`. Changing this forces a new resource to be created.

* `curve` - (Optional) Specifies the curve to use when creating an `EC` key. Possible values are `P-256`, `P-256K`, `P-384`, and `P-521`. This field will be required in a future release if `key_type` is `EC-HSM`. The API will default to `P-256` if nothing is specified. Changing this forces a new resource to be created.

* `key_options` - (Required) A list of JSON web key operations. Possible values are `decrypt`, `encrypt`, `import`, `export`, `sign`, `unwrapKey`, `verify` and `wrapKey`. Please note these values are case sensitive.

* `not_usable_before_date` - (Optional) Key not usable before the provided UTC datetime (Y-m-d'T'H:M:S'Z').

* `expiration_date` - (Optional) Expiration UTC datetime (Y-m-d'T'H:M:S'Z').

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `rotation_policy` - (Optional) A `rotation_policy` block as defined below.

---

A `rotation_policy` block supports the following:

* `expire_after_duration` - (Optional) Expire a Managed Hardware Security Module Key after given duration as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

* `automatic` - (Optional) An `automatic` block as defined below.

---

An `automatic` block supports the following:

* `duration_after_creation` - (Optional) Rotate automatically at a duration after create as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

* `duration_before_expiry` - (Optional) Rotate automatically at a duration before expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Managed Hardware Security Module Key ID.

* `resource_id` - The (Versioned) ID for this Managed Hardware Security Module Key. This property points to a specific version of a Managed Hardware Security Module Key, as such using this won't auto-rotate values if used in other Azure Services.

* `resource_versionless_id` - The Versionless ID of the Managed Hardware Security Module Key. This property allows other Azure Services (that support it) to auto-rotate their value when the Managed Hardware Security Module Key is updated.

* `version` - The current version of the Managed Hardware Security Module Key.

* `versionless_id` - The Base ID of the Managed Hardware Security Module Key.

* `n` - The RSA modulus of this Managed Hardware Security Module Key.

* `e` - The RSA public exponent of this Managed Hardware Security Module Key.

* `x` - The EC X component of this Managed Hardware Security Module Key.

* `y` - The EC Y component of this Managed Hardware Security Module Key.

* `public_key_pem` - The PEM encoded public key of this Managed Hardware Security Module Key.

* `public_key_openssh` - The OpenSSH encoded public key of this Managed Hardware Security Module Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Hardware Security Module Key.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Hardware Security Module Key.
* `read` - (Defaults to 30 minutes) Used when retrieving the Managed Hardware Security Module Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Hardware Security Module Key.

## Import

Managed Hardware Security Module Key which is Enabled can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module_key.example "https://example-mhsm.managedhsm.azure.net/keys/key-example/versionofthekey"
```
