---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_application_gateway"
description: |-
  Gets information about an existing Application Gateway.
---

# Data Source: azurerm_application_gateway

Use this data source to access information about an existing Application Gateway.

## Example Usage

```hcl
data "azurerm_application_gateway" "example" {
  name                = "existing-app-gateway"
  resource_group_name = "existing-resources"
}

output "id" {
  value = data.azurerm_application_gateway.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of this Application Gateway.

* `resource_group_name` - The name of the Resource Group where the Application Gateway exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Application Gateway.

* `location` - The Azure region where the Application Gateway should exist.

* `backend_address_pool` - One or more `backend_address_pool` blocks as defined below.

* `backend_http_settings` - One or more `backend_http_settings` blocks as defined below.

* `frontend_ip_configuration` - One or more `frontend_ip_configuration` blocks as defined below.

* `frontend_port` - One or more `frontend_port` blocks as defined below.

* `gateway_ip_configuration` - One or more `gateway_ip_configuration` blocks as defined below.

* `http_listener` - One or more `http_listener` blocks as defined below.

* `request_routing_rule` - One or more `request_routing_rule` blocks as defined below.

* `sku` - A `sku` block as defined below.

---

* `fips_enabled` - Is FIPS enabled on the Application Gateway?

* `global` - A `global` block as defined below.

* `identity` - An `identity` block as defined below.

* `private_link_configuration` - One or more `private_link_configuration` blocks as defined below.

* `zones` - Specifies a list of Availability Zones in which this Application Gateway should be located.

* `trusted_client_certificate` - One or more `trusted_client_certificate` blocks as defined below.

* `ssl_profile` - One or more `ssl_profile` blocks as defined below.

* `authentication_certificate` - One or more `authentication_certificate` blocks as defined below.

* `trusted_root_certificate` - One or more `trusted_root_certificate` blocks as defined below.

* `ssl_policy` - a `ssl_policy` block as defined below.

* `enable_http2` - Is HTTP2 enabled on the application gateway resource?

* `force_firewall_policy_association` - Is the Firewall Policy associated with the Application Gateway?

* `probe` - One or more `probe` blocks as defined below.

* `ssl_certificate` - One or more `ssl_certificate` blocks as defined below.

* `tags` - A mapping of tags to assign to the resource.

* `url_path_map` - One or more `url_path_map` blocks as defined below.

* `waf_configuration` - A `waf_configuration` block as defined below.

* `custom_error_configuration` - One or more `custom_error_configuration` blocks as defined below.

* `firewall_policy_id` - The ID of the Web Application Firewall Policy.

* `redirect_configuration` - One or more `redirect_configuration` blocks as defined below.

* `autoscale_configuration` - A `autoscale_configuration` block as defined below.

* `rewrite_rule_set` - One or more `rewrite_rule_set` blocks as defined below. Only valid for v2 SKUs.

---

A `authentication_certificate` block supports the following:

* `id` - The ID of the Authentication Certificate.

* `name` - The Name of the Authentication Certificate to use.

---

A `trusted_root_certificate` block supports the following:

* `name` - The Name of the Trusted Root Certificate to use.

* `key_vault_secret_id` - The Secret ID of (base-64 encoded unencrypted pfx) `Secret` or `Certificate` object stored in Azure KeyVault. You need to enable soft delete for the Key Vault to use this feature.

---

A `authentication_certificate` block, within the `backend_http_settings` block supports the following:

* `id` - The ID of the Authentication Certificate.

* `name` - The name of the Authentication Certificate.

---

A `backend_address_pool` block supports the following:

* `id` - The ID of the Backend Address Pool.

* `name` - The name of the Backend Address Pool.

* `fqdns` - A list of FQDN's which should be part of the Backend Address Pool.

* `ip_addresses` - A list of IP Addresses which should be part of the Backend Address Pool.

---

A `backend_http_settings` block supports the following:

* `id` - The ID of the Backend HTTP Settings Configuration.

* `probe_id` - The ID of the associated Probe.

* `cookie_based_affinity` - Is Cookie-Based Affinity enabled? Possible values are `Enabled` and `Disabled`.

* `affinity_cookie_name` - The name of the affinity cookie.

* `name` - The name of the Backend HTTP Settings Collection.

* `path` - The Path which should be used as a prefix for all HTTP requests.

* `port` - The port which should be used for this Backend HTTP Settings Collection.

* `probe_name` - The name of an associated HTTP Probe.

* `protocol` - The Protocol which should be used. Possible values are `Http` and `Https`.

* `request_timeout` - The request timeout in seconds, which must be between 1 and 86400 seconds.

* `host_name` - Host header to be sent to the backend servers. Cannot be set if `pick_host_name_from_backend_address` is set to `true`.

* `pick_host_name_from_backend_address` - Whether host header should be picked from the host name of the backend server.

* `authentication_certificate` - One or more `authentication_certificate` blocks as defined below.

* `trusted_root_certificate_names` - A list of `trusted_root_certificate` names.

* `connection_draining` - A `connection_draining` block as defined below.

---

A `connection_draining` block supports the following:

* `enabled` - If connection draining is enabled or not.

* `drain_timeout_sec` - The number of seconds connection draining is active. Acceptable values are from `1` second to `3600` seconds.

---

A `frontend_ip_configuration` block supports the following:

* `id` - The ID of the Frontend IP Configuration.

* `private_link_configuration_id` - The ID of the associated private link configuration.

* `name` - The name of the Frontend IP Configuration.

* `subnet_id` - The ID of the Subnet.

* `private_ip_address` - The Private IP Address to use for the Application Gateway.

* `public_ip_address_id` - The ID of a Public IP Address which the Application Gateway should use. The allocation method for the Public IP Address depends on the `sku` of this Application Gateway. Please refer to the [Azure documentation for public IP addresses](https://docs.microsoft.com/azure/virtual-network/public-ip-addresses#application-gateways) for details.

* `private_ip_address_allocation` - The Allocation Method for the Private IP Address. Possible values are `Dynamic` and `Static`.

* `private_link_configuration_name` - The name of the private link configuration to use for this frontend IP configuration.

---

A `frontend_port` block supports the following:

* `id` - The ID of the Frontend Port.

* `name` - The name of the Frontend Port.

* `port` - The port used for this Frontend Port.

---

A `gateway_ip_configuration` block supports the following:

* `id` - The ID of the Gateway IP Configuration.

* `name` - The Name of this Gateway IP Configuration.

* `subnet_id` - The ID of the Subnet which the Application Gateway should be connected to.

---

A `http_listener` block supports the following:

* `id` - The ID of the HTTP Listener.

* `frontend_ip_configuration_id` - The ID of the associated Frontend Configuration.

* `frontend_port_id` - The ID of the associated Frontend Port.

* `ssl_certificate_id` - The ID of the associated SSL Certificate.

* `ssl_profile_id` - The ID of the associated SSL Profile.

* `name` - The Name of the HTTP Listener.

* `frontend_ip_configuration_name` - The Name of the Frontend IP Configuration used for this HTTP Listener.

* `frontend_port_name` - The Name of the Frontend Port use for this HTTP Listener.

* `host_name` - The Hostname which should be used for this HTTP Listener.

* `host_names` - A list of Hostname(s) should be used for this HTTP Listener. It allows special wildcard characters.

* `protocol` - The Protocol to use for this HTTP Listener. Possible values are `Http` and `Https`.

* `require_sni` - Should Server Name Indication be Required?

* `ssl_certificate_name` - The name of the associated SSL Certificate which should be used for this HTTP Listener.

* `custom_error_configuration` - One or more `custom_error_configuration` blocks as defined below.

* `firewall_policy_id` - The ID of the Web Application Firewall Policy which should be used for this HTTP Listener.

* `ssl_profile_name` - The name of the associated SSL Profile which should be used for this HTTP Listener.

---

An `identity` block supports the following:

* `type` - Specifies the type of Managed Service Identity that should be configured on this Application Gateway. Only possible value is `UserAssigned`.

* `identity_ids` - Specifies a list of User Assigned Managed Identity IDs to be assigned to this Application Gateway.

---

A `private_endpoint_connection` block exports the following:

* `name` - The name of the private endpoint connection.

* `id` - The ID of the private endpoint connection.

---

A `private_link_configuration` block supports the following:

* `id` - The ID of the private link configuration.

* `name` - The name of the private link configuration.

* `ip_configuration` - One or more `ip_configuration` blocks as defined below.

---

An `ip_configuration` block supports the following:

* `name` - The name of the IP configuration.

* `subnet_id` - The ID of the subnet the private link configuration should connect to.

* `private_ip_address_allocation` - The allocation method used for the Private IP Address. Possible values are `Dynamic` and `Static`.

* `primary` - Is this the Primary IP Configuration?

* `private_ip_address` - The Static IP Address which should be used.

---

A `match` block supports the following:

* `body` - A snippet from the Response Body which must be present in the Response.

* `status_code` - A list of allowed status codes for this Health Probe.

---

A `path_rule` block supports the following:

* `id` - The ID of the Path Rule.

* `backend_address_pool_id` - The ID of the Backend Address Pool used in this Path Rule.

* `backend_http_settings_id` - The ID of the Backend HTTP Settings Collection used in this Path Rule.

* `redirect_configuration_id` - The ID of the Redirect Configuration used in this Path Rule.

* `rewrite_rule_set_id` - The ID of the Rewrite Rule Set used in this Path Rule.

* `name` - The Name of the Path Rule.

* `paths` - A list of Paths used in this Path Rule.

* `backend_address_pool_name` - The Name of the Backend Address Pool to use for this Path Rule. Cannot be set if `redirect_configuration_name` is set.

* `backend_http_settings_name` - The Name of the Backend HTTP Settings Collection to use for this Path Rule. Cannot be set if `redirect_configuration_name` is set.

* `redirect_configuration_name` - The Name of a Redirect Configuration to use for this Path Rule. Cannot be set if `backend_address_pool_name` or `backend_http_settings_name` is set.

* `rewrite_rule_set_name` - The Name of the Rewrite Rule Set which should be used for this URL Path Map. Only valid for v2 SKUs.

* `firewall_policy_id` - The ID of the Web Application Firewall Policy which should be used as an HTTP Listener.

---

A `probe` block support the following:

* `id` - The ID of the Probe.

* `host` - The Hostname used for this Probe. If the Application Gateway is configured for a single site, by default the Host name should be specified as `127.0.0.1`, unless otherwise configured in custom probe. Cannot be set if `pick_host_name_from_backend_http_settings` is set to `true`.

* `interval` - The Interval between two consecutive probes in seconds. Possible values range from 1 second to a maximum of 86,400 seconds.

* `name` - The Name of the Probe.

* `protocol` - The Protocol used for this Probe. Possible values are `Http` and `Https`.

* `path` - The Path used for this Probe.

* `timeout` - The Timeout used for this Probe, which indicates when a probe becomes unhealthy. Possible values range from 1 second to a maximum of 86,400 seconds.

* `unhealthy_threshold` - The Unhealthy Threshold for this Probe, which indicates the amount of retries which should be attempted before a node is deemed unhealthy. Possible values are from 1 to 20.

* `port` - Custom port which will be used for probing the backend servers. The valid value ranges from 1 to 65535. In case not set, port from HTTP settings will be used. This property is valid for Standard_v2 and WAF_v2 only.

* `pick_host_name_from_backend_http_settings` - Whether the host header should be picked from the backend HTTP settings.

* `match` - A `match` block as defined above.

* `minimum_servers` - The minimum number of servers that are always marked as healthy.

---

A `request_routing_rule` block supports the following:

* `id` - The ID of the Request Routing Rule.

* `http_listener_id` - The ID of the associated HTTP Listener.

* `backend_address_pool_id` - The ID of the associated Backend Address Pool.

* `backend_http_settings_id` - The ID of the associated Backend HTTP Settings Configuration.

* `redirect_configuration_id` - The ID of the associated Redirect Configuration.

* `rewrite_rule_set_id` - The ID of the associated Rewrite Rule Set.

* `url_path_map_id` - The ID of the associated URL Path Map.

* `name` - The Name of this Request Routing Rule.

* `rule_type` - The Type of Routing that should be used for this Rule. Possible values are `Basic` and `PathBasedRouting`.

* `http_listener_name` - The Name of the HTTP Listener which should be used for this Routing Rule.

* `backend_address_pool_name` - The Name of the Backend Address Pool which should be used for this Routing Rule. Cannot be set if `redirect_configuration_name` is set.

* `backend_http_settings_name` - The Name of the Backend HTTP Settings Collection which should be used for this Routing Rule. Cannot be set if `redirect_configuration_name` is set.

* `redirect_configuration_name` - The Name of the Redirect Configuration which should be used for this Routing Rule. Cannot be set if either `backend_address_pool_name` or `backend_http_settings_name` is set.

* `rewrite_rule_set_name` - The Name of the Rewrite Rule Set which should be used for this Routing Rule. Only valid for v2 SKUs.

* `url_path_map_name` - The Name of the URL Path Map which should be associated with this Routing Rule.

* `priority` - Rule evaluation order can be dictated by specifying an integer value from `1` to `20000` with `1` being the highest priority and `20000` being the lowest priority.

---

A `global` block supports the following:

* `request_buffering_enabled` - Whether Application Gateway's Request buffer is enabled.

* `response_buffering_enabled` - Whether Application Gateway's Response buffer is enabled.

---

A `sku` block supports the following:

* `name` - The Name of the SKU to use for this Application Gateway. Possible values are `Standard_Small`, `Standard_Medium`, `Standard_Large`, `Standard_v2`, `WAF_Medium`, `WAF_Large`, and `WAF_v2`.

* `tier` - The Tier of the SKU to use for this Application Gateway. Possible values are `Standard`, `Standard_v2`, `WAF` and `WAF_v2`.

* `capacity` - The Capacity of the SKU to use for this Application Gateway. When using a V1 SKU this value must be between 1 and 32, and 1 to 125 for a V2 SKU. This property is optional if `autoscale_configuration` is set.

---

A `ssl_certificate` block supports the following:

* `id` - The ID of the SSL Certificate.

* `public_cert_data` - The Public Certificate Data associated with the SSL Certificate.

* `name` - The Name of the SSL certificate that is unique within this Application Gateway

* `key_vault_secret_id` - Secret ID of (base-64 encoded unencrypted pfx) `Secret` or `Certificate` object stored in Azure KeyVault. You need to enable soft delete for keyvault to use this feature.

---

A `url_path_map` block supports the following:

* `id` - The ID of the URL Path Map.

* `default_backend_address_pool_id` - The ID of the Default Backend Address Pool.

* `default_backend_http_settings_id` - The ID of the Default Backend HTTP Settings Collection.

* `default_redirect_configuration_id` - The ID of the Default Redirect Configuration.

* `path_rule` - A list of `path_rule` blocks as defined above.

* `name` - The Name of the URL Path Map.

* `default_backend_address_pool_name` - The Name of the Default Backend Address Pool which should be used for this URL Path Map. Cannot be set if `default_redirect_configuration_name` is set.

* `default_backend_http_settings_name` - The Name of the Default Backend HTTP Settings Collection which should be used for this URL Path Map. Cannot be set if `default_redirect_configuration_name` is set.

* `default_redirect_configuration_name` - The Name of the Default Redirect Configuration which should be used for this URL Path Map. Cannot be set if either `default_backend_address_pool_name` or `default_backend_http_settings_name` is set.

* `default_rewrite_rule_set_name` - The Name of the Default Rewrite Rule Set which should be used for this URL Path Map. Only valid for v2 SKUs.

* `path_rule` - One or more `path_rule` blocks as defined above.

---

A `trusted_client_certificate` block supports the following:

* `name` - The name of the Trusted Client Certificate that is unique within this Application Gateway.

---

A `ssl_profile` block supports the following:

* `name` - The name of the SSL Profile that is unique within this Application Gateway.

* `trusted_client_certificate_names` - The name of the Trusted Client Certificate that will be used to authenticate requests from clients.

* `verify_client_cert_issuer_dn` - Should client certificate issuer DN be verified?

* `verify_client_certificate_revocation` - Specify the method to check client certificate revocation status. Possible value is `OCSP`.

* `ssl_policy` - a `ssl_policy` block as defined below.

---

A `ssl_policy` block supports the following:

* `disabled_protocols` - A list of SSL Protocols which should be disabled on this Application Gateway. Possible values are `TLSv1_0`, `TLSv1_1`, `TLSv1_2` and `TLSv1_3`.

* `policy_type` - The Type of the Policy. Possible values are `Predefined`, `Custom` and `CustomV2`.

* `policy_name` - The Name of the Policy e.g. AppGwSslPolicy20170401S. Possible values can change over time and are published here <https://docs.microsoft.com/azure/application-gateway/application-gateway-ssl-policy-overview>. Not compatible with `disabled_protocols`.

* `cipher_suites` - A List of accepted cipher suites. Possible values are: `TLS_DHE_DSS_WITH_3DES_EDE_CBC_SHA`, `TLS_DHE_DSS_WITH_AES_128_CBC_SHA`, `TLS_DHE_DSS_WITH_AES_128_CBC_SHA256`, `TLS_DHE_DSS_WITH_AES_256_CBC_SHA`, `TLS_DHE_DSS_WITH_AES_256_CBC_SHA256`, `TLS_DHE_RSA_WITH_AES_128_CBC_SHA`, `TLS_DHE_RSA_WITH_AES_128_GCM_SHA256`, `TLS_DHE_RSA_WITH_AES_256_CBC_SHA`, `TLS_DHE_RSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA`, `TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256`, `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA`, `TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384`, `TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384`, `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA`, `TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256`, `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`, `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA`, `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384`, `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384`, `TLS_RSA_WITH_3DES_EDE_CBC_SHA`, `TLS_RSA_WITH_AES_128_CBC_SHA`, `TLS_RSA_WITH_AES_128_CBC_SHA256`, `TLS_RSA_WITH_AES_128_GCM_SHA256`, `TLS_RSA_WITH_AES_256_CBC_SHA`, `TLS_RSA_WITH_AES_256_CBC_SHA256` and `TLS_RSA_WITH_AES_256_GCM_SHA384`.

* `min_protocol_version` - The minimal TLS version. Possible values are `TLSv1_0`, `TLSv1_1`, `TLSv1_2` and `TLSv1_3`.

---

A `waf_configuration` block supports the following:

* `enabled` - Is the Web Application Firewall enabled?

* `firewall_mode` - The Web Application Firewall Mode. Possible values are `Detection` and `Prevention`.

* `rule_set_type` - The Type of the Rule Set used for this Web Application Firewall. Possible values are `OWASP` and `Microsoft_BotManagerRuleSet`.

* `rule_set_version` - The Version of the Rule Set used for this Web Application Firewall. Possible values are `0.1`, `1.0`, `2.2.9`, `3.0`, `3.1` and `3.2`.

* `disabled_rule_group` - one or more `disabled_rule_group` blocks as defined below.

* `file_upload_limit_mb` - The File Upload Limit in MB. Accepted values are in the range `1`MB to `750`MB for the `WAF_v2` SKU, and `1`MB to `500`MB for all other SKUs.

* `request_body_check` - Is Request Body Inspection enabled?

* `max_request_body_size_kb` - The Maximum Request Body Size in KB. Accepted values are in the range `1`KB to `128`KB.

* `exclusion` - one or more `exclusion` blocks as defined below.

---

A `disabled_rule_group` block supports the following:

* `rule_group_name` - The rule group where specific rules should be disabled. Possible values are `BadBots`, `crs_20_protocol_violations`, `crs_21_protocol_anomalies`, `crs_23_request_limits`, `crs_30_http_policy`, `crs_35_bad_robots`, `crs_40_generic_attacks`, `crs_41_sql_injection_attacks`, `crs_41_xss_attacks`, `crs_42_tight_security`, `crs_45_trojans`, `General`, `GoodBots`, `Known-CVEs`, `REQUEST-911-METHOD-ENFORCEMENT`, `REQUEST-913-SCANNER-DETECTION`, `REQUEST-920-PROTOCOL-ENFORCEMENT`, `REQUEST-921-PROTOCOL-ATTACK`, `REQUEST-930-APPLICATION-ATTACK-LFI`, `REQUEST-931-APPLICATION-ATTACK-RFI`, `REQUEST-932-APPLICATION-ATTACK-RCE`, `REQUEST-933-APPLICATION-ATTACK-PHP`, `REQUEST-941-APPLICATION-ATTACK-XSS`, `REQUEST-942-APPLICATION-ATTACK-SQLI`, `REQUEST-943-APPLICATION-ATTACK-SESSION-FIXATION`, `REQUEST-944-APPLICATION-ATTACK-JAVA` and `UnknownBots`.

* `rules` - A list of rules which should be disabled in that group. Disables all rules in the specified group if `rules` is not specified.

---

A `exclusion` block supports the following:

* `match_variable` - Match variable of the exclusion rule to exclude header, cookie or GET arguments. Possible values are `RequestArgKeys`, `RequestArgNames`, `RequestArgValues`, `RequestCookieKeys`, `RequestCookieNames`, `RequestCookieValues`, `RequestHeaderKeys`, `RequestHeaderNames` and `RequestHeaderValues`

* `selector_match_operator` - Operator which will be used to search in the variable content. Possible values are `Contains`, `EndsWith`, `Equals`, `EqualsAny` and `StartsWith`. If empty will exclude all traffic on this `match_variable`

* `selector` - String value which will be used for the filter operation. If empty will exclude all traffic on this `match_variable`

---

A `custom_error_configuration` block supports the following:

* `id` - The ID of the Custom Error Configuration.

* `status_code` - Status code of the application gateway customer error. Possible values are `HttpStatus403` and `HttpStatus502`

* `custom_error_page_url` - Error page URL of the application gateway customer error.

---

A `redirect_configuration` block supports the following:

* `id` - The ID of the Redirect Configuration.

* `name` - Unique name of the redirect configuration block

* `redirect_type` - The type of redirect. Possible values are `Permanent`, `Temporary`, `Found` and `SeeOther`

* `target_listener_name` - The name of the listener to redirect to. Cannot be set if `target_url` is set.

* `target_url` - The Url to redirect the request to. Cannot be set if `target_listener_name` is set.

* `include_path` - Whether to include the path in the redirected Url.

* `include_query_string` - Whether to include the query string in the redirected Url. Default to `false`

---

A `autoscale_configuration` block supports the following:

* `min_capacity` - Minimum capacity for autoscaling. Accepted values are in the range `0` to `100`.

* `max_capacity` - Maximum capacity for autoscaling. Accepted values are in the range `2` to `125`.

---

A `rewrite_rule_set` block supports the following:

* `id` - The ID of the Rewrite Rule Set

* `name` - Unique name of the rewrite rule set block

* `rewrite_rule` - One or more `rewrite_rule` blocks as defined above.

---

A `rewrite_rule` block supports the following:

* `name` - Unique name of the rewrite rule block

* `rule_sequence` - Rule sequence of the rewrite rule that determines the order of execution in a set.

* `condition` - One or more `condition` blocks as defined above.

* `request_header_configuration` - One or more `request_header_configuration` blocks as defined above.

* `response_header_configuration` - One or more `response_header_configuration` blocks as defined above.

* `url` - One `url` block as defined below

---

A `condition` block supports the following:

* `variable` - The [variable](https://docs.microsoft.com/azure/application-gateway/rewrite-http-headers#server-variables) of the condition.

* `pattern` - The pattern, either fixed string or regular expression, that evaluates the truthfulness of the condition.

* `ignore_case` - Perform a case in-sensitive comparison.

* `negate` - Negate the result of the condition evaluation.

---

A `request_header_configuration` block supports the following:

* `header_name` - Header name of the header configuration.

* `header_value` - Header value of the header configuration. To delete a request header set this property to an empty string.

---

A `response_header_configuration` block supports the following:

* `header_name` - Header name of the header configuration.

* `header_value` - Header value of the header configuration. To delete a response header set this property to an empty string.

---

A `url` block supports the following:

* `path` - The URL path to rewrite.

* `query_string` - The query string to rewrite.

* `components` - The components used to rewrite the URL. Possible values are `path_only` and `query_string_only` to limit the rewrite to the URL Path or URL Query String only.

* `reroute` - Whether the URL path map should be reevaluated after this rewrite has been applied. [More info on rewrite configutation](https://docs.microsoft.com/azure/application-gateway/rewrite-http-headers-url#rewrite-configuration)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Application Gateway.
