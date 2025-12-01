---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_database_system"
description: |-
  Manages a Database System.
---

# azurerm_oracle_database_system

Manages a Database System.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}


resource "azurerm_virtual_network" "example" {
  name                = "example-virtual-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Oracle.Database/networkAttachments"
    }
  }
}

resource "azurerm_oracle_resource_anchor" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_oracle_network_anchor" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = "eastus"
  resource_anchor_id  = azurerm_oracle_resource_anchor.example.id
  subnet_id           = azurerm_subnet.example.id
  zones               = ["2"]
}

resource "azurerm_oracle_database_system" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = "eastus"
  zones               = ["1"]
  cpu_core_count      = 4
  source              = "None"
  database_edition    = "StandardEdition"
  database_version    = "19.27.0.0"
  hostname            = "hostname"
  network_anchor_id   = azurerm_oracle_network_anchor.example.id
  resource_anchor_id  = azurerm_oracle_resource_anchor.example.id
  shape               = "VM.Standard.x86"
  ssh_public_keys     = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC+wWK73dCr+jgQOAxNsHAnNNNMEMWOHYEccp6wJm2gotpr9katuF/ZAdou5AaW1C61slRkHRkpRRX9FA9CYBiitZgvCCz+3nWNN7l/Up54Zps/pHWGZLHNJZRYyAB6j5yVLMVHIHriY49d/GZTZVNB8GoJv9Gakwc/fuEZYYl4YDFiGMBP///TzlI4jhiJzjKnEvqPFki5p2ZRJqcbCiF4pJrxUQR/RXqVFQdbRLZgYfJ8xGB878RENq3yQ39d8dVOkq4edbkzwcUmwwwkYVPIoDGsYLaRHnG+To7FvMeyO7xDVQkMKzopTQV8AuKpyvpqu0a9pWOMaiCyDytO7GGN you@me.com"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Database system. Changing this forces a new Database system to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Database system should exist. Changing this forces a new Database system to be created.

* `location` - (Required) The Azure Region where the Database System should exist. Changing this forces a new Database system to be created.

* `admin_password` - (Required) A strong password for SYS, SYSTEM, and PDB Admin. The password must be at least nine characters and contain at least two uppercase, two lowercase, two numbers, and two special characters. The special characters must be _, #, or -. Changing this forces a new Database system to be created.

* `compute_count` - (Required) The number of compute servers for the Database system. Changing this forces a new Database system to be created.

* `compute_model` - (Required) The compute model for Base Database Service. This is required if using the `computeCount` parameter. If using `cpuCoreCount` then it is an error to specify `computeModel` to a non-null value. The ECPU compute model is the recommended model, and the OCPU compute model is legacy. Changing this forces a new Database system to be created.

* `database_edition` - (Required) The Oracle Database Edition that applies to all the databases on the Database system. Exadata Database systems and 2-node RAC Database systems require EnterpriseEditionExtremePerformance. Possible values are `EnterpriseEdition`, `EnterpriseEditionDeveloper`, `EnterpriseEditionExtreme`, `EnterpriseEditionHighPerformance` and `StandardEdition`. Changing this forces a new Database System to be created.

* `database_version` - (Required) A valid Oracle Database version. For a list of supported versions, use the ListDbVersions operation. Changing this forces a new Database System to be created.

* `hostname` - (Required) The hostname for the Database system. Changing this forces a new Database system to be created.

* `license_model` - (Required) The Oracle license model that applies to all the databases on the Database system. The default is LicenseIncluded. Possible values are `BringYourOwnLicense` and `LicenseIncluded`. Changing this forces a new Database system to be created.

* `network_anchor_id` - (Required) The ID of the Azure Network Anchor. Changing this forces a new Database system to be created.

* `resource_anchor_id` - (Required) The ID of the Azure Resource Anchor. Changing this forces a new Database system to be created.

* `shape` - (Required) The shape of the Database system. The shape determines resources to allocate to the Database system. For virtual machine shapes, the number of CPU cores and memory. For bare metal and Exadata shapes, the number of CPU cores, storage, and memory. The only possible value is `VM.Standard.x86`. Changing this forces a new Database system to be created.

* `source` - (Required) The source of the database: Use NONE for creating a new database. The default is `None`. Changing this forces a new Database System to be created.

* `ssh_public_keys` - (Required) The public key portion of one or more key pairs used for SSH access to the Database system. Changing this forces a new Database system to be created.

* `zones` - (Required) Specifies a list of Database System zones. Changing this forces a new Database System to be created.

---

* `cluster_name` - (Optional) The cluster name for Exadata and 2-node RAC virtual machine Database systems. The cluster name must begin with an alphabetic character, and may contain hyphens (-). Underscores (_) are not permitted. The cluster name can be no longer than 11 characters and is not case sensitive. Changing this forces a new Database system to be created.

* `database_system_options` - (Optional) One or more `database_system_options` blocks as defined below. Changing this forces a new resource to be created.

* `disk_redundancy` - (Optional) The type of redundancy configured for the Database system.  Possible values are `High` and `Normal`. NORMAL is 2-way redundancy. HIGH is 3-way redundancy. Changing this forces a new Database system to be created.

* `display_name` - (Optional) The user-friendly name for the Database system. Changing this forces a new Database system to be created. The name does not need to be unique.

* `domain` - (Optional) The domain name for the Database system. Changing this forces a new Database system to be created.

* `initial_data_storage_size_in_gb` - (Optional) Size in GB of the initial data volume that will be created and attached to a virtual machine Database system. You can scale up storage after provisioning, as needed. Note that the total storage size attached will be more than the amount you specify to allow for REDO/RECO space and software volume. Changing this forces a new Database system to be created.

* `node_count` - (Optional) The number of nodes in the Database system. For RAC Database systems, the value is greater than 1. Changing this forces a new Database system to be created.

* `pluggable_database_name` - (Optional) The name of the pluggable database. The name must begin with an alphabetic character and can contain a maximum of thirty alphanumeric characters. Special characters are not permitted. The pluggable database name should not be the same as the database name. Changing this forces a new Database system to be created.

* `storage_volume_performance_mode` - (Optional) The block storage volume performance level. Valid values are Balanced and HighPerformance. See [Block Volume Performance](/Content/Block/Concepts/blockvolumeperformance.htm) for more information. Changing this forces a new Database system to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Database System.

* `time_zone` - (Optional) The time zone of the Database system, e.g., UTC, to set the timeZone as UTC. Changing this forces a new Database System to be created.

---

A `database_system_options` block supports the following:

* `storage_management` - (Optional) The storage option used in Database system. ASM - Automatic storage management, LVM - Logical Volume management.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Database System.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Database System.
* `read` - (Defaults to 5 minutes) Used when retrieving the Database System.
* `update` - (Defaults to 30 minutes) Used when updating the Database System.
* `delete` - (Defaults to 60 minutes) Used when deleting the Database System.

## Import

Database Systems can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_database_system.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Oracle.Database/dbSystems/example
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
