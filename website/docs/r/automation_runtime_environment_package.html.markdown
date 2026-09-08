---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_runtime_environment_package"
description: |-
  Manages a Package within an Automation Runtime Environment.
---

# azurerm_automation_runtime_environment_package

Manages a Package within an Automation Runtime Environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "westeurope"
}

resource "azurerm_automation_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runtime_environment" "example" {
  name                  = "example-runtime-env"
  automation_account_id = azurerm_automation_account.example.id
  runtime_language      = "PowerShell"
  runtime_version       = "7.2"
  location              = azurerm_resource_group.example.location
}

resource "azurerm_automation_runtime_environment_package" "example" {
  name                              = "example-package"
  automation_runtime_environment_id = azurerm_automation_runtime_environment.example.id
  content_uri                       = "https://www.powershellgallery.com/api/v2/package/example-package/1.0.0"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the package. Changing this forces a new resource to be created.

* `automation_runtime_environment_id` - (Required) The ID of the Automation Runtime Environment in which to create this package. Changing this forces a new resource to be created.

* `content_uri` - (Required) The HTTPS URI of the package content. Changing this forces a new resource to be created.

---

* `content_version` - (Optional) The version of the package content. Changing this forces a new resource to be created.

~> **Note:** The `content_version` must be a version string with 2 to 4 segments (e.g. `1.0`, `1.0.0`, or `1.0.0.0`).

* `hash_algorithm` - (Optional) The hash algorithm used to hash the content. Changing this forces a new resource to be created.

~> **Note:** The argument `hash_algorithm` is required when `hash_value` is specified.

* `hash_value` - (Optional) The hash value of the content. Changing this forces a new resource to be created.

~> **Note:** The argument `hash_value` is required when `hash_algorithm` is specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Automation Runtime Environment Package.

* `default` - Whether this is a default package.

* `size_in_bytes` - The size of the package in bytes.

* `version` - The version of the package as reported by the platform.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Runtime Environment Package.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Runtime Environment Package.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Runtime Environment Package.

## Import

An Automation Runtime Environment Package can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_runtime_environment_package.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Automation/automationAccounts/automationAccount1/runtimeEnvironments/runtimeEnvironment1/packages/package1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Automation` - 2024-10-23
