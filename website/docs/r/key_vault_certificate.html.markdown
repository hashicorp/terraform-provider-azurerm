---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate"
description: |-
  Manages a Key Vault Certificate.

---

# azurerm_key_vault_certificate

Manages a Key Vault Certificate.

## Example Usage (Importing a PFX)

~> **Note:** this example assumed the PFX file is located in the same directory at `certificate-to-import.pfx`.

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "create",
      "delete",
      "deleteissuers",
      "get",
      "getissuers",
      "import",
      "list",
      "listissuers",
      "managecontacts",
      "manageissuers",
      "setissuers",
      "update",
    ]

    key_permissions = [
      "backup",
      "create",
      "decrypt",
      "delete",
      "encrypt",
      "get",
      "import",
      "list",
      "purge",
      "recover",
      "restore",
      "sign",
      "unwrapKey",
      "update",
      "verify",
      "wrapKey",
    ]

    secret_permissions = [
      "backup",
      "delete",
      "get",
      "list",
      "purge",
      "recover",
      "restore",
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "imported-cert"
  key_vault_id = azurerm_key_vault.example.id

  certificate {
    contents = filebase64("certificate-to-import.pfx")
    password = ""
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}
```

## Example Usage (Generating a new certificate)

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                       = "examplekeyvault"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "create",
      "delete",
      "deleteissuers",
      "get",
      "getissuers",
      "import",
      "list",
      "listissuers",
      "managecontacts",
      "manageissuers",
      "purge",
      "setissuers",
      "update",
    ]

    key_permissions = [
      "backup",
      "create",
      "decrypt",
      "delete",
      "encrypt",
      "get",
      "import",
      "list",
      "purge",
      "recover",
      "restore",
      "sign",
      "unwrapKey",
      "update",
      "verify",
      "wrapKey",
    ]

    secret_permissions = [
      "backup",
      "delete",
      "get",
      "list",
      "purge",
      "recover",
      "restore",
      "set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "generated-cert"
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
      # Server Authentication = 1.3.6.1.5.5.7.3.1
      # Client Authentication = 1.3.6.1.5.5.7.3.2
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]

      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject_alternative_names {
        dns_names = ["internal.contoso.com", "domain.hello.world"]
      }

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Certificate. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of the Key Vault where the Certificate should be created.

* `certificate` - (Optional) A `certificate` block as defined below, used to Import an existing certificate.

* `certificate_policy` - (Required) A `certificate_policy` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`certificate` supports the following:

* `contents` - (Required) The base64-encoded certificate contents. Changing this forces a new resource to be created.
* `password` - (Optional) The password associated with the certificate. Changing this forces a new resource to be created.

`certificate_policy` supports the following:

* `issuer_parameters` - (Required) A `issuer_parameters` block as defined below.
* `key_properties` - (Required) A `key_properties` block as defined below.
* `lifetime_action` - (Optional) A `lifetime_action` block as defined below.
* `secret_properties` - (Required) A `secret_properties` block as defined below.
* `x509_certificate_properties` - (Optional) A `x509_certificate_properties` block as defined below. Required when `certificate` block is not specified.

`issuer_parameters` supports the following:

* `name` - (Required) The name of the Certificate Issuer. Possible values include `Self` (for self-signed certificate), or `Unknown` (for a certificate issuing authority like `Let's Encrypt` and Azure direct supported ones). Changing this forces a new resource to be created.

`key_properties` supports the following:

* `curve` - (Optional) Specifies the curve to use when creating an `EC` key. Possible values are `P-256`, `P-256K`, `P-384`, and `P-521`. This field will be required in a future release if `key_type` is `EC` or `EC-HSM`. Changing this forces a new resource to be created.
* `exportable` - (Required) Is this certificate exportable? Changing this forces a new resource to be created.
* `key_size` - (Optional) The size of the key used in the certificate. Possible values include `2048`, `3072`, and `4096` for `RSA` keys, or `256`, `384`, and `521` for `EC` keys. This property is required when using RSA keys. Changing this forces a new resource to be created.
* `key_type` - (Required) Specifies the type of key, such as `RSA` or `EC`. Changing this forces a new resource to be created.
* `reuse_key` - (Required) Is the key reusable? Changing this forces a new resource to be created.

`lifetime_action` supports the following:

* `action` - (Required) A `action` block as defined below.
* `trigger` - (Required) A `trigger` block as defined below.

`action` supports the following:

* `action_type` - (Required) The Type of action to be performed when the lifetime trigger is triggerec. Possible values include `AutoRenew` and `EmailContacts`. Changing this forces a new resource to be created.

`trigger` supports the following:

* `days_before_expiry` - (Optional) The number of days before the Certificate expires that the action associated with this Trigger should run. Changing this forces a new resource to be created. Conflicts with `lifetime_percentage`.
* `lifetime_percentage` - (Optional) The percentage at which during the Certificates Lifetime the action associated with this Trigger should run. Changing this forces a new resource to be created. Conflicts with `days_before_expiry`.

`secret_properties` supports the following:

* `content_type` - (Required) The Content-Type of the Certificate, such as `application/x-pkcs12` for a PFX or `application/x-pem-file` for a PEM. Changing this forces a new resource to be created.

`x509_certificate_properties` supports the following:

* `extended_key_usage` - (Optional) A list of Extended/Enhanced Key Usages. Changing this forces a new resource to be created.
* `key_usage` - (Required) A list of uses associated with this Key. Possible values include `cRLSign`, `dataEncipherment`, `decipherOnly`, `digitalSignature`, `encipherOnly`, `keyAgreement`, `keyCertSign`, `keyEncipherment` and `nonRepudiation` and are case-sensitive. Changing this forces a new resource to be created.
* `subject` - (Required) The Certificate's Subject. Changing this forces a new resource to be created.
* `subject_alternative_names` - (Optional) A `subject_alternative_names` block as defined below.
* `validity_in_months` - (Required) The Certificates Validity Period in Months. Changing this forces a new resource to be created.

`subject_alternative_names` supports the following:

* `dns_names` - (Optional) A list of alternative DNS names (FQDNs) identified by the Certificate. Changing this forces a new resource to be created.
* `emails` - (Optional) A list of email addresses identified by this Certificate. Changing this forces a new resource to be created.
* `upns` - (Optional) A list of User Principal Names identified by the Certificate. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Certificate ID.
* `secret_id` - The ID of the associated Key Vault Secret.
* `version` - The current version of the Key Vault Certificate.
* `certificate_data` - The raw Key Vault Certificate data represented as a hexadecimal string.
* `certificate_data_base64` - The Base64 encoded Key Vault Certificate data.
* `thumbprint` - The X509 Thumbprint of the Key Vault Certificate represented as a hexadecimal string.
* `certificate_attribute` - A `certificate_attribute` block as defined below.

---

A `certificate_attribute` block exports the following:

* `created` - The create time of the Key Vault Certificate.
* `enabled` - whether the Key Vault Certificate is enabled.
* `expires` - The expires time of the Key Vault Certificate.
* `not_before` - The not before valid time of the Key Vault Certificate.
* `recovery_level` - The deletion recovery level of the Key Vault Certificate.
* `updated` - The recent update time of the Key Vault Certificate.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Certificate.

## Import

Key Vault Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_certificate.example "https://example-keyvault.vault.azure.net/certificates/example/fdf067c93bbb4b22bff4d8b7a9a56217"
```
