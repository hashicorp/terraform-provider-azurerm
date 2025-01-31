---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_compute_skus"
description: |-
  Lists available Compute SKUs
---

# Data Source: azurerm_compute_skus

This data source lists available Azure Compute SKUs.

This can be used together with a `precondition` to check if a Virtual Machine SKU is available before the `apply` phase.

## Example Usage

```hcl
data "azurerm_compute_skus" "available" {
  name     = "Standard_D2s_v3"
  location = "westus"
}

output "available_skus" {
  value = {
    for sku in data.azurerm_compute_skus.available.skus : sku.name => sku
  }
}

# Changes to Outputs:
#  + available_skus = {
#      + Standard_D2s_v3 = {
#          + capabilities          = {}
#          + location_restrictions = []
#          + name                  = "Standard_D2s_v3"
#          + resource_type         = "virtualMachines"
#          + size                  = "D2s_v3"
#          + tier                  = "Standard"
#          + zone_restrictions     = []
#          + zones                 = [
#              + "2",
#              + "1",
#              + "3",
#            ]
#        }
#    }
```

## Argument Reference

~> **Note:** Due to API limitations this data source will always get **ALL** available SKUs, regardless of any set filters.

* `location` - (Required) The Azure location of the SKU.

* `name` - (Optional) The name of the SKU, like `Standard_DS2_v2`.

* `include_capabilities` - (Optional) Set to `true` if the SKUs capabilities should be included in the result.

## Attributes Reference

* `skus` - One or more `sku` blocks as defined below.

---

The `sku` block exports the following:

* `name` - The name of the SKU.

* `resource_type` - The resource type of the SKU, like `virtualMachines` or `disks`.

* `tier` - The tier of the SKU.

* `size` - The size of the SKU.

* `capabilities` - If included, this provides a map of the SKUs capabilities.

* `zones` - If the SKU supports Availability Zones, this list contains the IDs of the zones at which the SKU is normally available.

* `location_restrictions` - A list of locations at which the SKU is currently not available. The availability is tied to your Azure subscription.

* `zone_restrictions` - A list of zones in which the SKU is currently not available. The availability is tied to your Azure subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SKUs.
