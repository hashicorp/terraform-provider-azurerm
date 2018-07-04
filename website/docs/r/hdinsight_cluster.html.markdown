---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_cluster"
sidebar_current: "docs-azurerm-resource-data-hdinsight-cluster"
description: |-
  Manages a HDInsight Cluster
---

# azurerm_hdinsight_cluster

Manages a HDInsight Cluster

~> **NOTE:** The HDInsights API isn't particularly descriptive when an error occurs. If you see the error `User input validation failed. Errors: The request payload is invalid.` - we'd suggest checking the machine configurations (e.g. sizes/counts) are valid. There's [an issue requesting better error handling for the HDInsights API](https://github.com/Azure/azure-sdk-for-go/issues/2179).

## Example Usage

```hcl
locals {
  username = "hdiadmin"
  password = "examplePass419"
}

resource "azurerm_resource_group" "test" {
  name     = "hdinsight-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "tfhdistor2018"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "data"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_hdinsight_cluster" "test" {
  name                = "terraform-hdi"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tier                = "standard"

  cluster {
    kind    = "Hadoop"
    version = "3.6"

    gateway {
      username = "hdsuperadmin"
      password = "Messina1234!"
    }
  }

  storage_profile {
    storage_account {
      storage_account_name = "${azurerm_storage_account.test.primary_blob_domain}"
      storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
      container_name       = "${azurerm_storage_container.test.name}"
      is_default           = true
    }
  }

  head_node {
    target_instance_count = 2

    hardware_profile {
      vm_size = "Standard_D3_V2"
    }

    os_profile {
      username = "${local.username}"
      password = "${local.password}"
    }
  }

  worker_node {
    target_instance_count = 4

    hardware_profile {
      vm_size = "Medium"
    }

    os_profile {
      username = "${local.username}"
      password = "${local.password}"
    }
  }

  zookeeper_node {
    target_instance_count = 3

    hardware_profile {
      vm_size = "A5"
    }

    os_profile {
      username = "${local.username}"
      password = "${local.password}"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the HDInsight Cluster. Changing this forces a new resource to be created.

-> **NOTE:** The name must be 59 characters or less and can contain letters, numbers, and hyphens (but the first and last character must be a letter or number) - and cannot contain a reserved word.

* `resource_group_name` - (Required) The name of the resource group in which to create the HDInsight Cluster. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the HDInsight Cluster exists. Changing this forces a new resource to be created.

* `tier` - (Required) The Pricing Tier to use for this HDInsight Cluster. Possible values include `Standard` and `Premium`. Changing this forces a new resource to be created.

* `cluster` - (Required) A `cluster` block as defined below.
* `storage_profile` - (Required) A `storage_profile` block as defined below.

* `head_node` - (Required) A `head_node` block as defined below.
* `worker_node` - (Required) A `worker_node` block as defined below.
* `zookeeper_node` - (Required) A `zookeeper_node` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `cluster` block contains:

* `kind` - (Required) The Kind of HDInsight Cluster which should be created. Possible values include `Hadoop`, `HBase`, `InteractiveHive`, `Kafka`, `RServer`, `Storm` and `Spark`. Changing this forces a new resource to be created.

* `version` - (Required) The version of the HDInsight Cluster which should be created, such as `3.6`. Changing this forces a new resource to be created.

* `gateway` - (Required) A `gateway` block as defined below. Changing this forces a new resource to be created.

---

A `gateway` block contains:

* `username` - (Required) The username associated with the Web Gateway. Changing this forces a new resource to be created.

* `password` - (Required) The password associated with the username used for the Web Gateway. Changing this forces a new resource to be created.

---

A `storage_profile` block contains:

* `storage_account` - (Required) One or mote `storage_account` blocks as defined below.

~> **NOTE:** Support for Data Lake Stores will be added in the future when an API is available

---

A `storage_account` block contains:

* `storage_account_name` - (Required) The Host Name associated with the Blob Storage Account. Changing this forces a new resource to be created.

-> This value can be found as the `blob_storage_domain` property on the `azurerm_storage_account` Data Source / Resource.

* `storage_account_key` - (Required) An Access Key associated with the Storage Account, used to push data. Changing this forces a new resource to be created.

* `container_name` - (Required) The name of the Container within the Storage Account where the data . Changing this forces a new resource to be created.

* `is_default` - (Required) Is this the default Storage Account?  Changing this forces a new resource to be created.

---

A `head_node` block as defined below:

* `target_instance_count` - (Required) The target instance count for the Head Node role. This can be a minimum of 2 nodes. Changing this forces a new resource to be created.
* `hardware_profile` - (Required) A `hardware_profile` block as defined below.
* `os_profile` - (Required) A `os_profile` block as defined below.
* `virtual_network_profile` - (Optional) A `virtual_network_profile` block as defined below.

---

A `worker_node` block as defined below:

* `target_instance_count` - (Required) The target instance count for the Worker Node role. This can be a minimum of 2 nodes.
* `hardware_profile` - (Required) A `hardware_profile` block as defined below.
* `os_profile` - (Required) A `os_profile` block as defined below.
* `virtual_network_profile` - (Optional) A `virtual_network_profile` block as defined below.

---

A `zookeeper_node` block as defined below:

* `target_instance_count` - (Required) The target instance count for the Zookeeper Node role. This can be a minimum of 3 nodes. Changing this forces a new resource to be created.
* `hardware_profile` - (Required) A `hardware_profile` block as defined below.
* `os_profile` - (Required) A `os_profile` block as defined below.
* `virtual_network_profile` - (Optional) A `virtual_network_profile` block as defined below.

---

A `hardware_profile` block contains:

* `vm_size` - (Required) The size of the Virtual Machine, such as `Standard_D3_v2`. Changing this forces a new resource to be created.

---

A `os_profile` block contains:

* `username` - (Required) The Username associated with the local administrator account on the Virtual Machine. Changing this forces a new resource to be created.
* `password` - (Required) The Password to use for the local administrator account on the Virtual Machine. Changing this forces a new resource to be created.

---

A `virtual_network_profile` block contains:

* `virtual_network_id` - (Required) The ID of the Virtual Network in which the HDInsight Cluster should be created. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet within the Virtual Network in which the HDInsight Cluster should be created. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HDInsight Cluster.

* `https_endpoint` - The HTTPS Endpoint for the Ambari Dashboard.

* `ssh_endpoint` - The SSH Endpoint for the Ambari Cluster.

## Import

Image can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hdinsight_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.HDInsight/clusters/cluster1
```
