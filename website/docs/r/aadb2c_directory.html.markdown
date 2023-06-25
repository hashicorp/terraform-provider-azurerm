---
subcategory: "AAD B2C"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_aadb2c_directory"
description: |-
  Manages an AAD B2C Directory.
---

# azurerm_aadb2c_directory

Manages an AAD B2C Directory.

## Example Usage

```hcl
resource "azurerm_aadb2c_directory" "example" {
  country_code            = "US"
  data_residency_location = "United States"
  display_name            = "example-b2c-tenant"
  domain_name             = "exampleb2ctenant.onmicrosoft.com"
  resource_group_name     = "example-rg"
  sku_name                = "PremiumP1"
}
```

## Arguments Reference

The following arguments are supported:

* `country_code` - (Optional) Country code of the B2C tenant. The `country_code` should be valid for the specified `data_residency_location`. See [official docs](https://aka.ms/B2CDataResidency) for valid country codes. Required when creating a new resource. Changing this forces a new AAD B2C Directory to be created.

* `data_residency_location` - (Required) Location in which the B2C tenant is hosted and data resides. The `data_residency_location` should be valid for the specified `country_code`. See [official docs](https://aka.ms/B2CDataResidenc) for more information. Changing this forces a new AAD B2C Directory to be created. Possible values are `Asia Pacific`, `Australia`, `Europe`, `Global` and `United States`.

* `display_name` - (Optional) The initial display name of the B2C tenant. Required when creating a new resource. Changing this forces a new AAD B2C Directory to be created.

* `domain_name` - (Required) Domain name of the B2C tenant, including the `.onmicrosoft.com` suffix. Changing this forces a new AAD B2C Directory to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the AAD B2C Directory should exist. Changing this forces a new AAD B2C Directory to be created.

* `sku_name` - (Required) Billing SKU for the B2C tenant. Must be one of: `PremiumP1` or `PremiumP2` (`Standard` is not supported). See [official docs](https://aka.ms/b2cBilling) for more information.

* `tags` - (Optional) A mapping of tags which should be assigned to the AAD B2C Directory.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the AAD B2C Directory.

* `billing_type` - The type of billing for the AAD B2C tenant. Possible values include: `MAU` or `Auths`.

* `effective_start_date` - The date from which the billing type took effect. May not be populated until after the first billing cycle.

* `tenant_id` - The Tenant ID for the AAD B2C tenant.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the AAD B2C Directory.
* `read` - (Defaults to 5 minutes) Used when retrieving the AAD B2C Directory.
* `update` - (Defaults to 30 minutes) Used when updating the AAD B2C Directory.
* `delete` - (Defaults to 30 minutes) Used when deleting the AAD B2C Directory.

## Import

AAD B2C Directories can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_aadb2c_directory.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resource-group/providers/Microsoft.AzureActiveDirectory/b2cDirectories/directory-name
```
