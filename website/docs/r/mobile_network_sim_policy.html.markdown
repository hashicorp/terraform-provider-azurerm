---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_sim_policy"
description: |-
  Manages a Mobile Network Sim Policy.
---

# azurerm_mobile_network_sim_policy

Manages a Mobile Network Sim Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_data_network" "example" {
  name              = "example-mndn"
  mobile_network_id = azurerm_mobile_network.example.id
  location          = azurerm_resource_group.example.location
}

resource "azurerm_mobile_network_service" "example" {
  name               = "example-mns"
  mobile_network_id  = azurerm_mobile_network.example.id
  location           = azurerm_resource_group.example.location
  service_precedence = 0

  pcc_rule {
    name                    = "default-rule"
    precedence              = 1
    traffic_control_enabled = true

    service_data_flow_template {
      direction      = "Uplink"
      name           = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }
  }
}

resource "azurerm_mobile_network_slice" "example" {
  name              = "example-mns"
  mobile_network_id = azurerm_mobile_network.example.id
  location          = azurerm_resource_group.example.location
  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }
}

resource "azurerm_mobile_network_sim_policy" "example" {
  name                          = "example-mnsp"
  mobile_network_id             = azurerm_mobile_network.example.id
  location                      = azurerm_resource_group.example.location
  registration_timer_in_seconds = 3240
  default_slice_id              = azurerm_mobile_network_slice.example.id

  slice {
    default_data_network_id = azurerm_mobile_network_data_network.example.id
    slice_id                = azurerm_mobile_network_slice.example.id
    data_network {
      data_network_id                         = azurerm_mobile_network_data_network.example.id
      allocation_and_retention_priority_level = 9
      default_session_type                    = "IPv4"
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"
      allowed_services_ids                    = [azurerm_mobile_network_service.example.id]
      session_aggregate_maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
  }

  user_equipment_aggregate_maximum_bit_rate {
    downlink = "1 Gbps"
    uplink   = "500 Mbps"
  }

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Mobile Network Sim Policies. Changing this forces a new Mobile Network Sim Policies to be created.

* `mobile_network_id` - (Required) The ID of the Mobile Network which the Sim Policy belongs to. Changing this forces a new Mobile Network Sim Policies to be created.

* `default_slice_id` - (Required) The ID of default slice to use if the UE does not explicitly specify it. This slice must exist in the `slice` block.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Sim Policy should exist. Changing this forces a new Mobile Network Sim Policies to be created.

* `registration_timer_in_seconds` - (Optional) Interval for the user equipment periodic registration update procedure. Defaults to `3240`.

* `rat_frequency_selection_priority_index` - (Optional) RAT/Frequency Selection Priority Index, defined in 3GPP TS 36.413.

* `user_equipment_aggregate_maximum_bit_rate` - (Required) A `user_equipment_aggregate_maximum_bit_rate` block as defined below.

* `slice` - (Required) An array of `slice` block as defined below. The allowed slices and the settings to use for them. The list must not contain duplicate items and must contain at least one item.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Sim Policies.

---

A `slice` block supports the following:

* `data_network` - (Required) An array of `data_network` block as defined below.

* `default_data_network_id` - (Required) The ID of default data network to use if the user equipment does not explicitly specify it. Configuration for this object must exist in the `data_network` block.

* `slice_id` - (Required) The ID of the slice that these settings apply to.

---

A `data_network` block supports the following:

* `allowed_services_ids` - (Required) An array of IDs of services that can be used as part of this SIM policy. The array must not contain duplicate items and must contain at least one item.

* `data_network_id` - (Required) The ID of Mobile Network Data Network which these settings apply to.

* `session_aggregate_maximum_bit_rate` - (Required) A `session_aggregate_maximum_bit_rate` block as defined below.

* `allocation_and_retention_priority_level` - (Optional) Default QoS Flow allocation and retention priority (ARP) level. Flows with higher priority preempt flows with lower priority, if the settings of `preemption_capability` and `preemption_vulnerability` allow it. `1` is the highest level of priority. If this field is not specified then `qos_indicator` is used to derive the ARP value. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `additional_allowed_session_types` - (Optional) Allowed session types in addition to the default session type. Must not duplicate the default session type. Possible values are `IPv4` and `IPv6`.

* `default_session_type` - (Optional) The default PDU session type, which is used if the user equipment does not request a specific session type. Possible values are `IPv4` and `IPv6`. Defaults to `IPv4`.

* `qos_indicator` - (Required) The QoS Indicator (5QI for 5G network /QCI for 4G net work) value identifies a set of QoS characteristics, it controls QoS forwarding treatment for QoS flows or EPS bearers. Recommended values: 5-9; 69-70; 79-80. Must be between `1` and `127`.

* `max_buffered_packets` - (Optional) The maximum number of downlink packets to buffer at the user plane for High Latency Communication - Extended Buffering. Defaults to `10`, Must be at least `0`, See 3GPP TS29.272 v15.10.0 section 7.3.188 for a full description. This maximum is not guaranteed because there is a internal limit on buffered packets across all PDU sessions.

* `preemption_capability` - (Optional) The Preemption Capability of a QoS Flow, it controls whether it can preempt another QoS Flow with a lower priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters. Possible values are `NotPreempt` and `MayPreempt`, Defaults to `NotPreempt`.

* `preemption_vulnerability` - (Optional) The Preemption Vulnerability of a QoS Flow, it controls whether it can be preempted by QoS Flow with a higher priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters. Possible values are `NotPreemptable` and `Preemptable`. Defaults to `NotPreemptable`.

---

A `session_aggregate_maximum_bit_rate` block supports the following:

* `downlink` - (Required) Downlink bit rate. Must be a number followed by `Kbps`, `Mbps`, `Gbps` or `Tbps`.

* `uplink` - (Required) Uplink bit rate. Must be a number followed by `Kbps`, `Mbps`, `Gbps` or `Tbps`.

---

A `user_equipment_aggregate_maximum_bit_rate` block supports the following:

* `downlink` - (Required) Downlink bit rate. Must be a number followed by `Kbps`, `Mbps`, `Gbps` or `Tbps`.

* `uplink` - (Required) Uplink bit rate. Must be a number followed by `Kbps`, `Mbps`, `Gbps` or `Tbps`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Sim Policies.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Sim Policies.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Sim Policies.
* `update` - (Defaults to 1 hour) Used when updating the Mobile Network Sim Policies.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Sim Policies.

## Import

Mobile Network Sim Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_sim_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/mobileNetworks/mobileNetwork1/simPolicies/simPolicy1
```
