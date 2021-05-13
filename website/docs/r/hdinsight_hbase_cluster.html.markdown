---
subcategory: "HDInsight"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_hbase_cluster"
description: |-
  Manages a HDInsight HBase Cluster.
---

# azurerm_hdinsight_hbase_cluster

Manages a HDInsight HBase Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "hdinsightstor"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "hdinsight"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_hdinsight_hbase_cluster" "example" {
  name                = "example-hdicluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    hbase = "1.1"
  }

  gateway {
    enabled  = true
    username = "acctestusrgw"
    password = "TerrAform123!"
  }

  storage_account {
    storage_container_id = azurerm_storage_container.example.id
    storage_account_key  = azurerm_storage_account.example.primary_access_key
    is_default           = true
  }

  roles {
    head_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_D3_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
    }

    zookeeper_node {
      vm_size  = "Standard_D3_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name for this HDInsight HBase Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which this HDInsight HBase Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region which this HDInsight HBase Cluster should exist. Changing this forces a new resource to be created.

* `cluster_version` - (Required) Specifies the Version of HDInsights which should be used for this Cluster. Changing this forces a new resource to be created.

* `component_version` - (Required) A `component_version` block as defined below.

* `gateway` - (Required) A `gateway` block as defined below.

* `roles` - (Required) A `roles` block as defined below.

* `storage_account` - (Required) One or more `storage_account` block as defined below.

* `storage_account_gen2` - (Required) A `storage_account_gen2` block as defined below.

* `tier` - (Required) Specifies the Tier which should be used for this HDInsight HBase Cluster. Possible values are `Standard` or `Premium`. Changing this forces a new resource to be created.

* `min_tls_version` - (Optional) The minimal supported TLS version. Possible values are 1.0, 1.1 or 1.2. Changing this forces a new resource to be created.

~> **NOTE:** Starting on June 30, 2020, Azure HDInsight will enforce TLS 1.2 or later versions for all HTTPS connections. For more information, see [Azure HDInsight TLS 1.2 Enforcement](https://azure.microsoft.com/en-us/updates/azure-hdinsight-tls-12-enforcement/).

---

* `tags` - (Optional) A map of Tags which should be assigned to this HDInsight HBase Cluster.

* `metastores` - (Optional) A `metastores` block as defined below.

* `monitor` - (Optional) A `monitor` block as defined below.

---

A `component_version` block supports the following:

* `hbase` - (Required) The version of HBase which should be used for this HDInsight HBase Cluster. Changing this forces a new resource to be created.

---

A `gateway` block supports the following:

* `enabled` - (Optional/ **Deprecated) Is the Ambari portal enabled? The HDInsight API doesn't support disabling gateway anymore.

* `password` - (Required) The password used for the Ambari Portal.

-> **NOTE:** This password must be different from the one used for the `head_node`, `worker_node` and `zookeeper_node` roles.

* `username` - (Required) The username used for the Ambari Portal. Changing this forces a new resource to be created.

---

A `head_node` block supports the following:

* `username` - (Required) The Username of the local administrator for the Head Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Head Nodes. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Head Nodes. Changing this forces a new resource to be created.

-> **NOTE:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Head Nodes. Changing this forces a new resource to be created.

-> **NOTE:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Head Nodes should be provisioned within. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Head Nodes should be provisioned within. Changing this forces a new resource to be created.

---

A `roles` block supports the following:

* `head_node` - (Required) A `head_node` block as defined above.

* `worker_node` - (Required) A `worker_node` block as defined below.

* `zookeeper_node` - (Required) A `zookeeper_node` block as defined below.

---

A `storage_account` block supports the following:

* `is_default` - (Required) Is this the Default Storage Account for the HDInsight Hadoop Cluster? Changing this forces a new resource to be created.

-> **NOTE:** One of the `storage_account` or `storage_account_gen2` blocks must be marked as the default.

* `storage_account_key` - (Required) The Access Key which should be used to connect to the Storage Account. Changing this forces a new resource to be created.

* `storage_container_id` - (Required) The ID of the Storage Container. Changing this forces a new resource to be created.

-> **NOTE:** This can be obtained from the `id` of the `azurerm_storage_container` resource.

---

A `storage_account_gen2` block supports the following:

* `is_default` - (Required) Is this the Default Storage Account for the HDInsight Hadoop Cluster? Changing this forces a new resource to be created.

-> **NOTE:** One of the `storage_account` or `storage_account_gen2` blocks must be marked as the default.

* `storage_resource_id` - (Required) The ID of the Storage Account. Changing this forces a new resource to be created.

* `filesystem_id` - (Required) The ID of the Gen2 Filesystem. Changing this forces a new resource to be created.

* `managed_identity_resource_id` - (Required) The ID of Managed Identity to use for accessing the Gen2 filesystem. Changing this forces a new resource to be created.

-> **NOTE:** This can be obtained from the `id` of the `azurerm_storage_container` resource.

---

A `worker_node` block supports the following:

* `username` - (Required) The Username of the local administrator for the Worker Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Worker Nodes. Changing this forces a new resource to be created.

* `min_instance_count` - (Optional / **Deprecated** ) The minimum number of instances which should be run for the Worker Nodes. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Worker Nodes. Changing this forces a new resource to be created.

-> **NOTE:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Worker Nodes. Changing this forces a new resource to be created.

-> **NOTE:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Worker Nodes should be provisioned within. Changing this forces a new resource to be created.

* `target_instance_count` - (Optional) The number of instances which should be run for the Worker Nodes.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Worker Nodes should be provisioned within. Changing this forces a new resource to be created.

* `autoscale` - (Optional) A `autoscale` block as defined below.

---

A `zookeeper_node` block supports the following:

* `username` - (Required) The Username of the local administrator for the Zookeeper Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Zookeeper Nodes. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Zookeeper Nodes. Changing this forces a new resource to be created.

-> **NOTE:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Zookeeper Nodes. Changing this forces a new resource to be created.

-> **NOTE:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Zookeeper Nodes should be provisioned within. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Zookeeper Nodes should be provisioned within. Changing this forces a new resource to be created.

--- 

A `metastores` block supports the following:

* `hive` - (Optional) A `hive` block as defined below.

* `oozie` - (Optional) An `oozie` block as defined below.

* `ambari` - (Optional) An `ambari` block as defined below.

---

A `hive` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Hive metastore.  Changing this forces a new resource to be created.

* `database_name` - (Required) The external Hive metastore's existing SQL database.  Changing this forces a new resource to be created.

* `username` - (Required) The external Hive metastore's existing SQL server admin username.  Changing this forces a new resource to be created.

* `password` - (Required) The external Hive metastore's existing SQL server admin password.  Changing this forces a new resource to be created.


---

An `oozie` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Oozie metastore.  Changing this forces a new resource to be created.

* `database_name` - (Required) The external Oozie metastore's existing SQL database.  Changing this forces a new resource to be created.

* `username` - (Required) The external Oozie metastore's existing SQL server admin username.  Changing this forces a new resource to be created.

* `password` - (Required) The external Oozie metastore's existing SQL server admin password.  Changing this forces a new resource to be created.

---

An `ambari` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Ambari metastore.  Changing this forces a new resource to be created.

* `database_name` - (Required) The external Hive metastore's existing SQL database.  Changing this forces a new resource to be created.

* `username` - (Required) The external Ambari metastore's existing SQL server admin username.  Changing this forces a new resource to be created.

* `password` - (Required) The external Ambari metastore's existing SQL server admin password.  Changing this forces a new resource to be created.

---

A `monitor` block supports the following:

* `log_analytics_workspace_id` - (Required) The Operations Management Suite (OMS) workspace ID.

* `primary_key` - (Required) The Operations Management Suite (OMS) workspace key.

---

An `autoscale` block supports the following:

* `recurrence` - (Required) A `recurrence` block as defined below.

-> **NOTE:** Capacity based autoscaling isn't supported to HBase clusters.

---

A `recurrence` block supports the following:

* `schedule` - (Required) A list of `schedule` blocks as defined below.

* `timezone` - (Required) The time zone for the autoscale schedule times.

---

A `schedule` block supports the following:

* `days` - (Required) The days of the week to perform autoscale.

* `target_instance_count` - (Required) The number of worker nodes to autoscale at the specified time.

* `time` - (Required) The time of day to perform the autoscale in 24hour format.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HDInsight HBase Cluster.

* `https_endpoint` - The HTTPS Connectivity Endpoint for this HDInsight HBase Cluster.

* `ssh_endpoint` - The SSH Connectivity Endpoint for this HDInsight HBase Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the HBase HDInsight Cluster.
* `update` - (Defaults to 60 minutes) Used when updating the HBase HDInsight Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the HBase HDInsight Cluster.
* `delete` - (Defaults to 60 minutes) Used when deleting the HBase HDInsight Cluster.

## Import

HDInsight HBase Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hdinsight_hbase_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.HDInsight/clusters/cluster1
```
