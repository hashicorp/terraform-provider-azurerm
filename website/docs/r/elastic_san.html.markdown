---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_san"
description: |-
  Manages an Azure Elastic SAN (Storage Area Network) resource.
---

# azurerm_elastic_san

Manages an Elastic SAN resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_elastic_san" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  base_size_in_tib     = 1
  extended_size_in_tib = 2
  sku {
    name = "example-value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Elastic SAN resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Elastic SAN resource should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Elastic SAN resource should exist. Changing this forces a new resource to be created.

* `base_size_in_tib` - (Required) Specifies the base size of the Elastic SAN resource in TiB. Possible values are between `1` and `100`.

-> **Note:** When updating `base_size_in_tib`, the new value should be greater than the existing one.

* `sku` - (Required) A `sku` block as defined below.

* `extended_size_in_tib` - (Optional) Specifies the extended size of the Elastic SAN resource in TiB. Possible values are between `1` and `100`.

-> **Note:** `extended_size_in_tib` cannot be removed and when updating, the new value should be greater than the existing one.

* `zones` - (Optional) Logical zone for the Elastic SAN resource. Changing this forces a new resource to be created.

-> **Note:** `zones` cannot be specified if `sku.name` is set to `Premium_ZRS`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Elastic SAN resource.

---

The `sku` block supports the following arguments:

* `name` - (Required) The SKU name. Possible values are `Premium_LRS` and `Premium_ZRS`. Changing this forces a new resource to be created.

-> **Note:** `Premium_ZRS` SKU is only available in limited Azure regions including `France Central`, `North Europe`, `West Europe`, and `West US 2`. Please refer to this [document](https://azure.microsoft.com/updates/regional-expansion-azure-elastic-san-public-preview-is-now-available-in-more-regions) for more details.

* `tier` - (Optional) The SKU tier. The only possible value is `Premium`. Defaults to `Premium`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elastic SAN resource.

* `total_iops` - Total Provisioned IOps of the Elastic SAN resource.

* `total_mbps` - Total Provisioned MBps Elastic SAN resource.

* `total_size_in_tib` - Total size of the Elastic SAN resource in TB.

* `total_volume_size_in_gib` - Total size of the provisioned Volumes in GiB.

* `volume_group_count` - Total number of volume groups in this Elastic SAN resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic SAN resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic SAN resource.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic SAN resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic SAN resource.

## Import

An existing Elastic SAN can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_san.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ElasticSan/elasticSans/esan1
```
