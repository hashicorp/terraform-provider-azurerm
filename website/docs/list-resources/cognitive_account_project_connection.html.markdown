---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project_connection"
description: |-
  Lists Cognitive Account Project Connection resources.
---

# List resource: azurerm_cognitive_account_project_connection

Lists Cognitive Account Project Connection resources, with optional filtering by authentication type.

## Example Usage

### List all connections for a specific project

```hcl
list "azurerm_cognitive_account_project_connection" "example" {
  provider = azurerm
  config {
    cognitive_account_name = "example-aiservices"
    project_name           = "example-project"
    resource_group_name    = "example-rg"
  }
}
```

### List connections filtered by authentication type

```hcl
list "azurerm_cognitive_account_project_connection" "example" {
  provider = azurerm
  config {
    cognitive_account_name = "example-aiservices"
    project_name           = "example-project"
    resource_group_name    = "example-rg"
    auth_types             = ["AAD"]
  }
}
```

### List connections filtered by multiple authentication types

```hcl
list "azurerm_cognitive_account_project_connection" "example" {
  provider = azurerm
  config {
    cognitive_account_name = "example-aiservices"
    project_name           = "example-project"
    resource_group_name    = "example-rg"
    auth_types             = ["AAD", "ApiKey", "ManagedIdentity"]
  }
}
```

### List all connections for all projects in an account

```hcl
list "azurerm_cognitive_account_project_connection" "example" {
  provider = azurerm
  config {
    cognitive_account_name = "example-aiservices"
    resource_group_name    = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `auth_types` - (Optional) A list of authentication types to filter by. When specified, only connections matching one of the given authentication types are returned. Possible values include `AAD`, `AccessKey`, `AccountKey`, `AccountManagedIdentity`, `AgentUserImpersonation`, `AgenticIdentityToken`, `AgenticUser`, `ApiKey`, `CustomKeys`, `DelegatedSAS`, `ManagedIdentity`, `None`, `OAuth2`, `PAT`, `ProjectManagedIdentity`, `SAS`, `ServicePrincipal`, `UserEntraToken`, and `UsernamePassword`.

* `cognitive_account_name` - (Required) The name of the Cognitive Services Account.

* `project_name` - (Optional) The name of the Cognitive Services Account Project. If specified, `cognitive_account_name` and `resource_group_name` must also be specified.

* `resource_group_name` - (Required) The name of the resource group containing the Cognitive Services Account.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.

## Attributes Reference

Each returned item includes the following attributes:

* `id` - The ID of the Cognitive Services Account Project Connection.

* `name` - The name of the connection.

* `cognitive_account_project_id` - The ID of the Cognitive Services Account Project.

* `auth_type` - The authentication type of the connection (e.g., `AAD`, `ApiKey`, `ManagedIdentity`).

* `category` - The category of the connection.

* `target` - The target endpoint or resource for the connection.

* `metadata` - A mapping of metadata key-value pairs for the connection.
