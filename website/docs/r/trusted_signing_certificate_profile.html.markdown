---
subcategory: "Trusted Signing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_trusted_signing_certificate_profile"
description: |-
  Manages a Trusted Signing Certificate Profile (Artifact Signing Certificate Profile).
---

# azurerm_trusted_signing_certificate_profile

Manages a Trusted Signing Certificate Profile (Artifact Signing Certificate Profile).

~> **Note:** Trusted Signing has been rebranded to [Artifact Signing](https://learn.microsoft.com/azure/artifact-signing/overview).

~> **Note:** Before creating a certificate profile, complete an [identity validation](https://learn.microsoft.com/azure/artifact-signing/quickstart#create-an-identity-validation-request) in the Azure Portal. Identity validations can be shared across Trusted Signing Accounts in the same subscription. Private identity validations support `PrivateTrust` and `PrivateTrustCIPolicy`; public identity validations support `PublicTrust`, `PublicTrustTest`, and `VBSEnclave`.

## Example Usage

```hcl
data "azurerm_trusted_signing_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resources"
}

resource "azurerm_trusted_signing_certificate_profile" "example" {
  name                       = "example-profile"
  trusted_signing_account_id = data.azurerm_trusted_signing_account.example.id
  identity_validation_id     = "00000000-0000-0000-0000-000000000000"
  profile_type               = "PrivateTrust"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Trusted Signing Certificate Profile. Must be between 5 and 100 characters, begin with a letter, end with a letter or digit, and contain only alphanumeric characters and non-consecutive hyphens. Changing this forces a new resource to be created.

* `trusted_signing_account_id` - (Required) The ID of the Trusted Signing Account in which to create this certificate profile. Changing this forces a new resource to be created.

* `identity_validation_id` - (Required) The ID of a completed identity validation in the same subscription as the Trusted Signing Account, used for the certificate subject name. Changing this forces a new resource to be created.

* `profile_type` - (Required) The type of the certificate profile. Possible values are `PrivateTrust`, `PrivateTrustCIPolicy`, `PublicTrust`, `PublicTrustTest`, and `VBSEnclave`. Changing this forces a new resource to be created.

* `include_city` - (Optional) Whether to include the city (`L`) in the certificate subject name. Applicable only to `PrivateTrust` and `PrivateTrustCIPolicy` profiles. Defaults to `false`.

* `include_country` - (Optional) Whether to include the country (`C`) in the certificate subject name. Applicable only to `PrivateTrust` and `PrivateTrustCIPolicy` profiles. Defaults to `false`.

* `include_postal_code` - (Optional) Whether to include the postal code (`PC`) in the certificate subject name. Defaults to `false`.

* `include_state` - (Optional) Whether to include the state (`S`) in the certificate subject name. Applicable only to `PrivateTrust` and `PrivateTrustCIPolicy` profiles. Defaults to `false`.

* `include_street_address` - (Optional) Whether to include the street address (`STREET`) in the certificate subject name. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Trusted Signing Certificate Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Trusted Signing Certificate Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Trusted Signing Certificate Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Trusted Signing Certificate Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Trusted Signing Certificate Profile.

## Import

Trusted Signing Certificate Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_trusted_signing_certificate_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.CodeSigning/codeSigningAccounts/example-account/certificateProfiles/example-profile
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CodeSigning` - 2025-10-13
