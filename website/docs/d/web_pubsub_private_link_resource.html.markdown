---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_web_pubsub_private_link_resource"
description: |-
  Gets information about the Private Link Resource supported by the Web Pubsub Resource.
---

# Data Source: azurerm_web_pubsub_private_link_resource

Use this data source to access information about the Private Link Resource supported by the Web Pubsub Resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "terraform-webpubsub"
  location = "east us"
}

resource "azurerm_web_pubsub" "test" {
  name                = "tfex-webpubsub"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_S1"
  capacity            = 1
}

data "azurerm_web_pubsub_private_link_resource" "test" {
  web_pubsub_id = azurerm_web_pubsub.test.id
}
```

## Argument Reference

* `web_pubsub_id` - The ID of an existing Web Pubsub Resource which Private Link Resource should be retrieved for.

## Attributes Reference

* `id` - The ID of an existing Web Pubsub Resource which supports the retrieved Private Link Resource list.

* `shared_private_link_resource_types` - A `shared_private_link_resource_types` block as defined below.

---

A `shared_private_link_resource_types` block exports the following:

* `subresource_name` - The  name for the resource that has been onboarded to private link service.

* `description` - The description of the resource type that has been onboarded to private link service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private Link Resource.
