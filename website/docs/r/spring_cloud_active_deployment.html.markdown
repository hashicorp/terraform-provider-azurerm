---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_active_deployment"
description: |-
  Manages an Active Azure Spring Cloud Deployment.
---

# azurerm_spring_cloud_active_deployment

Manages an Active Azure Spring Cloud Deployment.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_active_deployment` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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
  instance_count      = 2
  jvm_options         = "-XX:+PrintGC"
  runtime_version     = "Java_11"

  quota {
    cpu    = "2"
    memory = "4Gi"
  }

  environment_variables = {
    "Env" : "Staging"
  }
}

resource "azurerm_spring_cloud_active_deployment" "example" {
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
  deployment_name     = azurerm_spring_cloud_java_deployment.example.name
}
```

## Argument Reference

The following arguments are supported:

* `spring_cloud_app_id` - (Required) Specifies the id of the Spring Cloud Application. Changing this forces a new resource to be created.

* `deployment_name` - (Required) Specifies the name of Spring Cloud Deployment which is going to be active.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Active Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Active Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Active Deployment.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Active Deployment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Active Deployment.

## Import

Spring Cloud Active Deployment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_active_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AppPlatform/spring/service1/apps/app1
```
