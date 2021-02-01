---
subcategory: "Service Fabric Mesh"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_fabric_mesh_application"
description: |-
  Manages a Service Fabric Mesh Application.
---

# azurerm_service_fabric_mesh_application

Manages a Service Fabric Mesh Application.

!> **Note:** Service Fabric Mesh is being retired on `2021-04-28` and **new Clusters can no longer be provisioned**. More information about [the retirement can be found here](https://azure.microsoft.com/en-us/updates/azure-service-fabric-mesh-preview-retirement/). Azure recommends migrating to either a [Azure Container Instances](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/container_group), [Azure Kubernetes Service](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/kubernetes_cluster) or [a Service Fabric Managed Cluster](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/service_fabric_cluster).

!> **Note:** **This resource is deprecated** and will be removed in version 3.0 of the Azure Provider.

## Example Usage


```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_fabric_mesh_application" "example" {
  name                = "example-mesh-application"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  service {
    name    = "testservice1"
    os_type = "Linux"

    code_package {
      name       = "testcodepackage1"
      image_name = "seabreeze/sbz-helloworld:1.0-alpine"

      resources {
        requests {
          memory = 1
          cpu    = 1
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service Fabric Mesh Application. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Service Fabric Mesh Application exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Service Fabric Mesh Application should exist. Changing this forces a new resource to be created.

* `service` - (Required) Any number of `service` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `service` block supports the following:

* `name` - (Required) The name of the service resource.

* `os_type` - (Required) The operating system required by the code in service. Valid values are `Linux` or `Windows`.

* `code_package` - (Required) Any number `code_package` block as described below.

---

A `code_package` block supports the following:

* `name` - (Required) The name of the code package.

* `image_name` - (Required) The Container image the code package will use.

* `resources` - (Required) A `resources` block as defined below.

---

A `resources` block supports the following: 

* `requests` - (Required) A `requests` block as defined below.

* `limits` - (Optional) A `limits` block as defined below.

---

A `requests` block supports the following: 

* `cpu` - (Required) The minimum number of CPU cores the container requires. 

* `memory` - (Required) The minimum memory request in GB the container requires.

---

A `limits` block supports the following: 

* `cpu` - (Required) The maximum number of CPU cores the container can use. 

* `memory` - (Required) The maximum memory request in GB the container can use.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Service Fabric Mesh Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Fabric Mesh Application.
* `update` - (Defaults to 30 minutes) Used when updating the Service Fabric Mesh Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Fabric Mesh Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Fabric Mesh Application.

## Import

Service Fabric Mesh Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_service_fabric_mesh_application.application1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceFabricMesh/applications/application1
```
