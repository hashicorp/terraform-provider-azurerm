---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_hybrid_runbook_worker_group"
description: |-
  Manages a Automation.
---

# azurerm_automation_hybrid_runbook_worker_group

Manages a Automation Hybrid Runbook Worker Group.

## Example Usage

```hcl
resource "azurerm_automation_hybrid_runbook_worker_group" "example" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  credential_name         = azurerm_automation_credential.test.name
  name = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Automation should exist. Changing this forces a new Automation to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Worker Group is created. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Automation Hybrid Runbook Worker Group. Changing this forces a new Automation to be created.

* `credential_name` - (Required) The RunAs credential to user for Hybrid Worker.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Automation Hybrid Runbook Worker Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `update` - (Defaults to 10 minutes) Used when updating the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_hybrid_runbook_worker_group.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/hybridRunbookWorkerGroups/group1
```