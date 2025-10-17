---
subcategory: "Confluent"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_confluent_cluster"
description: |-
  Manages a Confluent Cluster.
---

# azurerm_confluent_cluster

Manages a Confluent Kafka Cluster within a Confluent Environment on Azure.

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
  display_name        = "production-cluster"
  availability        = "SINGLE_ZONE"
  cloud               = "AZURE"
  region              = "westeurope"
  package             = "ESSENTIALS"

  spec {
    zone = "westeurope-1"

    config {
      kind = "Basic"
    }

    environment {
      id = "env-12345"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `cluster_id` - (Required) The Confluent Cluster ID. Changing this forces a new resource to be created.

* `environment_id` - (Required) The ID of the parent Confluent Environment. Changing this forces a new resource to be created.

* `organization_id` - (Required) The name of the parent Confluent Organization. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Confluent Cluster should exist. Changing this forces a new resource to be created.

* `display_name` - (Optional) The display name for the cluster. Changing this forces a new resource to be created.

* `availability` - (Optional) The availability zone configuration (e.g., `SINGLE_ZONE`, `MULTI_ZONE`). Changing this forces a new resource to be created.

* `cloud` - (Optional) The cloud provider (e.g., `AZURE`). Changing this forces a new resource to be created.

* `region` - (Optional) The Azure region for the cluster. Changing this forces a new resource to be created.

* `package` - (Optional) The cluster package type. Possible values are `ESSENTIALS` and `ADVANCED`. Changing this forces a new resource to be created.

* `spec` - (Optional) A `spec` block as defined below. Changing this forces a new resource to be created.

---

A `spec` block supports the following:

* `zone` - (Optional) The availability zone.

* `config` - (Optional) A `config` block as defined below.

* `environment` - (Optional) An `environment` block as defined below.

* `network` - (Optional) A `network` block as defined below.

* `byok` - (Optional) A `byok` (Bring Your Own Key) block as defined below.

---

A `config` block supports the following:

* `kind` - (Optional) The configuration kind (e.g., `Basic`, `Standard`, `Dedicated`).

---

An `environment` block supports the following:

* `id` - (Optional) The environment ID.

* `environment` - (Optional) The environment reference.

---

A `network` block supports the following:

* `id` - (Optional) The network ID.

* `environment` - (Optional) The network environment reference.

---

A `byok` block supports the following:

* `id` - (Optional) The BYOK key ID.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Confluent Cluster.

* `kind` - The kind of the resource.

* `api_endpoint` - The API endpoint for the cluster.

* `http_endpoint` - The HTTP endpoint for the cluster.

* `kafka_bootstrap_endpoint` - The Kafka bootstrap endpoint for the cluster.

* `metadata` - A `metadata` block as defined below.

* `status` - A `status` block as defined below.

---

A `metadata` block exports the following:

* `self` - The self-referencing link for this cluster.

* `resource_name` - The Confluent resource name for this cluster.

* `created_timestamp` - The timestamp when the cluster was created.

* `updated_timestamp` - The timestamp when the cluster was last updated.

* `deleted_timestamp` - The timestamp when the cluster was deleted (if applicable).

---

A `status` block exports the following:

* `phase` - The current phase of the cluster (e.g., `PROVISIONING`, `RUNNING`).

* `cku` - The number of Confluent Kafka Units (CKU) allocated to the cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Confluent Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Confluent Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Confluent Cluster.

## Import

Confluent Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_confluent_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Confluent/organizations/org1/environments/env-12345/clusters/lkc-12345
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Confluent` - 2024-07-01
