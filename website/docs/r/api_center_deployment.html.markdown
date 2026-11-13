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

* `name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `api_center_api_definition_id` - (Optional) No description provided.

* `custom_properties` - (Optional) A JSON-encoded object containing custom metadata properties, keyed by metadata schema name. Defaults to `{}`. Removing this argument clears the deployment's custom properties. The corresponding metadata schemas must exist before assigning custom properties.

* `description` - (Optional) No description provided. String length must be between `1` and `500` characters inclusive.

* `display_name` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

* `environment_id` - (Optional) No description provided.

* `server` - (Required) (see [nested schema](#schema-6e6573746564-696e707574-736572766572))

* `state` - (Optional) No description provided. Possible values include `active` and `inactive`.

---

<a id="schema-6e6573746564-696e707574-736572766572"></a>
The `server` block supports the following:

* `runtime_uri` - (Required) No description provided.

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
