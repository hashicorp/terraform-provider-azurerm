---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_environment"
subcategory: "API Center"
description: |-
  Documentation for azurerm_api_center_environment.
---

# azurerm_api_center_environment

Documentation for azurerm_api_center_environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_api_center_service" "example" {
  name                = "example-service"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_api_center_environment" "example" {
  name             = "example-environment"
  service_id       = azurerm_api_center_service.example.id
  environment_type = "development"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the environment. Changing this value forces a new resource to be created.

* `api_center_service_id` - (Required) The ID of the API Center service. Changing this value forces a new resource to be created.

* `environment_type` - (Required) The type of the environment. Possible values include `development`, `production`, `staging`, and `testing`.

* `custom_properties_json` - (Optional) A JSON object containing custom properties. Defaults to `{}`.

* `description` - (Optional) A description of the environment.

* `development_portal_url` - (Optional) The HTTPS URL of the development portal.

* `display_name` - (Optional) The display name of the environment. String length must be between `1` and `50` characters inclusive.

* `instructions` - (Optional) Instructions for using the environment.

* `management_portal_url` - (Optional) The HTTPS URL of the management portal.

* `server_type` - (Optional) The type of server used by the environment. Possible values include `AWS API Gateway`, `Apigee API Management`, `Azure API Management`, `Azure compute service`, `Kong API Gateway`, `Kubernetes`, and `MuleSoft API Management`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

This resource can be imported using its resource ID, e.g.

```shell
terraform import azurerm_api_center_environment.example "<resource-id>"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiCenter` - 2024-03-01
