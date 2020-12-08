---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
description: |-
  Gets information about an existing NAT Gateway
---

# Data Source: azurerm_nat_gateway

Use this data source to access information about an existing NAT Gateway.

## Argument Reference

The following arguments are supported:

* `name` - Specifies the Name of the NAT Gateway.

* `resource_group_name` - Specifies the name of the Resource Group where the NAT Gateway exists.

## Attributes Reference

The following attributes are exported:

* `location` - The location where the NAT Gateway exists.

* `idle_timeout_in_minutes` - The idle timeout in minutes which is used for the NAT Gateway.

* `public_ip_address_ids` - A list of existing Public IP Address resource IDs which the NAT Gateway is using.

* `public_ip_prefix_ids` - A list of existing Public IP Prefix resource IDs which the NAT Gateway is using.

* `resource_guid` - The Resource GUID of the NAT Gateway.

* `sku_name` - The SKU used by the NAT Gateway.

* `tags` - A mapping of tags assigned to the resource.

* `zones` - A list of Availability Zones which the NAT Gateway exists in.

~> **NOTE:** The field `public_ip_address_ids` has been deprecated in favour of `azurerm_nat_gateway_public_ip_association`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NAT Gateway.
