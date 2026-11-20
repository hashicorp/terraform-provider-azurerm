---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_api"
subcategory: "API Center"
description: |-
  Documentation for azurerm_api_center_api.
---

# azurerm_api_center_api

Documentation for azurerm_api_center_api.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the API. Changing this value forces a new resource to be created.

* `api_center_service_id` - (Required) The ID of the API Center service. Changing this value forces a new resource to be created.

* `api_type` - (Required) The type of API. Possible values include `graphql`, `grpc`, `rest`, `soap`, `webhook`, and `websocket`.

* `contact` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-636f6e74616374))

* `custom_properties_json` - (Optional) A JSON object containing custom properties. Defaults to `{}`.

* `description` - (Optional) A description of the API. String length must be between `1` and `1000` characters inclusive.

* `display_name` - (Optional) The display name of the API. String length must be between `1` and `50` characters inclusive.

* `external_documentation` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-65787465726e616c5f646f63756d656e746174696f6e))

* `license` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-6c6963656e7365))

* `summary` - (Optional) A short summary of the API. String length must be between `1` and `200` characters inclusive.

* `terms_of_service_url` - (Optional) The URL of the API's terms of service.

---

<a id="schema-6e6573746564-696e707574-636f6e74616374"></a>
The `contact` block supports the following:

* `name` - (Required) The name of the API contact. String length must be between `1` and `100` characters inclusive.

* `email` - (Optional) The email address of the API contact. String length must be between `1` and `100` characters inclusive.

* `url` - (Optional) The URL for the API contact.

---

<a id="schema-6e6573746564-696e707574-65787465726e616c5f646f63756d656e746174696f6e"></a>
The `external_documentation` block supports the following:

* `title` - (Required) The title of the external documentation. String length must be between `1` and `50` characters inclusive.

* `description` - (Optional) A description of the external documentation. String length must be between `1` and `500` characters inclusive.

* `url` - (Optional) The URL of the external documentation.

---

<a id="schema-6e6573746564-696e707574-6c6963656e7365"></a>
The `license` block supports the following:

* `name` - (Required) The name of the API license. String length must be between `1` and `100` characters inclusive.

* `identifier` - (Optional) The SPDX identifier for the API license. String length must be between `1` and `50` characters inclusive.

* `url` - (Optional) The URL of the API license.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

* `lifecycle_stage` - The lifecycle stage of the API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

This resource can be imported using its resource ID, e.g.

```shell
terraform import azurerm_api_center_api.example "<resource-id>"
```
