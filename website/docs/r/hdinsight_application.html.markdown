---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_application"
sidebar_current: "docs-azurerm-resource-data-hdinsight-application"
description: |-
  Manages an Application within a HDInsight Cluster
---

# azurerm_hdinsight_application

Manages an Application within a HDInsight Cluster

~> **NOTE:** The HDInsights API isn't particularly descriptive when an error occurs. If you see the error `User input validation failed. Errors: The request payload is invalid.` - we'd suggest checking the machine configurations (e.g. sizes/counts) are valid. There's [an issue requesting better error handling for the HDInsights API](https://github.com/Azure/azure-sdk-for-go/issues/2179).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  # ...
}

resource "azurerm_hdinsight_cluster" "test" {
  # ...
}

resource "azurerm_hdinsight_application" "test" {
  name                   = "emptynodeapp"
  cluster_name           = "${azurerm_hdinsight_cluster.test.name}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  marketplace_identifier = "EmptyNode"

  edge_node {
    target_instance_count = 1

    hardware_profile {
      vm_size = "Standard_D3_v2"
    }
  }

  install_script_action {
    name  = "emptynode-sayhello"
    uri   = "https://gist.githubusercontent.com/tombuildsstuff/74ff75620a83cf2a737843920185dbc2/raw/8217fbbcf9728e23807c19a35f65136351e6da7a/hello.sh"
    roles = [ "edgenode" ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the HDInsight Application, which must be unique within the Cluster. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the HDInsight Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the HDInsight Cluster exists. Changing this forces a new resource to be created.

* `marketplace_identifier` - (Required) The Marketplace Identifier for this Application. Changing this forces a new resource to be created.

* `edge_node` - (Optional) A `edge_node` block as defined below.

* `install_script_action` - (Optional) One or more `install_script_action` blocks as defined below. Changing this forces a new resource to be created.

* `uninstall_script_action` - (Optional) One or more `uninstall_script_action` blocks as defined below. Changing this forces a new resource to be created.

---

A `edge_node` block supports the following arguments:

* `hardware_profile` - (Required) A `hardware_profile` block as defined below. Changing this forces a new resource to be created.

---

A `hardware_profile` block supports the following arguments:

* `vm_size` - (Required) The size of the Virtual Machine, such as `Standard_D3_v2`. Changing this forces a new resource to be created.

---

A `install_script_action` block supports the following arguments:

* `name` - (Required) The name of the install script action, which must be unique across script actions on the cluster. Changing this forces a new resource to be created.

* `uri` - (Required) The path to a publicly-accessible idempotent script which should be run. Changing this forces a new resource to be created.

* `roles` - (Required) One or more roles which this script should be run on. Possible values include `edgenode`, `headnode`, `workernode` and `zookeepernode`. Changing this forces a new resource to be created.

---

A `uninstall_script_action` block supports the following arguments:

* `name` - (Required) The name of the install script action, which must be unique across script actions on the cluster. Changing this forces a new resource to be created.

* `uri` - (Required) The path to a publicly-accessible idempotent script which should be run. Changing this forces a new resource to be created.

* `roles` - (Required) One or more roles which this script should be run on. Possible values include `edgenode`, `headnode`, `workernode` and `zookeepernode`. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HDInsight Application.

## Import

HDInsight Applications can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hdinsight_application.app1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.HDInsight/clusters/cluster1/applications/app1
```
