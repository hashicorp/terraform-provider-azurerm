---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_firewall_rules"
description: |-
  Manages a group of PostgreSQL Flexible Server Firewall Rules.
---

# azurerm_postgresql_flexible_server_firewall_rules

Manages a group of PostgreSQL Flexible Server Firewall Rules.

~> **NOTE:** It's possible to define Flexible Server Firewall Rules individually with [the `postgresql_flexible_server_firewall_rule` resource](postgresql_flexible_server_firewall_rule.html) or in groups with [the `postgresql_flexible_server_firewall_rules` resource](postgresql_flexible_server_firewall_rules.html). But, it's not possible to use both methods to manage Firewall Rules within a PostgreSQL Flexible Server, since there will be conflicts.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-psqlflexibleserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  version                = "12"
  administrator_login    = "psqladmin"
  administrator_password = "H@Sh1CoR3!"

  storage_mb = 32768

  sku_name = "GP_Standard_D4s_v3"
}

resource "azurerm_postgresql_flexible_server_firewall_rules" "example" {
  server_id        = azurerm_postgresql_flexible_server.example.id

  firewall_rule {
    name             = "ruleA"
    start_ip_address = "40.112.8.12"
    end_ip_address   = "40.112.8.12"
  }

  firewall_rule {
    name             = "ruleB"
    start_ip_address = "40.112.8.200"
    end_ip_address   = "40.112.8.205"
  }

  dynamic "firewall_rule" {
    for_each = var.firewall_rules
    content = {
      name             = firewall_rule.value.name
      start_ip_address = firewall_rule.value.key
      end_ip_address   = firewall_rule.value.key
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the PostgreSQL Flexible Server from which to create these PostgreSQL Flexible Server Firewall Rules. Changing this forces a new resource to be created.

* `firewall_rule` - (Optional) A `firewall_rule` object as defined below. Changing this forces a new resource to be created.

---

The `firewall_rule` block supports the following:

* `name` - (Required) The name which should be used for this PostgreSQL Flexible Server Firewall Rule. 

* `start_ip_address` - (Required) The Start IP Address associated with this PostgreSQL Flexible Server Firewall Rule.

* `end_ip_address` - (Required) The End IP Address associated with this PostgreSQL Flexible Server Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL Flexible Server Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Flexible Server Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL Flexible Server Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL Flexible Server Firewall Rule.
