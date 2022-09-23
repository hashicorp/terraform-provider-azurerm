---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_recover_plan"
description: |-
    Manages a Site Recovery Replication Recover Plan on Azure.
---

# azurerm_site_recovery_recover_plan

Manages a Azure Site Recovery Plan within a Recovery Services vault.A recovery plan gathers machines into recovery groups for the purpose of failover.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Recover Plan.

* `resource_group_name` - (Required) Name of the resource group where the vault that should be updated is located.

* `recovery_vault_name` - (Required) The name of the vault that should be updated.

* `source_recovery_fabric_id` - (Required) ID of source fabric to be recovered from. Changing this forces a new Recover Plan to be created.

* `target_recovery_fabric_id` - (Required) ID of target fabric to recover. Changing this forces a new Recover Plan to be created.

* `failover_deployment_model` - (Required) the deployment model of failover. Possible values are `Classic` and `ResourceManager`.  Changing this forces a new Recover Plan to be created.

* `recovery_groups` - (Required) One or more `recovery_groups` block. 

---

A `recovery_groups` block supports the following:

*  `group_type` - (Required) The Recovery Plan Group Type. Possible values are `Boot`, `Failover` and `Shutdown`.

* `replicated_protected_items` - (required) one or more id of protected VM.

* `pre_actions` - (Optional) one or more `action` block. which will be executed before the group recovery.

* `post_actions` - (Optional) one or more `action` block. which will be executed after the group recovery.

---

A `action` block supports the following:

* `name` - (Required) Name of the Action.

* `action_detail_type` - (Required) Type of the action detail. Possible values are `AutomationRunbookActionDetails`, `ManualActionDetails` and `ScriptActionDetails`.

* `fail_over_directions` - (Required) Directions of fail over. Must be one of `PrimaryToRecovery` or `RecoveryToPrimary`

* `fail_over_types` - (Required) Types of fail over. Possible values are `TestFailover`, `PlannedFailover` and `UnplannedFailover`

* `fabric_location` - (Optional) The fabric location of runbook or script. Must be one of `Primary` or `Recovery`.

-> **NOTE:** This is required when `action_detail_type` is set to `AutomationRunbookActionDetails` or `ScriptActionDetails`.

* `runbook_id` - (Optional) Id of runbook.

-> **NOTE:** This is required when `action_detail_type` is set to `AutomationRunbookActionDetails`.

* `manual_action_instruction` - (Optional) Instructions of manual action.

-> **NOTE:** This is required when `action_detail_type` is set to `ManualActionDetails`.

* `script_path` - (Optional) Path of action script.

-> **NOTE:** This is required when `action_detail_type` is set to `ScriptActionDetails`.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Fabric.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Site Recovery Fabric.
* `update` - (Defaults to 30 minutes) Used when updating the Site Recovery Fabric.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Fabric.
* `delete` - (Defaults to 30 minutes) Used when deleting the Site Recovery Fabric.

## Import

Site Recovery Fabric can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_fabric.myfabric /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resource-group-name/providers/Microsoft.RecoveryServices/vaults/recovery-vault-name/replicationFabrics/fabric-name
```
