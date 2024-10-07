---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_netapp_volume_quota_rule"
description: |-
  Gets information about an existing Volume Quota Rule.
---

# Data Source: azurerm_netapp_volume_quota_rule

Use this data source to access information about an existing Volume Quota Rule.

## Example Usage

```hcl
data "azurerm_netapp_volume_quota_rule" "example" {
  name      = "exampleQuotaRule"
  volume_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/vol1"
}

output "id" {
  value = data.azurerm_netapp_volume_quota_rule.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Volume Quota Rule.

* `volume_id` - (Required) The NetApp volume ID where the Volume Quota Rule is assigned to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Volume Quota Rule.

* `location` - The Azure Region where the Volume Quota Rule exists.

* `quota_size_in_kib` - The quota size in kibibytes.

* `quota_target` - The quota Target.

* `quota_type` - The quota type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Volume Quota Rule.
