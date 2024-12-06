---
subcategory: "Trusted Signing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_trusted_signing_certificate_profile"
description: |-
  Manages a Trusted Signing Certificate Profile.
---

# azurerm_trusted_signing_certificate_profile 

Manages a Trusted Signing Certificate Profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_trusted_signing_account" "example" {
  name                = "example-account"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}


resource "azurerm_trusted_signing_certificate_profile" "example" {
  name                       = "example-ccp"
  trusted_signing_account_id = azurerm_trusted_signing_account.example.id
  identity_validation_id     = "00000000-1111-2222-3333-444444444444"
  include_city               = false
  include_country            = false
  include_postal_code        = false
  include_state              = false
  include_street_address     = false
  profile_type               = "PublicTrust"

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Trusted Signing Certificate Profile. Changing this forces a new Trusted Signing Certificate Profile to be created.

* `trusted_signing_account_id ` - (Required) Specifies the ID of the Trusted Signing Account. Changing this forces a new Trusted Signing Certificate Profile to be created.

* `identity_validation_id` - (Required) Identity validation id used for the certificate subject name.

* `profile_type` - (Required) Profile type of the certificate.

* `include_city` - (Optional) Whether to include L in the certificate subject name. Applicable only for private trust, private trust ci profile types.

* `include_country` - (Optional) Whether to include C in the certificate subject name. Applicable only for private trust, private trust ci profile types.

* `include_postal_code` - (Optional) Whether to include PC in the certificate subject name.

* `include_state` - (Optional) Whether to include S in the certificate subject name. Applicable only for private trust, private trust ci profile types.

* `include_street_address` - (Optional) Whether to include STREET in the certificate subject name.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Trusted Signing Certificate Profile.

* `certificates` - A `certificates` block as defined below.

* `status` - Status of the certificate profile.

---

A `certificates` block exports the following:

* `created_date` - Certificate created date.

* `enhanced_key_usage` - Enhanced key usage of the certificate.

* `expiry_date` - Certificate expiry date.

* `revocation` - A `revocation` block as defined below.

* `serial_number` - Serial number of the certificate.

* `status` - Status of the certificate.

* `subject_name` - Subject name of the certificate.

* `thumbprint` - Thumbprint of the certificate.

---

A `revocation` block exports the following:

* `effective_at` - The timestamp when the revocation is effective.

* `failure_reason` - Reason for the revocation failure.

* `reason` - Reason for revocation.

* `remarks` - Remarks for the revocation.

* `requested_at` - The timestamp when the revocation is requested.

* `status` - Status of the revocation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Trusted Signing Certificate Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Trusted Signing Certificate Profile.

* `delete` - (Defaults to 30 minutes) Used when deleting the Trusted Signing Certificate Profile.

## Import

Trusted Signing Certificate Profile can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_trusted_signing_certificate_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CodeSigning/codeSigningAccounts/account1/certificateProfiles/profile1
```
