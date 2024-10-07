---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_sas"
description: |-
  Gets a Shared Access Signature (SAS Token) for an existing Event Hub.
---

# Data Source: azurerm_eventhub_sas

Use this data source to obtain a Shared Access Signature (SAS Token) for an existing Event Hub.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-ehn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "example" {
  name                = "example-eh"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 1
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "example" {
  name                = "example-ehar"
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name

  listen = true
  send   = true
  manage = true
}

data "azurerm_eventhub_authorization_rule" "example" {
  name                = azurerm_eventhub_authorization_rule.example.name
  namespace_name      = azurerm_eventhub_namespace.example.name
  eventhub_name       = azurerm_eventhub.example.name
  resource_group_name = azurerm_resource_group.example.name
}

data "azurerm_eventhub_sas" "example" {
  connection_string = data.azurerm_eventhub_authorization_rule.example.primary_connection_string
  expiry            = "2023-06-23T00:00:00Z"
}
```

## Argument Reference

* `connection_string` - The connection string for the Event Hub to which this SAS applies.

* `expiry` - The expiration time and date of this SAS. Must be a valid ISO-8601 format time/date string.

## Attributes Reference

* `sas` - The computed Event Hub Shared Access Signature (SAS).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SAS Token.
