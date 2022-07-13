---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_private_endpoint_connection"
description: |-
  Manages a Automation Private Endpoint Connection.
---

# azurerm_automation_private_endpoint_connection

Manages a Automation Private Endpoint Connection.

## Example Usage

```hcl
resource "azurerm_automation_private_endpoint_connection" "example" {
  name                    = azurerm_automation_account.test.private_endpoint_connection[0].name
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  link_status             = "Rejected"
  link_description        = "approved 2"
}
```

## Arguments Reference

The following arguments are supported:

* `automation_account_name` - (Required) TODO. Changing this forces a new Automation to be created.

* `name` - (Required) The name which should be used for this Automation Private Endpoint Connection. Changing this forces a new Automation to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Automation should exist. Changing this forces a new Automation to be created.

* `link_status` - (Required) The Operation on this connection, Possible values are `Aprroved` and `Rejected`.

---

* `link_description` - (Optional) Description of this operation.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the AutomationPrivate Endpoint Connection.

* `group_ids` - A list of stirng of the type of sub-resource your private endpoint will be able to access.

* `link_action_required` - Any action that is required beyond basic workflow (approve/ reject/ disconnect).

* `private_endpoint_id` - The ID of the the Private Endpoint of this connection belongs to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `update` - (Defaults to 10 minutes) Used when updating the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_private_endpoint_connection.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/privateEndpointConnections/uuid
```