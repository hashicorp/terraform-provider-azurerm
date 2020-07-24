---
subcategory: "attestation"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_attestation_attestation_provider"
description: |-
  Gets information about an existing attestation AttestationProvider.
---

# Data Source: azurerm_attestation_attestation_provider

Use this data source to access information about an existing attestation AttestationProvider.

## Example Usage

```hcl
data "azurerm_attestation_attestation_provider" "example" {
  name = "example-attestationprovider"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_attestation_attestation_provider.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this attestation AttestationProvider.

* `resource_group_name` - (Required) The name of the Resource Group where the attestation AttestationProvider exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the attestation AttestationProvider.

* `location` - The Azure Region where the attestation AttestationProvider exists.

* `tags` - A mapping of tags assigned to the attestation AttestationProvider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the attestation AttestationProvider.
