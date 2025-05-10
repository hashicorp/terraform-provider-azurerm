---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_packet_core_control_plane"
description: |-
  Manages a Mobile Network Packet Core Control Plane.
---

# azurerm_mobile_network_packet_core_control_plane

Manages a Mobile Network Packet Core Control Plane.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_site" "example" {
  name              = "example-mns"
  mobile_network_id = azurerm_mobile_network.test.id
  location          = azurerm_resource_group.example.location
}

resource "azurerm_databox_edge_device" "example" {
  name                = "example-device"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku_name = "EdgeP_Base-Standard"
}

resource "azurerm_mobile_network_packet_core_control_plane" "example" {
  name                              = "example-mnpccp"
  resource_group_name               = azurerm_resource_group.example.name
  location                          = azurerm_resource_group.example.location
  sku                               = "G0"
  control_plane_access_name         = "default-interface"
  control_plane_access_ipv4_address = "192.168.1.199"
  control_plane_access_ipv4_gateway = "192.168.1.1"
  control_plane_access_ipv4_subnet  = "192.168.1.0/25"
  site_ids                          = [azurerm_mobile_network_site.example.id]

  local_diagnostics_access {
    authentication_type = "AAD"
  }

  platform {
    type           = "AKS-HCI"
    edge_device_id = azurerm_databox_edge_device.example.id
  }

  interoperability_settings_json = jsonencode({
    "key" = "value"
  })

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies The name of the Mobile Network Packet Core Control Plane. Changing this forces a new Mobile Network Packet Core Control Plane to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Mobile Network Packet Core Control Plane should exist. Changing this forces a new Mobile Network Packet Core Control Plane to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Packet Core Control Plane should exist. Changing this forces a new Mobile Network Packet Core Control Plane to be created.

* `site_ids` - (Required) A list of Mobile Network Site IDs in which this packet core control plane should be deployed. The Sites must be in the same location as the packet core control plane.

* `sku` - (Required) The SKU defining the throughput and SIM allowances for this packet core control plane deployment. Possible values are `G0`, `G1`, `G2`, `G3`, `G4`, `G5` and `G10`.

* `local_diagnostics_access` - (Required) One or more `local_diagnostics_access` blocks as defined below. Specifies the Kubernetes ingress configuration that controls access to the packet core diagnostics through local APIs.

* `control_plane_access_name` - (Optional) Specifies the logical name for this interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `control_plane_access_ipv4_address` - (Optional) The IPv4 address for the control plane interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `control_plane_access_ipv4_subnet` - (Optional) The IPv4 subnet for the control plane interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `control_plane_access_ipv4_gateway` - (Optional) The default IPv4 gateway for the control plane interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `user_equipment_mtu_in_bytes` - (Optional) Specifies the MTU in bytes that can be sent to the user equipment. The same MTU is set on the user plane data links for all data networks. The MTU set on the user plane access link will be 60 bytes greater than this value to allow for GTP encapsulation.

* `core_network_technology` - (Optional) The core network technology generation. Possible values are `5GC` and `EPC`.

* `platform` - (Optional) A `platform` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `interoperability_settings_json` - (Optional) Settings in JSON format to allow interoperability with third party components e.g. RANs and UEs.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Packet Core Control Plane.

* `software_version` - (Optional) Specifies the version of the packet core software that is deployed.

---

A `local_diagnostics_access` block supports the following:

* `authentication_type` - (Required) How to authenticate users to access local diagnostics APIs. Possible values are `AAD` and `Password`.

* `https_server_certificate_url` - (Optional) The versionless certificate URL used to secure local access to packet core diagnostics over local APIs by the Kubernetes ingress.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Required) A list of the IDs for User Assigned Managed Identity resources to be assigned.

---

A `platform` block supports the following:

* `type` - (Required) Specifies the platform type where the packet core is deployed. Possible values are `AKS-HCI`, `3P-AZURE-STACK-HCI` and `BaseVM`.

* `edge_device_id` - (Optional) The ID of the Azure Stack Edge device where the packet core is deployed. If the device is part of a fault-tolerant pair, either device in the pair can be specified.

* `arc_kubernetes_cluster_id` - (Optional) The ID of the Azure Arc connected cluster where the packet core is deployed.

* `stack_hci_cluster_id` - (Optional) The ID of the Azure Stack HCI cluster where the packet core is deployed.

* `custom_location_id` - (Optional) The ID of the Azure Arc custom location where the packet core is deployed.

~> **Note:** At least one of `edge_device_id`, `arc_kubernetes_cluster_id`, `stack_hci_cluster_id` and `custom_location_id` should be specified. If multiple are set, they must be consistent with each other.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Packet Core Control Plane.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Packet Core Control Plane.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Packet Core Control Plane.
* `update` - (Defaults to 30 minutes) Used when updating the Mobile Network Packet Core Control Plane.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Packet Core Control Plane.

## Import

Mobile Network Packet Core Control Plane can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_packet_core_control_plane.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreControlPlane1
```
