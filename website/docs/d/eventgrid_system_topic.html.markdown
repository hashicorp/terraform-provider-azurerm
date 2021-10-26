---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_system_topic"
description: |-
  Gets information about an existing EventGrid System Topic

---

# Data Source: azurerm_eventgrid_system_topic

Use this data source to access information about an existing EventGrid System Topic

## Example Usage

```hcl
data "azurerm_eventgrid_system_topic" "example" {
  name                = "eventgrid-system-topic"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the EventGrid System Topic resource.

* `resource_group_name` - The name of the resource group in which the EventGrid System Topic exists.

## Attributes Reference

The following attributes are exported:

* `id` - The EventGrid System Topic ID.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Event Grid System Topic.

* `metric_arm_resource_id` - The Metric ARM Resource ID of the Event Grid System Topic.

* `source_arm_resource_id` - The ID of the Event Grid System Topic ARM Source.

* `topic_type` - The Topic Type of the Event Grid System Topic. Possible values are: `Microsoft.AppConfiguration.ConfigurationStores`, `Microsoft.Communication.CommunicationServices`
, `Microsoft.ContainerRegistry.Registries`, `Microsoft.Devices.IoTHubs`, `Microsoft.EventGrid.Domains`, `Microsoft.EventGrid.Topics`, `Microsoft.Eventhub.Namespaces`, `Microsoft.KeyVault.vaults`, `Microsoft.MachineLearningServices.Workspaces`, `Microsoft.Maps.Accounts`, `Microsoft.Media.MediaServices`, `Microsoft.Resources.ResourceGroups`, `Microsoft.Resources.Subscriptions`, `Microsoft.ServiceBus.Namespaces`, `Microsoft.SignalRService.SignalR`, `Microsoft.Storage.StorageAccounts`, `Microsoft.Web.ServerFarms` and `Microsoft.Web.Sites`. Changing this forces a new Event Grid System Topic to be created.

~> **NOTE:** Some `topic_type`s (e.g. **Microsoft.Resources.Subscriptions**) have location set to `Global` instead of a real location like `West US`.

* `tags` - A mapping of tags which are assigned to the Event Grid System Topic.

---

An `identity` block exports the following:

* `type` - Specifies the type of Managed Service Identity that is configured on this Event Grid System Topic.

* `principal_id` - Specifies the Principal ID of the System Assigned Managed Service Identity that is configured on this Event Grid System Topic.

* `tenant_id` - Specifies the Tenant ID of the System Assigned Managed Service Identity that is configured on this Event Grid System Topic.

* `identity_ids` - A list of IDs for User Assigned Managed Identity resources to be assigned.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid System Topic.
