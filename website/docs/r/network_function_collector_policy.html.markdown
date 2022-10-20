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
  emission_policies {
    emission_type = ""
    emission_destinations {
      destination_type = ""
    }
  }
  ingestion_policy {
    ingestion_type = ""
    ingestion_sources {
      resource_id = ""
      source_type = ""
    }
  }
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Network Function Collector Policy. It can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 80 characters, and it must begin with a letter or number and end with a letter, number or underscore. Changing this forces a new Network Function Collector Policy to be created.

* `network_function_azure_traffic_collector_id` - (Required) Specifies the ID of the Network Function Collector Policy. Changing this forces a new Network Function Collector Policy to be created.

* `location` - (Required) Specifies the Azure Region where the Network Function Collector Policy should exist. Changing this forces a new Network Function Collector Policy to be created.

* `emission_policies` - (Optional) An `emission_policies` block as defined below.

* `ingestion_policy` - (Optional) An `ingestion_policy` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Function Collector Policy.

---

An `emission_policies` block supports the following:

* `emission_destinations` - (Optional) An `emission_destinations` block as defined below.

* `emission_type` - (Optional) Emission format type.

---

An `emission_destinations` block supports the following:

* `destination_type` - (Optional) Emission destination type.

---

An `ingestion_policy` block supports the following:

* `ingestion_sources` - (Optional) An `ingestion_sources` block as defined below.

* `ingestion_type` - (Optional) Specifies the ingestion type.

---

An `ingestion_sources` block supports the following:

* `resource_id` - (Optional) Resource ID.

* `source_type` - (Optional) Ingestion source type.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Function Collector Policy.

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
