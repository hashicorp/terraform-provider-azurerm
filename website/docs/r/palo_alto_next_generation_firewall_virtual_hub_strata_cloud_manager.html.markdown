---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager"
description: |-
  Manages a Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.
---

# azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager

Manages a Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "acceptanceTestPublicIp1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.0.0/23"

  tags = {
    "hubSaaSPreview" = "true"
  }
}

resource "azurerm_palo_alto_virtual_network_appliance" "example" {
  name           = "example-appliance"
  virtual_hub_id = azurerm_virtual_hub.example.id
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager" "example" {
  name                             = "example"
  resource_group_name              = "example"
  location                         = "West Europe"
  strata_cloud_manager_tenant_name = "example"

  network_profile {
    public_ip_address_ids        = [azurerm_public_ip.example.id]
    virtual_hub_id               = azurerm_virtual_hub.example.id
    network_virtual_appliance_id = azurerm_palo_alto_virtual_network_appliance.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager should exist. Changing this forces a new Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager to be created.

* `name` - (Required) The name which should be used for this Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager. Changing this forces a new Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager to be created.

* `network_profile` - (Required) A `network_profile` block as defined below.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager should exist. Changing this forces a new Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager to be created.

* `strata_cloud_manager_tenant_name` - (Required) Strata Cloud Manager name which is intended to manage the policy for this firewall.

---

* `destination_nat` - (Optional) One or more `destination_nat` blocks as defined below.

* `dns_settings` - (Optional) A `dns_settings` block as defined below.

* `identity` - (Optional) A `identity` block as defined below.

* `marketplace_offer_id` - (Optional) The ID of the marketplace offer. Defaults to `pan_swfw_cloud_ngfw`. Changing this forces a new Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager to be created.

* `plan_id` - (Optional) The ID of the billing plan. Defaults to `panw-cngfw-payg`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.

---

A `destination_nat` block supports the following:

* `name` - (Required) The name which should be used for this Destination NAT rule.

* `protocol` - (Required) The protocol used for this Destination NAT. Possible values include `TCP` and `UDP`.

* `backend_config` - (Optional) One or more `backend_config` block as defined below.

* `frontend_config` - (Optional) One or more `frontend_config` block as defined below.

---

A `dns_settings` block supports the following:

* `dns_servers` - (Optional) Specifies a list of DNS servers to use. Conflicts with `dns_settings.0.use_azure_dns`.

* `use_azure_dns` - (Optional) Should Azure DNS servers be used? Conflicts with `dns_settings.0.dns_servers`. Defaults to `false`.

---

A `backend_config` block supports the following:

* `port` - (Required) The port number to send traffic to.

* `public_ip_address` - (Required) The public IP Address to send the traffic to.

---

A `frontend_config` block supports the following:

* `port` - (Required) The port on which traffic will be received.

* `public_ip_address_id` - (Required) The ID of the Public IP Address resource the traffic will be received on.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this App Configuration. Possible values are `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this App Configuration.

---

A `network_profile` block supports the following:

* `network_virtual_appliance_id` - (Required) The ID of the Palo Alto Network Virtual Appliance in the VHub. Changing this forces a new Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager to be created.

* `public_ip_address_ids` - (Required) Specifies a list of Public IP IDs to use for this Next Generation Firewall.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub this Next Generation Firewall will be deployed in. Changing this forces a new Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager to be created.

* `egress_nat_ip_address_ids` - (Optional) Specifies a list of Public IP IDs to use for Egress NAT.

* `trusted_address_ranges` - (Optional) Specifies a list of trusted ranges to use for the Network.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.

* `identity` - An `identity` block as defined below.

* `network_profile` - A `network_profile` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

A `network_profile` block exports the following:

* `egress_nat_ip_addresses` - A list of Egress NAT IP addresses.

* `ip_of_trust_for_user_defined_routes` - The IP of trusted subnet for UDR.

* `public_ip_addresses` - A list of public IPs associated with this Next Generation Firewall.

* `trusted_subnet_id` - The ID of trusted subnet.

* `untrusted_subnet_id` - The ID of untrusted subnet.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.
* `update` - (Defaults to 3 hours) Used when updating the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.
* `delete` - (Defaults to 2 hours) Used when deleting the Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Manager.

## Import

Palo Alto Next Generation Firewall Virtual Hub Strata Cloud Managers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/firewalls/myVNetStrataCloudManagerFW
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `PaloAltoNetworks.Cloudngfw` - 2025-05-23
