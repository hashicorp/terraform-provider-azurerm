---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_virtual_instance"
description: |-
  Manages a SAP Virtual Instance.
---

# azurerm_workloads_sap_virtual_instance

Manages a SAP Virtual Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_workloads_sap_virtual_instance" "example" {
  name                = "X00"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  environment         = "NonProd"
  sap_product         = "S4HANA"

  configuration {
    central_server_vm_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Virtual Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `environment` - (Required) The environment type for the SAP Virtual Instance. Possible values are `NonProd` and `Prod`. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

* `sap_product` - (Required) The SAP Product type for the SAP Virtual Instance. Possible values are `ECC`, `Other` and `S4HANA`. Changing this forces a new resource to be created.

* `deployment_configuration` - (Optional) A `deployment_configuration` block as defined below. Changing this forces a new resource to be created.

* `deployment_with_os_configuration` - (Optional) A `deployment_with_os_configuration` block as defined below. Changing this forces a new resource to be created.

* `discovery_configuration` - (Optional) A `discovery_configuration` block as defined below. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Optional) The name of the managed Resource Group for the SAP Virtual Instance. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Virtual Instance.

---

A `deployment_configuration` block supports the following:

* `app_location` - (Required) The Geo-Location where the SAP system is to be created. Changing this forces a new resource to be created.

* `single_server_configuration` - (Optional) A `single_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `three_tier_configuration` - (Optional) A `three_tier_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `deployment_with_os_configuration` block supports the following:

* `app_location` - (Required) The Geo-Location where the SAP system is to be created. Changing this forces a new resource to be created.

* `os_sap_configuration` - (Required) An `os_sap_configuration` block as defined below. Changing this forces a new resource to be created.

* `single_server_configuration` - (Optional) A `single_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `three_tier_configuration` - (Optional) A `three_tier_configuration` block as defined below. Changing this forces a new resource to be created.

---

An `os_sap_configuration` block supports the following:

* `sap_fqdn` - (Required) The FQDN of the SAP system. Changing this forces a new resource to be created.

* `deployer_vm_packages` - (Optional) A `deployer_vm_packages` block as defined below. Changing this forces a new resource to be created.

---

A `deployer_vm_packages` block supports the following:

* `storage_account_id` - (Required) A `deployer_vm_packages` block as defined below. Changing this forces a new resource to be created.

* `url` - (Required) The URL of the deployer VM packages file. Changing this forces a new resource to be created.

---

A `single_server_configuration` block supports the following:

* `app_resource_group_name` - (Required) The name of the application Resource Group where SAP system resources will be deployed. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the SAP Virtual Instance. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_type` - (Optional) The supported SAP database type. Possible values are `DB2` and `HANA`. Changing this forces a new resource to be created.

* `disk_volume_configuration` - (Optional) One or more `disk_volume_configuration` blocks as defined below. Changing this forces a new resource to be created.

* `is_secondary_ip_enabled` - (Optional) Is a secondary IP Address that should be added to the Network Interface on all VMs of the SAP system being deployed enabled? Defaults to `false`. Changing this forces a new resource to be created.

* `virtual_machine_full_resource_names` - (Optional) A `virtual_machine_full_resource_names` block as defined below. Changing this forces a new resource to be created.

---

A `disk_volume_configuration` block supports the following:

* `volume_name` - (Required) The name of the DB volume of the disk configuration. Possible values are `backup`, `hana/data`, `hana/log`, `hana/shared`, `os` and `usr/sap`. Changing this forces a new resource to be created.

* `count` - (Required) The total number of disks required for the concerned volume. Changing this forces a new resource to be created.

* `size_gb` - (Required) The size of the Disk in GB. Changing this forces a new resource to be created.

* `sku_name` - (Required) The name of the Disk SKU. Changing this forces a new resource to be created.

---

A `virtual_machine_configuration` block supports the following:

* `image_reference` - (Required) An `image_reference` block as defined below. Changing this forces a new resource to be created.

* `os_profile` - (Required) An `os_profile` block as defined below. Changing this forces a new resource to be created.

