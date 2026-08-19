---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_a_record"
description: |-
    Lists Private Dns A Record resources.
---

# List resource: azurerm_private_dns_a_record

Lists Private Dns A Record resources.

## Example Usage

### List Private Dns A Records in a Private Dns Zone

```hcl
list "azurerm_private_dns_a_record" "example" {
  provider = azurerm
  config {
    private_dns_zone_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/zone1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `private_dns_zone_id` - (Required) The ID of the Private Dns Zone to query.
