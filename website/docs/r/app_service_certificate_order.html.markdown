---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_certificate_order"
description: |-
  Manages an App Service Certificate Order.

---

# azurerm_app_service_certificate_order

Manages an App Service Certificate Order.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_certificate_order" "example" {
  name                = "example-cert-order"
  resource_group_name = azurerm_resource_group.example.name
  location            = "global"
  distinguished_name  = "CN=example.com"
  product_type        = "Standard"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the certificate. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the certificate. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created. Currently the only valid value is `global`.

* `auto_renew` - (Optional) true if the certificate should be automatically renewed when it expires; otherwise, false. Defaults to `true`.

* `csr` - (Optional) Last CSR that was created for this order.

* `distinguished_name` - (Optional) The Distinguished Name for the App Service Certificate Order.

-> **Note:** Either `csr` or `distinguished_name` must be set - but not both.

* `key_size` - (Optional) Certificate key size. Defaults to `2048`.

* `product_type` - (Optional) Certificate product type, such as `Standard` or `WildCard`. Defaults to `Standard`.

* `validity_in_years` - (Optional) Duration in years (must be between `1` and `3`). Defaults to `1`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The App Service Certificate Order ID.

* `certificates` - State of the Key Vault secret. A `certificates` block as defined below.

* `domain_verification_token` - Domain verification token.

* `status` - Current order status.

* `expiration_time` - Certificate expiration time.

* `is_private_key_external` - Whether the private key is external or not.

* `app_service_certificate_not_renewable_reasons` - Reasons why App Service Certificate is not renewable at the current moment.

* `signed_certificate_thumbprint` - Certificate thumbprint for signed certificate.

* `root_thumbprint` - Certificate thumbprint for root certificate.

* `intermediate_thumbprint` - Certificate thumbprint intermediate certificate.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `certificates` block supports the following:

* `certificate_name` - The name of the App Service Certificate.

* `key_vault_id` - Key Vault resource Id.

* `key_vault_secret_name` - Key Vault secret name.

* `provisioning_state` - Status of the Key Vault secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Certificate Order.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Certificate Order.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Certificate Order.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Certificate Order.

## Import

App Service Certificate Orders can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_certificate_order.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.CertificateRegistration/certificateOrders/certificateorder1
```
