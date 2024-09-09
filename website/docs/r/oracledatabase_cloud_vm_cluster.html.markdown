---
subcategory: "App Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracledatabase_cloud_vm_cluster"
description: |-
  Manages a Cloud VM Cluster.
---

# azurerm_oracledatabase_cloud_vm_cluster

Manages a Cloud VM Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_oracledatabase_exadata_infrastructure" "example" {
  name                = "example-exadata-infrastructure"
  display_name        = "example-exadata-infrastructure"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  shape               = "Exadata.X9M"
  storage_count       = "3"
  compute_count       = "2"
  zones               = ["3"]
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

data "azurerm_oracledatabase_db_servers" "example" {
  resource_group_name               = azurerm_resource_group.example.name
  cloud_exadata_infrastructure_name = azurerm_oracledatabase_exadata_infrastructure.example.name
}

resource "azurerm_oracledatabase_cloud_vm_cluster" "example" {
  name                            = "example-cloud-vm-cluster"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  gi_version                      = "23.0.0.0"
  vnet_id                         = azurerm_virtual_network.example.id
  license_model                   = "BringYourOwnLicense"
  db_servers                      = [for obj in data.azurerm_oracledatabase_db_servers.example.db_servers : obj.ocid]
  ssh_public_keys                 = [file("~/.ssh/id_rsa.pub")]
  display_name                    = "example-cloud-vm-cluster"
  cloud_exadata_infrastructure_id = azurerm_oracledatabase_exadata_infrastructure.example.id
  cpu_core_count                  = 2
  hostname                        = "hostname"
  subnet_id                       = azurerm_subnet.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `cloud_exadata_infrastructure_id` - (Required) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the cloud Exadata infrastructure.

* `cpu_core_count` - (Required) The number of CPU cores enabled on the cloud VM cluster.

* `db_servers` - (Required) The list of DB servers.

* `display_name` - (Required) The user-friendly name for the cloud VM cluster. The name does not need to be unique..

* `gi_version` - (Required) A valid Oracle Grid Infrastructure (GI) software version.

* `hostname` - (Required) The hostname for the cloud VM cluster.

* `license_model` - (Required) The Oracle license model that applies to the cloud VM cluster, either `BringYourOwnLicense` or `LicenseIncluded`.

* `location` - (Required) The Azure Region where the Cloud VM Cluster should exist.

* `name` - (Required) The name which should be used for this Cloud VM Cluster.

* `resource_group_name` - (Required) The name of the Resource Group where the Cloud VM Cluster should exist.

* `ssh_public_keys` - (Required) The public key portion of one or more key pairs used for SSH access to the cloud VM cluster.

* `subnet_id` - (Required) The ID of the subnet associated with the cloud VM cluster.

* `vnet_id` - (Required) The ID of the vnet associated with the cloud VM cluster.

---

* `backup_subnet_cidr` - (Optional) The backup subnet CIDR of the vnet associated with the cloud VM cluster..

* `cluster_name` - (Optional) The cluster name for cloud VM cluster.

* `data_collection_options` - (Optional) A `data_collection_options` block as defined below.

* `data_storage_percentage` - (Optional) The percentage assigned to DATA storage (user data and database files). The remaining percentage is assigned to RECO storage (database redo logs, archive logs, and recovery manager backups). Accepted values are 35, 40, 60 and 80.

* `data_storage_size_in_tbs` - (Optional) The data disk group size to be allocated in TBs.

* `db_node_storage_size_in_gbs` - (Optional) The local node storage to be allocated in GBs.

* `is_local_backup_enabled` - (Optional)  If true, database backup on local Exadata storage is configured for the cloud VM cluster. If false, database backup on local Exadata storage is not available in the cloud VM cluster.

* `is_sparse_diskgroup_enabled` - (Optional) If true, the sparse disk group is configured for the cloud VM cluster. If false, the sparse disk group is not created.

* `memory_size_in_gbs` - (Optional) The memory to be allocated in GBs.

* `tags` - (Optional) A mapping of tags which should be assigned to the Cloud VM Cluster.

* `time_zone` - (Optional) The time zone of the cloud VM cluster. For details, see [Exadata Infrastructure Time Zones](https://docs.cloud.oracle.com/iaas/Content/Database/References/timezones.htm).

---

A `data_collection_options` block supports the following:

* `is_diagnostics_events_enabled` - (Optional) Indicates whether diagnostic collection is enabled for the VM cluster/Cloud VM cluster/VMBM DBCS. Enabling diagnostic collection allows you to receive Events service notifications for guest VM issues. Diagnostic collection also allows Oracle to provide enhanced service and proactive support for your Exadata system. You can enable diagnostic collection during VM cluster/Cloud VM cluster provisioning. You can also disable or enable it at any time using the `UpdateVmCluster` or `updateCloudVmCluster` API.

* `is_health_monitoring_enabled` - (Optional) Indicates whether health monitoring is enabled for the VM cluster / Cloud VM cluster / VMBM DBCS. Enabling health monitoring allows Oracle to collect diagnostic data and share it with its operations and support personnel. You may also receive notifications for some events. Collecting health diagnostics enables Oracle to provide proactive support and enhanced service for your system. Optionally enable health monitoring while provisioning a system. You can also disable or enable health monitoring anytime using the `UpdateVmCluster`, `UpdateCloudVmCluster` or `updateDbsystem` API.

* `is_incident_logs_enabled` - (Optional) Indicates whether incident logs and trace collection are enabled for the VM cluster / Cloud VM cluster / VMBM DBCS. Enabling incident logs collection allows Oracle to receive Events service notifications for guest VM issues, collect incident logs and traces, and use them to diagnose issues and resolve them. Optionally enable incident logs collection while provisioning a system. You can also disable or enable incident logs collection anytime using the `UpdateVmCluster`, `updateCloudVmCluster` or `updateDbsystem` API.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cloud VM Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cloud VM Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cloud VM Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Cloud VM Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cloud VM Cluster.

## Import

Cloud VM Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracledatabase_cloud_vm_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/cloudVmClusters/cloudVmClusters1
```
