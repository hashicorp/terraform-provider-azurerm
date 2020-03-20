---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_partner"
description: |-
  Manages a Management Partner.
---

# azurerm_management_partner

Manages a Management Partner.

## Management Partner Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_management_partner" "example" {
  partner_id = 6080810
}
```

## Argument Reference

The following arguments are supported:

* `partner_id` - (Required) Specifies the ID of the Management Partner.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Partner.

* `partner_name` - The friendly name for the Management Partner.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Management Partner.
* `update` - (Defaults to 30 minutes) Used when updating the Management Partner.
* `read` - (Defaults to 5 minutes) Used when retrieving the Management Partner.
* `delete` - (Defaults to 30 minutes) Used when deleting the Management Partner.

## Import

Management Partner can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_management_partner.example /providers/Microsoft.ManagementPartner/partners/6080810
```
