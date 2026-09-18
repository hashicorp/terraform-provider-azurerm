---
subcategory: "Trusted Signing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_trusted_signing_certificate_profile"
description: |-
  Lists Trusted Signing Certificate Profiles (Artifact Signing Certificate Profiles).
---

# List resource: azurerm_trusted_signing_certificate_profile

Lists Trusted Signing Certificate Profiles (Artifact Signing Certificate Profiles) in a Trusted Signing Account.

## Example Usage

```hcl
list "azurerm_trusted_signing_certificate_profile" "example" {
  provider         = azurerm
  include_resource = true
  config {
    trusted_signing_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.CodeSigning/codeSigningAccounts/example-account"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `trusted_signing_account_id` - (Required) The ID of the Trusted Signing Account to query.
