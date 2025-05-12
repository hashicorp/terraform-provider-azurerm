---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_powershell72_module"
description: |-
  Manages a Automation Powershell 7.2 Module.
---

# azurerm_automation_powershell72_module

Manages a Automation Powershell 7.2 Module.

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

resource "azurerm_automation_powershell72_module" "example" {
  name                  = "xActiveDirectory"
  automation_account_id = azurerm_automation_account.example.id

  module_link {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Module. Changing this forces a new resource to be created.

* `automation_account_id` - (Required) The ID of Automation Account to manage this Watcher. Changing this forces a new Watcher to be created.

* `module_link` - (Required) A `module_link` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

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

* `create` - (Defaults to 30 minutes) Used when creating the Automation Powershell 7.2 Module.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Powershell 7.2 Module.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Powershell 7.2 Module.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Powershell 7.2 Module.

## Import

Automation Modules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_powershell72_module.module1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/powerShell72Modules/module1
```
