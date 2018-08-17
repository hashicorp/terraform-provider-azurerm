---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group"
sidebar_current: "docs-azurerm-management-group"
description: |-
  Manages a Management Group.
---

# azurerm_management_group

Create a management group with subscription assignments.

## Example Usage

```hcl

data "azurerm_subscription" "current" {}

resource "azurerm_management_group" "testmanagementgroup" {
    name = "TestManagementGroup"
    subscription_ids = [
        "${data.azurerm_subscription.current.id}"
    ]
}
```

## Example Usage with azurerm_role_assignment

```hcl
data "azurerm_subscription" "current" {}

data "azurerm_client_config" "test" {}

resource "azurerm_management_group" "testmanagementgroup" {
    name = "TestManagementGroup"
    subscription_ids = [
        "${data.azurerm_subscription.primary.id}"
    ]
}

resource "azurerm_role_assignment" "test" {
  scope                = "${data.azurerm_management_group.testmanagementgroup.id}"
  role_definition_name = "Reader"
  principal_id         = "${data.azurerm_client_config.test.service_principal_object_id}"
}
```

## Example Usage with azurerm_policy_assignment and azurerm_policy_definition

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_management_group" "testmanagementgroup" {
    name = "TestManagementGroup"
    subscription_ids = [
        "${data.azurerm_subscription.current.id}"
    ]
}

resource "azurerm_policy_definition" "test" {
  name         = "my-policy-definition"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%d"
  policy_rule  = <<POLICY_RULE
	{
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
	{
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}

resource "azurerm_policy_assignment" "test" {
  name                 = "example-policy-assignment"
  scope                = "${azurerm_management_group.testmanagementgroup.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "Acceptance Test Run %d"
  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "West Europe" ]
  }
}
PARAMETERS
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name & id of the management group. This needs to be unique across your AAD tenant.

* `subscription_ids` - (Optional) List of subscription IDs to be assigned to the management group.

## Attributes Reference

The following attributes are exported:

* `id` - The management group id.

## Import

Management groups can be imported using the `management group name`, e.g.

```shell
terraform import azurerm_management_group.testManagementGroup  /providers/Microsoft.Management/ManagementGroups/<MANAGEMENT_GROUP_NAME>
```
