---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine_group"
description: |-
  Manages a Microsoft SQL Virtual Machine Group.
---

# azurerm_mssql_virtual_machine_group

Manages a Microsoft SQL Virtual Machine Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mssql_virtual_machine_group" "example" {
  name                = "examplegroup"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sql_image_offer = "SQL2017-WS2016"
  sql_image_sku   = "Developer"

  wsfc_domain_profile {
    fqdn                = "testdomain.com"
    cluster_subnet_type = "SingleSubnet"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for the Microsoft SQL Virtual Machine Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Microsoft SQL Virtual Machine Group should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Microsoft SQL Virtual Machine Group should exist. Changing this forces a new resource to be created.

* `sql_image_offer` - (Required) The offer type of the marketplace image cluster to be used by the SQL Virtual Machine Group. Changing this forces a new resource to be created.

* `sql_image_sku` - (Required) The sku type of the marketplace image cluster to be used by the SQL Virtual Machine Group. Possible values are `Developer` and `Enterprise`.

* `wsfc_domain_profile` - (Required) A `wsfc_domain_profile` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Microsoft SQL Virtual Machine Group.

---

A `wsfc_domain_profile` block supports the following:

* `cluster_subnet_type` - (Required) The subnet type of the SQL Virtual Machine cluster. Possible values are `MultiSubnet` and `SingleSubnet`. Changing this forces a new resource to be created.

* `fqdn` - (Required) The fully qualified name of the domain. Changing this forces a new resource to be created.

* `cluster_bootstrap_account_name` - (Optional) The account name used for creating cluster. Changing this forces a new resource to be created.

* `cluster_operator_account_name` - (Optional) The account name used for operating cluster. Changing this forces a new resource to be created.

* `organizational_unit_path` - (Optional) The organizational Unit path in which the nodes and cluster will be present. Changing this forces a new resource to be created.

* `sql_service_account_name` - (Optional) The account name under which SQL service will run on all participating SQL virtual machines in the cluster. Changing this forces a new resource to be created.

* `storage_account_primary_key` - (Optional) The primary key of the Storage Account.

* `storage_account_url` - (Optional) The SAS URL to the Storage Container of the witness storage account. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Microsoft SQL Virtual Machine Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Microsoft SQL Virtual Machine Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Virtual Machine Group.
* `update` - (Defaults to 30 minutes) Used when updating the Microsoft SQL Virtual Machine Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Microsoft SQL Virtual Machine Group.

## Import

Microsoft SQL Virtual Machine Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_virtual_machine_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/vmgroup1
```



