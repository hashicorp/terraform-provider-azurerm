---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_certificate"
description: |-
  Gets information about an existing Key Vault Certificate.
---

# Data Source: azurerm_key_vault_certificate

Use this data source to access information about an existing Key Vault Certificate.

~> **Note:** All arguments including the secret value will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
data "azurerm_key_vault" "example" {
  name                = "examplekv"
  resource_group_name = "some-resource-group"
}

data "azurerm_key_vault_certificate" "example" {
  name         = "secret-sauce"
  key_vault_id = data.azurerm_key_vault.example.id
}

output "certificate_thumbprint" {
  value = data.azurerm_key_vault_certificate.example.thumbprint
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Key Vault Certificate.

* `key_vault_id` - Specifies the ID of the Key Vault instance where the Secret resides, available on the `azurerm_key_vault` Data Source / Resource.

* `version` - (Optional) Specifies the version of the certificate to look up.  (Defaults to latest)

**NOTE:** The vault must be in the same subscription as the provider. If the vault is in another subscription, you must create an aliased provider for that subscription.

## Attributes Reference

The following attributes are exported:


* `id` - The Key Vault Certificate ID.
*
* `name` - Specifies the name of the Key Vault Certificate.
*
* `secret_id` - The ID of the associated Key Vault Secret.
*
* `version` - The current version of the Key Vault Certificate.
*
* `certificate_data` - The raw Key Vault Certificate data represented as a hexadecimal string.
*
* `thumbprint` - The X509 Thumbprint of the Key Vault Certificate represented as a hexadecimal string.

* `certificate_policy` - A `certificate_policy` block as defined below.

* `tags` - A mapping of tags to assign to the resource.

---

`certificate_policy` exports the following:

* `issuer_parameters` - A `issuer_parameters` block as defined below.
* `key_properties` - A `key_properties` block as defined below.
* `lifetime_action` - A `lifetime_action` block as defined below.
* `secret_properties` - A `secret_properties` block as defined below.
* `x509_certificate_properties` - An `x509_certificate_properties` block as defined below.

---

`issuer_parameters` exports the following:

* `name` - The name of the Certificate Issuer.

---

`key_properties` exports the following:

* `exportable` - Is this Certificate Exportable?
* `key_size` - The size of the Key used in the Certificate.
* `key_type` - Specifies the Type of Key, for example `RSA`.
* `reuse_key` - Is the key reusable?

---

`lifetime_action` exports the following:

* `action` - A `action` block as defined below.
* `trigger` - A `trigger` block as defined below.

---

`action` exports the following:

* `action_type` - The Type of action to be performed when the lifetime trigger is triggerec.

---

`trigger` exports the following:

* `days_before_expiry` - The number of days before the Certificate expires that the action associated with this Trigger should run.
* `lifetime_percentage` - The percentage at which during the Certificates Lifetime the action associated with this Trigger should run.

---

`secret_properties` exports the following:

* `content_type` - The Content-Type of the Certificate, for example `application/x-pkcs12` for a PFX or `application/x-pem-file` for a PEM.

---

`x509_certificate_properties` exports the following:

* `extended_key_usage` - A list of Extended/Enhanced Key Usages.
* `key_usage` - A list of uses associated with this Key.
* `subject` - The Certificate's Subject.
* `subject_alternative_names` - A `subject_alternative_names` block as defined below.
* `validity_in_months` - The Certificates Validity Period in Months.

---

`subject_alternative_names` exports the following:

* `dns_names` - A list of alternative DNS names (FQDNs) identified by the Certificate.
* `emails` - A list of email addresses identified by this Certificate.
* `upns` - A list of User Principal Names identified by the Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Certificate.
