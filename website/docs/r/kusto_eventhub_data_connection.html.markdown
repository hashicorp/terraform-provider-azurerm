---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_eventhub_data_connection"
description: |-
  Manages Kusto / Data Explorer EventHub Data Connection
---

# azurerm_kusto_eventhub_data_connection

Manages a Kusto (also known as Azure Data Explorer) EventHub Data Connection

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "my-kusto-rg"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "kustocluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}

resource "azurerm_kusto_database" "database" {
  name                = "my-kusto-database"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
  hot_cache_period    = "P7D"
  soft_delete_period  = "P31D"
}

resource "azurerm_eventhub_namespace" "eventhub_ns" {
  name                = "my-eventhub-ns"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "eventhub" {
  name                = "my-eventhub"
  namespace_name      = azurerm_eventhub_namespace.eventhub_ns.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "consumer_group" {
  name                = "my-eventhub-consumergroup"
  namespace_name      = azurerm_eventhub_namespace.eventhub_ns.name
  eventhub_name       = azurerm_eventhub.eventhub.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_kusto_eventhub_data_connection" "eventhub_connection" {
  name                = "my-kusto-eventhub-data-connection"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
  database_name       = azurerm_kusto_database.database.name

  eventhub_id    = azurerm_eventhub.eventhub.id
  consumer_group = azurerm_eventhub_consumer_group.consumer_group.name

  table_name        = "my-table"         #(Optional)
  mapping_rule_name = "my-table-mapping" #(Optional)
  data_format       = "JSON"             #(Optional)
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto EventHub Data Connection to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Database should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Database should exist. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the Kusto Cluster this data connection will be added to. Changing this forces a new resource to be created.

* `compression` - (Optional) Specifies compression type for the connection. Allowed values: `GZip` and `None`. Defaults to `None`. Changing this forces a new resource to be created.

* `database_name` - (Required) Specifies the name of the Kusto Database this data connection will be added to. Changing this forces a new resource to be created.

* `eventhub_id` - (Required) Specifies the resource id of the EventHub this data connection will use for ingestion. Changing this forces a new resource to be created.

* `event_system_properties` - (Optional) Specifies a list of system properties for the Event Hub.

* `consumer_group` - (Required) Specifies the EventHub consumer group this data connection will use for ingestion. Changing this forces a new resource to be created.

* `table_name` - (Optional) Specifies the target table name used for the message ingestion. Table must exist before resource is created.

* `identity_id` - (Optional) The resource ID of a managed identity (system or user assigned) to be used to authenticate with event hub.

* `mapping_rule_name` - (Optional) Specifies the mapping rule used for the message ingestion. Mapping rule must exist before resource is created.

* `data_format` - (Optional) Specifies the data format of the EventHub messages. Allowed values: `APACHEAVRO`, `AVRO`, `CSV`, `JSON`, `MULTIJSON`, `ORC`, `PARQUET`, `PSV`, `RAW`, `SCSV`, `SINGLEJSON`, `SOHSV`, `TSVE`, `TSV`, `TXT`, and `W3CLOGFILE`.

* `database_routing_type` - (Optional) Indication for database routing information from the data connection, by default only database routing information is allowed. Allowed values: `Single`, `Multi`. Changing this forces a new resource to be created. Defaults to `Single`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kusto EventHub Data Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto EventHub Data Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto EventHub Data Connection.
* `update` - (Defaults to 1 hour) Used when updating the Kusto EventHub Data Connection.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto EventHub Data Connection.

## Import

Kusto EventHub Data Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_eventhub_data_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/dataConnections/eventHubConnection1
```
