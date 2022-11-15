---
subcategory: "Confidential Ledger"
layout: "azurerm"
page_title: "Confidential Ledger: azurerm_confidential_ledger"
description: |-
    Gets information about an existing Confidential Ledger.
---

# Data Source: azurerm_confidential_ledger

Gets information about an existing Confidential Ledger.

## Example Usage

```hcl
data "azurerm_confidential_ledger" "current" {
  name                = "example-ledger"
  resource_group_name = "example-resources"
}

output "ledger_endpoint" {
  value = data.azurerm_confidential_ledger.current.ledger_endpoint
}
```

## Argument Reference

* `name` - (Required) Specifies the name of this Confidential Ledger.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where this Confidential Ledger exists.

## Attributes Reference

* `azuread_service_principal` - A list of `azuread_service_principal` blocks as defined below.

* `cert_based_security_principals` - A list of `cert_based_security_principals` blocks as defined below.

* `id` - The ID of this Confidential Ledger.

* `identity_service_endpoint` - The Identity Service Endpoint for this Confidential Ledger.

* `ledger_endpoint` - The Endpoint for this Confidential Ledger.

* `location` - The supported Azure location where the Confidential Ledger exists.

* `ledger_type` - The type of Confidential Ledger.

* `tags` - A mapping of tags to assign to the Confidential Ledger.

---

A `azuread_based_service_principal` block exports the following:

* `ledger_role_name` - The Ledger Role to grant this AzureAD Service Principal.

* `principal_id` - The Principal ID of the AzureAD Service Principal.

* `tenant_id` - The Tenant ID for this AzureAD Service Principal.

---

A `certificate_based_security_principal` block exports the following:

* `ledger_role_name` - The Ledger Role to grant this Certificate Security Principal.

* `pem_public_key` - The public key, in PEM format, of the certificate used by this identity to authenticate with the Confidential Ledger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Confidential Ledger.
