---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_source_control"
description: |-
  Manages an Automation Source Control.
---

# azurerm_automation_source_control

Manages an Automation Source Control.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_source_control" "example" {
  name                  = "example"
  automation_account_id = azurerm_automation_account.example.id
  folder_path           = "runbook"

  security {
    token      = "ghp_xxx"
    token_type = "PersonalAccessToken"
  }
  repository_url      = "https://github.com/foo/bat.git"
  source_control_type = "GitHub"
  branch              = "main"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Automation Source Control. Changing this forces a new Automation Source Control to be created.

* `automation_account_id` - (Required) The ID of Automation Account to manage this Source Control. Changing this forces a new Automation Source Control to be created.

* `folder_path` - (Required) The folder path of the source control. This Path must be relative.

* `repository_url` - (Required) The Repository URL of the source control.

* `security` - (Required) A `security` block as defined below.

* `source_control_type` - (Required) The source type of Source Control, possible vaules are `VsoGit`, `VsoTfvc` and `GitHub`, and the value is case sensitive.

---

* `automatic_sync` - (Optional) Whether auto async the Source Control.

* `branch` - (Optional) Specify the repo branch of the Source Control. Empty value is valid only for `VsoTfvc`.

* `description` - (Optional) A short description of the Source Control.

* `publish_runbook_enabled` - (Optional) Whether auto publish the Source Control. Defaults to `true`.

---

A `security` block supports the following:

* `token` - (Required) The access token of specified repo.

* `token_type` - (Required) Specify the token type, possible values are `PersonalAccessToken` and `Oauth`.

* `refresh_token` - (Optional) The refresh token of specified rpeo.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Source Control.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation.
* `update` - (Defaults to 10 minutes) Used when updating the Automation.
* `delete` - (Defaults to 10 minutes) Used when deleting the Automation.

## Import

Automations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_source_control.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/sourceControls/sc1
```
