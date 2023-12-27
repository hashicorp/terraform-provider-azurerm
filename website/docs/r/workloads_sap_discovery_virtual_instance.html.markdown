---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_discovery_virtual_instance"
description: |-
  Manages a SAP Discovery Virtual Instance.
---

# azurerm_workloads_sap_discovery_virtual_instance

Manages a SAP Discovery Virtual Instance.

-> **Note:** Before using this resource, it's required to submit the request of registering the Resource Provider with Azure CLI `az provider register --namespace "Microsoft.Workloads"`. The Resource Provider can take a while to register, you can check the status by running `az provider show --namespace "Microsoft.Workloads" --query "registrationState"`. Once this outputs "Registered" the Resource Provider is available for use.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "tls_public_key" "example" {
  private_key_pem = tls_private_key.example.private_key_pem
}

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

resource "azurerm_workloads_sap_discovery_virtual_instance" "example" {
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
      secondary_ip_enabled    = true

      virtual_machine_configuration {
        virtual_machine_size = "Standard_E32ds_v4"

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

      virtual_machine_full_resource_names {
        host_name               = "apphostName0"
        os_disk_name            = "app0osdisk"
        virtual_machine_name    = "appvm0"
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

* `name` - (Required) The name which should be used for this SAP Discovery Virtual Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Discovery Virtual Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Discovery Virtual Instance should exist. Changing this forces a new resource to be created.

* `discovery_configuration` - (Required) A `discovery_configuration` block as defined below. Changing this forces a new resource to be created.

* `environment` - (Required) The environment type for the SAP Discovery Virtual Instance. Possible values are `NonProd` and `Prod`. Changing this forces a new resource to be created.

* `sap_product` - (Required) The SAP Product type for the SAP Discovery Virtual Instance. Possible values are `ECC`, `Other` and `S4HANA`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `managed_resource_group_name` - (Optional) The name of the managed Resource Group for the SAP Discovery Virtual Instance. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Discovery Virtual Instance.

---

A `discovery_configuration` block supports the following:

* `central_server_virtual_machine_id` - (Required) The ID of the Virtual Machine of the Central Server. Changing this forces a new resource to be created.

* `managed_storage_account_name` - (Optional) The name of the custom Storage Account created by the service in the managed Resource Group. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this SAP Discovery Virtual Instance. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this SAP Discovery Virtual Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Discovery Virtual Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the SAP Discovery Virtual Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Discovery Virtual Instance.
* `update` - (Defaults to 60 minutes) Used when updating the SAP Discovery Virtual Instance.
* `delete` - (Defaults to 60 minutes) Used when deleting the SAP Discovery Virtual Instance.

## Import

SAP Discovery Virtual Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_discovery_virtual_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Workloads/sapVirtualInstances/vis1
```
