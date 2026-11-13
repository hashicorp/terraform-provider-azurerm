---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_deployment"
subcategory: "ApiCenter"
description: |-
  Documentation for azurerm_api_center_deployment.
---

# azurerm_api_center_deployment

Documentation for azurerm_api_center_deployment.

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

resource "azurerm_api_center_api" "example" {
  name       = "example-api"
  service_id = azurerm_api_center_service.example.id
  kind       = "rest"
}

resource "azurerm_api_center_api_version" "example" {
  name            = "example-version"
  api_id          = azurerm_api_center_api.example.id
  lifecycle_stage = "development"
}

resource "azurerm_api_center_api_definition" "example" {
  name           = "example-definition"
  api_version_id = azurerm_api_center_api_version.example.id
}

resource "azurerm_api_center_environment" "example" {
  name             = "example-environment"
  service_id       = azurerm_api_center_service.example.id
  environment_type = "development"
}

resource "azurerm_api_center_deployment" "example" {
  name               = "example-deployment"
  api_definition_id  = azurerm_api_center_api_definition.example.id
  environment_id     = azurerm_api_center_environment.example.id
  server_runtime_uri = ["https://example.com/api"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the deployment. Changing this value forces a new resource to be created.

* `api_definition_id` - (Required) The ID of the API definition to deploy.

* `custom_properties` - (Optional) A JSON object containing custom properties. Defaults to `{}`.

* `description` - (Optional) A description of the deployment. String length must be between `1` and `500` characters inclusive.

* `display_name` - (Optional) The display name of the deployment. String length must be between `1` and `50` characters inclusive.

* `environment_id` - (Required) The ID of the environment for this deployment.

* `server_runtime_uri` - (Required) A set of server runtime URLs for this deployment.

* `state` - (Optional) The state of the deployment. Possible values include `active` and `inactive`.

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

API Center Deployments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_center_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.ApiCenter/services/example-service/workspaces/default/apis/example-api/deployments/example-deployment
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiCenter` - 2024-03-01
