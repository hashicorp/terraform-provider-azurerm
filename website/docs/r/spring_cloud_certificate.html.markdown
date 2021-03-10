---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_certificate"
description: |-
  Manages an Azure Spring Cloud Certificate.
---

# azurerm_spring_cloud_certificate

Manages an Azure Spring Cloud Certificate.

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
  display_name = "Azure Spring Cloud Domain-Management"
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
    secret_permissions      = ["set"]
    certificate_permissions = ["create", "delete", "get", "update"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.current.tenant_id
    object_id               = data.azuread_service_principal.example.object_id
    secret_permissions      = ["get", "list"]
    certificate_permissions = ["get", "list"]
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
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Certificate. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group in which to create the Spring Cloud Certificate. Changing this forces a new resource to be created.

* `service_name` - (Required) Specifies the name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `key_vault_certificate_id` - (Required) Specifies the ID of the Key Vault Certificate resource. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Certificate.

* `thumbprint` - The thumbprint of the Spring Cloud certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Certificate.

## Import

Spring Cloud Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AppPlatform/Spring/spring1/certificates/cert1
```
