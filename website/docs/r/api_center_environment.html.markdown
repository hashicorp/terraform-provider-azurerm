---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_environment"
subcategory: ""
description: |-
  Documentation for azurerm_api_center_environment.
---

# Resource: azurerm_api_center_environment

Documentation for azurerm_api_center_environment.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `api_center_service_id` - (Required) No description provided. Changing this value forces a new resource to be created.

* `description` - (Optional) No description provided.

* `development_portal_uri` - (Optional) No description provided.

* `display_name` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

* `environment_type` - (Required) No description provided. Possible values include `development`, `production`, `staging`, and `testing`. Changing this value forces a new resource to be created.

* `instructions` - (Optional) No description provided.

* `management_portal_uri` - (Optional) No description provided.

* `server_type` - (Optional) No description provided. Possible values include `AWS API Gateway`, `Apigee API Management`, `Azure API Management`, `Azure compute service`, `Kong API Gateway`, `Kubernetes`, and `MuleSoft API Management`.

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
