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

* `location` - The Azure region where the Application Gateway exists.

* `backend_address_pool` - One or more `backend_address_pool` blocks as defined below.

* `backend_http_settings` - One or more `backend_http_settings` blocks as defined below.

* `frontend_ip_configuration` - One or more `frontend_ip_configuration` blocks as defined below.

* `frontend_port` - One or more `frontend_port` blocks as defined below.

* `gateway_ip_configuration` - One or more `gateway_ip_configuration` blocks as defined below.

* `http_listener` - One or more `http_listener` blocks as defined below.

* `request_routing_rule` - One or more `request_routing_rule` blocks as defined below.

* `sku` - A `sku` block as defined below.

* `fips_enabled` - Is FIPS enabled on the Application Gateway?

* `global` - A `global` block as defined below.

* `identity` - An `identity` block as defined below.

* `private_link_configuration` - One or more `private_link_configuration` blocks as defined below.

* `zones` - The list of Availability Zones in which this Application Gateway can use.

* `trusted_client_certificate` - One or more `trusted_client_certificate` blocks as defined below.

* `ssl_profile` - One or more `ssl_profile` blocks as defined below.

* `authentication_certificate` - One or more `authentication_certificate` blocks as defined below.

* `trusted_root_certificate` - One or more `trusted_root_certificate` blocks as defined below.

* `ssl_policy` - An `ssl_policy` block as defined below.

* `http2_enabled` - Is HTTP2 enabled on the application gateway resource?

* `force_firewall_policy_association` - Is the Firewall Policy associated with the Application Gateway?

* `probe` - One or more `probe` blocks as defined below.

* `ssl_certificate` - One or more `ssl_certificate` blocks as defined below.

* `tags` - A mapping of tags to assign to the resource.

* `url_path_map` - One or more `url_path_map` blocks as defined below.

* `waf_configuration` - A `waf_configuration` block as defined below.

* `custom_error_configuration` - One or more `custom_error_configuration` blocks as defined below.

* `firewall_policy_id` - The ID of the Web Application Firewall Policy.

* `redirect_configuration` - One or more `redirect_configuration` blocks as defined below.

* `autoscale_configuration` - An `autoscale_configuration` block as defined below.

* `rewrite_rule_set` - One or more `rewrite_rule_set` blocks as defined below.

---

An `authentication_certificate` block exports the following:

* `id` - The ID of the Authentication Certificate.

* `name` - The Name of the Authentication Certificate in use.

---

A `trusted_root_certificate` block exports the following:

* `id` - The ID of the Trusted Root Certificate in use.

* `name` - The Name of the Trusted Root Certificate in use.

* `key_vault_secret_id` - The Secret ID of (base-64 encoded unencrypted pfx) `Secret` or `Certificate` object stored in Azure KeyVault.

---

A `authentication_certificate` block, within the `backend_http_settings` block exports the following:

* `id` - The ID of the Authentication Certificate.

* `name` - The name of the Authentication Certificate.

---

A `backend_address_pool` block exports the following:

* `id` - The ID of the Backend Address Pool.

* `name` - The name of the Backend Address Pool.

* `fqdns` - A list of FQDNs which are part of the Backend Address Pool.

* `ip_addresses` - A list of IP Addresses which are part of the Backend Address Pool.

---

A `backend_http_settings` block exports the following:

* `id` - The ID of the Backend HTTP Settings Configuration.

* `probe_id` - The ID of the associated Probe.

* `cookie_based_affinity` - Is Cookie-Based Affinity enabled?

* `affinity_cookie_name` - The name of the affinity cookie.

* `name` - The name of the Backend HTTP Settings Collection.

* `path` - The path which is used as a prefix for all HTTP requests.

* `port` - The port which is used for this Backend HTTP Settings Collection.

* `probe_name` - The name of the associated HTTP Probe.

* `protocol` - The protocol which will be used.

* `request_timeout` - The request timeout in seconds.

* `host_name` - Host header to be sent to the backend servers.

* `pick_host_name_from_backend_address` - Whether host header will be picked from the host name of the backend server.

* `authentication_certificate` - One or more `authentication_certificate` blocks as defined below.

* `trusted_root_certificate_names` - A list of `trusted_root_certificate` names.

* `connection_draining` - A `connection_draining` block as defined below.

---

A `connection_draining` block exports the following:

