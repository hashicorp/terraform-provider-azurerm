---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_certificate_order_key_vault_store"
description: |-
  Manages an App Service Certificate Order Certificate.

---

# azurerm_app_service_certificate_order_key_vault_store

Manages an App Service Certificate Order Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "test" {}

resource "azurerm_key_vault" "test" {
  name                = "example-keyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tenant_id = data.azurerm_client_config.test.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = "f8daea97-62e7-4026-becf-13c2ea98e8b4"

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = "ed47c2a1-bd23-4341-b39c-f4fd69138dd3"

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }
}

resource "azurerm_app_service_certificate_order" "example" {
  name                = "example-cert-order"
  resource_group_name = azurerm_resource_group.example.name
  location            = "global"
  distinguished_name  = "CN=example.com"
  product_type        = "Standard"
}

resource "azurerm_app_service_certificate_order_key_vault_store" "test" {
  name                  = "example-certorder-cert"
  certificate_order_id  = azurerm_app_service_certificate_order.example.id
  key_vault_id          = azurerm_key_vault.example.id
  key_vault_secret_name = "example-keyvault-secret"
}
```

-> **Note:** Please make sure the domain ownership is verified before configure the key vault.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Certificate Key Vault Store. Changing this forces a new resource to be created.

* `certificate_order_id` - (Required) The ID of the Certificate Order in which to configure the Certificate Key Vault Store Binding. Changing this forces a new resource to be created.

* `key_vault_id` - (Required) The ID of the Key Vault in which to bind the Certificate.

* `key_vault_secret_name` - (Required) The name of the Key Vault Secret to bind to the Certificate.

## Attributes Reference

* `location` - The location of the certificate.

* `type` - The type of the certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Certificate Order Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Certificate Order Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Certificate Order Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Certificate Order Certificate.

## Import

App Service Certificate Order Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_certificate_order_key_vault_store.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.CertificateRegistration/certificateOrders/certificateorder1/certificates/certificates1
```


