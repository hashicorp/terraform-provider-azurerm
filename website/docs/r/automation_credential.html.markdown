---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_credential"
sidebar_current: "docs-azurerm-resource-automation-credential"
description: |-
  Manages a Automation Credential.
---

# azurerm_automation_credential

Manages a Automation Credential.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
 name = "resourceGroup1"
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

resource "azurerm_automation_credential" "example" {
  name                = "credential1"
  resource_group_name = "${azurerm_resource_group.example.name}"
  account_name        = "${azurerm_automation_account.example.name}"
  username           = "example_user"
  password            = "example_pwd"
  description         = "This is an example credential"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Credential. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Credential is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the automation account in which the Credential is created. Changing this forces a new resource to be created.

* `username` - (Required) The username associated with this Automation Credential.

* `password` - (Required) The password associated with this Automation Credential.

* `description` -  (Optional) The description associated with this Automation Credential.

## Attributes Reference

The following attributes are exported:

* `id` - The Automation Credential ID.

## Import

Automation Credentials can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_credential.credential1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/credentials/credential1
```
