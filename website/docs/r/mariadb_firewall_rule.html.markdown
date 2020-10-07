---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mariadb_firewall_rule"
description: |-
  Manages a Firewall Rule for a MariaDB Server.
---

# azurerm_mariadb_firewall_rule

Manages a Firewall Rule for a MariaDB Server

## Example Usage (Single IP Address)

```hcl
resource "azurerm_mariadb_firewall_rule" "example" {
  name                = "test-rule"
  resource_group_name = "test-rg"
  server_name         = "test-server"
  start_ip_address    = "40.112.8.12"
  end_ip_address      = "40.112.8.12"
}
```

## Example Usage (IP Range)

```hcl
resource "azurerm_mariadb_firewall_rule" "example" {
  name                = "test-rule"
  resource_group_name = "test-rg"
  server_name         = "test-server"
  start_ip_address    = "40.112.0.0"
  end_ip_address      = "40.112.255.255"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Firewall Rule. Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MariaDB Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MariaDB Server exists. Changing this forces a new resource to be created.

* `start_ip_address` - (Required) Specifies the Start IP Address associated with this Firewall Rule. Changing this forces a new resource to be created.

* `end_ip_address` - (Required) Specifies the End IP Address associated with this Firewall Rule. Changing this forces a new resource to be created.

-> **NOTE:** The Azure feature `Allow access to Azure services` can be enabled by setting `start_ip_address` and `end_ip_address` to `0.0.0.0` which ([is documented in the Azure API Docs](https://docs.microsoft.com/en-us/rest/api/sql/firewallrules/createorupdate)).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MariaDB Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MariaDB Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the MariaDB Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the MariaDB Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the MariaDB Firewall Rule.

## Import

MariaDB Firewall rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mariadb_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMariaDB/servers/server1/firewallRules/rule1
```
