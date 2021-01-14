---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_java_deployment"
description: |-
  Manages an Azure Spring Cloud Deployment which runtime is Java.
---

# azurerm_spring_cloud_java_deployment

Manages an Azure Spring Cloud Deployment which runtime is Java.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "example-springcloudapp"
  resource_group_name = azurerm_resource_group.example.name
  service_name        = azurerm_spring_cloud_service.example.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_spring_cloud_java_deployment" "example" {
  name                = "deploy1"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
  cpu                 = 2
  instance_count      = 2
  jvm_options         = "-XX:+PrintGC"
  memory_in_gb        = 4
  runtime_version     = "Java_11"

  environment_variables = {
    "Foo" : "Bar"
    "Env" : "Staging"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Deployment. Changing this forces a new resource to be created.

* `spring_cloud_app_id` - (Required) Specifies the id of the Spring Cloud Application in which to create the Deployment. Changing this forces a new resource to be created.

* `cpu` - (Optional) Specifies the required cpu of the Spring Cloud Deployment. Possible Values are between `1` and `4`. Defaults to `1` if not specified.

* `environment_variables` - (Optional) Specifies the environment variables of the Spring Cloud Deployment as a map of key-value pairs.

* `instance_count` - (Optional) Specifies the required instance count of the Spring Cloud Deployment. Possible Values are between `1` and `500`. Defaults to `1` if not specified.

* `jvm_options` - (Optional) Specifies the jvm option of the Spring Cloud Deployment.

* `memory_in_gb` - (Optional) Specifies the required memory size of the Spring Cloud Deployment. Possible Values are between `1` and `8`. Defaults to `1` if not specified.

* `runtime_version` - (Optional) Specifies the runtime version of the Spring Cloud Deployment. Possible Values are `Java_8` and `Java_11`. Defaults to `Java_8`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Deployment.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Deployment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Deployment.

## Import

Spring Cloud Deployment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_java_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AppPlatform/Spring/service1/apps/app1/deployments/deploy1
```
