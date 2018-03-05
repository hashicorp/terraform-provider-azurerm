---
layout: "azurerm"
page_title: "AzureRM: Authenticating via Managed Service Identity"
sidebar_current: "docs-azurerm-index-authentication-msi"
description: |-
  The Azure Resource Manager provider supports authenticating via multiple means. This guide will cover configuring a Managed Service Identity which can be used to access Azure Resource Manager.

---

# Authenticating to Azure Resource Manager using Managed Service Identity

Terraform supports authenticating to Azure through Managed Service Identity, Service Principal or the Azure CLI.

Managed Service Identity can be used to access other Azure Services from within a Virtual Machine in Azure instead of specifying a Service Principal or Azure CLI credentials.

## Configuring Managed Service Identity

Managed Service Identity allows an Azure virtual machine to retrieve a token to access the Azure API without needing to pass in credentials. This works by creating a service principal in Azure Active Directory that is associated to a virtual machine. This service principal can then be granted permissions to Azure resources.
There are various ways to configure managed service identity - see the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/active-directory/msi-overview) for details.
You can then run Terraform from the MSI enabled virtual machine by setting the `use_msi` provider option to `true`.

### Configuring Managed Service Identity using Terraform

Managed service identity can also be configured using Terraform. The following template shows how. Note that for a Linux VM you must use the `ManagedIdentityExtensionForLinux` extension.

```hcl
resource "azurerm_virtual_machine" "virtual_machine" {
  name                  = "test"
  location              = "${var.location}"
  resource_group_name   = "test"
  network_interface_ids = ["${azurerm_network_interface.test.id}"]
  vm_size               = "Standard_DS1_v2"

  identity = {
    type = "SystemAssigned"
  }

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-smalldisk"
    version   = "latest"
  }

  storage_os_disk {
    name              = "test"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "test"
    admin_username = "username"
    admin_password = "password"
  }

  os_profile_windows_config {
    provision_vm_agent        = true
    enable_automatic_upgrades = false
  }
}

resource "azurerm_virtual_machine_extension" "virtual_machine_extension" {
  name                 = "test"
  location             = "${var.location}"
  resource_group_name  = "test"
  virtual_machine_name = "${azurerm_virtual_machine.virtual_machine.name}"
  publisher            = "Microsoft.ManagedIdentity"
  type                 = "ManagedIdentityExtensionForWindows"
  type_handler_version = "1.0"

  settings = <<SETTINGS
    {
        "port": 50342
    }
SETTINGS
}

data "azurerm_subscription" "subscription" {}

data "azurerm_builtin_role_definition" "builtin_role_definition" {
  name = "Contributor"
}

# Grant the VM identity contributor rights to the current subscription
resource "azurerm_role_assignment" "role_assignment" {
  scope              = "${data.azurerm_subscription.subscription.id}"
  role_definition_id = "${data.azurerm_subscription.subscription.id}${data.azurerm_builtin_role_definition.builtin_role_definition.id}"
  principal_id       = "${lookup(azurerm_virtual_machine.virtual_machine.identity[0], "principal_id")}"

  lifecycle {
    ignore_changes = ["name"]
  }
}
```