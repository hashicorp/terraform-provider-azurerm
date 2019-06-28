---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_front_door"
sidebar_current: "docs-azurerm-resource-front-door"
description: |-
  Manage Azure FrontDoor instance.
---

# azurerm_front_door

Manage Azure FrontDoor instance.


## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Front Door which is globally unique. Changing this forces a new resource to be created.

* `resource_group` - (Required) Name of the Resource group within the Azure subscription. Changing this forces a new resource to be created.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `backend_pools` - (Optional) One `backend_pool` block defined below.

* `backend_pools_settings` - (Optional) One `backend_pools_setting` block defined below.

* `enabled_state` - (Optional) Operational status of the Front Door load balancer. Permitted values are 'Enabled' or 'Disabled' Defaults to `Enabled`.

* `friendly_name` - (Optional) A friendly name for the frontDoor

* `frontend_endpoints` - (Optional) One `frontend_endpoint` block defined below.

* `health_probe_settings` - (Optional) One `health_probe_setting` block defined below.

* `load_balancing_settings` - (Optional) One `load_balancing_setting` block defined below.

* `resource_state` - (Optional) Resource status of the Front Door. Defaults to `Creating`.

* `routing_rules` - (Optional) One `routing_rule` block defined below.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

The `backend_pool` block supports the following:

* `id` - (Optional) Resource ID.

* `backends` - (Optional) One `backend` block defined below.

* `load_balancing_settings` - (Optional) One `load_balancing_setting` block defined below.

* `health_probe_settings` - (Optional) One `health_probe_setting` block defined below.

* `resource_state` - (Optional) Resource status. Defaults to `Creating`.

* `name` - (Optional) Resource name.


---

The `backend` block supports the following:

* `address` - (Optional) Location of the backend (IP address or FQDN)

* `http_port` - (Optional) The HTTP TCP port number. Must be between 1 and 65535.

* `https_port` - (Optional) The HTTPS TCP port number. Must be between 1 and 65535.

* `enabled_state` - (Optional) Whether to enable use of this backend. Permitted values are 'Enabled' or 'Disabled' Defaults to `Enabled`.

* `priority` - (Optional) Priority to use for load balancing. Higher priorities will not be used for load balancing if any lower priority backend is healthy.

* `weight` - (Optional) Weight of this endpoint for load balancing purposes.

* `backend_host_header` - (Optional) The value to use as the host header sent to the backend. If blank or unspecified, this defaults to the incoming host.

---

The `load_balancing_setting` block supports the following:

* `id` - (Optional) Resource ID.

---

The `health_probe_setting` block supports the following:

* `id` - (Optional) Resource ID.

---

The `backend_pools_setting` block supports the following:

* `enforce_certificate_name_check` - (Optional) Whether to enforce certificate name check on HTTPS requests to all backend pools. No effect on non-HTTPS requests. Defaults to `Enabled`.

---

The `frontend_endpoint` block supports the following:

* `id` - (Optional) Resource ID.

* `host_name` - (Optional) The host name of the frontendEndpoint. Must be a domain name.

* `session_affinity_enabled_state` - (Optional) Whether to allow session affinity on this host. Valid options are 'Enabled' or 'Disabled' Defaults to `Enabled`.

* `session_affinity_ttl_seconds` - (Optional) UNUSED. This field will be ignored. The TTL to use in seconds for session affinity, if applicable.

* `web_application_firewall_policy_link` - (Optional) One `web_application_firewall_policy_link` block defined below.

* `resource_state` - (Optional) Resource status. Defaults to `Creating`.

* `name` - (Optional) Resource name.


---

The `web_application_firewall_policy_link` block supports the following:

* `id` - (Optional) Resource ID.

---

The `health_probe_setting` block supports the following:

* `id` - (Optional) Resource ID.

* `path` - (Optional) The path to use for the health probe. Default is /

* `protocol` - (Optional) Protocol scheme to use for this probe Defaults to `Http`.

* `interval_in_seconds` - (Optional) The number of seconds between health probes.

* `resource_state` - (Optional) Resource status. Defaults to `Creating`.

* `name` - (Optional) Resource name.

---

The `load_balancing_setting` block supports the following:

* `id` - (Optional) Resource ID.

* `sample_size` - (Optional) The number of samples to consider for load balancing decisions

* `successful_samples_required` - (Optional) The number of samples within the sample period that must succeed

* `additional_latency_milliseconds` - (Optional) The additional latency in milliseconds for probes to fall into the lowest latency bucket

* `resource_state` - (Optional) Resource status. Defaults to `Creating`.

* `name` - (Optional) Resource name.

---

The `routing_rule` block supports the following:

* `id` - (Optional) Resource ID.

* `frontend_endpoints` - (Optional) One `frontend_endpoint` block defined below.

* `accepted_protocols` - (Optional) Protocol schemes to match for this rule Defaults to `Http`.

* `patterns_to_match` - (Optional) The route patterns of the rule.

* `enabled_state` - (Optional) Whether to enable use of this rule. Permitted values are 'Enabled' or 'Disabled' Defaults to `Enabled`.

* `resource_state` - (Optional) Resource status. Defaults to `Creating`.

* `name` - (Optional) Resource name.


---

The `frontend_endpoint` block supports the following:

* `id` - (Optional) Resource ID.

## Attributes Reference

The following attributes are exported:

* `provisioning_state` - Provisioning state of the Front Door.

* `cname` - The host that each frontendEndpoint must CNAME to.

* `id` - Resource ID.

* `type` - Resource type.
