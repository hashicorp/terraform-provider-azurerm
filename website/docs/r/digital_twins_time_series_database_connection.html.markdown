---
subcategory: "Digital Twins"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_digital_twins_time_series_database_connection"
description: |-
  Manages a Digital Twins Time Series Database Connection.
---

# azurerm_digital_twins_time_series_database_connection

Manages a Digital Twins Time Series Database Connection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_digital_twins_instance" "example" {
  name                = "example-DT"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "exampleEventHubNamespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "example" {
  name                = "exampleEventHub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 7
}

resource "azurerm_eventhub_consumer_group" "example" {
  name                = "example-consumergroup"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_kusto_cluster" "example" {
  name                = "examplekc"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "example-kusto-database"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name
}

resource "azurerm_role_assignment" "database_contributor" {
  scope                = azurerm_kusto_database.example.id
  principal_id         = azurerm_digital_twins_instance.example.identity.0.principal_id
  role_definition_name = "Contributor"
}

resource "azurerm_role_assignment" "eventhub_data_owner" {
  scope                = azurerm_eventhub.example.id
  principal_id         = azurerm_digital_twins_instance.example.identity.0.principal_id
  role_definition_name = "Azure Event Hubs Data Owner"
}

resource "azurerm_kusto_database_principal_assignment" "example" {
  name                = "dataadmin"
  resource_group_name = azurerm_resource_group.example.name
  cluster_name        = azurerm_kusto_cluster.example.name
  database_name       = azurerm_kusto_database.example.name

  tenant_id      = azurerm_digital_twins_instance.example.identity.0.tenant_id
  principal_id   = azurerm_digital_twins_instance.example.identity.0.principal_id
  principal_type = "App"
  role           = "Admin"
}

resource "azurerm_digital_twins_time_series_database_connection" "example" {
  name                            = "example-connection"
  digital_twins_id                = azurerm_digital_twins_instance.example.id
  eventhub_name                   = azurerm_eventhub.example.name
  eventhub_namespace_id           = azurerm_eventhub_namespace.example.id
  eventhub_namespace_endpoint_uri = "sb://${azurerm_eventhub_namespace.example.name}.servicebus.windows.net"
  eventhub_consumer_group_name    = azurerm_eventhub_consumer_group.example.name
  kusto_cluster_id                = azurerm_kusto_cluster.example.id
  kusto_cluster_uri               = azurerm_kusto_cluster.example.uri
  kusto_database_name             = azurerm_kusto_database.example.name
  kusto_table_name                = "exampleTable"

  depends_on = [
    azurerm_role_assignment.database_contributor,
    azurerm_role_assignment.eventhub_data_owner,
    azurerm_kusto_database_principal_assignment.example
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Digital Twins Time Series Database Connection. Changing this forces a new resource to be created.

* `digital_twins_id` - (Required) The ID of the Digital Twins. Changing this forces a new resource to be created.

* `eventhub_name` - (Required) Name of the Event Hub. Changing this forces a new resource to be created.

* `eventhub_namespace_endpoint_uri` - (Required) URI of the Event Hub Namespace. Changing this forces a new resource to be created.

* `eventhub_namespace_id` - (Required) The ID of the Event Hub Namespace. Changing this forces a new resource to be created.

* `kusto_cluster_id` - (Required) The ID of the Kusto Cluster. Changing this forces a new resource to be created.

* `kusto_cluster_uri` - (Required) URI of the Kusto Cluster. Changing this forces a new resource to be created.

* `kusto_database_name` - (Required) Name of the Kusto Database. Changing this forces a new resource to be created.

---

* `eventhub_consumer_group_name` - (Optional) Name of the Event Hub Consumer Group. Changing this forces a new resource to be created. Defaults to `$Default`.

* `kusto_table_name` - (Optional) Name of the Kusto Table. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Digital Twins Time Series Database Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Digital Twins Time Series Database Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Digital Twins Time Series Database Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Digital Twins Time Series Database Connection.

## Import

Digital Twins Time Series Database Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_digital_twins_time_series_database_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DigitalTwins/digitalTwinsInstances/dt1/timeSeriesDatabaseConnections/connection1
```
