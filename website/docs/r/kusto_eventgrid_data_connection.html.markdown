---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_eventgrid_data_connection"
description: |-
  Manages Kusto / Data Explorer Event Grid Data Connection
---

# azurerm_kusto_eventgrid_data_connection

Manages a Kusto (also known as Azure Data Explorer) Event Grid Data Connection

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kusto_cluster" "example" {
  name                = "examplekustocluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }
}

resource "azurerm_kusto_database" "example" {
  name                = "example-kusto-database"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name
  hot_cache_period    = "P7D"
  soft_delete_period  = "P31D"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "eventhubnamespace-example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "example" {
  name                = "eventhub-example"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_consumer_group" "example" {
  name                = "consumergroup-example"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_eventgrid_event_subscription" "example" {
  name                  = "eventgrid-example"
  scope                 = azurerm_storage_account.example.id
  eventhub_endpoint_id  = azurerm_eventhub.example.id
  event_delivery_schema = "EventGridSchema"
  included_event_types  = ["Microsoft.Storage.BlobCreated", "Microsoft.Storage.BlobRenamed"]

  retry_policy {
    event_time_to_live    = 144
    max_delivery_attempts = 10
  }
}

resource "azurerm_kusto_eventgrid_data_connection" "example" {
  name                         = "my-kusto-eventgrid-data-connection"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  cluster_name                 = azurerm_kusto_cluster.example.name
  database_name                = azurerm_kusto_database.example.name
  storage_account_id           = azurerm_storage_account.example.id
  eventhub_id                  = azurerm_eventhub.example.id
  eventhub_consumer_group_name = azurerm_eventhub_consumer_group.example.name

  table_name        = "my-table"         #(Optional)
  mapping_rule_name = "my-table-mapping" #(Optional)
  data_format       = "JSON"             #(Optional)

  depends_on = [azurerm_eventgrid_event_subscription.example]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto Event Grid Data Connection to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Database should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Database should exist. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the Kusto Cluster this data connection will be added to. Changing this forces a new resource to be created.

* `database_name` - (Required) Specifies the name of the Kusto Database this data connection will be added to. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) Specifies the resource id of the Storage Account this data connection will use for ingestion. Changing this forces a new resource to be created.

* `eventhub_id` - (Required) Specifies the resource id of the Event Hub this data connection will use for ingestion. Changing this forces a new resource to be created.

* `eventhub_consumer_group_name` - (Required) Specifies the Event Hub consumer group this data connection will use for ingestion. Changing this forces a new resource to be created.

* `blob_storage_event_type` - (Optional) Specifies the blob storage event type that needs to be processed. Possible Values are `Microsoft.Storage.BlobCreated` and `Microsoft.Storage.BlobRenamed`. Defaults to `Microsoft.Storage.BlobCreated`.

* `data_format` - (Optional) Specifies the data format of the EventHub messages. Allowed values: `APACHEAVRO`, `AVRO`, `CSV`, `JSON`, `MULTIJSON`, `ORC`, `PARQUET`, `PSV`, `RAW`, `SCSV`, `SINGLEJSON`, `SOHSV`, `TSV`, `TSVE`, `TXT` and `W3CLOGFILE`.

* `database_routing_type` - (Optional) Indication for database routing information from the data connection, by default only database routing information is allowed. Allowed values: `Single`, `Multi`. Changing this forces a new resource to be created. Defaults to `Single`.

* `eventgrid_event_subscription_id` - (Optional) The resource ID of the event grid that is subscribed to the storage account events.

* `managed_identity_id` - (Optional) Empty for non-managed identity based data connection. For system assigned identity, provide cluster resource Id. For user assigned identity (UAI) provide the UAI resource Id.

* `mapping_rule_name` - (Optional) Specifies the mapping rule used for the message ingestion. Mapping rule must exist before resource is created.

* `table_name` - (Optional) Specifies the target table name used for the message ingestion. Table must exist before resource is created.

* `skip_first_record` - (Optional) is the first record of every file ignored? Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kusto Event Grid Data Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Kusto Event Grid Data Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Event Grid Data Connection.
* `update` - (Defaults to 1 hour) Used when updating the Kusto Event Grid Data Connection.
* `delete` - (Defaults to 1 hour) Used when deleting the Kusto Event Grid Data Connection.

## Import

Kusto Event Grid Data Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_eventgrid_data_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/clusters/cluster1/databases/database1/dataConnections/dataConnection1
```
