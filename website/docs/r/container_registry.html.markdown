---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry"
sidebar_current: "docs-azurerm-resource-container-registry"
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
  resource_group_name      = "${azurerm_resource_group.rg.name}"
  location                 = "${azurerm_resource_group.rg.location}"
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

* `storage_account_id` - (Required for `Classic` Sku - Forbidden otherwise) The ID of a Storage Account which must be located in the same Azure Region as the Container Registry.

* `sku` - (Optional) The SKU name of the container registry. Possible values are `Classic` (which was previously `Basic`), `Basic`, `Standard` and `Premium`.

~> **NOTE:** The `Classic` SKU is Deprecated and will no longer be available for new resources from the end of March 2019.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `georeplication_locations` - (Optional) A list of Azure locations where the container registry should be geo-replicated.

* `network_rule_set` - (Optional) A `network_rule_set` block as documented below.

`network_rule_set` supports the following:

* `default_action` - (Optional) The behaviour for requests matching no rules. Either `Allow` or `Deny`. Defaults to `Allow`

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

* `virtual_network` - (Optional) One or more `virtual_network` blocks as defined below.

~> **NOTE:** `network_rule_set ` is only supported with the `Premium` SKU at this time.

`ip_rule` supports the following:

* `action` - (Required) The behaviour for requests matching this rule. At this time the only supported value is `Allow`

* `ip_range` - (Required) The CIDR block from which requests will match the rule.

`virtual_network` supports the following:

* `action` - (Required) The behaviour for requests matching this rule. At this time the only supported value is `Allow`

* `subnet_id` - (Required) The subnet id from which requests will match the rule.


---
## Attributes Reference

The following attributes are exported:

* `id` - The Container Registry ID.

* `login_server` - The URL that can be used to log into the container registry.

* `admin_username` - The Username associated with the Container Registry Admin account - if the admin account is enabled.

* `admin_password` - The Password associated with the Container Registry Admin account - if the admin account is enabled.

## Import

Container Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1
```
