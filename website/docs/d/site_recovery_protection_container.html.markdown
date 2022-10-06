---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_protection_container"
description: |-
  Gets information about an existing site recovery services protection container on Azure.
---

# Data Source: azurerm_site_recovery_protection_container

Use this data source to access information about an existing site recovery services protection container.

## Example Usage

```hcl
data "azurerm_site_recovery_protection_container" "container" {
  name                 = "primary-container"
  recovery_vault_name  = "tfex-recovery_vault"
  resource_group_name  = "tfex-resource_group"
  recovery_fabric_name = "primary-fabric"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the protection container.

* `recovery_vault_name` - (Required) The name of the Recovery Services Vault that the protection container is associated witth.

* `resource_group_name` - (Required) The name of the resource group in which the associated protection container resides.

* `recovery_fabric_name` - (Required) The name of the fabric that contains the protection container.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the protection container.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services Vault.
