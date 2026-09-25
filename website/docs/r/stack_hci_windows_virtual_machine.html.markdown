---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_windows_virtual_machine"
description: |-
  Manages an Azure Stack HCI Windows Virtual Machine.
---

# azurerm_stack_hci_windows_virtual_machine

Manages an Azure Stack HCI Windows Virtual Machine.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_stack_hci_logical_network" "example" {
  name                = "example-ln"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  virtual_switch_name = "ConvergedSwitch(managementcompute)"
  dns_servers         = ["192.168.1.254"]

  subnet {
    ip_allocation_method = "Static"
    address_prefix       = "192.168.1.0/24"
    ip_pool {
      start = "192.168.1.0"
      end   = "192.168.1.255"
    }

    route {
      address_prefix      = "0.0.0.0/0"
      next_hop_ip_address = "192.168.1.1"
    }
  }
}

resource "azurerm_stack_hci_network_interface" "example" {
  name                = "example-ni"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  dns_servers         = ["192.168.1.254"]
  mac_address         = "02:ec:01:0c:00:08"

  ip_configuration {
    private_ip_address = "192.168.1.15"
    subnet_id          = azurerm_stack_hci_logical_network.example.id
  }
}

// service principal of 'Microsoft.AzureStackHCI Resource Provider'
data "azuread_service_principal" "hciRp" {
  client_id = "1412d89f-b8a8-4111-b4fd-e82905cbd85d"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Azure Connected Machine Resource Manager"
  principal_id         = data.azuread_service_principal.hciRp.object_id
}

resource "azurerm_stack_hci_virtual_hard_disk" "example" {
  name                = "example-vhd"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  disk_size_in_gb     = 2

  lifecycle {
    ignore_changes = [storage_path_id]
  }
}

resource "azurerm_stack_hci_marketplace_gallery_image" "example" {
  name                = "example-mgi"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  hyperv_generation   = "V2"
  os_type             = "Windows"
  version             = "20348.2655.240905"
  identifier {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
  }
  tags = {
    foo = "bar"
    env = "example"
  }
  depends_on = [azurerm_role_assignment.example]
}

resource "azurerm_arc_machine" "example" {
  name                = "example-hcivm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  kind                = "HCI"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_stack_hci_windows_virtual_machine" "example" {
  arc_machine_id     = azurerm_stack_hci_windows_virtual_machine.example.arc_machine_id
  custom_location_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"

  hardware_profile {
    vm_size          = "Custom"
    processor_number = 2
    memory_in_mb     = 8192
  }

  network_profile {
    network_interface_ids = [azurerm_stack_hci_network_interface.example.id]
  }

  os_profile {
    admin_username = "adminuser"
    admin_password = "!password!@#$"
    computer_name  = "examplevm"
  }

  storage_profile {
    data_disk_ids = [azurerm_stack_hci_virtual_hard_disk.example.id]
    image_id      = azurerm_stack_hci_marketplace_gallery_image.example.id
  }

  lifecycle {
    ignore_changes = [storage_profile.0.vm_config_storage_path_id]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `arc_machine_id` - (Required) The ID of the Arc Machine. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `custom_location_id` - (Required) The ID of the Custom Location. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `hardware_profile` - (Required) A `hardware_profile` block as defined below. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `network_profile` - (Required) A `network_profile` block as defined below.

* `os_profile` - (Required) An `os_profile` block as defined below. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `storage_profile` - (Required) A `storage_profile` block as defined below.

---

* `http_proxy_configuration` - (Optional) A `http_proxy_configuration` block as defined below. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `secure_boot_enabled` - (Optional) Whether to enable secure boot. Defaults to `true`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `security_type` - (Optional) The security type of the virtual machine. Possible values are `TrustedLaunch` and `ConfidentialVM`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `tpm_enabled` - (Optional) Whether to enable the TPM (Trusted Platform Module). Defaults to `false`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

---

A `dynamic_memory` block supports the following:

* `maximum_memory_in_mb` - (Required) The maximum memory in Megabytes . Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `minimum_memory_in_mb` - (Required) The minimum memory in Megabytes. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `target_memory_buffer_percentage` - (Required) The percentage of total memory to reserve as extra memory for a virtual machine instance during runtime. Possible value can be in the range of `5` to `2000`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

---

A `hardware_profile` block supports the following:

* `memory_in_mb` - (Required) The memory in Megabytes. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `processor_number` - (Required) The number of processors. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `vm_size` - (Required) The size of virtual machine. Possible values are
`Custom`, `Default`, `Standard_A4_v2`, `Standard_A2_v2`, `Standard_D8s_v3`, `Standard_D4s_v3`, `Standard_D16s_v3`, `Standard_DS5_v2`, `Standard_DS4_v2`, `Standard_DS13_v2`, `Standard_DS3_v2`, `Standard_DS2_v2`, `Standard_D32s_v3`, `Standard_D2s_v3`, `Standard_K8S5_v1`, `Standard_K8S4_v1`, `Standard_K8S3_v1`, `Standard_K8S2_v1`, `Standard_K8S_v1`, `Standard_NK12`, `Standard_NK6`, `Standard_NV12` and `Standard_NV6`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `dynamic_memory` - (Optional) A `dynamic_memory` block as defined above. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

---

A `http_proxy_configuration` block supports the following:

* `http_proxy` - (Required) The HTTP proxy server endpoint. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `https_proxy` - (Required) The HTTPS proxy server endpoint. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `no_proxy` - (Optional) Specifies a list of endpoints that should not go through proxy. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `trusted_ca` - (Optional) Alternative CA cert to use for connecting to proxy servers. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

---

A `network_profile` block supports the following:

* `network_interface_ids` - (Required) Specifies a list of Azure Stack HCI Network Interface IDs.

---

An `os_profile` block supports the following:

* `admin_username` - (Required) The username of admin. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `computer_name` - (Required) The computer name. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `admin_password` - (Optional) The password of admin. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `automatic_update_enabled` - (Optional) Whether automatic update is enabled. Defaults to `false`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `provision_vm_agent_enabled` - (Optional) Whether the Arc VM Agent should be installed during the virtual machine creation process. Defaults to `false`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `provision_vm_config_agent_enabled` - (Optional) Whether the VM Config Agent should be installed during the virtual machine creation process. Defaults to `false`. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `ssh_public_key` - (Optional) One or more `ssh_public_key` blocks as defined below. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `time_zone` - (Optional) The timezone for the virtual machine. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

---

A `ssh_public_key` block supports the following:

* `key_data` - (Required) SSH public key certificate used to authenticate with the VM through SSH. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `path` - (Required) The full path on the created VM where ssh public key is stored. If the file already exists, the specified key is appended to the file. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

---

A `storage_profile` block supports the following:

* `data_disk_ids` - (Required) Specifies a list of Azure Stack HCI Virtual Hard Disk IDs.

* `image_id` - (Required) The ID of the Stack HCI VM Image. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

* `vm_config_storage_path_id` - (Optional) The ID of the Azure Stack HCI Storage Path to host the VM configuration file. Changing this forces a new Azure Stack HCI Windows Virtual Machine to be created.

-> **Note:** If `vm_config_storage_path_id` is not specified, it will be assigned by the server. If you experience a diff you may need to add this to `ignore_changes`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Stack HCI Windows Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Azure Stack HCI Windows Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Windows Virtual Machine.
* `update` - (Defaults to 1 hour and 30 minutes) Used when updating the Azure Stack HCI Windows Virtual Machine.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Windows Virtual Machine.

## Import

Azure Stack HCI Windows Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_windows_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default
```
