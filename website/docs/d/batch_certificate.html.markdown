---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_certificate"
description: |-
  Get information about an existing certificate in a Batch Account

---

# Data Source: azurerm_batch_certificate

Use this data source to access information about an existing certificate in a Batch Account.

## Example Usage

```hcl
data "azurerm_batch_certificate" "example" {
  name                = "SHA1-42C107874FD0E4A9583292A2F1098E8FE4B2EDDA"
  account_name        = "examplebatchaccount"
  resource_group_name = "example"
}

output "thumbprint" {
  value = data.azurerm_batch_certificate.example.thumbprint
}
```

## Argument Reference

* `name` - The name of the Batch certificate.

* `account_name` - The name of the Batch account.

* `resource_group_name` - The Name of the Resource Group where this Batch account exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Batch certificate ID.

* `public_data` - The public key of the certificate.

* `format` - The format of the certificate, such as `Cer` or `Pfx`.

* `thumbprint` - The thumbprint of the certificate.

* `thumbprint_algorithm` - The algorithm of the certificate thumbprint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the certificate.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Batch`: 2024-07-01
