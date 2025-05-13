---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_container_deployment"
description: |-
  Manages a Spring Cloud Container Deployment.
---

# azurerm_spring_cloud_container_deployment

Manages a Spring Cloud Container Deployment.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_container_deployment` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "example"
  resource_group_name = azurerm_spring_cloud_service.example.resource_group_name
  service_name        = azurerm_spring_cloud_service.example.name
}

resource "azurerm_spring_cloud_container_deployment" "example" {
  name                = "example"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
  instance_count      = 2
  arguments           = ["-cp", "/app/resources:/app/classes:/app/libs/*", "hello.Application"]
  commands            = ["java"]
  environment_variables = {
    "Foo" : "Bar"
    "Env" : "Staging"
  }
  server             = "docker.io"
  image              = "springio/gs-spring-boot-docker"
  language_framework = "springboot"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Container Deployment. Changing this forces a new Spring Cloud Container Deployment to be created.

* `spring_cloud_app_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new Spring Cloud Container Deployment to be created.

* `image` - (Required) Container image of the custom container. This should be in the form of `<repository>:<tag>` without the server name of the registry.

* `server` - (Required) The name of the registry that contains the container image.

---

* `addon_json` - (Optional) A JSON object that contains the addon configurations of the Spring Cloud Container Deployment.

* `application_performance_monitoring_ids` - (Optional) Specifies a list of Spring Cloud Application Performance Monitoring IDs.

* `arguments` - (Optional) Specifies the arguments to the entrypoint. The docker image's `CMD` is used if not specified.

* `commands` - (Optional) Specifies the entrypoint array. It will not be executed within a shell. The docker image's `ENTRYPOINT` is used if not specified.

* `environment_variables` - (Optional) Specifies the environment variables of the Spring Cloud Deployment as a map of key-value pairs.

* `instance_count` - (Optional) Specifies the required instance count of the Spring Cloud Deployment. Possible Values are between `1` and `500`. Defaults to `1` if not specified.

* `language_framework` - (Optional) Specifies the language framework of the container image. The only possible value is `springboot`.

* `quota` - (Optional) A `quota` block as defined below.

---

A `quota` block supports the following:

* `cpu` - (Optional) Specifies the required cpu of the Spring Cloud Deployment. Possible Values are `500m`, `1`, `2`, `3` and `4`. Defaults to `1` if not specified.

-> **Note:** `cpu` supports `500m` and `1` for Basic tier, `500m`, `1`, `2`, `3` and `4` for Standard tier.

* `memory` - (Optional) Specifies the required memory size of the Spring Cloud Deployment. Possible Values are `512Mi`, `1Gi`, `2Gi`, `3Gi`, `4Gi`, `5Gi`, `6Gi`, `7Gi`, and `8Gi`. Defaults to `1Gi` if not specified.

-> **Note:** `memory` supports `512Mi`, `1Gi` and `2Gi` for Basic tier, `512Mi`, `1Gi`, `2Gi`, `3Gi`, `4Gi`, `5Gi`, `6Gi`, `7Gi`, and `8Gi` for Standard tier.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Container Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Container Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Container Deployment.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Container Deployment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Container Deployment.

## Import

Spring Cloud Container Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_container_deployment.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/apps/app1/deployments/deploy1
```
