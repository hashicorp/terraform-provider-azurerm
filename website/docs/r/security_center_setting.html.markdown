---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_setting"
description: |-
    Manages the Data Access Settings for Azure Security Center.
---

# azurerm_security_center_setting

Manages the Data Access Settings for Azure Security Center.

~> **Note:** This resource requires the `Owner` permission on the Subscription.

~> **Note:** Deletion of this resource disables the setting.

## Example Usage

```hcl
resource "azurerm_security_center_setting" "example" {
  setting_name = "MCAS"
  enabled      = true
}
```

## Argument Reference

The following arguments are supported:

* `setting_name` - (Required) The setting to manage. Possible values are `MCAS` , `WDATP`, `WDATP_EXCLUDE_LINUX_PUBLIC_PREVIEW`, `WDATP_UNIFIED_SOLUTION` and `Sentinel`. Changing this forces a new resource to be created.
* `enabled` - (Required) Boolean flag to enable/disable data access.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The subscription security center setting id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Security Center Setting.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Setting.
* `update` - (Defaults to 10 minutes) Used when updating the Security Center Setting.
* `delete` - (Defaults to 10 minutes) Used when deleting the Security Center Setting.

## Import

The setting can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_setting.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/settings/<setting_name>
```
