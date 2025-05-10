---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_three_tier_virtual_instance"
description: |-
  Manages an SAP Three Tier Virtual Instance with a new SAP System.
---

# azurerm_workloads_sap_three_tier_virtual_instance

Manages an SAP Three Tier Virtual Instance with a new SAP System.

-> **Note:** Before using this resource, it's required to submit the request of registering the Resource Provider with Azure CLI `az provider register --namespace "Microsoft.Workloads"`. The Resource Provider can take a while to register, you can check the status by running `az provider show --namespace "Microsoft.Workloads" --query "registrationState"`. Once this outputs "Registered" the Resource Provider is available for use.

## Example Usage

```hcl
resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "tls_public_key" "example" {
  private_key_pem = tls_private_key.example.private_key_pem
}

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

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_workloads_sap_three_tier_virtual_instance" "example" {
  name                        = "X05"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "exampleManagedRG"
  app_location                = azurerm_resource_group.app.location
  sap_fqdn                    = "sap.bpaas.com"

  three_tier_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    secondary_ip_enabled    = true

    application_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.example.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

        image {
          offer     = "RHEL-SAP-HA"
          publisher = "RedHat"
          sku       = "82sapha-gen2"
          version   = "latest"
        }

        os_profile {
          admin_username  = "testAdmin"
          ssh_private_key = tls_private_key.example.private_key_pem
          ssh_public_key  = data.tls_public_key.example.public_key_openssh
        }
      }
    }

    central_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.example.id

      virtual_machine_configuration {
        virtual_machine_size = "Standard_D16ds_v4"

        image {
          offer     = "RHEL-SAP-HA"
          publisher = "RedHat"
          sku       = "82sapha-gen2"
          version   = "latest"
        }

        os_profile {
          admin_username  = "testAdmin"
          ssh_private_key = tls_private_key.example.private_key_pem
          ssh_public_key  = data.tls_public_key.example.public_key_openssh
        }
      }
    }

    database_server_configuration {
      instance_count = 1
      subnet_id      = azurerm_subnet.example.id
      database_type  = "HANA"

      virtual_machine_configuration {
        virtual_machine_size = "Standard_E16ds_v4"

        image {
          offer     = "RHEL-SAP-HA"
          publisher = "RedHat"
          sku       = "82sapha-gen2"
          version   = "latest"
        }

        os_profile {
          admin_username  = "testAdmin"
          ssh_private_key = tls_private_key.example.private_key_pem
          ssh_public_key  = data.tls_public_key.example.public_key_openssh
        }
      }

      disk_volume_configuration {
        volume_name     = "hana/data"
        number_of_disks = 3
        size_in_gb      = 128
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "hana/log"
        number_of_disks = 3
        size_in_gb      = 128
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "hana/shared"
        number_of_disks = 1
        size_in_gb      = 256
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "usr/sap"
        number_of_disks = 1
        size_in_gb      = 128
        sku_name        = "Premium_LRS"
      }

      disk_volume_configuration {
        volume_name     = "backup"
        number_of_disks = 2
        size_in_gb      = 256
        sku_name        = "StandardSSD_LRS"
      }

      disk_volume_configuration {
        volume_name     = "os"
        number_of_disks = 1
        size_in_gb      = 64
        sku_name        = "StandardSSD_LRS"
      }
    }

    resource_names {
      application_server {
        availability_set_name = "appAvSet"

        virtual_machine {
          host_name               = "apphostName0"
          os_disk_name            = "app0osdisk"
          virtual_machine_name    = "appvm0"
          network_interface_names = ["appnic0"]

          data_disk {
            volume_name = "default"
            names       = ["app0disk0"]
          }
        }
      }

      central_server {
        availability_set_name = "csAvSet"

        load_balancer {
          name                            = "ascslb"
          backend_pool_names              = ["ascsBackendPool"]
          frontend_ip_configuration_names = ["ascsip0"]
          health_probe_names              = ["ascsHealthProbe"]
        }

        virtual_machine {
          host_name               = "ascshostName"
          os_disk_name            = "ascsosdisk"
          virtual_machine_name    = "ascsvm"
          network_interface_names = ["ascsnic"]

          data_disk {
            volume_name = "default"
            names       = ["ascsdisk"]
          }
        }
      }

      database_server {
        availability_set_name = "dbAvSet"

        load_balancer {
          name                            = "dblb"
          backend_pool_names              = ["dbBackendPool"]
          frontend_ip_configuration_names = ["dbip"]
          health_probe_names              = ["dbHealthProbe"]
        }

        virtual_machine {
          host_name               = "dbprhost"
          os_disk_name            = "dbprosdisk"
          virtual_machine_name    = "dbvmpr"
          network_interface_names = ["dbprnic"]

          data_disk {
            volume_name = "hanaData"
            names       = ["hanadatapr0", "hanadatapr1"]
          }

          data_disk {
            volume_name = "hanaLog"
            names       = ["hanalogpr0", "hanalogpr1", "hanalogpr2"]
          }

          data_disk {
            volume_name = "usrSap"
            names       = ["usrsappr0"]
          }

          data_disk {
            volume_name = "hanaShared"
            names       = ["hanasharedpr0", "hanasharedpr1"]
          }
        }
      }

      shared_storage {
        account_name          = "sharedexamplesa"
        private_endpoint_name = "examplePE"
      }
    }

    transport_create_and_mount {
      resource_group_id    = azurerm_resource_group.app.id
      storage_account_name = "exampletranssa"
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.example.id,
    ]
  }

  tags = {
    Env = "Test"
  }

  depends_on = [
    azurerm_role_assignment.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this SAP Three Tier Virtual Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Three Tier Virtual Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Three Tier Virtual Instance should exist. Changing this forces a new resource to be created.

* `app_location` - (Required) The Geo-Location where the SAP system is to be created. Changing this forces a new resource to be created.

* `environment` - (Required) The environment type for the SAP Three Tier Virtual Instance. Possible values are `NonProd` and `Prod`. Changing this forces a new resource to be created.

* `sap_fqdn` - (Required) The FQDN of the SAP system. Changing this forces a new resource to be created.

* `sap_product` - (Required) The SAP Product type for the SAP Three Tier Virtual Instance. Possible values are `ECC`, `Other` and `S4HANA`. Changing this forces a new resource to be created.

* `three_tier_configuration` - (Required) A `three_tier_configuration` block as defined below. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `managed_resource_group_name` - (Optional) The name of the managed Resource Group for the SAP Three Tier Virtual Instance. Changing this forces a new resource to be created.

* `managed_resources_network_access_type` - (Optional) The network access type for managed resources. Possible values are `Private` and `Public`. Defaults to `Public`.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Three Tier Virtual Instance.

---

A `disk_volume_configuration` block supports the following:

* `volume_name` - (Required) Specifies the volumn name of the database disk. Possible values are `backup`, `hana/data`, `hana/log`, `hana/shared`, `os` and `usr/sap`. Changing this forces a new resource to be created.

* `number_of_disks` - (Required) The total number of disks required for the concerned volume. Possible values are at least `1`. Changing this forces a new resource to be created.

* `size_in_gb` - (Required) The size of the Disk in GB. Changing this forces a new resource to be created.

* `sku_name` - (Required) The name of the Disk SKU. Possible values are `Premium_LRS`, `PremiumV2_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS` and `UltraSSD_LRS`. Changing this forces a new resource to be created.

---

A `virtual_machine_configuration` block supports the following:

* `image` - (Required) An `image` block as defined below. Changing this forces a new resource to be created.

* `os_profile` - (Required) An `os_profile` block as defined below. Changing this forces a new resource to be created.

* `virtual_machine_size` - (Required) The size of the Virtual Machine. Changing this forces a new resource to be created.

---

An `image` block supports the following:

* `offer` - (Required) Specifies the offer of the platform image or marketplace image used to create the virtual machine. Changing this forces a new resource to be created.

* `publisher` - (Required) The publisher of the Image. Possible values are `RedHat` and `SUSE`. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU of the Image. Changing this forces a new resource to be created.

* `version` - (Required) Specifies the version of the platform image or marketplace image used to create the virtual machine. Changing this forces a new resource to be created.

---

An `os_profile` block supports the following:

* `admin_username` - (Required) The name of the administrator account. Changing this forces a new resource to be created.

* `ssh_private_key` - (Required) The SSH public key that is used to authenticate with the Virtual Machine. Changing this forces a new resource to be created.

* `ssh_public_key` - (Required) The SSH private key that is used to authenticate with the Virtual Machine. Changing this forces a new resource to be created.

---

A `three_tier_configuration` block supports the following:

* `app_resource_group_name` - (Required) The name of the application Resource Group where SAP system resources will be deployed. Changing this forces a new resource to be created.

~> **Note:** While creating an SAP Three Tier Virtual Instance, the service will provision the extra SAP systems/components in the `app_resource_group_name` that are not defined in the HCL Configuration. At this time, if the `app_resource_group_name` is different from the Resource Group where SAP Three Tier Virtual Instance exists, you can set `prevent_deletion_if_contains_resources` to `false` to delete all resources defined in the HCL Configurations and the resources created in the `app_resource_group_name` with `terraform destroy`. However, if the `app_resource_group_name` is the same with the Resource Group where SAP Three Tier Virtual Instance exists, some resources, such as the subnet defined in the HCL Configuration, cannot be deleted with `terraform destroy` since the resources defined in the HCL Configuration are being referenced by the SAP system/component. In this case, you have to manually delete the SAP system/component before deleting the resources in the HCL Configuration.

* `application_server_configuration` - (Required) An `application_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `central_server_configuration` - (Required) A `central_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_server_configuration` - (Required) A `database_server_configuration` block as defined below. Changing this forces a new resource to be created.

