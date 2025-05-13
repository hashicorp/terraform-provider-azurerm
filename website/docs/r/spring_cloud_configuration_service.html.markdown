---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_configuration_service"
description: |-
  Manages a Spring Cloud Configuration Service.
---

# azurerm_spring_cloud_configuration_service

Manages a Spring Cloud Configuration Service.

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_configuration_service` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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

resource "azurerm_spring_cloud_configuration_service" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
  repository {
    name                     = "fake"
    label                    = "master"
    patterns                 = ["app/dev"]
    uri                      = "https://github.com/Azure-Samples/piggymetrics"
    search_paths             = ["dir1", "dir2"]
    strict_host_key_checking = false
    username                 = "adminuser"
    password                 = "H@Sh1CoR3!"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Configuration Service. The only possible value is `default`. Changing this forces a new Spring Cloud Configuration Service to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new Spring Cloud Configuration Service to be created.

---

* `generation` - (Optional) The generation of the Spring Cloud Configuration Service. Possible values are `Gen1` and `Gen2`.

* `refresh_interval_in_seconds` - (Optional) Specifies how often to check repository updates. Minimum value is 0.

* `repository` - (Optional) One or more `repository` blocks as defined below.

---

A `repository` block supports the following:

* `label` - (Required) Specifies the label of the repository.

* `name` - (Required) Specifies the name which should be used for this repository.

* `patterns` - (Required) Specifies the collection of patterns of the repository.

* `uri` - (Required) Specifies the URI of the repository.

* `ca_certificate_id` - (Optional) Specifies the ID of the Certificate Authority used when retrieving the Git Repository via HTTPS.

* `host_key` - (Optional) Specifies the SSH public key of git repository.

* `host_key_algorithm` - (Optional) Specifies the SSH key algorithm of git repository.

* `password` - (Optional) Specifies the password of git repository basic auth.

* `private_key` - (Optional) Specifies the SSH private key of git repository.

* `search_paths` - (Optional) Specifies a list of searching path of the repository

* `strict_host_key_checking` - (Optional) Specifies whether enable the strict host key checking.

* `username` - (Optional) Specifies the username of git repository basic auth.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Configuration Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Configuration Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Configuration Service.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Configuration Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Configuration Service.

## Import

Spring Cloud Configuration Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_configuration_service.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/configurationServices/configurationService1
```
