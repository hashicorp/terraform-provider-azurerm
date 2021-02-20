---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_azure_active_directory"
description: |-
  Manages a Azure Active Directory Data Connector.
---

# azurerm_sentinel_data_connector_azure_active_directory

Manages a Azure Active Directory Data Connector.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_data_connector_azure_active_directory" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

- `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this Azure Active Directory Data Connector resides in. Changing this forces a new Azure Active Directory Data Connector to be created.

* `name` - (Required) The name which should be used for this Azure Active Directory Data Connector. Changing this forces a new Azure Active Directory Data Connector to be created.

---

- `tenant_id` - (Optional) The ID of the tenant that this Azure Active Directory Data Connector connects to. Changing this forces a new Azure Active Directory Data Connector to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Azure Active Directory Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `create` - (Defaults to 30 minutes) Used when creating the Azure Active Directory Data Connector.
- `read` - (Defaults to 5 minutes) Used when retrieving the Azure Active Directory Data Connector.
- `delete` - (Defaults to 30 minutes) Used when deleting the Azure Active Directory Data Connector.

## Import

Azure Active Directory Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_azure_active_directory.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
