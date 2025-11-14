---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_oracle_network_anchor"
description: |-
  Manages an Oracle Network Anchor.
---

# azurerm_oracle_network_anchor

Manages an Oracle Network Anchor.

## Example Usage

```hcl
resource "azurerm_oracle_network_anchor" "example" {
  name                = "example-network-anchor"
  resource_group_name = "example-resource-group"
  location            = "West Europe"
  resource_anchor_id  = "example-azure-resource-anchor-id"
  subnet_id           = "example-azure-delegated-subnet-id"
  zones               = ["2"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Oracle Network Anchor. Changing this forces a new Oracle Network Anchor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Oracle Network Anchor should exist. Changing this forces a new Oracle Network Anchor to be created.

* `location` - (Required) The Azure Region where the Oracle Network Anchor should exist. Changing this forces a new Oracle Network Anchor to be created.

* `resource_anchor_id` - (Required) The ID of the Azure Resource Anchor. Changing this forces a new Oracle Network Anchor to be created.

* `subnet_id` - (Required) The ID of the Delegated Azure subnet. Changing this forces a new Oracle Network Anchor to be created.

* `zones` - (Required) A list of availability zones for the Network Anchor. Changing this forces a new Oracle Network Anchor to be created.

---

* `dns_forwarding_rule` - (Optional) One or more `dns_forwarding_rule` blocks as defined below. Changing this forces a new Oracle Network Anchor to be created.

* `dns_listening_endpoint_allowed_cidrs` - (Optional) A list of CIDRs that are allowed to send requests to the DNS listening endpoint. Changing this forces a new Oracle Network Anchor to be created.

* `oci_backup_cidr_block` - (Optional) Oracle Cloud Infrastructure backup subnet cidr block.

* `oci_vcn_dns_label` - (Optional) The DNS label for the Oracle Cloud Infrastructure VCN. If not specified, defaults to the Network Anchor name. Changing this forces a new Oracle Network Anchor to be created.

* `oracle_dns_forwarding_endpoint_enabled` - (Optional) Whether to enable the Oracle DNS forwarding endpoint.

* `oracle_dns_listening_endpoint_enabled` - (Optional) Whether to enable the Oracle DNS listening endpoint.

* `oracle_to_azure_dns_zone_sync_enabled` - (Optional) Whether to enable DNS zone sync from OCI to Azure.

* `tags` - (Optional) A mapping of tags which should be assigned to the Oracle Network Anchor.

---

A `dns_forwarding_rule` block supports the following:

* `domain_names` - (Required) Comma-separated domain names.

* `forwarding_ip_address` - (Required) Forwarding ip address.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Oracle Network Anchor.

* `dns_forwarding_endpoint_ip_address` - A DNS forwarding endpoint IP address.

* `dns_forwarding_endpoint_nsg_rule_url` - A link to OCI console DNS Forwarding endpoint NSG rules.

* `dns_forwarding_rule_url` - A link to OCI console DNS Forwarding rules page.

* `dns_listening_endpoint_ip_address` - A DNS listening endpoint IP address

* `dns_listening_endpoint_nsg_rule_url` - A link to OCI console DNS Listening endpoint NSG rules

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Oracle Network Anchor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle Network Anchor.
* `update` - (Defaults to 30 minutes) Used when updating the Oracle Network Anchor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Oracle Network Anchor.

## Import

Oracle Network Anchors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_network_anchor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testRG/providers/Oracle.Database/networkAnchors/testNetworAnchor1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
