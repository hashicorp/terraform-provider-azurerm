---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_encryption_set"
sidebar_current: "docs-azurerm-datasource-disk-encryption-set"
description: |-
  Gets information about an existing Disk Encryption Set
---

# Data Source: azurerm_disk_encryption_set

Use this data source to access information about an existing Disk Encryption Set.



## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Disk Encryption Set exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Disk Encryption Set exists.


## Attributes Reference

The following attributes are exported:

* `id` - Resource Id

* `location` - Resource location

* `active_key` - One `active_key` block defined below.

* `identity` - A `identity` block defined below.

* `previous_keys` - One or more `previous_key` block defined below.

* `tags` - Resource tags


---

The `active_key` block contains the following:

* `source_vault_id` - The resource id of the KeyVault containing the key or secret which the Disk Encryption Set is using.

* `key_url` - The URL pointing to a key or secret in KeyVault.

---

The `identity` block contains the following:

* `type` - The type of Managed Service Identity used by this Disk Encryption Set. Only SystemAssigned is supported.

* `principal_id` - The object ID of the Managed Service Identity created by Azure.

* `tenant_id` - The tenant ID of the Managed Service Identity created by Azure.

---

The `previous_key` block contains the following:

* `source_vault_id` - The resource id of the KeyVault containing the key or secret which the Disk Encryption Set is using.

* `key_url` - The URL pointing to a key or secret in KeyVault.
