---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_quota_rule"
description: |-
  Manages a Volume Quota Rule.
---

# azurerm_netapp_volume_quota_rule

Manages a Volume Quota Rule.

## Example Usage

```hcl
resource "azurerm_netapp_volume_quota_rule" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
  quota_size_in_kib   = 42
  account_name        = "example"
  pool_name           = "example"
  volume_name         = "example"
  quota_type          = "IndividualUserQuota"
}
```

## Arguments Reference

The following arguments are supported:

* `account_name` - (Required) The name of the NetApp Account where the volume is located. Changing this forces a new Volume Quota Rule to be created.

* `location` - (Required) The Azure Region where the Volume Quota Rule should exist. Changing this forces a new Volume Quota Rule to be created.

* `name` - (Required) The name which should be used for this Volume Quota Rule. Changing this forces a new Volume Quota Rule to be created.

* `pool_name` - (Required) The name of the NetApp pool in which the NetApp Volume belongs to. Changing this forces a new Volume Quota Rule to be created.

* `quota_size_in_kib` - (Required) Quota size in kibibytes.

* `quota_type` - (Required) Quota type. Possible values are `DefaultGroupQuota`, `DefaultUserQuota`, `IndividualGroupQuota` and `IndividualUserQuota`. Please note that `IndividualGroupQuota` and `DefaultGroupQuota` are not applicable to SMB and dual-protocol volumes.

* `resource_group_name` - (Required) The name of the Resource Group where the Volume Quota Rule should exist. Changing this forces a new Volume Quota Rule to be created.

* `volume_name` - (Required) The name of the NetApp Volume where the quota will be assigned to. Changing this forces a new Volume Quota Rule to be created.

* `quota_target` - (Optional) Quota Target. This can be Unix UID/GID for NFSv3/NFSv4.1 volumes and Windows User SID for CIFS based volumes. Use user id quota targets for `InidividualUserQuota` and group id quota targets for `IndividualGroupQuota` quota types.

---

~> **NOTE:** more information about this resource and Azure NetApp Files feature can be found at [Understand default and individual user and group quotas](https://learn.microsoft.com/en-us/azure/azure-netapp-files/default-individual-user-group-quotas-introduction)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Volume Quota Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour and 30 minutes) Used when creating the Volume Quota Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Volume Quota Rule.
* `update` - (Defaults to 2 hours) Used when updating the Volume Quota Rule.
* `delete` - (Defaults to 2 hours) Used when deleting the Volume Quota Rule.

## Import

Volume Quota Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_volume_quota_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/vol1/volumeQuotaRules/quota1
```
