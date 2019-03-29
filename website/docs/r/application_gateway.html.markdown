---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_gateway"
sidebar_current: "docs-azurerm-resource-network-application-gateway"
description: |-
  Manages an Application Gateway.
---

# azurerm_application_gateway

Manages an Application Gateway.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_virtual_network" "test" {
  name                = "example-network"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = ["10.254.0.0/16"]
}

resource "azurerm_subnet" "frontend" {
  name                 = "frontend"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_subnet" "backend" {
  name                 = "backend"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.254.2.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "example-pip"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  allocation_method   = "Dynamic"
}

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.test.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.test.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.test.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.test.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.test.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.test.name}-rqrt"
  redirect_configuration_name    = "${azurerm_virtual_network.test.name}-rdrcfg"
}

resource "azurerm_application_gateway" "network" {
  name                = "example-appgateway"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = "${azurerm_subnet.frontend.id}"
  }

  frontend_port {
    name = "${local.frontend_port_name}"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${local.frontend_ip_configuration_name}"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }

  backend_address_pool {
    name = "${local.backend_address_pool_name}"
  }

  backend_http_settings {
    name                  = "${local.http_setting_name}"
    cookie_based_affinity = "Disabled"
    path         = "/path1/"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = "${local.listener_name}"
    frontend_ip_configuration_name = "${local.frontend_ip_configuration_name}"
    frontend_port_name             = "${local.frontend_port_name}"
    protocol                       = "Http"
  }

  request_routing_rule {
    name                        = "${local.request_routing_rule_name}"
    rule_type                   = "Basic"
    http_listener_name          = "${local.listener_name}"
    backend_address_pool_name   = "${local.backend_address_pool_name}"
    backend_http_settings_name  = "${local.http_setting_name}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Application Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to the Application Gateway should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure region where the Application Gateway should exist. Changing this forces a new resource to be created.

* `backend_address_pool` - (Required) One or more `backend_address_pool` blocks as defined below.

* `backend_http_settings` - (Required) One or more `backend_http_settings` blocks as defined below.

* `frontend_ip_configuration` - (Required) One or more `frontend_ip_configuration` blocks as defined below.

* `frontend_port` - (Required) One or more `frontend_port` blocks as defined below.

* `gateway_ip_configuration` - (Required) One or more `gateway_ip_configuration` blocks as defined below.

* `http_listener` - (Required) One or more `http_listener` blocks as defined below.

* `request_routing_rule` - (Required) One or more `request_routing_rule` blocks as defined below.

* `sku` - (Required) A `sku` block as defined below.

* `zones` - (Optional) A collection of availability zones to spread the Application Gateway over.

-> **Please Note**: Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview).  They are also only supported for [v2 SKUs](https://docs.microsoft.com/en-us/azure/application-gateway/application-gateway-autoscaling-zone-redundant)

---

* `authentication_certificate` - (Optional) One or more `authentication_certificate` blocks as defined below.

* `disabled_ssl_protocols` - (Optional) A list of SSL Protocols which should be disabled on this Application Gateway. Possible values are `TLSv1_0`, `TLSv1_1` and `TLSv1_2`.

* `enable_http2` - (Optional) Is HTTP2 enabled on the application gateway resource? Defaults to `false`.

* `probe` - (Optional) One or more `probe` blocks as defined below.

* `ssl_certificate` - (Optional) One or more `ssl_certificate` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `url_path_map` - (Optional) One or more `url_path_map` blocks as defined below.

* `waf_configuration` - (Optional) A `waf_configuration` block as defined below.

* `custom_error_configuration` - (Optional) One or more `custom_error_configuration` blocks as defined below.

* `redirect_configuration` - (Optional) A `redirect_configuration` block as defined below.

---

A `authentication_certificate` block supports the following:

* `name` - (Required) The Name of the Authentication Certificate to use.

* `data` - (Required) The contents of the Authentication Certificate which should be used.

---

A `authentication_certificate` block, within the `backend_http_settings` block supports the following:

* `name` - (Required) The name of the Authentication Certificate.

---

A `backend_address_pool` block supports the following:

* `name` - (Required) The name of the Backend Address Pool.

* `fqdns` - (Optional) A list of FQDN's which should be part of the Backend Address Pool.

* `fqdn_list` - (Optional **Deprecated**) A list of FQDN's which should be part of the Backend Address Pool. This field has been deprecated in favour of `fqdns` and will be removed in v2.0 of the AzureRM Provider.

* `ip_addresses` - (Optional) A list of IP Addresses which should be part of the Backend Address Pool.

* `ip_address_list` - (Optional **Deprecated**) A list of IP Addresses which should be part of the Backend Address Pool. This field has been deprecated in favour of `ip_addresses` and will be removed in v2.0 of the AzureRM Provider.

---

A `backend_http_settings` block supports the following:

* `cookie_based_affinity` - (Required) Is Cookie-Based Affinity enabled? Possible values are `Enabled` and `Disabled`.

* `name` - (Required) The name of the Backend HTTP Settings Collection.

* `path` - (Optional) The Path which should be used as a prefix for all HTTP requests.

* `port`- (Required) The port which should be used for this Backend HTTP Settings Collection.

* `probe_name` - (Required) The name of an associated HTTP Probe.

* `protocol`- (Required) The Protocol which should be used. Possible values are `Http` and `Https`.

* `request_timeout` - (Required) The request timeout in seconds, which must be between 1 and 86400 seconds.

* `host_name` - (Optional) Host header to be sent to the backend servers. Cannot be set if `pick_host_name_from_backend_address` is set to `true`.

* `pick_host_name_from_backend_address` - (Optional) Whether host header should be picked from the host name of the backend server. Defaults to `false`.

* `authentication_certificate` - (Optional) One or more `authentication_certificate` blocks.

* `connection_draining` - (Optional) A `connection_draining` block as defined below.

---

A `connection_draining` block supports the following:

* `enabled` - (Required) If connection draining is enabled or not.

* `drain_timeout_sec` - (Required) The number of seconds connection draining is active. Acceptable values are from `1` second to `3600` seconds.

---

      
A `frontend_ip_configuration` block supports the following:

* `name` - (Required) The name of the Frontend IP Configuration.

* `subnet_id` - (Required) The ID of the Subnet which the Application Gateway should be connected to.

* `private_ip_address` - (Optional) The Private IP Address to use for the Application Gateway.

* `public_ip_address_id` - (Optional) The ID of a Public IP Address which the Application Gateway should use.

-> **NOTE:** The Allocation Method for this Public IP Address should be set to `Dynamic`.

* `private_ip_address_allocation` - (Optional) The Allocation Method for the Private IP Address. Possible values are `Dynamic` and `Static`.

---

A `frontend_port` block supports the following:

* `name` - (Required) The name of the Frontend Port.

* `port` - (Required) The port used for this Frontend Port.

---

A `gateway_ip_configuration` block supports the following:

* `name` - (Required) The Name of this Gateway IP Configuration.

* `subnet_id` - (Required) The ID of a Subnet.

---

A `http_listener` block supports the following:

* `name` - (Required) The Name of the HTTP Listener.

* `frontend_ip_configuration_name` - (Required) The Name of the Frontend IP Configuration used for this HTTP Listener.

* `frontend_port_name` - (Required) The Name of the Frontend Port use for this HTTP Listener.

* `host_name` - (Optional) The Hostname which should be used for this HTTP Listener.

* `protocol` - (Required) The Protocol to use for this HTTP Listener. Possible values are `Http` and `Https`.

* `require_sni` - (Optional) Should Server Name Indication be Required? Defaults to `false`.

* `ssl_certificate_name` - (Optional) The name of the associated SSL Certificate which should be used for this HTTP Listener.

* `custom_error_configuration` - (Optional) One or more `custom_error_configuration` blocks as defined below.

---

A `match` block supports the following:

* `body` - (Optional) A snippet from the Response Body which must be present in the Response. Defaults to `*`.

* `status_code` - (Optional) A list of allowed status codes for this Health Probe.

---

A `path_rule` block supports the following:

* `name` - (Required) The Name of the Path Rule.

* `paths` - (Required) A list of Paths used in this Path Rule.

* `backend_address_pool_name` - (Optional) The Name of the Backend Address Pool to use for this Path Rule. Cannot be set if `redirect_configuration_name` is set.

* `backend_http_settings_name` - (Optional) The Name of the Backend HTTP Settings Collection to use for this Path Rule. Cannot be set if `redirect_configuration_name` is set.

* `redirect_configuration_name` - (Optional) The Name of a Redirect Configuration to use for this Path Rule. Cannot be set if `backend_address_pool_name` or `backend_http_settings_name` is set.

---

A `probe` block support the following:

* `host` - (Optional) The Hostname used for this Probe. If the Application Gateway is configured for a single site, by default the Host name should be specified as ‘127.0.0.1’, unless otherwise configured in custom probe. Cannot be set if `pick_host_name_from_backend_http_settings` is set to `true`.

* `interval` - (Required) The Interval between two consecutive probes in seconds. Possible values range from 1 second to a maximum of 86,400 seconds.

* `name` - (Required) The Name of the Probe.

* `protocol` - (Required) The Protocol used for this Probe. Possible values are `Http` and `Https`.

* `path` - (Required) The Path used for this Probe.

* `timeout` - (Required) The Timeout used for this Probe, which indicates when a probe becomes unhealthy. Possible values range from 1 second to a maximum of 86,400 seconds.

* `unhealthy_threshold` - (Required) The Unhealthy Threshold for this Probe, which indicates the amount of retries which should be attempted before a node is deemed unhealthy. Possible values are from 1 - 20 seconds.

* `pick_host_name_from_backend_http_settings` - (Optional) Whether the host header should be picked from the backend http settings. Defaults to `false`.

* `match` - (Optional) A `match` block as defined above.

* `minimum_servers` - (Optional) The minimum number of servers that are always marked as healthy. Defaults to `0`.

---

A `request_routing_rule` block supports the following:

* `name` - (Required) The Name of this Request Routing Rule.

* `rule_type` - (Required) The Type of Routing that should be used for this Rule. Possible values are `Basic` and `PathBasedRouting`.

* `http_listener_name` - (Required) The Name of the HTTP Listener which should be used for this Routing Rule.

* `backend_address_pool_name` - (Optional) The Name of the Backend Address Pool which should be used for this Routing Rule. Cannot be set if `redirect_configuration_name` is set.

* `backend_http_settings_name` - (Optional) The Name of the Backend HTTP Settings Collection which should be used for this Routing Rule. Cannot be set if `redirect_configuration_name` is set.

* `redirect_configuration_name` - (Optional) The Name of the Redirect Configuration which should be used for this Routing Rule. Cannot be set if either `backend_address_pool_name` or `backend_http_settings_name` is set.

* `url_path_map_name` - (Optional) The Name of the URL Path Map which should be associated with this Routing Rule.

---

A `sku` block supports the following:

* `name` - (Required) The Name of the SKU to use for this Application Gateway. Possible values are `Standard_Small`, `Standard_Medium`, `Standard_Large`, `Standard_v2`, `WAF_Medium`, `WAF_Large`, and `WAF_v2`.

* `tier` - (Required) The Tier of the SKU to use for this Application Gateway. Possible values are `Standard`, `Standard_v2`, `WAF` and `WAF_v2`.

* `capacity` - (Required) The Capacity of the SKU to use for this Application Gateway - which must be between 1 and 10.

---

A `ssl_certificate` block supports the following:

* `name` - (Required) The Name of the SSL certificate that is unique within this Application Gateway

* `data` - (Required) PFX certificate.

* `password` - (Required) Password for the pfx file specified in data.

---

A `url_path_map` block supports the following:

* `name` - (Required) The Name of the URL Path Map.

* `default_backend_address_pool_name` - (Optional) The Name of the Default Backend Address Pool which should be used for this URL Path Map. Cannot be set if there are path_rules with re-direct configurations set.

* `default_backend_http_settings_name` - (Optional) The Name of the Default Backend HTTP Settings Collection which should be used for this URL Path Map. Cannot be set if there are path_rules with re-direct configurations set.

* `default_redirect_configuration_name` - (Optional) The Name of the Default Redirect Configuration which should be used for this URL Path Map. Cannot be set if there are path_rules with Backend Address Pool or HTTP Settings set.

* `path_rule` - (Required) One or more `path_rule` blocks as defined above.

---

A `waf_configuration` block supports the following:

* `enabled` - (Required) Is the Web Application Firewall be enabled?

* `firewall_mode` - (Required) The Web Application Firewall Mode. Possible values are `Detection` and `Prevention`.

* `rule_set_type` - (Required) The Type of the Rule Set used for this Web Application Firewall.

* `rule_set_version` - (Required) The Version of the Rule Set used for this Web Application Firewall.

* `file_upload_limit_mb` - (Optional) The File Upload Limit in MB. Accepted values are in the range `1`MB to `500`MB. Defaults to `100`MB.

* `request_body_check` - (Optional) Is Request Body Inspection enabled?  Defaults to `true`.

* `max_request_body_size_kb` - (Optional) The Maximum Request Body Size in KB.  Accepted values are in the range `1`KB to `128`KB.  Defaults to `128`KB.

---

A `custom_error_configuration` block supports the following:

* `status_code` - (Required) Status code of the application gateway customer error. Possible values are `HttpStatus403` and `HttpStatus502`

* `custom_error_page_url` - (Required) Error page URL of the application gateway customer error.

---

A `redirect_configuration` block supports the following:

* `name` - (Required) Unique name of the redirect configuration block

* `redirect_type` - (Required) The type of redirect. Possible values are `Permanent`, `Temporary`, `Found` and `SeeOther`

* `target_listener_name` - (Optional) The name of the listener to redirect to. Cannot be set if `target_url` is set.

* `target_url` - (Optional) The Url to redirect the request to. Cannot be set if `target_listener_name` is set.

* `include_path` - (Optional) Whether or not to include the path in the redirected Url. Defaults to `false`

* `include_query_string` - (Optional) Whether or not to include the query string in the redirected Url. Default to `false`

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Gateway.

* `authentication_certificate` - A list of `authentication_certificate` blocks as defined below.

* `backend_address_pool` - A list of `backend_address_pool` blocks as defined below.

* `backend_http_settings` - A list of `backend_http_settings` blocks as defined below.

* `frontend_ip_configuration` - A list of `frontend_ip_configuration` blocks as defined below.

* `frontend_port` - A list of `frontend_port` blocks as defined below.

* `gateway_ip_configuration` - A list of `gateway_ip_configuration` blocks as defined below.

* `enable_http2` - (Optional) Is HTTP2 enabled on the application gateway resource? Defaults to `false`.

* `http_listener` - A list of `http_listener` blocks as defined below.

* `probe` - A `probe` block as defined below.

* `request_routing_rule` - A list of `request_routing_rule` blocks as defined below.

* `ssl_certificate` - A list of `ssl_certificate` blocks as defined below.

* `url_path_map` - A list of `url_path_map` blocks as defined below.

* `custom_error_configuration` - A list of `custom_error_configuration` blocks as defined below.

* `redirect_configuration` - A list of `redirect_configuration` blocks as defined below.
---

A `authentication_certificate` block exports the following:

* `id` - The ID of the Authentication Certificate.

---

A `authentication_certificate` block, within the `backend_http_settings` block exports the following:

* `id` - The ID of the Authentication Certificate.

---

A `backend_address_pool` block exports the following:

* `id` - The ID of the Backend Address Pool.

---

A `backend_http_settings` block exports the following:

* `id` - The ID of the Backend HTTP Settings Configuration.

* `probe_id` - The ID of the associated Probe.

---

A `frontend_ip_configuration` block exports the following:

* `id` - The ID of the Frontend IP Configuration.

---

A `frontend_port` block exports the following:

* `id` - The ID of the Frontend Port.

---

A `gateway_ip_configuration` block exports the following:

* `id` - The ID of the Gateway IP Configuration.

---

A `http_listener` block exports the following:

* `id` - The ID of the HTTP Listener.

* `frontend_ip_configuration_id` - The ID of the associated Frontend Configuration.

* `frontend_port_id` - The ID of the associated Frontend Port.

* `ssl_certificate_id` - The ID of the associated SSL Certificate.

---

A `path_rule` block exports the following:

* `id` - The ID of the Path Rule.

* `backend_address_pool_id` - The ID of the Backend Address Pool used in this Path Rule.

* `backend_http_settings_id` - The ID of the Backend HTTP Settings Collection used in this Path Rule.

* `redirect_configuration_id` - The ID of the Redirect Configuration used in this Path Rule.

---

A `probe` block exports the following:

* `id` - The ID of the Probe.

---

A `request_routing_rule` block exports the following:

* `id` - The ID of the Request Routing Rule.

* `http_listener_id` - The ID of the associated HTTP Listener.

* `backend_address_pool_id` - The ID of the associated Backend Address Pool.

* `backend_http_settings_id` - The ID of the associated Backend HTTP Settings Configuration.

* `redirect_configuration_id` - The ID of the associated Redirect Configuration.

* `url_path_map_id` - The ID of the associated URL Path Map.

---

A `ssl_certificate` block exports the following:

* `id` - The ID of the SSL Certificate.

* `public_cert_data` - The Public Certificate Data associated with the SSL Certificate.

---

A `url_path_map` block exports the following:

* `id` - The ID of the URL Path Map.

* `default_backend_address_pool_id` - The ID of the Default Backend Address Pool.

* `default_backend_http_settings_id` - The ID of the Default Backend HTTP Settings Collection.

* `default_redirect_configuration_id` - The ID of the Default Redirect Configuration.

* `path_rule` - A list of `path_rule` blocks as defined above.

---

A `custom_error_configuration` block exports the following:

* `id` - The ID of the Custom Error Configuration.

---

A `redirect_configuration` block exports the following:

* `id` - The ID of the Redirect Configuration.

## Import

Application Gateway's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_gateway.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/applicationGateways/myGateway1
```
