---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway_public_ip_association"
description: |-
  Lists Public IP Address Associations on a NAT Gateway.
---

# List resource: azurerm_nat_gateway_public_ip_association

Lists Public IP Address Associations on a NAT Gateway.

## Example Usage

### List all Public IP Address Associations in the Subscription

```hcl
list "azurerm_nat_gateway_public_ip_association" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Public IP Address Associations in a Resource Group

```hcl
list "azurerm_nat_gateway_public_ip_association" "example" {
  provider = azurerm
  config {
    resource_group_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1"
  }
}
```

### List all Public IP Address Associations on a Specific NAT Gateway

```hcl
list "azurerm_nat_gateway_public_ip_association" "example" {
  provider = azurerm
  config {
    nat_gateway_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/natGateways/natgateway1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_id` - (Optional) The ID of the Resource Group to query NAT Gateway Public IP Address Associations for. Conflicts with `nat_gateway_id`.

* `nat_gateway_id` - (Optional) The ID of the NAT Gateway to query Public IP Address Associations for. Conflicts with `resource_group_id`.
