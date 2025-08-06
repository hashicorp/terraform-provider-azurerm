---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_gateway_route_config"
description: |-
  Manages a Spring Cloud Gateway Route Config.
---

# azurerm_spring_cloud_gateway_route_config

Manages a Spring Cloud Gateway Route Config.

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_gateway_route_config` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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
  resource_group_name = azurerm_resource_group.example.name
  service_name        = azurerm_spring_cloud_service.example.name
}

resource "azurerm_spring_cloud_gateway" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
}

resource "azurerm_spring_cloud_gateway_route_config" "example" {
  name                    = "example"
  spring_cloud_gateway_id = azurerm_spring_cloud_gateway.example.id
  spring_cloud_app_id     = azurerm_spring_cloud_app.example.id
  protocol                = "HTTPS"
  route {
    description            = "example description"
    filters                = ["StripPrefix=2", "RateLimit=1,1s"]
    order                  = 1
    predicates             = ["Path=/api5/customer/**"]
    sso_validation_enabled = true
    title                  = "myApp route config"
    token_relay            = true
    uri                    = "https://www.example.com"
    classification_tags    = ["tag1", "tag2"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Gateway Route Config. Changing this forces a new Spring Cloud Gateway Route Config to be created.

* `spring_cloud_gateway_id` - (Required) The ID of the Spring Cloud Gateway. Changing this forces a new Spring Cloud Gateway Route Config to be created.

* `protocol` - (Required) Specifies the protocol of routed Spring Cloud App. Allowed values are `HTTP` and `HTTPS`.

~> **Note:** You likely want to use `HTTPS` in a production environment, since `HTTP` offers no encryption.

* `filters` - (Optional) Specifies a list of filters which are used to modify the request before sending it to the target endpoint, or the received response in app level.

* `predicates` - (Optional) Specifies a list of conditions to evaluate a route for each request in app level. Each predicate may be evaluated against request headers and parameter values. All of the predicates associated with a route must evaluate to true for the route to be matched to the request.

* `sso_validation_enabled` - (Optional) Should the sso validation be enabled in app level?

---

* `route` - (Optional) One or more `route` blocks as defined below.

* `open_api` - (Optional) One or more `open_api` blocks as defined below.

* `spring_cloud_app_id` - (Optional) The ID of the Spring Cloud App.

---

A `route` block supports the following:

* `order` - (Required) Specifies the route processing order.

* `description` - (Optional) Specifies the description which will be applied to methods in the generated OpenAPI documentation.

* `filters` - (Optional) Specifies a list of filters which are used to modify the request before sending it to the target endpoint, or the received response.

* `predicates` - (Optional) Specifies a list of conditions to evaluate a route for each request. Each predicate may be evaluated against request headers and parameter values. All of the predicates associated with a route must evaluate to true for the route to be matched to the request.

* `sso_validation_enabled` - (Optional) Should the sso validation be enabled?

* `classification_tags` - (Optional) Specifies the classification tags which will be applied to methods in the generated OpenAPI documentation.

* `title` - (Optional) Specifies the title which will be applied to methods in the generated OpenAPI documentation.

* `token_relay` - (Optional) Should pass currently-authenticated user's identity token to application service?

* `uri` - (Optional) Specifies the full uri which will override `appName`.

---

A `open_api` block supports the following:

* `uri` - (Optional) The URI of OpenAPI specification.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Gateway Route Config.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Gateway Route Config.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Gateway Route Config.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Gateway Route Config.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Gateway Route Config.

## Import

Spring Cloud Gateway Route Configs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_gateway_route_config.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/routeConfig1
```
