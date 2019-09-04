---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_protection_plan"
sidebar_current: "docs-azurerm-datasource-network-ddos-protection-plan-x"
description: |-
  Use this data source to access information about an existing Azure Network DDoS Protection Plan.

---

# Data Source: azurerm_network_ddos_protection_plan

Manages an AzureNetwork DDoS Protection Plan.

-> **NOTE** Azure only allow `one` DDoS Protection Plan per region.

## Example Usage

```hcl
data "azurerm_network_ddos_protection_plan" "example" {
  name                = "azurerm_network_ddos_protection_plan.test.name"
  resource_group_name = "${azurerm_network_ddos_protection_plan.test.resource_group_name}"
}

output "ddos_protection_plan_id" {
  value = "${data.azurerm_network_ddos_protection_plan.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Network DDoS Protection Plan.

* `location` - (Required) Specifies the supported Azure location where the resource exists.

* `resource_group_name` - (Required) The name of the resource group where the Network DDoS Protection Plan exists.

* `tags` - (Optional) A mapping of tags assigned to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the DDoS Protection Plan

* `virtual_network_ids` - The Resource ID list of the Virtual Networks associated with DDoS Protection Plan.
