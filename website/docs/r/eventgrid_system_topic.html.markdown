---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_system_topic"
description: |-
  Manages an Event Grid System Topic
---

# azurerm_eventgrid_system_topic

Manages an Event Grid System Topic.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracct"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_eventgrid_system_topic" "example" {
  name                   = "example-topic"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  source_arm_resource_id = azurerm_storage_account.example.id
  topic_type             = "Microsoft.Storage.StorageAccounts"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Event Grid System Topic should exist. Changing this forces a new Event Grid System Topic to be created.

* `name` - (Required) The name which should be used for this Event Grid System Topic. Changing this forces a new Event Grid System Topic to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Event Grid System Topic should exist. Changing this forces a new Event Grid System Topic to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `source_arm_resource_id` - (Required) The ID of the Event Grid System Topic ARM Source. Changing this forces a new Event Grid System Topic to be created.

* `topic_type` - (Required) The Topic Type of the Event Grid System Topic. Possible values are: `Microsoft.AppConfiguration.ConfigurationStores`, `Microsoft.Communication.CommunicationServices`
, `Microsoft.ContainerRegistry.Registries`, `Microsoft.Devices.IoTHubs`, `Microsoft.EventGrid.Domains`, `Microsoft.EventGrid.Topics`, `Microsoft.Eventhub.Namespaces`, `Microsoft.KeyVault.vaults`, `Microsoft.MachineLearningServices.Workspaces`, `Microsoft.Maps.Accounts`, `Microsoft.Media.MediaServices`, `Microsoft.Resources.ResourceGroups`, `Microsoft.Resources.Subscriptions`, `Microsoft.ServiceBus.Namespaces`, `Microsoft.SignalRService.SignalR`, `Microsoft.Storage.StorageAccounts`, `Microsoft.Web.ServerFarms` and `Microsoft.Web.Sites`. Changing this forces a new Event Grid System Topic to be created.

~> **NOTE:** Some `topic_type`s (e.g. **Microsoft.Resources.Subscriptions**) requires location to be set to `Global` instead of a real location like `West US`.

---

A `identity` block supports the following:

* `type` - Specifies the identity type of Event Grid System Topic. Possible values are `SystemAssigned` (where Azure will generate a Principal for you) or `UserAssigned` where you can specify the User Assigned Managed Identity IDs in the `identity_ids` field.

~> **NOTE:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the Event Grid System Topic has been created. More details are available below.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Event Grid System Topic.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Event Grid System Topic.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Event Grid System Topic.

* `metric_arm_resource_id` - The Metric ARM Resource ID of the Event Grid System Topic.

---

A `identity` block supports the following:

* `type` - Specifies the type of Managed Service Identity that is configured on this Event Grid System Topic.

* `principal_id` - Specifies the Principal ID of the System Assigned Managed Service Identity that is configured on this Event Grid System Topic.

* `tenant_id` - Specifies the Tenant ID of the System Assigned Managed Service Identity that is configured on this Event Grid System Topic.

* `identity_ids` - A list of IDs for User Assigned Managed Identity resources to be assigned.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Event Grid System Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the Event Grid System Topic.
* `update` - (Defaults to 30 minutes) Used when updating the Event Grid System Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the Event Grid System Topic.

## Import

Event Grid System Topic can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_system_topic.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/systemTopics/systemTopic1
```
