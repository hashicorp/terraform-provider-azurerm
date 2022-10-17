---
subcategory: "Network Function"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_function_collector_policy"
description: |-
  Manages a Network Function Collector Policy.
---

# azurerm_network_function_collector_policy

Manages a Network Function Collector Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_function_azure_traffic_collector" "example" {
  name                = "example-nfatc"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_function_collector_policy" "example" {
  name                                        = "example-nfcp"
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "West Europe"
  ingestion_policy {
    ingestion_sources {
      resource_id = ""
    }
  }
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Function Collector Policy. Changing this forces a new Network Function Collector Policy to be created.

* `network_function_azure_traffic_collector_id` - (Required) Specifies the ID of the Network Function Collector Policy. Changing this forces a new Network Function Collector Policy to be created.

* `location` - (Required) Specifies the Azure Region where the Network Function Collector Policy should exist. Changing this forces a new Network Function Collector Policy to be created.

* `ingestion_policy` - (Optional) An `ingestion_policy` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Function Collector Policy.

---

An `ingestion_policy` block supports the following:

* `ingestion_sources` - (Optional) An `ingestion_sources` block as defined below.

---

An `ingestion_sources` block supports the following:

* `resource_id` - (Optional) Resource ID.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Function Collector Policy.

* `emission_policies` - An `emission_policies` block as defined below.

* `ingestion_policy` - An `ingestion_policy` block as defined below.

* `system_data` - A `system_data` block as defined below.

---

An `emission_policies` block exports the following:

* `emission_destinations` - An `emission_destinations` block as defined below.

* `emission_type` - Emission format type.

---

An `emission_destinations` block exports the following:

* `destination_type` - Emission destination type.

---

An `ingestion_policy` block exports the following:

* `ingestion_sources` - An `ingestion_sources` block as defined below.

* `ingestion_type` - The ingestion type.

---

An `ingestion_sources` block exports the following:

* `source_type` - Ingestion source type.

---

A `system_data` block exports the following:

* `created_at` - The timestamp of resource creation (UTC).

* `created_by` - The identity that created the resource.

* `created_by_type` - The type of identity that created the resource.

* `last_modified_by` - The identity that last modified the resource.

* `last_modified_by_type` - The type of identity that last modified the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Function Collector Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Function Collector Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Network Function Collector Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Function Collector Policy.

## Import

Network Function Collector Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_function_collector_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.NetworkFunction/azureTrafficCollectors/azureTrafficCollector1/collectorPolicies/collectorPolicy1
```