* `enabled` - If connection draining is enabled or not.

* `drain_timeout_sec` - The number of seconds connection draining is active.

---

A `frontend_ip_configuration` block exports the following:

* `id` - The ID of the Frontend IP Configuration.

* `private_link_configuration_id` - The ID of the associated Private Link configuration.

* `name` - The name of the Frontend IP Configuration.

* `subnet_id` - The ID of the Subnet.

* `private_ip_address` - The Private IP Address to use for the Application Gateway.

* `public_ip_address_id` - The ID of the Public IP Address which the Application Gateway will use.

* `private_ip_address_allocation` - The Allocation Method for the Private IP Address.

* `private_link_configuration_name` - The name of the Private Link configuration in use by this Frontend IP Configuration.

---

A `frontend_port` block exports the following:

* `id` - The ID of the Frontend Port.

* `name` - The name of the Frontend Port.

* `port` - The port used for this Frontend Port.

---

A `gateway_ip_configuration` block exports the following:

* `id` - The ID of the Gateway IP Configuration.

* `name` - The Name of this Gateway IP Configuration.

* `subnet_id` - The ID of the Subnet which the Application Gateway is connected to.

---

A `http_listener` block exports the following:

* `id` - The ID of the HTTP Listener.

* `frontend_ip_configuration_id` - The ID of the associated Frontend Configuration.

* `frontend_port_id` - The ID of the associated Frontend Port.

* `ssl_certificate_id` - The ID of the associated SSL Certificate.

* `ssl_profile_id` - The ID of the associated SSL Profile.

* `name` - The Name of the HTTP Listener.

* `frontend_ip_configuration_name` - The Name of the Frontend IP Configuration used for this HTTP Listener.

* `frontend_port_name` - The Name of the Frontend Port used for this HTTP Listener.

* `host_name` - The Hostname which is used for this HTTP Listener.

* `host_names` - A list of Hostname(s) used for this HTTP Listener. It allows special wildcard characters.

* `protocol` - The Protocol to use for this HTTP Listener.

* `require_sni` - Is Server Name Indication required?

* `ssl_certificate_name` - The name of the associated SSL Certificate which is used for this HTTP Listener.

* `custom_error_configuration` - One or more `custom_error_configuration` blocks as defined below.

* `firewall_policy_id` - The ID of the Web Application Firewall Policy which is used for this HTTP Listener.

* `ssl_profile_name` - The name of the associated SSL Profile which is used for this HTTP Listener.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Application Gateway.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Application Gateway.

---

A `private_endpoint_connection` block exports the following:

* `name` - The name of the private endpoint connection.

* `id` - The ID of the private endpoint connection.

---

A `private_link_configuration` block exports the following:

* `id` - The ID of the private link configuration.

* `name` - The name of the private link configuration.

* `ip_configuration` - One or more `ip_configuration` blocks as defined below.

---

An `ip_configuration` block exports the following:

* `name` - The name of the IP configuration.

* `subnet_id` - The ID of the subnet the private link configuration is connected to.

* `private_ip_address_allocation` - The allocation method used for the Private IP Address.

* `primary` - Is this the Primary IP Configuration?

* `private_ip_address` - The Static IP Address which is used.

---

A `match` block exports the following:

* `body` - A snippet from the Response Body which must be present in the Response.

* `status_code` - A list of allowed status codes for this Health Probe.

---

A `path_rule` block exports the following:

* `id` - The ID of the Path Rule.

* `backend_address_pool_id` - The ID of the Backend Address Pool used in this Path Rule.

* `backend_http_settings_id` - The ID of the Backend HTTP Settings Collection used in this Path Rule.

* `redirect_configuration_id` - The ID of the Redirect Configuration used in this Path Rule.

* `rewrite_rule_set_id` - The ID of the Rewrite Rule Set used in this Path Rule.

* `name` - The Name of the Path Rule.

* `paths` - A list of Paths used in this Path Rule.

* `backend_address_pool_name` - The Name of the Backend Address Pool used for this Path Rule.

* `backend_http_settings_name` - The Name of the Backend HTTP Settings Collection used for this Path Rule.

* `redirect_configuration_name` - The Name of a Redirect Configuration used for this Path Rule.

* `rewrite_rule_set_name` - The Name of the Rewrite Rule Set which is used for this Path Rule.

