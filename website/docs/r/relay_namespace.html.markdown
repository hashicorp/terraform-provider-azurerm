---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_relay_namespace"
sidebar_current: "docs-azurerm-resource-relay-namespace"
description: |-
  Manages an Azure Relay Namespace.

---

# azurerm_relay_namespace

Manages an Azure Relay Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_relay_namespace" "test" {
  name                = "example-relay"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
    name = "Standard"
  }

  tags {
    source = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure Relay Namespace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure Relay Namespace.

* `location` - (Required) Specifies the supported Azure location where the Azure Relay Namespace exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block contains:

* `name` - (Required) The name of the SKU to use. At this time the only supported value is `Standard`.


## Attributes Reference

The following attributes are exported:

* `id` - The Azure Relay Namespace ID.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `primary_connection_string` - The primary connection string for the authorization rule `RootManageSharedAccessKey`.

* `secondary_connection_string` - The secondary connection string for the authorization rule `RootManageSharedAccessKey`.

* `primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

* `metric_id` - The Identifier for Azure Insights metrics.

## Import

Azure Relay Namespace's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_relay_namespace.relay1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Relay/namespaces/relay1
```
