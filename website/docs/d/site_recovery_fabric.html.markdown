---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_fabric"
description: |-
  Gets information about an existing Site Recovery Replication Fabric on Azure.
---

# Data Source: azurerm_site_recovery_fabric

Use this data source to access information about an existing Site Recovery Replication Fabric.

## Example Usage

```hcl
data "azurerm_site_recovery_fabric" "fabric" {
  name                = "primary-fabric"
  recovery_vault_name = "tfex-recovery_vault"
  resource_group_name = "tfex-resource_group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Site Recovery Replication Fabric.

* `recovery_vault_name` - (Required) The name of the Recovery Services Vault that the Site Recovery Replication Fabric is associated witth.

* `resource_group_name` - (Required) The name of the resource group in which the associated Recovery Services Vault resides.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Site Recovery Replication Fabric.

* `location` - The Azure location where the Site Recovery Replication Fabric resides.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services Vault.
