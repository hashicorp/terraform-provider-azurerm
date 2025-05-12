---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_webhook"
description: |-
  Manages an Azure Container Registry Webhook.

---

# azurerm_container_registry_webhook

Manages an Azure Container Registry Webhook.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_registry" "acr" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Standard"
  admin_enabled       = false
}

resource "azurerm_container_registry_webhook" "webhook" {
  name                = "mywebhook"
  resource_group_name = azurerm_resource_group.example.name
  registry_name       = azurerm_container_registry.acr.name
  location            = azurerm_resource_group.example.location

  service_uri = "https://mywebhookreceiver.example/mytag"
  status      = "enabled"
  scope       = "mytag:*"
  actions     = ["push"]
  custom_headers = {
    "Content-Type" = "application/json"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container Registry Webhook. Only Alphanumeric characters allowed. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Registry Webhook. Changing this forces a new resource to be created.

* `registry_name` - (Required) The Name of Container registry this Webhook belongs to. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `service_uri` - (Required) Specifies the service URI for the Webhook to post notifications.

* `actions` - (Required) A list of actions that trigger the Webhook to post notifications. At least one action needs to be specified. Valid values are: `push`, `delete`, `quarantine`, `chart_push`, `chart_delete`

* `status` - (Optional) Specifies if this Webhook triggers notifications or not. Valid values: `enabled` and `disabled`. Default is `enabled`.

* `scope` - (Optional) Specifies the scope of repositories that can trigger an event. For example, `foo:*` means events for all tags under repository `foo`. `foo:bar` means events for 'foo:bar' only. `foo` is equivalent to `foo:latest`. Empty means all events. Defaults to `""`.

* `custom_headers` - (Optional) Custom headers that will be added to the webhook notifications request.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry Webhook.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry Webhook.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Webhook.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry Webhook.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Webhook.

## Import

Container Registry Webhooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_webhook.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1/webHooks/mywebhook1
```
