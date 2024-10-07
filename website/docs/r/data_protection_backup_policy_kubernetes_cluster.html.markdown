---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_policy_kubernetes_cluster"
description: |-
  Manages a Backup Policy to back up Kubernetes Cluster.
---

# azurerm_data_protection_backup_policy_kubernetes_cluster

Manages a Backup Policy to back up Kubernetes Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-backup-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}

resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "example" {
  name                = "example-backup-policy"
  resource_group_name = azurerm_resource_group.example.name
  vault_name          = azurerm_data_protection_backup_vault.example.name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  time_zone                       = "India Standard Time"
  default_retention_duration      = "P4M"

  retention_rule {
    name     = "Daily"
    priority = 25

    life_cycle {
      duration        = "P84D"
      data_store_type = "OperationalStore"
    }

    criteria {
      absolute_criteria = "FirstOfDay"
    }
  }

  default_retention_rule {
    life_cycle {
      duration        = "P7D"
      data_store_type = "OperationalStore"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the Backup Policy Kubernetes Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Backup Policy Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `vault_name` - (Required) The name of the Backup Vault where the Backup Policy Kubernetes Cluster should exist. Changing this forces a new resource to be created.

* `backup_repeating_time_intervals` - (Required) Specifies a list of repeating time interval. It supports weekly back. It should follow `ISO 8601` repeating time interval. Changing this forces a new resource to be created.

* `default_retention_rule` - (Required) A `default_retention_rule` block as defined below. Changing this forces a new resource to be created.

* `retention_rule` - (Optional) One or more `retention_rule` blocks as defined below. Changing this forces a new resource to be created.

* `time_zone` - (Optional) Specifies the Time Zone which should be used by the backup schedule. Changing this forces a new resource to be created.

---

A `default_retention_rule` block supports the following:

* `life_cycle` - (Required) A `life_cycle` block as defined below. Changing this forces a new resource to be created.

---

A `retention_rule` block supports the following:

* `name` - (Required) The name which should be used for this retention rule. Changing this forces a new resource to be created.

* `criteria` - (Required) A `criteria` block as defined below. Changing this forces a new resource to be created.

* `life_cycle` - (Required) A `life_cycle` block as defined below. Changing this forces a new resource to be created.

* `priority` - (Required) Specifies the priority of the rule. The priority number must be unique for each rule. The lower the priority number, the higher the priority of the rule. Changing this forces a new resource to be created.

---

A `criteria` block supports the following:

* `absolute_criteria` - (Optional) Possible values are `AllBackup`, `FirstOfDay`, `FirstOfWeek`, `FirstOfMonth` and `FirstOfYear`. These values mean the first successful backup of the day/week/month/year. Changing this forces a new resource to be created.

* `days_of_week` - (Optional) Possible values are `Monday`, `Tuesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`. Changing this forces a new resource to be created.

* `months_of_year` - (Optional) Possible values are `January`, `February`, `March`, `April`, `May`, `June`, `July`, `August`, `September`, `October`, `November` and `December`. Changing this forces a new resource to be created.

* `scheduled_backup_times` - (Optional) Specifies a list of backup times for backup in the `RFC3339` format. Changing this forces a new resource to be created.

* `weeks_of_month` - (Optional) Possible values are `First`, `Second`, `Third`, `Fourth` and `Last`. Changing this forces a new resource to be created.

---

A `life_cycle` block supports the following:

* `data_store_type` - (Required) The type of data store. The only possible value is `OperationalStore`. Changing this forces a new resource to be created.

* `duration` - (Required) The retention duration up to which the backups are to be retained in the data stores. It should follow `ISO 8601` duration format. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Policy Kubernetes Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Policy Kubernetes Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Policy Kubernetes Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Policy Kubernetes Cluster.

## Import

Backup Policy Kubernetes Cluster's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_policy_kubernetes_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1/backupPolicies/backupPolicy1
```
