---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_partner"
sidebar_current: "docs-azurerm-datasource-management-partner"
description: |-
Gets information about an existing Management Partner
---

# Data Source: azurerm_management_partner

Uses this data source to access information about an existing Management Partner.

## Management Partner Usage

```hcl
data "azurerm_management_partner" "example" {
  partner_id = 6080810
}

output "management_partner_id" {
  value = "${data.azurerm_management_partner.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `partner_id` - (Required) Specifies the ID of this Management Partner.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Management Partner.

* `partner_name` - The friendly name for the Management Partner.
