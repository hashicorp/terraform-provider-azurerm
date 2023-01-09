---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_service"
description: |-
  Get information about a Mobile Network Service.
---

# azurerm_mobile_network_service

Get information about a Mobile Network Service.

## Example Usage

```hcl
data "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = "example-rg"
}

resource "azurerm_mobile_network_service" "example" {
  name              = "example-mns"
  mobile_network_id = data.azurerm_mobile_network.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Service. Changing this forces a new Mobile Network Service to be created.

* `mobile_network_id` - (Required) Specifies the ID of the Mobile Network Service. Changing this forces a new Mobile Network Service to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Service.

* `location` - The Azure Region where the Mobile Network Service should exist.

* `service_precedence` - A precedence value that is used to decide between services when identifying the QoS values to use for a particular SIM. A lower value means a higher priority. This value should be unique among all services configured in the mobile network.

* `pcc_rules` - A `pcc_rules` block as defined below. The set of PCC Rules that make up this service.

* `service_qos_policy` - A `service_qos_policy` block as defined below. The QoS policy to use for packets matching this service. This can be overridden for particular flows using the ruleQosPolicy field in a `pcc_rules`. If this field is not specified then the `sim_policy` of User Equipment (UE) will define the QoS settings.

* `tags` - A mapping of tags which should be assigned to the Mobile Network Service.

---

A `pcc_rules` block supports the following:

* `rule_name` - Specifies the name of the rule. This must be unique within the parent service.

* `rule_precedence` - A precedence value that is used to decide between data flow policy rules when identifying the QoS values to use for a particular SIM. A lower value means a higher priority. This value should be unique among all data flow policy rules configured in the mobile network.

* `rule_qos_policy` - A `rule_qos_policy` block as defined below. The QoS policy to use for packets matching this rule. If this field is not specified then the Service will define the QoS settings.

* `service_data_flow_templates` - A `service_data_flow_templates` block as defined below. The set of service data flow templates to use for this PCC rule.

* `traffic_control_enabled` - Determines whether flows that match this data flow policy rule are permitted.

---

A `rule_qos_policy` block supports the following:

* `allocation_and_retention_priority_level` - QoS Flow allocation and retention priority (ARP) level. Flows with higher priority preempt flows with lower priority, if the settings of `preemption_capability` and `preemption_vulnerability` allow it. 1 is the highest level of priority. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `qos_indicator` - The QoS Indicator (5QI for 5G network /QCI for 4G net work) value identifies a set of QoS characteristics that control QoS forwarding treatment for QoS flows or EPS bearers.

* `guaranteed_bit_rate` - A `guaranteed_bit_rate` block as defined below. The Guaranteed Bit Rate (GBR) for all service data flows that use this PCC Rule. If it's not specified, there will be no GBR set for the PCC Rule that uses this QoS definition.

* `maximum_bit_rate` - A `maximum_bit_rate` block as defined below. The Maximum Bit Rate (MBR) for all service data flows that use this PCC Rule or Service.

* `preemption_capability` - The Preemption Capability of a QoS Flow controls whether it can preempt another QoS Flow with a lower priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `preemption_vulnerability` - The Preemption Vulnerability of a QoS Flow controls whether it can be preempted by QoS Flow with a higher priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

---

A `guaranteed_bit_rate` block supports the following:

* `downlink` - Downlink bit rate.

* `uplink` - Uplink bit rate.

---

A `service_data_flow_templates` block supports the following:

* `template_name` - Specifies the name of the data flow template. This must be unique within the parent data flow policy rule.

* `direction` - Specifies the direction of this flow. Possible values are `Uplink`, `Downlink` and `Bidirectional`.

* `protocol` - A list of the allowed protocol(s) for this flow. If you want this flow to be able to use any protocol within the internet protocol suite, use the value `ip`. If you only want to allow a selection of protocols, you must use the corresponding IANA Assigned Internet Protocol Number for each protocol, as described in https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml. For example, for UDP, you must use 17. If you use the value `ip` then you must leave the field `port` unspecified.

* `remote_ip_list` - The remote IP address(es) to which UEs will connect for this flow. If you want to allow connections on any IP address, use the value `any`. Otherwise, you must provide each of the remote IP addresses to which the packet core instance will connect for this flow. You must provide each IP address in CIDR notation, including the netmask (for example, `192.0.2.54/24`).

* `ports` - The port(s) to which UEs will connect for this flow. You can specify zero or more ports or port ranges. If you specify one or more ports or port ranges then you must specify a value other than `ip` in the `protocol` field. If it is not specified then connections will be allowed on all ports. Port ranges must be specified as <FirstPort>-<LastPort>. For example: [`8080`, `8082-8085`].

---

A `service_qos_policy` block supports the following:

* `allocation_and_retention_priority_level` - QoS Flow allocation and retention priority (ARP) level. Flows with higher priority preempt flows with lower priority, if the settings of `preemption_capability` and `preemption_vulnerability` allow it. 1 is the highest level of priority. If this field is not specified then `qos_indicator` is used to derive the ARP value. Defaults to `9`. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `qos_indicator` - The QoS Indicator (5QI for 5G network /QCI for 4G net work) value identifies a set of QoS characteristics that control QoS forwarding treatment for QoS flows or EPS bearers.

* `maximum_bit_rate` - A `maximum_bit_rate` block as defined below. The Maximum Bit Rate (MBR) for all service data flows that use this PCC Rule or Service.

* `preemption_capability` - The Preemption Capability of a QoS Flow controls whether it can preempt another QoS Flow with a lower priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `preemption_vulnerability` - The Preemption Vulnerability of a QoS Flow controls whether it can be preempted by QoS Flow with a higher priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters. Possible values are `NotPreemptable` and `Preemptable`. Defaults to `Preemptable`.

---

A `maximum_bit_rate` block supports the following:

* `downlink` - Downlink bit rate.

* `uplink` - Uplink bit rate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Service.
