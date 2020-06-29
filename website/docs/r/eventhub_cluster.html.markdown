---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_cluster"
description: |-
  Manages an EventHub Cluster

---

# azurerm_eventhub_cluster

Manages an EventHub Cluster

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US 2"
}

resource "azurerm_eventhub_cluster" "test" {
  name                = "acctesteventhubcluster-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Dedicated_1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Cluster resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventHub Cluster exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) The sku name of the EventHub Cluster. The only supported value at this time is `Dedicated_1`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub Cluster ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Cluster.
* `delete` - (Defaults to 300 minutes) Used when deleting the EventHub Cluster.

## Import

EventHub Cluster's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_cluster.cluster1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1
```
