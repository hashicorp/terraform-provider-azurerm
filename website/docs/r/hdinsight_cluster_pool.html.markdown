---
subcategory: "HDInsight"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hdinsight_cluster_pool"
description: |-
  Manages a HDInsight HBase Cluster Pool.
---

# azurerm_hdinsight_cluster_pool

Manages a HDInsight HBase Cluster Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_hdinsight_cluster_pool" "example" {
  name                        = "example-pool"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  managed_resource_group_name = "test-rg"

  cluster_pool_profile {
    cluster_pool_version = "1.1"
  }
  compute_profile {
    vm_size = "Standard_F4s_v2"
  }
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) Specifies the name of the HDInsight Cluster Pool resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the HDInsight Cluster Pool. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Required) A resource group created by RP, to hold the resources created by RP on-behalf of customers. Changing this forces a new resource to be created.

* `cluster_pool_profile` - (Required) The profile of the HDInsight Cluster Pool. A `cluster_pool_profile` block supports the following:

  * `cluster_pool_version` - (Required) The version of the HDInsight Cluster Pool. Changing this forces a new resource to be created.

* `compute_profile` - (Required) The compute profile of the HDInsight Cluster Pool. Chaning this forces a new resource to be created. A `compute_profile` block supports the following:
    
  * `vm_size` - (Required) The size of the Virtual Machine. Changing this forces a new resource to be created.

* `log_analytics_profile` - (Optional) A `log_analytics_profile` block as defined below.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `log_analytics_profile` block supports the following:

* `log_analytics_profile_enabled` - (Required) Specifies whether the Log Analytics is enabled.

* `workspace_id` - (Required) The ID of the Log Analytics Workspace.

---

A `network_profile` block supports the following:

* `subnet_id` - (Required) The cluster pool subnet ID.

* `private_api_server_enabled` - (Optional) Specifies whether the private API server is enabled.

* `outbound_type` - (Optional) The outbound type of the HDInsight Cluster Pool. Possible values are `loadBalancer` and `userDefinedRouting`.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - The ID of the HDInsight Cluster Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the HDInsight Cluster Pool.
* `update` - (Defaults to 60 minutes) Used when updating the HDInsight Cluster Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the HDInsight Cluster Pool.
* `delete` - (Defaults to 60 minutes) Used when deleting the HDInsight Cluster Pool.

## Import

HDInsight Cluster Pool can be imported using the resource id, e.g.

```shell
terraform import azurerm_hdinsight_cluster_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.HDInsight/clusterPools/myclusterpool1
```
