---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_cloud_vm_cluster"
description: |-
  Manages a Cloud VM Cluster.
---

# azurerm_oracle_cloud_vm_cluster

Manages a Cloud VM Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_oracle_exadata_infrastructure" "example" {
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

data "azurerm_oracle_db_servers" "example" {
  resource_group_name               = azurerm_resource_group.example.name
  cloud_exadata_infrastructure_name = azurerm_oracle_exadata_infrastructure.example.name
}

resource "azurerm_oracle_cloud_vm_cluster" "example" {
  name                            = "example-cloud-vm-cluster"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  gi_version                      = "23.0.0.0"
  virtual_network_id              = azurerm_virtual_network.example.id
  license_model                   = "BringYourOwnLicense"
  db_servers                      = [for obj in data.azurerm_oracle_db_servers.example.db_servers : obj.ocid]
  ssh_public_keys                 = [file("~/.ssh/id_rsa.pub")]
  display_name                    = "example-cloud-vm-cluster"
  cloud_exadata_infrastructure_id = azurerm_oracle_exadata_infrastructure.example.id
  cpu_core_count                  = 2
  hostname                        = "hostname"
  subnet_id                       = azurerm_subnet.example.id
  system_version                  = "23.1.19.0.0.241015"
}
```

## Arguments Reference

The following arguments are supported:

* `cloud_exadata_infrastructure_id` - (Required) The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Cloud Exadata infrastructure. Changing this forces a new Cloud VM Cluster to be created.

* `cpu_core_count` - (Required) The number of CPU cores enabled on the Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

* `db_servers` - (Required) The list of DB servers. Changing this forces a new Cloud VM Cluster to be created.

* `display_name` - (Required) The user-friendly name for the Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created. The name does not need to be unique.

* `gi_version` - (Required) A valid Oracle Grid Infrastructure (GI) software version. Changing this forces a new Cloud VM Cluster to be created.

* `hostname` - (Required) The hostname for the Cloud VM Cluster without suffix. Changing this forces a new Cloud VM Cluster to be created.

* `license_model` - (Required) The Oracle license model that applies to the Cloud VM Cluster, either `BringYourOwnLicense` or `LicenseIncluded`. Changing this forces a new Cloud VM Cluster to be created.

* `location` - (Required) The Azure Region where the Cloud VM Cluster should exist. Changing this forces a new Cloud VM Cluster to be created.

* `name` - (Required) The name which should be used for this Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cloud VM Cluster should exist. Changing this forces a new Cloud VM Cluster to be created.

* `ssh_public_keys` - (Required) The public key portion of one or more key pairs used for SSH access to the Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

* `subnet_id` - (Required) The ID of the subnet associated with the Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network associated with the Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

---

* `backup_subnet_cidr` - (Optional) The backup subnet CIDR of the Virtual Network associated with the Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

* `cluster_name` - (Optional) The cluster name for Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

* `data_collection_options` - (Optional) A `data_collection_options` block as defined below. Changing this forces a new Cloud VM Cluster to be created.

* `data_storage_percentage` - (Optional) The percentage assigned to DATA storage (user data and database files). Changing this forces a new Cloud VM Cluster to be created. The remaining percentage is assigned to RECO storage (database redo logs, archive logs, and recovery manager backups). Accepted values are `35`, `40`, `60` and `80`.

* `data_storage_size_in_tbs` - (Optional) The data disk group size to be allocated in TBs. Changing this forces a new Cloud VM Cluster to be created.

* `db_node_storage_size_in_gbs` - (Optional) The local node storage to be allocated in GBs. Changing this forces a new Cloud VM Cluster to be created.

* `domain` - (Optional) The name of the OCI Private DNS Zone to be associated with the Cloud VM Cluster. This is required for specifying your own private domain name. Changing this forces a new Cloud VM Cluster to be created.

* `local_backup_enabled` - (Optional)  If true, database backup on local Exadata storage is configured for the Cloud VM Cluster. If `false`, database backup on local Exadata storage is not available in the Cloud VM Cluster. Changing this forces a new Cloud VM Cluster to be created.

* `sparse_diskgroup_enabled` - (Optional) If true, the sparse disk group is configured for the Cloud VM Cluster. If `false`, the sparse disk group is not created. Changing this forces a new Cloud VM Cluster to be created.

* `memory_size_in_gbs` - (Optional) The memory to be allocated in GBs. Changing this forces a new Cloud VM Cluster to be created.

* `scan_listener_port_tcp` - (Optional) The TCP Single Client Access Name (SCAN) port. The default port to 1521. Changing this forces a new Cloud VM Cluster to be created.

* `scan_listener_port_tcp_ssl` - (Optional) The TCPS Single Client Access Name (SCAN) port. The default port to 2484. Changing this forces a new Cloud VM Cluster to be created.

* `system_version` - (Optional) Operating system version of the Exadata image. System version must be <= Db server major version (the first two parts of the DB server version eg 23.1.X.X.XXXX). Accepted Values for Grid Infrastructure (GI) version 19.0.0.0 are 22.1.30.0.0.241204, 22.1.32.0.0.250205, 22.1.31.0.0.250110, 23.1.20.0.0.241112, 23.1.21.0.0.241204, 23.1.22.0.0.250119, 23.1.23.0.0.250207. For Grid Infrastructure (GI) version 23.0.0.0 allowed system versions are 23.1.19.0.0.241015, 23.1.20.0.0.241112, 23.1.22.0.0.250119, 23.1.21.0.0.241204, 23.1.23.0.0.250207.

* `tags` - (Optional) A mapping of tags which should be assigned to the Cloud VM Cluster.

* `time_zone` - (Optional) The time zone of the Cloud VM Cluster. For details, see [Exadata Infrastructure Time Zones](https://docs.cloud.oracle.com/iaas/Content/Database/References/timezones.htm). Changing this forces a new Cloud VM Cluster to be created.

* `zone_id` - (Optional) The OCID of the OCI Private DNS Zone to be associated with the Cloud VM Cluster. This is required for specifying your own private domain name. Changing this forces a new Cloud VM Cluster to be created.

---

A `data_collection_options` block supports the following:

* `diagnostics_events_enabled` - (Optional) Indicates whether diagnostic collection is enabled for the VM Cluster/Cloud VM Cluster/VMBM DBCS. Enabling diagnostic collection allows you to receive Events service notifications for guest VM issues. Diagnostic collection also allows Oracle to provide enhanced service and proactive support for your Exadata system. You can enable diagnostic collection during VM Cluster/Cloud VM Cluster provisioning. You can also disable or enable it at any time using the `UpdateVmCluster` or `updateCloudVmCluster` API. Changing this forces a new Cloud VM Cluster to be created.

* `health_monitoring_enabled` - (Optional) Indicates whether health monitoring is enabled for the VM Cluster / Cloud VM Cluster / VMBM DBCS. Enabling health monitoring allows Oracle to collect diagnostic data and share it with its operations and support personnel. You may also receive notifications for some events. Collecting health diagnostics enables Oracle to provide proactive support and enhanced service for your system. Optionally enable health monitoring while provisioning a system. You can also disable or enable health monitoring anytime using the `UpdateVmCluster`, `UpdateCloudVmCluster` or `updateDbsystem` API. Changing this forces a new Cloud VM Cluster to be created.

* `incident_logs_enabled` - (Optional) Indicates whether incident logs and trace collection are enabled for the VM Cluster / Cloud VM Cluster / VMBM DBCS. Enabling incident logs collection allows Oracle to receive Events service notifications for guest VM issues, collect incident logs and traces, and use them to diagnose issues and resolve them. Optionally enable incident logs collection while provisioning a system. You can also disable or enable incident logs collection anytime using the `UpdateVmCluster`, `updateCloudVmCluster` or `updateDbsystem` API. Changing this forces a new Cloud VM Cluster to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cloud VM Cluster.

* `hostname_actual` - The hostname for the Cloud VM Cluster with suffix.

* `ocid` - The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Cloud VM Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Cloud VM Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cloud VM Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Cloud VM Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cloud VM Cluster.

## Import

Cloud VM Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_cloud_vm_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup/providers/Oracle.Database/cloudVmClusters/cloudVmClusters1
```
