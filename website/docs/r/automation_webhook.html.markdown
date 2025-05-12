---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_webhook"
description: |-
  Manages an Automation Runbook's Webhook.
---

# azurerm_automation_webhook

Manages an Automation Runbook's Webhook.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Basic"
}

resource "azurerm_automation_runbook" "example" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  log_verbose             = "true"
  log_progress            = "true"
  description             = "This is an example runbook"
  runbook_type            = "PowerShellWorkflow"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }
}

resource "azurerm_automation_webhook" "example" {
  name                    = "TestRunbook_webhook"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  expiry_time             = "2021-12-31T00:00:00Z"
  enabled                 = true
  runbook_name            = azurerm_automation_runbook.example.name
  parameters = {
    input = "parameter"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Webhook. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Webhook is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Webhook is created. Changing this forces a new resource to be created.

* `expiry_time` - (Required) Timestamp when the webhook expires. Changing this forces a new resource to be created.

* `enabled` - (Optional) Controls if Webhook is enabled. Defaults to `true`.

* `runbook_name` - (Required) Name of the Automation Runbook to execute by Webhook.

* `run_on_worker_group` - (Optional) Name of the hybrid worker group the Webhook job will run on.

* `parameters` - (Optional) Map of input parameters passed to runbook.

* `uri` - (Optional) URI to initiate the webhook. Can be generated using [Generate URI API](https://docs.microsoft.com/rest/api/automation/webhook/generate-uri). By default, new URI is generated on each new resource creation. Changing this forces a new resource to be created. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Automation Webhook ID.

* `uri` - (Sensitive) Generated URI for this Webhook. Changing this forces a new resource to be created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Webhook.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Webhook.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Webhook.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Webhook.

## Import

Automation Webhooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_webhook.TestRunbook_webhook /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/webHooks/TestRunbook_webhook
```
