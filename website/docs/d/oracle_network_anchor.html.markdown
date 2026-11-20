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
  name                = "example-network-anchor"
  resource_group_name = "example-resources"
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

* `backup_subnet_cidr` - The CIDR block for the Oracle Cloud Infrastructure backup subnet.

* `cidr_block` - The CIDR block of the delegated Azure client subnet.

* `dns_forwarding_endpoint_enabled` - Whether a DNS forwarding endpoint is enabled for Oracle database instances to resolve Azure private domain names.

* `dns_forwarding_endpoint_ip_address` - The DNS forwarding endpoint IP address.

* `dns_forwarding_endpoint_nsg_rule_url` - A link to OCI console DNS Forwarding endpoint NSG rules.

* `dns_forwarding_rule_url` - A link to OCI console DNS Forwarding rules page.

* `dns_listening_endpoint_allowed_cidrs` - A set of CIDR blocks allowed to send requests to the DNS listening endpoint.

* `dns_listening_endpoint_enabled` - Whether a DNS listening endpoint is enabled for Azure applications to resolve Oracle database domain names.

* `dns_listening_endpoint_ip_address` - The DNS listening endpoint IP address.

* `dns_listening_endpoint_nsg_rule_url` - A link to OCI console DNS Listening endpoint NSG rules.

* `dns_private_zone_replication_enabled` - Whether database private DNS zones are replicated from Oracle Cloud Infrastructure to Azure during database creation.

* `location` - The Azure Region where the Oracle Network Anchor exists.

* `oci_subnet_id` - Oracle Cloud Infrastructure subnet OCID.

* `oci_vcn_dns_label` - The Oracle Cloud Infrastructure VCN DNS label.

* `oci_vcn_id` - Oracle Cloud Infrastructure VCN OCID.

* `provisioning_state` - Oracle Network Anchor provisioning state.

* `resource_anchor_id` - The ID of the corresponding Azure Resource Anchor.

* `subnet_id` - The ID of the delegated Azure subnet.

* `tags` - A mapping of tags assigned to the Oracle Network Anchor.

* `virtual_network_id` - The Azure virtual network ID for network connectivity.

* `zones` - A list of availability `zones`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle Network Anchor.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
