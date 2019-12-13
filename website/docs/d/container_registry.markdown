---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry"
sidebar_current: "docs-azurerm-datasource-container-registry"
description: |-
  Get information about an existing Container Registry

---

# Data Source: azurerm_container_registry

Use this data source to access information about an existing Container Registry.

## Example Usage

```hcl
data "azurerm_container_registry" "example" {
  name                = "testacr"
  resource_group_name = "test"
}

output "login_server" {
  value = "${data.azurerm_container_registry.example.login_server}"
}
```

## Argument Reference

* `name` - (Required) The name of the Container Registry.
* `resource_group_name` - (Required) The Name of the Resource Group where this Container Registry exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Container Registry ID.

* `login_server` - The URL that can be used to log into the container registry.

* `admin_username` - The Username associated with the Container Registry Admin account - if the admin account is enabled.

* `admin_password` - The Password associated with the Container Registry Admin account - if the admin account is enabled.

* `location` - The Azure Region in which this Container Registry exists.

* `admin_enabled` - Is the Administrator account enabled for this Container Registry.

* `sku` - The SKU of this Container Registry, such as `Basic`.

* `storage_account_id` - The ID of the Storage Account used for this Container Registry. This is only returned for `Classic` SKU's.

* `tags` - A map of tags assigned to the Container Registry.
