---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_dsc_configuration"
sidebar_current: "docs-azurerm-resource-automation-dsc-configuration"
description: |-
  Manages a Automation DSC Configuration.
---

# azurerm_automation_dsc_configuration

Manages a Automation DSC Configuration.

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

resource "azurerm_automation_dsc_configuration" "example" {
  name                    = "test"
  resource_group_name     = "${azurerm_resource_group.example.name}"
  automation_account_name = "${azurerm_automation_account.example.name}"
  location                = "${azurerm_resource_group.example.location}"
  content_embedded        = "configuration test {}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the DSC Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the DSC Configuration is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the DSC Configuration is created. Changing this forces a new resource to be created.

* `content_embedded` - (Required) The PowerShell DSC Configuration script.

* `location` - (Required) Must be the same location as the Automation Account.

* `log_verbose` - (Optional) Verbose log option.

* `description` - (Optional) Description to go with DSC Configuration.

## Attributes Reference

The following attributes are exported:

* `id` - The DSC Configuration ID.
