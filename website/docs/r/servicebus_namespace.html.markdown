---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace"
sidebar_current: "docs-azurerm-resource-messaging-servicebus-namespace-x"
description: |-
  Manages a ServiceBus Namespace.
---

# azurerm_servicebus_namespace

Manages a ServiceBus Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "terraform-servicebus"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "tfex_sevicebus_namespace"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "Standard"

  tags = {
    source = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the ServiceBus Namespace resource . Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the namespace.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Defines which tier to use. Options are basic, standard or premium.

* `capacity` - (Optional) Specifies the capacity. When `sku` is `Premium` can be `1`, `2` or `4`. When `sku` is `Basic` or `Standard` can be `0` only.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ServiceBus Namespace ID.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Import

Service Bus Namespace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.servicebus/namespaces/sbns1
```
