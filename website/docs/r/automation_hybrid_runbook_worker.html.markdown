---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_hybrid_runbook_worker"
description: |-
  Manages a Automation.
---

# azurerm_automation_hybrid_runbook_worker

Manages a Automation Hybrid Runbook Worker.

## Example Usage

```hcl
resource "azurerm_automation_hybrid_runbook_worker" "example" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  worker_group_name       = azurerm_automation_hybrid_runbook_worker_group.test.name
  vm_resource_id          = azurerm_linux_virtual_machine.test.id
  worker_id               = "00000000-0000-0000-0000-000000000000" #unique uuid
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Automation should exist. Changing this forces a new Automation to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Hybrid Worker is created. Changing this forces a new resource to be created.

* `worker_group_name` - (Required) The name of the HybridWorker Group. Changing this forces a new Automation to be created.

* `worker_id` - (Required) Specify the ID of this HybridWorker in UUID notation. Changing this forces a new Automation to be created.

* `vm_resource_id` - (Required) The ID of the virtual machine used for this HybridWorker. Changing this forces a new Automation to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Hybrid Runbook Worker.

* `ip` - The IP address of assigned machine.

* `last_seen_date_time` - Last Heartbeat from the Worker.

* `registration_date_time` - The registration time of the worker machine.

* `worker_name` - The name of HybridWorker.

* `worker_type` - The type of the HybridWorker, the possible values are `HybridV1` and `HybridV2`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_hybrid_runbook_worker.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/hybridRunbookWorkerGroups/group1/hybridRunbookWorkers/00000000-0000-0000-0000-000000000000
```
