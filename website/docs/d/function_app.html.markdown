---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app"
description: |-
  Gets information about an existing Function App.

---

# Data Source: azurerm_function_app

Use this data source to access information about a Function App.

## Example Usage

```hcl
data "azurerm_function_app" "example" {
  name                = "test-azure-functions"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Function App resource.

* `resource_group_name` - The name of the Resource Group where the Function App exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Function App

* `app_service_plan_id` - The ID of the App Service Plan within which to create this Function App.

* `app_settings` - A key-value pair of App Settings.

* `connection_string` - An `connection_string` block as defined below.

* `custom_domain_verification_id` - An identifier used by App Service to perform domain ownership verification via DNS TXT record.

* `default_hostname` - The default hostname associated with the Function App.

* `enabled` - Is the Function App enabled?

* `identity` - A `identity` block as defined below.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

* `os_type` - A string indicating the Operating System type for this function app.

~> **NOTE:** This value will be `linux` for Linux Derivatives, or an empty string for Windows. 

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses.

* `possible_outbound_ip_addresses` - A comma separated list of outbound IP addresses, not all of which are necessarily in use. Superset of `outbound_ip_addresses`.

* `source_control` - A `source_control` block as defined below.

---

The `connection_string` supports the following:

* `name` - The name of the Connection String.
* `type` - The type of the Connection String. 
* `value` - The value for the Connection String.

---

The `site_credential` block exports the following:

* `username` - The username which can be used to publish to this App Service
* `password` - The password associated with the username, which can be used to publish to this App Service.

---

An `ip_restriction` block exports the following:

* `ip_address` - The IP Address used for this IP Restriction.

* `service_tag` - The Service Tag used for this IP Restriction.

* `subnet_mask` - The Subnet mask used for this IP Restriction.

* `name` - The name for this IP Restriction.

* `priority` - The priority for this IP Restriction.

* `action` - Does this restriction `Allow` or `Deny` access for this IP range?

---
An `scm_ip_restriction` block exports the following:  

* `ip_address` - The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - The Virtual Network Subnet ID used for this IP Restriction.

* `name` - The name for this IP Restriction.

* `priority` - The priority for this IP Restriction.

* `action` - Allow or Deny access for this IP range. Defaults to Allow.  

---

A `site_config` block exports the following:

* `always_on` - Is the app loaded at all times?

* `cors` - A `cors` block as defined above.

* `http2_enabled` - Is HTTP2 Enabled on this App Service?

* `ftps_state` - State of FTP / FTPS service for this AppService.

* `ip_restriction` - One or more `ip_restriction` blocks as defined above.

* `java_version` - Java version hosted by the function app in Azure.

* `pre_warmed_instance_count` - The number of pre-warmed instances for this function app. Only applicable to apps on the Premium plan.

* `scm_use_main_ip_restriction` - IP security restrictions for scm to use main.  

* `scm_ip_restriction` - One or more `scm_ip_restriction` blocks as defined above.

* `linux_fx_version` - Linux App Framework and version for the AppService.

* `min_tls_version` - The minimum supported TLS version for this App Service.

* `scm_type` - The type of Source Control enabled for this App Service.

* `use_32_bit_worker_process` - Does the App Service run in 32 bit mode, rather than 64 bit mode?

* `websockets_enabled` - Are WebSockets enabled for this App Service?

---

A `source_control` block exports the following:

* `repo_url` -  The URL of the source code repository.

* `branch` - The branch of the remote repository in use. 

* `manual_integration` - Limits to manual integration.  

* `rollback_enabled` - Is roll-back enabled for the repository.

* `use_mercurial` - Uses Mercurial if `true`, otherwise uses Git. 

---

An `identity` block exports the following:

* `principal_id` - The ID of the System Managed Service Principal assigned to the function app.

* `tenant_id` - The ID of the Tenant of the System Managed Service Principal assigned to the function app.

* `type` - The identity type of the Managed Identity assigned to the function app.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Function App.
