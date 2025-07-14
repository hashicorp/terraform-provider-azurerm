---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center_catalog"
description: |-
  Gets information about an existing Dev Center Catalog.
---

# Data Source: azurerm_dev_center_catalog

Use this data source to access information about an existing Dev Center Catalog.

## Example Usage

```hcl
data "azurerm_dev_center_catalog" "example" {
  name          = azurerm_dev_center_catalog.example.name
  dev_center_id = azurerm_dev_center_catalog.example.dev_center_id
}

output "id" {
  value = data.azurerm_dev_center_catalog.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center Catalog.

* `dev_center_id` - (Required) Specifies the Dev Center Id within which this Dev Center Catalog should exist.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Catalog.

* `catalog_github` - A `catalog_github` block as defined below.

* `catalog_adogit` - A `catalog_adogit` block as defined below.

---

The `catalog_github` block exports the following:

* `branch` - The Git branch of the Dev Center Catalog.

* `path` - The folder where the catalog items can be found inside the repository.

* `key_vault_key_url` - A reference to the Key Vault secret containing a security token to authenticate to a Git repository.

* `uri` - The Git URI of the Dev Center Catalog.

---

The `catalog_adogit` block exports the following:

* `branch` - The Git branch of the Dev Center Catalog.

* `path` - The folder where the catalog items can be found inside the repository.

* `key_vault_key_url` - A reference to the Key Vault secret containing a security token to authenticate to a Git repository.

* `uri` - The Git URI of the Dev Center Catalog.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Catalog.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevCenter`: 2025-02-01
