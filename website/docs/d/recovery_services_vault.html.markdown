---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_vault"
description: |-
  Gets information about an existing Recovery Services Vault.
---

# Data Source: azurerm_recovery_services_vault

Use this data source to access information about an existing Recovery Services Vault.

## Example Usage

```hcl
data "azurerm_recovery_services_vault" "vault" {
  name                = "tfex-recovery_vault"
  resource_group_name = "tfex-resource_group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Recovery Services Vault.

* `resource_group_name` - The name of the resource group in which the Recovery Services Vault resides.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Recovery Services Vault.

* `location` - The Azure location where the resource resides.

* `tags` - A mapping of tags assigned to the resource.

* `sku` - The vault's current SKU.

* `identity` - (Optional) An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

* `type` - The identity type of this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services Vault.
