---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_gateway"
sidebar_current: "docs-azurerm-resource-application-gateway"
description: |-
  Manages a application gateway based on a previously created virtual network with configured subnets.
---

# azurerm_application_gateway

Manages a application gateway based on a previously created virtual network with configured subnets.

## Example Usage

```hcl
# Create a resource group
resource "azurerm_resource_group" "rg" {
  name     = "my-rg-application-gateway-12345"
  location = "West US"
}

# Create a application gateway in the web_servers resource group
resource "azurerm_virtual_network" "vnet" {
  name                = "my-vnet-12345"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.rg.location}"
}

resource "azurerm_subnet" "sub1" {
  name                 = "my-subnet-1"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_subnet" "sub2" {
  name                 = "my-subnet-2"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  address_prefix       = "10.254.2.0/24"
}

resource "azurerm_public_ip" "pip" {
  name                         = "my-pip-12345"
  location                     = "${azurerm_resource_group.rg.location}"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  public_ip_address_allocation = "dynamic"
}

# Create an application gateway
resource "azurerm_application_gateway" "network" {
  name                = "my-application-gateway-12345"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  location            = "West US"

  sku {
    name           = "Standard_Small"
    tier           = "Standard"
    capacity       = 2
  }

  gateway_ip_configuration {
      name         = "my-gateway-ip-configuration"
      subnet_id    = "${azurerm_virtual_network.vnet.id}/subnets/${azurerm_subnet.sub1.name}"
  }

  frontend_port {
      name         = "${azurerm_virtual_network.vnet.name}-feport"
      port         = 80
  }

  frontend_ip_configuration {
      name         = "${azurerm_virtual_network.vnet.name}-feip"
      public_ip_address_id = "${azurerm_public_ip.pip.id}"
  }

  backend_address_pool {
      name = "${azurerm_virtual_network.vnet.name}-beap"
  }

  backend_http_settings {
      name                  = "${azurerm_virtual_network.vnet.name}-be-htst"
      cookie_based_affinity = "Disabled"
      port                  = 80
      protocol              = "Http"
     request_timeout        = 1
  }

  http_listener {
        name                                  = "${azurerm_virtual_network.vnet.name}-httplstn"
        frontend_ip_configuration_name        = "${azurerm_virtual_network.vnet.name}-feip"
        frontend_port_name                    = "${azurerm_virtual_network.vnet.name}-feport"
        protocol                              = "Http"
  }

  request_routing_rule {
          name                       = "${azurerm_virtual_network.vnet.name}-rqrt"
          rule_type                  = "Basic"
          http_listener_name         = "${azurerm_virtual_network.vnet.name}-httplstn"
          backend_address_pool_name  = "${azurerm_virtual_network.vnet.name}-beap"
          backend_http_settings_name = "${azurerm_virtual_network.vnet.name}-be-htst"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the application gateway. Changing this forces a
  new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
  create the application gateway.

* `location` - (Required) The location/region where the application gateway is
  created. Changing this forces a new resource to be created.

* `sku` - (Required) Specifies size, tier and capacity of the application gateway. Must be specified once. The `sku` block fields documented below.

* `gateway_ip_configuration` - (Required) List of subnets that the application gateway is deployed into. The application gateway must be deployed into an existing virtual network/subnet. No other resource can be deployed in a subnet where application gateway is deployed. The `gateway_ip_configuration` block supports fields documented below.

* `frontend_port` - (Required) Front-end port for the application gateway. The `frontend_port` block supports fields documented below.

* `frontend_ip_configuration` - (Required) Specifies lists of frontend IP configurations. Currently only one Public and/or one Private IP address can be specified. Also one frontendIpConfiguration element can specify either Public or Private IP address, not both. The `frontend_ip_configuration` block supports fields documented below.

* `backend_address_pool` - (Required) Backend pools can be composed of NICs, virtual machine scale sets, public IPs, internal IPs, fully qualified domain names (FQDN), and multi-tenant back-ends like Azure Web Apps. Application Gateway backend pool members are not tied to an availability set. Members of backend pools can be across clusters, data centers, or outside of Azure as long as they have IP connectivity. The `backend_address_pool` block supports fields documented below.

* `backend_http_settings` - (Required) Related group of backend http and/or https features to be applied when routing to backend address pools. The `backend_http_settings` block supports fields documented below.

* `http_listener` - (Required) 1 or more listeners specifying port, http or https and SSL certificate (if configuring SSL offload) Each `http_listener` is attached to a `frontend_ip_configuration`. The `http_listener` block supports fields documented below.

* `probe` - (Optional) Specifies list of URL probes. The `probe` block supports fields documented below.

* `request_routing_rule` - (Required) Request routing rules can be either Basic or Path Based. Request routing rules are order sensitive. The `request_routing_rule` block supports fields documented below.

* `url_path_map` - (Optional) UrlPathMaps give url Path to backend mapping information for PathBasedRouting specified in `request_routing_rule`. The `url_path_map` block supports fields documented below.

* `authentication_certificate` - (Optional) List of authentication certificates. The `authentication_certificate` block supports fields documented below.

* `ssl_certificate` - (Optional) List of ssl certificates. The `ssl_certificate` block supports fields documented below.

* `waf_configuration` - (Optional) Web Application Firewall configuration settings. The `waf_configuration` block supports fields documented below.

* `disabled_ssl_protocols` - TODO - based on "sslPolicy": {"disabledSslProtocols": []}

The `sku` block supports:

* `name` - (Required) Supported values are:

  * `Standard_Small`
  * `Standard_Medium`
  * `Standard_Large`
  * `WAF_Medium`
  * `WAF_Large`

* `tier` - (Required) Supported values are:

  * `Standard`
  * `WAF`

* `capacity` - (Required) Specifies instance count. Can be 1 to 10.

The `gateway_ip_configuration` block supports:

* `name` - (Required) User defined name of the gateway ip configuration.

* `subnet_id` - (Required) Reference to a Subnet. Application Gateway is deployed in this subnet. No other resource can be deployed in a subnet where Application Gateway is deployed.

The `frontend_port` block supports:

* `name` - (Required) User defined name for frontend Port.

* `port` - (Required) Port number.

The `frontend_ip_configuration` block supports:

* `name` - (Required) User defined name for a frontend IP configuration.

* `subnet_id` - (Optional) Reference to a Subnet.

* `private_ip_address` - (Optional) Private IP Address.

* `public_ip_address_id`- (Optional) Specifies resource Id of a Public Ip Address resource. IPAllocationMethod should be Dynamic.

* `private_ip_address_allocation` - (Optional) Valid values are:
  * `Dynamic`
  * `Static`

The `backend_address_pool` block supports:

* `name` - (Required) User defined name for a backend address pool.

* `ip_address_list` - (Optional) List of public IPAdresses, or internal IP addresses in a backend address pool.

* `fqdn_list` - (Optional) List of FQDNs in a backend address pool.

The `backend_http_settings` block supports:

* `name` - (Required) User defined name for a backend http setting.

* `port` - (Required) Backend port for backend address pool.

* `protocol` - (Required) Valid values are:

  * `Http`
  * `Https`

* `cookie_based_affinity` - (Required) Valid values are:

  * `Enabled`
  * `Disabled`

* `request_timeout` - (Required) RequestTimeout in second. Application Gateway fails the request if response is not received within RequestTimeout. Minimum 1 second and Maximum 86400 secs.

* `probe_name` - (Optional) Reference to URL probe.

* `authentication_certificate` - (Optional) - A list of `authentication_certificate` references for the `backend_http_setting` to use. Each element consists of:

  * `name` (Required)
  * `id` (Calculated)

The `http_listener` block supports:

* `name` - (Required) User defined name for a backend http setting.

* `frontend_ip_configuration_name` - (Required) Reference to frontend Ip configuration.

* `frontend_port_name` - (Required) Reference to frontend port.

* `protocol` - (Required) Valid values are:

  * `Http`
  * `Https`

* `host_name` - (Optional) HostName for `http_listener`. It has to be a valid DNS name.

* `ssl_certificate_name` - (Optional) Reference to ssl certificate. Valid only if protocol is https.

* `require_sni` - (Optional) Applicable only if protocol is https. Enables SNI for multi-hosting.
  Valid values are:
* true
* false

The `probe` block supports:

* `name` - (Required) User defined name for a probe.

* `protocol` - (Required) Protocol used to send probe. Valid values are:

  * `Http`
  * `Https`

* `path` - (Required) Relative path of probe. Valid path starts from '/'. Probe is sent to \{Protocol}://\{host}:\{port}\{path}. The port used will be the same port as defined in the `backend_http_settings`.

* `host` - (Required) Host name to send probe to. If Application Gateway is configured for a single site, by default the Host name should be specified as ‘127.0.0.1’, unless otherwise configured in custom probe.

* `interval` - (Required) Probe interval in seconds. This is the time interval between two consecutive probes. Minimum 1 second and Maximum 86,400 secs.

* `timeout` - (Required) Probe timeout in seconds. Probe marked as failed if valid response is not received with this timeout period. Minimum 1 second and Maximum 86,400 secs.

* `unhealthy_threshold` - (Required) Probe retry count. Backend server is marked down after consecutive probe failure count reaches UnhealthyThreshold. Minimum 1 second and Maximum 20.

The `request_routing_rule` block supports:

* `name` - (Required) User defined name for a request routing rule.

* `rule_type' - (Required) Routing rule type. Valid values are:

  * `Basic`
  * `PathBasedRouting`

