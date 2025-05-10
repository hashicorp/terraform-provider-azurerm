---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate"
description: |-
  Manages a Key Vault Certificate.

---

# azurerm_key_vault_certificate

Manages a Key Vault Certificate.

~> **Note:** The Azure Provider includes a Feature Toggle which will purge a Key Vault Certificate resource on destroy, rather than the default soft-delete. See [`purge_soft_deleted_certificates_on_destroy`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block#purge_soft_deleted_certificates_on_destroy) for more information.

## Example Usage (Importing a PFX)

~> **Note:** this example assumed the PFX file is located in the same directory at `certificate-to-import.pfx`.

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_certificates_on_destroy = true
      recover_soft_deleted_certificates          = true
    }
  }
}

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
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "GetIssuers",
      "Import",
      "List",
      "ListIssuers",
      "ManageContacts",
      "ManageIssuers",
      "SetIssuers",
      "Update",
    ]

    key_permissions = [
      "Backup",
      "Create",
      "Decrypt",
      "Delete",
      "Encrypt",
      "Get",
      "Import",
      "List",
      "Purge",
      "Recover",
      "Restore",
      "Sign",
      "UnwrapKey",
      "Update",
      "Verify",
      "WrapKey",
    ]

    secret_permissions = [
      "Backup",
      "Delete",
      "Get",
      "List",
      "Purge",
      "Recover",
      "Restore",
      "Set",
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
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "GetIssuers",
      "Import",
      "List",
      "ListIssuers",
      "ManageContacts",
      "ManageIssuers",
      "Purge",
      "SetIssuers",
      "Update",
    ]

    key_permissions = [
      "Backup",
      "Create",
      "Decrypt",
      "Delete",
      "Encrypt",
      "Get",
      "Import",
      "List",
      "Purge",
      "Recover",
      "Restore",
      "Sign",
      "UnwrapKey",
      "Update",
      "Verify",
      "WrapKey",
    ]

    secret_permissions = [
      "Backup",
      "Delete",
      "Get",
      "List",
      "Purge",
      "Recover",
      "Restore",
      "Set",
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

* `key_vault_id` - (Required) The ID of the Key Vault where the Certificate should be created. Changing this forces a new resource to be created.

* `certificate` - (Optional) A `certificate` block as defined below, used to Import an existing certificate. Changing this will create a new version of the Key Vault Certificate.

* `certificate_policy` - (Optional) A `certificate_policy` block as defined below. Changing this (except the `lifetime_action` field) will create a new version of the Key Vault Certificate.

~> **Note:** When creating a Key Vault Certificate, at least one of `certificate` or `certificate_policy` is required. Provide `certificate` to import an existing certificate, `certificate_policy` to generate a new certificate.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `certificate` block supports the following:

* `contents` - (Required) The base64-encoded certificate contents.
* `password` - (Optional) The password associated with the certificate.

~> **Note:** A PEM certificate is already base64 encoded. To successfully import, the `contents` property should include a PEM encoded X509 certificate and a private_key in pkcs8 format. There should only be linux style `\n` line endings and the whole block should have the PEM begin/end blocks around the certificate data and the private key data.

To convert a private key to pkcs8 format with openssl use:
```shell
openssl pkcs8 -topk8 -nocrypt -in private_key.pem > private_key_pk8.pem
```

The PEM content should look something like:
```text
-----BEGIN CERTIFICATE-----
aGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8K
:
aGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8KaGVsbG8K
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
d29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQK
:
d29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQKd29ybGQK
-----END PRIVATE KEY-----
```

---

The `certificate_policy` block supports the following:

* `issuer_parameters` - (Required) A `issuer_parameters` block as defined below.
* `key_properties` - (Required) A `key_properties` block as defined below.
* `lifetime_action` - (Optional) A `lifetime_action` block as defined below.
* `secret_properties` - (Required) A `secret_properties` block as defined below.
* `x509_certificate_properties` - (Optional) A `x509_certificate_properties` block as defined below. Required when `certificate` block is not specified.

---

The `issuer_parameters` block supports the following:

* `name` - (Required) The name of the Certificate Issuer. Possible values include `Self` (for self-signed certificate), or `Unknown` (for a certificate issuing authority like `Let's Encrypt` and Azure direct supported ones).

---

The `key_properties` block supports the following:

* `curve` - (Optional) Specifies the curve to use when creating an `EC` key. Possible values are `P-256`, `P-256K`, `P-384`, and `P-521`. This field will be required in a future release if `key_type` is `EC` or `EC-HSM`.
* `exportable` - (Required) Is this certificate exportable?
* `key_size` - (Optional) The size of the key used in the certificate. Possible values include `2048`, `3072`, and `4096` for `RSA` keys, or `256`, `384`, and `521` for `EC` keys. This property is required when using RSA keys.
* `key_type` - (Required) Specifies the type of key. Possible values are `EC`, `EC-HSM`, `RSA`, `RSA-HSM` and `oct`.
* `reuse_key` - (Required) Is the key reusable?

---

The `lifetime_action` block supports the following:

* `action` - (Required) A `action` block as defined below.
* `trigger` - (Required) A `trigger` block as defined below.

---

The `action` block supports the following:

* `action_type` - (Required) The Type of action to be performed when the lifetime trigger is triggerec. Possible values include `AutoRenew` and `EmailContacts`.

---

The `trigger` block supports the following:

* `days_before_expiry` - (Optional) The number of days before the Certificate expires that the action associated with this Trigger should run. Conflicts with `lifetime_percentage`.
* `lifetime_percentage` - (Optional) The percentage at which during the Certificates Lifetime the action associated with this Trigger should run. Conflicts with `days_before_expiry`.

---

The `secret_properties` block supports the following:

* `content_type` - (Required) The Content-Type of the Certificate, such as `application/x-pkcs12` for a PFX or `application/x-pem-file` for a PEM.

---

The `x509_certificate_properties` block supports the following:

* `extended_key_usage` - (Optional) A list of Extended/Enhanced Key Usages.
* `key_usage` - (Required) A list of uses associated with this Key. Possible values include `cRLSign`, `dataEncipherment`, `decipherOnly`, `digitalSignature`, `encipherOnly`, `keyAgreement`, `keyCertSign`, `keyEncipherment` and `nonRepudiation` and are case-sensitive.
* `subject` - (Required) The Certificate's Subject.
* `subject_alternative_names` - (Optional) A `subject_alternative_names` block as defined below.
* `validity_in_months` - (Required) The Certificates Validity Period in Months.

---

The `subject_alternative_names` block supports the following:

* `dns_names` - (Optional) A list of alternative DNS names (FQDNs) identified by the Certificate.
* `emails` - (Optional) A list of email addresses identified by this Certificate.
* `upns` - (Optional) A list of User Principal Names identified by the Certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Key Vault Certificate ID.
* `secret_id` - The ID of the associated Key Vault Secret.
* `version` - The current version of the Key Vault Certificate.
* `versionless_id` - The Base ID of the Key Vault Certificate.
* `versionless_secret_id` - The Base ID of the Key Vault Secret.
* `certificate_data` - The raw Key Vault Certificate data represented as a hexadecimal string.
* `certificate_data_base64` - The Base64 encoded Key Vault Certificate data.
* `thumbprint` - The X509 Thumbprint of the Key Vault Certificate represented as a hexadecimal string.
* `certificate_attribute` - A `certificate_attribute` block as defined below.
 
* `resource_manager_id` - The (Versioned) ID for this Key Vault Certificate. This property points to a specific version of a Key Vault Certificate, as such using this won't auto-rotate values if used in other Azure Services.

* `resource_manager_versionless_id` - The Versionless ID of the Key Vault Certificate. This property allows other Azure Services (that support it) to auto-rotate their value when the Key Vault Certificate is updated.

---

A `certificate_attribute` block exports the following:

* `created` - The create time of the Key Vault Certificate.
* `enabled` - whether the Key Vault Certificate is enabled.
* `expires` - The expires time of the Key Vault Certificate.
* `not_before` - The not before valid time of the Key Vault Certificate.
* `recovery_level` - The deletion recovery level of the Key Vault Certificate.
* `updated` - The recent update time of the Key Vault Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Key Vault Certificate.
* `read` - (Defaults to 30 minutes) Used when retrieving the Key Vault Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Key Vault Certificate.

## Import

Key Vault Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_certificate.example "https://example-keyvault.vault.azure.net/certificates/example/fdf067c93bbb4b22bff4d8b7a9a56217"
```
