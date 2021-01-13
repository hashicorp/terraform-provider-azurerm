---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_center_auto_provisioning"
description: |-
  Manages the subscription's Security Center Auto Provisioning.
---

# azurerm_security_center_auto_provisioning

Enables or disables the Security Center Auto Provisioning feature for the subscription

~> **NOTE:** There is no resource name required, it will always be "default"

## Example Usage

```hcl
resource "azurerm_security_center_auto_provisioning" "example" {
  auto_provision = "On"
}
```

## Arguments Reference

The following arguments are supported:

* `auto_provision` - (Required) Should the security agent be automatically provisioned on Virtual Machines in this subscription? Possible values are `On` (to install the security agent automatically, if it's missing) or `Off` (to not install the security agent automatically).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Security Center Auto Provisioning.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Security Center Auto Provisioning.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Center Auto Provisioning.
* `update` - (Defaults to 1 hour) Used when updating the Security Center Auto Provisioning.
* `delete` - (Defaults to 1 hour) Used when deleting the Security Center Auto Provisioning.

## Import

Security Center Auto Provisioning can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_center_auto_provisioning.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Security/autoProvisioningSettings/default
```
