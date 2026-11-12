---
subcategory: "SRE Agents"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sre_agent"
description: |-
  Manages an Azure SRE Agent.
---

# azurerm_sre_agent

Manages an Azure SRE Agent.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-user-assigned-identity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-log-analytics-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_application_insights" "example" {
  name                = "example-application-insights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_sre_agent" "example" {
  name                = "example-sre-agent"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  upgrade_channel     = "Stable"

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  action {
    access_level              = "Low"
    mode                      = "Review"
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }

  application_insights {
    application_id    = azurerm_application_insights.example.app_id
    connection_string = azurerm_application_insights.example.connection_string
  }

  sandbox {
    egress {
      mode               = "Limited"
      allowed_hosts      = ["*.example.com"]
      allowed_registries = ["pypi"]
    }

    package {
      name            = "requests"
      package_manager = "pip"
      version         = "2.32.5"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure SRE Agent. Must be between 2 and 32 characters, start with a letter, end with a letter or digit, and contain only letters, digits, and hyphens. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure SRE Agent should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure SRE Agent should exist. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

* `action` - (Optional) An `action` block as defined below.

* `application_insights` - (Optional) An `application_insights` block as defined below.

* `default_model` - (Optional) A `default_model` block as defined below.

-> **Note:** When `default_model` is omitted during creation, Azure selects the model and its provider.

* `incident_management` - (Optional) An `incident_management` block as defined below.

* `initial_sponsor_group_id` - (Optional) The object ID (UUID) of the Microsoft Entra group that initially sponsors the agent identity.

~> **Note:** Use a Microsoft 365 group or a dynamic membership group as the sponsor. Fixed-membership security groups are not supported for new sponsor assignments.

* `knowledge_graph` - (Optional) A `knowledge_graph` block as defined below.

* `sandbox` - (Optional) A `sandbox` block as defined below.

~> **Note:** A configured `sandbox` block must specify `egress` or `package`. When omitted during creation, Azure supplies the sandbox configuration.

* `subnet_id` - (Optional) The ID of the Subnet used by the Azure SRE Agent.

~> **Note:** The Subnet must be delegated to `Microsoft.App/environments`.

~> **Note:** `subnet_id` and `sandbox.egress.mode` set to `AzureVNet` must be configured together. When removing `subnet_id`, also remove the `AzureVNet` egress mode or change it to `Limited` or `Unrestricted`.

* `upgrade_channel` - (Optional) The upgrade channel for the Azure SRE Agent. Possible values are `Preview` and `Stable`. Defaults to `Stable`.

* `tags` - (Optional) A mapping of tags to assign to the Azure SRE Agent.

---

An `action` block supports the following:

* `user_assigned_identity_id` - (Required) The ID of the User Assigned Identity used for agent actions.

* `access_level` - (Optional) The permission level for agent actions. Possible values are `High` (Privileged, allowing actions using the User Assigned Identity) and `Low` (Reader, intended for read-only diagnostics).

-> **Note:** Resource access depends on the RBAC roles assigned to the User Assigned Identity. For more information about permission levels and authorization, see [agent permissions](https://learn.microsoft.com/azure/sre-agent/permissions).

* `mode` - (Optional) The action execution mode. Possible values are `Autonomous` and `Review`.

-> **Note:** `Review` requires approval before Azure infrastructure write actions, while `Autonomous` executes them without waiting for approval. These modes control approvals independently of the Reader or Privileged permission level. For more information, see [run modes](https://learn.microsoft.com/azure/sre-agent/run-modes).

---

An `application_insights` block supports the following:

* `application_id` - (Required) The application UUID of the Application Insights component.

* `connection_string` - (Required) The connection string of the Application Insights component.

---

A `default_model` block supports the following:

* `name` - (Required) The name of the model used by the Azure SRE Agent.

* `provider` - (Required) The provider of the model used by the Azure SRE Agent.

-> **Note:** For available providers and model families, see [model provider selection](https://learn.microsoft.com/azure/sre-agent/model-provider-selection#available-providers).

---

An `egress` block supports the following:

* `allowed_code_repositories` - (Optional) A list of allowed code repository providers. Possible values are `AzureDevOps`, `GitLab`, `Github`, and `GithubEnterprise`.

* `allowed_hosts` - (Optional) A list of up to `100` allowed DNS hostnames or leading wildcard hostname patterns. Values must be between 1 and 253 characters, with labels of 1 to 63 letters, digits, or hyphens that start and end with a letter or digit.

~> **Note:** A wildcard, such as `*.example.com`, must be followed by at least two hostname labels. A bare `*` is not accepted.

* `allowed_registries` - (Optional) A list of allowed package registries. Possible values are `npmjs`, `nuget-org`, `pypi`, and `ubuntu`.

* `http_mcp_server_network_access_enabled` - (Optional) Whether the sandbox can access remote Model Context Protocol (MCP) servers over HTTP through the managed network. Defaults to `false`.

* `http_proxy` - (Optional) An `http_proxy` block as defined below.

* `mode` - (Optional) The sandbox egress mode. Possible values are `AzureVNet`, `Limited`, and `Unrestricted`.

---

An `http_proxy` block supports the following:

* `https_url` - (Required) The proxy URL for outbound HTTPS traffic. Must start with `http://` or `https://`.

* `http_url` - (Optional) The proxy URL for outbound HTTP traffic. Must start with `http://` or `https://`.

* `no_proxy` - (Optional) A list of hosts or addresses that bypass the proxy.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity assigned to the Azure SRE Agent. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) A set of User Assigned Identity IDs assigned to the Azure SRE Agent.

~> **Note:** `identity_ids` can only be set when `type` is `UserAssigned` or `SystemAssigned, UserAssigned`.

~> **Note:** The identities specified in `action.user_assigned_identity_id` and `knowledge_graph.user_assigned_identity_id` must be included in `identity_ids`.

---

An `incident_management` block supports the following:

* `type` - (Required) The type of the incident management integration.

* `connection_key` - (Optional) The key used to authenticate the incident management connection.

* `connection_name` - (Optional) The name of the incident management connection.

* `connection_url` - (Optional) The URL of the incident management connection. Must start with `http://` or `https://`.

* `obo_user` - (Optional) The user on whose behalf the incident management connection operates.

---

A `knowledge_graph` block supports the following:

* `user_assigned_identity_id` - (Required) The ID of the User Assigned Identity used for the knowledge graph.

* `managed_resource_ids` - (Optional) A list of Azure Resource Manager IDs identifying the resources managed by the Azure SRE Agent.

---

A `package` block supports the following:

* `name` - (Required) The package name. Must be between 1 and 214 characters, contain letters, digits, periods, underscores, or hyphens, and may include a leading `@` and a single `/`-separated scope. Each component must start with a letter or digit.

* `package_manager` - (Required) The package manager. Possible values are `npm`, `nuget`, and `pip`.

* `version` - (Optional) The package version or version constraint. Must be between 1 and 64 characters and contain only letters, digits, or the characters `._+-~^<>=,*[]()`.

---

A `sandbox` block supports the following:

* `egress` - (Optional) An `egress` block as defined above.

~> **Note:** A configured `egress` block must specify at least one of `allowed_code_repositories`, `allowed_hosts`, `allowed_registries`, `http_mcp_server_network_access_enabled`, `http_proxy`, or `mode`.

* `package` - (Optional) One or more `package` blocks as defined above. A maximum of `200` blocks can be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure SRE Agent.

* `endpoint` - The endpoint of the Azure SRE Agent.

* `identity` - An `identity` block as defined below.

* `outbound_ip_addresses` - A list of outbound IP addresses of the Azure SRE Agent.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID of the System Assigned Identity.

* `tenant_id` - The Tenant ID of the System Assigned Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure SRE Agent.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure SRE Agent.
* `update` - (Defaults to 30 minutes) Used when updating the Azure SRE Agent.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure SRE Agent.

## Import

An Azure SRE Agent can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sre_agent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.App/agents/agent1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2026-07-01
