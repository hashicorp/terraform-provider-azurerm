---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry"
description: |-
  Manages an Azure Container Registry.

---

# azurerm_container_registry

Manages an Azure Container Registry.

~> **Note:** All arguments including the access key will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_registry" "acr" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
  admin_enabled       = false
  georeplications {
    location                = "East US"
    zone_redundancy_enabled = true
    tags                    = {}
  }
  georeplications {
    location                = "North Europe"
    zone_redundancy_enabled = true
    tags                    = {}
  }
}
```

## Example Usage (Encryption)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_registry" "acr" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }

  encryption {
    key_vault_key_id   = data.azurerm_key_vault_key.example.id
    identity_client_id = azurerm_user_assigned_identity.example.client_id
  }

}

resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  name = "registry-uai"
}

data "azurerm_key_vault_key" "example" {
  name         = "super-secret"
  key_vault_id = data.azurerm_key_vault.existing.id
}



```

## Example Usage (Attaching a Container Registry to a Kubernetes Cluster)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-aks1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "exampleaks1"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Production"
  }
}

resource "azurerm_role_assignment" "example" {
  principal_id                     = azurerm_kubernetes_cluster.example.kubelet_identity[0].object_id
  role_definition_name             = "AcrPull"
  scope                            = azurerm_container_registry.example.id
  skip_service_principal_aad_check = true
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container Registry. Only Alphanumeric characters allowed. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Registry. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) The SKU name of the container registry. Possible values are `Basic`, `Standard` and `Premium`.

* `admin_enabled` - (Optional) Specifies whether the admin user is enabled. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `georeplications` - (Optional) One or more `georeplications` blocks as documented below.

~> **Note:** The `georeplications` is only supported on new resources with the `Premium` SKU.

~> **Note:** The `georeplications` list cannot contain the location where the Container Registry exists.

~> **Note:** If more than one `georeplications` block is specified, they are expected to follow the alphabetic order on the `location` property.

* `network_rule_set` - (Optional) A `network_rule_set` block as documented below.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for the container registry. Defaults to `true`.

* `quarantine_policy_enabled` - (Optional) Boolean value that indicates whether quarantine policy is enabled.

* `retention_policy_in_days` - (Optional) The number of days to retain and untagged manifest after which it gets purged. Defaults to `7`.

* `trust_policy_enabled` - (Optional) Boolean value that indicated whether trust policy is enabled. Defaults to `false`.

* `zone_redundancy_enabled` - (Optional) Whether zone redundancy is enabled for this Container Registry? Changing this forces a new resource to be created. Defaults to `false`.

* `export_policy_enabled` - (Optional) Boolean value that indicates whether export policy is enabled. Defaults to `true`. In order to set it to `false`, make sure the `public_network_access_enabled` is also set to `false`.

~> **Note:** `quarantine_policy_enabled`, `retention_policy_in_days`, `trust_policy_enabled`, `export_policy_enabled` and `zone_redundancy_enabled` are only supported on resources with the `Premium` SKU.

* `identity` - (Optional) An `identity` block as defined below.

* `encryption` - (Optional) An `encryption` block as documented below.

* `anonymous_pull_enabled` - (Optional) Whether allows anonymous (unauthenticated) pull access to this Container Registry? This is only supported on resources with the `Standard` or `Premium` SKU.

* `data_endpoint_enabled` - (Optional) Whether to enable dedicated data endpoints for this Container Registry? This is only supported on resources with the `Premium` SKU.

* `network_rule_bypass_option` - (Optional) Whether to allow trusted Azure services to access a network restricted Container Registry? Possible values are `None` and `AzureServices`. Defaults to `AzureServices`.

---

The `georeplications` block supports the following:

* `location` - (Required) A location where the container registry should be geo-replicated.

* `regional_endpoint_enabled` - (Optional) Whether regional endpoint is enabled for this Container Registry?

* `zone_redundancy_enabled` - (Optional) Whether zone redundancy is enabled for this replication location? Defaults to `false`.

~> **Note:** Changing the `zone_redundancy_enabled` forces the a underlying replication to be created.

* `tags` - (Optional) A mapping of tags to assign to this replication location.

---

The `network_rule_set` block supports the following:

* `default_action` - (Optional) The behaviour for requests matching no rules. Either `Allow` or `Deny`. Defaults to `Allow`

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

~> **Note:** `network_rule_set` is only supported with the `Premium` SKU at this time.

~> **Note:** Azure automatically configures Network Rules - to remove these you'll need to specify an `network_rule_set` block with `default_action` set to `Deny`.

---

The `ip_rule` block supports the following:

* `action` - (Required) The behaviour for requests matching this rule. At this time the only supported value is `Allow`

* `ip_range` - (Required) The CIDR block from which requests will match the rule.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Container Registry. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Container Registry.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

The `encryption` block supports the following:

* `key_vault_key_id` - (Required) The ID of the Key Vault Key.

* `identity_client_id` - (Required) The client ID of the managed identity associated with the encryption key.

~> **Note:** The managed identity used in `encryption` also needs to be part of the `identity` block under `identity_ids`

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry.

* `login_server` - The URL that can be used to log into the container registry.

* `admin_username` - The Username associated with the Container Registry Admin account - if the admin account is enabled.

* `admin_password` - The Password associated with the Container Registry Admin account - if the admin account is enabled.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

-> **Note:** You can access the Principal ID via `azurerm_container_registry.example.identity[0].principal_id` and the Tenant ID via `azurerm_container_registry.example.identity[0].tenant_id`

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry.

## Import

Container Registries can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1
```
