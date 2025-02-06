---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_app_configuration"
description: |-
  Gets information about an existing App Configuration.
---

# Data Source: azurerm_app_configuration

Use this data source to access information about an existing App Configuration.

## Example Usage

```hcl
data "azurerm_app_configuration" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_app_configuration.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name of this App Configuration.

* `resource_group_name` - (Required) The name of the Resource Group where the App Configuration exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Configuration.

* `endpoint` - The Endpoint used to access this App Configuration.

* `encryption` - An `encryption` block as defined below.

* `local_auth_enabled` - Whether local authentication methods is enabled.

* `location` - The Azure Region where the App Configuration exists.

* `primary_read_key` - A `primary_read_key` block as defined below containing the primary read access key.

* `primary_write_key` - A `primary_write_key` block as defined below containing the primary write access key.

* `public_network_access` - The Public Network Access setting of this App Configuration.

* `purge_protection_enabled` - Whether Purge Protection is enabled.

* `replica` - One or more `replica` blocks as defined below.

* `secondary_read_key` - A `secondary_read_key` block as defined below containing the secondary read access key.

* `secondary_write_key` - A `secondary_write_key` block as defined below containing the secondary write access key.

* `sku` - The name of the SKU used for this App Configuration.

* `soft_delete_retention_days` - The number of days that items should be retained for once soft-deleted.

* `tags` - A mapping of tags assigned to the App Configuration.

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

A `replica` block exports the following:

* `id` - The ID of the App Configuration Replica.

* `endpoint` - The URL of the App Configuration Replica.

* `location` - The supported Azure location where the App Configuration Replica exists.

* `name` - The name of the App Configuration Replica.


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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the App Configuration.
