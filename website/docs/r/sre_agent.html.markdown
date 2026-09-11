---
subcategory: "SRE Agent"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sre_agent"
description: |-
  Manages an SRE Agent.
---

# azurerm_sre_agent

Manages an SRE Agent (`Microsoft.App/agents`).

~> **Note:** This local draft is not included in a released AzureRM provider. Its examples and acceptance test definitions have not been exercised against Azure.

## Prerequisites

Register `Microsoft.App` in the target subscription before using this resource. The user-assigned identity example also requires `Microsoft.ManagedIdentity`.

For subscriptions registered separately, set `resource_provider_registrations = "none"` in your provider configuration to disable automatic registration.

Use a subscription and region where SRE Agent is available. The Terraform principal needs permission to manage agents and their resource group, plus any identities and role assignments it creates.

## Example Usage

### Basic

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-sre"
  location = "East US 2"
}

resource "azurerm_sre_agent" "example" {
  name                = "example-sre-agent"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}
```

### Separate identities and an explicit resource scope

This example assigns separate user-assigned managed identities to action and resource configuration. Each identity has an explicit `Reader` role assignment on the monitored resource group.

`Reader` grants access to resource properties. Additional permissions are needed for other operations, including log queries and remediation. See [Agent permissions](https://learn.microsoft.com/azure/sre-agent/permissions) before selecting roles for a workload.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-sre"
  location = "East US 2"
}

resource "azurerm_resource_group" "monitored" {
  name     = "example-monitored"
  location = azurerm_resource_group.example.location
}

resource "azurerm_user_assigned_identity" "actions" {
  name                = "example-sre-actions"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_user_assigned_identity" "resources" {
  name                = "example-sre-resources"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_role_assignment" "actions_reader" {
  scope                = azurerm_resource_group.monitored.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.actions.principal_id
}

resource "azurerm_role_assignment" "resources_reader" {
  scope                = azurerm_resource_group.monitored.id
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.resources.principal_id
}

resource "azurerm_sre_agent" "example" {
  name                = "example-sre-agent"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.actions.id,
      azurerm_user_assigned_identity.resources.id,
    ]
  }

  action_configuration {
    access_level = "Low"
    identity_id  = azurerm_user_assigned_identity.actions.id
    mode         = "Review"
  }

  resources_configuration {
    identity_id  = azurerm_user_assigned_identity.resources.id
    resource_ids = [azurerm_resource_group.monitored.id]
  }

  tags = {
    environment = "example"
  }

  depends_on = [
    azurerm_role_assignment.actions_reader,
    azurerm_role_assignment.resources_reader,
  ]
}
```

The root `identity` block attaches managed identities to the agent; each nested `identity_id` selects one for that configuration. These references remain separate in Terraform state. The provider does not create identities or grant role-based access control (RBAC) permissions implicitly.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the SRE Agent. Must contain 2 to 32 letters, digits, or hyphens, start with a letter, and end with a letter or digit. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SRE Agent should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SRE Agent should exist. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below. At least one managed identity must remain assigned to the agent.

---

* `action_configuration` - (Optional, Computed) An `action_configuration` block as defined below. When omitted, the provider reads the service value into state without changing it.

* `resources_configuration` - (Optional, Computed) A `resources_configuration` block as defined below. When omitted, the provider reads the service value into state without changing it.

* `tags` - (Optional) A mapping of tags to assign to the SRE Agent.

---

An `action_configuration` block supports the following:

* `access_level` - (Required) The action access level. Possible values are `Low` and `High`. Setting this value does not create an RBAC role assignment.

* `identity_id` - (Required) The Azure resource ID of the managed identity used for actions.

* `mode` - (Required) The action mode. Possible values are `Review` and `Autonomous`.

---

An `identity` block supports the following:

* `type` - (Required) The type of managed identity assigned to the SRE Agent. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) A set of user-assigned managed identity resource IDs to attach to the SRE Agent. Required when `type` includes `UserAssigned`.

---

A `resources_configuration` block supports the following:

* `identity_id` - (Required) The Azure resource ID of the identity for resource configuration. Maps to `knowledgeGraphConfiguration.identity`; this reference does not grant permissions.

* `resource_ids` - (Required) A nonempty set of Azure resource IDs managed by the agent. The provider maps this set to `knowledgeGraphConfiguration.managedResources`. Resource group IDs can define the monitored scope.

~> **Note:** This draft requires all inner fields when either configuration block is supplied. The ARM API permits a broader set of configurations.

## Local Draft Limitations

* This draft exposes no virtual network (VNet) arguments. Public networking examples use fields absent from the [stable Agent schema](https://learn.microsoft.com/azure/templates/microsoft.app/2026-01-01/agents). An experimental request adapter is technically possible; upstream support requires public contract evidence and maintainer agreement.

* Removing an `action_configuration` or `resources_configuration` block leaves its remote value unchanged. The provider continues to read that value into state. A dynamic block with no instances has the same behavior.

* These settings use nested-block syntax; assignments such as `action_configuration = []` or `resources_configuration = []` are unsupported. This draft has no block-reset operation because the stable API reset contract has not been established.

* `resource_ids` must contain at least one ID. An empty set is invalid and is never treated as a request to remove resource scope.

* The public API includes `ReadOnly`, which this draft preserves during read, import and discovery. Only `Review` and `Autonomous` can be configured. Review generated configuration before applying it.

* `default_model` is computed-only. This draft has no model-selection setter. The endpoint and runtime status attributes are also computed-only.

* AgentSpaces, connectors, and implicit RBAC assignments are outside this resource's scope.

* At least one managed identity must remain assigned. Before removing a user-assigned identity, update any action or resource configuration that references it. Identity updates remove the assignments omitted from the desired identity set.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SRE Agent.

* `default_model` - A `default_model` block as defined below, when returned by the service.

* `endpoint` - The agent endpoint returned by the service.

* `power_state` - The power state returned by the service.

* `provisioning_state` - The provisioning state returned by the service.

* `running_state` - The running state returned by the service.

---

A `default_model` block exports the following:

* `name` - The default model name returned by the service.

* `provider` - The default model provider returned by the service.

---

An `identity` block exports the following:

* `principal_id` - The principal ID of the system-assigned managed identity, when configured.

* `tenant_id` - The tenant ID of the system-assigned managed identity, when configured.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SRE Agent.
* `read` - (Defaults to 5 minutes) Used when retrieving the SRE Agent.
* `update` - (Defaults to 30 minutes) Used when updating the SRE Agent.
* `delete` - (Defaults to 30 minutes) Used when deleting the SRE Agent.

## Import

SRE Agents can be imported using the `resource id`, for example:

```shell
terraform import azurerm_sre_agent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-sre/providers/Microsoft.App/agents/example-sre-agent
```

Import reads the existing resource. It does not configure unsupported fields or establish that the imported configuration can be applied unchanged.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2026-01-01
