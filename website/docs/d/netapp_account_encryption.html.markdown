---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_netapp_account_encryption"
description: |-
  Gets information about an existing NetApp Account Encryption Resource.
---

# Data Source: azurerm_netapp_account_encryption

Use this data source to access information about an existing NetApp Account Encryption Resource.

## Example Usage

```hcl
data "azurerm_netapp_account_encryption" "example" {
  netapp_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1"
}

output "id" {
  value = data.azurerm_netapp_account_encryption.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `netapp_account_id` - (Required) The ID of the NetApp account where customer managed keys-based encryption is enabled.

---

* `encryption_key` - The key vault encryption key.

* `system_assigned_identity_principal_id` - The ID of the System Assigned Manged Identity.

* `user_assigned_identity_id` - The ID of the User Assigned Managed Identity.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Account Encryption Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Account Encryption Resource.
