---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_replication_recovery_plan"
description: |-
    Manages an Azure Site Recovery Plan within a Recovery Services vault.
---

# azurerm_site_recovery_replication_recovery_plan

Manages an Azure Site Recovery Plan within a Recovery Services vault. A recovery plan gathers machines into recovery groups for the purpose of failover.

## Example Usage

```hcl
resource "azurerm_resource_group" "source" {
  name     = "example-source-rg"
  location = "west us"
}


resource "azurerm_resource_group" "target" {
  name     = "example-target-rg"
  location = "east us"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "example-kv"
  location            = azurerm_resource_group.target.location
  resource_group_name = azurerm_resource_group.target.name
  sku                 = "Standard"
}

resource "azurerm_site_recovery_fabric" "source" {
  resource_group_name = azurerm_resource_group.example.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name
  name                = "example-fabric-source"
  location            = azurerm_resource_group.source.location
}

resource "azurerm_site_recovery_fabric" "target" {
  resource_group_name = azurerm_resource_group.target.name
  recovery_vault_name = azurerm_recovery_services_vault.example.name
  name                = "example-fabric-target"
  location            = azurerm_resource_group.target.location
  depends_on          = [azurerm_site_recovery_fabric.source]
}

resource "azurerm_site_recovery_replication_recovery_plan" "example" {
  name                      = "example-recover-plan"
  recovery_vault_id         = azurerm_recovery_services_vault.target.id
  source_recovery_fabric_id = azurerm_site_recovery_fabric.source.id
  target_recovery_fabric_id = azurerm_site_recovery_fabric.target.id

  recovery_group {
    type                       = "Boot"
    replicated_protected_items = [azurerm_site_recovery_replicated_vm.test.id]
  }
  recovery_group {
    type = "Failover"
  }
  recovery_group {
    type = "Shutdown"
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Replication Plan. The name can contain only letters, numbers, and hyphens. It should start with a letter and end with a letter or a number. Can be a maximum of 63 characters.

* `recovery_vault_id` - (Required) The ID of the vault that should be updated.

* `source_recovery_fabric_id` - (Required) ID of source fabric to be recovered from. Changing this forces a new Replication Plan to be created.

* `target_recovery_fabric_id` - (Required) ID of target fabric to recover. Changing this forces a new Replication Plan to be created.

* `recovery_group` - (Required) Three or more `recovery_group` block.

---

A `recovery_groups` block supports the following:

*  `type` - (Required) The Recovery Plan Group Type. Possible values are `Boot`, `Failover` and `Shutdown`.

* `replicated_protected_items` - (required) one or more id of protected VM.

* `pre_action` - (Optional) one or more `action` block. which will be executed before the group recovery.

* `post_action` - (Optional) one or more `action` block. which will be executed after the group recovery.

---

An `action` block supports the following:

* `name` - (Required) Name of the Action.

* `type` - (Required) Type of the action detail. Possible values are `AutomationRunbookActionDetails`, `ManualActionDetails` and `ScriptActionDetails`.

* `fail_over_directions` - (Required) Directions of fail over. Possible values are `PrimaryToRecovery` and `RecoveryToPrimary`

* `fail_over_types` - (Required) Types of fail over. Possible values are `TestFailover`, `PlannedFailover` and `UnplannedFailover`

* `fabric_location` - (Optional) The fabric location of runbook or script. Possible values are `Primary` and `Recovery`.

-> **NOTE:** This is required when `type` is set to `AutomationRunbookActionDetails` or `ScriptActionDetails`.

* `runbook_id` - (Optional) Id of runbook.

-> **NOTE:** This property is required when `type` is set to `AutomationRunbookActionDetails`.

* `manual_action_instruction` - (Optional) Instructions of manual action.

-> **NOTE:** This property is required when `type` is set to `ManualActionDetails`.

* `script_path` - (Optional) Path of action script.

-> **NOTE:** This property is required when `type` is set to `ScriptActionDetails`.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the Site Recovery Fabric.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Site Recovery Replication Plan.
* `update` - (Defaults to 30 minutes) Used when updating the Site Recovery Replication Plan.
* `read` - (Defaults to 5 minutes) Used when retrieving the Site Recovery Replication Plan.
* `delete` - (Defaults to 30 minutes) Used when deleting the Site Recovery Replication Plan.

## Import

Site Recovery Fabric can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_fabric.myfabric-id=/subscriptions/00000000-0000-0000-0000-00000000000/resourceGroups/groupName/providers/Microsoft.RecoveryServices/vaults/vaultName/replicationRecoveryPlans/planName
```