* `firewall_policy_id` - The ID of the Web Application Firewall Policy which is used as an HTTP Listener for this Path Rule.

---

A `probe` block exports the following:

* `id` - The ID of the Probe.

* `host` - The Hostname used for this Probe.

* `interval` - The Interval between two consecutive probes in seconds.

* `name` - The Name of the Probe.

* `protocol` - The Protocol used for this Probe.

* `path` - The Path used for this Probe.

* `timeout` - The Timeout used for this Probe, indicating when a probe becomes unhealthy.

* `unhealthy_threshold` - The Unhealthy Threshold for this Probe, which indicates the amount of retries which will be attempted before a node is deemed unhealthy.

* `port` - Custom port which is used for probing the backend servers.

* `pick_host_name_from_backend_http_settings` - Whether the host header is picked from the backend HTTP settings.

* `match` - A `match` block as defined above.

* `minimum_servers` - The minimum number of servers that are always marked as healthy.

---

A `request_routing_rule` block exports the following:

* `id` - The ID of the Request Routing Rule.

* `http_listener_id` - The ID of the associated HTTP Listener.

* `backend_address_pool_id` - The ID of the associated Backend Address Pool.

* `backend_http_settings_id` - The ID of the associated Backend HTTP Settings Configuration.

* `redirect_configuration_id` - The ID of the associated Redirect Configuration.

* `rewrite_rule_set_id` - The ID of the associated Rewrite Rule Set.

* `url_path_map_id` - The ID of the associated URL Path Map.

* `name` - The Name of this Request Routing Rule.

* `rule_type` - The Type of Routing that is used for this Rule.

* `http_listener_name` - The Name of the HTTP Listener which is used for this Routing Rule.

* `backend_address_pool_name` - The Name of the Backend Address Pool which is used for this Routing Rule.

* `backend_http_settings_name` - The Name of the Backend HTTP Settings Collection which is used for this Routing Rule.

* `redirect_configuration_name` - The Name of the Redirect Configuration which is used for this Routing Rule.

* `rewrite_rule_set_name` - The Name of the Rewrite Rule Set which is used for this Routing Rule.

* `url_path_map_name` - The Name of the URL Path Map which is associated with this Routing Rule.

* `priority` - The Priority of this Routing Rule.

---

A `global` block exports the following:

* `request_buffering_enabled` - Whether Application Gateway's Request buffer is enabled.

* `response_buffering_enabled` - Whether Application Gateway's Response buffer is enabled.

---

A `sku` block exports the following:

* `name` - The Name of the SKU in use for this Application Gateway.

* `tier` - The Tier of the SKU in use for this Application Gateway.

* `capacity` - The Capacity of the SKU in use for this Application Gateway.

---

A `ssl_certificate` block exports the following:

* `id` - The ID of the SSL Certificate.

* `public_cert_data` - The Public Certificate Data associated with the SSL Certificate.

* `name` - The Name of the SSL certificate that is unique within this Application Gateway

* `key_vault_secret_id` - The Secret ID of (base-64 encoded unencrypted pfx) the `Secret` or `Certificate` object stored in Azure KeyVault.

---

A `url_path_map` block exports the following:

* `id` - The ID of the URL Path Map.

* `default_backend_address_pool_id` - The ID of the Default Backend Address Pool.

* `default_backend_http_settings_id` - The ID of the Default Backend HTTP Settings Collection.

* `default_redirect_configuration_id` - The ID of the Default Redirect Configuration.

* `path_rule` - A list of `path_rule` blocks as defined above.

* `name` - The Name of the URL Path Map.

* `default_backend_address_pool_name` - The Name of the Default Backend Address Pool which is used for this URL Path Map.

* `default_backend_http_settings_name` - The Name of the Default Backend HTTP Settings Collection which is used for this URL Path Map.

* `default_redirect_configuration_name` - The Name of the Default Redirect Configuration which is used for this URL Path Map.

* `default_rewrite_rule_set_name` - The Name of the Default Rewrite Rule Set which is used for this URL Path Map.

* `path_rule` - One or more `path_rule` blocks as defined above.

---

A `trusted_client_certificate` block exports the following:

* `id` - The ID of the Trusted Client Certificate in use.

* `name` - The name of the Trusted Client Certificate in use.

* `data` - The content of the Trusted Client Certificate in use.

