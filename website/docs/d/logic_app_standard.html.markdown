---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_logic_app_standard"
description: |-
  Gets information about an existing Logic App Standard instance.
---

# Data Source: azurerm_logic_app_standard

Use this data source to access information about an existing Logic App Standard instance.

## Example Usage

```hcl
data "azurerm_logic_app_standard" "example" {
  name                = "example-logic-app"
  resource_group_name = "example-rg"
}

output "id" {
  value = data.azurerm_logic_app_standard.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Logic App.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Logic App Standard ID.

* `location` - The Azure location where the Logic App Standard exists.

* `identity` - An `identity` block as defined below.

* `app_service_plan_id` - The ID of the App Service Plan.

* `app_settings` - A map of key-value pairs for [App Settings](https://docs.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.

* `use_extension_bundle` - Whether the logic app should use the bundled extension package.

* `bundle_version` - Controls the allowed range for bundle versions.

* `client_affinity_enabled` - Should the Logic App send session affinity cookies, which route client requests in the same session to the same instance.

* `client_certificate_mode` - The mode of the Logic App's client certificates requirement for incoming requests.

* `connection_string` - A `connection_string` block as defined below.

* `custom_domain_verification_id` - The custom domain verification of the Logic App.

* `default_hostname` - The default hostname of the Logic App.

* `enabled` - Whether the Logic App is enabled.

* `ftp_publish_basic_authentication_enabled` - Whether the default FTP basic authentication publishing profile is enabled.

* `https_only` - Whether the Logic App can only be accessed via HTTPS.

* `kind` - The kind of the Logic App.

* `outbound_ip_addresses` - The outbound IP addresses of the Logic App.

* `possible_outbound_ip_addresses` - The possible outbound IP addresses of the Logic App.

* `public_network_access` - Whether Public Network Access should be enabled or not.

* `scm_publish_basic_authentication_enabled` - Whether the default SCM basic authentication publishing profile is enabled.

* `site_config` - A `site_config` object as defined below.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this Logic App.

* `storage_account_name` - The backend storage account name which will be used by this Logic App (e.g. for Stateful workflows data).

* `storage_account_access_key` - The access key which will be used to access the backend storage account for the Logic App.

* `storage_account_share_name` - The name of the share used by the logic app.

* `tags` - A mapping of tags assigned to the resource.

* `version` - The runtime version associated with the Logic App.

* `virtual_network_subnet_id` - The subnet ID for the Logic App.

---

The `identity` block exports the following:

* `type` - The Type of Managed Identity assigned to this Logic App Workflow.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Logic App Workflow.

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Logic App Workflow.

---

he `site_credential` block exports the following:

* `username` - The username which can be used to publish to this Logic App.

* `password` - The password associated with the username, which can be used to publish to this Logic App.

---

The `site_config` block exports the following:

* `always_on` - Should the Logic App be loaded at all times?

* `app_scale_limit` - The number of workers this Logic App can scale out to. Only applicable to apps on the Consumption and Premium plan.

* `auto_swap_slot_name` - The Auto-swap slot name.

* `cors` - A `cors` block as defined below.

* `dotnet_framework_version` - The version of the .NET framework's CLR used in this Logic App.

* `elastic_instance_minimum` - The number of minimum instances for this Logic App Only affects apps on the Premium plan.

* `ftps_state` - The state of FTP / FTPS service for this Logic App.

* `health_check_path` - Path which will be checked for this Logic App health.

* `http2_enabled` - Specifies whether the HTTP2 protocol should be enabled.

* `ip_restriction` - A list of `ip_restriction` objects representing IP restrictions as defined below.

* `scm_ip_restriction` - A list of `scm_ip_restriction` objects representing SCM IP restrictions as defined below.

* `scm_use_main_ip_restriction` - Should the Logic App `ip_restriction` configuration be used for the SCM too.

* `scm_min_tls_version` - The minimum version of TLS required for SSL requests to the SCM site.

* `scm_type` - The type of Source Control used by the Logic App in use by the Windows Function App.

* `linux_fx_version` - Linux App Framework and version for the Logic App.

* `min_tls_version` - The minimum supported TLS version for the Logic App.

* `pre_warmed_instance_count` - The number of pre-warmed instances for this Logic App Only affects apps on the Premium plan.

* `runtime_scale_monitoring_enabled` - Should Runtime Scale Monitoring be enabled?. Only applicable to apps on the Premium plan.

* `use_32_bit_worker_process` - Should the Logic App run in 32 bit mode, rather than 64 bit mode?

* `vnet_route_all_enabled` - Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied.

* `websockets_enabled` - Should WebSockets be enabled?

---

A `cors` block supports the following:

* `allowed_origins` - A list of origins which should be able to make cross-origin calls.

* `support_credentials` - Are credentials supported?

---

A `ip_restriction` block supports the following:

* `ip_address` - The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - The Virtual Network Subnet ID used for this IP Restriction.

* `name` - The name for this IP Restriction.

* `priority` - The priority for this IP Restriction. Restrictions are enforced in priority order.

* `action` - Does this restriction `Allow` or `Deny` access for this IP range.

* `headers` - The `headers` block for this specific as a `ip_restriction` block as defined below.

---

A `scm_ip_restriction` block supports the following:

* `ip_address` - The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - The Virtual Network Subnet ID used for this IP Restriction.

* `name` - The name for this IP Restriction.

* `priority` - The priority for this IP Restriction. Restrictions are enforced in priority order.

* `action` - Does this restriction `Allow` or `Deny` access for this IP range.

* `headers` - The `headers` block for this specific `ip_restriction` as defined below.

---

A `headers` block supports the following:

* `x_azure_fdid` - A list of allowed Azure FrontDoor IDs in UUID notation.

* `x_fd_health_probe` - A list to allow the Azure FrontDoor health probe header.

* `x_forwarded_for` - A list of allowed 'X-Forwarded-For' IPs in CIDR notation.

* `x_forwarded_host` - A list of allowed 'X-Forwarded-Host' domains.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Workflow.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Web`: 2023-12-01