* `resource_names` - (Optional) A `resource_names` block as defined below. Changing this forces a new resource to be created.

* `high_availability_type` - (Optional) The high availability type for the three tier configuration. Possible values are `AvailabilitySet` and `AvailabilityZone`. Changing this forces a new resource to be created.

* `secondary_ip_enabled` - (Optional) Specifies whether a secondary IP address should be added to the network interface on all VMs of the SAP system being deployed. Defaults to `false`. Changing this forces a new resource to be created.

* `transport_create_and_mount` - (Optional) A `transport_create_and_mount` block as defined below. Changing this forces a new resource to be created.

~> **Note:** The file share configuration uses `skip` by default when `transport_create_and_mount` isn't set.

~> **Note:** Due to [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/25209) where the Storage File Share Id is not defined correctly, it is not currently possible to support using Transport Mount.

---

A `transport_create_and_mount` block supports the following:

* `resource_group_id` - (Optional) The ID of the Resource Group of the transport File Share. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) The name of the Storage Account of the File Share. Changing this forces a new resource to be created.

---

An `application_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Application Server. Possible values are at least `1`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Application Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `central_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Central Server. Possible values are at least `1`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Central Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

---

A `database_server_configuration` block supports the following:

* `instance_count` - (Required) The number of instances for the Database Server. Possible values are at least `1`. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The resource ID of the Subnet for the Database Server. Changing this forces a new resource to be created.

