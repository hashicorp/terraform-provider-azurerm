---
subcategory: "Confidential Ledger"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_confidential_ledger"
description: |-
  Manages a Confidential Ledger.
---

# azurerm_confidential_ledger

Manages a Confidential Ledger.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_confidential_ledger" "ledger" {
  name                = "example-ledger"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  ledger_type         = "Private"

  azuread_based_service_principal {
    principal_id     = data.azurerm_client_config.current.object_id
    tenant_id        = data.azurerm_client_config.current.tenant_id
    ledger_role_name = "Administrator"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Confidential Ledger. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Confidential Ledger exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Confidential Ledger exists. Changing this forces a new resource to be created.

* `azuread_based_service_principal` - (Required) A list of `azuread_based_service_principal` blocks as defined below.

* `ledger_type` - (Required) Specifies the type of Confidential Ledger. Possible values are `Private` and `Public`. Changing this forces a new resource to be created.

---

* `certificate_based_security_principal` - (Optional) A list of `certificate_based_security_principal` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Confidential Ledger.

---

A `azuread_based_service_principal` block supports the following:

* `ledger_role_name` - (Required) Specifies the Ledger Role to grant this AzureAD Service Principal. Possible values are `Administrator`, `Contributor` and `Reader`.

* `principal_id` - (Required) Specifies the Principal ID of the AzureAD Service Principal.

* `tenant_id` - (Required) Specifies the Tenant ID for this AzureAD Service Principal.

---

A `certificate_based_security_principal` block supports the following:

* `ledger_role_name` - (Required) Specifies the Ledger Role to grant this Certificate Security Principal. Possible values are `Administrator`, `Contributor` and `Reader`.

* `pem_public_key` - (Required) The public key, in PEM format, of the certificate used by this identity to authenticate with the Confidential Ledger.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this Confidential Ledger.

* `identity_service_endpoint` - The Identity Service Endpoint for this Confidential Ledger.

* `ledger_endpoint` - The Endpoint for this Confidential Ledger.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Confidential Ledger.
* `read` - (Defaults to 5 minutes) Used when retrieving the Confidential Ledger.
* `update` - (Defaults to 30 minutes) Used when updating the Confidential Ledger.
* `delete` - (Defaults to 30 minutes) Used when deleting the Confidential Ledger.

## Import

Confidential Ledgers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_confidential_ledger.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-group/providers/Microsoft.ConfidentialLedger/ledgers/example-ledger
```
