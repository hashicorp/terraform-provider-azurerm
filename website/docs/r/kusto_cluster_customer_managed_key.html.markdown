---
subcategory: "Data Explorer"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kusto_cluster_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a Kusto Cluster.
---

# azurerm_kusto_cluster_customer_managed_key

Manages a Customer Managed Key for a Kusto Cluster.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "cluster" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_kusto_cluster.example.identity.0.principal_id

  key_permissions = ["get", "unwrapkey", "wrapkey"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["get", "list", "create", "delete", "recover"]
}

resource "azurerm_key_vault_key" "example" {
  name         = "tfex-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.cluster,
  ]
}

resource "azurerm_kusto_cluster" "example" {
  name                = "kustocluster"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Standard_D13_v2"
    capacity = 2
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kusto_cluster_customer_managed_key" "example" {
  cluster_id   = azurerm_kusto_cluster.example.id
  key_vault_id = azurerm_key_vault.example.id
  key_name     = azurerm_key_vault_key.example.name
  key_version  = azurerm_key_vault_key.example.version
}
```

## Arguments Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the Kusto Cluster. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of the Key Vault. Changing this forces a new resource to be created.

* `key_name` - (Required) The name of Key Vault Key.

* `key_version` - (Required) The version of Key Vault Key.

* `user_identity` - (Optional) The user assigned identity that has access to the Key Vault Key. If not specified, system assigned identity will be used.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Kusto Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kusto Cluster Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kusto Cluster Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the Kusto Cluster Customer Managed Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kusto Cluster Customer Managed Key.

## Import

Customer Managed Keys for a Kusto Cluster can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kusto_cluster_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Kusto/Clusters/cluster1
```
