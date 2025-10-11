---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_slot_connection"
description: |-
  Manages an App Service Slot Connection for a service.
---

# azurerm_app_service_slot_connection

Manages an App Service Slot Connection for a service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacct"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-appserviceplan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "example-app-service"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_app_service_slot" "example" {
  name                = "example-slot"
  app_service_name    = azurerm_app_service.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_app_service_slot_connection" "example" {
  name                = "example-serviceconnector"
  app_service_slot_id = azurerm_app_service_slot.example.id
  target_resource_id  = azurerm_storage_account.example.id
  client_type         = "dotnet"

  authentication {
    type = "systemAssignedIdentity"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the service connection. Changing this forces a new resource to be created.

* `app_service_slot_id` - (Required) The ID of the App Service Slot. Changing this forces a new resource to be created.

* `authentication` - (Required) The authentication info. An `authentication` block as defined below. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The ID of the target resource. Changing this forces a new resource to be created.

* `client_type` - (Optional) The application client type. Possible values are `none`, `dotnet`, `java`, `python`, `go`, `php`, `ruby`, `django`, `nodejs` and `springBoot`. Defaults to `none`. Changing this forces a new resource to be created.

* `secret_store` - (Optional) An option to store secret value in secure place. A `secret_store` block as defined below. Changing this forces a new resource to be created.

* `vnet_solution` - (Optional) The type of the VNet solution. Possible values are `serviceEndpoint` and `privateLink`. Changing this forces a new resource to be created.

---

An `authentication` block supports the following:

* `type` - (Required) The authentication type. Possible values are `systemAssignedIdentity`, `userAssignedIdentity`, `servicePrincipalSecret`, `servicePrincipalCertificate`, `secret`. Changing this forces a new resource to be created.

* `certificate` - (Optional) The service principal certificate. It's used when `type` is set to `servicePrincipalCertificate`. Changing this forces a new resource to be created.

* `client_id` - (Optional) The application client ID. It's used when `type` is set to `servicePrincipalSecret` or `servicePrincipalCertificate`. Changing this forces a new resource to be created.

* `name` - (Optional) The username or account name for authentication. It's used when `type` is set to `secret`. Changing this forces a new resource to be created.

* `principal_id` - (Optional) The principal ID for user assigned identity. It's used when `type` is set to `userAssignedIdentity`. Changing this forces a new resource to be created.

* `secret` - (Optional) The password or key for authentication. It's used when `type` is set to `servicePrincipalSecret` or `secret`. Changing this forces a new resource to be created.

* `subscription_id` - (Optional) The subscription ID. It's used when `type` is set to `servicePrincipalSecret` or `servicePrincipalCertificate`. Changing this forces a new resource to be created.

---

A `secret_store` block supports the following:

* `key_vault_id` - (Required) The key vault ID to store secret. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Slot Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Slot Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Slot Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Slot Connection.

## Import

App Service Slot Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_slot_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/sites/site1/slots/slot1/providers/Microsoft.ServiceLinker/linkers/linker1
```
