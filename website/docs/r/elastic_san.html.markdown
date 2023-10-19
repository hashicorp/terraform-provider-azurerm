---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_san"
description: |-
  Manages a Azure Elastic SAN (Storage Area Network).
---

# azurerm_elastic_san

Manages an Elastic SAN instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_elastic_san" "example" {
  name                          = "example"
  resource_group_name           = azurerm_resource_group.example.name
  location                      = azurerm_resource_group.example.location
  base_size_in_tib              = 1
  extended_capacity_size_in_tib = 2
  sku {
    name = "example-value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Elastic SAN. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Elastic SAN should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Elastic SAN should exist. Changing this forces a new resource to be created.

* `base_size_in_tib` - (Required) Base size of the Elastic SAN appliance in TiB. Possible value is integer between `1` and `100`.

* `extended_capacity_size_in_tib` - (Required) Extended size of the Elastic SAN appliance in TiB. Possible value is integer between `1` and `100`.

* `sku` - (Required) A `sku` block as defined below.

* `zones` - (Optional) Logical zone for the Elastic SAN resource. Changing this forces a new Elastic SAN to be created.

-> **NOTE** `zones` cannot be specified if `sku.name` is set to `Premium_ZRS`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Elastic SAN.

---

The `sku` block supports the following arguments:

* `name` - (Required) The SKU name. Possible values are `Premium_LRS` and `Premium_ZRS`.

-> **NOTE** `Premium_ZRS` SKU is currently only available in several Azure regions, please see the [link](https://azure.microsoft.com/en-us/updates/regional-expansion-azure-elastic-san-public-preview-is-now-available-in-more-regions/) for more information.

* `tier` - (Optional) The SKU tier. The only possible value is `Premium`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elastic SAN.

* `total_iops` - Total Provisioned IOps of the Elastic SAN appliance.

* `total_mbps` - Total Provisioned MBps Elastic SAN appliance.

* `total_size_in_tib` - Total size of the Elastic SAN appliance in TB.

* `total_volume_size_in_gib` - Total size of the provisioned Volumes in GiB.

* `volume_group_count` - Total number of volume groups in this Elastic SAN appliance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Elastic SAN.
* `delete` - (Defaults to 30 minutes) Used when deleting this Elastic SAN.
* `read` - (Defaults to 5 minutes) Used when retrieving this Elastic SAN.
* `update` - (Defaults to 30 minutes) Used when updating this Elastic SAN.

## Import

An existing Elastic SAN can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_san.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ElasticSan/elasticSans/esan1
```