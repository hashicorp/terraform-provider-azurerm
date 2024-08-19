---
subcategory: "Elastic SAN"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_elastic_san_volume_group"
description: |-
  Gets information about an existing Elastic SAN Volume Group.
---

# Data Source: azurerm_elastic_san_volume_group

Use this data source to access information about an existing Elastic SAN Volume Group.

## Example Usage

```hcl
data "azurerm_elastic_san" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

data "azurerm_elastic_san_volume_group" "example" {
  name           = "existing"
  elastic_san_id = data.azurerm_elastic_san.example.id
}

output "id" {
  value = data.azurerm_elastic_san_volume_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the Elastic SAN Volume Group.

* `elastic_san_id` - The Elastic SAN ID within which the Elastic SAN Volume Group exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Elastic SAN.

* `encryption_type` - The type of the key used to encrypt the data of the disk.

* `encryption` - An `encryption` block as defined below.

* `identity` - An `identity` block as defined below.

* `network_rule` - One or more `network_rule` blocks as defined below.

* `protocol_type` - The type of the storage target.

---

An `encryption` block exports the following arguments:

* `key_vault_key_id` - The Key Vault Key URI for Customer Managed Key encryption, which can be either a full URI or a versionless URI.

* `user_assigned_identity_id` - The ID of the User Assigned Identity used by this Elastic SAN Volume Group.

* `current_versioned_key_expiration_timestamp` - The timestamp of the expiration time for the current version of the Customer Managed Key.

* `current_versioned_key_id` - The ID of the current versioned Key Vault Key in use.

* `last_key_rotation_timestamp` - The timestamp of the last rotation of the Key Vault Key.

---

An `identity` block exports the following arguments:

* `type` - The type of Managed Identity assigned to this Elastic SAN Volume Group.

* `identity_ids` - A list of the User Assigned Identity IDs assigned to this Elastic SAN Volume Group.

* `principal_id` - The Principal ID associated with the Managed Service Identity assigned to this Elastic SAN Volume Group.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity assigned to this Elastic SAN Volume Group.

---

A `network_rule` block exports the following arguments:

* `subnet_id` - The ID of the Subnet from which access to this Elastic SAN Volume Group is allowed.

* `action` - The action to take when an access attempt to this Elastic SAN Volume Group from this Subnet is made.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic SAN Volume Group.
