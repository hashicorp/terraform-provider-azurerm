---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_secret"
description: |-
  Manages a Front Door (standard/premium) Secret.
---

# azurerm_cdn_frontdoor_secret

Manages a Front Door (standard/premium) Secret.

## Required Key Vault Permissions

!> **Note:** You must add an `Access Policy` to your `azurerm_key_vault` for the `Microsoft.AzurefrontDoor-Cdn` Enterprise Application Object ID.

This can be created by running Az Powershell command like this:

```New-AzADServicePrincipal -ApplicationId "00000000-0000-0000-0000-000000000000"```

| Object ID                                | Key Permissions | Secret Permissions   | Certificate Permissions                       |
|:-----------------------------------------|:---------------:|:--------------------:|:---------------------------------------------:|
| `Microsoft.Azure.Cdn` Object ID          | -               | **Get**              | -                                             |
| Your Personal AAD Object ID              | -               | **Get** and **List** | **Get**, **List**, **Purge** and **Recover**  |
| Terraform Service Principal              | -               | **Get**              | **Get**, **Import**, **Delete** and **Purge** |

-> **Note:** You only need to add the `Access Policy` for your personal AAD Object ID if you are planning to view the `secrets` via the Azure Portal.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}
data "azuread_service_principal" "frontdoor" {
  display_name = "Microsoft.AzurefrontDoor-Cdn"
}

resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

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

  # CDN Front Door Enterprise Application Object ID(e.g. Microsoft.Azure.Cdn)
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.frontdoor.object_id

    secret_permissions = [
      "Get",
    ]
  }

  # Terraform Service Principal
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id # <- Object Id of the Service Principal that Terraform is running as

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
  key_vault_id = azurerm_key_vault.example.id

  certificate {
    contents = filebase64("my-certificate.pfx")
  }
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-cdn-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_secret" "example" {
  name                     = "example-customer-managed-secret"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  secret {
    customer_certificate {
      key_vault_certificate_id = azurerm_key_vault_certificate.example.id
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Front Door Secret. Possible values must start with a letter or a number, only contain letters, numbers and hyphens and have a length of between 2 and 260 characters. Changing this forces a new Front Door Secret to be created.

* `cdn_frontdoor_profile_id` - (Required) The Resource ID of the Front Door Profile. Changing this forces a new Front Door Secret to be created.

* `secret` - (Required) A `secret` block as defined below. Changing this forces a new Front Door Secret to be created.

---

A `secret` block supports the following:

* `customer_certificate` - (Required) A `customer_certificate` block as defined below. Changing this forces a new Front Door Secret to be created.

---

A `customer_certificate` block supports the following:

* `key_vault_certificate_id` - (Required) The ID of the Key Vault certificate resource to use. Changing this forces a new Front Door Secret to be created.

-> **Note:** If you would like to use the **latest version** of the Key Vault Certificate use the Key Vault Certificates `versionless_id` attribute as the `key_vault_certificate_id` fields value(e.g. `key_vault_certificate_id = azurerm_key_vault_certificate.example.versionless_id`).

* `subject_alternative_names` - (Computed) One or more `subject alternative names` contained within the key vault certificate.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Front Door Secret.

* `cdn_frontdoor_profile_name` - The name of the Front Door Profile containing this Front Door Secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Secret.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Secret.
* `delete` - (Defaults to 30 minutes) Used when deleting the Front Door Secret.

## Import

Front Door Secrets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_secret.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/secrets/secrets1
```
