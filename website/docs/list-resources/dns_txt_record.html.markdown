---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_txt_record"
description: |-
    Lists DNS TXT Record resources.
---

# List resource: azurerm_dns_txt_record

Lists DNS TXT Record resources.

## Example Usage

### List all DNS TXT Records in a DNS Zone

```hcl
list "azurerm_dns_txt_record" "example" {
  provider = azurerm
  config {
    dns_zone_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `dns_zone_id` - (Required) The ID of the DNS Zone to query.
