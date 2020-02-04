---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_eventhub_data_connection"
description: |-
  Manages Kusto / Data Explorer EventGrid Data Connection
---

# azurerm_kusto_eventgrid_data_connection

Manages a Kusto (also known as Azure Data Explorer) EventGrid Data Connection

## Example Usage

Assuming that Kusto Cluster, Database, Table & Mapping are created.
   * Table & Mapping are part of the data plan

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "my-kusto-rg"
  location = "East US"
}
resource "azurerm_eventhub_namespace" "eventhub_ns" {
  name                = "my-eventhub-ns"
  location            = "${azurerm_resource_group.rg.location}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  sku                 = "Standard"
}
resource "azurerm_eventhub" "eventhub" {
  name                = "my-eventhub"
  namespace_name      = "${azurerm_eventhub_namespace.eventhub_ns.name}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  partition_count     = 1
  message_retention   = 1
}
resource "azurerm_storage_account" "storageaccount" {
  name                     = "my-storage-account"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_eventhub_consumer_group" "consumer_group" {
  name                = "my-eventhub-consumergroup"
  namespace_name      = "${azurerm_eventhub_namespace.eventhub_ns.name}"
  eventhub_name       = "${azurerm_eventhub.eventhub.name}"
  resource_group_name = "${azurerm_resource_group.rg.name}"
}
resource "azurerm_kusto_eventhub_data_connection" "eventhub_connection" {
  name                = "my-kusto-eventhub-data-connection"
  resource_group_name = "my-cluster-rg"
  location            = "East US" # Cluster location
  cluster_name        = "my-cluster-name"
  database_name       = "my-database-name"

  storage_account_id = "${azurerm_storage_account.storageaccount.id}"
  eventhub_id        = "${azurerm_eventhub.evenhub.id}"
  consumer_group     = "${azurerm_eventhub_consumer_group.consumer_group.name}"

  table_name        = "my-table"
  mapping_rule_name = "my-table-mapping"
  data_format       = "JSON"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto EventGrid Data Connection to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Database should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Database should exist. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the Kusto Cluster this data connection will be added to. Changing this forces a new resource to be created.

* `database_name` - (Required) Specifies the name of the Kusto Database this data connection will be added to. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) Specifies the resource id of the Storage Account this data connection will use for ingestion. Changing this forces a new resource to be created.

* `eventhub_id` - (Required) Specifies the resource id of the EventHub this data connection will use for ingestion. Changing this forces a new resource to be created.
* `consumer_group` - (Required) Specifies the EventHub consumer group this data connection will use for ingestion. Changing this forces a new resource to be created.

* `table_name` - (Required) Specifies the target table name used for the message ingestion. Table must exist before resource is created.

* `data_format` - (Required) Specifies the data format of the EventHub messages. Allowed values: `AVRO`, `CSV`, `JSON`, `MULTIJSON`, `PSV`, `RAW`, `SCSV`, `SINGLEJSON`, `SOHSV`, `TSV` and `TXT`

* `mapping_rule_name` - (Required) Specifies the mapping rule used for the message ingestion. Mapping rule must exist before resource is created.

## Attributes Reference

The following attributes are exported:

- `id` - The EventGrid Data Connection ID.


### Timeouts

~> **Note:** Custom Timeouts is available [as an opt-in Beta in version 1.43 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Kusto EventGrid Data Connection.
* `update` - (Defaults to 60 minutes) Used when updating the Kusto EventGrid Data Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto EventGrid Data Connection.
* `delete` - (Defaults to 60 minutes) Used when deleting the Kusto EventGrid Data Connection.

## Import

Kusto EventGrid Data Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_eventgrid_data_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1/DataConnections/eventGridConnection1
```
