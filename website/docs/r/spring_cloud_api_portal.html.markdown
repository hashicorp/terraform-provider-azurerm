---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_api_portal"
description: |-
  Manages a Spring Cloud API Portal.
---

# azurerm_spring_cloud_api_portal

Manages a Spring Cloud API Portal.

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_api_portal` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_gateway" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
}

resource "azurerm_spring_cloud_api_portal" "example" {
  name                          = "default"
  spring_cloud_service_id       = azurerm_spring_cloud_service.example.id
  gateway_ids                   = [azurerm_spring_cloud_gateway.example.id]
  https_only_enabled            = false
  public_network_access_enabled = true
  instance_count                = 1
  api_try_out_enabled           = true
  sso {
    client_id     = "test"
    client_secret = "secret"
    issuer_uri    = "https://www.example.com/issueToken"
    scope         = ["read"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud API Portal. Changing this forces a new Spring Cloud API Portal to be created. The only possible value is `default`.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new Spring Cloud API Portal to be created.

---

* `api_try_out_enabled` - (Optional) Specifies whether the API try-out feature is enabled. When enabled, users can try out the API by sending requests and viewing responses in API portal.

* `gateway_ids` - (Optional) Specifies a list of Spring Cloud Gateway.

* `https_only_enabled` - (Optional) is only https is allowed?

* `instance_count` - (Optional) Specifies the required instance count of the Spring Cloud API Portal. Possible Values are between `1` and `500`. Defaults to `1` if not specified.

* `public_network_access_enabled` - (Optional) Is the public network access enabled?

* `sso` - (Optional) A `sso` block as defined below.

---

A `sso` block supports the following:

* `client_id` - (Optional) The public identifier for the application.

* `client_secret` - (Optional) The secret known only to the application and the authorization server.

* `issuer_uri` - (Optional) The URI of Issuer Identifier.

* `scope` - (Optional) It defines the specific actions applications can be allowed to do on a user's behalf.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud API Portal.

* `url` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud API Portal.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud API Portal.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud API Portal.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud API Portal.

## Import

Spring Cloud API Portals can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_api_portal.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/apiPortals/apiPortal1
```
