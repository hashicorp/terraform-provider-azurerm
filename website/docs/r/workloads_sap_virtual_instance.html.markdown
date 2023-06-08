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
data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Azure Center for SAP solutions service role"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_resource_group" "app" {
  name     = "example-sapapp"
  location = "West Europe"

  depends_on = [
    azurerm_subnet.example
  ]
}

resource "azurerm_workloads_sap_virtual_instance" "example" {
  name                        = "X05"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "managedTestRG"

  deployment_with_os_configuration {
    app_location = azurerm_resource_group.app.location

    os_sap_configuration {
      sap_fqdn = "sap.bpaas.com"
    }

    single_server_configuration {
      app_resource_group_name = azurerm_resource_group.app.name
      subnet_id               = azurerm_subnet.example.id
      database_type           = "HANA"
      is_secondary_ip_enabled = true

      virtual_machine_configuration {
        vm_size = "Standard_E32ds_v4"

        image_reference {
          offer     = "RHEL-SAP-HA"
          publisher = "RedHat"
          sku       = "82sapha-gen2"
          version   = "latest"
        }

        os_profile {
          admin_username = "testAdmin"

          ssh_key_pair {
            private_key = "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
            public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC8k1K0mjpBtyBRcr+74p6zShq09hMXmQZJJem41RlqV132kbGUy/EHcQLoOmp2qdv6Kz/jgz49YWM2WxxhB8G+Uh8G58feYFaGNkT+WTrfE/XLHxFnXMghZgCEiV8vfNUROuHux7UfsxiDt6Nw2215fkJ8lQjQEduKPfS7OmoQQDUvGeg8rxUbwouYdR9z5MmcNKSiBC0Z3EEBWVkY0jsbASarE9iRLDT5JrXeeZWCSe+/lg3Ofe7b28f+nXo9LH/voRVc+eKsx6T6reW9wVtm6ygEBXZbqvuEeqISgLvSFWXD7SsO2CbyJpmaOkzjIMbzlaHV126+PYdAY2vf0gVhCwRt5wU/ueFloMmTb3HbLkyVczluBUsSfcJ3H5DYPECiCy5qKaPsGHDzFGByc/0Yy+fNn2L/WrNgNuC3LaeszTZ7YqNVo/6tXTGSqOMrv8AD1rP71sX2HDSdlLf/HyLRgNjfJn/KR3ZE9QZ2Hl9TRcZ9DN+V/wvc/lci/7nsSx0= generated-by-azure"
          }
        }
      }

      disk_volume_configuration {
        volume_name = "hana/data"
        count       = 3
        size_gb     = 128
        sku_name    = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name = "hana/log"
        count       = 3
        size_gb     = 128
        sku_name    = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name = "hana/shared"
        count       = 1
        size_gb     = 256
        sku_name    = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name = "usr/sap"
        count       = 1
        size_gb     = 128
        sku_name    = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name = "backup"
        count       = 2
        size_gb     = 256
        sku_name    = "StandardSSD_LRS"
      }

      disk_volume_configuration {
        volume_name = "os"
        count       = 1
        size_gb     = 64
        sku_name    = "StandardSSD_LRS"
      }

      virtual_machine_full_resource_names {
        host_name               = "apphostName0"
        os_disk_name            = "app0osdisk"
        vm_name                 = "appvm0"
        network_interface_names = ["appnic0"]

        data_disk_names = {
          default = "app0disk0"
        }
      }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.example.id,
    ]
  }

  depends_on = [
    azurerm_role_assignment.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Virtual Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `environment` - (Required) The environment type for the SAP Virtual Instance. Possible values are `NonProd` and `Prod`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `sap_product` - (Required) The SAP Product type for the SAP Virtual Instance. Possible values are `ECC`, `Other` and `S4HANA`. Changing this forces a new resource to be created.

* `deployment_with_os_configuration` - (Optional) A `deployment_with_os_configuration` block as defined below. Changing this forces a new resource to be created.

* `discovery_configuration` - (Optional) A `discovery_configuration` block as defined below. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Optional) The name of the managed Resource Group for the SAP Virtual Instance. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Virtual Instance.

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

~> **Note:** While creating SAP Virtual Instance, service would provision the extra SAP system/component in `app_resource_group_name` which aren't defined in tf config. At this time, when `app_resource_group_name` is different with the Resource Group where SAP Virtual Instance exists, we can set `prevent_deletion_if_contains_resources` to `false` to delete all resources defined in tf config and the resources created in `app_resource_group_name` with `tf destroy`. But when `app_resource_group_name` is same with the Resource Group where SAP Virtual Instance exists, some resources like the subnet defined in tf config cannot be deleted with `tf destroy` since the resources defined in tf config are being referenced by the SAP system/component created in `app_resource_group_name`. So it has to manually delete the resources in `app_resource_group_name` first after the SAP Virtual Instance is deleted, and then the resources in tf config can be deleted successfully for this situation.

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

* `vm_size` - (Required) The size of the Virtual Machine. Changing this forces a new resource to be created.

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

* `transport_create_and_mount` - (Optional) A `transport_create_and_mount` block as defined below. Changing this forces a new resource to be created.

* `transport_mount` - (Optional) A `transport_mount` block as defined below. Changing this forces a new resource to be created.

~> **Note:** The `Skip` configuration type would be enabled when the `transport_create_and_mount` and the `transport_mount` aren't set.

---

A `transport_create_and_mount` block supports the following:

* `resource_group_name` - (Optional) The name of Resource Group of the transport File Share. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) The name of the Storage Account of the File Share. Changing this forces a new resource to be created.

---

A `transport_mount` block supports the following:

* `share_file_id` - (Required) The resource ID of the Share File resource. Changing this forces a new resource to be created.

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
