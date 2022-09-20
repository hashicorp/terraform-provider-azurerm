---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_connection"
description: |-
  Manages a service connector for app service.
---

# azurerm_app_service_connection

Manages a service connector for app service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                       = "example-key-vault"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
}

resource "azurerm_service_plan" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-service-plan"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "P1v2"
  os_type             = "Linux"
}

resource "azurerm_linux_web_app" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-web-app"
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id
  site_config {}
}

resource "azurerm_app_service_connection" "example" {
  name               = "example-app-service-connection"
  app_service_id     = azurerm_linux_web_app.example.id
  target_resource_id = azurerm_key_vault.example.id
  authentication {
    type = "systemAssignedIdentity"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service connection. Changing this forces a new resource to be created.

* `app_service_id` - (Required) The ID of the data source web app. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the target resource. Changing this forces a new resource to be created. Possible values are `Postgres`, `PostgresFlexible`, `Mysql`, `Sql`, `Redis`, `RedisEnterprise`, `CosmosCassandra`, `CosmosGremlin`, `CosmosMongo`, `CosmosSql`, `CosmosTable`, `StorageBlob`, `StorageQueue`, `StorageFile`, `StorageTable`, `AppConfig`, `EventHub`, `ServiceBus`, `SignalR`, `WebPubSub`, `ConfluentKafka`.

* `authentication` - (Required) The authentication info. An `authentication` block as defined below.
---
* `type` - (Required) The authentication type. Possible values are `systemAssignedIdentity`, `userAssignedIdentity`, `servicePrincipalSecret`, `servicePrincipalCertificate`, `secret`.

* `name` - (Optional) Username or account name for secret auth. `name` and `secret` should be either both specified or both not specified when `type` is set to `secret`.

* `secret` - (Optional) Password or account key for secret auth. `secret` and `name` should be either both specified or both not specified when `type` is set to `secret`.

* `client_id` - (Optional) Client ID for `userAssignedIdentity` or `servicePrincipal` auth. Should be specified when `type` is set to `servicePrincipalSecret` or `servicePrincipalCertificate`. When `type` is set to `userAssignedIdentity`, `client_id` and `subscription_id` should be either both specified or both not specified.

* `subscription_id` - (Optional) Subscription ID for `userAssignedIdentity`. `subscription_id` and `client_id` should be either both specified or both not specified.

* `principal_id` - (Optional) Principal ID for `servicePrincipal` auth. Should be specified when `type` is set to `servicePrincipalSecret` or `servicePrincipalCertificate`.

* `certificate` - (Optional) Service principal certificate for `servicePrincipal` auth. Should be specified when `type` is set to `servicePrincipalCertificate`.
---

* `client_type` - (Optional) The application client type. Possible values are `dotnet`, `java`, `python`, `go`, `php`, `ruby`, `django`, `nodejs`, `springBoot`.

* `vnet_solution` - (Optional) The type of the VNet solution. Possible values are `serviceEndpoint`, `privateLink`.

## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the service connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Connector for app service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Connector for app service.
* `update` - (Defaults to 30 minutes) Used when updating the Service Connector for app service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Connector for app service.

## Import

Service Connector for app service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/webapp/providers/Microsoft.ServiceLinker/linkers/serviceconnector1
```
