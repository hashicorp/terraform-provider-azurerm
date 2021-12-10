---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_iothub_data_connection"
description: |-
  Manages Kusto / Data Explorer IotHub Data Connection
---

# azurerm_kusto_iothub_data_connection

Manages a Kusto (also known as Azure Data Explorer) IotHub Data Connection

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

resource "azurerm_iothub" "example" {
  name                = "exampleIoTHub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "B1"
    capacity = "1"
  }
}

resource "azurerm_iothub_shared_access_policy" "example" {
  name                = "example-shared-access-policy"
  resource_group_name = azurerm_resource_group.example.name
  iothub_name         = azurerm_iothub.example.name

  registry_read = true
}

resource "azurerm_iothub_consumer_group" "example" {
  name                   = "example-consumer-group"
  resource_group_name    = azurerm_resource_group.example.name
  iothub_name            = azurerm_iothub.example.name
  eventhub_endpoint_name = "events"
}

resource "azurerm_kusto_iothub_data_connection" "example" {
  name                = "my-kusto-iothub-data-connection"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  cluster_name        = azurerm_kusto_cluster.example.name
  database_name       = azurerm_kusto_database.example.name

  iothub_id                 = azurerm_iothub.example.id
  consumer_group            = azurerm_iothub_consumer_group.example.name
  shared_access_policy_name = azurerm_iothub_shared_access_policy.example.name
  event_system_properties   = ["message-id", "sequence-number", "to"]

  table_name        = "my-table"
  mapping_rule_name = "my-table-mapping"
  data_format       = "JSON"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Kusto IotHub Data Connection to create. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Kusto Database should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Resource Group where the Kusto Database should exist. Changing this forces a new resource to be created.

* `cluster_name` - (Required) Specifies the name of the Kusto Cluster this data connection will be added to. Changing this forces a new resource to be created.

* `database_name` - (Required) Specifies the name of the Kusto Database this data connection will be added to. Changing this forces a new resource to be created.

* `iothub_id` - (Required) Specifies the resource id of the IotHub this data connection will use for ingestion. Changing this forces a new resource to be created.

* `consumer_group` - (Required) Specifies the IotHub consumer group this data connection will use for ingestion. Changing this forces a new resource to be created.

* `shared_access_policy_name` - (Required) Specifies the IotHub Shared Access Policy this data connection will use for ingestion, which must have read permission. Changing this forces a new resource to be created.

* `event_system_properties` - (Optional) Specifies the System Properties that each IoT Hub message should contain. Changing this forces a new resource to be created.

* `table_name` - (Optional) Specifies the target table name used for the message ingestion. Table must exist before resource is created.

* `mapping_rule_name` - (Optional) Specifies the mapping rule used for the message ingestion. Mapping rule must exist before resource is created.

* `data_format` - (Optional) Specifies the data format of the IoTHub messages. Allowed values: `APACHEAVRO`, `AVRO`, `CSV`, `JSON`, `MULTIJSON`, `ORC`, `PARQUET`, `PSV`, `RAW`, `SCSV`, `SINGLEJSON`, `SOHSV`, `TSV`, `TSVE`, `TXT` and `W3CLOGFILE`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Kusto IotHub Data Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Kusto IotHub Data Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto IotHub Data Connection.
* `delete` - (Defaults to 60 minutes) Used when deleting the Kusto IotHub Data Connection.

## Import

Kusto IotHub Data Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_iothub_data_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1/Databases/database1/DataConnections/dataConnection1
```
