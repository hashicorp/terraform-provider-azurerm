---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container"
description: |-
  Gets information about an existing Storage Container.
---

# Data Source: azurerm_storage_container

Use this data source to access information about an existing Storage Container.

## Example Usage

```hcl

resource "azurerm_storage_container" "test" {
  name                  = "containerdstest-%s"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
  metadata = {
    key1 = "value1"
    key2 = "value2"
   }
}

data "azurerm_storage_container" "example" {
	storage_container_id = "${azurerm_storage_container.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `storage_container_id` - (Required) Specifies the id of the storage container.

## Attributes Reference

* `container_access_type` - The Access Level configured for this Container.
* `has_immutability_policy` - Is there an Immutability Policy configured on this Storage Container?
* `has_legal_hold` - Is there a Legal Hold configured on this Storage Container?
* `storage_account_name` - The name of the Storage Account where the Container is created.
* `metadata`  - A mapping of MetaData for this Container.
