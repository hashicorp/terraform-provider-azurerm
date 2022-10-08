---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_replication_policy"
description: |-
  Gets information about an existing Azure Site Recovery replication policy.
---

# Data Source: azurerm_site_recovery_replication_policy

Use this data source to access information about an existing Azure Site Recovery replication policy.

## Example Usage

```hcl
data "azurerm_site_recovery_replication_policy" "policy" {
  name                = "replication-policy"
  recovery_vault_name = "tfex-recovery_vault"
  resource_group_name = "tfex-resource_group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure Site Recovery replication policy.

* `recovery_vault_name` - (Required) The name of the Recovery Services Vault that the Azure Site Recovery replication policy is associated witth.

* `resource_group_name` - (Required) The name of the resource group in which the associated Azure Site Recovery replication policy resides.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Site Recovery replication policy.

* `recovery_point_retention_in_minutes` - The duration in minutes for which the recovery points need to be stored.

* `application_consistent_snapshot_frequency_in_minutes` - Specifies the frequency (in minutes) at which to create application consistent recovery.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services Vault.
