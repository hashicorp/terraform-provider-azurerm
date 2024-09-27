---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_workspace"
description: |-
  Manages an Azure Monitor Workspace.
---

# azurerm_monitor_workspace

Manages an Azure Monitor Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_workspace" "example" {
  name                = "example-mamw"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Azure Monitor Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Azure Monitor Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Azure Monitor Workspace should exist. Changing this forces a new resource to be created.

* `public_network_access_enabled` - (Optional) Is public network access enabled? Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Monitor Workspace.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Monitor Workspace.

* `query_endpoint` - The query endpoint for the Azure Monitor Workspace.

* `default_data_collection_endpoint_id` - The ID of the managed default Data Collection Endpoint created with the Azure Monitor Workspace.

* `default_data_collection_rule_id` - The ID of the managed default Data Collection Rule created with the Azure Monitor Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Monitor Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Monitor Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Monitor Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Monitor Workspace.

## Import

Azure Monitor Workspace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Monitor/accounts/azureMonitorWorkspace1
```
