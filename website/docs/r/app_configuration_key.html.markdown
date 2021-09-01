---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration_key"
description: |-
  Manages an Azure App Configuration Key.

---

# azurerm_app_configuration_key

Manages an Azure App Configuration Key.

## Example Usage of `kv` type

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

## Example Usage of `vault` type
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

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "kv" {
  name                       = "kv"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "get",
    ]

    secret_permissions = [
      "set",
      "get",
      "delete",
      "purge",
      "recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "kvs" {
  name         = "kvs"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.kv.id
}

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "key1"
  type                   = "vault"
  label                  = "label1"
  vault_key_reference    = azurerm_key_vault_secret.kvs.id
}

```

## Argument Reference

The following arguments are supported:

* `configuration_store_id` - (Required) Specifies the id of the App Configuration. Changing this forces a new resource to be created.

* `key` - (Required) The name of the App Configuration Key to create. Changing this forces a new resource to be created.

* `content_type` - (Optional) The content type of the App Configuration Key. This should only be set when type is set to `kv`.

* `label` - (Optional) The label of the App Configuration Key.  Changing this forces a new resource to be created.

* `value` - (Optional) The value of the App Configuration Key. This should only be set when type is set to `kv`.

* `locked` - (Optional) Should this App Configuration Key be Locked to prevent changes?

* `type` - (Optional) The type of the App Configuration Key. It can either be `kv` (simple [key/value](https://docs.microsoft.com/en-us/azure/azure-app-configuration/concept-key-value)) or `vault` (where the value is a reference to a [Key Vault Secret](https://azure.microsoft.com/en-gb/services/key-vault/). 

* `vault_key_reference` - (Optional) The ID of the vault secret this App Configuration Key refers to, when `type` is set to `vault`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---
## Attributes Reference

The following attributes are exported:

* `id` - The App Configuration Key ID.

* `etag` - The ETag of the key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Configuration Key.
* `update` - (Defaults to 30 minutes) Used when updating the App Configuration Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Configuration Key.

## Import

App Configuration Keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration_key.test /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resourceGroup1/providers/Microsoft.AppConfiguration/configurationStores/appConf1/AppConfigurationKey/appConfKey1/Label/label1
```
