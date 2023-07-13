---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_sim_policy"
description: |-
  Get information about a Mobile Network Sim Policy.
---

# azurerm_mobile_network_sim_policy

Get information about a Mobile Network Sim Policy.

## Example Usage

```hcl
data "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = "example-rg"
}

data "azurerm_mobile_network_sim_policy" "example" {
  name              = "example-mnsp"
  mobile_network_id = data.azurerm_mobile_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name which should be used for this Mobile Network Sim Policies.

* `mobile_network_id` - The ID of the Mobile Network which the Sim Policy belongs to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Sim Policies.

* `default_slice_id` - The ID of default slice to use if the UE does not explicitly specify it.

* `location` - The Azure Region where the Mobile Network Sim Policy should exist.

* `registration_timer_in_seconds` - Interval for the UE periodic registration update procedure.

* `rat_frequency_selection_priority_index` - RAT/Frequency Selection Priority Index, defined in 3GPP TS 36.413.

* `user_equipment_aggregate_maximum_bit_rate` - A `user_equipment_aggregate_maximum_bit_rate` block as defined below.

* `slice` - An array of `slice` block as defined below. The allowed slices and the settings to use for them.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Sim Policies.

---

A `slice` block supports the following:

* `data_network` - An array of `data_network` block as defined below.

* `default_data_network_id` - The ID of default data network to use if the UE does not explicitly specify it.

* `slice_id` - The ID of the slice that these settings apply to.

---

A `data_network` block supports the following:

* `allowed_services_ids` - An array of IDs of services that can be used as part of this SIM policy.

* `data_network_id` - The ID of Mobile Network Data Network which these settings apply to.

* `session_aggregate_maximum_bit_rate` - A `session_aggregate_maximum_bit_rate` block as defined below.

* `allocation_and_retention_priority_level` - Default QoS Flow allocation and retention priority (ARP) level. Flows with higher priority preempt flows with lower priority, if the settings of `preemption_capability` and `preemption_vulnerability` allow it. 1 is the highest level of priority. If this field is not specified then `qos_indicator` is used to derive the ARP value. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `additional_allowed_session_types` - Allowed session types in addition to the default session type.

* `default_session_type` - The default PDU session type, which is used if the UE does not request a specific session type.

* `qos_indicator` - The QoS Indicator (5QI for 5G network /QCI for 4G net work) value identifies a set of QoS characteristics that control QoS forwarding treatment for QoS flows or EPS bearers.

* `max_buffered_packets` - The maximum number of downlink packets to buffer at the user plane for High Latency Communication - Extended Buffering.

* `preemption_capability` - The Preemption Capability of a QoS Flow controls whether it can preempt another QoS Flow with a lower priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `preemption_vulnerability` - The Preemption Vulnerability of a QoS Flow controls whether it can be preempted by QoS Flow with a higher priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

---

A `session_aggregate_maximum_bit_rate` block supports the following:

* `downlink` - Downlink bit rate.

* `uplink` - Uplink bit rate.

---

A `user_equipment_aggregate_maximum_bit_rate` block supports the following:

* `downlink` - Downlink bit rate.

* `uplink` - Uplink bit rate.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Sim Policies.
