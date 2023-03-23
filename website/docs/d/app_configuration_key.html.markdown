---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_app_configuration_key"
description: |-
  Gets information about an existing Azure App Configuration Key.
---

# Data Source: azurerm_app_configuration_key

Use this data source to access information about an existing Azure App Configuration Key.

-> **Note:** App Configuration Keys are provisioned using a Data Plane API which requires the role `App Configuration Data Owner` on either the App Configuration or a parent scope (such as the Resource Group/Subscription). [More information can be found in the Azure Documentation for App Configuration](https://docs.microsoft.com/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration).

## Example Usage

```hcl
data "azurerm_app_configuration_key" "test" {
  configuration_store_id = azurerm_app_configuration.appconf.id
  key                    = "appConfKey1"
  label                  = "somelabel"
}

output "value" {
  value = data.azurerm_app_configuration_key.test.value
}
```

## Argument Reference

The following arguments are supported:

* `configuration_store_id` - (Required) Specifies the id of the App Configuration.

* `key` - (Required) The name of the App Configuration Key.

* `label` - (Optional) The label of the App Configuration Key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `content_type` - The content type of the App Configuration Key.

* `value` - The value of the App Configuration Key.

* `locked` - Is this App Configuration Key be Locked to prevent changes.

* `type` - The type of the App Configuration Key. It can either be `kv` (simple [key/value](https://docs.microsoft.com/azure/azure-app-configuration/concept-key-value)) or `vault` (where the value is a reference to a [Key Vault Secret](https://azure.microsoft.com/en-gb/services/key-vault/).

* `vault_key_reference` - The ID of the vault secret this App Configuration Key refers to, when `type` is `vault`.

* `tags` - A mapping of tags assigned to the resource.

* `id` - The App Configuration Key ID.

* `etag` - The ETag of the key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration Key.
