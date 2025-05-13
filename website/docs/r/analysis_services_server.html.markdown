---
subcategory: "Analysis Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_analysis_services_server"
description: |-
  Manages an Analysis Services Server.
---

# azurerm_analysis_services_server

Manages an Analysis Services Server.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "analysis-services-server-test"
  location = "West Europe"
}

resource "azurerm_analysis_services_server" "server" {
  name                     = "analysisservicesserver"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  sku                      = "S0"
  admin_users              = ["myuser@domain.tld"]
  power_bi_service_enabled = true

  ipv4_firewall_rule {
    name        = "myRule1"
    range_start = "210.117.252.0"
    range_end   = "210.117.252.255"
  }

  tags = {
    abc = 123
  }
}
```

-> **Note:** The server resource will automatically be started and stopped during an update if it is in `paused` state.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Analysis Services Server. Only lowercase Alphanumeric characters allowed, starting with a letter. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Analysis Services Server exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Analysis Services Server should be exist. Changing this forces a new resource to be created.

* `sku` - (Required) SKU for the Analysis Services Server. Possible values are: `D1`, `B1`, `B2`, `S0`, `S1`, `S2`, `S4`, `S8`, `S9`, `S8v2` and `S9v2`.

* `admin_users` - (Optional) List of email addresses of admin users.

* `querypool_connection_mode` - (Optional) Controls how the read-write server is used in the query pool. If this value is set to `All` then read-write servers are also used for queries. Otherwise with `ReadOnly` these servers do not participate in query operations. Defaults to `All`.

* `backup_blob_container_uri` - (Optional) URI and SAS token for a blob container to store backups.

* `power_bi_service_enabled` - (Optional) Indicates if the Power BI service is allowed to access or not.

* `ipv4_firewall_rule` - (Optional) One or more `ipv4_firewall_rule` block(s) as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `ipv4_firewall_rule` block supports the following:

* `name` - (Required) Specifies the name of the firewall rule.

* `range_start` - (Required) Start of the firewall rule range as IPv4 address.

* `range_end` - (Required) End of the firewall rule range as IPv4 address.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Analysis Services Server.

* `server_full_name` - The full name of the Analysis Services Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Analysis Services Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Analysis Services Server.
* `update` - (Defaults to 30 minutes) Used when updating the Analysis Services Server.
* `delete` - (Defaults to 30 minutes) Used when deleting the Analysis Services Server.

## Import

Analysis Services Server can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_analysis_services_server.server /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AnalysisServices/servers/server1
```
