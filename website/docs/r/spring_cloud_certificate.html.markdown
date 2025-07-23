---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_certificate"
description: |-
  Manages an Azure Spring Cloud Certificate.
---

# azurerm_spring_cloud_certificate

Manages an Azure Spring Cloud Certificate.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_certificate` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "example" {
  display_name = "Azure Spring Cloud Resource Provider"
}

resource "azurerm_key_vault" "example" {
  name                = "keyvaultcertexample"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azurerm_client_config.current.object_id
    secret_permissions      = ["Set"]
    certificate_permissions = ["Create", "Delete", "Get", "Update"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azuread_service_principal.example.object_id
    secret_permissions      = ["Get", "List"]
    certificate_permissions = ["Get", "List"]
  }
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "cert-example"
  key_vault_id = azurerm_key_vault.example.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=contoso.com"
      validity_in_months = 12
    }
  }
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_spring_cloud_certificate" "example" {
  name                     = "example-scc"
  resource_group_name      = azurerm_spring_cloud_service.example.resource_group_name
  service_name             = azurerm_spring_cloud_service.example.name
  key_vault_certificate_id = azurerm_key_vault_certificate.example.id
  exclude_private_key      = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Certificate. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group in which to create the Spring Cloud Certificate. Changing this forces a new resource to be created.

* `service_name` - (Required) Specifies the name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `exclude_private_key` - (Optional) Specifies whether the private key should be excluded from the Key Vault Certificate. Changing this forces a new resource to be created. Defaults to `false`.

* `key_vault_certificate_id` - (Optional) Specifies the ID of the Key Vault Certificate resource. Changing this forces a new resource to be created.

* `certificate_content` - (Optional) The content of uploaded certificate. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Certificate.

* `thumbprint` - The thumbprint of the Spring Cloud certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Certificate.

## Import

Spring Cloud Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AppPlatform/spring/spring1/certificates/cert1
```
