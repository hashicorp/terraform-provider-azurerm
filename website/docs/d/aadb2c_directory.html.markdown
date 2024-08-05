---
subcategory: "AAD B2C"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_aadb2c_directory"
description: |-
  Gets information about an existing AAD B2C Directory.
---

# Data Source: azurerm_aadb2c_directory

Use this data source to access information about an existing AAD B2C Directory.

## Example Usage

```hcl
data "azurerm_aadb2c_directory" "example" {
  resource_group_name = "example-rg"
  domain_name         = "exampleb2ctenant.onmicrosoft.com"
}

output "tenant_id" {
  value = data.azurerm_aadb2c_directory.example.tenant_id
}
```

## Arguments Reference

The following arguments are supported:

* `domain_name` - (Required) Domain name of the B2C tenant, including the `.onmicrosoft.com` suffix.

* `resource_group_name` - (Required) The name of the Resource Group where the AAD B2C Directory exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the AAD B2C Directory.

* `billing_type` - The type of billing for the AAD B2C tenant. Possible values include: `MAU` or `Auths`.

* `data_residency_location` - Location in which the B2C tenant is hosted and data resides. See [official docs](https://aka.ms/B2CDataResidenc) for more information.

* `effective_start_date` - The date from which the billing type took effect. May not be populated until after the first billing cycle.

* `sku_name` - Billing SKU for the B2C tenant. See [official docs](https://aka.ms/b2cBilling) for more information.

* `tags` - A mapping of tags assigned to the AAD B2C Directory.

* `tenant_id` - The Tenant ID for the AAD B2C tenant.

~> **Note:** The `country_code` and `display_name` are not returned by this data source due to API limitations.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the AAD B2C Directory.