* `http_listener_name` - (Required) Reference to `http_listener`.

* `backend_address_pool_name` - (Optional) Reference to `backend_address_pool_name`. Valid for Basic Rule only.

* `backend_http_settings_name` - (Optional) Reference to `backend_http_settings`. Valid for Basic Rule only.

* `url_path_map_name` - (Optional) Reference to `url_path_map`. Valid for PathBasedRouting Rule only.

The `url_path_map` block supports:

* `name` - (Required) User defined name for a url path map.

* `default_backend_address_pool_name` - (Required) Reference to `backend_address_pool_name`.

* `default_backend_http_settings_name` - (Required) Reference to `backend_http_settings`.

* `path_rule` - (Required) List of pathRules. pathRules are order sensitive. Are applied in order they are specified.

The `path_rule` block supports:

* `name` - (Required) User defined name for a path rule.

* `paths` - (Required) The list of path patterns to match. Each must start with / and the only place a \* is allowed is at the end following a /. The string fed to the path matcher does not include any text after the first ? or #, and those chars are not allowed here.

* `backend_address_pool_name` - (Required) Reference to `backend_address_pool_name`.

* `backend_http_settings_name` - (Required) Reference to `backend_http_settings`.

The `authentication_certificate` block supports:

* `name` - (Required) User defined name for an authentication certificate.

* `data` - (Required) Base-64 encoded cer certificate. Only applicable in PUT Request.

The `ssl_certificate` block supports:

* `name` - (Required) User defined name for an SSL certificate.

* `data` - (Required) Base-64 encoded Public cert data corresponding to pfx specified in data. Only applicable in GET request.

* `password` - (Required) Password for the pfx file specified in data. Only applicable in PUT request.

The `waf_configuration` block supports:

* `firewall_mode` - (Required) Firewall mode. Valid values are:

  * `Detection`
  * `Prevention`

* `rule_set_type` - (Required) Rule set type. Must be set to `OWASP`

* `rule_set_version` - (Required) Ruleset version. Supported values:
  * `2.2.9`
  * `3.0`

* `enabled` - (Required) Is the Web Application Firewall enabled?

## Attributes Reference

The following attributes are exported:

* `id` - The application gatewayConfiguration ID.

* `name` - The name of the application gateway.

* `resource_group_name` - The name of the resource group in which to create the application gateway.

* `location` - The location/region where the application gateway is created

## Import

application gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_gateway.testApplicationGateway /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/applicationGateways/myGateway1
```
