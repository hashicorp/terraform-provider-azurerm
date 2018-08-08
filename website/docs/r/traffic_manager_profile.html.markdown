---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_profile"
sidebar_current: "docs-azurerm-resource-network-traffic-manager-profile"
description: |-
  Manages a Traffic Manager Profile.

---

# azurerm_traffic_manager_profile

Manages a Traffic Manager Profile to which multiple endpoints can be attached.

## Example Usage


```hcl
resource "azurerm_resource_group" "test" {
  name     = "trafficmanagerProfile"
  location = "West US"
}

resource "azurerm_traffic_manager_profile" "test" {
  name                   = "${random_id.server.hex}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "${random_id.server.hex}"
    ttl           = 100
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }

  tags {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual network. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the virtual network.

* `dns_config` - (Required) A `dns_config` block as defined below.

* `monitor_config` - (Required) A `monitor_config` block as defined below.

* `traffic_routing_method` - (Required) Specifies the algorithm used to route traffic, possible values are:

    - `Geographic` - where traffic is routed based on Geographic regions specified in the Endpoint.
    - `Performance` - where traffic is routed via the User's closest Endpoint.
    - `Priority` - where traffic is routed to the Endpoint with the lowest `priority` value.
    - `Weighted` - where traffic is spread across Endpoints proportional to their `weight` value.

---

* `profile_status` - (Optional) The status of the profile. Possible values are `Enabled` or `Disabled`. Defaults to `Enabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `dns_config` block supports the following:

* `relative_name` - (Required) The relative domain name, this is combined with the domain name used by Traffic Manager to form the FQDN which is exported as documented below. Changing this forces a new resource to be created.

* `ttl` - (Required) The TTL value of the Profile used by Local DNS resolvers and clients.

---

A `monitor_config` block supports the following:

* `protocol` - (Required) The protocol used by the monitoring checks, supported values are `HTTP`, `HTTPS` and `TCP`.

* `port` - (Required) The port number used by the monitoring checks.

* `path` - (Optional) The path used by the monitoring checks. Required when `protocol` is set to `HTTP` or `HTTPS` - cannot be set when `protocol` is set to `TCP`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Traffic Manager Profile.
* `fqdn` - The FQDN of the created Profile.

## Import

Traffic Manager Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_traffic_manager_profile.testProfile /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/trafficManagerProfiles/mytrafficmanagerprofile1
```
