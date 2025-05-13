---
subcategory: "Dashboard"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dashboard_grafana_managed_private_endpoint"
description: |-
  Manages a Dashboard Grafana Managed Private Endpoint.
---

# azurerm_dashboard_grafana_managed_private_endpoint

Manages a Dashboard Grafana Managed Private Endpoint.

~> **Note:** This resource will _not_ approve the managed private endpoint connection on the linked resource. This will need to be done manually via Azure CLI, PowerShell, or AzAPI resources. See [here](https://github.com/hashicorp/terraform-provider-azurerm/issues/23950#issuecomment-2035109970) for an example that uses AzAPI.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Canada Central"
}

resource "azurerm_monitor_workspace" "example" {
  name                          = "example-mamw"
  resource_group_name           = azurerm_resource_group.example.name
  location                      = azurerm_resource_group.example.location
  public_network_access_enabled = false
}

resource "azurerm_dashboard_grafana" "example" {
  name                          = "example-dg"
  resource_group_name           = azurerm_resource_group.example.name
  location                      = azurerm_resource_group.example.location
  grafana_major_version         = 11
  public_network_access_enabled = false

  azure_monitor_workspace_integrations {
    resource_id = azurerm_monitor_workspace.example.id
  }
}

resource "azurerm_dashboard_grafana_managed_private_endpoint" "example" {
  grafana_id                   = azurerm_dashboard_grafana.example.id
  name                         = "example-mpe"
  location                     = azurerm_dashboard_grafana.example.location
  private_link_resource_id     = azurerm_monitor_workspace.example.id
  group_ids                    = ["prometheusMetrics"]
  private_link_resource_region = azurerm_dashboard_grafana.example.location
}
```

## Arguments Reference

The following arguments are supported:

- `grafana_id` - (Required) The id of the associated managed Grafana. Changing this forces a new Dashboard Grafana Managed Private Endpoint to be created.

- `location` - (Required) The Azure Region where the Dashboard Grafana Managed Private Endpoint should exist. Changing this forces a new Dashboard Grafana Managed Private Endpoint to be created.

- `name` - (Required) The name which should be used for this Dashboard Grafana Managed Private Endpoint. Must be between 2 and 20 alphanumeric characters or dashes, must begin with letter and end with a letter or number. Changing this forces a new Dashboard Grafana Managed Private Endpoint to be created.

- `private_link_resource_id` - (Required) The ID of the resource to which this Dashboard Grafana Managed Private Endpoint will connect. Changing this forces a new Dashboard Grafana Managed Private Endpoint to be created.

---

- `group_ids` - (Optional) Specifies a list of private link group IDs. The value of this will depend on the private link resource to which you are connecting. Changing this forces a new Dashboard Grafana Managed Private Endpoint to be created.

- `private_link_resource_region` - (Optional) The region in which to create the private link. Changing this forces a new Dashboard Grafana Managed Private Endpoint to be created.

- `request_message` - (Optional) A message to provide in the request which will be seen by approvers.

- `tags` - (Optional) A mapping of tags which should be assigned to the Dashboard Grafana Managed Private Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Dashboard Grafana Managed Private Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dashboard Grafana Managed Private Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dashboard Grafana Managed Private Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Dashboard Grafana Managed Private Endpoint.
* `delete` - (Defaults to 5 minutes) Used when deleting the Dashboard Grafana Managed Private Endpoint.

## Import

Dashboard Grafana Managed Private Endpoint Examples can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dashboard_grafana_managed_private_endpoint.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Dashboard/grafana/workspace1/managedPrivateEndpoints/endpoint1
```
