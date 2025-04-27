---
subcategory: "NGINX"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_nginx_deployment"
description: |-
  Gets information about an existing NGINX Deployment.
---

# Data Source: azurerm_nginx_deployment

Use this data source to access information about an existing NGINX Deployment.

## Example Usage

```hcl
data "azurerm_nginx_deployment" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_nginx_deployment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this NGINX Deployment.

* `resource_group_name` - (Required) The name of the Resource Group where the NGINX Deployment exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NGINX Deployment.

* `capacity` - The number of NGINX capacity units for this NGINX Deployment.

* `auto_scale_profile` - An `auto_scale_profile` block as defined below.

* `diagnose_support_enabled` - Whether metrics are exported to Azure Monitor.

* `email` - Preferred email associated with the NGINX Deployment.

* `frontend_private` - A `frontend_private` block as defined below.

* `frontend_public` - A `frontend_public` block as defined below.

* `identity` - A `identity` block as defined below.

* `ip_address` - The IP address of the NGINX Deployment.

* `location` - The Azure Region where the NGINX Deployment exists.

* `network_interface` - A `network_interface` block as defined below.

* `nginx_version` - NGINX version of the Deployment.

* `sku` - The NGINX Deployment SKU.

* `automatic_upgrade_channel` - The automatic upgrade channel for this NGINX deployment.

* `web_application_firewall` - A `web_application_firewall` block as defined below.

* `dataplane_api_endpoint` - The dataplane API endpoint of the NGINX Deployment.

* `tags` - A mapping of tags assigned to the NGINX Deployment.

---

A `frontend_private` block exports the following:

* `allocation_method` - The method of allocating the private IP to the NGINX Deployment.

* `ip_address` - Private IP address of the NGINX Deployment.

* `subnet_id` - The subnet ID of the NGINX Deployment.

---

A `frontend_public` block exports the following:

* `ip_address` - The list of Public IP Resource IDs for this NGINX Deployment.

---

A `identity` block exports the following:

* `identity_ids` - List of identities attached to the NGINX Deployment.

* `type` - Type of identity attached to the NGINX Deployment.

---

A `network_interface` block exports the following:

* `subnet_id` - The subnet resource ID of the NGINX Deployment.

---

An `auto_scale_profile` block exports the following:

* `name` - Name of the autoscaling profile.

* `min_capacity` - The minimum number of NGINX capacity units for this NGINX Deployment.

* `max_capacity` - The maximum number of NGINX capacity units for this NGINX Deployment.

---

A `web_application_firewall` block exports the following:

* `activation_state_enabled` - Whether WAF is enabled/disabled for this NGINX Deployment.
* `status` - A `status` block as defined below.

---

A `web_application_firewall.status` block exports the following:

* `attack_signatures_package` - Indicates the version of the attack signatures package used by NGINX App Protect.

* `bot_signatures_package` - Indicates the version of the bot signatures package used by NGINX App Protect.

* `threat_campaigns_package` - Indicates the version of the threat campaigns package used by NGINX App Protect.

* `component_versions` - Indicates the version of the WAF Engine and Nginx WAF Module used by NGINX App Protect.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NGINX Deployment.
