---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration"
sidebar_current: "docs-azurerm-resource-app-configuration"
description: |-
  Manages an Azure App Configuration.

---

# azurerm_app_configuration

Manages an Azure App Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "resourceGroup1"
  location = "West Europe"
}

resource "azurerm_app_configuration" "appconf" {
  name                = "appConf1"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  location            = "${azurerm_resource_group.rg.location}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Configuration. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Optional) The SKU name of the the App Configuration. Possible values are `free` and `standard`.

~> **NOTE:** Azure does not allow a downgrade from `standard` to `free`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---
## Attributes Reference

The following attributes are exported:

* `id` - The App Configuration ID.

* `endpoint` - The URL of the App Configuration.

* `primary_write_key` - An `access_key` block as defined below containing the primary write access key.

* `secondary_write_key` - An `access_key` block as defined below containing the secondary write access key.

* `primary_read_key` - An `access_key` block as defined below containing the primary read access key.

* `secondary_read_key` - An `access_key` block as defined below containing the secondary read access key.

---

A `access_key` block exports the following:

* `id` - The ID of the access key.

* `secret` - The secret of the access key.

* `connection_string` - The connection string including the endpoint, id and secret.

## Import

App Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration.appconf /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1
```
