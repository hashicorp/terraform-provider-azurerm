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
  location = "West US 2"
}

resource "azurerm_express_route_port" "example" {
  name                = "example-erp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  peering_location    = "Equinix-Seattle-SE2"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "example" {
  name                  = "example-erc"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  express_route_port_id = azurerm_express_route_port.example.id
  bandwidth_in_gbps     = 1

  sku {
    tier   = "Standard"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_circuit_peering" "example" {
  peering_type                  = "MicrosoftPeering"
  express_route_circuit_name    = azurerm_express_route_circuit.example.name
  resource_group_name           = azurerm_resource_group.example.name
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.199.0/30"
  secondary_peer_address_prefix = "192.168.200.0/30"
  vlan_id                       = 300

  microsoft_peering_config {
    advertised_public_prefixes = ["123.6.0.0/24"]
  }
}

resource "azurerm_network_function_azure_traffic_collector" "example" {
  name                = "example-nfatc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  depends_on = [
    azurerm_express_route_circuit_peering.example
  ]
}

resource "azurerm_network_function_collector_policy" "example" {
  name                 = "example-nfcp"
  traffic_collector_id = azurerm_network_function_azure_traffic_collector.example.id
  location             = azurerm_resource_group.example.location

  emission_policy {
    emission_type = "IPFIX"
    emission_destination {
      destination_type = "AzureMonitor"
    }
  }

  ingestion_policy {
    ingestion_type = "IPFIX"
    ingestion_source {
      resource_id = azurerm_express_route_circuit.example.id
      source_type = "Resource"
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

* `traffic_collector_id` - (Required) Specifies the Azure Traffic Collector ID of the Network Function Collector Policy. Changing this forces a new Network Function Collector Policy to be created.

* `location` - (Required) Specifies the Azure Region where the Network Function Collector Policy should exist. Changing this forces a new Network Function Collector Policy to be created.

* `emission_policy` - (Required) An `emission_policy` block as defined below. Changing this forces a new Network Function Collector Policy to be created.

* `ingestion_policy` - (Required) An `ingestion_policy` block as defined below. Changing this forces a new Network Function Collector Policy to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Network Function Collector Policy.

---

An `emission_policy` block supports the following:

* `emission_destination` - (Required) An `emission_destination` block as defined below. Changing this forces a new Network Function Collector Policy to be created.

* `emission_type` - (Required) Emission format type. The only possible value is `IPFIX`. Changing this forces a new Network Function Collector Policy to be created.

---

An `emission_destination` block supports the following:

* `destination_type` - (Required) Emission destination type. The only possible value is `AzureMonitor`. Changing this forces a new Network Function Collector Policy to be created.

---

An `ingestion_policy` block supports the following:

* `ingestion_source` - (Required) An `ingestion_source` block as defined below. Changing this forces a new Network Function Collector Policy to be created.

* `ingestion_type` - (Required) Specifies the ingestion type. The only possible value is `IPFIX`. Changing this forces a new Network Function Collector Policy to be created.

---

An `ingestion_source` block supports the following:

* `resource_id` - (Required) Resource ID.

* `source_type` - (Required) Ingestion source type. The only possible value is `Resource`.

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
