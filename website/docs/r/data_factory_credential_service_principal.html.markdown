---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_credential_service_principal"
description: |-
  Manage a Data Factory Service Principal credential resource
---

# azurerm_data_factory_credential_service_principal

Manage a Data Factory Service Principal credential resource. These resources are used by Data Factory to access data sources.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "westeurope"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_key_vault" "example" {
  name                       = "example"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
    ]

    secret_permissions = [
      "Set",
      "Get",
      "Delete",
      "Purge",
      "Recover"
    ]
  }
}

resource "azurerm_key_vault_secret" "example" {
  name         = "example"
  value        = "example-secret"
  key_vault_id = azurerm_key_vault.example.id
}

resource "azurerm_data_factory_linked_service_key_vault" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id
  key_vault_id    = azurerm_key_vault.example.id
}

resource "azurerm_data_factory_credential_service_principal" "example" {
  name                 = "example"
  description          = "example description"
  data_factory_id      = azurerm_data_factory.example.id
  tenant_id            = data.azurerm_client_config.current.tenant_id
  service_principal_id = data.azurerm_client_config.current.client_id
  service_principal_key {
    linked_service_name = azurerm_data_factory_linked_service_key_vault.example.name
    secret_name         = azurerm_key_vault_secret.example.name
    secret_version      = azurerm_key_vault_secret.example.version
  }
  annotations = ["1", "2"]
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Credential. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The Data Factory ID in which to associate the Credential with. Changing this forces a new resource.

* `tenant_id` - (Required) The Tenant ID of the Service Principal.

* `service_principal_id` - (Required) The Client ID of the Service Principal.

* `service_principal_key` - (Required) A `service_principal_key` block as defined below.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Credential.

* `description` - (Optional) The description for the Data Factory Credential.

---

A `service_principal_key` block supports the following:

* `linked_service_name` - (Required) The name of the Linked Service to use for the Service Principal Key.

* `secret_name` - (Required) The name of the Secret in the Key Vault.

* `secret_version` - (Optional) The version of the Secret in the Key Vault.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Credential.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Data Factory Credential.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Credential.
* `update` - (Defaults to 5 minutes) Used when updating the Data Factory Credential.
* `delete` - (Defaults to 5 minutes) Used when deleting the Data Factory Credential.

## Import

Data Factory Credentials can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_credential_service_principal.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.DataFactory/factories/example/credentials/credential1
```
