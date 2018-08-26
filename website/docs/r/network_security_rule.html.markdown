---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_rule"
sidebar_current: "docs-azurerm-resource-network-security-rule"
description: |-
  Manages a Network Security Rule.

---

# azurerm_network_security_rule

Manages a Network Security Rule.

~> **NOTE on Network Security Groups and Network Security Rules:** Terraform currently
provides both a standalone [Network Security Rule resource](network_security_rule.html), and allows for Network Security Rules to be defined in-line within the [Network Security Group resource](network_security_group.html).
At this time you cannot use a Network Security Group with in-line Network Security Rules in conjunction with any Network Security Rule resources. Doing so will cause a conflict of rule settings and will overwrite rules.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_security_rule" "test" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  network_security_group_name = "${azurerm_network_security_group.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security rule. This needs to be unique across all Rules in the Network Security Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Network Security Rule. Changing this forces a new resource to be created.

* `network_security_group_name` - (Required) The name of the Network Security Group that we want to attach the rule to. Changing this forces a new resource to be created.

* `description` - (Optional) A description for this rule. Restricted to 140 characters.

* `protocol` - (Required) Network protocol this rule applies to. Possible values include `Tcp`, `Udp` or `*` (which matches both).

* `source_port_range` - (Optional) Source Port or Range. Integer or range between `0` and `65535` or `*` to match any. This is required if `source_port_ranges` is not specified.

* `source_port_ranges` - (Optional) List of source ports or port ranges. This is required if `source_port_range` is not specified.

* `destination_port_range` - (Optional) Destination Port or Range. Integer or range between `0` and `65535` or `*` to match any. This is required if `destination_port_ranges` is not specified.

* `destination_port_ranges` - (Optional) List of destination ports or port ranges. This is required if `destination_port_range` is not specified.

* `source_address_prefix` - (Optional) CIDR or source IP range or * to match any IP. Tags such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be used. This is required if `source_address_prefixes` is not specified.

* `source_address_prefixes` - (Optional) List of source address prefixes. Tags may not be used. This is required if `source_address_prefix` is not specified.

* `source_application_security_group_ids` - (Optional) A List of source Application Security Group ID's

* `destination_address_prefix` - (Optional) CIDR or destination IP range or * to match any IP. Tags such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be used. This is required if `destination_address_prefixes` is not specified.

* `destination_address_prefixes` - (Optional) List of destination address prefixes. Tags may not be used. This is required if `destination_address_prefix` is not specified.

* `destination_application_security_group_ids` - (Optional) A List of destination Application Security Group ID's

* `access` - (Required) Specifies whether network traffic is allowed or denied. Possible values are `Allow` and `Deny`.

* `priority` - (Required) Specifies the priority of the rule. The value can be between 100 and 4096. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.

* `direction` - (Required) The direction specifies if rule will be evaluated on incoming or outgoing traffic. Possible values are `Inbound` and `Outbound`.

## Attributes Reference

The following attributes are exported:

* `id` - The Network Security Rule ID.


## Import

Network Security Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_security_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkSecurityGroups/mySecurityGroup/securityRules/rule1
```
