---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_webhook"
sidebar_current: "docs-azurerm-resource-container-registry-webhook"
description: |-
  Manages an Azure Container Registry Webhook.

---

# azurerm_container_registry_webhook

Manages an Azure Container Registry Webhook.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_container_registry" "acr" {
  name                     = "containerRegistry1"
  resource_group_name      = "${azurerm_resource_group.rg.name}"
  location                 = "${azurerm_resource_group.rg.location}"
  sku                      = "Standard"
  admin_enabled            = false
}

resource "azurerm_container_registry_webhook" "webhook" {
  name                = "mywebhook"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  registry_name       = "${azurerm_container_registry.acr.name}"
  location            = "${azurerm_resource_group.rg.location}"
  
  service_uri    = "https://mywebhookreceiver.example/mytag"
  status         = "enabled"
  scope          = "mytag:*"
  actions        = ["push"]
  custom_headers = { "Content-Type" = "application/json" }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container Registry Webhook. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Registry Webhook. Changing this forces a new resource to be created.

* `registry_name` - (Required) The Name of Container registry this Webhook belongs to. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `service_uri` - (Required) Specifies the service URI for the Webhook to post notifications.

* `actions` - (Required) A list of actions that trigger the Webhook to post notifications. At least one action needs to be specified. Valid values are: `push`, `delete`, `quarantine`, `chart_push`, `chart_delete`

* `status` - (Optional) Specifies if this Webhook triggers notifications or not. Valid values: `enabled` and `disabled`. Default is `enabled`. 

* `scope` - (Optional) Specifies the scope of repositories that can trigger an event. For example, 'foo:*' means events for all tags under repository 'foo'. 'foo:bar' means events for 'foo:bar' only. 'foo' is equivalent to 'foo:latest'. Empty means all events.

* `custom_headers` - (Optional) Custom headers that will be added to the webhook notifications request.

---
## Attributes Reference

The following attributes are exported:

* `id` - The Container Registry Webhook ID.

## Import

Container Registry Webhooks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_webhook.test /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1/webhooks/mywebhook1
```
