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
data "azurerm_storage_container" "example" {
  name                 = "example-container-name"
  storage_account_name = "example-storage-account-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Container.
* `storage_account_name` - (Required) The name of the Storage Account where the Container was created.

## Attributes Reference

* `container_access_type` - The Access Level configured for this Container.
* `has_immutability_policy` - Is there an Immutability Policy configured on this Storage Container?
* `has_legal_hold` - Is there a Legal Hold configured on this Storage Container?
* `metadata`  - A mapping of MetaData for this Container.
