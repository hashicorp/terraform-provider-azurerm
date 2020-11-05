---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration"
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
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Configuration. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Optional) The SKU name of the the App Configuration. Possible values are `free` and `standard`.

* `identity` - (Optional) An `identity` block as defined below.

~> **NOTE:** Azure does not allow a downgrade from `standard` to `free`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the App Configuration. At this time the only allowed value is `SystemAssigned`.

---
## Attributes Reference

The following attributes are exported:

* `id` - The App Configuration ID.

* `endpoint` - The URL of the App Configuration.

* `primary_read_key` - A `primary_read_key` block as defined below containing the primary read access key.

* `primary_write_key` - A `primary_write_key` block as defined below containing the primary write access key.

* `secondary_read_key` - A `secondary_read_key` block as defined below containing the secondary read access key.

* `secondary_write_key` - A `secondary_write_key` block as defined below containing the secondary write access key.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory.

* `tenant_id` - The ID of the Azure Active Directory Tenant.

---

A `primary_read_key` block exports the following:

* `connection_string` - The Connection String for this Access Key - comprising of the Endpoint, ID and Secret.

* `id` - The ID of the Access Key.

* `secret` - The Secret of the Access Key.

---

A `primary_write_key` block exports the following:

* `connection_string` - The Connection String for this Access Key - comprising of the Endpoint, ID and Secret.

* `id` - The ID of the Access Key.

* `secret` - The Secret of the Access Key.

---

A `secondary_read_key` block exports the following:

* `connection_string` - The Connection String for this Access Key - comprising of the Endpoint, ID and Secret.

* `id` - The ID of the Access Key.

* `secret` - The Secret of the Access Key.

---

A `secondary_write_key` block exports the following:

* `connection_string` - The Connection String for this Access Key - comprising of the Endpoint, ID and Secret.

* `id` - The ID of the Access Key.

* `secret` - The Secret of the Access Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the App Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Configuration.

## Import

App Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration.appconf /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1
```
