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
  name                = "example-attestationprovider"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_attestation.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Attestation Provider.

* `resource_group_name` - (Required) The name of the Resource Group where the Attestation Provider exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Attestation Provider.

* `location` - The Azure Region where the Attestation Provider exists.

* `attestation_uri` - The (Endpoint|URI) of the Attestation Service.

* `trust_model` - Trust model used for the Attestation Service.

* `tags` - A mapping of tags assigned to the Attestation Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Attestation Provider.
