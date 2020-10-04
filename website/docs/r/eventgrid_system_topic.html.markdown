c---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_system_topic"
description: |-
  Manages an EventGrid System Topic
---

# azurerm_eventgrid_system_topic

Manages an EventGrid System Topic.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "ExampleResourceGroup"
  location = "westus2"
}

resource "azurerm_storage_account" "example" {
  name                     = "ExampleStorageAccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_eventgrid_system_topic" "example" {
  name                = "ExampleSystemTopic"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  source_resource_id  = azurerm_storage_account.example.id
  topic_type          = "Microsoft.Storage.StorageAccounts"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the EventGrid System Topic should exist. Changing this forces a new EventGrid System Topic to be created.

* `name` - (Required) The name which should be used for this EventGrid System Topic. Changing this forces a new EventGrid System Topic to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the EventGrid System Topic should exist. Changing this forces a new EventGrid System Topic to be created.

* `source_resource_id` - (Required) The ID of the EventGrid System Topic Source. Changing this forces a new EventGrid System Topic to be created.

* `topic_type` - (Required) The Topic Type of the EventGrid System Topic. Changing this forces a new EventGrid System Topic to be created. Possible values are:
  * `Microsoft.Eventhub.Namespaces`
  * `Microsoft.Storage.StorageAccounts`
  * `Microsoft.Resources.Subscriptions`
  * `Microsoft.Resources.ResourceGroups`
  * `Microsoft.Devices.IoTHubs`
  * `Microsoft.EventGrid.Topics`
  * `Microsoft.ServiceBus.Namespaces`
  * `Microsoft.ContainerRegistry.Registries`
  * `Microsoft.Media.MediaServices`
  * `Microsoft.Maps.Accounts`
  * `Microsoft.EventGrid.Domains`
  * `Microsoft.AppConfiguration.ConfigurationStores`
  * `Microsoft.KeyVault.vaults`
  * `Microsoft.Web.Sites`
  * `Microsoft.Web.ServerFarms`
  * `Microsoft.SignalRService.SignalR`
  * `Microsoft.MachineLearningServices.Workspaces`
  * `Microsoft.Communication.CommunicationServices`
---

* `tags` - (Optional) A mapping of tags which should be assigned to the EventGrid System Topic.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the EventGrid System Topic.

* `metric_resource_id` - The ID of the Resource used by Azure Monitor.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid System Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid System Topic.
* `update` - (Defaults to 30 minutes) Used when updating the EventGrid System Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid System Topic.

## Import

EventGrid System Topic can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_system_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/systemTopics/systemTopic1
```
