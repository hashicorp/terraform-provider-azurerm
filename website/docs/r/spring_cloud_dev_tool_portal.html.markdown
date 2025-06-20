---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_dev_tool_portal"
description: |-
  Manages a Spring Cloud Dev Tool Portal.
---

# azurerm_spring_cloud_dev_tool_portal

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

Manages a Spring Cloud Dev Tool Portal.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_dev_tool_portal` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

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

resource "azurerm_spring_cloud_dev_tool_portal" "example" {
  name                          = "default"
  spring_cloud_service_id       = azurerm_spring_cloud_service.example.id
  public_network_access_enabled = true

  sso {
    client_id     = "example id"
    client_secret = "example secret"
    metadata_url  = "https://login.microsoftonline.com/${data.azurerm_client_config.current.tenant_id}/v2.0/.well-known/openid-configuration"
    scope         = ["openid", "profile", "email"]
  }

  application_accelerator_enabled = true
  application_live_view_enabled   = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Dev Tool Portal. The only possible value is `default`. Changing this forces a new Spring Cloud Dev Tool Portal to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new Spring Cloud Dev Tool Portal to be created.

---

* `application_accelerator_enabled` - (Optional) Should the Accelerator plugin be enabled?

* `application_live_view_enabled` - (Optional) Should the Application Live View be enabled?

* `public_network_access_enabled` - (Optional) Is public network access enabled?

* `sso` - (Optional) A `sso` block as defined below.

---

A `sso` block supports the following:

* `client_id` - (Optional) Specifies the public identifier for the application.

* `client_secret` - (Optional) Specifies the secret known only to the application and the authorization server.

* `metadata_url` - (Optional) Specifies the URI of a JSON file with generic OIDC provider configuration.

* `scope` - (Optional) Specifies a list of specific actions applications can be allowed to do on a user's behalf.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Spring Cloud Dev Tool Portal.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Dev Tool Portal.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Dev Tool Portal.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Dev Tool Portal.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Dev Tool Portal.

## Import

Spring Cloud Dev Tool Portals can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_dev_tool_portal.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/Spring/service1/DevToolPortals/default
```
