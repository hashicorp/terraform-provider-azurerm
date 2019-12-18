---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_certificate_order"
sidebar_current: "docs-azurerm-datasource-app-service-x"
description: |-
  Gets information about an existing App Service Certificate Order.
---

# Data Source: azurerm_app_service_certificate_order

Use this data source to access information about an existing App Service Certificate Order.

## Example Usage

```hcl
data "azurerm_app_service_certificate_order" "example" {
  name                = "example-cert-order"
  resource_group_name = "example-resources"
}

output "certificate_order_id" {
  value = "${data.azurerm_app_service_certificate_order.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the App Service.

* `resource_group_name` - (Required) The Name of the Resource Group where the App Service exists.

## Attributes Reference

* `id` - The ID of the App Service.

* `location` - The Azure location where the App Service exists.

* `auto_renew` - true if the certificate should be automatically renewed when it expires; otherwise, false.

* `certificates` - State of the Key Vault secret. A `certificates` block as defined below.

* `csr` - Last CSR that was created for this order.

* `distinguished_name` - The Distinguished Name for the App Service Certificate Order.

* `key_size` - Certificate key size.

* `product_type` - Certificate product type, such as `Standard` or `WildCard`.

* `validity_in_years` - Duration in years (must be between 1 and 3).

* `domain_verification_token` - Domain verification token.

* `status` - Current order status.

* `expiration_time` - Certificate expiration time.

* `is_private_key_external` - Whether the private key is external or not.

* `app_service_certificate_not_renewable_reasons` - Reasons why App Service Certificate is not renewable at the current moment.

* `signed_certificate_thumbprint` - Certificate thumbprint for signed certificate.

* `root_thumbprint` - Certificate thumbprint for root certificate.

* `intermediate_thumbprint` - Certificate thumbprint intermediate certificate.

* `tags` - A mapping of tags to assign to the resource.

---

* `certificate_name` - The name of the App Service Certificate.

* `key_vault_id` - Key Vault resource Id.

* `key_vault_secret_name` - Key Vault secret name.

* `provisioning_state` - Status of the Key Vault secret.
