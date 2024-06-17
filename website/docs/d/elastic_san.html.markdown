---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_elastic_san"
description: |-
  Gets information about an existing Elastic SAN.
---

# Data Source: azurerm_elastic_san

Use this data source to access information about an existing Elastic SAN.

## Example Usage

```hcl
data "azurerm_elastic_san" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_elastic_san.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Elastic SAN.

* `resource_group_name` - (Required) The name of the Resource Group where the Elastic SAN exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic SAN.

* `base_size_in_tib` - The base size of the Elastic SAN resource in TiB.

* `extended_size_in_tib` - The base size of the Elastic SAN resource in TiB.

* `location` - The Azure Region where the Elastic SAN exists.

* `sku` - A `sku` block as defined below.

* `tags` - A mapping of tags assigned to the Elastic SAN.

* `total_iops` - Total Provisioned IOps of the Elastic SAN resource.

* `total_mbps` - Total Provisioned MBps Elastic SAN resource.

* `total_size_in_tib` - Total size of the Elastic SAN resource in TB.

* `total_volume_size_in_gib` - Total size of the provisioned Volumes in GiB.

* `volume_group_count` - Total number of volume groups in this Elastic SAN resource.

* `zones` - Logical zone for the Elastic SAN resource.

---

A `sku` block exports the following:

* `name` - The SKU name.

* `tier` - The SKU tier.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic SAN.
