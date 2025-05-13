---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventgrid_domain_topic"
description: |-
  Manages an EventGrid Domain Topic
---

# azurerm_eventgrid_domain_topic

Manages an EventGrid Domain Topic

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_eventgrid_domain" "example" {
  name                = "my-eventgrid-domain"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    environment = "Production"
  }
}
resource "azurerm_eventgrid_domain_topic" "example" {
  name                = "my-eventgrid-domain-topic"
  domain_name         = azurerm_eventgrid_domain.example.name
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventGrid Domain Topic resource. Changing this forces a new resource to be created.

* `domain_name` - (Required) Specifies the name of the EventGrid Domain. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the EventGrid Domain exists. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventGrid Domain Topic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventGrid Domain Topic.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventGrid Domain Topic.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventGrid Domain Topic.

## Import

EventGrid Domain Topics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventgrid_domain_topic.topic1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventGrid/domains/domain1/topics/topic1
```
