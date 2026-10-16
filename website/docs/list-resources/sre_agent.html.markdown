---
subcategory: "SRE Agent"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sre_agent"
description: |-
  Lists SRE Agent resources.
---

# List resource: azurerm_sre_agent

Lists SRE Agent resources (`Microsoft.App/agents`) in a subscription or resource group.

~> **Note:** This local draft is not included in a released AzureRM provider. Azure acceptance test definitions have not been run.

List queries require Terraform 1.14 or later.

Register `Microsoft.App` in the subscription and grant the caller permission to list agents in the selected scope. For subscriptions registered separately, set `resource_provider_registrations = "none"` in your provider configuration.

## Example Usage

### List all SRE Agents in the provider subscription

```hcl
list "azurerm_sre_agent" "example" {
  provider = azurerm
  config {}
}
```

### List all SRE Agents in a specific subscription

```hcl
list "azurerm_sre_agent" "example" {
  provider = azurerm
  config {
    subscription_id = "00000000-0000-0000-0000-000000000000"
  }
}
```

### Include resource state for agents in a resource group

```hcl
list "azurerm_sre_agent" "example" {
  provider         = azurerm
  include_resource = true
  config {
    resource_group_name = "example-sre"
  }
}
```

## Argument Reference

The `config` block supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query. When omitted, the query lists agents across the selected subscription.

* `subscription_id` - (Optional) The subscription ID to query. Defaults to the value in the provider configuration. When both arguments are supplied, the resource group is queried in this subscription.

Set the Terraform list argument `include_resource = true` to include resource state with each result.

## Results

Each result includes the agent name as its display name and a Terraform resource identity containing `subscription_id`, `resource_group_name`, and `name`.

With `include_resource = true`, results use the same state mapping as [`azurerm_sre_agent`](/docs/providers/azurerm/r/sre_agent.html). This includes managed identities, action and resource configuration, experimental networking, tags, and computed attributes.

Discovery leaves the agent and its permissions unchanged. Results preserve returned `ReadOnly` modes and empty resource scopes, although this draft cannot configure those values. Review generated configuration before applying it.

Omitting an action or resource configuration block preserves its remote value, including when a dynamic block produces no instances. Empty-list assignments are unsupported nested-block syntax; this draft has no block-reset operation.

The resource's other limitations also apply. Discovery uses the proposed stable `2026-10-01` API and a separately generated SDK. This version has not been deployed or live-validated by this contribution; earlier local `2026-01-01` results are historical only. DNS-selector support is deferred; API-returned DNS settings are not exposed as Terraform controls or copied into generated configuration. Returned mode and subnet settings do not prove runtime routing or DNS behavior. Omitting the networking block preserves the remote settings; detachment requires an explicit non-VNet mode and an omitted subnet ID. Default-model setters, AgentSpaces and connectors remain outside this draft's scope.
