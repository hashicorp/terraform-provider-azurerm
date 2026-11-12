---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_workspace"
subcategory: ""
description: |-
  Documentation for azurerm_api_center_workspace.
---

# Resource: azurerm_api_center_workspace

Documentation for azurerm_api_center_workspace.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `api_center_service_id` - (Required) No description provided. Changing this value forces a new resource to be created.

* `description` - (Optional) No description provided.

* `title` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

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
terraform import azurerm_api_center_workspace.example "<resource-id>"
```
