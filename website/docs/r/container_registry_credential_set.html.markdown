---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_credential_set"
description: |-
  Manages a Container Registry Credential Set.
---

# azurerm_container_registry_credential_set

Manages a Container Registry Credential Set.

## Example Usage (minimal)

~> **Note:** Be aware that you will need to permit the Identity that is created for the Container Registry to have `get` on secrets to the Key Vault, e.g. using the `azurerm_key_vault_access_policy` resource.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                = "exampleContainerRegistry"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Basic"
}

resource "azurerm_container_registry_credential_set" "example" {
  name                  = "exampleCredentialSet"
  container_registry_id = azurerm_container_registry.example.id
  login_server          = "docker.io"
  identity {
    type = "SystemAssigned"
  }
  authentication_credentials {
    username_secret_id = "https://example-keyvault.vault.azure.net/secrets/example-user-name"
    password_secret_id = "https://example-keyvault.vault.azure.net/secrets/example-user-password"
  }
}
```

## Example Usage (full)

This example provisions a key vault with two secrets, a container registry, a container registry credential set, and an access policy to allow the container registry to read the secrets from the key vault.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                       = "examplekeyvault"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azurerm_client_config.current.object_id
    certificate_permissions = []
    key_permissions         = []
    secret_permissions      = ["Get", "Set", "Delete", "Purge"]
  }
}

resource "azurerm_key_vault_secret" "example_user" {
  key_vault_id = azurerm_key_vault.example.id
  name         = "example-user-name"
  value        = "name"
}

resource "azurerm_key_vault_secret" "example_password" {
  key_vault_id = azurerm_key_vault.example.id
  name         = "example-user-password"
  value        = "password"
}

resource "azurerm_container_registry" "example" {
  name                = "exampleContainerRegistry"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Basic"
}

resource "azurerm_container_registry_credential_set" "example" {
  name                  = "exampleCredentialSet"
  container_registry_id = azurerm_container_registry.example.id
  login_server          = "docker.io"
  identity {
    type = "SystemAssigned"
  }
  authentication_credentials {
    username_secret_id = azurerm_key_vault_secret.example_user.versionless_id
    password_secret_id = azurerm_key_vault_secret.example_password.versionless_id
  }
}

resource "azurerm_key_vault_access_policy" "read_secrets" {
  key_vault_id       = azurerm_key_vault.example.id
  tenant_id          = azurerm_container_registry_credential_set.example.identity[0].tenant_id
  object_id          = azurerm_container_registry_credential_set.example.identity[0].principal_id
  secret_permissions = ["Get"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Container Registry Credential Set. Changing this forces a new Container Registry Credential Set to be created.

* `container_registry_id` - (Required) The ID of the Container Registry. Changing this forces a new Container Registry Credential Set to be created.

* `login_server` - (Required) The login server for the Credential Set. Changing this forces a new Container Registry Credential Set to be created.

* `authentication_credentials` - (Required) A `authentication_credentials` block as defined below.

* `identity` - (Required) An `identity` block as defined below.

---

A `authentication_credentials` block supports the following:

* `username_secret_id` - (Required) The URI of the secret containing the username in a Key Vault.

* `password_secret_id` - (Required) The URI of the secret containing the password in a Key Vault.

~> **Note:** Be aware that you will need to permit the Identity that is created for the Container Registry to have `get` on secrets to the Key Vault, e.g. using the `azurerm_key_vault_access_policy` resource.

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that is configured on for the Container Registry Credential Set. Currently the only possible value is `SystemAssigned`.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Container Registry Credential Set.

---

A `identity` block exports the following:

* `principal_id` - The principal ID of the Identity.

* `tenant_id` - The tenant ID of the Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry Credential Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Credential Set.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry Credential Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Credential Set.

## Import

Container Registry Credential Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_credential_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/credentialSets/credentialSet1
```
