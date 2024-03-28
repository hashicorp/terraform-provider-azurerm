---
subcategory: "Nginx"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_nginx_certificate"
description: |-
  Gets information about an existing Nginx Certificate.
---

# Data Source: azurerm_nginx_certificate

Use this data source to access information about an existing Nginx Certificate.

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

* `name` - (Required) The name of this Nginx Certificate.

* `nginx_deployment_id` - (Required) The ID of the Nginx Deployment that this certificate is associated with.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Nginx Certificate.

* `certificate_virtual_path` - The path to the certificate file of this certificate.

* `key_virtual_path` - The path to the key file of this certificate.

* `key_vault_secret_id` - The ID of the Key Vault Secret for this certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Nginx Certificate.
