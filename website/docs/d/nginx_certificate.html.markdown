---
subcategory: "NGINX"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_nginx_certificate"
description: |-
  Gets information about an existing NGINX Certificate.
---

# Data Source: azurerm_nginx_certificate

Use this data source to access information about an existing NGINX Certificate.

## Example Usage

```hcl
data "azurerm_nginx_certificate" "example" {
  name                = "existing"
  nginx_deployment_id = azurerm_nginx_deployment.example.id
}

output "id" {
  value = data.azurerm_nginx_certificate.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the NGINX Certificate.

* `nginx_deployment_id` - (Required) The ID of the NGINX Deployment that the certificate is associated with.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NGINX Certificate.

* `certificate_virtual_path` - The path to the certificate file of the certificate.

* `key_virtual_path` - The path to the key file of the certificate.

* `key_vault_secret_id` - The ID of the Key Vault Secret for the certificate.

* `sha1_thumbprint` - The SHA-1 thumbprint of the certificate.

* `key_vault_secret_version` - The version of the certificate.

* `key_vault_secret_creation_date` - The date/time the certificate was created in Azure Key Vault.

* `error_code` - The error code of the certificate error, if any.

* `error_message` - The error message of the certificate error, if any.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NGINX Certificate.
