---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration_key"
description: |-
  Manages an Azure App Configuration Key.

---

# azurerm_app_configuration_key

Manages an Azure App Configuration Key.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_configuration" "appconf" {
  name                = "appConf1"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
}

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.appconf.id
  key                    = "appConfKey1"
  label                  = "somelabel"
  value                  = "a test"
}

```

## Argument Reference

The following arguments are supported:

* `configuration_store_id` - (Required) Specifies the id of the App Configuration. Changing this forces a new resource to be created.

* `key` - (Required) The name of the App Configuration Key to create. Changing this forces a new resource to be created.

* `content_type` - (Optional) The content type of the App Configuration Key.

* `label` - (Optional) The label of the App Configuration Key.  Changing this forces a new resource to be created.

* `value` - (Optional) The value of the App Configuration Key.

* `locked` - (Optional) The lock status of the App Configuration Key. Can be either `true` or `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---
## Attributes Reference

The following attributes are exported:

* `id` - The App Configuration Key ID.

* `etag` - The ETag of the key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Configuration Key.
* `update` - (Defaults to 30 minutes) Used when updating the App Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Configuration.

## Import

App Configuration Keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration_key.test /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationKey/appConfKey1/Label/someLabel
```
