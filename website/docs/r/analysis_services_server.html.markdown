---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_analysis_services_server"
sidebar_current: "docs-azurerm-resource-analysis_services_server-x"
description: |-
  Manages an Analysis Services Server.
---

# azurerm_analysis_services_server

Manages an Analysis Services Server.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "analysis-services-server-test"
  location = "northeurope"
}

resource "azurerm_analysis_services_server" "server" {
  name                      = "analysisservicesserver"
  location                  = "northeurope"
  resource_group_name       = "${azurerm_resource_group.rg.name}"
  sku                       = "S0"
  admin_users               = ["myuser@domain.tld"]
  
  tags {
    abc = 123
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Analysis Services Server. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Analysis Services Server exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Analysis Services Server should be exist. Changing this forces a new resource to be created.

* `sku` - (Required) SKU for the Analysis Services Server. Possible values are: `D1`, `B1`, `B2`, `S0`, `S1`, `S2`, `S4`, `S8` and `S9`

* `admin_users` - (Optional) List of email addresses of admin users.

* `gateway_resource_id` - (Optional) ID of the Gateway resource to be associated with the server.

* `querypool_connection_mode` - (Optional) Controls how the read-write server is used in the query pool. If this values is set to `All` then read-write servers are also used for queries. Otherwise with `ReadOnly` theses servers do not participate in query operations.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Analysis Services Server.

## Import

Analysis Services Server can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_anaylsis_services_server.server /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AnalysisServices/server/server1
```