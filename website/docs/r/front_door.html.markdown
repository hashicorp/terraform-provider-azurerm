---
subcategory: "Front Door"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor"
sidebar_current: "docs-azurerm-resource-front-door"
description: |-
  Manages an Azure Front Door instance.
---

# azurerm_frontdoor

Manages an Azure Front Door instance.

Azure Front Door Service is Microsoft's highly available and scalable web application acceleration platform and global HTTP(s) load balancer. It provides built-in DDoS protection and application layer security and caching. Front Door enables you to build applications that maximize and automate high-availability and performance for your end-users. Use Front Door with Azure services including Web/Mobile Apps, Cloud Services and Virtual Machines – or combine it with on-premises services for hybrid deployments and smooth cloud migration.

Below are some of the key scenarios that Azure Front Door Service addresses: 
* Use Front Door to improve application scale and availability with instant multi-region failover
* Use Front Door to improve application performance with SSL offload and routing requests to the fastest available application backend.
* Use Front Door for application layer security and DDoS protection for your application.

## Example Usage

```hcl
resource "azurerm_frontdoor" "example" {
  name                                         = "example-FrontDoor"
  location                                     = "${azurerm_resource_group.example.location}"
  resource_group_name                          = "${azurerm_resource_group.example.name}"
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "exampleRoutingRule1"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = ["exampleFrontendEndpoint1"]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "exampleBackendBing"
    }
  }

  backend_pool_load_balancing {
    name = "exampleLoadBalancingSettings1"
  }

  backend_pool_health_probe {
    name = "exampleHealthProbeSetting1"
  }

  backend_pool {
    name = "exampleBackendBing"
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = "exampleLoadBalancingSettings1"
    health_probe_name   = "exampleHealthProbeSetting1"
  }

  frontend_endpoint {
    name                              = "exampleFrontendEndpoint1"
    host_name                         = "example-FrontDoor.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Front Door which is globally unique. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Name of the Resource group within the Azure subscription. Changing this forces a new resource to be created.

* `location` - (Required) Resource location. Changing this forces a new resource to be created.

* `backend_pool` - (Required) A `backend_pool` block as defined below.

* `backend_pool_health_probe` - (Required) A `backend_pool_health_probe` block as defined below.

* `backend_pool_load_balancing` - (Required) A `backend_pool_load_balancing` block as defined below.

* `enforce_backend_pools_certificate_name_check` - (Required) Whether to enforce certificate name check on HTTPS requests to all backend pools. No effect on non-HTTPS requests. Permitted values are `true` or `false`.

* `load_balancer_enabled` - (Optional) Operational status of the Front Door load balancer. Permitted values are `true` or `false` Defaults to `true`.

* `friendly_name` - (Optional) A friendly name for the Front Door service.

* `frontend_endpoint` - (Required) A `frontend_endpoint` block as defined below.

* `routing_rule` - (Required) A `routing_rule` block as defined below.

* `tags` - (Optional) Resource tags.

---

The `backend_pool` block supports the following:

* `name` - (Required) The name of the `Backend Pool`.

* `backend` - (Required) A `backend` block as defined below.

* `load_balancing_name` - (Required) The name property of the `backend_pool_load_balancing` block whithin this resource to use for the `Backend Pool`.

* `health_probe_name` - (Required) The name property of a `backend_pool_health_probe` block whithin this resource to use for the `Backend Pool`.

---

The `backend` block supports the following:

* `address` - (Required) Location of the backend (IP address or FQDN)

* `host_header` - (Required) The value to use as the host header sent to the backend.

* `http_port` - (Required) The HTTP TCP port number. Possible values are between `1` - `65535`.

* `https_port` - (Required) The HTTPS TCP port number. Possible values are between `1` - `65535`.

* `priority` - (Optional) Priority to use for load balancing. Higher priorities will not be used for load balancing if any lower priority backend is healthy. Defaults to `1`.

* `weight` - (Optional) Weight of this endpoint for load balancing purposes. Defaults to `50`.

---

The `frontend_endpoint` block supports the following:

* `name` - (Required) The name of the Frontend Endpoint.

* `host_name` - (Required) The host name of the Frontend Endpoint. Must be a domain name.

* `custom_https_provisioning_enabled` - (Required) Whether to allow HTTPS protocol for a custom domain that's associated with Front Door to ensure sensitive data is delivered securely via TLS/SSL encryption when sent across the internet. Valid options are `true` or `false`.

* `session_affinity_enabled` - (Optional) Whether to allow session affinity on this host. Valid options are `true` or `false` Defaults to `false`.

* `session_affinity_ttl_seconds` - (Optional) The TTL to use in seconds for session affinity, if applicable. Defaults to `0`.

* `web_application_firewall_policy_link_id` - (Optional) Defines the Web Application Firewall policy `ID` for each host.

---

The `backend_pool_health_probe` block supports the following:

* `name` - (Required) The name of the Azure Front Door Backend Health Probe.

* `path` - (Optional) The path to use for the Backend Health Probe. Default is `/`.

* `protocol` - (Optional) Protocol scheme to use for the Backend Health Probe. Defaults to `Http`.

* `interval_in_seconds` - (Optional) The number of seconds between health probes. Defaults to `120`.

---

The `backend_pool_load_balancing` block supports the following:

* `name` - (Required) The name of the Azure Front Door Backend Load Balancer.

* `sample_size` - (Optional) The number of samples to consider for load balancing decisions. Defaults to `4`.

* `successful_samples_required` - (Optional) The number of samples within the sample period that must succeed. Defaults to `2`.

* `additional_latency_milliseconds` - (Optional) The additional latency in milliseconds for probes to fall into the lowest latency bucket. Defaults to `0`.

---

The `routing_rule` block supports the following:

* `name` - (Required) The name of the Front Door Backend Routing Rule.

* `frontend_endpoints` - (Required) The names of the `frontend_endpoint` blocks whithin this resource to associate with this `routing_rule`.

* `accepted_protocols` - (Optional) Protocol schemes to match for the Backend Routing Rule. Defaults to `Http`.

* `patterns_to_match` - (Optional) The route patterns for the Backend Routing Rule. Defaults to `/*`.

* `enabled` - (Optional) `Enable` or `Disable` use of this Backend Routing Rule. Permitted values are `true` or `false`. Defaults to `true`.

* `forwarding_configuration` - (Optional) A `forwarding_configuration` block as defined below.

* `redirect_configuration`   - (Optional) A `redirect_configuration` block as defined below.

---

The `forwarding_configuration` block supports the following:

* `backend_pool_name` - (Required) The name of the Front Door Backend Pool. 

* `cache_use_dynamic_compression` - (Optional) Whether to use dynamic compression when caching. Valid options are `true` or `false`. Defaults to `true`.

* `cache_query_parameter_strip_directive` - (Optional) Defines cache behavior in releation to query string parameters. Valid options are `StripAll` or `StripNone`. Defaults to `StripNone`

* `custom_forwarding_path` - (Optional) Path to use when constructing the request to forward to the backend. This functions as a URL Rewrite. Default behavior preserves the URL path.

* `forwarding_protocol` - (Optional) Protocol to use when redirecting. Valid options are `HttpOnly`, `HttpsOnly`, or `MatchRequest`. Defaults to `MatchRequest`.

---

The `redirect_configuration` block supports the following:

* `custom_host` - (Optional)  Set this to change the URL for the redirection. 

* `redirect_protocol` - (Optional) Protocol to use when redirecting. Valid options are `HttpOnly`, `HttpsOnly`, `MatchRequest`. Defaults to `MatchRequest`

* `redirect_type` - (Optional) Status code for the redirect. Valida options are `Moved`, `Found`, `TemporaryRedirect`, `PermanentRedirect`. Defaults to `Found`

* `custom_fragment` - (Optional) The destination fragment in the portion of URL after '#'. Set this to add a fragment to the redirect URL.

* `custom_path` - (Optional) The path to retain as per the incoming request, or update in the URL for the redirection.

* `custom_query_string` - (Optional) Replace any existing query string from the incoming request URL.

---

The `custom_https_configuration` block supports the following:

* `certificate_source` - (Optional) Certificate source to encrypted `HTTPS` traffic with. Allowed values are `FrontDoor` or `AzureKeyVault`. Defaults to `FrontDoor`.

The following attributes are only valid if `certificate_source` is set to `AzureKeyVault`:

* `azure_key_vault_certificate_vault_id` - (Required) The `id` of the Key Vault containing the SSL certificate.

* `azure_key_vault_certificate_secret_name` - (Required) The name of the Key Vault secret representing the full certificate PFX.

* `azure_key_vault_certificate_secret_version` - (Required) The version of the Key Vault secret representing the full certificate PFX.

~> **Note:** In order to enable the use of your own custom `HTTPS certificate` you must grant `Azure Front Door Service` access to your key vault. For instuctions on how to configure your `Key Vault` correctly please refer to the [product documentation](https://docs.microsoft.com/en-us/azure/frontdoor/front-door-custom-domain-https#option-2-use-your-own-certificate). 

---

## Attributes Reference

`backend_pool` exports the following:

* `id` - The Resource ID of the Azure Front Door Backend Pool.


`backend` exports the following:

* `id` - The Resource ID of the Azure Front Door Backend.


`frontend_endpoint` exports the following:

* `id` - The Resource ID of the Azure Front Door Frontend Endpoint.

* `provisioning_state` - Provisioning state of the Front Door.

* `provisioning_substate` - Provisioning substate of the Front Door

[//]: * "* `web_application_firewall_policy_link_id` - (Optional) The `id` of the `web_application_firewall_policy_link` to use for this Frontend Endpoint."


`backend_pool_health_probe` exports the following:

* `id` - The Resource ID of the Azure Front Door Backend Health Probe.


`backend_pool_load_balancing` exports the following:

* `id` - The Resource ID of the Azure Front Door Backend Load Balancer.


`routing_rule`  exports the following:

* `id` - The Resource ID of the Azure Front Door Backend Routing Rule.


The following attributes are exported:

* `cname` - The host that each frontendEndpoint must CNAME to.

* `id` - Resource ID.

## Import

Front Doors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.Network/frontdoors/frontdoor1
```
