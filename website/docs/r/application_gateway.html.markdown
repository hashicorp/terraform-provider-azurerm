---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_gateway"
description: |-
  Manages an Application Gateway.
---

# azurerm_application_gateway

Manages an Application Gateway.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.254.0.0/16"]
}

resource "azurerm_subnet" "frontend" {
  name                 = "frontend"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.254.0.0/24"]
}

resource "azurerm_subnet" "backend" {
  name                 = "backend"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.254.2.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Dynamic"
}

# since these variables are re-used - a locals block makes this more maintainable
locals {
  backend_address_pool_name      = "${azurerm_virtual_network.example.name}-beap"
  frontend_port_name             = "${azurerm_virtual_network.example.name}-feport"
  frontend_ip_configuration_name = "${azurerm_virtual_network.example.name}-feip"
  http_setting_name              = "${azurerm_virtual_network.example.name}-be-htst"
  listener_name                  = "${azurerm_virtual_network.example.name}-httplstn"
  request_routing_rule_name      = "${azurerm_virtual_network.example.name}-rqrt"
  redirect_configuration_name    = "${azurerm_virtual_network.example.name}-rdrcfg"
}

resource "azurerm_application_gateway" "network" {
  name                = "example-appgateway"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "Standard_Small"
    tier     = "Standard"
    capacity = 2
  }

  gateway_ip_configuration {
    name      = "my-gateway-ip-configuration"
    subnet_id = azurerm_subnet.frontend.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 80
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
    public_ip_address_id = azurerm_public_ip.example.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    path                  = "/path1/"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 60
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "Http"
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
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

* `identity` - (Optional) A `identity` block.

* `request_routing_rule` - (Required) One or more `request_routing_rule` blocks as defined below.

* `sku` - (Required) A `sku` block as defined below.

* `zones` - (Optional) A collection of availability zones to spread the Application Gateway over.

-> **Please Note**: Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview).  They are also only supported for [v2 SKUs](https://docs.microsoft.com/en-us/azure/application-gateway/application-gateway-autoscaling-zone-redundant)

---

* `authentication_certificate` - (Optional) One or more `authentication_certificate` blocks as defined below.

* `trusted_root_certificate` - (Optional) One or more `trusted_root_certificate` blocks as defined below.

* `ssl_policy` (Optional) a `ssl policy` block as defined below.

* `enable_http2` - (Optional) Is HTTP2 enabled on the application gateway resource? Defaults to `false`.

* `probe` - (Optional) One or more `probe` blocks as defined below.

* `ssl_certificate` - (Optional) One or more `ssl_certificate` blocks as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `url_path_map` - (Optional) One or more `url_path_map` blocks as defined below.

* `waf_configuration` - (Optional) A `waf_configuration` block as defined below.

* `custom_error_configuration` - (Optional) One or more `custom_error_configuration` blocks as defined below.

* `firewall_policy_id` - (Optional) The ID of the Web Application Firewall Policy.

* `redirect_configuration` - (Optional) A `redirect_configuration` block as defined below.

* `autoscale_configuration` - (Optional) A `autoscale_configuration` block as defined below.

* `rewrite_rule_set` - (Optional) One or more `rewrite_rule_set` blocks as defined below. Only valid for v2 SKUs.

---

A `authentication_certificate` block supports the following:

* `name` - (Required) The Name of the Authentication Certificate to use.

* `data` - (Required) The contents of the Authentication Certificate which should be used.

---

A `trusted_root_certificate` block supports the following:

* `name` - (Required) The Name of the Trusted Root Certificate to use.

* `data` - (Required) The contents of the Trusted Root Certificate which should be used.

---

A `authentication_certificate` block, within the `backend_http_settings` block supports the following:

* `name` - (Required) The name of the Authentication Certificate.

---

A `backend_address_pool` block supports the following:

* `name` - (Required) The name of the Backend Address Pool.

* `fqdns` - (Optional) A list of FQDN's which should be part of the Backend Address Pool.

* `ip_addresses` - (Optional) A list of IP Addresses which should be part of the Backend Address Pool.

---

A `backend_http_settings` block supports the following:

* `cookie_based_affinity` - (Required) Is Cookie-Based Affinity enabled? Possible values are `Enabled` and `Disabled`.

* `affinity_cookie_name` - (Optional) The name of the affinity cookie.

* `name` - (Required) The name of the Backend HTTP Settings Collection.

* `path` - (Optional) The Path which should be used as a prefix for all HTTP requests.

* `port`- (Required) The port which should be used for this Backend HTTP Settings Collection.

* `probe_name` - (Optional) The name of an associated HTTP Probe.

* `protocol`- (Required) The Protocol which should be used. Possible values are `Http` and `Https`.

* `request_timeout` - (Required) The request timeout in seconds, which must be between 1 and 86400 seconds.

* `host_name` - (Optional) Host header to be sent to the backend servers. Cannot be set if `pick_host_name_from_backend_address` is set to `true`.

* `pick_host_name_from_backend_address` - (Optional) Whether host header should be picked from the host name of the backend server. Defaults to `false`.

* `authentication_certificate` - (Optional) One or more `authentication_certificate` blocks.

* `trusted_root_certificate_names` - (Optional) A list of `trusted_root_certificate` names.

* `connection_draining` - (Optional) A `connection_draining` block as defined below.

---

A `connection_draining` block supports the following:

* `enabled` - (Required) If connection draining is enabled or not.

* `drain_timeout_sec` - (Required) The number of seconds connection draining is active. Acceptable values are from `1` second to `3600` seconds.

---


A `frontend_ip_configuration` block supports the following:

* `name` - (Required) The name of the Frontend IP Configuration.

* `subnet_id` - (Optional) The ID of the Subnet.

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

* `subnet_id` - (Required) The ID of the Subnet which the Application Gateway should be connected to.

---

A `http_listener` block supports the following:

* `name` - (Required) The Name of the HTTP Listener.

* `frontend_ip_configuration_name` - (Required) The Name of the Frontend IP Configuration used for this HTTP Listener.

* `frontend_port_name` - (Required) The Name of the Frontend Port use for this HTTP Listener.

* `host_name` - (Optional) The Hostname which should be used for this HTTP Listener. Setting this value changes Listener Type to 'Multi site'.

* `host_names` - (Optional) A list of Hostname(s) should be used for this HTTP Listener. It allows special wildcard characters.

-> **NOTE** The `host_names` and `host_name` are mutually exclusive and cannot both be set.

* `protocol` - (Required) The Protocol to use for this HTTP Listener. Possible values are `Http` and `Https`.

* `require_sni` - (Optional) Should Server Name Indication be Required? Defaults to `false`.

* `ssl_certificate_name` - (Optional) The name of the associated SSL Certificate which should be used for this HTTP Listener.

* `custom_error_configuration` - (Optional) One or more `custom_error_configuration` blocks as defined below.

* `firewall_policy_id` - (Optional) The ID of the Web Application Firewall Policy which should be used as a HTTP Listener.

---

A `identity` block supports the following:

* `type` - (Optional) The Managed Service Identity Type of this Application Gateway. The only possible value is `UserAssigned`. Defaults to `UserAssigned`.

* `identity_ids` - (Required) Specifies a list with a single user managed identity id to be assigned to the Application Gateway.

---

A `match` block supports the following:

* `body` - (Optional) A snippet from the Response Body which must be present in the Response..

* `status_code` - (Optional) A list of allowed status codes for this Health Probe.

---

A `path_rule` block supports the following:

* `name` - (Required) The Name of the Path Rule.

* `paths` - (Required) A list of Paths used in this Path Rule.

* `backend_address_pool_name` - (Optional) The Name of the Backend Address Pool to use for this Path Rule. Cannot be set if `redirect_configuration_name` is set.

* `backend_http_settings_name` - (Optional) The Name of the Backend HTTP Settings Collection to use for this Path Rule. Cannot be set if `redirect_configuration_name` is set.

* `redirect_configuration_name` - (Optional) The Name of a Redirect Configuration to use for this Path Rule. Cannot be set if `backend_address_pool_name` or `backend_http_settings_name` is set.

* `rewrite_rule_set_name` - (Optional) The Name of the Rewrite Rule Set which should be used for this URL Path Map. Only valid for v2 SKUs.

---

A `probe` block support the following:

* `host` - (Optional) The Hostname used for this Probe. If the Application Gateway is configured for a single site, by default the Host name should be specified as ‘127.0.0.1’, unless otherwise configured in custom probe. Cannot be set if `pick_host_name_from_backend_http_settings` is set to `true`.

* `interval` - (Required) The Interval between two consecutive probes in seconds. Possible values range from 1 second to a maximum of 86,400 seconds.

* `name` - (Required) The Name of the Probe.

* `protocol` - (Required) The Protocol used for this Probe. Possible values are `Http` and `Https`.

* `path` - (Required) The Path used for this Probe.

* `timeout` - (Required) The Timeout used for this Probe, which indicates when a probe becomes unhealthy. Possible values range from 1 second to a maximum of 86,400 seconds.

* `unhealthy_threshold` - (Required) The Unhealthy Threshold for this Probe, which indicates the amount of retries which should be attempted before a node is deemed unhealthy. Possible values are from 1 - 20 seconds.

* `port` - (Optional) Custom port which will be used for probing the backend servers. The valid value ranges from 1 to 65535. In case not set, port from http settings will be used. This property is valid for Standard_v2 and WAF_v2 only.

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

* `rewrite_rule_set_name` - (Optional) The Name of the Rewrite Rule Set which should be used for this Routing Rule. Only valid for v2 SKUs.

-> **NOTE:** `backend_address_pool_name`, `backend_http_settings_name`, `redirect_configuration_name`, and `rewrite_rule_set_name` are applicable only when `rule_type` is `Basic`.

* `url_path_map_name` - (Optional) The Name of the URL Path Map which should be associated with this Routing Rule.

---

A `sku` block supports the following:

* `name` - (Required) The Name of the SKU to use for this Application Gateway. Possible values are `Standard_Small`, `Standard_Medium`, `Standard_Large`, `Standard_v2`, `WAF_Medium`, `WAF_Large`, and `WAF_v2`.

* `tier` - (Required) The Tier of the SKU to use for this Application Gateway. Possible values are `Standard`, `Standard_v2`, `WAF` and `WAF_v2`.

* `capacity` - (Required) The Capacity of the SKU to use for this Application Gateway. When using a V1 SKU this value must be between 1 and 32, and 1 to 125 for a V2 SKU. This property is optional if `autoscale_configuration` is set.

---

A `ssl_certificate` block supports the following:

* `name` - (Required) The Name of the SSL certificate that is unique within this Application Gateway

* `data` - (Optional) PFX certificate. Required if `key_vault_secret_id` is not set.

* `password` - (Optional) Password for the pfx file specified in data.  Required if `data` is set.

* `key_vault_secret_id` - (Optional) Secret Id of (base-64 encoded unencrypted pfx) `Secret` or `Certificate` object stored in Azure KeyVault. You need to enable soft delete for keyvault to use this feature. Required if `data` is not set.

-> **NOTE:** TLS termination with Key Vault certificates is limited to the [v2 SKUs](https://docs.microsoft.com/en-us/azure/application-gateway/key-vault-certs).

-> **NOTE:** For TLS termination with Key Vault certificates to work properly existing user-assigned managed identity, which Application Gateway uses to retrieve certificates from Key Vault, should be defined via `identity` block. Additionally, access policies in the Key Vault to allow the identity to be granted *get* access to the secret should be defined.

---

A `url_path_map` block supports the following:

* `name` - (Required) The Name of the URL Path Map.

* `default_backend_address_pool_name` - (Optional) The Name of the Default Backend Address Pool which should be used for this URL Path Map. Cannot be set if `default_redirect_configuration_name` is set.

* `default_backend_http_settings_name` - (Optional) The Name of the Default Backend HTTP Settings Collection which should be used for this URL Path Map. Cannot be set if `default_redirect_configuration_name` is set.

* `default_redirect_configuration_name` - (Optional) The Name of the Default Redirect Configuration which should be used for this URL Path Map. Cannot be set if either `default_backend_address_pool_name` or `default_backend_http_settings_name` is set.

* `default_rewrite_rule_set_name` - (Optional) The Name of the Default Rewrite Rule Set which should be used for this URL Path Map. Only valid for v2 SKUs.

* `path_rule` - (Required) One or more `path_rule` blocks as defined above.

---

A `ssl_policy` block supports the following:

* `disabled_protocols` - (Optional) A list of SSL Protocols which should be disabled on this Application Gateway. Possible values are `TLSv1_0`, `TLSv1_1` and `TLSv1_2`.

~> **NOTE:** `disabled_protocols` cannot be set when `policy_name` or `policy_type` are set.

* `policy_type` - (Optional) The Type of the Policy. Possible values are `Predefined` and `Custom`.

~> **NOTE:** `policy_type` is Required when `policy_name` is set - cannot be set if `disabled_protocols` is set.

When using a `policy_type` of `Predefined` the following fields are supported:

* `policy_name` - (Optional) The Name of the Policy e.g AppGwSslPolicy20170401S. Required if `policy_type` is set to `Predefined`. Possible values can change over time and
are published here https://docs.microsoft.com/en-us/azure/application-gateway/application-gateway-ssl-policy-overview. Not compatible with `disabled_protocols`.

When using a `policy_type` of `Custom` the following fields are supported:

* `cipher_suites` - (Optional) A List of accepted cipher suites. Possible values are: `TLS_DHE_DSS_WITH_AES_128_CBC_SHA`, `TLS_DHE_DSS_WITH_AES_128_CBC_SHA256`, `TLS_DHE_DSS_WITH_AES_256_CBC_SHA`, `TLS_DHE_DSS_WITH_AES_256_CBC_SHA256`, `TLS_DHE_RSA_WITH_AES_128_CBC_SHA`, `TLS_DHE_RSA_WITH_AES_128_GCM_SHA256`, `TLS_DHE_RSA_WITH_AES_256_CBC_SHA`, `TLS_DHE_RSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA`, `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256`, `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA`, `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384`, `TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA`, `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA`, `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384`, `TLS_RSA_WITH_3DES_EDE_CBC_SHA`, `TLS_RSA_WITH_AES_128_CBC_SHA`, `TLS_RSA_WITH_AES_128_CBC_SHA256`, `TLS_RSA_WITH_AES_128_GCM_SHA256`, `TLS_RSA_WITH_AES_256_CBC_SHA`, `TLS_RSA_WITH_AES_256_CBC_SHA256` and `TLS_RSA_WITH_AES_256_GCM_SHA384`.

* `min_protocol_version` - (Optional) The minimal TLS version. Possible values are `TLSv1_0`, `TLSv1_1` and `TLSv1_2`.



---

A `waf_configuration` block supports the following:

* `enabled` - (Required) Is the Web Application Firewall be enabled?

* `firewall_mode` - (Required) The Web Application Firewall Mode. Possible values are `Detection` and `Prevention`.

* `rule_set_type` - (Required) The Type of the Rule Set used for this Web Application Firewall. Currently, only `OWASP` is supported.

* `rule_set_version` - (Required) The Version of the Rule Set used for this Web Application Firewall. Possible values are `2.2.9`, `3.0`, and `3.1`.

* `disabled_rule_group` - (Optional) one or more `disabled_rule_group` blocks as defined below.

* `file_upload_limit_mb` - (Optional) The File Upload Limit in MB. Accepted values are in the range `1`MB to `750`MB for the `WAF_v2` SKU, and `1`MB to `500`MB for all other SKUs. Defaults to `100`MB.

* `request_body_check` - (Optional) Is Request Body Inspection enabled?  Defaults to `true`.

* `max_request_body_size_kb` - (Optional) The Maximum Request Body Size in KB.  Accepted values are in the range `1`KB to `128`KB.  Defaults to `128`KB.

* `exclusion` - (Optional) one or more `exclusion` blocks as defined below.

---

A `disabled_rule_group` block supports the following:

* `rule_group_name` - (Required) The rule group where specific rules should be disabled. Accepted values are:  `crs_20_protocol_violations`, `crs_21_protocol_anomalies`, `crs_23_request_limits`, `crs_30_http_policy`, `crs_35_bad_robots`, `crs_40_generic_attacks`, `crs_41_sql_injection_attacks`, `crs_41_xss_attacks`, `crs_42_tight_security`, `crs_45_trojans`, `General`, `REQUEST-911-METHOD-ENFORCEMENT`, `REQUEST-913-SCANNER-DETECTION`, `REQUEST-920-PROTOCOL-ENFORCEMENT`, `REQUEST-921-PROTOCOL-ATTACK`, `REQUEST-930-APPLICATION-ATTACK-LFI`, `REQUEST-931-APPLICATION-ATTACK-RFI`, `REQUEST-932-APPLICATION-ATTACK-RCE`, `REQUEST-933-APPLICATION-ATTACK-PHP`, `REQUEST-941-APPLICATION-ATTACK-XSS`, `REQUEST-942-APPLICATION-ATTACK-SQLI`, `REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION`

* `rules` - (Optional) A list of rules which should be disabled in that group. Disables all rules in the specified group if `rules` is not specified.

---

A `exclusion` block supports the following:

* `match_variable` - (Required) Match variable of the exclusion rule to exclude header, cookie or GET arguments. Possible values are `RequestHeaderNames`, `RequestArgNames` and `RequestCookieNames`

* `selector_match_operator` - (Optional) Operator which will be used to search in the variable content. Possible values are `Equals`, `StartsWith`, `EndsWith`, `Contains`. If empty will exclude all traffic on this `match_variable`

* `selector` - (Optional) String value which will be used for the filter operation. If empty will exclude all traffic on this `match_variable`

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

---

A `autoscale_configuration` block supports the following:

* `min_capacity` - (Required) Minimum capacity for autoscaling. Accepted values are in the range `0` to `100`.

* `max_capacity` - (Optional) Maximum capacity for autoscaling. Accepted values are in the range `2` to `125`.

---

A `rewrite_rule_set` block supports the following:

* `name` - (Required) Unique name of the rewrite rule set block

* `rewrite_rule` - (Required) One or more `rewrite_rule` blocks as defined above.

---

A `rewrite_rule` block supports the following:

* `name` - (Required) Unique name of the rewrite rule block

* `rule_sequence` - (Required) Rule sequence of the rewrite rule that determines the order of execution in a set.

* `condition` - (Optional) One or more `condition` blocks as defined above.

* `request_header_configuration` - (Optional) One or more `request_header_configuration` blocks as defined above.

* `response_header_configuration` - (Optional) One or more `response_header_configuration` blocks as defined above.

---

A `condition` block supports the following:

* `variable` - (Required) The [variable](https://docs.microsoft.com/en-us/azure/application-gateway/rewrite-http-headers#server-variables) of the condition.

* `pattern` - (Required) The pattern, either fixed string or regular expression, that evaluates the truthfulness of the condition.

* `ignore_case` - (Optional) Perform a case in-sensitive comparison. Defaults to `false`

* `negate` - (Optional) Negate the result of the condition evaluation. Defaults to `false`

---

A `request_header_configuration` block supports the following:

* `header_name` - (Required) Header name of the header configuration.

* `header_value` - (Required) Header value of the header configuration. To delete a request header set this property to an empty string.

---

A `response_header_configuration` block supports the following:

* `header_name` - (Required) Header name of the header configuration.

* `header_value` - (Required) Header value of the header configuration. To delete a response header set this property to an empty string.

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

* `rewrite_rule_set_id` - The ID of the Rewrite Rule Set used in this Path Rule.

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

* `rewrite_rule_set_id` - The ID of the associated Rewrite Rule Set.

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

---

A `rewrite_rule_set` block exports the following:

* `id` - The ID of the Rewrite Rule Set

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Application Gateway.
* `update` - (Defaults to 90 minutes) Used when updating the Application Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Gateway.
* `delete` - (Defaults to 90 minutes) Used when deleting the Application Gateway.

## Import

Application Gateway's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/applicationGateways/myGateway1
```
