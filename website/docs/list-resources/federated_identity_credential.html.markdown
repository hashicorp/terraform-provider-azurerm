---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_federated_identity_credential"
description: |-
  Lists Federated Identity Credential resources.
---

# List resource: azurerm_federated_identity_credential

Lists Federated Identity Credential resources.

## Example Usage

### List all Federated Identity Credentials for a User Assigned Identity

```hcl
list "azurerm_federated_identity_credential" "example" {
  provider = azurerm
  config {
    user_assigned_identity_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ManagedIdentity/userAssignedIdentities/identity1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `user_assigned_identity_id` - (Required) The ID of the User Assigned Identity whose Federated Identity Credentials should be listed.
