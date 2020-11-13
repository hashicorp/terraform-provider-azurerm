---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_cluster_customer_managed_key"
description: |-
  Manages a Log Analytics Cluster Customer Managed Key.
---

# azurerm_log_analytics_cluster_customer_managed_key

Manages a Log Analytics Cluster Customer Managed Key.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault" "example" {
  name                = "keyvaultkeyexample"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  sku_name = "premium"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "get",
    ]

    secret_permissions = [
      "set",
    ]
  }

  tags = {
    environment = "Production"
  }

  access_policy {
    tenant_id = azurerm_log_analytics_cluster.example.identity.0.tenant_id
    object_id = azurerm_log_analytics_cluster.example.identity.0.principal_id

    key_permissions = [
      "get",
      "unwrapkey",
      "wrapkey",
    ]
  }

}

resource "azurerm_key_vault_key" "example" {
  name         = "generated-certificate"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_log_analytics_cluster_customer_managed_key" "example" {
  log_analytics_cluster_id = azurerm_log_analytics_cluster.example.id
  key_vault_key_id         = azurerm_key_vault_key.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `key_vault_key_id` - (Required) The ID of the Key Vault Key to use for encryption.

* `log_analytics_cluster_id` - (Required) The ID of the Log Analytics Cluster. Changing this forces a new Log Analytics Cluster Custoemr Managed Key to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Log Analytics Cluster Custoemr Managed Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the Log Analytics Cluster Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Cluster Customer Managed Key.
* `update` - (Defaults to 6 hours) Used when updating the Log Analytics Cluster Customer Managed Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Cluster Customer Managed Key.

## Import

Log Analytics Cluster Customer Managed Keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_cluster_customer_managed_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/group1/providers/Microsoft.OperationalInsights/clusters/cluster1/CMK
```
