---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_cluster"
sidebar_current: "docs-azurerm-resource-hdinsight-cluster"
description: |-
  Create an HDInsight cluster component.
---

# azurerm_hdinsight_cluster

Create an HDInsight cluster component.

## Example Usage with two Head and four Worker nodes

```hcl
resource "azurerm_resource_group" "main" {
  name     = "resourceGroupName"
  location = "West US 2"
}

resource "azurerm_storage_account" "main" {
  name                     = "storageaccountname"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  location                 = "${azurerm_resource_group.main.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_storage_container" "main" {
  name                  = "containername"
  resource_group_name   = "${azurerm_resource_group.main.name}"
  storage_account_name  = "${azurerm_storage_account.main.name}"
  container_access_type = "private"
}

resource "azurerm_hdinsight_cluster" "main" {
  name                = "hdinsightclustername"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  cluster_type        = "hadoop"
  cluster_version     = "3.6"
  login_username   = "testadmin"
  login_password   = "Password1234!"

  storage_account {
    blob_endpoint = "${azurerm_storage_account.main.primary_blob_endpoint}"
    container     = "${azurerm_storage_container.main.name}"
    access_key    = "${azurerm_storage_account.main.primary_access_key}"
  }

  head_node {
    target_instance_count = 2
    vm_size               = "Large"

    linux_os_profile {
      username = "super"

      ssh_keys {
        key_data = "ssh-rsa <Public SSH Key Data>"
      }
    }
  }

  worker_node {
    target_instance_count = 4
    vm_size               = "Large"

    linux_os_profile {
      username = "super"

      ssh_keys {
        key_data = "ssh-rsa <Public SSH Key Data>"
      }
    }
  }
}

output "endpoints" {
  value = "${azurerm_hdinsight_cluster.main.connectivity_endpoints}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights component. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Application Insights component.
* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.
* `cluster_type` - (Required) Specifies the type of this HDInsight cluster. Possible values are: `hadoop`, `hbase`, `storm` and `spark`. Changing this forces a new resource to be created.
* `cluster_version` - (Required) Specifies the version of this HDInsight cluster. It should be in the form of `<major>.<minor>`. Changing this forces a new resource to be created.
* `component_version` - (Required) Changing this forces a new resource to be created.
* `login_username` - (Required) Specifies the username of this HDInsight cluster. It can be used to log in to the cluster dashboard and submit jobs. Changing this forces a new resource to be created.
* `login_password` - (Required) Specifies the password of this HDInsight cluster. It can be used to log in to the cluster dashboard and submit jobs. Changing this forces a new resource to be created.
* `storage_account` - (Required) A Storage Account Reference block as documented below.
* `head_node` - (Required) A head compute node block as documented below.
* `worker_node` - (Required) A worker compute node block as documented below.
* `zookeeper_node` - (Optional) A zookeeper compute node block as documented below.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `tier` - (Optional) Specifies the tier for the cluster. Possible values are: `Standard` and `Premium`. Changing this forces a new resource to be created.

`storage_account` supports the following:

* `blob_endpoint`: (Required) The blob endpoint of an existing storage account. Changing this forces a new resource to be created.
* `access_key`: (Required) Specifies the access key of your storage account. Changing this forces a new resource to be created.
* `container`: (Required) Specifies a storage container to be used by this HDInsight cluster. Changing this forces a new resource to be created.

`head_node`, `worker_node` and `zookeeper_node` supports the following:

* `linux_os_profile`: (Required) A Linux Profile block as documented below.
* `target_instance_count`: (Required) The target instances count used by this HDInsight cluster. Changing this forces a new resource to be created.
* `min_instance_count`: (Optional) The minimum instances count used by this HDInsight cluster. Changing this forces a new resource to be created.

`linux_os_profile` supports the following:

* `username`: (Required) The username used to login to the Linux machine.
* `ssh_keys`: (Required) The public SSH key block used to login to the Linux machine as documented below.
* `password`: (Optional) The password used to login to the Linux machine.

`ssh_keys` supports the following:

* `key_data`: (Required) One or more SSH public key data used in this Linux machine.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Insights component.
* `connectivity_endpoints` - The Connectivity Endpoints block as documented below.

`connectivity_endpoints` supports the following:

* `name`: The type (such as `SSH` or `HTTPS`) of this endpoint.
* `location`: The host URL of this endpoint.
* `port`: The host port number of this endpoint.
* `protocol`: The protocol (such as `TCP`) used by this endpoint.


## Import

Application Insights instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hdinsight_cluster /subscriptions/<subscription ID>/resourceGroups/<resource group name>/providers/Microsoft.HDInsight/clusters/<HDInsight cluster name>
```
