---
subcategory: "Trusted Signing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_trusted_signing_account"
description: |-
  Gets information about an existing Trusted Signing Account.
---

# Data Source: azurerm_trusted_signing_account

Use this data source to access information about an existing Trusted Signing Account.

## Example Usage

```hcl

data "azurerm_trusted_signing_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resource-group"
}

output "trusted_signing_account_id" {
  value = data.azurerm_trusted_signing_account.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Trusted Signing Account.

* `resource_group_name` - (Required) The name of the Resource Group where the Trusted Signing Account exists.

## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Trusted Signing Account.

* `location` - The Azure Region where the Trusted Signing Account exists.

* `account_uri` - The URI of the Trusted Signing Account.

* `sku_name` - The sku name of the Trusted Signing Account.

* `tags` - A mapping of tags assigned to the Trusted Signing Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Trusted Signing Account.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.CodeSigning` - 2024-09-30-preview
