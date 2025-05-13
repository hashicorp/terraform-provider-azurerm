---
subcategory: "HDInsight"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_interactive_query_cluster"
description: |-
  Manages a HDInsight Interactive Query Cluster.
---

# azurerm_hdinsight_interactive_query_cluster

Manages a HDInsight Interactive Query Cluster.

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

resource "azurerm_hdinsight_interactive_query_cluster" "example" {
  name                = "example-hdicluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_version     = "3.6"
  tier                = "Standard"

  component_version {
    interactive_hive = "2.1"
  }

  gateway {
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
      vm_size  = "Standard_D13_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }

    worker_node {
      vm_size               = "Standard_D14_V2"
      username              = "acctestusrvm"
      password              = "AccTestvdSC4daf986!"
      target_instance_count = 3
    }

    zookeeper_node {
      vm_size  = "Standard_A4_V2"
      username = "acctestusrvm"
      password = "AccTestvdSC4daf986!"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name for this HDInsight Interactive Query Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which this HDInsight Interactive Query Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region which this HDInsight Interactive Query Cluster should exist. Changing this forces a new resource to be created.

* `cluster_version` - (Required) Specifies the Version of HDInsights which should be used for this Cluster. Changing this forces a new resource to be created.

* `component_version` - (Required) A `component_version` block as defined below.

* `encryption_in_transit_enabled` - (Optional) Whether encryption in transit is enabled for this Cluster. Changing this forces a new resource to be created.

* `disk_encryption` - (Optional) A `disk_encryption` block as defined below.

* `gateway` - (Required) A `gateway` block as defined below.

* `roles` - (Required) A `roles` block as defined below.

* `network` - (Optional) A `network` block as defined below.

* `private_link_configuration` - (Optional) A `private_link_configuration` block as defined below.

* `compute_isolation` - (Optional) A `compute_isolation` block as defined below.

* `storage_account` - (Optional) One or more `storage_account` block as defined below.

* `storage_account_gen2` - (Optional) A `storage_account_gen2` block as defined below.

* `tier` - (Required) Specifies the Tier which should be used for this HDInsight Interactive Query Cluster. Possible values are `Standard` or `Premium`. Changing this forces a new resource to be created.

* `tls_min_version` - (Optional) The minimal supported TLS version. Possible values are 1.0, 1.1 or 1.2. Changing this forces a new resource to be created.

~> **Note:** Starting on June 30, 2020, Azure HDInsight will enforce TLS 1.2 or later versions for all HTTPS connections. For more information, see [Azure HDInsight TLS 1.2 Enforcement](https://azure.microsoft.com/en-us/updates/azure-hdinsight-tls-12-enforcement/).

---

* `tags` - (Optional) A map of Tags which should be assigned to this HDInsight Interactive Query Cluster.

* `metastores` - (Optional) A `metastores` block as defined below.

* `monitor` - (Optional) A `monitor` block as defined below.

* `extension` - (Optional) An `extension` block as defined below.

* `security_profile` - (Optional) A `security_profile` block as defined below. Changing this forces a new resource to be created.

---

A `component_version` block supports the following:

* `interactive_hive` - (Required) The version of Interactive Query which should be used for this HDInsight Interactive Query Cluster. Changing this forces a new resource to be created.

---

A `gateway` block supports the following:

* `password` - (Required) The password used for the Ambari Portal.

-> **Note:** This password must be different from the one used for the `head_node`, `worker_node` and `zookeeper_node` roles.

* `username` - (Required) The username used for the Ambari Portal. Changing this forces a new resource to be created.

---

A `head_node` block supports the following:

* `username` - (Required) The Username of the local administrator for the Head Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Head Nodes. Possible values are `ExtraSmall`, `Small`, `Medium`, `Large`, `ExtraLarge`, `A5`, `A6`, `A7`, `A8`, `A9`, `A10`, `A11`, `Standard_A1_V2`, `Standard_A2_V2`, `Standard_A2m_V2`, `Standard_A3`, `Standard_A4_V2`, `Standard_A4m_V2`, `Standard_A8_V2`, `Standard_A8m_V2`, `Standard_D1`, `Standard_D2`, `Standard_D3`, `Standard_D4`, `Standard_D11`, `Standard_D12`, `Standard_D13`, `Standard_D14`, `Standard_D1_V2`, `Standard_D2_V2`, `Standard_D3_V2`, `Standard_D4_V2`, `Standard_D5_V2`, `Standard_D11_V2`, `Standard_D12_V2`, `Standard_D13_V2`, `Standard_D14_V2`, `Standard_DS1_V2`, `Standard_DS2_V2`, `Standard_DS3_V2`, `Standard_DS4_V2`, `Standard_DS5_V2`, `Standard_DS11_V2`, `Standard_DS12_V2`, `Standard_DS13_V2`, `Standard_DS14_V2`, `Standard_E2_V3`, `Standard_E4_V3`, `Standard_E8_V3`, `Standard_E16_V3`, `Standard_E20_V3`, `Standard_E32_V3`, `Standard_E64_V3`, `Standard_E64i_V3`, `Standard_E2s_V3`, `Standard_E4s_V3`, `Standard_E8s_V3`, `Standard_E16s_V3`, `Standard_E20s_V3`, `Standard_E32s_V3`, `Standard_E64s_V3`, `Standard_E64is_V3`, `Standard_D2a_V4`, `Standard_D4a_V4`, `Standard_D8a_V4`, `Standard_D16a_V4`, `Standard_D32a_V4`, `Standard_D48a_V4`, `Standard_D64a_V4`, `Standard_D96a_V4`, `Standard_E2a_V4`, `Standard_E4a_V4`, `Standard_E8a_V4`, `Standard_E16a_V4`, `Standard_E20a_V4`, `Standard_E32a_V4`, `Standard_E48a_V4`, `Standard_E64a_V4`, `Standard_E96a_V4`, `Standard_D2ads_V5`, `Standard_D4ads_V5`, `Standard_D8ads_V5`, `Standard_D16ads_V5`, `Standard_D32ads_V5`, `Standard_D48ads_V5`, `Standard_D64ads_V5`, `Standard_D96ads_V5`, `Standard_E2ads_V5`, `Standard_E4ads_V5`, `Standard_E8ads_V5`, `Standard_E16ads_V5`, `Standard_E20ads_V5`, `Standard_E32ads_V5`, `Standard_E48ads_V5`, `Standard_E64ads_V5`, `Standard_E96ads_V5`, `Standard_G1`, `Standard_G2`, `Standard_G3`, `Standard_G4`, `Standard_G5`, `Standard_F2s_V2`, `Standard_F4s_V2`, `Standard_F8s_V2`, `Standard_F16s_V2`, `Standard_F32s_V2`, `Standard_F64s_V2`, `Standard_F72s_V2`, `Standard_GS1`, `Standard_GS2`, `Standard_GS3`, `Standard_GS4`, `Standard_GS5` and `Standard_NC24`. Changing this forces a new resource to be created.

-> **Note:** High memory instances must be specified for the Head Node (Azure suggests a `Standard_D13_V2`).

* `password` - (Optional) The Password associated with the local administrator for the Head Nodes. Changing this forces a new resource to be created.

-> **Note:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Head Nodes. Changing this forces a new resource to be created.

-> **Note:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Head Nodes should be provisioned within. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Head Nodes should be provisioned within. Changing this forces a new resource to be created.

* `script_actions` - (Optional) The script action which will run on the cluster. One or more `script_actions` blocks as defined below.

---

A `script_actions` block supports the following:

* `name` - (Required) The name of the script action.

* `uri` - (Required) The URI to the script.

* `parameters` - (Optional) The parameters for the script provided.

---

A `roles` block supports the following:

* `head_node` - (Required) A `head_node` block as defined above.

* `worker_node` - (Required) A `worker_node` block as defined below.

* `zookeeper_node` - (Required) A `zookeeper_node` block as defined below.

---

A `network` block supports the following:

* `connection_direction` - (Optional) The direction of the resource provider connection. Possible values include `Inbound` or `Outbound`. Defaults to `Inbound`. Changing this forces a new resource to be created.

-> **Note:** To enabled the private link the `connection_direction` must be set to `Outbound`.

* `private_link_enabled` - (Optional) Is the private link enabled? Possible values include `true` or `false`. Defaults to `false`. Changing this forces a new resource to be created.

---

A `compute_isolation` block supports the following:

* `compute_isolation_enabled` - (Optional) This field indicates whether enable compute isolation or not. Possible values are `true` or `false`.

* `host_sku` - (Optional) The name of the host SKU.

---

A `storage_account` block supports the following:

* `is_default` - (Required) Is this the Default Storage Account for the HDInsight Hadoop Cluster? Changing this forces a new resource to be created.

-> **Note:** One of the `storage_account` or `storage_account_gen2` blocks must be marked as the default.

* `storage_account_key` - (Required) The Access Key which should be used to connect to the Storage Account. Changing this forces a new resource to be created.

* `storage_container_id` - (Required) The ID of the Storage Container. Changing this forces a new resource to be created.

-> **Note:** This can be obtained from the `id` of the `azurerm_storage_container` resource.

* `storage_resource_id` - (Optional) The ID of the Storage Account. Changing this forces a new resource to be created.

---

A `storage_account_gen2` block supports the following:

* `is_default` - (Required) Is this the Default Storage Account for the HDInsight Hadoop Cluster? Changing this forces a new resource to be created.

-> **Note:** One of the `storage_account` or `storage_account_gen2` blocks must be marked as the default.

* `storage_resource_id` - (Required) The ID of the Storage Account. Changing this forces a new resource to be created.

* `filesystem_id` - (Required) The ID of the Gen2 Filesystem. Changing this forces a new resource to be created.

* `managed_identity_resource_id` - (Required) The ID of Managed Identity to use for accessing the Gen2 filesystem. Changing this forces a new resource to be created.

-> **Note:** This can be obtained from the `id` of the `azurerm_storage_container` resource.

---

A `private_link_configuration` block supports the following:

* `name` - (Required) The name of the private link configuration.

* `group_id` - (Required) The ID of the private link service group.

* `private_link_service_connection` - (Required) A `private_link_service_connection` block as defined below.

---

A `private_link_service_connection` block supports the following:

* `name` - (Required) The name of the private link service connection.

* `primary` - (Optional) Indicates whether this IP configuration is primary.

* `private_ip_allocation_method` - (Optional) The private IP allocation method. The only possible value now is `Dynamic`.

* `private_ip_address` - (Optional) The private IP address of the IP configuration.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the private link service connection should be provisioned within.

---

A `worker_node` block supports the following:

* `script_actions` - (Optional) The script action which will run on the cluster. One or more `script_actions` blocks as defined above.

* `username` - (Required) The Username of the local administrator for the Worker Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Worker Nodes. Possible values are `ExtraSmall`, `Small`, `Medium`, `Large`, `ExtraLarge`, `A5`, `A6`, `A7`, `A8`, `A9`, `A10`, `A11`, `Standard_A1_V2`, `Standard_A2_V2`, `Standard_A2m_V2`, `Standard_A3`, `Standard_A4_V2`, `Standard_A4m_V2`, `Standard_A8_V2`, `Standard_A8m_V2`, `Standard_D1`, `Standard_D2`, `Standard_D3`, `Standard_D4`, `Standard_D11`, `Standard_D12`, `Standard_D13`, `Standard_D14`, `Standard_D1_V2`, `Standard_D2_V2`, `Standard_D3_V2`, `Standard_D4_V2`, `Standard_D5_V2`, `Standard_D11_V2`, `Standard_D12_V2`, `Standard_D13_V2`, `Standard_D14_V2`, `Standard_DS1_V2`, `Standard_DS2_V2`, `Standard_DS3_V2`, `Standard_DS4_V2`, `Standard_DS5_V2`, `Standard_DS11_V2`, `Standard_DS12_V2`, `Standard_DS13_V2`, `Standard_DS14_V2`, `Standard_E2_V3`, `Standard_E4_V3`, `Standard_E8_V3`, `Standard_E16_V3`, `Standard_E20_V3`, `Standard_E32_V3`, `Standard_E64_V3`, `Standard_E64i_V3`, `Standard_E2s_V3`, `Standard_E4s_V3`, `Standard_E8s_V3`, `Standard_E16s_V3`, `Standard_E20s_V3`, `Standard_E32s_V3`, `Standard_E64s_V3`, `Standard_E64is_V3`, `Standard_D2a_V4`, `Standard_D4a_V4`, `Standard_D8a_V4`, `Standard_D16a_V4`, `Standard_D32a_V4`, `Standard_D48a_V4`, `Standard_D64a_V4`, `Standard_D96a_V4`, `Standard_E2a_V4`, `Standard_E4a_V4`, `Standard_E8a_V4`, `Standard_E16a_V4`, `Standard_E20a_V4`, `Standard_E32a_V4`, `Standard_E48a_V4`, `Standard_E64a_V4`, `Standard_E96a_V4`, `Standard_D2ads_V5`, `Standard_D4ads_V5`, `Standard_D8ads_V5`, `Standard_D16ads_V5`, `Standard_D32ads_V5`, `Standard_D48ads_V5`, `Standard_D64ads_V5`, `Standard_D96ads_V5`, `Standard_E2ads_V5`, `Standard_E4ads_V5`, `Standard_E8ads_V5`, `Standard_E16ads_V5`, `Standard_E20ads_V5`, `Standard_E32ads_V5`, `Standard_E48ads_V5`, `Standard_E64ads_V5`, `Standard_E96ads_V5`, `Standard_G1`, `Standard_G2`, `Standard_G3`, `Standard_G4`, `Standard_G5`, `Standard_F2s_V2`, `Standard_F4s_V2`, `Standard_F8s_V2`, `Standard_F16s_V2`, `Standard_F32s_V2`, `Standard_F64s_V2`, `Standard_F72s_V2`, `Standard_GS1`, `Standard_GS2`, `Standard_GS3`, `Standard_GS4`, `Standard_GS5` and `Standard_NC24`. Changing this forces a new resource to be created.

-> **Note:** High memory instances must be specified for the Head Node (Azure suggests a `Standard_D14_V2`).

* `password` - (Optional) The Password associated with the local administrator for the Worker Nodes. Changing this forces a new resource to be created.

-> **Note:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Worker Nodes. Changing this forces a new resource to be created.

-> **Note:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Worker Nodes should be provisioned within. Changing this forces a new resource to be created.

* `target_instance_count` - (Required) The number of instances which should be run for the Worker Nodes.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Worker Nodes should be provisioned within. Changing this forces a new resource to be created.

* `autoscale` - (Optional) A `autoscale` block as defined below.

---

A `disk_encryption` block supports the following:

* `encryption_algorithm` - (Optional) This is an algorithm identifier for encryption. Possible values are `RSA1_5`, `RSA-OAEP`, `RSA-OAEP-256`.

* `encryption_at_host_enabled` - (Optional) This is indicator to show whether resource disk encryption is enabled.

* `key_vault_key_id` - (Optional) The ID of the key vault key.

* `key_vault_managed_identity_id` - (Optional) This is the resource ID of Managed Identity used to access the key vault.

---

A `zookeeper_node` block supports the following:

* `script_actions` - (Optional) The script action which will run on the cluster. One or more `script_actions` blocks as defined above.

* `username` - (Required) The Username of the local administrator for the Zookeeper Nodes. Changing this forces a new resource to be created.

* `vm_size` - (Required) The Size of the Virtual Machine which should be used as the Zookeeper Nodes. Possible values are `ExtraSmall`, `Small`, `Medium`, `Large`, `ExtraLarge`, `A5`, `A6`, `A7`, `A8`, `A9`, `A10`, `A11`, `Standard_A1_V2`, `Standard_A2_V2`, `Standard_A2m_V2`, `Standard_A3`, `Standard_A4_V2`, `Standard_A4m_V2`, `Standard_A8_V2`, `Standard_A8m_V2`, `Standard_D1`, `Standard_D2`, `Standard_D3`, `Standard_D4`, `Standard_D11`, `Standard_D12`, `Standard_D13`, `Standard_D14`, `Standard_D1_V2`, `Standard_D2_V2`, `Standard_D3_V2`, `Standard_D4_V2`, `Standard_D5_V2`, `Standard_D11_V2`, `Standard_D12_V2`, `Standard_D13_V2`, `Standard_D14_V2`, `Standard_DS1_V2`, `Standard_DS2_V2`, `Standard_DS3_V2`, `Standard_DS4_V2`, `Standard_DS5_V2`, `Standard_DS11_V2`, `Standard_DS12_V2`, `Standard_DS13_V2`, `Standard_DS14_V2`, `Standard_E2_V3`, `Standard_E4_V3`, `Standard_E8_V3`, `Standard_E16_V3`, `Standard_E20_V3`, `Standard_E32_V3`, `Standard_E64_V3`, `Standard_E64i_V3`, `Standard_E2s_V3`, `Standard_E4s_V3`, `Standard_E8s_V3`, `Standard_E16s_V3`, `Standard_E20s_V3`, `Standard_E32s_V3`, `Standard_E64s_V3`, `Standard_E64is_V3`, `Standard_D2a_V4`, `Standard_D4a_V4`, `Standard_D8a_V4`, `Standard_D16a_V4`, `Standard_D32a_V4`, `Standard_D48a_V4`, `Standard_D64a_V4`, `Standard_D96a_V4`, `Standard_E2a_V4`, `Standard_E4a_V4`, `Standard_E8a_V4`, `Standard_E16a_V4`, `Standard_E20a_V4`, `Standard_E32a_V4`, `Standard_E48a_V4`, `Standard_E64a_V4`, `Standard_E96a_V4`, `Standard_D2ads_V5`, `Standard_D4ads_V5`, `Standard_D8ads_V5`, `Standard_D16ads_V5`, `Standard_D32ads_V5`, `Standard_D48ads_V5`, `Standard_D64ads_V5`, `Standard_D96ads_V5`, `Standard_E2ads_V5`, `Standard_E4ads_V5`, `Standard_E8ads_V5`, `Standard_E16ads_V5`, `Standard_E20ads_V5`, `Standard_E32ads_V5`, `Standard_E48ads_V5`, `Standard_E64ads_V5`, `Standard_E96ads_V5`, `Standard_G1`, `Standard_G2`, `Standard_G3`, `Standard_G4`, `Standard_G5`, `Standard_F2s_V2`, `Standard_F4s_V2`, `Standard_F8s_V2`, `Standard_F16s_V2`, `Standard_F32s_V2`, `Standard_F64s_V2`, `Standard_F72s_V2`, `Standard_GS1`, `Standard_GS2`, `Standard_GS3`, `Standard_GS4`, `Standard_GS5` and `Standard_NC24`. Changing this forces a new resource to be created.

* `password` - (Optional) The Password associated with the local administrator for the Zookeeper Nodes. Changing this forces a new resource to be created.

-> **Note:** If specified, this password must be at least 10 characters in length and must contain at least one digit, one uppercase and one lower case letter, one non-alphanumeric character (except characters ' " ` \).

* `ssh_keys` - (Optional) A list of SSH Keys which should be used for the local administrator on the Zookeeper Nodes. Changing this forces a new resource to be created.

-> **Note:** Either a `password` or one or more `ssh_keys` must be specified - but not both.

* `subnet_id` - (Optional) The ID of the Subnet within the Virtual Network where the Zookeeper Nodes should be provisioned within. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of the Virtual Network where the Zookeeper Nodes should be provisioned within. Changing this forces a new resource to be created.

---

A `metastores` block supports the following:

* `hive` - (Optional) A `hive` block as defined below.

* `oozie` - (Optional) An `oozie` block as defined below.

* `ambari` - (Optional) An `ambari` block as defined below.

---

A `hive` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Hive metastore. Changing this forces a new resource to be created.

* `database_name` - (Required) The external Hive metastore's existing SQL database. Changing this forces a new resource to be created.

* `username` - (Required) The external Hive metastore's existing SQL server admin username. Changing this forces a new resource to be created.

* `password` - (Required) The external Hive metastore's existing SQL server admin password. Changing this forces a new resource to be created.

---

An `oozie` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Oozie metastore. Changing this forces a new resource to be created.

* `database_name` - (Required) The external Oozie metastore's existing SQL database. Changing this forces a new resource to be created.

* `username` - (Required) The external Oozie metastore's existing SQL server admin username. Changing this forces a new resource to be created.

* `password` - (Required) The external Oozie metastore's existing SQL server admin password. Changing this forces a new resource to be created.

---

An `ambari` block supports the following:

* `server` - (Required) The fully-qualified domain name (FQDN) of the SQL server to use for the external Ambari metastore. Changing this forces a new resource to be created.

* `database_name` - (Required) The external Hive metastore's existing SQL database. Changing this forces a new resource to be created.

* `username` - (Required) The external Ambari metastore's existing SQL server admin username. Changing this forces a new resource to be created.

* `password` - (Required) The external Ambari metastore's existing SQL server admin password. Changing this forces a new resource to be created.

---

A `monitor` block supports the following:

* `log_analytics_workspace_id` - (Required) The Operations Management Suite (OMS) workspace ID.

* `primary_key` - (Required) The Operations Management Suite (OMS) workspace key.

---

A `extension` block supports the following:

* `log_analytics_workspace_id` - (Required) The workspace ID of the log analytics extension.

* `primary_key` - (Required) The workspace key of the log analytics extension.

---

An `autoscale` block supports the following:

* `recurrence` - (Optional) A `recurrence` block as defined below.

---

A `recurrence` block supports the following:

* `schedule` - (Required) A list of `schedule` blocks as defined below.

* `timezone` - (Required) The time zone for the autoscale schedule times.

---

A `schedule` block supports the following:

* `days` - (Required) The days of the week to perform autoscale. Possible values are `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` and `Sunday`.

* `target_instance_count` - (Required) The number of worker nodes to autoscale at the specified time.

* `time` - (Required) The time of day to perform the autoscale in 24hour format.

---

A `security_profile` block supports the following:

* `aadds_resource_id` - (Required) The resource ID of the Azure Active Directory Domain Service. Changing this forces a new resource to be created.

* `domain_name` - (Required) The name of the Azure Active Directory Domain. Changing this forces a new resource to be created.

* `domain_username` - (Required) The username of the Azure Active Directory Domain. Changing this forces a new resource to be created.

* `domain_user_password` - (Required) The user password of the Azure Active Directory Domain. Changing this forces a new resource to be created.

* `ldaps_urls` - (Required) A list of the LDAPS URLs to communicate with the Azure Active Directory. Changing this forces a new resource to be created.

* `msi_resource_id` - (Required) The User Assigned Identity for the HDInsight Cluster. Changing this forces a new resource to be created.

* `cluster_users_group_dns` - (Optional) A list of the distinguished names for the cluster user groups. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the HDInsight Interactive Query Cluster.

* `https_endpoint` - The HTTPS Connectivity Endpoint for this HDInsight Interactive Query Cluster.

* `ssh_endpoint` - The SSH Connectivity Endpoint for this HDInsight Interactive Query Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Interactive Query HDInsight Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Interactive Query HDInsight Cluster.
* `update` - (Defaults to 1 hour) Used when updating the Interactive Query HDInsight Cluster.
* `delete` - (Defaults to 1 hour) Used when deleting the Interactive Query HDInsight Cluster.

## Import

HDInsight Interactive Query Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hdinsight_interactive_query_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.HDInsight/clusters/cluster1
```
