---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_deployment"
subcategory: ""
description: |-
  Documentation for azurerm_api_center_deployment.
---

# Resource: azurerm_api_center_deployment

Documentation for azurerm_api_center_deployment.

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

This resource can be imported using its resource ID, e.g.

```shell
terraform import azurerm_api_center_deployment.example "<resource-id>"
```
