---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_python2_package"
description: |-
  Manages a Automation.
---

# azurerm_automation_python2_package

Manages a Automation Python 2 Package.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Basic"
}

resource "azurerm_automation_python2_package" "example" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  name                    = "acctest-py-example"
  tags = {
    foo = "bar"
  }
  content {
    uri            = "https://files.pythonhosted.org/packages/3e/93/02056aca45162f9fc275d1eaad12a2a07ef92375afb48eabddc4134b8315/azure_graphrbac-0.61.1-py2.py3-none-any.whl"
    hash_algorithm = "sha256"
    hash_value     = "7b4e0f05676acc912f2b33c71c328d9fb2e4dc8e70ebadc9d3de8ab08bf0b175"
    version        = "1.0.0.0"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Automation Python 2 Package. Changing this forces a new Automation to be created.
 
* `resource_group_name` - (Required) The name of the resource group in which the Runbook is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Runbook is created. Changing this forces a new resource to be created.
 
* `content` - (Required) One or more `content` blocks as defined below.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Automation.

---

A `content` block supports the following:

* `uri` - (Required) The URI of the runbook content.

* `hash_algorithm` - (Optional) The algorithm used to hash the content.

* `hash_value` - (Optional) The expected hash value of the content.

* `version` - (Optional) The Version of the content.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Automation Python2.

* `activity_count` - The activity count of the module.

* `creation_time` - The creation time of the module.

* `error_code` - The error code of the module.

* `error_meesage` - The error message of the module.

* `is_composite` - Whethe the module is composite or not.

* `is_global` - Whether the module's isGlobal flag set.

* `last_modified_time` - The last modified time of the module.

* `location` - The Azure Region where the Automation Python 2 Package exists.

* `size_in_bytes` - The Size in byte of the module.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `update` - (Defaults to 10 minutes) Used when updating the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_python2_package.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/python2Packages/pkg1
```
