---
subcategory: "Network Function"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_function_azure_traffic_collector"
description: |-
  Manages a Network Function Azure Traffic Collectors.
---

# azurerm_network_function_azure_traffic_collector

Manages a Network Function Azure Traffic Collectors.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_function_azure_traffic_collector" "example" {
  name                = "example-nfatc"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  virtual_hub {

  }
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Function Azure Traffic Collectors. Changing this forces a new Network Function Azure Traffic Collectors to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Network Function Azure Traffic Collectors should exist. Changing this forces a new Network Function Azure Traffic Collectors to be created.

* `location` - (Required) Specifies the Azure Region where the Network Function Azure Traffic Collectors should exist.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Function Azure Traffic Collectors.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Function Azure Traffic Collectors.

* `collector_policies` - A `collector_policies` block as defined below.

* `system_data` - A `system_data` block as defined below.

* `virtual_hub` - A `virtual_hub` block as defined below.

---

A `collector_policies` block exports the following:

* `id` - Resource ID.

---

A `system_data` block exports the following:

* `created_at` - The timestamp of resource creation (UTC).

* `created_by` - The identity that created the resource.

* `created_by_type` - The type of identity that created the resource.

* `last_modified_by` - The identity that last modified the resource.

* `last_modified_by_type` - The type of identity that last modified the resource.

---

A `virtual_hub` block exports the following:

* `id` - Resource ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Function Azure Traffic Collectors.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Function Azure Traffic Collectors.
* `update` - (Defaults to 30 minutes) Used when updating the Network Function Azure Traffic Collectors.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Function Azure Traffic Collectors.

## Import

Network Function Azure Traffic Collectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_function_azure_traffic_collector.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.NetworkFunction/azureTrafficCollectors/azureTrafficCollector1
```
