---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_partner"
sidebar_current: "docs-azurerm-resource-management-partner"
description: |-
  Manages a Management Partner.
---

# azurerm_management_partner

Manages a Management Partner.

## Management Partner Usage

```hcl
resource "azurerm_management_partner" "example" {
  partner_id = 512725
}
```

## Argument Reference

The following arguments are supported:

* `partner_id` - (Required) The ID of the Management Partner. Changing this forces a new resource to be created.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Partner.

* `partner_name` - A friendly name for the Management Partner.

## Import

Management Partner can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_management_partner.example /providers/Microsoft.ManagementPartner/partners/5127255
```
