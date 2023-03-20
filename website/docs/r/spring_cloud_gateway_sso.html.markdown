---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_gateway_sso"
description: |-
  Manages a Spring Cloud Gateway SSO.
---

# azurerm_spring_cloud_gateway

-> **NOTE:** This resource is applicable only for Spring Cloud Service with enterprise tier.

Manages a Spring Cloud Gateway SSO.

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

resource "azurerm_spring_cloud_gateway" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
}

resource "azurerm_spring_cloud_gateway_sso" "example" {
  spring_cloud_gateway_id = azurerm_spring_cloud_gateway.example.id
  client_id               = "example id"
  client_secret           = "example secret"
  issuer_uri              = "https://www.test.com/issueToken"
  scope                   = ["read"]
}
```

## Arguments Reference

The following arguments are supported:

* `spring_cloud_gateway_id` - (Required) The ID of the Spring Cloud Gateway SSO. Changing this forces a new Spring Cloud Gateway SSO to be created.

---

* `client_id` - (Optional) The public identifier for the application.

* `client_secret` - (Optional) The secret known only to the application and the authorization server.

* `issuer_uri` - (Optional) The URI of Issuer Identifier.

* `scope` - (Optional) It defines the specific actions applications can be allowed to do on a user's behalf.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Gateway SSO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Gateway SSO.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Gateway SSO.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Gateway SSO.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Gateway SSO.

## Import

Spring Cloud Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_gateway_sso.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1
```
