---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration_key"
description: |-
  Manages an Azure App Configuration Key.

---

# azurerm_app_configuration_key

Manages an Azure App Configuration Key.

-> **Note:** App Configuration Keys are provisioned using a Data Plane API which requires the role `App Configuration Data Owner` on either the App Configuration or a parent scope (such as the Resource Group/Subscription). [More information can be found in the Azure Documentation for App Configuration](https://docs.microsoft.com/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration).

## Example Usage of `kv` type

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

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.appconf.id
  key                    = "appConfKey1"
  label                  = "somelabel"
  value                  = "a test"

  depends_on = [
    azurerm_role_assignment.appconf_dataowner
  ]
}
```

## Example Usage of `vault` type

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
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "kvs" {
  name         = "kvs"
  value        = "szechuan"
  key_vault_id = azurerm_key_vault.kv.id
}

resource "azurerm_role_assignment" "appconf_dataowner" {
  scope                = azurerm_app_configuration.appconf.id
  role_definition_name = "App Configuration Data Owner"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.test.id
  key                    = "key1"
  type                   = "vault"
  label                  = "label1"
  vault_key_reference    = azurerm_key_vault_secret.kvs.versionless_id

  depends_on = [
    azurerm_role_assignment.appconf_dataowner
  ]
}
```

## Argument Reference

The following arguments are supported:

* `configuration_store_id` - (Required) Specifies the id of the App Configuration. Changing this forces a new resource to be created.

* `key` - (Required) The name of the App Configuration Key to create. Changing this forces a new resource to be created.

* `content_type` - (Optional) The content type of the App Configuration Key. This should only be set when type is set to `kv`.

* `label` - (Optional) The label of the App Configuration Key. Changing this forces a new resource to be created.

* `value` - (Optional) The value of the App Configuration Key. This should only be set when type is set to `kv`.

~> **Note:** `value` and `vault_key_reference` are mutually exclusive.

* `locked` - (Optional) Should this App Configuration Key be Locked to prevent changes?

* `type` - (Optional) The type of the App Configuration Key. It can either be `kv` (simple [key/value](https://docs.microsoft.com/azure/azure-app-configuration/concept-key-value)) or `vault` (where the value is a reference to a [Key Vault Secret](https://azure.microsoft.com/en-gb/services/key-vault/). Defaults to `kv`.

* `vault_key_reference` - (Optional) The ID of the vault secret this App Configuration Key refers to. This should only be set when `type` is set to `vault`.

~> **Note:** `vault_key_reference` and `value` are mutually exclusive.

~> **Note:** When setting the `vault_key_reference` using the `id` will pin the value to specific version of the secret, to reference latest secret value use `versionless_id`

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The App Configuration Key ID.

* `etag` - (Optional) The ETag of the key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the App Configuration Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration Key.
* `update` - (Defaults to 30 minutes) Used when updating the App Configuration Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Configuration Key.

## Import

App Configuration Keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration_key.test https://appconfname1.azconfig.io/kv/keyName?label=labelName
```

If you wish to import a key with an empty label then simply leave label's name blank:

```shell
terraform import azurerm_app_configuration_key.test https://appconfname1.azconfig.io/kv/keyName?label=
```