* `virtual_machine_configuration` - (Required) A `virtual_machine_configuration` block as defined below. Changing this forces a new resource to be created.

* `database_type` - (Optional) The database type for the Database Server. Possible values are `DB2` and `HANA`. Changing this forces a new resource to be created.

* `disk_volume_configuration` - (Optional) One or more `disk_volume_configuration` blocks as defined below. Changing this forces a new resource to be created.

---

A `resource_names` block supports the following:

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

* `data_disk` - (Optional) One or more `data_disk` blocks as defined below. Changing this forces a new resource to be created.

* `host_name` - (Optional) The full name of the host of the Virtual Machine. Changing this forces a new resource to be created.

* `network_interface_names` - (Optional) A list of full names for the Network Interface of the Virtual Machine. Changing this forces a new resource to be created.

* `os_disk_name` - (Optional) The full name of the OS Disk attached to the Virtual Machine. Changing this forces a new resource to be created.

* `virtual_machine_name` - (Optional) The full name of the Virtual Machine in a single server SAP system. Changing this forces a new resource to be created.

---

A `data_disk` block supports the following:

* `volume_name` - (Required) The name of the Volume. Changing this forces a new resource to be created.

~> **Note:** Possible value for Application Server and Central Server is `default`.

~> **Note:** Possible values for Database Server are `hanaData`, `hanaLog`, `hanaShared` and `usrSap`.

* `names` - (Required) A list of full names of Data Disks per Volume. Changing this forces a new resource to be created.

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

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this SAP Three Tier Virtual Instance. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this SAP Three Tier Virtual Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Three Tier Virtual Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the SAP Three Tier Virtual Instance with new SAP System.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Three Tier Virtual Instance with new SAP System.
* `update` - (Defaults to 1 hour) Used when updating the SAP Three Tier Virtual Instance with new SAP System.
* `delete` - (Defaults to 1 hour) Used when deleting the SAP Three Tier Virtual Instance with new SAP System.

## Import

SAP Three Tier Virtual Instances with new SAP Systems can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_three_tier_virtual_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Workloads/sapVirtualInstances/vis1
```
