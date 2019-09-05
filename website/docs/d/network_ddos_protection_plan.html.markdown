---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_protection_plan"
sidebar_current: "docs-azurerm-datasource-network-ddos-protection-plan-x"
description: |-
  Use this data source to access information about an existing Azure Network DDoS Protection Plan.

---

# Data Source: azurerm_network_ddos_protection_plan

Use this data source to access information about an existing Azure Network DDoS Protection Plan.

## Example Usage

```hcl
data "azurerm_network_ddos_protection_plan" "example" {
  name                = "${azurerm_network_ddos_protection_plan.example.name}"
  resource_group_name = "${azurerm_network_ddos_protection_plan.example.resource_group_name}"
}

output "ddos_protection_plan_id" {
  value = "${data.azurerm_network_ddos_protection_plan.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Network DDoS Protection Plan.

* `resource_group_name` - (Required) The name of the resource group where the Network DDoS Protection Plan exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the DDoS Protection Plan

* `location` - Specifies the supported Azure location where the resource exists.

* `tags` - A mapping of tags assigned to the resource.

* `virtual_network_ids` - The Resource ID list of the Virtual Networks associated with DDoS Protection Plan.
