---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate"
sidebar_current: "docs-azurerm-resource-key-vault-certificate"
description: |-
  Manages a Key Vault Certificate.

---

# azurerm_key_vault_certificate

Manages a Key Vault Certificate.

## Example Usage (Importing a PFX)

~> **Note:** this example assumed the PFX file is located in the same directory at `certificate-to-import.pfx`.

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "key-vault-certificate-example"
  location = "West Europe"
}

resource "azurerm_key_vault" "test" {
  name                = "keyvaultcertexample"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "all",
    ]

    key_permissions = [
      "all",
    ]

    secret_permissions = [
      "all",
    ]
  }

  tags {
    environment = "Production"
  }
}


resource "azurerm_key_vault_certificate" "test" {
  name      = "imported-cert"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"

  certificate {
    contents = "${base64encode(file("certificate-to-import.pfx"))}"
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

resource "azurerm_resource_group" "test" {
  name     = "key-vault-certificate-example"
  location = "West Europe"
}

resource "azurerm_key_vault" "test" {
  name                = "keyvaultcertexample"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.service_principal_object_id}"

    certificate_permissions = [
      "all",
    ]

    key_permissions = [
      "all",
    ]

    secret_permissions = [
      "all",
    ]
  }

  tags {
    environment = "Production"
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name      = "generated-cert"
  vault_uri = "${azurerm_key_vault.test.vault_uri}"

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
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Key Vault Certificate. Changing this forces a new resource to be created.

* `vault_uri` - (Required) Specifies the URI used to access the Key Vault instance, available on the `azurerm_key_vault` resource.

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
* `x509_certificate_properties` - (Optional) A `x509_certificate_properties` block as defined below.

`issuer_parameters` supports the following:

* `name` - (Required) The name of the Certificate Issuer. Possible values include `Self`, or the name of a certificate issuing authority supported by Azure. Changing this forces a new resource to be created.

`key_properties` supports the following:

* `exportable` - (Required) Is this Certificate Exportable? Changing this forces a new resource to be created.
* `key_size` - (Required) The size of the Key used in the Certificate. Possible values include `2048` and `4096`. Changing this forces a new resource to be created.
* `key_type` - (Required) Specifies the Type of Key, such as `RSA`. Changing this forces a new resource to be created.
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

* `key_usage` - (Required) A list of uses associated with this Key. Possible values include `cRLSign`, `dataEncipherment`, `decipherOnly`, `digitalSignature`, `encipherOnly`, `keyAgreement`, `keyCertSign`, `keyEncipherment` and `nonRepudiation` and are case-sensitive. Changing this forces a new resource to be created.
* `subject` - (Required) The Certificate's Subject. Changing this forces a new resource to be created.
* `validity_in_months` - (Required) The Certificates Validity Period in Months. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The Key Vault Certificate ID.
* `version` - The current version of the Key Vault Certificate.


## Import

Key Vault Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_certificate.test https://example-keyvault.vault.azure.net/certificates/example/fdf067c93bbb4b22bff4d8b7a9a56217
```
