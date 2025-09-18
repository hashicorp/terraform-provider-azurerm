---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_connection"
description: |-
  Manages a Service Connector for a Container App.
---

# azurerm_container_app_connection

Manages a Service Connector for a Container App.

## Example Usage

### Container App Connection to CosmosDB with System Assigned Identity

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmosdb-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "exampledb"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  throughput          = 400
}

resource "azurerm_container_app_environment" "example" {
  name                = "Example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_container_app" "example" {
  name                         = "example-app"
  container_app_environment_id = azurerm_container_app_environment.example.id
  resource_group_name          = azurerm_resource_group.example.name
  revision_mode                = "Single"

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

resource "azurerm_container_app_connection" "example" {
  name               = "example-serviceconnector"
  container_app_id   = azurerm_container_app.example.id
  target_resource_id = azurerm_cosmosdb_sql_database.example.id
  scope              = "container"

  authentication {
    type = "systemAssignedIdentity"
  }
}
```

### Container App Connection to Storage Blob with Secret Authentication

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_container_app_environment" "example" {
  name                = "Example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_container_app" "example" {
  name                         = "example-app"
  container_app_environment_id = azurerm_container_app_environment.example.id
  resource_group_name          = azurerm_resource_group.example.name
  revision_mode                = "Single"

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

resource "azurerm_container_app_connection" "example" {
  name               = "example-serviceconnector"
  container_app_id   = azurerm_container_app.example.id
  target_resource_id = azurerm_storage_account.example.id
  scope              = "container"

  authentication {
    type   = "secret"
    name   = "accesskey"
    secret = azurerm_storage_account.example.primary_access_key
  }
}
```

### Container App Connection with User Assigned Identity

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "example" {}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "example-cosmosdb-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "BoundedStaleness"
    max_interval_in_seconds = 10
    max_staleness_prefix    = 200
  }

  geo_location {
    location          = azurerm_resource_group.example.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "example" {
  name                = "exampledb"
  resource_group_name = azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  throughput          = 400
}

resource "azurerm_container_app_environment" "example" {
  name                = "Example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_container_app" "example" {
  name                         = "example-app"
  container_app_environment_id = azurerm_container_app_environment.example.id
  resource_group_name          = azurerm_resource_group.example.name
  revision_mode                = "Single"

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

resource "azurerm_container_app_connection" "example" {
  name               = "example-serviceconnector"
  container_app_id   = azurerm_container_app.example.id
  target_resource_id = azurerm_cosmosdb_sql_database.example.id
  scope              = azurerm_container_app.example.template[0].container[0].name

  authentication {
    type            = "userAssignedIdentity"
    subscription_id = data.azurerm_subscription.example.subscription_id
    client_id       = azurerm_user_assigned_identity.example.client_id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the service connection. Changing this forces a new resource to be created.

* `container_app_id` - (Required) The ID of the Container App. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the target resource. Changing this forces a new resource to be created.

* `scope` - (Required) The scope of the connection. Changing this forces a new resource to be created.

* `authentication` - (Required) The authentication info. An `authentication` block as defined below.

* `client_type` - (Optional) The application client type. Possible values are `none`, `dotnet`, `java`, `python`, `go`, `php`, `ruby`, `django`, `nodejs` and `springBoot`. Defaults to `none`.

* `secret_store` - (Optional) An option to store secret value in secure place. An `secret_store` block as defined below.

* `vnet_solution` - (Optional) The type of the VNet solution. Possible values are `serviceEndpoint` and `privateLink`.

---

An `authentication` block supports the following:

* `type` - (Required) The authentication type. Possible values are `systemAssignedIdentity`, `userAssignedIdentity`, `servicePrincipalSecret`, `servicePrincipalCertificate` and `secret`.

* `client_id` - (Optional) Client ID for userAssignedIdentity or servicePrincipal auth. Required when `type` is set to `userAssignedIdentity` or `servicePrincipalSecret` or `servicePrincipalCertificate`.

* `subscription_id` - (Optional) Subscription ID for userAssignedIdentity. Required when `type` is set to `userAssignedIdentity`.

* `principal_id` - (Optional) Principal ID for servicePrincipal auth. Required when `type` is set to `servicePrincipalSecret` or `servicePrincipalCertificate`.

* `secret` - (Optional) Password or account key for secret auth. Required when `type` is set to `servicePrincipalSecret` or `secret`.

* `name` - (Optional) Username or account name for secret auth. Required when `type` is set to `secret`.

* `certificate` - (Optional) Service principal certificate for servicePrincipal auth. Required when `type` is set to `servicePrincipalCertificate`.

---

A `secret_store` block supports the following:

* `key_vault_id` - (Required) The key vault id to store secret.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Service Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Service Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Service Connector.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Service Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Service Connector.

## Import

Container App Service Connector can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.App/containerApps/containerApp1/providers/Microsoft.ServiceLinker/linkers/link1
```
