---
subcategory: "Hybrid Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_arc_machine"
description: |-
  Gets information about an existing Azure Arc machine.
---

# Data Source: azurerm_arc_machine

Use this data source to access information about an existing Azure Arc machine.

## Example Usage

```hcl
data "azurerm_arc_machine" "example" {
  name                = "existing-hcmachine"
  resource_group_name = "existing-rg"
}

output "id" {
  value = data.azurerm_arc_machine.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Azure Arc machine.

* `resource_group_name` - (Required) The name of the Resource Group where the Hybrid Compute exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Arc machine.

* `active_directory_fqdn` - Specifies the AD fully qualified display name.

* `agent` - A `agent` block as defined below.

* `agent_version` - The Azure Arc machine agent full version.

* `client_public_key` - Public Key that the client provides to be used during initial resource onboarding.

* `cloud_metadata` - A `cloud_metadata` block as defined below.

* `detected_properties` - A `detected_properties` block as defined below.

* `display_name` - Specifies the Azure Arc machine display name.

* `dns_fqdn` - Specifies the DNS fully qualified display name.

* `domain_name` - Specifies the Windows domain name.

* `error_details` - A `error_details` block as defined below.

* `identity` - A `identity` block as defined below.

* `last_status_change_time` - The time of the last status change.

* `location` - The Azure Region where the Azure Arc machine exists.

* `location_data` - A `location_data` block as defined below.

* `machine_fqdn` - Specifies the Azure Arc machine fully qualified display name.

* `mssql_discovered` - Specifies whether any MS SQL instance is discovered on the machine.

* `os_name` - The Operating System running on the Azure Arc machine.

* `os_profile` - A `os_profile` block as defined below.

* `os_sku` - Specifies the Operating System product SKU.

* `os_type` - The type of Operating System. Possible values are `windows` and `linux`.

* `os_version` - The version of Operating System running on the Azure Arc machine.

* `parent_cluster_resource_id` - The resource id of the parent cluster (Azure HCI) this machine is assigned to, if any.

* `private_link_scope_resource_id` - The resource id of the parent cluster (Azure HCI) this machine is assigned to, if any.

* `service_status` - A `service_status` block as defined below.

* `status` - The status of the Azure Arc machine agent.

* `tags` - A mapping of tags assigned to the Hybrid Compute.

* `vm_id` - Specifies the Azure Arc machine unique ID.

* `vm_uuid` - Specifies the Arc Machine's unique SMBIOS ID.

---

A `agent` block exports the following:

* `extensions_allow_list` - A `extensions_allow_list` block as defined below.

* `extensions_block_list` - A `extensions_block_list` block as defined below.

* `extensions_enabled` - Specifies whether the extension service is enabled or disabled.

* `guest_configuration_enabled` - Specified whether the guest configuration service is enabled or disabled.

* `incoming_connections_ports` - Specifies the list of ports that the agent will be able to listen on.

* `proxy_bypass` - List of service names which should not use the specified proxy server.

* `proxy_url` - Specifies the URL of the proxy to be used.

---

A `cloud_metadata` block exports the following:

* `provider` - Specifies the cloud provider. For example `Azure`, `AWS` and `GCP`.

---

A `error_details` block exports the following:

* `additional_info` - A `additional_info` block as defined above.

* `code` - The error code.

* `message` - The error message.

* `target` - The error target.

---

A `extension_service` block exports the following:

* `startup_type` - The behavior of the service when the Arc-enabled machine starts up.

* `status` - The current status of the service.

---

A `extensions_allow_list` block exports the following:

* `publisher` - Publisher of the extension.

* `type` - Type of the extension.

---

A `extensions_block_list` block exports the following:

* `publisher` - Publisher of the extension.

* `type` - Type of the extension.

---

A `guest_configuration_service` block exports the following:

* `startup_type` - The behavior of the service when the Arc-enabled machine starts up.

* `status` - The current status of the service.

---

A `identity` block exports the following:

* `principal_id` - The principal ID of resource identity.

* `tenant_id` - The tenant ID of resource.

* `type` - The identity type.

---

A `linux` block exports the following:

* `patch` - A `patch` block as defined below.

---

A `location_data` block exports the following:

* `city` - The city or locality where the resource is located.

* `country_or_region` - The country or region where the resource is located.

* `district` - The district, state, or province where the resource is located.

* `name` - A canonical name for the geographic or physical location.

---

A `os_profile` block exports the following:

* `computer_name` - Specifies the host OS name of the Azure Arc machine.

* `linux` - A `linux` block as defined above.

* `windows` - A `windows` block as defined below.

---

A `patch` block exports the following:

* `assessment_mode` - Specifies the assessment mode.

* `patch_mode` - Specifies the patch mode.

---

A `service_status` block exports the following:

* `extension_service` - A `extension_service` block as defined above.

* `guest_configuration_service` - A `guest_configuration_service` block as defined above.

---

A `windows` block exports the following:

* `patch` - A `patch` block as defined above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Arc machine.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.HybridCompute`: 2022-11-10
