---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_python3_package"
description: |-
  Manages a Automation Python3 Package.
---

# azurerm_automation_python3_package

Manages a Automation Python3 Package.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "%[2]s"
}

resource "azurerm_automation_account" "example" {
  name                = "accexample"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_python3_package" "example" {
  name                    = "example"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  content_uri             = "https://pypi.org/packages/source/r/requests/requests-2.31.0.tar.gz"
  content_version         = "2.31.0"
  hash_algorithm          = "sha256"
  hash_value              = "942c5a758f98d790eaed1a29cb6eefc7ffb0d1cf7af05c3d2791656dbd6ad1e1"
  tags = {
    key = "foo"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Python3 Package is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Python3 Package is created. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Automation Python3 Package. Changing this forces a new Automation Python3 Package to be created.

* `content_uri` - (Required) The URL of the python package. Changing this forces a new Automation Python3 Package to be created.


---

* `content_version` - (Optional) Specify the version of the python3 package. The value should meet the system.version class format like `1.1.1`. Changing this forces a new Automation Python3 Package to be created.

* `hash_algorithm` - (Optional) Specify the hash algorithm used to hash the content of the python3 package. Changing this forces a new Automation Python3 Package to be created.

* `hash_value` - (Optional) Specity the hash value of the content. Changing this forces a new Automation Python3 Package to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Automation Python3 Package.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Automation Python3 Package.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Python3 Package.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Python3 Package.
* `update` - (Defaults to 10 minutes) Used when updating the Automation Python3 Package.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation Python3 Package.

## Import

Automation Python3 Packages can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_python3_package.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/python3Packages/pkg
```
