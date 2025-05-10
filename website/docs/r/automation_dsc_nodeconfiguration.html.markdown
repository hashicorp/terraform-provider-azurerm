---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_dsc_nodeconfiguration"
description: |-
  Manages a Automation DSC Node Configuration.
---

# azurerm_automation_dsc_nodeconfiguration

Manages a Automation DSC Node Configuration.

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

resource "azurerm_automation_dsc_configuration" "example" {
  name                    = "test"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  location                = azurerm_resource_group.example.location
  content_embedded        = "configuration test {}"
}

resource "azurerm_automation_dsc_nodeconfiguration" "example" {
  name                    = "test.localhost"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name
  depends_on              = [azurerm_automation_dsc_configuration.example]

  content_embedded = <<mofcontent
instance of MSFT_FileDirectoryConfiguration as $MSFT_FileDirectoryConfiguration1ref
{
  ResourceID = "[File]bla";
  Ensure = "Present";
  Contents = "bogus Content";
  DestinationPath = "c:\\bogus.txt";
  ModuleName = "PSDesiredStateConfiguration";
  SourceInfo = "::3::9::file";
  ModuleVersion = "1.0";
  ConfigurationName = "bla";
};
instance of OMI_ConfigurationDocument
{
  Version="2.0.0";
  MinimumCompatibleVersion = "1.0.0";
  CompatibleVersionAdditionalProperties= {"Omi_BaseResource:ConfigurationName"};
  Author="bogusAuthor";
  GenerationDate="06/15/2018 14:06:24";
  GenerationHost="bogusComputer";
  Name="test";
};
mofcontent

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the DSC Node Configuration. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the DSC Node Configuration is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the DSC Node Configuration is created. Changing this forces a new resource to be created.

* `content_embedded` - (Required) The PowerShell DSC Node Configuration (mof content).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The DSC Node Configuration ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation DSC Node Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation DSC Node Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Automation DSC Node Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation DSC Node Configuration.

## Import

Automation DSC Node Configuration's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_dsc_nodeconfiguration.configuration1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/nodeConfigurations/configuration1
```
