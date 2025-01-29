---
subcategory: "ApiCenter"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_environment"
description: |-
  Manages an API Center Environment.
---

# azurerm_api_center_environment

Manages an API Center Environment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_apicenter_service" "example" {
  name                = "apicenter-example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_api_center_environment" "example" {
  name                   = "test"
  service_id             = azurerm_apicenter_service.example.id
  identification         = "exampleid"
  environment_type       = "testing"
  description            = "example environment"
  development_portal_uri = "https://developer.com"
  instructions           = "Use this wonderful API to CRUD brilliant data."
  server_type            = "Azure API Management"
  management_portal_uri  = "https://azure-apim-mgmt-portal.azure.com"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this API Center Environment. Changing this forces a new API Center Environment to be created.

* `service_id` - (Required) The ID of the API Center Service where this API Center Environment should exist. Changing this forces a new API Center Environment to be created.

* `identification` - (Required) Identifier of this API Center Environment. Changing this forces a new API Center Environment to be created.

* `environment_type` - (Required) Type of this API Center Environment. Possible values are `development`, `testing`, `staging` and `production`. Changing this forces a new API Center Environment to be created.

* `description` - (Optional) A description of this API Center Environment.

* `development_portal_uri` - (Optional) URI of development portal for this API Center Environment.

* `instructions` - (Optional) Arbitrary auxiliary instructions for this API Center Environment.

* `server_type` - (Optional) Server type for this API Center Environment. Possible values are `AWS API Gateway`, `Apigee API Management`, `Azure API Management`, `Azure compute service`, `Kong API Gateway`, `Kubernetes` and `MuleSoft API Management`.

* `management_portal_uri` - (Optional) URI of management portal for this API Center Environment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Center Environment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Center Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Center Environment.
* `update` - (Defaults to 30 minutes) Used when updating the API Center Environment.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Center Environment.

## Import

API Center Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_center_environment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.ApiCenter/services/example/workspaces/default/environments/example
```
