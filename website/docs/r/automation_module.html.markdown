---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_module"
sidebar_current: "docs-azurerm-resource-automation-module"
description: |-
  Manages a Automation Module.
---

# azurerm_automation_module

Manages a Automation Module.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_module" "example" {
  name                    = "xActiveDirectory"
  resource_group_name     = "${azurerm_resource_group.example.name}"
  automation_account_name = "${azurerm_automation_account.example.name}"

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

* `module_link` - (Required) The published Module link.

`module_link` supports the following:

* `uri` - (Required) The uri of the module content (zip or nupkg).

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Module ID.
