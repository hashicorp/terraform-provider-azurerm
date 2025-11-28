---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_workspace_private_endpoint_connection_approval"
description: |-
  Approves a Private Endpoint Connection for an Azure Monitor Workspace.
---

# azurerm_monitor_workspace_private_endpoint_connection_approval

Approves a Private Endpoint Connection for an Azure Monitor Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_workspace" "example" {
  name                = "example-mamw"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_dashboard_grafana" "example" {
  name                  = "example-grafana"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  grafana_major_version = "10"
}

resource "azurerm_dashboard_grafana_managed_private_endpoint" "example" {
  grafana_id                   = azurerm_dashboard_grafana.example.id
  name                         = "example-mpe"
  location                     = azurerm_dashboard_grafana.example.location
  private_link_resource_id     = azurerm_monitor_workspace.example.id
  group_ids                    = ["prometheusMetrics"]
  private_link_resource_region = azurerm_dashboard_grafana.example.location
}

data "azurerm_monitor_workspace" "example" {
  name                = azurerm_monitor_workspace.example.name
  resource_group_name = azurerm_monitor_workspace.example.resource_group_name

  depends_on = [azurerm_dashboard_grafana_managed_private_endpoint.example]
}

resource "azurerm_monitor_workspace_private_endpoint_connection_approval" "example" {
  workspace_id                     = azurerm_monitor_workspace.example.id
  private_endpoint_connection_name = data.azurerm_monitor_workspace.example.private_endpoint_connections[0].name
}
```

## Arguments Reference

The following arguments are supported:

* `workspace_id` - (Required) The ID of the Azure Monitor Workspace. Changing this forces a new resource to be created.

* `private_endpoint_connection_name` - (Required) The name of the private endpoint connection to approve. Changing this forces a new resource to be created.

* `approval_message` - (Optional) The approval message for the private endpoint connection. Defaults to `Approved via Terraform`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Monitor Workspace.

* `name` - The name of the private endpoint connection.

* `private_endpoint_id` - The ID of the private endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when approving the Private Endpoint Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private Endpoint Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Private Endpoint Connection.
* `delete` - (Defaults to 30 minutes) Used when rejecting the Private Endpoint Connection.

## Import

Azure Monitor Workspace Private Endpoint Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_workspace_private_endpoint_connection_approval.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Monitor/accounts/azureMonitorWorkspace1
```
