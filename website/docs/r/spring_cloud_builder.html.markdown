---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_builder"
description: |-
  Manages a Spring Cloud Builder.
---

# azurerm_spring_cloud_builder

Manages a Spring Cloud Builder.

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_builder` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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
    build_pack_ids = ["tanzu-buildpacks/java-azure"]
  }

  stack {
    id      = "io.buildpacks.stacks.bionic"
    version = "base"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Builder. Changing this forces a new Spring Cloud Builder to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new Spring Cloud Builder to be created.

* `build_pack_group` - (Required) One or more `build_pack_group` blocks as defined below.

* `stack` - (Required) A `stack` block as defined below.

---

A `build_pack_group` block supports the following:

* `name` - (Required) The name which should be used for this build pack group.

* `build_pack_ids` - (Optional) Specifies a list of the build pack's ID.

---

A `stack` block supports the following:

* `id` - (Required) Specifies the ID of the ClusterStack.

* `version` - (Required) Specifies the version of the ClusterStack

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Builder.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Builder.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Builder.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Builder.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Builder.

## Import

Spring Cloud Builders can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_builder.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1
```
