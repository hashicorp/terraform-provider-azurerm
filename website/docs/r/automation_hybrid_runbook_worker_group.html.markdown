---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_hybrid_runbook_worker_group"
description: |-
  Manages a Automation Account Runbook Worker Group.
---

# azurerm_automation_hybrid_runbook_worker_group

Manages a Automation Hybrid Runbook Worker Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_hybrid_runbook_worker_group" "example" {
  name                    = "example"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `automation_account_name` - (Required) The name of the Automation Account in which the Runbook Worker Group is created. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Automation Account Runbook Worker Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Automation should exist. Changing this forces a new Automation to be created.

---

* `credential_name` - (Optional) The name of resource type `azurerm_automation_credential` to use for hybrid worker.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Hybrid Runbook Worker Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `update` - (Defaults to 10 minutes) Used when updating the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_hybrid_runbook_worker_group.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/hybridRunbookWorkerGroups/grp1
```
