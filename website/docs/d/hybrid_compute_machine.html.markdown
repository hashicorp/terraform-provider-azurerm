---
subcategory: "Hybrid Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_hybrid_compute_machine"
description: |-
  Gets information about an existing Hybrid Compute.
---

# Data Source: azurerm_hybrid_compute_machine

Use this data source to access information about an existing Hybrid Compute.

## Example Usage

```hcl
data "azurerm_hybrid_compute_machine" "example" {
  name = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_hybrid_compute_machine.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Hybrid Compute. Changing this forces a new Hybrid Compute to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Hybrid Compute exists. Changing this forces a new Hybrid Compute to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Hybrid Compute.

* `ad_fqdn` - TODO.

* `agent_configuration` - A `agent_configuration` block as defined below.

* `agent_version` - TODO.

* `client_public_key` - TODO.

* `cloud_metadata` - A `cloud_metadata` block as defined below.

* `detected_properties` - A `detected_properties` block as defined below.

* `display_name` - TODO.

* `dns_fqdn` - TODO.

* `domain_name` - TODO.

* `error_details` - A `error_details` block as defined below.

* `extensions` - A `extensions` block as defined below.

* `identity` - A `identity` block as defined below.

* `last_status_change` - TODO.

* `location` - The Azure Region where the Hybrid Compute exists.

* `location_data` - A `location_data` block as defined below.

* `machine_fqdn` - TODO.

* `mssql_discovered` - TODO.

* `os_name` - TODO.

* `os_profile` - A `os_profile` block as defined below.

* `os_sku` - TODO.

* `os_type` - TODO.

* `os_version` - TODO.

* `parent_cluster_resource_id` - The ID of the TODO.

* `private_link_scope_resource_id` - The ID of the TODO.

* `service_statuses` - A `service_statuses` block as defined below.

* `status` - TODO.

* `tags` - A mapping of tags assigned to the Hybrid Compute.

* `vm_id` - The ID of the TODO.

* `vm_uuid` - TODO.

---

A `additional_info` block exports the following:

* `info` - TODO.

* `type` - TODO.

---

A `agent_configuration` block exports the following:

* `extensions_allow_list` - A `extensions_allow_list` block as defined below.

* `extensions_block_list` - A `extensions_block_list` block as defined below.

* `extensions_enabled` - Is the TODO enabled?

* `guest_configuration_enabled` - Is the TODO enabled?

* `incoming_connections_ports` - A `incoming_connections_ports` block as defined below.

* `proxy_bypass` - A `proxy_bypass` block as defined below.

* `proxy_url` - TODO.

---

A `cloud_metadata` block exports the following:

* `provider` - TODO.

---

A `error_details` block exports the following:

* `additional_info` - A `additional_info` block as defined above.

* `code` - TODO.

* `message` - TODO.

* `target` - TODO.

---

A `extension_service` block exports the following:

* `startup_type` - TODO.

* `status` - TODO.

---

A `extensions` block exports the following:

* `name` - The name of this TODO.

* `status` - A `status` block as defined below.

* `type` - TODO.

* `type_handler_version` - TODO.

---

A `extensions_allow_list` block exports the following:

* `publisher` - TODO.

* `type` - TODO.

---

A `extensions_block_list` block exports the following:

* `publisher` - TODO.

* `type` - TODO.

---

A `guest_configuration_service` block exports the following:

* `startup_type` - TODO.

* `status` - TODO.

---

A `identity` block exports the following:

* `principal_id` - The ID of the TODO.

* `tenant_id` - The ID of the TODO.

* `type` - TODO.

---

A `linux_configuration` block exports the following:

* `patch_settings` - A `patch_settings` block as defined below.

---

A `location_data` block exports the following:

* `city` - TODO.

* `country_or_region` - TODO.

* `district` - TODO.

* `name` - The name of this TODO.

---

A `os_profile` block exports the following:

* `computer_name` - TODO.

* `linux_configuration` - A `linux_configuration` block as defined above.

* `windows_configuration` - A `windows_configuration` block as defined below.

---

A `patch_settings` block exports the following:

* `assessment_mode` - TODO.

* `patch_mode` - TODO.

---

A `service_statuses` block exports the following:

* `extension_service` - A `extension_service` block as defined above.

* `guest_configuration_service` - A `guest_configuration_service` block as defined above.

---

A `status` block exports the following:

* `code` - TODO.

* `display_status` - TODO.

* `level` - TODO.

* `message` - TODO.

* `time` - TODO.

---

A `windows_configuration` block exports the following:

* `patch_settings` - A `patch_settings` block as defined above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Hybrid Compute.