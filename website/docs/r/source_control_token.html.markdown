---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_source_control_token"
description: |-
  Manages an App Service GitHub Token.
---

# azurerm_source_control_token

Manages an App Service Source Control Token.

~> **Note:** This resource can only manage the token for the user currently running Terraform. Managing tokens for another user is not supported by the service.

## Example Usage

```hcl
resource "azurerm_source_control_token" "example" {
  type  = "GitHub"
  token = "ghp_sometokenvaluesecretsauce"
}
```

## Arguments Reference

The following arguments are supported:

* `type` - (Required) The Token type. Possible values include `Bitbucket`, `Dropbox`, `Github`, and `OneDrive`.

* `token` - (Required) The Access Token.

* `token_secret` - (Optional) The Access Token Secret.

~> **Note:** The token used for deploying App Service needs the following permissions: `repo` and `workflow`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Source GitHub Token.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the App Service Source GitHub Token.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Source GitHub Token.
* `update` - (Defaults to 5 minutes) Used when updating the App Service Source GitHub Token.
* `delete` - (Defaults to 5 minutes) Used when deleting the App Service Source GitHub Token.

## Import

App Service Source GitHub Tokens can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_source_control_token.example /providers/Microsoft.Web/sourceControls/GitHub
```
