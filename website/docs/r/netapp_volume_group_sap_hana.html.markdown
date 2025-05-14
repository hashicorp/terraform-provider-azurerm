---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_group_sap_hana"
description: |-
  Manages a Application Volume Group for SAP HANA application.
---

# azurerm_netapp_volume_group_sap_hana

Manages a Application Volume Group for SAP HANA application.

-> **Note:** This feature is intended to be used for SAP-HANA workloads only, with several requirements, please refer to [Understand Azure NetApp Files application volume group for SAP HANA](https://learn.microsoft.com/en-us/azure/azure-netapp-files/application-volume-group-introduction) document as the starting point to understand this feature before using it with Terraform.

## Example Usage

```hcl
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "random_string" "example" {
  length  = 12
  special = true
}

locals {
  admin_username = "exampleadmin"
  admin_password = random_string.example.result
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.88.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-delegated-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.88.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "example1" {
  name                 = "${var.prefix}-hosts-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.88.1.0/24"]
}

resource "azurerm_proximity_placement_group" "example" {
  name                = "${var.prefix}-ppg"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_availability_set" "example" {
  name                = "${var.prefix}-avset"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  proximity_placement_group_id = azurerm_proximity_placement_group.example.id
}

resource "azurerm_network_interface" "example" {
  name                = "${var.prefix}-nic"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example1.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                            = "${var.prefix}-vm"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  size                            = "Standard_M8ms"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false
  proximity_placement_group_id    = azurerm_proximity_placement_group.example.id
  availability_set_id             = azurerm_availability_set.example.id
  network_interface_ids = [
    azurerm_network_interface.example.id
  ]

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netapp-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  depends_on = [
    azurerm_subnet.example,
    azurerm_subnet.example1
  ]
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-netapp-pool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 8
  qos_type            = "Manual"
}

resource "azurerm_netapp_volume_group_sap_hana" "example" {
  name                   = "${var.prefix}-netapp-volumegroup"
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  account_name           = azurerm_netapp_account.example.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                         = "${var.prefix}-netapp-volume-1"
    volume_path                  = "my-unique-file-path-1"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.example.id
    subnet_id                    = azurerm_subnet.example.id
    proximity_placement_group_id = azurerm_proximity_placement_group.example.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "foo" = "bar"
    }
  }

  volume {
    name                         = "${var.prefix}-netapp-volume-2"
    volume_path                  = "my-unique-file-path-2"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.example.id
    subnet_id                    = azurerm_subnet.example.id
    proximity_placement_group_id = azurerm_proximity_placement_group.example.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "foo" = "bar"
    }
  }

  volume {
    name                         = "${var.prefix}-netapp-volume-3"
    volume_path                  = "my-unique-file-path-3"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.example.id
    subnet_id                    = azurerm_subnet.example.id
    proximity_placement_group_id = azurerm_proximity_placement_group.example.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  depends_on = [
    azurerm_linux_virtual_machine.example,
    azurerm_proximity_placement_group.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `account_name` - (Required) Name of the account where the application volume group belong to. Changing this forces a new Application Volume Group to be created and data will be lost.

* `application_identifier` - (Required) The SAP System ID, maximum 3 characters, e.g. `SH9`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `group_description` - (Required) Volume group description. Changing this forces a new Application Volume Group to be created and data will be lost.

* `location` - (Required) The Azure Region where the Application Volume Group should exist. Changing this forces a new Application Volume Group to be created and data will be lost.

* `name` - (Required) The name which should be used for this Application Volume Group. Changing this forces a new Application Volume Group to be created and data will be lost.

* `resource_group_name` - (Required) The name of the Resource Group where the Application Volume Group should exist. Changing this forces a new Application Volume Group to be created and data will be lost.

* `volume` - (Required) One or more `volume` blocks as defined below.

---

A `volume` block supports the following:

* `capacity_pool_id` - (Required) The ID of the Capacity Pool. Changing this forces a new Application Volume Group to be created and data will be lost.

* `name` - (Required) The name which should be used for this volume. Changing this forces a new Application Volume Group to be created and data will be lost.

* `protocols` - (Required) The target volume protocol expressed as a list. Changing this forces a new Application Volume Group to be created and data will be lost. Supported values for Application Volume Group include `NFSv3` or `NFSv4.1`, multi-protocol is not supported and there are certain rules on which protocol is supporteed per volume spec, please check [Configure application volume groups for the SAP HANA REST API](https://learn.microsoft.com/en-us/azure/azure-netapp-files/configure-application-volume-group-sap-hana-api) document for details.

* `proximity_placement_group_id` - (Optional) The ID of the proximity placement group. Changing this forces a new Application Volume Group to be created and data will be lost. For SAP-HANA application, it is required to have PPG enabled so Azure NetApp Files can pin the volumes next to your compute resources, please check [Requirements and considerations for application volume group for SAP HANA](https://learn.microsoft.com/en-us/azure/azure-netapp-files/application-volume-group-considerations) for details and other requirements.

* `security_style` - (Required) Volume security style. Possible values are `ntfs` and `unix`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `service_level` - (Required) Volume security style. Possible values are `Premium`, `Standard` and `Ultra`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `snapshot_directory_visible` - (Required) Specifies whether the .snapshot (NFS clients) path of a volume is visible. Changing this forces a new Application Volume Group to be created and data will be lost.

* `storage_quota_in_gb` - (Required) The maximum Storage Quota allowed for a file system in Gigabytes.

* `subnet_id` - (Required) The ID of the Subnet the NetApp Volume resides in, which must have the `Microsoft.NetApp/volumes` delegation. Changing this forces a new Application Volume Group to be created and data will be lost.

* `throughput_in_mibps` - (Required) Throughput of this volume in Mibps.

* `volume_path` - (Required) A unique file path for the volume. Changing this forces a new Application Volume Group to be created and data will be lost.

* `volume_spec_name` - (Required) Volume specification name. Possible values are `data`, `log`, `shared`, `data-backup` and `log-backup`. Changing this forces a new Application Volume Group to be created and data will be lost.

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Volume Group.

* `export_policy_rule` - (Required) One or more `export_policy_rule` blocks as defined below.

* `data_protection_replication` - (Optional) A `data_protection_replication` block as defined below. Changing this forces a new Application Volume Group to be created and data will be lost.

* `data_protection_snapshot_policy` - (Optional) A `data_protection_snapshot_policy` block as defined below.

---

A `data_protection_replication` block is used when enabling the Cross-Region Replication (CRR) data protection option by deploying two Azure NetApp Files Volumes, one to be a primary volume and the other one will be the secondary, the secondary will have this block and will reference the primary volume, not all volume spec types are supported, please refer to  [Configure application volume groups for the SAP HANA REST API](https://learn.microsoft.com/en-us/azure/azure-netapp-files/configure-application-volume-group-sap-hana-api) for detauls. Each volume must be in a supported [region pair](https://docs.microsoft.com/azure/azure-netapp-files/cross-region-replication-introduction#supported-region-pairs).

This block supports the following:

* `remote_volume_location` - (Required) Location of the primary volume.

* `remote_volume_resource_id` - (Required) Resource ID of the primary volume.

* `replication_frequency` - (Required) eplication frequency. Possible values are `10minutes`, `daily` and `hourly`.

* `endpoint_type` - (Optional) The endpoint type. Possible values are `dst` and `src`. Defaults to `dst`.

---

A `data_protection_snapshot_policy` block supports the following:

* `snapshot_policy_id` - (Required) Resource ID of the snapshot policy to apply to the volume.

---

A `export_policy_rule` block supports the following:

* `allowed_clients` - (Required) A comma-sperated list of allowed client IPv4 addresses.

* `nfsv3_enabled` - (Required) Enables NFSv3. Please note that this cannot be enabled if volume has NFSv4.1 as its protocol.

* `nfsv41_enabled` - (Required) Enables NFSv4.1. Please note that this cannot be enabled if volume has NFSv3 as its protocol.

* `root_access_enabled` - (Optional) Is root access permitted to this volume? Defaults to `true`.

* `rule_index` - (Required) The index number of the rule, must start at 1 and maximum 5.

* `unix_read_only` - (Optional) Is the file system on unix read only? Defaults to `false.

* `unix_read_write` - (Optional) Is the file system on unix read and write? Defaults to `true`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Volume Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Application Volume Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Volume Group.
* `update` - (Defaults to 2 hours) Used when updating the Application Volume Group.
* `delete` - (Defaults to 2 hours) Used when deleting the Application Volume Group.

## Import

Application Volume Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_volume_group_sap_hana.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mytest-rg/providers/Microsoft.NetApp/netAppAccounts/netapp-account-test/volumeGroups/netapp-volumegroup-test
```
