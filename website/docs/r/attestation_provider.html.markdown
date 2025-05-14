---
subcategory: "Attestation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_attestation_provider"
description: |-
  Manages an Attestation Provider.
---

# azurerm_attestation_provider

Manages an Attestation Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_attestation_provider" "example" {
  name                = "exampleprovider"
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

* `policy_signing_certificate_data` - (Optional) A valid X.509 certificate (Section 4 of [RFC4648](https://tools.ietf.org/html/rfc4648)). Changing this forces a new resource to be created.

-> **Note:** If the `policy_signing_certificate_data` argument contains more than one valid X.509 certificate only the first certificate will be used.

* `open_enclave_policy_base64` - (Optional) Specifies the base64 URI Encoded RFC 7519 JWT that should be used for the Attestation Policy.

* `sgx_enclave_policy_base64` - (Optional) Specifies the base64 URI Encoded RFC 7519 JWT that should be used for the Attestation Policy.

* `tpm_policy_base64` - (Optional) Specifies the base64 URI Encoded RFC 7519 JWT that should be used for the Attestation Policy.

* `sev_snp_policy_base64` - (Optional) Specifies the base64 URI Encoded RFC 7519 JWT that should be used for the Attestation Policy.

-> **Note:** [More information on the JWT Policies can be found in this article on `learn.microsoft.com`](https://learn.microsoft.com/azure/attestation/author-sign-policy).

* `tags` - (Optional) A mapping of tags which should be assigned to the Attestation Provider.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Attestation Provider.

* `attestation_uri` - The URI of the Attestation Service.

* `trust_model` - Trust model used for the Attestation Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Attestation Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the Attestation Provider.
* `update` - (Defaults to 30 minutes) Used when updating the Attestation Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the Attestation Provider.

## Import

Attestation Providers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_attestation_provider.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Attestation/attestationProviders/provider1
```
