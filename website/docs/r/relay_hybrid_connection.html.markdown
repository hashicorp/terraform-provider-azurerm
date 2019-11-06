---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_relay_hybrid_connection"
sidebar_current: "docs-azurerm-resource-messaging-relay-hybrid-connection"
description: |-
  Manages an Azure Relay Hybrid Connection.

---

# azurerm_relay_hybrid_connection

Manages an Azure Relay Hybrid Connection.

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

  sku_name = "Standard"

  tags = {
    source = "terraform"
  }
}

resource "azurerm_relay_hybrid_connection" "test" {
  name                 = "acctestrnhc-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  relay_namespace_name   = "${azurerm_relay_namespace.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure Relay Hybrid Connection. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Azure Relay Hybrid Connection. Changing this forces a new resource to be created.

* `relay_namespace_name` - (Required) The name of the Azure Relay in which to create the Azure Relay Hybrid Connection. Changing this forces a new resource to be created.

* `requires_client_authorization` - (Optional) Specify if client authorization is needed for this hybrid connection. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Azure Relay Hybrid Connection ID.

## Import

Azure Relay Namespace's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_relay_namespace.relay1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Relay/namespaces/relay1/hybridConnections/hconn1
```
