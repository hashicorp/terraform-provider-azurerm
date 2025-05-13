---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_ddos_protection_plan"
description: |-
  Manages an Azure Network DDoS Protection Plan.

---

# azurerm_network_ddos_protection_plan

Manages an AzureNetwork DDoS Protection Plan.

-> **Note:** Azure only allows `one` DDoS Protection Plan per region.

~> **Note:** The DDoS Protection Plan is a [high-cost service](https://azure.microsoft.com/en-us/pricing/details/ddos-protection/#pricing). Please keep this in mind while testing and learning.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_ddos_protection_plan" "example" {
  name                = "example-protection-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Network DDoS Protection Plan. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the DDoS Protection Plan

* `virtual_network_ids` - A list of Virtual Network IDs associated with the DDoS Protection Plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DDoS Protection Plan.
* `read` - (Defaults to 5 minutes) Used when retrieving the DDoS Protection Plan.
* `update` - (Defaults to 30 minutes) Used when updating the DDoS Protection Plan.
* `delete` - (Defaults to 30 minutes) Used when deleting the DDoS Protection Plan.

## Import

Azure DDoS Protection Plan can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_ddos_protection_plan.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/ddosProtectionPlans/testddospplan
```
