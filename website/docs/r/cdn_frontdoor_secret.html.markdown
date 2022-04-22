---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_secret"
description: |-
  Manages a Frontdoor Secret.
---

# azurerm_cdn_frontdoor_secret

Manages a Frontdoor Secret.

## Required Key Vault Permissions

!>**IMPORTANT:** You must add an `Access Policy` to your `azurerm_key_vault` for the `Microsoft.AzureFrontDoor-Cdn` Enterprise Application Object ID.

| Object ID                                | Key Permissions | Secret Permissions   | Certificate Permissions                       |
|:-----------------------------------------|:---------------:|:--------------------:|:---------------------------------------------:|
| `Microsoft.AzureFrontDoor-Cdn` Object ID | -               | **Get**              | -                                             |
| Your Personal AAD Object ID              | -               | **Get** and **List** | **Get**, **List**, **Purge** and **Recover**  |
| Terraform Service Principal              | -               | **Get**              | **Get**, **Import**, **Delete** and **Purge** |

->**NOTE:** You only need to add the `Access Policy` for your personal AAD Object ID if you are planning to view the `secrets` via the Azure Portal.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                       = "example-keyvault"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "premium"
  soft_delete_retention_days = 7

  network_acls {
    default_action = "Deny"
    bypass         = "AzureServices"
    ip_rules       = ["10.0.0.0/24"]
  }

  # Frontdoor Enterprise Application Object ID(e.g. Microsoft.AzureFrontDoor-Cdn)
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = "00000000-0000-0000-0000-000000000000" # <- Object Id for the Microsoft.AzureFrontDoor-Cdn Enterprise Application

    secret_permissions = [
      "Get",
    ]
  }

  # Terraform Service Principal
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = "00000000-0000-0000-0000-000000000000" # <- Object Id of the Service Principal that Terraform is running as

    certificate_permissions = [
      "Get",
      "Import",
      "Delete",
      "Purge"
    ]

    secret_permissions = [
      "Get",
    ]
  }
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "example-cert"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("my-certificate.pfx")
  }
}

resource "azurerm_cdn_frontdoor_secret" "example" {
  name                     = "example-customer-managed-secret"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.test.id

  secret_parameters {
    customer_certificate {
      key_vault_id                  = azurerm_key_vault_certificate.test.key_vault_id
      key_vault_certificate_name    = azurerm_key_vault_certificate.test.name
      key_vault_certificate_version = azurerm_key_vault_certificate.test.version
      use_latest                    = false
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Secret. Changing this forces a new Frontdoor Secret to be created.

* `cdn_frontdoor_profile_id` - (Required) The Resource ID of the Frontdoor Profile. Changing this forces a new Frontdoor Secret to be created.

* `secret_parameters` - (Required) A `secret_parameters` block as defined below. Changing this forces a new Frontdoor Secret to be created.

---

A `secret_parameters` block supports the following:

* `customer_certificate` - (Required) A `customer_certificate` block as defined below. Changing this forces a new Frontdoor Secret to be created.

---

A `customer_certificate` - (Required)  block supports the following:

* `key_vault_id`- (Required) The Resource ID of the Azure Key Vault which contains the certificate.

* `key_vault_certificate_name` - (Required) The Name of the Azure Key Vault certificate.
​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​
* `key_vault_certificate_version` - (Optional) The version of the Azure Key Vault certificate to be used.

->**NOTE:** The `key_vault_certificate_version` field should be removed from the configuration file if the `use_latest` field is set to `true`.

* `use_latest` - (Optional) Should the latest version of the certificate be used? Defaults to `true`.

* `subject_alternative_names` - (Computed) One or more `subject alternative names` contained within the key vault certificate.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Secret.

* `cdn_frontdoor_profile_name` - The name of the Frontdoor Profile containing this Frontdoor Secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Secret.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Secret.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Secret.

## Import

Frontdoor Secrets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_secret.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/secrets/secrets1
```
