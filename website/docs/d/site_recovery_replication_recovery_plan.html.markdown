---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_replication_recovery_plan"
description: |-
    Get information about an Azure Site Recovery Plan within a Recovery Services vault.
---

# azurerm_site_recovery_replication_recovery_plan

Get information about an Azure Site Recovery Plan within a Recovery Services vault. A recovery plan gathers machines into recovery groups for the purpose of failover.

## Example Usage

```hcl
data "azurerm_recovery_services_vault" "vault" {
  name                = "tfex-recovery_vault"
  resource_group_name = "tfex-resource_group"
}

data "azurerm_site_recovery_replication_recovery_plan" "example" {
  name              = "example-recovery-plan"
  recovery_vault_id = data.azurerm_recovery_services_vault.vault.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Replication Plan.

* `recovery_vault_id` - (Required) The ID of the vault that should be updated.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Fabric.

* `source_recovery_fabric_id` - The ID of source fabric to be recovered from.

* `target_recovery_fabric_id` - The ID of target fabric to recover. 

* `recovery_group` - `recovery_group` block defined as below.
---

A `recovery_groups` block supports the following:

*  `type` - The Recovery Plan Group Type. Possible values are `Boot`, `Failover` and `Shutdown`.

* `replicated_protected_items` - one or more id of protected VM.

* `pre_action` - one or more `action` block. which will be executed before the group recovery.

* `post_action` - one or more `action` block. which will be executed after the group recovery.

---

An `action` block supports the following:

* `name` - Name of the Action.

* `type` - Type of the action detail. 

* `fail_over_directions` - Directions of fail over.

* `fail_over_types` - Types of fail over. 

* `fabric_location` - The fabric location of runbook or script. 

* `runbook_id` - Id of runbook.

* `manual_action_instruction` - Instructions of manual action.

* `script_path` - Path of action script.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Replication Plan.
