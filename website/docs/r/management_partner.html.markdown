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

## Import

Management Partner can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_management_partner.example /providers/Microsoft.ManagementPartner/partners/6080810
```
