---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_set_definition"
sidebar_current: "docs-azurerm-resource-policy-set-definition"
description: |-
  Manages a policy set definition.
---

# azurerm_policy_set_definition

Manages a policy set definition. 

-> **NOTE:**  Policy set definitions (also known as policy initiatives) do not take effect until they are assigned to a scope using a Policy Set Assignment.

## Example Usage

```hcl
resource "azurerm_policy_set_definition" "test" {
  name         = "testPolicySet"
  policy_type  = "Custom"
  display_name = "Test Policy Set"

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

  policy_definitions = <<POLICY_DEFINITIONS
    [
        {
            "parameters": {
                "listOfAllowedLocations": {
                    "value": "[parameters('allowedLocations')]"
                }
            },
            "policyDefinitionId": "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
        }
    ]
POLICY_DEFINITIONS
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the policy set definition. Changing this forces a new resource to be created.

* `policy_type` - (Required) The policy set type. The value can be `BuiltIn`, `Custom` or `NotSpecified`
. Changing this forces a new resource to be created.

* `display_name` - (Required) The display set name of the policy definition.

* `policy_definitions` - (Required) The policy definitions for the policy set definition. This is a json object representing the bundled policy definitions .

* `description` - (Optional) The description of the policy set definition.

* `metadata` - (Optional) The metadata for the policy definition. This is a json object representing additional metadata that should be stored with the policy definition.

* `parameters` - (Optional) Parameters for the policy definition. This field is a json object that allows you to parameterize your policy definition.

## Attributes Reference

The following attributes are exported:

* `id` - The policy set definition id.

## Import

Policy Set Definitions can be imported using the `policy set name`, e.g.

```shell
terraform import azurerm_policy_set_definition.test  /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policySetDefinitions/testPolicySet
```
