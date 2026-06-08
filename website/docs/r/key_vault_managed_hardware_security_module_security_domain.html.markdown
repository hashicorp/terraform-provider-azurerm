---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module_security_domain"
description: |-
  Manages a Key Vault Managed Hardware Security Module Security Domain.
---

# azurerm_key_vault_managed_hardware_security_module_security_domain

Manages a Key Vault Managed Hardware Security Module Security Domain.

~> **Note:** The Security Domain download is an activation step that enables the Managed HSM. There is no Azure API operation to "delete" or "deactivate" a Security Domain once downloaded. Therefore, deleting this resource from Terraform will only remove it from the Terraform state and will not modify the underlying Managed HSM.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault_managed_hardware_security_module" "example" {
  name                       = "exampleKVHsm"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  sku_name                   = "Standard_B1"
  purge_protection_enabled   = false
  soft_delete_retention_days = 90
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]

  tags = {
    Env = "Test"
  }
}

resource "azurerm_key_vault" "example" {
  name                       = "exampleKV"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id
    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
}

resource "azurerm_key_vault_certificate" "example" {
  count        = 3
  name         = "example-cert-${count.index}"
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
      extended_key_usage = []
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}

resource "azurerm_key_vault_managed_hardware_security_module_security_domain" "example" {
  managed_hsm_id                            = azurerm_key_vault_managed_hardware_security_module.example.id
  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.example : cert.id]
  security_domain_quorum                    = 2
}
```

## Arguments Reference

The following arguments are supported:

* `managed_hsm_id` - (Required) The ID of the Key Vault Managed Hardware Security Module that should be activated. Changing this forces a new resource to be created.

* `security_domain_key_vault_certificate_ids` - (Required) A list of KeyVault certificates resource IDs (minimum of three and up to a maximum of 10) to activate this Managed HSM. More information see [activate-your-managed-hsm](https://learn.microsoft.com/azure/key-vault/managed-hsm/quick-create-cli#activate-your-managed-hsm)

* `security_domain_quorum` - (Required) Specifies the minimum number of shares required to decrypt the security domain for recovery. Valid values are between 2 and 10.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Key Vault Managed Hardware Security Module.

* `security_domain_encrypted_data` - This attribute can be used for disaster recovery or when creating another Managed HSM that shares the same security domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Key Vault Managed Hardware Security Module Security Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Key Vault Managed Hardware Security Module Security Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Key Vault Managed Hardware Security Module Security Domain.
* `delete` - (Defaults to 5 minutes) Used when deleting the Key Vault Managed Hardware Security Module Security Domain.

## Import

Key Vault Managed Hardware Security Module Security Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module_security_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.KeyVault/managedHSMs/hsm1
```
