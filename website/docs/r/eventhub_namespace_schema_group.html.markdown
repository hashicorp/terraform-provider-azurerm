---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace_schema_group"
description: |-
  Manages a Schema Group for a EventHub Namespace.
---

# azurerm_eventhub_namespace_schema_group

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG-ehn-schemaGroup"
  location = "East US"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "example-ehn-schemaGroup"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_eventhub_namespace_schema_group" "test" {
  name                 = "example-schemaGroup"
  namespace_id         = azurerm_eventhub_namespace.test.id
  schema_compatibility = "Forward"
  schema_type          = "Avro"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this schema group. Changing this forces a new resource to be created.

* `namespace_id` - (Required) Specifies the ID of the EventHub Namespace. Changing this forces a new resource to be created.

* `schema_compatibility` - (Required) Specifies the compatibility of this schema group. Possible values are `None`, `Backward`, `Forward`. Changing this forces a new resource to be created.

* `schema_type` - (Required) Specifies the Type of this schema group. Possible values are `Avro`, `Unknown`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the EventHub Namespace Schema Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Namespace Schema Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Namespace Schema Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Namespace Schema Group.

## Import

Schema Group for a EventHub Namespace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_namespace_schema_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/schemaGroups/group1
```
