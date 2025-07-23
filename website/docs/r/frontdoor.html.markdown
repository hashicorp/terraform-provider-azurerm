---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor"
description: |-
  Manages an Azure Front Door (classic) instance.
---

# azurerm_frontdoor

!> **Note:** This deploys an Azure Front Door (classic) resource which has been deprecated and will receive security updates only. Please migrate your existing Azure Front Door (classic) deployments to the new [Azure Front Door (standard/premium) resources](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/cdn_frontdoor_custom_domain). For your convenience, the service team has exposed a `Front Door Classic` to `Front Door Standard/Premium` [migration tool](https://learn.microsoft.com/azure/frontdoor/tier-migration) to allow you to migrate your existing `Front Door Classic` instances to the new `Front Door Standard/Premium` product tiers.

Manages an Azure Front Door (classic) instance.

Azure Front Door Service is Microsoft's highly available and scalable web application acceleration platform and global HTTP(S) load balancer. It provides built-in DDoS protection and application layer security and caching. Front Door enables you to build applications that maximize and automate high-availability and performance for your end-users. Use Front Door with Azure services including Web/Mobile Apps, Cloud Services and Virtual Machines – or combine it with on-premises services for hybrid deployments and smooth cloud migration.

Below are some of the key scenarios that Azure Front Door Service addresses:

* Use Front Door to improve application scale and availability with instant multi-region failover
* Use Front Door to improve application performance with SSL offload and routing requests to the fastest available application backend.
* Use Front Door for application layer security and DDoS protection for your application.

!> **Note:** The `custom_https_provisioning_enabled` field and the `custom_https_configuration` block have been removed from the `azurerm_frontdoor` resource in the `v2.58.0` provider due to changes made by the service team. If you wish to enable the custom HTTPS configuration functionality within your `azurerm_frontdoor` resource moving forward you will need to define a separate `azurerm_frontdoor_custom_https_configuration` block in your configuration file.

!> **Note:** With the release of the `v2.58.0` provider, if you run the `apply` command against an existing Front Door resource it **will not** apply the detected changes. Instead it will persist the `explicit_resource_order` mapping structure to the state file. Once this operation has completed the resource will resume functioning normally.This change in behavior in Terraform is due to an issue where the underlying service teams API is now returning the response JSON out of order from the way it was sent to the resource via Terraform causing unexpected discrepancies in the `plan` after the resource has been provisioned. If your pre-existing Front Door instance contains `custom_https_configuration` blocks there are additional steps that will need to be completed to successfully migrate your Front Door onto the `v2.58.0` provider which [can be found in this guide](https://registry.terraform.io/providers/hashicorp/azurerm/2.59.0/docs/guides/2.58.0-frontdoor-upgrade-guide).

!> **Note:** On `1 April 2025`, Azure Front Door (classic) will be retired for the public cloud, existing Azure Front Door (classic) resources must be migrated out of Azure Front Door (classic) to Azure Front Door Standard/Premium before `1 October 2025` to avoid potential disruptions in service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "FrontDoorExampleResourceGroup"
  location = "West Europe"
}

resource "azurerm_frontdoor" "example" {
  name                = "example-FrontDoor"
  resource_group_name = azurerm_resource_group.example.name

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
    name      = "exampleFrontendEndpoint1"
    host_name = "example-FrontDoor.azurefd.net"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door service. Must be globally unique. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) Specifies the name of the Resource Group in which the Front Door service should exist. Changing this forces a new resource to be created.

* `backend_pool` - (Required) A `backend_pool` block as defined below.

-> **Note:** Azure by default allows specifying up to 50 Backend Pools - but this quota can be increased via Microsoft Support.

* `backend_pool_health_probe` - (Required) A `backend_pool_health_probe` block as defined below.

* `backend_pool_load_balancing` - (Required) A `backend_pool_load_balancing` block as defined below.

* `load_balancer_enabled` - (Optional) Should the Front Door Load Balancer be Enabled? Defaults to `true`.

* `friendly_name` - (Optional) A friendly name for the Front Door service.

* `backend_pool_settings` - (Optional) A `backend_pool_settings` block as defined below.

* `frontend_endpoint` - (Required) A `frontend_endpoint` block as defined below.

* `routing_rule` - (Required) A `routing_rule` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `backend` block supports the following:

* `enabled` - (Optional) Specifies if the backend is enabled or not. Valid options are `true` or `false`. Defaults to `true`.

* `address` - (Required) Location of the backend (IP address or FQDN)

* `host_header` - (Required) The value to use as the host header sent to the backend.

* `http_port` - (Required) The HTTP TCP port number. Possible values are between `1` - `65535`.

* `https_port` - (Required) The HTTPS TCP port number. Possible values are between `1` - `65535`.

* `priority` - (Optional) Priority to use for load balancing. Higher priorities will not be used for load balancing if any lower priority backend is healthy. Defaults to `1`.

* `weight` - (Optional) Weight of this endpoint for load balancing purposes. Defaults to `50`.

---

The `backend_pool` block supports the following:

* `name` - (Required) Specifies the name of the Backend Pool.

* `backend` - (Required) A `backend` block as defined below.

* `load_balancing_name` - (Required) Specifies the name of the `backend_pool_load_balancing` block within this resource to use for this `Backend Pool`.

* `health_probe_name` - (Required) Specifies the name of the `backend_pool_health_probe` block within this resource to use for this `Backend Pool`.

---

The `backend_pool_settings` block supports the following:

* `backend_pools_send_receive_timeout_seconds` - (Optional) Specifies the send and receive timeout on forwarding request to the backend. When the timeout is reached, the request fails and returns. Possible values are between `0` - `240`. Defaults to `60`.

* `enforce_backend_pools_certificate_name_check` - (Required) Enforce certificate name check on `HTTPS` requests to all backend pools, this setting will have no effect on `HTTP` requests. Permitted values are `true` or `false`.

-> **Note:** `backend_pools_send_receive_timeout_seconds` and `enforce_backend_pools_certificate_name_check` apply to all backend pools.

---

The `frontend_endpoint` block supports the following:

* `name` - (Required) Specifies the name of the `frontend_endpoint`.

* `host_name` - (Required) Specifies the host name of the `frontend_endpoint`. Must be a domain name. In order to use a name.azurefd.net domain, the name value must match the Front Door name.

* `session_affinity_enabled` - (Optional) Whether to allow session affinity on this host. Valid options are `true` or `false` Defaults to `false`.

* `session_affinity_ttl_seconds` - (Optional) The TTL to use in seconds for session affinity, if applicable. Defaults to `0`.

* `web_application_firewall_policy_link_id` - (Optional) Defines the Web Application Firewall policy `ID` for each host.

---

The `backend_pool_health_probe` block supports the following:

* `name` - (Required) Specifies the name of the Health Probe.

* `enabled` - (Optional) Is this health probe enabled? Defaults to `true`.

* `path` - (Optional) The path to use for the Health Probe. Default is `/`.

* `protocol` - (Optional) Protocol scheme to use for the Health Probe. Possible values are `Http` and `Https`. Defaults to `Http`.

* `probe_method` - (Optional) Specifies HTTP method the health probe uses when querying the backend pool instances. Possible values include: `GET` and `HEAD`. Defaults to `GET`.

-> **Note:** Use the `HEAD` method if you do not need to check the response body of your health probe.

* `interval_in_seconds` - (Optional) The number of seconds between each Health Probe. Defaults to `120`.

---

The `backend_pool_load_balancing` block supports the following:

* `name` - (Required) Specifies the name of the Load Balancer.

* `sample_size` - (Optional) The number of samples to consider for load balancing decisions. Defaults to `4`.

* `successful_samples_required` - (Optional) The number of samples within the sample period that must succeed. Defaults to `2`.

* `additional_latency_milliseconds` - (Optional) The additional latency in milliseconds for probes to fall into the lowest latency bucket. Defaults to `0`.

---

The `routing_rule` block supports the following:

* `name` - (Required) Specifies the name of the Routing Rule.

* `frontend_endpoints` - (Required) The names of the `frontend_endpoint` blocks within this resource to associate with this `routing_rule`.

* `accepted_protocols` - (Required) Protocol schemes to match for the Backend Routing Rule. Possible values are `Http` and `Https`.

* `patterns_to_match` - (Required) The route patterns for the Backend Routing Rule.

* `enabled` - (Optional) `Enable` or `Disable` use of this Backend Routing Rule. Permitted values are `true` or `false`. Defaults to `true`.

* `forwarding_configuration` - (Optional) A `forwarding_configuration` block as defined below.

* `redirect_configuration` - (Optional) A `redirect_configuration` block as defined below.

---

The `forwarding_configuration` block supports the following:

* `backend_pool_name` - (Required) Specifies the name of the Backend Pool to forward the incoming traffic to.

* `cache_enabled` - (Optional) Specifies whether to Enable caching or not. Valid options are `true` or `false`. Defaults to `false`.

* `cache_use_dynamic_compression` - (Optional) Whether to use dynamic compression when caching. Valid options are `true` or `false`. Defaults to `false`.

* `cache_query_parameter_strip_directive` - (Optional) Defines cache behaviour in relation to query string parameters. Valid options are `StripAll`, `StripAllExcept`, `StripOnly` or `StripNone`. Defaults to `StripAll`.

* `cache_query_parameters` - (Optional) Specify query parameters (array). Works only in combination with `cache_query_parameter_strip_directive` set to `StripAllExcept` or `StripOnly`.

* `cache_duration` - (Optional) Specify the minimum caching duration (in ISO8601 notation e.g. `P1DT2H` for 1 day and 2 hours). Needs to be greater than 0 and smaller than 365 days. `cache_duration` works only in combination with `cache_enabled` set to `true`.

* `custom_forwarding_path` - (Optional) Path to use when constructing the request to forward to the backend. This functions as a URL Rewrite. Default behaviour preserves the URL path.

* `forwarding_protocol` - (Optional) Protocol to use when redirecting. Valid options are `HttpOnly`, `HttpsOnly`, or `MatchRequest`. Defaults to `HttpsOnly`.

---

The `redirect_configuration` block supports the following:

* `custom_host` - (Optional) Set this to change the URL for the redirection.

* `redirect_protocol` - (Required) Protocol to use when redirecting. Valid options are `HttpOnly`, `HttpsOnly`, or `MatchRequest`.

* `redirect_type` - (Required) Status code for the redirect. Valida options are `Moved`, `Found`, `TemporaryRedirect`, `PermanentRedirect`.

* `custom_fragment` - (Optional) The destination fragment in the portion of URL after '#'. Set this to add a fragment to the redirect URL.

* `custom_path` - (Optional) The path to retain as per the incoming request, or update in the URL for the redirection.

* `custom_query_string` - (Optional) Replace any existing query string from the incoming request URL.

---

## Attributes Reference

-> **Note:** UPCOMING BREAKING CHANGE: In order to address the ordering issue we have changed the design on how to retrieve existing sub resources such as backend pool health probes, backend pool loadbalancer settings, backend pools, frontend endpoints and routing rules. Existing design will be deprecated and will result in an incorrect configuration. Please refer to the updated documentation below for more information.

* `backend_pool_health_probes` - A map/dictionary of Backend Pool Health Probe Names (key) to the Backend Pool Health Probe ID (value)
* `backend_pool_load_balancing_settings` - A map/dictionary of Backend Pool Load Balancing Setting Names (key) to the Backend Pool Load Balancing Setting ID (value)
* `backend_pools` - A map/dictionary of Backend Pool Names (key) to the Backend Pool ID (value)
* `frontend_endpoints` - A map/dictionary of Frontend Endpoint Names (key) to the Frontend Endpoint ID (value)
* `routing_rules` - A map/dictionary of Routing Rule Names (key) to the Routing Rule ID (value)

---

`backend` exports the following:

* `id` - The ID of the Azure Front Door Backend.

---

`backend_pool` exports the following:

* `id` - The ID of the Azure Front Door Backend Pool.

---

`backend_pool_health_probe` exports the following:

* `id` - The ID of the Azure Front Door Backend Health Probe.

---

`backend_pool_load_balancing` exports the following:

* `id` - The ID of the Azure Front Door Backend Load Balancer.

---

`frontend_endpoint` exports the following:

* `id` - The ID of the Azure Front Door Frontend Endpoint.

---

`routing_rule` exports the following:

* `id` - The ID of the Azure Front Door Backend Routing Rule.

---

In addition to the Arguments listed above - the following Attributes are exported:

* `cname` - The host that each frontendEndpoint must CNAME to.

* `header_frontdoor_id` - The unique ID of the Front Door which is embedded into the incoming headers `X-Azure-FDID` attribute and maybe used to filter traffic sent by the Front Door to your backend.

* `id` - The ID of the FrontDoor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the FrontDoor.
* `read` - (Defaults to 5 minutes) Used when retrieving the FrontDoor.
* `update` - (Defaults to 6 hours) Used when updating the FrontDoor.
* `delete` - (Defaults to 6 hours) Used when deleting the FrontDoor.

## Import

Front Doors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/frontDoors/frontdoor1
```
