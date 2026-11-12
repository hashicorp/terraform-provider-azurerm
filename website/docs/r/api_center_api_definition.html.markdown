---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_api_definition"
subcategory: ""
description: |-
  Documentation for azurerm_api_center_api_definition.
---

# Resource: azurerm_api_center_api_definition

Documentation for azurerm_api_center_api_definition.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `api_center_api_version_id` - (Required) No description provided. Changing this value forces a new resource to be created.

* `description` - (Optional) No description provided.

* `display_name` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

* `specification` - (see [nested schema](#schema-6e6573746564-6f7574707574-73706563696669636174696f6e))

---

<a id="schema-6e6573746564-6f7574707574-73706563696669636174696f6e"></a>
The `specification` object exports the following:

* `name` - No description provided.

* `version` - No description provided.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

This resource can be imported using its resource ID, e.g.

```shell
terraform import azurerm_api_center_api_definition.example "<resource-id>"
```
