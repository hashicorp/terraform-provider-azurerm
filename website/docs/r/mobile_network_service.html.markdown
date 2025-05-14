---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_service"
description: |-
  Manages a Mobile Network Service.
---

# azurerm_mobile_network_service

Manages a Mobile Network Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "east us"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  mobile_country_code = "001"
  mobile_network_code = "01"
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

    qos_policy {
      allocation_and_retention_priority_level = 9
      qos_indicator                           = 9
      preemption_capability                   = "NotPreempt"
      preemption_vulnerability                = "Preemptable"

      guaranteed_bit_rate {
        downlink = "100 Mbps"
        uplink   = "10 Mbps"
      }

      maximum_bit_rate {
        downlink = "1 Gbps"
        uplink   = "100 Mbps"
      }
    }

    service_data_flow_template {
      direction      = "Uplink"
      name           = "IP-to-server"
      ports          = []
      protocol       = ["ip"]
      remote_ip_list = ["10.3.4.0/24"]
    }
  }

  service_qos_policy {
    allocation_and_retention_priority_level = 9
    qos_indicator                           = 9
    preemption_capability                   = "NotPreempt"
    preemption_vulnerability                = "Preemptable"
    maximum_bit_rate {
      downlink = "1 Gbps"
      uplink   = "100 Mbps"
    }
  }

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Service. Changing this forces a new Mobile Network Service to be created.

* `mobile_network_id` - (Required) Specifies the ID of the Mobile Network Service. Changing this forces a new Mobile Network Service to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Service should exist. Changing this forces a new Mobile Network Service to be created.

* `service_precedence` - (Required) A precedence value that is used to decide between services when identifying the QoS values to use for a particular SIM. A lower value means a higher priority. This value should be unique among all services configured in the mobile network. Must be between `0` and `255`.

* `pcc_rule` - (Required) A `pcc_rule` block as defined below. The set of PCC Rules that make up this service.

* `service_qos_policy` - (Optional) A `service_qos_policy` block as defined below. The QoS policy to use for packets matching this service. This can be overridden for particular flows using the ruleQosPolicy field in a `pcc_rule`. If this field is not specified then the `sim_policy` of User Equipment (UE) will define the QoS settings.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Service.

---

A `pcc_rule` block supports the following:

* `name` - (Required) Specifies the name of the rule. This must be unique within the parent service. You must not use any of the following reserved strings - `default`, `requested` or `service`.

* `precedence` - (Required) A precedence value that is used to decide between data flow policy rules when identifying the QoS values to use for a particular SIM. A lower value means a higher priority. This value should be unique among all data flow policy rules configured in the mobile network. Must be between `0` and `255`.

* `qos_policy` - (Optional) A `qos_policy` block as defined below. The QoS policy to use for packets matching this rule. If this field is not specified then the Service will define the QoS settings.

* `service_data_flow_template` - (Required) A `service_data_flow_template` block as defined below. The set of service data flow templates to use for this PCC rule.

* `traffic_control_enabled` - (Optional) Determines whether flows that match this data flow policy rule are permitted. Defaults to `true`.

---

A `qos_policy` block supports the following:

* `allocation_and_retention_priority_level` - (Optional) QoS Flow allocation and retention priority (ARP) level. Flows with higher priority preempt flows with lower priority, if the settings of `preemption_capability` and `preemption_vulnerability` allow it. 1 is the highest level of priority. If this field is not specified then `qos_indicator` is used to derive the ARP value. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `qos_indicator` - (Required) The QoS Indicator (5QI for 5G network /QCI for 4G net work) value identifies a set of QoS characteristics that control QoS forwarding treatment for QoS flows or EPS bearers. Recommended values: 5-9; 69-70; 79-80. Must be between `1` and `127`.

* `guaranteed_bit_rate` - (Optional) A `guaranteed_bit_rate` block as defined below. The Guaranteed Bit Rate (GBR) for all service data flows that use this PCC Rule. If it's not specified, there will be no GBR set for the PCC Rule that uses this QoS definition.

* `maximum_bit_rate` - (Required) A `maximum_bit_rate` block as defined below. The Maximum Bit Rate (MBR) for all service data flows that use this PCC Rule or Service.

