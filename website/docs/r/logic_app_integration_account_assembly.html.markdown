---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_assembly"
description: |-
  Manages a Logic App Integration Account Assembly.
---

# azurerm_logic_app_integration_account_assembly

Manages a Logic App Integration Account Assembly.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "example" {
  name                = "example-ia"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_integration_account_assembly" "example" {
  name                     = "example-assembly"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name
  assembly_name            = "TestAssembly"
  content                  = filebase64("testdata/log4net.dll")
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Assembly Artifact. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Assembly Artifact should exist. Changing this forces a new resource to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new resource to be created.

* `assembly_name` - (Required) The name of the Logic App Integration Account Assembly.

* `assembly_version` - (Optional) The version of the Logic App Integration Account Assembly. Defaults to `0.0.0.0`.

* `content` - (Optional) The content of the Logic App Integration Account Assembly.

* `content_link_uri` - (Optional) The content link URI of the Logic App Integration Account Assembly.

* `metadata` - (Optional) The metadata of the Logic App Integration Account Assembly.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Assembly.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Assembly.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Assembly.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Assembly.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Assembly.

## Import

Logic App Integration Account Assemblies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_assembly.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/assemblies/assembly1
```
