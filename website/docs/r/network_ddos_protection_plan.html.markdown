---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_protection_plan"
sidebar_current: "docs-azurerm-resource-network-ddos-protection-plan-x"
description: |-
  Manages an Azure Network DDoS Protection Plan.

---

# azurerm_network_ddos_protection_plan

Manages an AzureNetwork DDoS Protection Plan.

-> **NOTE** Azure only allows `one` DDoS Protection Plan per region.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "example-protection-plan"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Network DDoS Protection Plan. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the DDoS Protection Plan

* `virtual_network_ids` - The Resource ID list of the Virtual Networks associated with DDoS Protection Plan.

## Import

Azure DDoS Protection Plan can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_ddos_protection_plan.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/ddosProtectionPlans/testddospplan
```
