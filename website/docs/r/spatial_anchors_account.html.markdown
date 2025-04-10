---
subcategory: "Mixed Reality"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spatial_anchors_account"
description: |-
  Manages an Azure Spatial Anchors Account.
---

# azurerm_spatial_anchors_account

Manages an Azure Spatial Anchors Account.

~> **Note:** The `azurerm_spatial_anchors_account` resource has been deprecated because the service is retiring from 2024-11-20 and will be removed in v5.0 of the AzureRM Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_spatial_anchors_account" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spatial Anchors Account. Changing this forces a new resource to be created. Must be globally unique.

* `resource_group_name` - (Required) The name of the resource group in which to create the Spatial Anchors Account. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spatial Anchors Account.

* `account_domain` - The domain of the Spatial Anchors Account.

* `account_id` - The account ID of the Spatial Anchors Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spatial Anchors Account.
* `update` - (Defaults to 30 minutes) Used when updating the Spatial Anchors Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spatial Anchors Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spatial Anchors Account.

## Import

Spatial Anchors Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spatial_anchors_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.MixedReality/spatialAnchorsAccounts/example
```
