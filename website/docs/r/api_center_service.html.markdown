---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_service"
subcategory: ""
description: |-
  Documentation for azurerm_api_center_service.
---

# Resource: azurerm_api_center_service

Documentation for azurerm_api_center_service.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `resource_group_name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `location` - (Required) No description provided. Changing this value forces a new resource to be created.

* `identity` - (Optional) (see [nested schema](#schema-6e6573746564-696e707574-6964656e74697479))

* `tags` - (Optional) No description provided.

---

<a id="schema-6e6573746564-696e707574-6964656e74697479"></a>
The `identity` block supports the following:

* `type` - (Required) No description provided. Possible values include `UserAssigned`, `SystemAssigned`, and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) No description provided.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Resource

* `identity` - (see [nested schema](#schema-6e6573746564-6f7574707574-6964656e74697479))

---

<a id="schema-6e6573746564-6f7574707574-6964656e74697479"></a>
The `identity` block exports the following:

* `principal_id` - No description provided.

* `tenant_id` - No description provided.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the resource.
* `update` - (Defaults to 30 minutes) Used when updating the resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the resource.

## Import

This resource can be imported using its resource ID, e.g.

```shell
terraform import azurerm_api_center_service.example "<resource-id>"
```
