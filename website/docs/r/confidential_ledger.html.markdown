---
subcategory: "Confidential Ledger"
layout: "azurerm"
page_title: "Azure Confidential Ledger: azurerm_confidential_ledger"
description: |-
  Manages a Confidential Ledger.
---

# azurerm_confidential_ledger

Manages an Azure Confidential Ledger.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_confidential_ledger" "ledger" {
  name                = "MyConfidentialLedger"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  ledger_type         = "Public"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Confidential Ledger. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the Confidential Ledger exists.

* `location` - (Required) Specifies the supported Azure location where the Confidential Ledger exists.

* `ledger_type` - (Required) Specifies the type of Confidential Ledger. Possible values are "Public" and "Private".

~> **NOTE:** `ledger_type` cannot be changed after the Confidential Ledger has been created.

* `aad_based_security_principals` - (Optional) An `aadBasedSecurityPrincipal` block as defined below.

* `cert_based_security_principals` - (Optional) A `certBasedSecurityPrincipal` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `aadBasedSecurityPrincipal` block supports the following:

* `principal_id` - (Required) The identifier for the Azure Activate Directory service principal.

* `tenant_id` - (Required) The identifier for the tenant containing the specificed service principal.

* `ledger_role_name` - (Required) The role to assign to the identity. Possible values are "Administrator", "Contributor", and "Reader".

---

A `certBasedSecurityPrincipal` block supports the following:

* `cert` - (Required) The public key, in PEM format, of the certificate used by this identity to authenticate with the Confidential Ledger.

* `ledger_role_name` - (Required) The role to assign to the identity. Possible values are "Administrator", "Contributor", and "Reader".

---
## Attributes Reference

The following attributes are exported:

* `name` - Specifies the name of the Confidential Ledger. Changing this forces a new resource to be created.

* `resource_group_name` - The name of the resource group where the Confidential Ledger exists.

* `location` - Specifies the supported Azure location where the Confidential Ledger exists.

* `ledger_type` - Specifies the type of Confidential Ledger.

* `aad_based_security_principals` - An `aadBasedSecurityPrincipal` block as defined below.

* `cert_based_security_principals` - A `certBasedSecurityPrincipal` block as defined below.

* `tags` - A mapping of tags to assign to the resource.

---

A `aadBasedSecurityPrincipal` block supports the following:

* `principal_id` - The identifier for the Azure Activate Directory service principal.

* `tenant_id` - The identifier for the tenant containing the specificed service principal.

* `ledger_role_name` - The role to assign to the identity.

---

A `certBasedSecurityPrincipal` block supports the following:

* `cert` - The public key, in PEM format, of the certificate used by this identity to authenticate with the Confidential Ledger.

* `ledger_role_name` - The role to assign to the identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Confidential Ledger.
* `update` - (Defaults to 30 minutes) Used when updating the Confidential Ledger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Confidential Ledger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Confidential Ledger.

## Import

Confidential Ledgers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_confidential_ledger.testLedger /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ledgerRG/providers/Microsoft.ConfidentialLedger/Ledgers/testLedger
```