* `vm_size` - (Required) The size of the Virtual Machine. Possible values are `Standard_E32ds_v4`, `Standard_E48ds_v4`, `Standard_E64ds_v4`, `Standard_M128ms`,`Standard_M128s`, `Standard_M208ms_v2`, `Standard_M208s_v2`, `Standard_M32Is`, `Standard_M32ts`, `Standard_M416ms_v2`, `Standard_M416s_v2`, `Standard_M64Is`, `Standard_M64ms` and `Standard_M64s`. Changing this forces a new resource to be created.

---

An `image_reference` block supports the following:

* `offer` - (Required) The offer of the platform image or marketplace image used to create the Virtual Machine. Changing this forces a new resource to be created.

* `publisher` - (Required) The publisher of the Image. Possible values are `RedHat` and `SUSE`. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Image. Changing this forces a new resource to be created.

* `version` - (Required) The version of the platform image or marketplace image used to create the Virtual Machine. Changing this forces a new resource to be created.

---

An `os_profile` block supports the following:

* `admin_username` - (Required) The name of the administrator account. Changing this forces a new resource to be created.

* `ssh_key_pair` - (Required) A `ssh_key_pair` block as defined below. Changing this forces a new resource to be created.

---

A `ssh_key_pair` block supports the following:

* `private_key` - (Required) The SSH public key that is used to authenticate with the VM. Changing this forces a new resource to be created.

* `public_key` - (Required) The SSH private key that is used to authenticate with the VM. Changing this forces a new resource to be created.

---

A `virtual_machine_full_resource_names` block supports the following:

* `data_disk_names` - (Optional) A mapping of Data Disk names to pass to the backend host. The keys are Volume names and the values are a comma separated string of full names for Data Disks belonging to the specific Volume. This is converted to a list before being passed to the API. Changing this forces a new resource to be created.

* `host_name` - (Optional) The full name of the host of the Virtual Machine. Changing this forces a new resource to be created.

* `network_interface_names` - (Optional) A list of full names for the Network Interface of the Virtual Machine. Changing this forces a new resource to be created.

* `os_disk_name` - (Optional) The full name of the OS Disk attached to the VM. Changing this forces a new resource to be created.

* `vm_name` - (Optional) The full name of the Virtual Machine in a single server SAP system. Changing this forces a new resource to be created.

---

A `three_tier_configuration` block supports the following:

* `app_resource_group_name` - (Required) The name of the application Resource Group where SAP system resources will be deployed. Changing this forces a new resource to be created.

* `application_server_configuration` - (Required) An `application_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `central_server_configuration` - (Required) A `central_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_server_configuration` - (Required) A `database_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `full_resource_names` - (Optional) A `full_resource_names` block as defined below. Changing this forces a new resource to be created.

* `high_availability_type` - (Optional) The high availability type for the three tier configuration. Possible values are `AvailabilitySet` and `AvailabilityZone`. Changing this forces a new resource to be created.

* `is_secondary_ip_enabled` - (Optional) Is a secondary IP Address that should be added to the Network Interface on all VMs of the SAP system being deployed enabled? Defaults to `false`. Changing this forces a new resource to be created.

* `storage_configuration` - (Optional) A `storage_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `storage_configuration` block supports the following:

* `transport_create_and_mount` - (Optional) A `transport_create_and_mount` block as defined below. Changing this forces a new resource to be created.

* `transport_mount` - (Optional) A `transport_mount` block as defined below. Changing this forces a new resource to be created.

~> **Note:** The `Skip` configuration type is enabled when `storage_configuration` isn't set.

---

A `transport_create_and_mount` block supports the following:

* `resource_group_name` - (Optional) The name of Resource Group of the transport File Share. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) The name of the Storage Account of the File Share. Changing this forces a new resource to be created.

---

A `transport_mount` block supports the following:

* `file_share_id` - (Required) The resource ID of the File Share resource. Changing this forces a new resource to be created.

* `private_endpoint_id` - (Required) The resource ID of the Private Endpoint. Changing this forces a new resource to be created.

---

An `application_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Application Server. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Application Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Optional) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `central_server_configuration` block supports the following:

* `instance_count` - (Optional) The number of instances for the Central Server. Changing this forces a new resource to be created.

