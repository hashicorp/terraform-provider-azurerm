---
subcategory: "Attestation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_attestation"
description: |-
  Manages a Attestation Provider.
---

# azurerm_attestation

Manages a Attestation Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "UK South"
}

resource "azurerm_attestation_provider" "example" {
  name                = "example-attestationprovider"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  policy_signing_certificate_data = file("./example/cert.pem")
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Attestation Provider. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the attestation provider should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Attestation Provider should exist. Changing this forces a new resource to be created.

-> **NOTE:** Currently only supported in the `East US 2`, `West Central US`, and `UK South` regions.

* `policy_signing_certificate_data` - (Optional) A valid X.509 certificate (Section 4 of [RFC4648](https://tools.ietf.org/html/rfc4648)). Changing this forces a new resource to be created.

-> **NOTE:** If the `policy_signing_certificate_data` argument contains more than one valid X.509 certificate only the first certificate will be used.

* `tags` - (Optional) A mapping of tags which should be assigned to the Attestation Provider.

## Attributes Reference

The following Attributes are exported: 

* `id` - The ID of the Attestation Provider.

* `attestation_uri` - The URI of the Attestation Service.

* `trust_model` - Trust model used for the Attestation Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Attestation Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the Attestation Provider.
* `update` - (Defaults to 30 minutes) Used when updating the Attestation Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the Attestation Provider.

## Import

Attestation Providers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_attestation_provider.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Attestation/attestationProviders/provider1
```
