---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_definition"
description: |-
  Manages a Policy Definition in a Subscription or Management Group.
---

# azurerm_policy_definition

Manages a Policy Definition in a Subscription or Management Group.

When `management_group_id` is omitted, the Policy Definition is created in the Subscription configured in the provider.

Policy definitions do not take effect until they are assigned to a scope using a Policy Assignment.

## Example Usage

```hcl
resource "azurerm_policy_definition" "example" {
  name         = "example-policy-definition"
  policy_type  = "Custom"
  mode         = "Indexed"
  display_name = "Example Policy Definition"

  metadata = <<METADATA
  {
    "category": "General"
  }
METADATA

  policy_rule = <<POLICY_RULE
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Definition. Changing this forces a new resource to be created.

* `management_group_id` - (Optional) The ID of the Management Group where the Policy Definition should exist. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of the Policy Definition.

* `mode` - (Required) The Resource Manager mode that determines which resource types are evaluated by the Policy Definition. Possible values are `All`, `Indexed`, `Microsoft.ContainerService.Data`, `Microsoft.CustomerLockbox.Data`, `Microsoft.DataCatalog.Data`, `Microsoft.KeyVault.Data`, `Microsoft.Kubernetes.Data`, `Microsoft.MachineLearningServices.Data`, `Microsoft.Network.Data`, and `Microsoft.Synapse.Data`.

-> **Note:** For service-specific availability, see [resource provider modes](https://learn.microsoft.com/azure/governance/policy/concepts/definition-structure#resource-provider-modes).

* `policy_type` - (Required) The type of the Policy Definition. Possible values are `BuiltIn`, `Custom`, `NotSpecified`, and `Static`. Changing this forces a new resource to be created.

* `description` - (Optional) The description of the Policy Definition.

* `metadata` - (Optional) A JSON string containing additional metadata for the Policy Definition.

-> **Note:** When omitted, `metadata` is populated from the Policy Definition returned by Azure.

* `parameters` - (Optional) A JSON string defining the parameters for the Policy Definition.

~> **Note:** Reducing the number of parameters or clearing `parameters` forces a new resource to be created.

* `policy_rule` - (Optional) A JSON string defining the Policy Definition rule, containing `if` and `then` blocks.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Policy Definition.

* `role_definition_ids` - A list of Role Definition IDs extracted from `policy_rule` for remediation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Definition.

## Import

A Policy Definition can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_policy_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/policyDefinitions/policyDefinition1
```

A Policy Definition at Management Group scope can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_policy_definition.example /providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.Authorization/policyDefinitions/name
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Authorization` - 2025-01-01

* `Microsoft.Management` - 2025-01-01
