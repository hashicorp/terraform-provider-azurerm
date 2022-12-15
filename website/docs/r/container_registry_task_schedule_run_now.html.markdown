---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_task_schedule_run_now"
description: |-
  Runs a Container Registry Task Schedule.
---

# azurerm_container_registry_task_schedule_run_now

Runs a Container Registry Task Schedule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}
resource "azurerm_container_registry" "example" {
  name                = "example-acr"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Basic"
}
resource "azurerm_container_registry_task" "example" {
  name                  = "example-task"
  container_registry_id = azurerm_container_registry.example.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "https://github.com/<user name>/acr-build-helloworld-node#main"
    context_access_token = "<github personal access token>"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
}
resource "azurerm_container_registry_task_schedule_run_now" "example" {
  container_registry_task_id = azurerm_container_registry_task.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `container_registry_task_id` - (Required) The ID of the Container Registry Task that to be scheduled. Changing this forces a new Container Registry Task Schedule to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry Task Schedule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry Task Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Task Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Task Schedule.
