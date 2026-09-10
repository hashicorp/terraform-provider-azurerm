---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_elastic_san_volume_group"
description: |-
  Manages a Data Protection Backup Policy for Elastic SAN Volume Groups.
---

# azurerm_data_protection_backup_policy_elastic_san_volume_group

Manages a Data Protection Backup Policy for Elastic SAN Volume Groups.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-data-protection-backup-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

resource "azurerm_data_protection_backup_policy_elastic_san_volume_group" "example" {
  name                            = "example-data-protection-backup-policy-elastic-san-volume-group"
  vault_id                        = azurerm_data_protection_backup_vault.example.id
  backup_repeating_time_intervals = ["R/2024-02-08T13:00:00+00:00/P1D"]

  default_retention_rule {
    life_cycle {
      data_store_type = "OperationalStore"
      duration        = "P7D"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Protection Backup Policy for Elastic SAN Volume Groups. Changing this forces a new resource to be created.

-> **Note:** The `name` must be between 3 and 150 characters, contain only letters, numbers, and hyphens, and start with a letter.

* `vault_id` - (Required) The ID of the Data Protection Backup Vault where the policy should exist. Changing this forces a new resource to be created.

* `backup_repeating_time_intervals` - (Required) A list of repeating backup intervals. Changing this forces a new resource to be created.

-> **Note:** Each interval must use ISO 8601 repeating time interval format.

* `default_retention_rule` - (Required) A `default_retention_rule` block as defined below. Changing this forces a new resource to be created.

* `retention_rule` - (Optional) One or more `retention_rule` blocks as defined below. Changing this forces a new resource to be created.

* `time_zone` - (Optional) The time zone used by the backup schedule. Changing this forces a new resource to be created.

---

A `criteria` block supports the following:

* `absolute_criteria` - (Optional) The absolute criterion used to identify retained backups. Changing this forces a new resource to be created.

-> **Note:** Possible values are `AllBackup`, `FirstOfDay`, `FirstOfWeek`, `FirstOfMonth`, and `FirstOfYear`.

* `days_of_week` - (Optional) The days of the week on which backups are retained. Changing this forces a new resource to be created.

-> **Note:** Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, and `Sunday`.

* `months_of_year` - (Optional) The months in which backups are retained. Changing this forces a new resource to be created.

-> **Note:** Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November`, and `December`.

* `scheduled_backup_times` - (Optional) A list of scheduled backup times in RFC3339 format. Changing this forces a new resource to be created.

* `weeks_of_month` - (Optional) The weeks of the month in which backups are retained. Changing this forces a new resource to be created.

-> **Note:** Possible values are `First`, `Second`, `Third`, `Fourth`, and `Last`.

---

A `default_retention_rule` block supports the following:

* `life_cycle` - (Required) One or more `life_cycle` blocks as defined below. Changing this forces a new resource to be created.

---

A `life_cycle` block supports the following:

* `data_store_type` - (Required) The type of data store. Changing this forces a new resource to be created.

-> **Note:** The only possible value is `OperationalStore`.

* `duration` - (Required) The duration for which backups are retained. Changing this forces a new resource to be created.

-> **Note:** The duration must use ISO 8601 duration format.

---

A `retention_rule` block supports the following:

* `criteria` - (Required) A `criteria` block as defined above. Changing this forces a new resource to be created.

* `life_cycle` - (Required) One or more `life_cycle` blocks as defined above. Changing this forces a new resource to be created.

* `name` - (Required) The name of the retention rule. Changing this forces a new resource to be created.

* `priority` - (Required) The priority of the retention rule. Lower values have higher priority. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Protection Backup Policy for Elastic SAN Volume Groups.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Protection Backup Policy for Elastic SAN Volume Groups.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Protection Backup Policy for Elastic SAN Volume Groups.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Protection Backup Policy for Elastic SAN Volume Groups.

## Import

A Data Protection Backup Policy for Elastic SAN Volume Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_elastic_san_volume_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DataProtection/backupVaults/backupVault1/backupPolicies/backupPolicy1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2026-06-01
