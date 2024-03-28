---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_san_snapshot"
description: |-
  Manages a Elastic SAN Snapshot.
---

# azurerm_elastic_san_snapshot

manages an Elastic SAN Snapshot.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_elastic_san" "example" {
  name                = "example-es"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  base_size_in_tib    = 1
  sku {
    name = "Premium_LRS"
  }
}

resource "azurerm_elastic_san_volume_group" "example" {
  name           = "example-esvg"
  elastic_san_id = azurerm_elastic_san.example.id
}

resource "azurerm_elastic_san_volume" "example" {
  name            = "example-esv"
  volume_group_id = azurerm_elastic_san_volume_group.example.id
  size_in_gib     = 1
}

resource "azurerm_elastic_san_snapshot" "test" {
  name            = "acctestess-${var.random_string}"
  volume_group_id = azurerm_elastic_san_volume_group.test.id
  creation_source {
    source_id = azurerm_elastic_san_volume.test.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Elastic SAN Snapshot. Changing this forces a new Elastic SAN Snapshot to be created.

* `volume_group_id` - (Required) Specifies the Volume Group Id within which this Elastic SAN Snapshot should exist. Changing this forces a new Elastic SAN Snapshot to be created.

* `creation_source` - (Required) A `creation_source` block as defined below. Changing this forces a new Elastic SAN Snapshot to be created.

---

A `creation_source` block supports the following arguments:

* `source_id` - (Required) Specifies the Resource ID to create a Snapshot from. Current this must be ID of an Elastic SAN Volume. Changing this forces a new Elastic SAN Snapshot to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elastic SAN Snapshot.

* `source_volume_size_gib` - Size of Source Volume.

* `volume_name` - Source Volume Name of a snapshot.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this Elastic SAN Snapshot.
* `delete` - (Defaults to 30 minutes) Used when deleting this Elastic SAN Snapshot.
* `read` - (Defaults to 5 minutes) Used when retrieving this Elastic SAN Snapshot.

## Import

An existing Elastic SAN Snapshot can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_san_snapshot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ElasticSan/elasticSans/esan1/volumeGroups/vg1/snapshots/sp1
```