---

A `ssl_profile` block exports the following:

* `name` - The name of the SSL Profile that is unique within this Application Gateway.

* `trusted_client_certificate_names` - The name of the Trusted Client Certificate that will be used to authenticate requests from clients.

* `verify_client_cert_issuer_dn` - Will the client certificate issuer DN be verified?

* `verify_client_certificate_revocation` - The method used to check client certificate revocation status.

* `ssl_policy` - a `ssl_policy` block as defined below.

---

A `ssl_policy` block exports the following:

* `disabled_protocols` - A list of SSL Protocols which are disabled on this Application Gateway.

* `policy_type` - The Type of the Policy.

* `policy_name` - The Name of the Policy.

* `cipher_suites` - A List of accepted cipher suites.

* `min_protocol_version` - The minimum TLS version.

---

A `waf_configuration` block exports the following:

* `enabled` - Is the Web Application Firewall enabled?

* `firewall_mode` - The Web Application Firewall Mode.

* `rule_set_type` - The Type of the Rule Set used for this Web Application Firewall.

* `rule_set_version` - The Version of the Rule Set used for this Web Application Firewall.

* `disabled_rule_group` - One or more `disabled_rule_group` blocks as defined below.

* `file_upload_limit_mb` - The File Upload Limit in MB.

* `request_body_check` - Is Request Body Inspection enabled?

* `max_request_body_size_kb` - The Maximum Request Body Size in KB.

* `exclusion` - One or more `exclusion` blocks as defined below.

---

A `disabled_rule_group` block exports the following:

* `rule_group_name` - The rule group where specific rules are disabled.

* `rules` - A list of rules which will be disabled in that group.

---

A `exclusion` block exports the following:

* `match_variable` - Match variable of the exclusion rule.

* `selector_match_operator` - Operator which will be used to search in the variable content.

* `selector` - String value which will be used for the filter operation.

---

A `custom_error_configuration` block exports the following:

* `id` - The ID of the Custom Error Configuration.

* `status_code` - Status code of the application gateway custom error.

* `custom_error_page_url` - Error page URL of the application gateway custom error.

---

A `redirect_configuration` block exports the following:

* `id` - The ID of the Redirect Configuration.

* `name` - Unique name of the redirect configuration block

* `redirect_type` - The type of redirect.

* `target_listener_name` - The name of the listener to redirect to.

* `target_url` - The URL to redirect the request to.

* `include_path` - Whether the path is included in the redirected URL.

* `include_query_string` - Whether to include the query string in the redirected URL.

---

An `autoscale_configuration` block exports the following:

* `min_capacity` - Minimum capacity for autoscaling.

* `max_capacity` - Maximum capacity for autoscaling.

---

A `rewrite_rule_set` block exports the following:

* `id` - The ID of the Rewrite Rule Set

* `name` - Unique name of the Rewrite Rule Set

* `rewrite_rule` - One or more `rewrite_rule` blocks as defined below.

---

A `rewrite_rule` block exports the following:

* `name` - Unique name of the Rewrite Rule

* `rule_sequence` - Rule sequence of the Rewrite Rule that determines the order of execution in a set.

* `condition` - One or more `condition` blocks as defined above.

* `request_header_configuration` - One or more `request_header_configuration` blocks as defined above.

* `response_header_configuration` - One or more `response_header_configuration` blocks as defined above.

* `url` - One `url` block as defined below

---

A `condition` block exports the following:

* `variable` - The [variable](https://docs.microsoft.com/azure/application-gateway/rewrite-http-headers#server-variables) of the condition.

* `pattern` - The pattern, either fixed string or regular expression, that evaluates the truthfulness of the condition.

* `ignore_case` - Whether a case insensitive comparison is performed.

* `negate` - Whether the result of the condition evaluation is negated.

---

A `request_header_configuration` block exports the following:

* `header_name` - Header name of the header configuration.

* `header_value` - Header value of the header configuration.

---

A `response_header_configuration` block exports the following:

* `header_name` - Header name of the header configuration.

* `header_value` - Header value of the header configuration.

---

A `url` block exports the following:

* `path` - The URL path to rewrite.

* `query_string` - The query string to rewrite.

* `components` - The components used to rewrite the URL.

* `reroute` - Whether the URL path map is reevaluated after this rewrite has been applied.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Application Gateway.
