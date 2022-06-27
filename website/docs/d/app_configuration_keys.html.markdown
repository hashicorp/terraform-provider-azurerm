---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_app_configuration_keys"
description: |-
  Gets information about existing Azure App Configuration Keys.
---

# Data Source: azurerm_app_configuration_keys

Use this data source to access information about existing Azure App Configuration Keys.

-> **Note:** App Configuration Keys are provisioned using a Data Plane API which requires the role `App Configuration Data Owner` on either the App Configuration or a parent scope (such as the Resource Group/Subscription). [More information can be found in the Azure Documentation for App Configuration](https://docs.microsoft.com/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration).

## Example Usage

```hcl
data "azurerm_app_configuration_keys" "test" {
  configuration_store_id = azurerm_app_configuration.appconf.id
}

output "value" {
  value = data.azurerm_app_configuration_keys.test.items
}
```

## Argument Reference

The following arguments are supported:

* `configuration_store_id` - (Required) Specifies the id of the App Configuration.

* `key` - (Optional) The name of the App Configuration Keys to look up.

* `label` - (Optional) The label of the App Configuration Keys tp look up.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `items` - A list of `items` blocks as defined below.

---

Each element in `items` block exports the following:

* `key` - The name of the App Configuration Key.

* `label` - The label of the App Configuration Key.

* `content_type` - The content type of the App Configuration Key.

* `value` - The value of the App Configuration Key.

* `locked` - Is this App Configuration Key be Locked to prevent changes.

* `type` - The type of the App Configuration Key. It can either be `kv` (simple [key/value](https://docs.microsoft.com/azure/azure-app-configuration/concept-key-value)) or `vault` (where the value is a reference to a [Key Vault Secret](https://azure.microsoft.com/en-gb/services/key-vault/).

* `vault_key_reference` - The ID of the vault secret this App Configuration Key refers to, when `type` is `vault`.

* `tags` - A mapping of tags assigned to the resource.

* `id` - The App Configuration Key ID.

* `etag` - The ETag of the key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration Key.
