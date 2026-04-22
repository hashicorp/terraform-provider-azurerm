---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_certificate"
description: |-
  Manages a Container App Environment Certificate.
---

# azurerm_container_app_environment_certificate

Manages a Container App Environment Certificate.

## Example Usage

### Certificate from .pfx file

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "myEnvironment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_container_app_environment_certificate" "example" {
  name                         = "myfriendlyname"
  container_app_environment_id = azurerm_container_app_environment.example.id
  certificate_blob_base64      = filebase64("path/to/certificate_file.pfx")
  certificate_password         = "$3cretSqu1rreL"
}
```

### Certificate from Key Vault

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_container_app_environment" "example" {
  name                       = "example-environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }
}

resource "azurerm_key_vault" "example" {
  name                      = "example-keyvault"
  location                  = azurerm_resource_group.example.location
  resource_group_name       = azurerm_resource_group.example.name
  tenant_id                 = data.azurerm_client_config.current.tenant_id
  sku_name                  = "standard"
  enable_rbac_authorization = true
}

resource "azurerm_role_assignment" "user_keyvault_admin" {
  scope                = azurerm_key_vault.example.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_key_vault.example.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_container_app_environment.example.identity[0].principal_id
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "example-certificate"
  key_vault_id = azurerm_key_vault.example.id

  certificate {
    contents = filebase64("path/to/certificate_file.pfx")
    password = ""
  }

  depends_on = [azurerm_role_assignment.user_keyvault_admin, azurerm_role_assignment.example]
}

resource "azurerm_container_app_environment_certificate" "example" {
  name                         = "example-certificate"
  container_app_environment_id = azurerm_container_app_environment.example.id

  certificate_key_vault {
    identity            = azurerm_user_assigned_identity.example.id
    key_vault_secret_id = azurerm_key_vault_certificate.example.versionless_secret_id
  }

  depends_on = [azurerm_role_assignment.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Environment Certificate. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The Container App Managed Environment ID to configure this Certificate on. Changing this forces a new resource to be created.

* `certificate_blob_base64` - (Optional) The Certificate Private Key as a base64 encoded PFX or PEM. Changing this forces a new resource to be created.

~> **Note:** One of `certificate_blob_base64` and `certificate_key_vault` must be set.

* `certificate_password` - (Optional) The password for the Certificate. Changing this forces a new resource to be created.

~> **Note:** required if `certificate_blob_base64` is specified.

* `certificate_key_vault` - (Optional) A `certificate_key_vault` block as defined below. Changing this forces a new resource to be created.

~> **Note:** one of `certificate_blob_base64` and `certificate_key_vault` must be set.

---

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `certificate_key_vault` block supports the following:

* `identity` - (Optional) The managed identity to authenticate with Azure Key Vault. Possible values are the resource ID of user-assigned identity, and `System` for system-assigned identity. Defaults to `System`. Changing this forces a new resource to be created.

~> **Note:** Please make sure [required permissions](https://learn.microsoft.com/en-us/azure/container-apps/key-vault-certificates-manage) are correctly configured for your Key Vault and managed identity.

* `key_vault_secret_id` - (Required) The ID of the Key Vault Secret containing the certificate. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Certificate

* `expiration_date` - The expiration date for the Certificate.

* `issue_date` - The date of issue for the Certificate.

* `issuer` - The Certificate Issuer.

* `subject_name` - The Subject Name for the Certificate.

* `thumbprint` - The Thumbprint of the Certificate.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment Certificate.

## Import

A Container App Environment Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_certificate.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myenv/certificates/mycertificate"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
