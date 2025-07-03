---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_verifier_workspace"
description: |-
  Manages a Network Manager Verifier Workspace.
---

# azurerm_network_manager_verifier_workspace

Manages a Network Manager Verifier Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {}

resource "azurerm_network_manager" "example" {
  name                = "example-nm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scope {
    subscription_ids = [data.azurerm_subscription.current.id]
  }
  scope_accesses = ["Connectivity"]
}
resource "azurerm_network_manager_verifier_workspace" "example" {
  name               = "example"
  network_manager_id = azurerm_network_manager.example.id
  location           = azurerm_resource_group.example.location
  description        = "This is an example verifier workspace"

  tags = {
    foo = "bar"
    env = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Manager Verifier Workspace. Changing this forces a new Network Manager Verifier Workspace to be created.

* `location` - (Required) The Azure Region where the Network Manager Verifier Workspace should exist. Changing this forces a new Network Manager Verifier Workspace to be created.

* `network_manager_id` - (Required) The ID of the Network Manager. Changing this forces a new Network Manager Verifier Workspace to be created.

---

* `description` - (Optional) The Description of the Network Manager Verifier Workspace.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Manager Verifier Workspace.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Verifier Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Verifier Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Verifier Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Network Manager Verifier Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Verifier Workspace.

## Import

Network Manager Verifier Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_verifier_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/manager1/verifierWorkspaces/workspace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network`: 2024-05-01
