---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project_connection_entra_id"
description: |-
  Lists Cognitive Account Project Connection (Entra ID) resources.
---

# List resource: azurerm_cognitive_account_project_connection_entra_id

Lists Cognitive Account Project Connection resources that use Entra ID (formerly Azure Active Directory) authentication.

## Example Usage

### List all Entra ID connections for a specific project

```hcl
list "azurerm_cognitive_account_project_connection_entra_id" "example" {
  provider = azurerm
  config {
    cognitive_account_name = "example-aiservices"
    project_name           = "example-project"
    resource_group_name    = "example-rg"
  }
}
```

### List all Entra ID connections for all projects in an account

```hcl
list "azurerm_cognitive_account_project_connection_entra_id" "example" {
  provider = azurerm
  config {
    cognitive_account_name = "example-aiservices"
    resource_group_name    = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cognitive_account_name` - (Required) The name of the Cognitive Services Account.

* `project_name` - (Optional) The name of the Cognitive Services Account Project. If specified, `cognitive_account_name` and `resource_group_name` must also be specified.

* `resource_group_name` - (Required) The name of the resource group containing the Cognitive Services Account.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.

## Attributes Reference

Each returned item includes the following attributes:

* `id` - The ID of the Cognitive Services Account Project Connection.

* `name` - The name of the connection.

* `cognitive_account_project_id` - The ID of the Cognitive Services Account Project.

* `category` - The category of the connection.

* `target` - The target endpoint or resource for the connection.

* `metadata` - A mapping of metadata key-value pairs for the connection.
