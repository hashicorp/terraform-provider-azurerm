---
subcategory: "Service Fabric Mesh"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_fabric_mesh_local_network"
description: |-
  Manages a Service Fabric Mesh Local Network.
---

# azurerm_service_fabric_mesh_local_network

Manages a Service Fabric Mesh Local Network.

## Example Usage


```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_fabric_mesh_local_network" "example" {
  name                = "example-mesh-local-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  network_address_prefix = "10.0.0.0/22"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Fabric Mesh Local Network. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Service Fabric Mesh Local Network exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Service Fabric Mesh Local Network should exist. Changing this forces a new resource to be created.

* `network_address_prefix` - (Required) The address space for the local container network.

* `description` - (Optional) A description of this Service Fabric Mesh Local Network.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Service Fabric Mesh Local Network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Fabric Mesh Local Network.
* `update` - (Defaults to 30 minutes) Used when updating the Service Fabric Mesh Local Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Fabric Mesh Local Network.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Fabric Mesh Local Network.

## Import

Service Fabric Mesh Local Network can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_fabric_mesh_local_network.network1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceFabricMesh/networks/network1
```
