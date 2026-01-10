---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_network_anchor"
description: |-
  Gets information about an existing Oracle Network Anchor.
---

# Data Source: azurerm_oracle_network_anchor

Use this data source to access information about an existing Oracle Network Anchor.

## Example Usage

```hcl
data "azurerm_oracle_network_anchor" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_oracle_network_anchor.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Oracle Network Anchor.

* `resource_group_name` - (Required) The name of the Resource Group where the Oracle Network Anchor exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Oracle Network Anchor.

* `cidr_block` - Delegated Azure subnet cidr block.

* `dns_forwarding_endpoint_ip_address` - A DNS forwarding endpoint IP address.

* `dns_forwarding_endpoint_nsg_rule_url` - A link to OCI console DNS Forwarding endpoint NSG rules.

* `dns_forwarding_rule_url` - A link to OCI console DNS Forwarding rules page.

* `dns_listening_endpoint_allowed_cidrs` - A list of CIDRs that are allowed to send requests to the DNS listening endpoint.

* `dns_listening_endpoint_ip_address` - A DNS listening endpoint IP address.

* `dns_listening_endpoint_nsg_rule_url` - A link to OCI console DNS Listening endpoint NSG rules.

* `location` - The Azure Region where the Oracle Network Anchor exists.

* `oci_backup_cidr_block` - Oracle Cloud Infrastructure backup subnet cidr block.

* `oci_subnet_id` - Oracle Cloud Infrastructure subnet OCID.

* `oci_vcn_dns_label` - Oracle Cloud Infrastructure DNS label. This is optional if DNS config is provided.

* `oci_vcn_id` - Oracle Cloud Infrastructure VCN OCID.

* `oracle_dns_forwarding_endpoint_enabled` - Whether the Oracle DNS forwarding endpoint is enabled.

* `oracle_dns_listening_endpoint_enabled` - Whether the Oracle DNS listening endpoint is enabled.

* `oracle_to_azure_dns_zone_sync_enabled` - Whether DNS zone sync from OCI to Azure is enabled.

* `provisioning_state` - Oracle Network Anchor provisioning state.

* `resource_anchor_id` - The ID of the corresponding Azure Resource Anchor.

* `subnet_id` - Delegated Azure subnet.

* `vnet_id` - Oracle Cloud Infrastructure VNET for network connectivity.

* `tags` - A mapping of tags assigned to the Oracle Network Anchor.

* `zones` - A list of availability `zones`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle Network Anchor.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
