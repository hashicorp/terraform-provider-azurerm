---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store"
sidebar_current: "docs-azurerm-resource-data-lake-store-x"
description: |-
  Manage an Azure Data Lake Store.
---

# azurerm_data_lake_store

Manage an Azure Data Lake Store.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "northeurope"
}

resource "azurerm_data_lake_store" "example" {
  name                = "consumptiondatalake"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  
  encrytpion {
    type = "UserManaged"
    key_vault_id = "${azurerm_key_vault.example.id}"
    key_name     = "${azurerm_key_vault_key.example.name}"
    key_version  = "${azurerm_key_vault_key.example.version}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Lake Store. Changing this forces a new resource to be created. Has to be between 3 to 24 characters.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Lake Store.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tier` - (Optional) The monthly commitment tier for Data Lake Store. Accepted values are `Consumption`, `Commitment_1TB`, `Commitment_10TB`, `Commitment_100TB`, `Commitment_500TB`, `Commitment_1PB` or `Commitment_5PB`.

* `encryption` - (Optional/Computed) A block detailing the current encryption settings.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`encryption` block supports the following:

* `enabled` - (Optional) Sets if encryption is enabled or disabled for this store. Defaults to `false` 

* `type` - (Optional) This property determins the source of the encryption keys used, can be one of `ServiceManaged` or `UserManaged`. Defaults to `ServiceManaged`.
    
* `key_vault_id` - (Optional) The id of a key vault to get an encryption key from. This must be specified when the encryption type is `UserManaged`

* `key_name` - (Optional) The name of a key in the key vault to use for encryption. This must be specified when the encryption type is `UserManaged`

* `key_version` - (Optional) The version of the key to use. This must be specified when the encryption type is `UserManaged`
    
## Attributes Reference

The following attributes are exported:

* `id` - The Date Lake Store ID.

## Import

Date Lake Store can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_store.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DataLakeStore/accounts/mydatalakeaccount
```
