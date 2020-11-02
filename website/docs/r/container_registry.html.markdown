---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry"
description: |-
  Manages an Azure Container Registry.

---

# azurerm_container_registry

Manages an Azure Container Registry.

~> **Note:** All arguments including the access key will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_container_registry" "acr" {
  name                     = "containerRegistry1"
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  sku                      = "Premium"
  admin_enabled            = false
  georeplication_locations = ["East US", "West Europe"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container Registry. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Registry. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `admin_enabled` - (Optional) Specifies whether the admin user is enabled. Defaults to `false`.

* `storage_account_id` - (Required for `Classic` Sku - Forbidden otherwise) The ID of a Storage Account which must be located in the same Azure Region as the Container Registry.  Changing this forces a new resource to be created.

* `sku` - (Optional) The SKU name of the container registry. Possible values are  `Basic`, `Standard` and `Premium`. `Classic` (which was previously `Basic`) is supported only for existing resources.

~> **NOTE:** The `Classic` SKU is Deprecated and will no longer be available for new resources from the end of March 2019.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `georeplication_locations` - (Optional) A list of Azure locations where the container registry should be geo-replicated.

~> **NOTE:** The `georeplication_locations` is only supported on new resources with the `Premium` SKU.

~> **NOTE:** The `georeplication_locations` list cannot contain the location where the Container Registry exists.

* `network_rule_set` - (Optional) A `network_rule_set` block as documented below.

* `retention_policy` - (Optional) A `retention_policy` block as documented below.

* `trust_policy` - (Optional) A `trust_policy` block as documented below.

~> **NOTE:** `retention_policy` and `trust_policy` are only supported on resources with the `Premium` SKU.

`network_rule_set` supports the following:

* `default_action` - (Optional) The behaviour for requests matching no rules. Either `Allow` or `Deny`. Defaults to `Allow`

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

* `virtual_network` - (Optional) One or more `virtual_network` blocks as defined below.

~> **NOTE:** `network_rule_set ` is only supported with the `Premium` SKU at this time.

~> **NOTE:** Azure automatically configures Network Rules - to remove these you'll need to specify an `network_rule_set` block with `default_action` set to `Deny`.

`ip_rule` supports the following:

* `action` - (Required) The behaviour for requests matching this rule. At this time the only supported value is `Allow`

* `ip_range` - (Required) The CIDR block from which requests will match the rule.

`virtual_network` supports the following:

* `action` - (Required) The behaviour for requests matching this rule. At this time the only supported value is `Allow`

* `subnet_id` - (Required) The subnet id from which requests will match the rule.

`trust_policy` supports the following:

* `enabled` - (Optional) Boolean value that indicates whether the policy is enabled.

`retention_policy` supports the following:

* `days` - (Optional) The number of days to retain an untagged manifest after which it gets purged. Default is `7`.

* `enabled` - (Optional) Boolean value that indicates whether the policy is enabled.

---
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Container Registry.

* `login_server` - The URL that can be used to log into the container registry.

* `admin_username` - The Username associated with the Container Registry Admin account - if the admin account is enabled.

* `admin_password` - The Password associated with the Container Registry Admin account - if the admin account is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry.

## Import

Container Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1
```
