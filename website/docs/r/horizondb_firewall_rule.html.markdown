---
subcategory: "HorizonDB"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_horizondb_firewall_rule"
description: |-
  Manages a HorizonDB Firewall Rule.
---

# azurerm_horizondb_firewall_rule

Manages a HorizonDB Firewall Rule.

## Example Usage

```hcl
variable "horizondb_administrator_password" {
  type        = string
  description = "The PostgreSQL administrator password for the HorizonDB Cluster."
  sensitive   = true
}

resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "East US 2"
}

resource "azurerm_horizondb_cluster" "example" {
  name                = "example-horizondb-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  postgresql_version  = "18"

  administrator_login    = "exampleadmin"
  administrator_password = var.horizondb_administrator_password

  read_write_pool {
    name                   = "example-pool"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }
}

resource "azurerm_horizondb_firewall_rule" "example" {
  name              = "example-horizondb-firewall-rule"
  horizondb_pool_id = azurerm_horizondb_cluster.example.read_write_pool[0].id
  start_ip_address  = "192.0.2.1"
  end_ip_address    = "192.0.2.10"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the HorizonDB Firewall Rule. Changing this forces a new resource to be created.

~> **Note:** The name must contain between 2 and 80 alphanumeric characters or hyphens and must start and end with an alphanumeric character.

* `horizondb_pool_id` - (Required) The ID of the HorizonDB Pool. Changing this forces a new resource to be created.

* `end_ip_address` - (Required) The last IPv4 address in the allowed range.

~> **Note:** `end_ip_address` must be greater than or equal to `start_ip_address`.

* `start_ip_address` - (Required) The first IPv4 address in the allowed range.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the HorizonDB Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HorizonDB Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the HorizonDB Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the HorizonDB Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the HorizonDB Firewall Rule.

## Import

A HorizonDB Firewall Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_horizondb_firewall_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HorizonDb/clusters/cluster1/pools/pool1/firewallRules/firewallRule1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.HorizonDb` - 2026-10-01-preview
