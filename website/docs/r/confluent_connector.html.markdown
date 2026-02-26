---
subcategory: "Confluent"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_confluent_connector"
description: |-
  Manages a Confluent Connector.
---

# azurerm_confluent_connector

Manages a Confluent Kafka Connector within a Confluent Cluster on Azure.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_confluent_organization" "example" {
  name                = "example-confluent-org"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  offer_detail {
    id           = "confluent-cloud-azure-prod"
    plan_id      = "confluent-cloud-azure-payg-prod"
    plan_name    = "Confluent Cloud - Pay as you Go"
    publisher_id = "confluentinc"
    term_unit    = "P1M"
  }

  user_detail {
    email_address = "user@example.com"
  }
}

resource "azurerm_confluent_environment" "example" {
  environment_id      = "env-12345"
  organization_id     = azurerm_confluent_organization.example.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_confluent_cluster" "example" {
  cluster_id          = "lkc-12345"
  environment_id      = azurerm_confluent_environment.example.environment_id
  organization_id     = azurerm_confluent_organization.example.name
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_confluent_connector" "example" {
  connector_name      = "azure-blob-sink-connector"
  cluster_id          = azurerm_confluent_cluster.example.cluster_id
  environment_id      = azurerm_confluent_environment.example.environment_id
  organization_id     = azurerm_confluent_organization.example.name
  resource_group_name = azurerm_resource_group.example.name
  connector_type      = "SINK"
  connector_class     = "AZUREBLOBSINK"
}
```

## Arguments Reference

The following arguments are supported:

* `connector_name` - (Required) The name of the Kafka connector. Changing this forces a new resource to be created.

* `cluster_id` - (Required) The ID of the parent Confluent Cluster. Changing this forces a new resource to be created.

* `environment_id` - (Required) The ID of the parent Confluent Environment. Changing this forces a new resource to be created.

* `organization_id` - (Required) The name of the parent Confluent Organization. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Confluent Connector should exist. Changing this forces a new resource to be created.

* `connector_type` - (Optional) The type of connector. Possible values are `SINK` and `SOURCE`. Changing this forces a new resource to be created.

* `connector_class` - (Optional) The connector class. Possible values are `AZUREBLOBSINK` and `AZUREBLOBSOURCE`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Confluent Connector.

* `connector_id` - The Confluent Connector ID.

* `connector_state` - The current state of the connector. Possible values include `RUNNING`, `FAILED`, `PAUSED`, and `PROVISIONING`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Confluent Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Confluent Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Confluent Connector.

## Import

Confluent Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_confluent_connector.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Confluent/organizations/org1/environments/env-12345/clusters/lkc-12345/connectors/azure-blob-sink-connector
```

~> **Note:** The current implementation provides basic connector management. For detailed connector configuration (such as Azure Blob Storage connection strings, Cosmos DB settings, etc.), additional configuration may need to be performed through the Confluent Cloud portal or CLI.
