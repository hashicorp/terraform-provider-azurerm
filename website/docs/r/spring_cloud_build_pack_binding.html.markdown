---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_build_pack_binding"
description: |-
  Manages a Spring Cloud Build Pack Binding.
---

# azurerm_spring_cloud_build_pack_binding

Manages a Spring Cloud Build Pack Binding.

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_build_pack_binding` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_builder" "example" {
  name                    = "example"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id

  build_pack_group {
    name           = "mix"
    build_pack_ids = ["tanzu-Build Packs/java-azure"]
  }

  stack {
    id      = "io.Build Packs.stacks.bionic"
    version = "base"
  }
}

resource "azurerm_spring_cloud_build_pack_binding" "example" {
  name                    = "example"
  spring_cloud_builder_id = azurerm_spring_cloud_builder.example.id
  binding_type            = "ApplicationInsights"
  launch {
    properties = {
      abc           = "def"
      any-string    = "any-string"
      sampling-rate = "12.0"
    }

    secrets = {
      connection-string = "XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXX;XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXXXXXXXX"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Build Pack Binding. Changing this forces a new Spring Cloud Build Pack Binding to be created.

* `spring_cloud_builder_id` - (Required) The ID of the Spring Cloud Builder. Changing this forces a new Spring Cloud Build Pack Binding to be created.

---

* `binding_type` - (Optional) Specifies the Build Pack Binding Type. Allowed values are `ApacheSkyWalking`, `AppDynamics`, `ApplicationInsights`, `Dynatrace`, `ElasticAPM` and `NewRelic`.

* `launch` - (Optional) A `launch` block as defined below.

---

A `launch` block supports the following:

* `properties` - (Optional) Specifies a map of non-sensitive properties for launchProperties.

* `secrets` - (Optional) Specifies a map of sensitive properties for launchProperties.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Build Pack Binding.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Build Pack Binding.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Build Pack Binding.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Build Pack Binding.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Build Pack Binding.

## Import

Spring Cloud Build Pack Bindings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_build_pack_binding.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/buildPackBindings/binding1
```
