---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_protection_plan"
description: |-
  Use this data source to access information about an existing Azure Network DDoS Protection Plan.

---

# Data Source: azurerm_network_ddos_protection_plan

Use this data source to access information about an existing Azure Network DDoS Protection Plan.

## Example Usage

```hcl
data "azurerm_network_ddos_protection_plan" "example" {
  name                = azurerm_network_ddos_protection_plan.example.name
  resource_group_name = azurerm_network_ddos_protection_plan.example.resource_group_name
}

output "ddos_protection_plan_id" {
  value = data.azurerm_network_ddos_protection_plan.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Network DDoS Protection Plan.

* `resource_group_name` - The name of the resource group where the Network DDoS Protection Plan exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DDoS Protection Plan

* `location` - Specifies the supported Azure location where the resource exists.

* `tags` - A mapping of tags assigned to the resource.

* `virtual_network_ids` - A list of ID's of the Virtual Networks associated with this DDoS Protection Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Protection Plan.
