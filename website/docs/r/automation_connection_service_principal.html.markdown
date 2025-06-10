---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_connection_service_principal"
description: |-
  Manages an Automation Connection with type `AzureServicePrincipal`.
---

# azurerm_automation_connection_service_principal

Manages an Automation Connection with type `AzureServicePrincipal`.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "resourceGroup-example"
  location = "West Europe"
}

data "azurerm_client_config" "example" {}

resource "azurerm_automation_account" "example" {
  name                = "account-example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_connection_service_principal" "example" {
  name                    = "connection-example"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  application_id          = "00000000-0000-0000-0000-000000000000"
  tenant_id               = data.azurerm_client_config.example.tenant_id
  subscription_id         = data.azurerm_client_config.example.subscription_id
  certificate_thumbprint  = file("automation_certificate_test.thumb")
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Connection. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Connection is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Connection is created. Changing this forces a new resource to be created.

* `application_id` - (Required) The (Client) ID of the Service Principal.

* `certificate_thumbprint` - (Required) The thumbprint of the Service Principal Certificate.

* `subscription_id` - (Required) The subscription GUID.
  
* `tenant_id` - (Required) The ID of the Tenant the Service Principal is assigned in.

* `description` - (Optional) A description for this Connection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Automation Connection ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Connection.

## Import

Automation Connection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_connection_service_principal.conn1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/connections/conn1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Automation`: 2023-11-01
