---
subcategory: "App Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_configuration_replica"
description: |-
  Manages a App Configuration Replica.

---

# azurerm_app_configuration_replica

Manages an App Configuration Replica.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_configuration" "example" {
  name                       = "example-app-conf"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  sku                        = "standard"
  soft_delete_retention_days = 1
}
resource "azurerm_app_configuration_replica" "example" {
  configuration_store_id = azurerm_app_configuration.example.id
  location                 = azurerm_resource_group.example.location
  name                     = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `configuration_store_id` - (Required) Specifies the ID of the Configuration Store within which this App Configuration Replica should exist. Changing this forces a new App Configuration Replica to be created.

* `location` - (Required) The Azure Region where the App Configuration Replica should exist. Changing this forces a new App Configuration Replica to be created.

* `name` - (Required) Specifies the name of this App Configuration Replica. Changing this forces a new App Configuration Replica to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Configuration Replica.

* `endpoint` - The URI of the replica where the replica API will be available.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this App Configuration Replica.
* `delete` - (Defaults to 30 minutes) Used when deleting this App Configuration Replica.
* `read` - (Defaults to 5 minutes) Used when retrieving this App Configuration Replica.

## Import

An existing App Configuration Replica can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_app_configuration_replica.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.AppConfiguration/configurationStores/{configurationStoreName}/replicas/{replicaName}
```
