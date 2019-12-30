---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
sidebar_current: "docs-azurerm-datasource-nat-gateway"
description: |-
  Gets information about an existing NAT Gateway
---

# Data Source: azurerm_nat_gateway

Use this data source to access information about an existing NAT Gateway.

-> **NOTE:** The Azure NAT Gateway service is currently in private preview. Your subscription must be on the NAT Gateway private preview whitelist for this resource to be provisioned correctly. If you attempt to provision this resource and receive an `InvalidResourceType` error may mean that your subscription is not part of the NAT Gateway private preview or you are using a region which does not yet support the NAT Gateway private preview service. The NAT Gateway private preview service is currently available in a limited set of regions. Private preview resources may have multiple breaking changes over their lifecycle until they GA. You can opt into the Private Preview by contacting your Microsoft Representative.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the NAT Gateway.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the NAT Gateway exists.

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
