---
subcategory: "Mongo Cluster"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mongo_cluster_firewall_rule"
description: |-
  Manages a Mongo Cluster Firewall Rule.
---

# azurerm_mongo_cluster_firewall_rule

Manages a Mongo Cluster Firewall Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mongo_cluster" "example" {
  name                   = "example-mongocluster"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_username = "adminuser"
  administrator_password = "P@ssw0rd1234!"
  shard_count            = 1
  compute_tier           = "M30"
  high_availability_mode = "Disabled"
  storage_size_in_gb     = 32
  version                = "7.0"
}

resource "azurerm_mongo_cluster_firewall_rule" "example" {
  name             = "example-firewall-rule"
  mongo_cluster_id = azurerm_mongo_cluster.example.id
  start_ip_address = "10.0.0.1"
  end_ip_address   = "10.0.0.255"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Mongo Cluster Firewall Rule. Changing this forces a new resource to be created.

* `mongo_cluster_id` - (Required) The ID of the Mongo Cluster. Changing this forces a new resource to be created.

* `end_ip_address` - (Required) The end IP address of the Mongo Cluster Firewall Rule.

* `start_ip_address` - (Required) The start IP address of the Mongo Cluster Firewall Rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mongo Cluster Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Mongo Cluster Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mongo Cluster Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Mongo Cluster Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Mongo Cluster Firewall Rule.

## Import

Mongo Cluster Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mongo_cluster_firewall_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DocumentDB/mongoClusters/cluster1/firewallRules/rule1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DocumentDB` - 2025-09-01
