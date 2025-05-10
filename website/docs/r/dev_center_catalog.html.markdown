---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center_catalog"
description: |-
  Manages a Dev Center Catalog.
---
# azurerm_dev_center_catalog

Manages a Dev Center Catalog.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_dev_center" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dev_center_catalog" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.test.name
  dev_center_id       = azurerm_dev_center.test.id
  catalog_github {
    branch            = "foo"
    path              = ""
    uri               = "example URI"
    key_vault_key_url = "secret"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Dev Center Catalog. Changing this forces a new Dev Center to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Dev Center Catalog should exist. Changing this forces a new Dev Center to be created.

* `dev_center_id` - (Required) Specifies the Dev Center Id within which this Dev Center Catalog should exist. Changing this forces a new Dev Center Catalog to be created.

* `catalog_github` - (Optional) A `catalog_github` block as defined below.

* `catalog_adogit` - (Optional) A `catalog_adogit` block as defined below.

---

The `catalog_github` block supports the following:

* `branch` - (Required) The Git branch of the Dev Center Catalog.

* `path` - (Required) The folder where the catalog items can be found inside the repository.

* `key_vault_key_url` - (Required) A reference to the Key Vault secret containing a security token to authenticate to a Git repository.

* `uri` - (Required) The Git URI of the Dev Center Catalog.

---

The `catalog_adogit` block supports the following:

* `branch` - (Required) The Git branch of the Dev Center Catalog.

* `path` - (Required) The folder where the catalog items can be found inside the repository.

* `key_vault_key_url` - (Required) A reference to the Key Vault secret containing a security token to authenticate to a Git repository.

* `uri` - (Required) The Git URI of the Dev Center Catalog.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center Catalog.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center Catalog.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center Catalog.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center Catalog.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center Catalog.

## Import

An existing Dev Center Catalog can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center_catalog.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{devCenterName}/catalogs/{catalogName}
```

* Where `{subscriptionId}` is the ID of the Azure Subscription where the Dev Center exists. For example `12345678-1234-9876-4563-123456789012`.
* Where `{resourceGroupName}` is the name of Resource Group where this Dev Center exists. For example `example-resource-group`.
* Where `{devCenterName}` is the name of the Dev Center. For example `devCenterValue`.
* Where `{catalogName}` is the name of the Dev Center Catalog. For example `catalogValue`.
