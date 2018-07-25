---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry"
sidebar_current: "docs-azurerm-datasource-image"
description: |-
  Get information about an Image

---

# Data Source: azurerm_container_registry

Use this data source to access information about a Container Registry.

## Example Usage

```hcl
data "azurerm_container_registry" "test" {
  name                = "testacr"
  resource_group_name = "test"
}

output "image_id" {
  value = "${data.azurerm_container_registry.test.id}"
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
