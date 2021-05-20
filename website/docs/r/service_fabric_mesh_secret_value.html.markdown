---
subcategory: "Service Fabric Mesh"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_fabric_mesh_secret_value"
description: |-
  Manages a Service Fabric Mesh Secret Value.
---

# azurerm_service_fabric_mesh_secret_value

Manages a Service Fabric Mesh Secret Value.

!> **Note:** Service Fabric Mesh is being retired on `2021-04-28` and **new Clusters can no longer be provisioned**. More information about [the retirement can be found here](https://azure.microsoft.com/en-us/updates/azure-service-fabric-mesh-preview-retirement/). Azure recommends migrating to either a [Azure Container Instances](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/container_group), [Azure Kubernetes Service](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/kubernetes_cluster) or [a Service Fabric Managed Cluster](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/service_fabric_cluster).

!> **Note:** **This resource is deprecated** and will be removed in version 3.0 of the Azure Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_fabric_mesh_secret" "example" {
  name                = "example-mesh-secret"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_service_fabric_mesh_secret_value" "example" {
  name                          = "example-secret-value"
  service_fabric_mesh_secret_id = azurerm_service_fabric_mesh_secret_inline.test.id
  location                      = azurerm_resource_group.test.location
  value                         = "testValue"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Fabric Mesh Secret Value. Changing this forces a new resource to be created.

* `service_fabric_mesh_secret_id` - (Required) The id of the Service Fabric Mesh Secret in which the value will be applied to. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Service Fabric Mesh Secret Value should exist. Changing this forces a new resource to be created.

* `value` - (Required) Specifies the value that will be applied to the Service Fabric Mesh Secret. Changing this forces a new resource to be created.

* `description` - (Optional) A description of this Service Fabric Mesh Secret Value.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Service Fabric Mesh Secret Value.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Fabric Mesh Secret Value.
* `update` - (Defaults to 30 minutes) Used when updating the Service Fabric Mesh Secret Value.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Fabric Mesh Secret Value.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Fabric Mesh Secret Value.

## Import

Service Fabric Mesh Secret Value can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_fabric_mesh_secret_value.value1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceFabricMesh/secrets/secret1/values/value1
```