* `subnet_id` - (Optional) The resource ID of the Subnet for the Central Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Optional) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `database_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Database Server. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Database Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_type` - (Optional) The database type for the Database Server. Possible values are `DB2` and `HANA`. Changing this forces a new resource to be created.

* `disk_volume_configuration` - (Optional) One or more `disk_volume_configuration` blocks as defined below. Changing this forces a new resource to be created.

---

A `full_resource_names` block supports the following:

* `application_server` - (Optional) An `application_server` block as defined below. Changing this forces a new resource to be created.

* `central_server` - (Optional) A `central_server` block as defined below. Changing this forces a new resource to be created.

* `database_server` - (Optional) A `database_server` block as defined below. Changing this forces a new resource to be created.

* `shared_storage` - (Optional) A `shared_storage` block as defined below. Changing this forces a new resource to be created.

---

An `application_server` block supports the following:

* `availability_set_name` - (Optional) The full name for the availability set. Changing this forces a new resource to be created.

* `virtual_machine` - (Optional) One or more `virtual_machine` blocks as defined below. Changing this forces a new resource to be created.

---

A `virtual_machine` block supports the following:

* `data_disk_names` - (Optional) A mapping of Data Disk names to pass to the backend host. The keys are Volume names and the values are a comma separated string of full names for Data Disks belonging to the specific Volume. This is converted to a list before being passed to the API. Changing this forces a new resource to be created.

* `host_name` - (Optional) The full name of the host of the Virtual Machine. Changing this forces a new resource to be created.

* `network_interface_names` - (Optional) A list of full names for the Network Interface of the Virtual Machine. Changing this forces a new resource to be created.

* `os_disk_name` - (Optional) The full name of the OS Disk attached to the VM. Changing this forces a new resource to be created.

* `vm_name` - (Optional) The full name of the Virtual Machine in a single server SAP system. Changing this forces a new resource to be created.

---

A `central_server` block supports the following:

* `availability_set_name` - (Optional) The full name for the availability set. Changing this forces a new resource to be created.

* `load_balancer` - (Optional) A `load_balancer` block as defined below. Changing this forces a new resource to be created.

* `virtual_machine` - (Optional) One or more `virtual_machine` blocks as defined below. Changing this forces a new resource to be created.

---

A `load_balancer` block supports the following:

* `name` - (Optional) The full resource name of the Load Balancer. Changing this forces a new resource to be created.

* `backend_pool_names` - (Optional) A list of Backend Pool names for the Load Balancer. Changing this forces a new resource to be created.

* `frontend_ip_configuration_names` - (Optional) A list of Frontend IP Configuration names. Changing this forces a new resource to be created.

* `health_probe_names` - (Optional) A list of Health Probe names. Changing this forces a new resource to be created.

---

A `database_server` block supports the following:

* `availability_set_name` - (Optional) The full name for the availability set. Changing this forces a new resource to be created.

* `load_balancer` - (Optional) A `load_balancer` block as defined below. Changing this forces a new resource to be created.

* `virtual_machine` - (Optional) One or more `virtual_machine` blocks as defined below. Changing this forces a new resource to be created.

---

A `shared_storage` block supports the following:

* `account_name` - (Optional) The full name of the Shared Storage Account. Changing this forces a new resource to be created.

* `private_endpoint_name` - (Optional) The full name of Private Endpoint for the Shared Storage Account. Changing this forces a new resource to be created.

---

A `discovery_configuration` block supports the following:

* `central_server_vm_id` - (Required) The resource ID of the Virtual Machine of the Central Server. Changing this forces a new resource to be created.

* `managed_storage_account_name` - (Optional) The name of the custom Storage Account created by the service in the managed Resource Group. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this SAP Virtual Instance. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this SAP Virtual Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Virtual Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the SAP Virtual Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Virtual Instance.
* `update` - (Defaults to 60 minutes) Used when updating the SAP Virtual Instance.
* `delete` - (Defaults to 60 minutes) Used when deleting the SAP Virtual Instance.

## Import

SAP Virtual Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_virtual_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Workloads/sapVirtualInstances/vis1
```
