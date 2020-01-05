---
subcategory: "Data Lake"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_store"
sidebar_current: "docs-azurerm-resource-data-lake-store-x"
description: |-
  Manages an Azure Data Lake Store.
---

# azurerm_data_lake_store

Manages an Azure Data Lake Store.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_data_lake_store" "example" {
  name                = "consumptiondatalake"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  encryption_state    = "Enabled"
  encryption_type     = "ServiceManaged"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Lake Store. Changing this forces a new resource to be created. Has to be between 3 to 24 characters.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Lake Store.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tier` - (Optional) The monthly commitment tier for Data Lake Store. Accepted values are `Consumption`, `Commitment_1TB`, `Commitment_10TB`, `Commitment_100TB`, `Commitment_500TB`, `Commitment_1PB` or `Commitment_5PB`.

* `encryption_state` - (Optional) Is Encryption enabled on this Data Lake Store Account? Possible values are `Enabled` or `Disabled`. Defaults to `Enabled`.

* `encryption_type` - (Optional) The Encryption Type used for this Data Lake Store Account. Currently can be set to `ServiceManaged` when `encryption_state` is `Enabled` - and must be a blank string when it's Disabled.

-> **NOTE:** Support for User Managed encryption will be supported in the future once a bug in the API is fixed.

* `firewall_allow_azure_ips` - are Azure Service IP's allowed through the firewall? Possible values are `Enabled` and `Disabled`. Defaults to `Enabled.`

* `firewall_state` - the state of the Firewall. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled.`

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Date Lake Store ID.

* `endpoint` - The Endpoint for the Data Lake Store.

## Import

Date Lake Store can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_store.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DataLakeStore/accounts/mydatalakeaccount
```
