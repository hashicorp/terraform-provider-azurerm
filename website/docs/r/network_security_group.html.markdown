---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_group"
description: |-
  Manages a network security group that contains a list of network security rules. Network security groups enable inbound or outbound traffic to be enabled or denied.

---

# azurerm_network_security_group

Manages a network security group that contains a list of network security rules.  Network security groups enable inbound or outbound traffic to be enabled or denied.

~> **Note:** Terraform currently
provides both a standalone [Network Security Rule resource](network_security_rule.html), and allows for Network Security Rules to be defined in-line within the [Network Security Group resource](network_security_group.html).
At this time you cannot use a Network Security Group with in-line Network Security Rules in conjunction with any Network Security Rule resources. Doing so will cause a conflict of rule settings and will overwrite rules.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the network security group. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the resource group in which to create the network security group. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `security_rule` - (Optional) List of `security_rule` objects representing security rules, as defined below.

-> **Note:** Since `security_rule` can be configured both inline and via the separate `azurerm_network_security_rule` resource, we have to explicitly set it to empty slice (`[]`) to remove it.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `security_rule` block support:

* `name` - (Required) The name of the security rule.

* `description` - (Optional) A description for this rule. Restricted to 140 characters.

* `protocol` - (Required) Network protocol this rule applies to. Possible values include `Tcp`, `Udp`, `Icmp`, `Esp`, `Ah` or `*` (which matches all).

* `source_port_range` - (Optional) Source Port or Range. Integer or range between `0` and `65535` or `*` to match any. This is required if `source_port_ranges` is not specified.

* `source_port_ranges` - (Optional) List of source ports or port ranges. This is required if `source_port_range` is not specified.

* `destination_port_range` - (Optional) Destination Port or Range. Integer or range between `0` and `65535` or `*` to match any. This is required if `destination_port_ranges` is not specified.

* `destination_port_ranges` - (Optional) List of destination ports or port ranges. This is required if `destination_port_range` is not specified.

* `source_address_prefix` - (Optional) CIDR or source IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used. This is required if `source_address_prefixes` is not specified.

* `source_address_prefixes` - (Optional) List of source address prefixes. Tags may not be used. This is required if `source_address_prefix` is not specified.

* `source_application_security_group_ids` - (Optional) A List of source Application Security Group IDs

* `destination_address_prefix` - (Optional) CIDR or destination IP range or * to match any IP. Tags such as `VirtualNetwork`, `AzureLoadBalancer` and `Internet` can also be used. This is required if `destination_address_prefixes` is not specified.

* `destination_address_prefixes` - (Optional) List of destination address prefixes. Tags may not be used. This is required if `destination_address_prefix` is not specified.

* `destination_application_security_group_ids` - (Optional) A List of destination Application Security Group IDs

* `access` - (Required) Specifies whether network traffic is allowed or denied. Possible values are `Allow` and `Deny`.

* `priority` - (Required) Specifies the priority of the rule. The value can be between 100 and 4096. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.

* `direction` - (Required) The direction specifies if rule will be evaluated on incoming or outgoing traffic. Possible values are `Inbound` and `Outbound`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Network Security Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Security Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Group.
* `update` - (Defaults to 30 minutes) Used when updating the Network Security Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Security Group.

## Import

Network Security Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_security_group.group1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/networkSecurityGroups/mySecurityGroup
```
