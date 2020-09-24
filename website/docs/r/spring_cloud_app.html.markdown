---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app"
description: |-
  Manage an Azure Spring Cloud Application.
---

# azurerm_spring_cloud_app

Manage an Azure Spring Cloud Application.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Southeast Asia"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Application. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group in which to create the Spring Cloud Application. Changing this forces a new resource to be created.

* `service_name` - (Required) Specifies the name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Spring Cloud Application. Possible value is `SystemAssigned`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Application.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application.

## Import

Spring Cloud Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/Spring/myservice/apps/myapp
```
