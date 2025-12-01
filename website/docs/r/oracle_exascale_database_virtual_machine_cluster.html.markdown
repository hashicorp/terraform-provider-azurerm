---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_exascale_database_virtual_machine_cluster"
description: |-
  Manages a Exadata VM Cluster.
---

# azurerm_oracle_exascale_database_virtual_machine_cluster

Manages a Exadata VM Cluster on Exascale Infrastructure.

## Example Usage

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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

resource "azurerm_oracle_exascale_database_storage_vault" "example" {
  name                = "example-exascale-db-storage-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  display_name        = "ExampleExascaleDbStorage"
  description         = "description"
  high_capacity_database_storage {
    total_size_in_gb = 300
  }
  additional_flash_cache_percentage = 100
  zones                             = ["3"]
}

resource "azurerm_oracle_exascale_database_virtual_machine_cluster" "example" {
  name                = "example-exadb-vm-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  exascale_database_storage_vault_id = azurerm_oracle_exascale_database_storage_vault.example.id
  display_name                       = "ExampleExascaleVmCluster"
  enabled_ecpu_count                 = 4
  hostname                           = "host"
  node_count                         = 2
  shape                              = "EXADBXS"
  ssh_public_keys                    = [file("~/.ssh/id_rsa.pub")]
  subnet_id                          = azurerm_subnet.example.id
  total_ecpu_count                   = 8
  virtual_machine_file_system_storage {
    total_size_in_gb = 440
  }
  virtual_network_id = azurerm_virtual_network.example.id
  zones              = ["3"]
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Exadata VM Cluster. Changing this forces a new Exadata VM Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Exadata VM Cluster should exist. Changing this forces a new Exadata VM Cluster to be created.

* `location` - (Required) The Azure Region where the Exadata VM Cluster should exist. Changing this forces a new Exadata VM Cluster to be created.

* `display_name` - (Required) The user-friendly name for the Exadata VM Cluster. Changing this forces a new Exadata VM Cluster to be created. The name does not need to be unique.

* `enabled_ecpu_count` - (Required) The number of ECPUs to enable for an Exadata VM cluster on Exascale Infrastructure. Possible values range between `8` and `200`, and must be divisible by `4`. Changing this forces a new resource to be created.

* `exascale_database_storage_vault_id` - (Required) The Azure Resource ID of the Exadata Database Storage Vault. Changing this forces a new Exadata VM Cluster to be created.

* `hostname` - (Required) The hostname for the Exadata VM Cluster on Exascale Infrastructure. Changing this forces a new Exadata VM Cluster to be created.

* `grid_image_ocid` - (Required) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the grid image that is used for grid setup. Changing this forces a new resource to be created.

* `node_count` - (Required) The number of nodes in the Exadata VM cluster on Exascale Infrastructure. Possible values range between `2` and `10`.

* `shape` - (Required) The shape of the Exadata VM cluster on Exascale Infrastructure resource.Possible values is `EXADBXS`.  Changing this forces a new resource to be created.

* `ssh_public_keys` - (Required) The public key portion of one or more key pairs used for SSH access to the Exadata VM Cluster. Changing this forces a new Exadata VM Cluster to be created.

* `subnet_id` - (Required) The ID of the subnet associated with the Exadata VM Cluster. This subnet must belong to the specified virtual_network_id. Changing this value forces a new Exadata VM Cluster to be created.

* `total_ecpu_count` - (Required) The number of Total ECPUs for an Exadata VM cluster on Exascale Infrastructure. Possible values range between `8` and `200`, and must be divisible by `4`. Changing this forces a new resource to be created.

* `virtual_machine_file_system_storage` - (Required) A `virtual_machine_file_system_storage` block as defined below. Changing this forces a new resource to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network associated with the Exadata VM Cluster that the specified subnet belongs to. Changing this value forces a new Exadata VM Cluster to be created.

* `zones` - (Required) Specifies a list of Availability Zones in which this Exadata VM Cluster should be located. Changing this forces a new Exadata VM Cluster to be created.

---

* `backup_subnet_cidr` - (Optional) The backup subnet CIDR of the Virtual Network associated with the Exadata VM Cluster. Changing this forces a new Exadata VM Cluster to be created.

* `cluster_name` - (Optional) The cluster name for Exadata VM Cluster. Changing this forces a new Exadata VM Cluster to be created.

* `data_collection` - (Optional) A `data_collection` block as defined below. Changing this forces a new Exadata VM Cluster to be created.

* `domain` - (Optional) The name of the OCI Private DNS Zone to be associated with the Exadata VM Cluster. This is required for specifying your own private domain name. Changing this forces a new Exadata VM Cluster to be created.

* `license_model` - (Optional) The Oracle license model that applies to the Exadata VM Cluster. Possible values are `BringYourOwnLicense` and `LicenseIncluded`. Changing this forces a new Exadata VM Cluster to be created.

* `inbound_network_security_group_rule` - (Optional) A `inbound_network_security_group_rule` block as defined below. Changing this forces a new Exadata VM Cluster to be created.

* `private_zone_ocid` - (Optional) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the private zone in which you want DNS records to be created. Changing this forces a new Exadata VM Cluster to be created.

* `single_client_access_name_listener_port_tcp` - (Optional) The TCP Single Client Access Name (SCAN) port. Possible values range between `1024` and `8999`. Defaults to `1521`. Changing this forces a new Exadata VM Cluster to be created.

* `single_client_access_name_listener_port_tcp_ssl` - (Optional) The TCPS Single Client Access Name (SCAN) port. Possible values range between `1024` and `8999`. Defaults to `2484`. Changing this forces a new Exadata VM Cluster to be created.

* `system_version` - (Optional) Operating system version of the Exadata image. Changing this forces a new resource to be created.

~> **Note:** `system_version` must be less than or equal to the database server's major version (the first two parts of the database server version, e.g. 23.1.X.X.XXXX).

* `tags` - (Optional) A mapping of tags which should be assigned to the Exadata VM Cluster.

 `time_zone` - (Optional) The time zone of the Exadata VM Cluster. Changing this forces a new Exadata VM Cluster to be created.

 -> **Note:** For more information on `time_zone`, see [Exadata Infrastructure Time Zones](https://docs.cloud.oracle.com/iaas/Content/Database/References/timezones.htm).

---

A `virtual_machine_file_system_storage` block supports the following:

* `total_size_in_gb` - (Required) Total Capacity. Changing this forces a new resource to be created.

---

A `data_collection` block supports the following:

* `diagnostics_events_enabled` - (Optional) Indicates whether diagnostic collection is enabled for the VM cluster. Defaults to `false`. Changing this forces a new Exadata VM Cluster to be created.

* `health_monitoring_enabled` - (Optional) Indicates whether health monitoring is enabled for the VM cluster. Defaults to `false`. Changing this forces a new Exadata VM Cluster to be created.

* `incident_logs_enabled` - (Optional) Indicates whether incident logs and trace collection are enabled for the VM cluster. Defaults to `false`. Changing this forces a new Exadata VM Cluster to be created.

---

A `inbound_network_security_group_rule` block supports the following:

* `destination_port_range` - (Required) A `destination_port_range` block as defined below.

* `source_ip_range` - (Required) It is a range of IP addresses that a packet coming into the instance can come from. Changing this forces a new resource to be created.

---

A `destination_port_range` block supports the following:

* `maximum` - (Required) The maximum port number, which must not be less than the minimum port number. Possible values range between `1` and `65535`. Changing this forces a new resource to be created.

* `minimum` - (Required) The minimum port number, which must not be greater than the maximum port number. Possible values range between `1` and `65535`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Exadata VM Cluster.

* `ocid` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Exadata VM Cluster.

* `zone_ocid` - The [OCID](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/identifiers.htm) of the zone the Exadata VM Cluster is associated with.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Exadata VM Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Exadata VM Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Exadata VM Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Exadata VM Cluster.

## Import

Exadata VM Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_exascale_database_virtual_machine_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/exadbVmClusters/exadbVmClusters1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
