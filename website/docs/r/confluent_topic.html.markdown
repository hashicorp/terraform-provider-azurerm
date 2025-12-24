---
subcategory: "Confluent"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_confluent_topic"
description: |-
  Manages a Confluent Topic.
---

# azurerm_confluent_topic

Manages a Confluent Kafka Topic within a Confluent Cluster on Azure.

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

resource "azurerm_confluent_topic" "example" {
  topic_name          = "orders"
  cluster_id          = azurerm_confluent_cluster.example.cluster_id
  environment_id      = azurerm_confluent_environment.example.environment_id
  organization_id     = azurerm_confluent_organization.example.name
  resource_group_name = azurerm_resource_group.example.name
  partitions_count    = "6"
  replication_factor  = "3"

  configs {
    name  = "cleanup.policy"
    value = "compact"
  }

  configs {
    name  = "retention.ms"
    value = "604800000"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `topic_name` - (Required) The name of the Kafka topic. Changing this forces a new resource to be created.

* `cluster_id` - (Required) The ID of the parent Confluent Cluster. Changing this forces a new resource to be created.

* `environment_id` - (Required) The ID of the parent Confluent Environment. Changing this forces a new resource to be created.

* `organization_id` - (Required) The name of the parent Confluent Organization. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Confluent Topic should exist. Changing this forces a new resource to be created.

* `partitions_count` - (Optional) The number of partitions for the topic. Changing this forces a new resource to be created.

* `replication_factor` - (Optional) The replication factor for the topic. Changing this forces a new resource to be created.

* `configs` - (Optional) One or more `configs` blocks as defined below. Topic configuration settings. Changing this forces a new resource to be created.

---

A `configs` block supports the following:

* `name` - (Required) The configuration property name (e.g., `cleanup.policy`, `retention.ms`).

* `value` - (Required) The configuration property value.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Confluent Topic.

* `topic_id` - The Confluent Topic ID.

* `kind` - The kind of the resource.

* `metadata` - A `metadata` block as defined below.

---

A `metadata` block exports the following:

* `self` - The self-referencing link for this topic.

* `resource_name` - The Confluent resource name for this topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Confluent Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the Confluent Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the Confluent Topic.

## Import

Confluent Topics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_confluent_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Confluent/organizations/org1/environments/env-12345/clusters/lkc-12345/topics/orders
```
