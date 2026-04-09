---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_manager_verifier_workspace_reachability_analysis_intent"
description: |-
  Manages a Network Manager Verifier Workspace Reachability Analysis Intent.
---

# azurerm_network_manager_verifier_workspace_reachability_analysis_intent

Manages a Network Manager Verifier Workspace Reachability Analysis Intent.

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
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                            = "example-machine"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  size                            = "Standard_B1ls"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssw0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_network_manager_verifier_workspace_reachability_analysis_intent" "example" {
  name                    = "example-intent"
  verifier_workspace_id   = azurerm_network_manager_verifier_workspace.example.id
  source_resource_id      = azurerm_linux_virtual_machine.example.id
  destination_resource_id = azurerm_linux_virtual_machine.example.id
  description             = "example"
  ip_traffic {
    source_ips        = ["10.0.2.1"]
    source_ports      = ["80"]
    destination_ips   = ["10.0.2.2"]
    destination_ports = ["*"]
    protocols         = ["Any"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Manager Verifier Workspace Reachability Analysis Intent. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `verifier_workspace_id` - (Required) The ID of the Network Manager Verifier Workspace. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `destination_resource_id` - (Required) The ID of the destination resource. The value can be the ID of either Public internet, Cosmos DB, Storage Account, SQL Server, Virtual machines, or Subnet. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `ip_traffic` - (Required) An `ip_traffic` block as defined below. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `source_resource_id` - (Required) The ID of the source resource. The value can be the ID of either Public internet, Virtual machines, or Subnet. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

---

* `description` - (Optional) The description of the resource. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

---

A `ip_traffic` block supports the following:

* `destination_ips` - (Required) Specifies a list of IPv4 or IPv6 addresses or ranges using CIDR notation of the source you want to verify. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `destination_ports` - (Required) Specifies a list of ports or ranges of the destination you want to verify. To specify any port, use `["*"]`. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `protocols` - (Required) Specifies a list of network protocols. Possible values are `Any`, `TCP`, `UDP` and `ICMP`. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `source_ips` - (Required) Specifies a list of IPv4 or IPv6 addresses or ranges using CIDR notation of the source you want to verify. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

* `source_ports` - (Required) Specifies a list of ports or ranges of the source you want to verify. To specify any port, use `["*"]`. Changing this forces a new Network Manager Verifier Workspace Reachability Analysis Intent to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Manager Verifier Workspace Reachability Analysis Intent.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Manager Verifier Workspace Reachability Analysis Intent.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Manager Verifier Workspace Reachability Analysis Intent.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Manager Verifier Workspace Reachability Analysis Intent.

## Import

Network Manager Verifier Workspace Reachability Analysis Intents can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_manager_verifier_workspace_reachability_analysis_intent.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/networkManagers/manager1/verifierWorkspaces/workspace1/reachabilityAnalysisIntents/intent1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2025-01-01
