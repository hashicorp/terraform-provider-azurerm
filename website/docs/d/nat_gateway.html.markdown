---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
sidebar_current: "docs-azurerm-datasource-nat-gateway"
description: |-
  Gets information about an existing Nat Gateway
---

# Data Source: azurerm_nat_gateway

Use this data source to access information about an existing Nat Gateway.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of the Resource Group where the Nat Gateway exists.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Nat Gateway exists.

## Attributes Reference

The following attributes are exported:

* `location` - The location where the NAT Gateway exists.

* `idle_timeout_in_minutes` - The idle timeout of the Nat Gateway.

* `public_ip_address_ids` - A list of existing Public IP Address resource IDs which the NAT Gateway is using.

* `public_ip_prefix_ids` - A list of existing Public IP Prefix resource IDs which the NAT Gateway is using.

* `resource_guid` - The resource GUID property of the Nat Gateway.

* `sku_name` - The SKU used by the NAT Gateway.

* `subnet_ids` - A list of existing Subnet resource IDs which the NAT Gateway is using.

* `zones` - A list of Availability Zones which the NAT Gateway gets created in.

* `tags` - A mapping of tags assigned to the resource.
