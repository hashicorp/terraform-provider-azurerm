---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis_cluster"
description: |-
  Manages a Managed Redis Cluster.
---

# azurerm_managed_redis_cluster

Manages a Managed Redis Cluster. Accompanying [azurerm_managed_redis_database](managed_redis_database.html) has to be created for a fully functional Managed Redis feature.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_redis_cluster" "example" {
  name                = "example-managed-redis"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Balanced_B3"
}
```

## Example Usage with Customer Managed Key

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_key_vault" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  purge_protection_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "List",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.example.principal_id

    key_permissions = [
      "Get",
      "WrapKey",
      "UnwrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "managedrediscmk"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts = [
    "unwrapKey", "wrapKey"
  ]
}

resource "azurerm_managed_redis_cluster" "example" {
  name                = "example-managed-redis"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Balanced_B3"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  customer_managed_key {
    encryption_key_url        = azurerm_key_vault_key.example.id
    user_assigned_identity_id = azurerm_user_assigned_identity.example.id
  }
}
```

## Example Usage with Private Endpoint

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_managed_redis_cluster" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Balanced_B3"
}

resource "azurerm_private_endpoint" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id

  private_service_connection {
    name                           = "example"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_managed_redis_cluster.example.id
    subresource_names              = ["redisEnterprise"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Managed Redis Cluster. Changing this forces a new Managed Redis Cluster to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Redis Cluster should exist. Changing this forces a new Managed Redis Cluster to be created.

* `location` - (Required) The Azure Region where the Managed Redis Cluster should exist. Refer to [product availability documentation](https://azure.microsoft.com/explore/global-infrastructure/products-by-region/table) for supported locations. Changing this forces a new Managed Redis Cluster to be created.

* `sku_name` - (Required) The features and specification of the Managed Redis Cluster to deploy. Refer to [the documentation](https://learn.microsoft.com/en-us/rest/api/redis/redisenterprisecache/redis-enterprise/create?view=rest-redis-redisenterprisecache-2025-04-01&tabs=HTTP#skuname) for valid values. `Enterprise_` and `EnterpriseFlash_` prefixed SKUs are [no longer supported](https://learn.microsoft.com/azure/redis/migrate/migrate-overview). Changing this forces a new Managed Redis Cluster to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as defined below.

* `high_availability_enabled` - (Optional) Whether to enable high availability for the Managed Redis Cluster. Defaults to `true`. Changing this forces a new Managed Redis Cluster to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Managed Redis Cluster should be located. Only needed for legacy Redis Enterprise SKU. Changing this forces a new Managed Redis Cluster to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Managed Redis Cluster.

---

A `customer_managed_key` block supports the following:

* `encryption_key_url` - (Required) The URL of the Key Vault Key used for encryption. For example: `https://example-vault-name.vault.azure.net/keys/example-key-name/a1b2c3d4`. The id of [`azurerm_key_vault_key` resource](key_vault_key.html) can be used.

* `user_assigned_identity_id` - (Required) The ID of the User Assigned Identity that has access to the Key Vault Key.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Managed Redis Cluster. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Managed Redis Cluster.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis Cluster.

* `hostname` - DNS name of the cluster endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Redis Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Redis Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Redis Cluster.

## Import

Managed Redis Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_redis_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cache` - 2025-04-01
