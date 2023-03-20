---
subcategory: "PostgreSQL HyperScale"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_hyperscale_firewall_rule"
description: |-
  Manages a PostgreSQL HyperScale Firewall Rule.
---

# azurerm_postgresql_hyperscale_firewall_rule

Manages a PostgreSQL HyperScale Firewall Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_hyperscale_cluster" "example" {
  name                = "example-postgresqlhscsg"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_postgresql_hyperscale_firewall_rule" "example" {
  name             = "example-postgresqlhscfwr"
  server_group_id  = azurerm_postgresql_hyperscale_cluster.example.id
  end_ip_address   = "10.0.17.62"
  start_ip_address = "10.0.17.64"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the PostgreSQL HyperScale Firewall Rule. Changing this forces a new resource to be created.

* `server_group_id` - (Required) The ID of the PostgreSQL HyperScale Cluster. Changing this forces a new resource to be created.

* `end_ip_address` - (Required) The end IP address of the PostgreSQL HyperScale Firewall Rule. Must be IPv4 format.

* `start_ip_address` - (Required) The start IP address of the PostgreSQL HyperScale Firewall Rule. Must be IPv4 format.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL HyperScale Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL HyperScale Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL HyperScale Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL HyperScale Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL HyperScale Firewall Rule.

## Import

PostgreSQL HyperScale Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_hyperscale_firewall_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/cluster1/firewallRules/firewallRule1
```
