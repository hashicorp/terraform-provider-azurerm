---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_san_volume"
description: |-
  Manages an Elastic SAN Volume resource.
---

# azurerm_elastic_san_volume

Manages an Elastic SAN Volume resource.

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

output "target_iqn" {
  value = azurerm_elastic_san_volume.example.target_iqn
}
```

## Example of creating an Elastic SAN Volume from a Disk Snapshot
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

resource "azurerm_managed_disk" "example" {
  name                 = "example-disk"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  create_option        = "Empty"
  storage_account_type = "Standard_LRS"
  disk_size_gb         = 2
}

resource "azurerm_snapshot" "example" {
  name                = "example-ss"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.example.id
}

resource "azurerm_elastic_san_volume" "example2" {
  name            = "example-esv2"
  volume_group_id = azurerm_elastic_san_volume_group.example.id
  size_in_gib     = 2
  create_source {
    source_type = "DiskSnapshot"
    source_id   = azurerm_snapshot.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Elastic SAN Volume. Changing this forces a new resource to be created.

* `volume_group_id` - (Required) Specifies the Volume Group ID within which this Elastic SAN Volume should exist. Changing this forces a new resource to be created.

* `size_in_gib` - (Required) Specifies the size of the Elastic SAN Volume in GiB. The size should be within the remaining capacity of the parent Elastic SAN. Possible values are between `1` and `65536` (16 TiB).

-> **Note:** The size can only be increased. If `create_source` is specified, then the size must be equal to or greater than the source's size.

* `create_source` - (Optional) A `create_source` block as defined below.

---

A `create_source` block supports the following:

* `source_id` - (Required) Specifies the ID of the source to create the Elastic SAN Volume from. Changing this forces a new resource to be created.

* `source_type` - (Required) Specifies the type of the source to create the Elastic SAN Volume from. Possible values are `Disk`, `DiskRestorePoint`, `DiskSnapshot` and `VolumeSnapshot`. Changing this forces a new resource to be created.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elastic SAN Volume.

* `target_iqn` - The iSCSI Target IQN of the Elastic SAN Volume.

* `target_portal_hostname` - The iSCSI Target Portal Host Name of the Elastic SAN Volume.

* `target_portal_port` - The iSCSI Target Portal Port of the Elastic SAN Volume.

* `volume_id` - The UUID of the Elastic SAN Volume.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic SAN Volume.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic SAN Volume.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic SAN Volume.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic SAN Volume.

## Import

An existing Elastic SAN Volume can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_san_volume.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ElasticSan/elasticSans/esan1/volumeGroups/vg1/volumes/vol1
```

