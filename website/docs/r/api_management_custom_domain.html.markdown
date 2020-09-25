---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_custom_domain"
description: |-
  Manages a API Management Custom Domain.
---

# azurerm_api_management_custom_domain

Manages a API Management Custom Domain.

## Disclaimers

~> **Note:** It's possible to define Custom Domains both within [the `azurerm_api_management` resource](api_management.html) via the `hostname_configurations` block and by using this resource. However it's not possible to use both methods to manage Custom Domains within an API Management Service, since there will be conflicts.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "example-certificate"
  key_vault_id = azurerm_key_vault.test.id

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

      subject            = "CN=api.example.com"
      validity_in_months = 12

      subject_alternative_names {
        dns_names = [
          "api.example.com",
          "portal.example.com",
        ]
      }
    }
  }
}

resource "azurerm_api_management_custom_domain" "example" {
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name

  proxy {
    host_name    = "api.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }

  developer_portal {
    host_name    = "portal.example.com"
    key_vault_id = azurerm_key_vault_certificate.test.secret_id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_name` - (Required) The name of the API Management service to configure Custom Domain for. Changing this forces a new API Management Custom Domain to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the API Management resource exists. Changing this forces a new API Management Custom Domain to be created.

---

* `developer_portal` - (Optional) One or more `developer_portal` blocks as defined below.

* `management` - (Optional) One or more `management` blocks as defined below.

* `portal` - (Optional) One or more `portal` blocks as defined below.

* `proxy` - (Optional) One or more `proxy` blocks as defined below.

* `scm` - (Optional) One or more `scm` blocks as defined below.

---

A `developer_portal`, `management`, `portal` or `scm` block supports the following:

* `host_name` - (Required) The Hostname to use for the corresponding endpoint.

* `certificate` - (Optional) The Base64 Encoded Certificate. (Mutually exlusive with `key_vault_id`.)

* `certificate_password` - (Optional) The password associated with the certificate provided above.

* `key_vault_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be should be of the type application/x-pkcs12.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to false.

---

A `proxy` block supports the following:

* `host_name` - (Required) The Hostname to use for the API Proxy Endpoint.

* `certificate` - (Optional) The Base64 Encoded Certificate. (Mutually exlusive with `key_vault_id`.)

* `certificate_password` - (Optional) The password associated with the certificate provided above.

* `default_ssl_binding` - (Optional) Is the certificate associated with this Hostname the Default SSL Certificate? This is used when an SNI header isn't specified by a client. Defaults to false.

* `key_vault_id` - (Optional) The ID of the Key Vault Secret containing the SSL Certificate, which must be should be of the type application/x-pkcs12.

* `negotiate_client_certificate` - (Optional) Should Client Certificate Negotiation be enabled for this Hostname? Defaults to false.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Custom Domain.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Custom Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Custom Domain.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Custom Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Custom Domain.

## Import

API Management Custom Domains can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_custom_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1
```
