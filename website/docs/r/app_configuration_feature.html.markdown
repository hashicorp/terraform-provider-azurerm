---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration_feature"
description: |-
  Manages an Azure App Configuration Feature.

---

# azurerm_app_configuration_feature

Manages an Azure App Configuration Feature.

-> **Note:** App Configuration Features are provisioned using a Data Plane API which requires the role `App Configuration Data Owner` on either the App Configuration or a parent scope (such as the Resource Group/Subscription). [More information can be found in the Azure Documentation for App Configuration](https://docs.microsoft.com/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration). This is similar to providing App Configuration Keys.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_configuration" "appconf" {
  name                = "appConf1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

data "azurerm_client_config" "current" {}

resource "azurerm_role_assignment" "appconf_dataowner" {
  scope                = azurerm_app_configuration.appconf.id
  role_definition_name = "App Configuration Data Owner"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_app_configuration_feature" "test" {
  configuration_store_id = azurerm_app_configuration.appconf.id
  description            = "test description"
  name                   = "test-ackey"
  label                  = "test-ackeylabel"
  enabled                = true
}
```

## Argument Reference

The following arguments are supported:

* `configuration_store_id` - (Required) Specifies the id of the App Configuration. Changing this forces a new resource to be created.

* `description` - (Optional) The description of the App Configuration Feature. 

* `enabled` - (Optional) The status of the App Configuration Feature. By default, this is set to false.

* `key` - (Optional) The key of the App Configuration Feature. The value for `name` will be used if this is unspecified. Changing this forces a new resource to be created.

* `label` - (Optional) The label of the App Configuration Feature. Changing this forces a new resource to be created.

* `locked` - (Optional) Should this App Configuration Feature be Locked to prevent changes?

* `name` - (Required) The name of the App Configuration Feature. Changing this forces a new resource to be created.

* `percentage_filter_value` - (Optional) A number representing the value of the percentage required to enable this feature.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `targeting_filter` - (Optional) A `targeting_filter` block as defined below.

* `timewindow_filter` - (Optional) A `timewindow_filter` block as defined below.

---

A `targeting_filter` block represents a feature filter of type `Microsoft.Targeting` and takes the following attributes:

* `default_rollout_percentage` - (Required) A number representing the percentage of the entire user base.

* `groups` - (Optional) One or more `groups` blocks as defined below.

* `users` - (Optional) A list of users to target for this feature.

---

A `groups` block represents a group that can be used in a `targeting_filter` and takes the following attributes:

* `name` - (Required) The name of the group.

* `rollout_percentage` - (Required) Rollout percentage of the group.

---

A `timewindow_filter` block represents a feature filter of type `Microsoft.TimeWindow` and takes the following attributes:

* `start` - (Optional) The earliest timestamp the feature is enabled. The timestamp must be in RFC3339 format.

* `end` - (Optional) The latest timestamp the feature is enabled. The timestamp must be in RFC3339 format.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The App Configuration Feature ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the App Configuration Feature.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration Feature.
* `update` - (Defaults to 30 minutes) Used when updating the App Configuration Feature.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Configuration Feature.

## Import

App Configuration Features can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration_feature.test https://appconfname1.azconfig.io/kv/.appconfig.featureflag%2FkeyName?label=labelName
```

If you wish to import with an empty label then simply leave the label's name blank:

```shell
terraform import azurerm_app_configuration_feature.test https://appconfname1.azconfig.io/kv/.appconfig.featureflag%2FkeyName?label=
```
