---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_api"
subcategory: ""
description: |-
  Documentation for azurerm_api_center_api.
---

# Resource: azurerm_api_center_api

Documentation for azurerm_api_center_api.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `api_center_service_id` - (Required) No description provided. Changing this value forces a new resource to be created.

* `contacts` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-636f6e7461637473))

* `description` - (Optional) No description provided. String length must be between `1` and `1000` characters inclusive.

* `display_name` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

* `external_documentation` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-65787465726e616c5f646f63756d656e746174696f6e))

* `kind` - (Required) No description provided. Possible values include `graphql`, `grpc`, `rest`, `soap`, `webhook`, and `websocket`.

* `license` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-6c6963656e7365))

* `summary` - (Optional) No description provided. String length must be between `1` and `200` characters inclusive.

* `terms_of_service_url` - (Optional) No description provided.

---

<a id="schema-6e6573746564-696e707574-636f6e7461637473"></a>
The `contacts` block supports the following:

* `email` - (Optional) No description provided. String length must be between `1` and `100` characters inclusive.

* `name` - (Optional) No description provided. String length must be between `1` and `100` characters inclusive.

* `url` - (Optional) No description provided.

---

<a id="schema-6e6573746564-696e707574-65787465726e616c5f646f63756d656e746174696f6e"></a>
The `external_documentation` block supports the following:

* `url` - (Required) No description provided.

* `description` - (Optional) No description provided. String length must be between `1` and `500` characters inclusive.

* `title` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

---

<a id="schema-6e6573746564-696e707574-6c6963656e7365"></a>
The `license` block supports the following:

* `identifier` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

* `name` - (Optional) No description provided. String length must be between `1` and `100` characters inclusive.

* `url` - (Optional) No description provided.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

* `lifecycle_stage` - No description provided.

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