* `preemption_capability` - (Optional) The Preemption Capability of a QoS Flow controls whether it can preempt another QoS Flow with a lower priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters. Possible values are `NotPreempt` and `MayPreempt`, Defaults to `NotPreempt`.

* `preemption_vulnerability` - (Optional) The Preemption Vulnerability of a QoS Flow controls whether it can be preempted by QoS Flow with a higher priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters. Possible values are `NotPreemptable` and `Preemptable`. Defaults to `Preemptable`.

---

A `guaranteed_bit_rate` block supports the following:

* `downlink` - (Required) Downlink bit rate. Must be a number followed by `Kbps`, `Mbps`, `Gbps` or `Tbps`.

* `uplink` - (Required) Uplink bit rate. Must be a number followed by `Kbps`, `Mbps`, `Gbps` or `Tbps`.

---

A `service_data_flow_template` block supports the following:

* `name` - (Required) Specifies the name of the data flow template. This must be unique within the parent data flow policy rule. You must not use any of the following reserved strings - `default`, `requested` or `service`.

* `direction` - (Required) Specifies the direction of this flow. Possible values are `Uplink`, `Downlink` and `Bidirectional`.

* `protocol` - (Required) A list of the allowed protocol(s) for this flow. If you want this flow to be able to use any protocol within the internet protocol suite, use the value `ip`. If you only want to allow a selection of protocols, you must use the corresponding IANA Assigned Internet Protocol Number for each protocol, as described in https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml. For example, for UDP, you must use 17. If you use the value `ip` then you must leave the field `port` unspecified.

* `remote_ip_list` - (Required) Specifies the remote IP address(es) to which UEs will connect for this flow. If you want to allow connections on any IP address, use the value `any`. Otherwise, you must provide each of the remote IP addresses to which the packet core instance will connect for this flow. You must provide each IP address in CIDR notation, including the netmask (for example, `192.0.2.54/24`).

* `ports` - (Optional) The port(s) to which UEs will connect for this flow. You can specify zero or more ports or port ranges. If you specify one or more ports or port ranges then you must specify a value other than `ip` in the `protocol` field. If it is not specified then connections will be allowed on all ports. Port ranges must be specified as <FirstPort>-<LastPort>. For example: [`8080`, `8082-8085`].

---

A `service_qos_policy` block supports the following:

* `allocation_and_retention_priority_level` - (Optional) QoS Flow allocation and retention priority (ARP) level. Flows with higher priority preempt flows with lower priority, if the settings of `preemption_capability` and `preemption_vulnerability` allow it. 1 is the highest level of priority. If this field is not specified then `qos_indicator` is used to derive the ARP value. Defaults to `9`. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters.

* `qos_indicator` - (Optional) The QoS Indicator (5QI for 5G network /QCI for 4G net work) value identifies a set of QoS characteristics that control QoS forwarding treatment for QoS flows or EPS bearers. Recommended values: 5-9; 69-70; 79-80. Must be between `1` and `127`.

* `maximum_bit_rate` - (Required) A `maximum_bit_rate` block as defined below. The Maximum Bit Rate (MBR) for all service data flows that use this PCC Rule or Service.

* `preemption_capability` - (Optional) The Preemption Capability of a QoS Flow controls whether it can preempt another QoS Flow with a lower priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters. Possible values are `NotPreempt` and `MayPreempt`,.

* `preemption_vulnerability` - (Optional) The Preemption Vulnerability of a QoS Flow controls whether it can be preempted by QoS Flow with a higher priority level. See 3GPP TS23.501 section 5.7.2.2 for a full description of the ARP parameters. Possible values are `NotPreemptable` and `Preemptable`.

---

A `maximum_bit_rate` block supports the following:

* `downlink` - (Required) Downlink bit rate. Must be a number followed by `bps`, `Kbps`, `Mbps`, `Gbps` or `Tbps`.

* `uplink` - (Required) Uplink bit rate. Must be a number followed by `bps`, `Kbps`, `Mbps`, `Gbps` or `Tbps`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Service.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Service.
* `update` - (Defaults to 3 hours) Used when updating the Mobile Network Service.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Service.

## Import

Mobile Network Service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/mobileNetworks/mobileNetwork1/services/service1
```
