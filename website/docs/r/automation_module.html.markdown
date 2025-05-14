---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_module"
description: |-
  Manages a Automation Module.
---

# azurerm_automation_module

Manages a Automation Module.

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
  sku_name            = "Basic"
}

resource "azurerm_automation_module" "example" {
  name                    = "xActiveDirectory"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Module. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Module is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Module is created. Changing this forces a new resource to be created.

* `module_link` - (Required) A `module_link` block as defined below.

---

The `module_link` block supports the following:

* `uri` - (Required) The URI of the module content (zip or nupkg).

* `hash` - (Optional) A `hash` block as defined below.

---

The `hash` block supports the following:

* `algorithm` - (Required) Specifies the algorithm used for the hash content.

* `value` - (Required) The hash value of the content.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Automation Module ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Module.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Module.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Module.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Module.

## Import

Automation Modules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_module.module1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/modules/module1
```
