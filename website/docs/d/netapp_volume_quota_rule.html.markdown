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
  name = "existing"
  resource_group_name = "existing"
  volume_name = "existing"
  account_name = "existing"
  pool_name = "existing"
}

output "id" {
  value = data.azurerm_netapp_volume_quota_rule.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `account_name` - (Required) The name of the NetApp Account where the volume is located.

* `name` - (Required) The name of this Volume Quota Rule.

* `pool_name` - (Required) The name of the NetApp pool in which the NetApp Volume belongs to.

* `resource_group_name` - (Required) The name of the Resource Group where the Volume Quota Rule exists. Changing this forces a new Volume Quota Rule to be created.

* `volume_name` - (Required) The name of the NetApp Volume where the quota is assigned to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Volume Quota Rule.

* `location` - The Azure Region where the Volume Quota Rule exists.

* `quota_size_in_kib` - Quota size in kibibytes.

* `quota_target` -Quota Target.

* `quota_type` - Quota type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Volume Quota Rule.