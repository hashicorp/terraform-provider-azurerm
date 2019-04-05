---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_openid_connect_provider"
sidebar_current: "docs-azurerm-resource-api-management-openid-connect-provider"
description: |-
  Manages an OpenID Connect Provider within a API Management Service.
---

# azurerm_api_management_openid_connect_provider

Manages an OpenID Connect Provider within a API Management Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "test" {
  name                = "example-apim"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "example-provider"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  client_id           = "00001111-2222-3333-4444-555566667777"
  display_name        = "Example Provider"
  metadata_endpoint   = "https://example.com/example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) the Name of the OpenID Connect Provider which should be created within the API Management Service. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service in which this OpenID Connect Provider should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `client_id` - (Required) The Client ID used for the Client Application.

* `client_secret` - (Required) The Client Secret used for the Client Application.

* `display_name` - (Required) A user-friendly name for this OpenID Connect Provider.

* `metadata_endpoint` - (Required) The URI of the Metadata endpoint.

---

* `description` - (Optional) A description of this OpenID Connect Provider.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management OpenID Connect Provider.

## Import

API Management OpenID Connect Providers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_openid_connect_provider.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/openidConnectProviders/provider1
```
