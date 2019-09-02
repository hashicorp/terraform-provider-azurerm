---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_source_control_token"
sidebar_current: "docs-azurerm-resource-app-service-source-control-token"
description: |-
  Manages an App Service source control token.

---

# azurerm_app_service_source_control_token

Manages an App Service source control token.

## Example Usage

```hcl
resource "azurerm_app_service_source_control_token" "test" {
  type  = "GitHub"
  token = "7e57735e77e577e57"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The source control type. Possible values are `BitBucket`, `Dropbox`, `GitHub` and `OneDrive`.

* `token` - (Required) The OAuth access token.

## Import

App Service source control tokens can be imported using the `type`, e.g.

```shell
terraform import azurerm_app_service_source_control_token.test GitHub
```
