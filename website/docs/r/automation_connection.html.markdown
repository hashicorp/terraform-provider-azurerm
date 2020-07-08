---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_connection"
description: |-
  Manages an Automation Connection.
---

# azurerm_automation_connection

Manages an Automation Connection.

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

  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_connection" "example" {
  name                    = "connection-example"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  type                    = "AzureServicePrincipal"

  values = {
    "ApplicationId" : "00000000-0000-0000-0000-000000000000",
    "TenantId" : data.azurerm_client_config.example.tenant_id,
    "SubscriptionId" : data.azurerm_client_config.example.subscription_id,
    "CertificateThumbprint" : "sample-certificate-thumbprint",
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Connection. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Connection is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Connection is created. Changing this forces a new resource to be created.

* `type` - (Required) The type of the Connection - can be either builtin type such as `Azure`, `AzureClassicCertificate`, and `AzureServicePrincipal`, or a user defined types. Changing this forces a new resource to be created.

* `values` - (Optional) A mapping of key value pairs passed to the connection. Different `type` needs different parameters in the `values`. Builtin types have required field values as below:

  - `Azure`: parameters `AutomationCertificateName` and `SubscriptionID`.

  - `AzureClassicCertificate`: parameters `SubscriptionName`, `SubscriptionId` and `CertificateAssetName`.

  - `AzureServicePrincipal`: parameters `ApplicationId`, `CertificateThumbprint`, `SubscriptionId` and `TenantId`.

* `description` - (Optional) A description for this Connection.

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Connection ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Connection.

## Import

Automation Connection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/connections/conn1
```
