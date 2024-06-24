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
  volume_group_id = data.azurerm_elastic_san.example.id
}

output "id" {
  value = data.azurerm_elastic_san_volume_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Elastic SAN Volume Snapshot.

* `volume_group_id` - The Elastic SAN Volume Group ID within which the Elastic SAN Volume Snapshot exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic SAN.

* `creation_source` - A `creation_source` block as defined below.

---

A `creation_source` block exports the following arguments:

* `source_id` - The Resource ID from which the Snapshot is created.
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic SAN Volume Snapshot.
