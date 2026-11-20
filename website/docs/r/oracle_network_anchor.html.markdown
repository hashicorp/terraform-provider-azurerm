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
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "virtual_network" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "virtual_network_subnet" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.virtual_network.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      actions = [
        "Microsoft.Network/networkinterfaces/*",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
      name = "Oracle.Database/networkAttachments"
    }
  }
}

resource "azurerm_oracle_resource_anchor" "example" {
  name                = "example-resource-anchor"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_oracle_network_anchor" "example" {
  name                = "example-network-anchor"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  resource_anchor_id = azurerm_oracle_resource_anchor.example.id
  subnet_id          = azurerm_subnet.virtual_network_subnet.id
  zones              = ["2"]

  dns_forwarding_rule {
    domain_names          = ["app.example.com", "corp.example.com"]
    forwarding_ip_address = "10.0.3.4"
  }

  dns_listening_endpoint_allowed_cidrs = ["10.0.2.0/24"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Oracle Network Anchor. Changing this forces a new Oracle Network Anchor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Oracle Network Anchor should exist. Changing this forces a new Oracle Network Anchor to be created.

* `location` - (Required) The Azure Region where the Oracle Network Anchor should exist. Changing this forces a new Oracle Network Anchor to be created.

* `resource_anchor_id` - (Required) The ID of the Azure Resource Anchor. Changing this forces a new Oracle Network Anchor to be created.

* `subnet_id` - (Required) The ID of the delegated Azure subnet. Changing this forces a new resource to be created.

-> **Note:** The subnet must be delegated to `Oracle.Database/networkAttachments`. The virtual network is determined by this subnet ID.

* `zones` - (Required) A list of availability zones for the Network Anchor. Changing this forces a new Oracle Network Anchor to be created.

---

* `backup_subnet_cidr` - (Optional) The CIDR block for the Oracle Cloud Infrastructure backup subnet.

* `dns_forwarding_rule` - (Optional) One or more `dns_forwarding_rule` blocks as defined below. Configuring at least one rule enables the DNS forwarding endpoint so Oracle database instances can resolve Azure private domain names. Omitting the rules or specifying an empty list disables forwarding. Changing this forces a new resource to be created.

* `dns_listening_endpoint_allowed_cidrs` - (Optional) A set of CIDR blocks allowed to send requests to the DNS listening endpoint. Configuring at least one CIDR enables the DNS listening endpoint so Azure applications can resolve Oracle database domain names. Omitting the CIDRs or specifying an empty set disables listening. Changing this forces a new resource to be created.

* `dns_private_zone_replication_enabled` - (Optional) Whether to replicate database private DNS zones from Oracle Cloud Infrastructure to Azure during database creation. Defaults to `false`. Changing this forces a new resource to be created.

-> **Note:** DNS listening and private-zone replication provide overlapping ways to resolve Oracle database domain names from Azure. Choose the option that matches your DNS design. They can be enabled together. Private-zone replication occurs during database creation; it does not imply continuous DNS synchronization.

* `tags` - (Optional) A mapping of tags which should be assigned to the Oracle Network Anchor.

---

A `dns_forwarding_rule` block supports the following:

* `domain_names` - (Required) A set of domain names whose DNS queries should be forwarded.

* `forwarding_ip_address` - (Required) Forwarding IP address.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Oracle Network Anchor.

* `dns_forwarding_endpoint_ip_address` - A DNS forwarding endpoint IP address.

* `dns_forwarding_endpoint_nsg_rule_url` - A link to OCI console DNS Forwarding endpoint NSG rules.

* `dns_forwarding_rule_url` - A link to OCI console DNS Forwarding rules page.

* `dns_listening_endpoint_ip_address` - A DNS listening endpoint IP address.

* `dns_listening_endpoint_nsg_rule_url` - A link to OCI console DNS Listening endpoint NSG rules.

* `oci_vcn_dns_label` - The Oracle Cloud Infrastructure VCN DNS label.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Oracle Network Anchor.
* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle Network Anchor.
* `update` - (Defaults to 30 minutes) Used when updating the Oracle Network Anchor.
* `delete` - (Defaults to 30 minutes) Used when deleting the Oracle Network Anchor.

## Import

Oracle Network Anchors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_oracle_network_anchor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testRG/providers/Oracle.Database/networkAnchors/testNetworkAnchor1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Oracle.Database` - 2025-09-01
