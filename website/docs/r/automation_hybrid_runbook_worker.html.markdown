---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_hybrid_runbook_worker"
description: |-
  Manages a Automation.
---

# azurerm_automation_hybrid_runbook_worker

Manages a Automation Hybrid Runbook Worker.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_hybrid_runbook_worker_group" "example" {
  name                    = "example"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["192.168.1.0/24"]
  location            = azurerm_resource_group.example.location
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["192.168.1.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "vm-example"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                = "example-vm"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  size                            = "Standard_B1s"
  admin_username                  = "testadmin"
  admin_password                  = "Password1234!"
  disable_password_authentication = false

  source_image_reference {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7.5"
    version   = "latest"
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  network_interface_ids = [azurerm_network_interface.example.id]
}

resource "azurerm_automation_hybrid_runbook_worker" "example" {
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  worker_group_name       = azurerm_automation_hybrid_runbook_worker_group.example.name
  vm_resource_id          = azurerm_linux_virtual_machine.example.id
  worker_id               = "00000000-0000-0000-0000-000000000000" #unique uuid
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Automation should exist. Changing this forces a new Automation to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Hybrid Worker is created. Changing this forces a new resource to be created.

* `worker_group_name` - (Required) The name of the HybridWorker Group. Changing this forces a new Automation to be created.

* `worker_id` - (Required) Specify the ID of this HybridWorker in UUID notation. Changing this forces a new Automation to be created.

* `vm_resource_id` - (Required) The ID of the virtual machine used for this HybridWorker. Changing this forces a new Automation to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Hybrid Runbook Worker.

* `ip` - The IP address of assigned machine.

* `last_seen_date_time` - Last Heartbeat from the Worker.

* `registration_date_time` - The registration time of the worker machine.

* `worker_name` - The name of HybridWorker.

* `worker_type` - The type of the HybridWorker, the possible values are `HybridV1` and `HybridV2`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_hybrid_runbook_worker.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/hybridRunbookWorkerGroups/group1/hybridRunbookWorkers/00000000-0000-0000-0000-000000000000
```
