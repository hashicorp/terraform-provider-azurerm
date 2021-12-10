---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_setting"
description: |-
    Manages the Data Access Settings for Azure Security Center.
---

# azurerm_security_center_setting

Manages the Data Access Settings for Azure Security Center.

~> **NOTE:** This resource requires the `Owner` permission on the Subscription.

~> **NOTE:** Deletion of this resource does not change or reset the data access settings

## Example Usage

```hcl
resource "azurerm_security_center_setting" "example" {
  setting_name = "MCAS"
  enabled      = true
}
```

## Argument Reference

The following arguments are supported:

* `setting_name` - (Required) The setting to manage. Possible values are `MCAS` and `WDATP`.
* `enabled` - (Required) Boolean flag to enable/disable data access.

## Attributes Reference

The following attributes are exported:

* `id` - The subscription security center setting id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Security Center Setting.
* `update` - (Defaults to 60 minutes) Used when updating the Security Center Setting.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Setting.
* `delete` - (Defaults to 60 minutes) Used when deleting the Security Center Setting.

## Import

The setting can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_setting.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/settings/<setting_name>
```
