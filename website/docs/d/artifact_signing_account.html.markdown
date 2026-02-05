---
subcategory: "Artifact Signing"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_artifact_signing_account"
description: |-
  Gets information about an existing Artifact Signing Account.
---

# Data Source: azurerm_artifact_signing_account

Use this data source to access information about an existing Artifact Signing Account.

## Example Usage

```hcl

data "azurerm_artifact_signing_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resource-group"
}

output "artifact_signing_account_id" {
  value = data.azurerm_artifact_signing_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Artifact Signing Account.

* `resource_group_name` - (Required) The name of the Resource Group where the Artifact Signing Account exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Artifact Signing Account.

* `location` - The Azure Region where the Artifact Signing Account exists.

* `account_uri` - The URI of the Artifact Signing Account.

* `sku_name` - The sku name of the Artifact Signing Account.

* `tags` - A mapping of tags assigned to the Artifact Signing Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Artifact Signing Account.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.CodeSigning` - 2025-10-13
