---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_packet_core_control_plane"
description: |-
  Get information about a Mobile Network Packet Core Control Plane.
---

# azurerm_mobile_network_packet_core_control_plane

Get information about a Mobile Network Packet Core Control Plane.

## Example Usage

```hcl
data "azurerm_mobile_network_packet_core_control_plane" "example" {
  name                = "example-mnpccp"
  resource_group_name = "example-rg"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Mobile Network Packet Core Control Plane. 

* `resource_group_name` - The name of the Resource Group where the Mobile Network Packet Core Control Plane exists.  

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Packet Core Control Plane.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Mobile Network Packet Core Control Plane exists. 

* `site_ids` - The list of Mobile Network Site IDs in which this packet core control plane is deployed.

* `sku` - The SKU defining the throughput and SIM allowances for this packet core control plane deployment. 

* `local_diagnostics_access` - One or more `local_diagnostics_access` blocks as defined below. The Kubernetes ingress configuration that controls access to the packet core diagnostics through local APIs.

* `control_plane_access_interface` - A `control_plane_access_interface` block as defined below. The control plane interface on the access network.

* `user_equipment_mtu_in_bytes` - The MTU in bytes that can be sent to the user equipment.

* `core_network_technology` - The core network technology generation.

* `platform` - A `platform` block as defined below.

* `identity` - An `identity` block as defined below.

* `interoperability_settings_json` - Settings in JSON format to allow interoperability with third party components e.g. RANs and UEs.

* `tags` - A mapping of tags assigned to the Mobile Network Packet Core Control Plane.

* `software_version` - The version of the packet core software that is deployed.

---

A `control_plane_access_interface` block supports the following:

* `name` - The logical name for this interface.

* `ipv4_address` - The IPv4 address.

* `ipv4_subnet` - The IPv4 subnet.

* `ipv4_gateway` - The default IPv4 gateway (router).

---

A `local_diagnostics_access` block supports the following:

* `authentication_type` - How to authenticate users who access local diagnostics APIs. 

* `https_server_certificate_url` - The versionless certificate URL used to secure local access to packet core diagnostics over local APIs by the Kubernetes ingress.

---

An `identity` block supports the following:

* `type` - The type of Managed Service Identity.

* `identity_ids` - A list of User Assigned Managed Identity IDs assigned to this resource.

---

A `platform` block supports the following:

* `type` - The platform type where the packet core is deployed.

* `edge_device_id` - The ID of Azure Stack Edge device where the packet core is deployed. 

* `arc_kubernetes_cluster_id` - The ID of Azure Arc connected cluster where the packet core is deployed.

* `stack_hci_cluster_id` - The ID of Azure Stack HCI cluster where the packet core is deployed.

* `custom_location_id` -  The ID of Azure Arc custom location where the packet core is deployed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Packet Core Control Plane.
