---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_elastic_san_volume_snapshot"
description: |-
  Gets information about an existing Elastic SAN Volume Snapshot.
---

# Data Source: azurerm_elastic_san_volume_snapshot

Use this data source to access information about an existing Elastic SAN Volume Snapshot.

## Example Usage

```hcl
data "azurerm_elastic_san" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

data "azurerm_elastic_san_volume_group" "example" {
  name           = "existing"
  elastic_san_id = data.azurerm_elastic_san.example.id
}

data "azurerm_elastic_san_volume_snapshot" "example" {
  name            = "existing"
  volume_group_id = data.azurerm_elastic_san_volume_group.example.id
}

output "id" {
  value = data.azurerm_elastic_san_volume_snapshot.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Elastic SAN Volume Snapshot.

* `volume_group_id` - The Elastic SAN Volume Group ID within which the Elastic SAN Volume Snapshot exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic SAN Volume Snapshot.

* `source_id` - The resource ID from which the Snapshot is created.

* `source_volume_size_in_gib` - The size of source volume.

* `volume_name` - The source volume name of the Snapshot.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic SAN Volume Snapshot.
