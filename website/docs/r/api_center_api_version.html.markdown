---
layout: "azurerm"
page_title: "Resource: azurerm_api_center_api_version"
subcategory: ""
description: |-
  Documentation for azurerm_api_center_api_version.
---

# Resource: azurerm_api_center_api_version

Documentation for azurerm_api_center_api_version.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) No description provided. Changing this value forces a new resource to be created.

* `api_center_api_id` - (Required) No description provided. Changing this value forces a new resource to be created.

* `display_name` - (Optional) No description provided. String length must be between `1` and `50` characters inclusive.

* `lifecycle_stage` - (Required) No description provided. Possible values include `deprecated`, `design`, `development`, `preview`, `production`, `retired`, and `testing`.

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
terraform import azurerm_api_center_api_version.example "<resource-id>"
```
