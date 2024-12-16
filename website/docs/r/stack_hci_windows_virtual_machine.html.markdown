---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_windows_virtual_machine"
description: |-
  Manages a Stack HCI Windows Virtual Machine.
---

# azurerm_stack_hci_windows_virtual_machine

Manages a Stack HCI Windows Virtual Machine.

## Example Usage

```hcl
resource "azurerm_stack_hci_windows_virtual_machine" "example" {

  network_profile {
    network_interface_ids = [ "example" ]    
  }
  custom_location_id = "TODO"
  arc_machine_id = "TODO"

  hardware_profile {
    vm_size = "TODO"
    processor_number = 42
    memory_mb = 42    
  }

  os_profile {
    admin_username = "TODO"
    computer_name = "example"    
  }

  storage_profile {
    data_disk_ids = [ "example" ]
    image_id = "TODO"    
  }
}
```

## Arguments Reference

The following arguments are supported:

* `arc_machine_id` - (Required) The ID of the TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `custom_location_id` - (Required) The ID of the TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `hardware_profile` - (Required) A `hardware_profile` block as defined below. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `network_profile` - (Required) A `network_profile` block as defined below.

* `os_profile` - (Required) A `os_profile` block as defined below. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `storage_profile` - (Required) A `storage_profile` block as defined below.

---

* `http_proxy_configuration` - (Optional) A `http_proxy_configuration` block as defined below. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `secure_boot_enabled` - (Optional) Should the TODO be enabled? Defaults to `true`. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `security_type` - (Optional) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `tpm_enabled` - (Optional) Should the TODO be enabled? Defaults to `false`. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

---

A `dynamic_memory` block supports the following:

* `maximum_memory_mb` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `minimum_memory_mb` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `target_memory_buffer` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

---

A `hardware_profile` block supports the following:

* `memory_mb` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `processor_number` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `vm_size` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `dynamic_memory` - (Optional) A `dynamic_memory` block as defined above. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

---

A `http_proxy_configuration` block supports the following:

* `http_proxy` - (Optional) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `https_proxy` - (Optional) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `no_proxy` - (Optional) Specifies a list of TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `trusted_ca` - (Optional) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

---

A `network_profile` block supports the following:

* `network_interface_ids` - (Required) Specifies a list of TODO.

---

A `os_profile` block supports the following:

* `admin_username` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `computer_name` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `admin_password` - (Optional) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `automatic_update_enabled` - (Optional) Should the TODO be enabled? Defaults to `false`. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `provision_vm_agent_enabled` - (Optional) Should the TODO be enabled? Defaults to `false`. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `provision_vm_config_agent_enabled` - (Optional) Should the TODO be enabled? Defaults to `false`. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `ssh_public_key` - (Optional) One or more `ssh_public_key` blocks as defined below. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `time_zone` - (Optional) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

---

A `ssh_public_key` block supports the following:

* `key_data` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `path` - (Required) TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

---

A `storage_profile` block supports the following:

* `data_disk_ids` - (Required) Specifies a list of TODO.

* `image_id` - (Required) The ID of the TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `os_disk_id` - (Optional) The ID of the TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

* `vm_config_storage_path_id` - (Optional) The ID of the TODO. Changing this forces a new Stack HCI Windows Virtual Machine to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Stack HCI Windows Virtual Machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour and 30 minutes) Used when creating the Stack HCI Windows Virtual Machine.
* `read` - (Defaults to 5 minutes) Used when retrieving the Stack HCI Windows Virtual Machine.
* `update` - (Defaults to 1 hour and 30 minutes) Used when updating the Stack HCI Windows Virtual Machine.
* `delete` - (Defaults to 30 minutes) Used when deleting the Stack HCI Windows Virtual Machine.

## Import

Stack HCI Windows Virtual Machines can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_windows_virtual_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default
```