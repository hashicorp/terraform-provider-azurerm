---
subcategory: "Mixed Reality"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spatial_anchors_account"
description: |-
  Get information about an Azure Spatial Anchors Account.
---

# azurerm_spatial_anchors_account

Get information about an Azure Spatial Anchors Account.

~> **Note:** The `azurerm_spatial_anchors_account` data source has been deprecated because the service is retiring from 2024-11-20 and will be removed in v5.0 of the AzureRM Provider.

## Example Usage

```hcl
data "azurerm_spatial_anchors_account" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

output "account_domain" {
  value = data.azurerm_spatial_anchors_account.account_domain
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spatial Anchors Account. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Spatial Anchors Account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spatial Anchors Account.

* `account_domain` - The domain of the Spatial Anchors Account.

* `account_id` - The account ID of the Spatial Anchors Account.

* `tags` - The Tags assigned to this Spatial Anchors Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Spatial Anchors Account.
