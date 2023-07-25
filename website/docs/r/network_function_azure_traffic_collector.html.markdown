---
subcategory: "Network Function"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_function_azure_traffic_collector"
description: |-
  Manages a Network Function Azure Traffic Collector.
---

# azurerm_network_function_azure_traffic_collector

Manages a Network Function Azure Traffic Collector.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_network_function_azure_traffic_collector" "example" {
  name                = "example-nfatc"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West US"

  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Function Azure Traffic Collector. Changing this forces a new Network Function Azure Traffic Collector to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Network Function Azure Traffic Collector should exist. Changing this forces a new Network Function Azure Traffic Collector to be created.

* `location` - (Required) Specifies the Azure Region where the Network Function Azure Traffic Collector should exist. Changing this forces a new Network Function Azure Traffic Collector to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Function Azure Traffic Collector.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Function Azure Traffic Collector.

* `collector_policy_ids` - The list of Resource IDs of collector policies.

* `virtual_hub_id` - The Resource ID of virtual hub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Function Azure Traffic Collector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Function Azure Traffic Collector.
* `update` - (Defaults to 30 minutes) Used when updating the Network Function Azure Traffic Collector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Function Azure Traffic Collector.

## Import

Network Function Azure Traffic Collector can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_function_azure_traffic_collector.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.NetworkFunction/azureTrafficCollectors/azureTrafficCollector1
```
