---
subcategory: "Attestation"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_attestation"
description: |-
  Gets information about an existing Attestation Provider.
---

# Data Source: azurerm_attestation

Use this data source to access information about an existing Attestation Provider.

## Example Usage

```hcl
data "azurerm_attestation" "example" {
  name = "example-attestationprovider"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_attestation.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this attestation provider.

* `resource_group_name` - (Required) The name of the Resource Group where the attestation provider exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the attestation provider.

* `location` - The Azure Region where the attestation provider exists.

* `tags` - A mapping of tags assigned to the attestation provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the attestation provider.
